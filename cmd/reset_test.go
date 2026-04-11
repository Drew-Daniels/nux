package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Drew-Daniels/nux/internal/config"
)

func resetDeps(t *testing.T) (*deps, string) {
	t.Helper()
	d := testDeps(t)
	cfgDir := filepath.Dir(d.projectCfgDir)
	cfgPath := filepath.Join(cfgDir, "config.yaml")
	if err := os.WriteFile(cfgPath, config.ScaffoldGlobalConfig(), 0o644); err != nil {
		t.Fatal(err)
	}
	return d, cfgDir
}

func TestRunResetWith_RemovesConfig(t *testing.T) {
	d, cfgDir := resetDeps(t)
	resetForce = true
	resetProjects = false
	t.Cleanup(func() { resetForce = false })

	_ = d.store.Save("blog", &config.ProjectConfig{Windows: []config.Window{{Name: "main", Panes: []config.Pane{{Command: "vim"}}}}})

	if err := runResetWith(d); err != nil {
		t.Fatalf("runResetWith: %v", err)
	}

	cfgPath := filepath.Join(cfgDir, "config.yaml")
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		t.Error("expected config to be removed")
	}

	if _, _, err := d.store.Load("blog"); err != nil {
		t.Error("project configs should be kept by default")
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "Done.") {
		t.Errorf("expected Done in output, got %q", out)
	}
}

func TestRunResetWith_RemovesProjects(t *testing.T) {
	d, cfgDir := resetDeps(t)
	resetForce = true
	resetProjects = true
	t.Cleanup(func() { resetForce = false; resetProjects = false })

	_ = d.store.Save("blog", &config.ProjectConfig{Windows: []config.Window{{Name: "main", Panes: []config.Pane{{Command: "vim"}}}}})

	if err := runResetWith(d); err != nil {
		t.Fatalf("runResetWith: %v", err)
	}

	cfgPath := filepath.Join(cfgDir, "config.yaml")
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		t.Error("expected config to be removed")
	}
	if _, err := os.Stat(d.projectCfgDir); !os.IsNotExist(err) {
		t.Error("expected projects dir to be removed")
	}
}

func TestRunResetWith_PreviewShowsWillRemoveAndKeep(t *testing.T) {
	d, _ := resetDeps(t)
	resetForce = true
	resetProjects = false
	t.Cleanup(func() { resetForce = false })

	_ = d.store.Save("blog", &config.ProjectConfig{Windows: []config.Window{{Name: "main", Panes: []config.Pane{{Command: "vim"}}}}})

	if err := runResetWith(d); err != nil {
		t.Fatalf("runResetWith: %v", err)
	}

	out := stdoutStr(d)
	for _, want := range []string{
		"Will remove:",
		"config:",
		"Will keep:",
		"--projects",
		"Running tmux sessions are not affected",
		"Done.",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in output:\n%s", want, out)
		}
	}
}

func TestRunResetWith_PreviewWithProjects(t *testing.T) {
	d, _ := resetDeps(t)
	resetForce = true
	resetProjects = true
	t.Cleanup(func() { resetForce = false; resetProjects = false })

	_ = d.store.Save("blog", &config.ProjectConfig{Windows: []config.Window{{Name: "main", Panes: []config.Pane{{Command: "vim"}}}}})

	if err := runResetWith(d); err != nil {
		t.Fatalf("runResetWith: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "projects:") {
		t.Errorf("expected projects in Will remove section:\n%s", out)
	}
	if strings.Contains(out, "--projects") {
		t.Errorf("should not suggest --projects when already purging:\n%s", out)
	}
}

func TestRunResetWith_Cancelled(t *testing.T) {
	d, cfgDir := resetDeps(t)
	resetForce = false
	resetProjects = false
	d.confirm = func(string) (bool, error) { return false, nil }

	if err := runResetWith(d); err != nil {
		t.Fatalf("runResetWith: %v", err)
	}

	cfgPath := filepath.Join(cfgDir, "config.yaml")
	if _, err := os.Stat(cfgPath); err != nil {
		t.Error("config should still exist after cancel")
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "Cancelled") {
		t.Errorf("expected 'Cancelled' in output, got %q", out)
	}
}

func TestRunResetWith_ConfigMissing(t *testing.T) {
	d := testDeps(t)
	resetForce = true
	t.Cleanup(func() { resetForce = false })

	err := runResetWith(d)
	if err == nil {
		t.Fatal("expected error for missing config")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error = %q, expected 'not found'", err.Error())
	}
}
