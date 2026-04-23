package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/Drew-Daniels/nux/internal/config"
)

func TestRunCopyWith_Copies(t *testing.T) {
	d := testDeps(t)
	d.openEditor = func(string) error { return nil }
	_ = d.store.Save("blog", &config.ProjectConfig{Windows: []config.Window{{Name: "editor", Panes: []config.Pane{{Command: "vim"}}}}})

	if err := runCopyWith(d, []string{"blog", "blog2"}); err != nil {
		t.Fatalf("runCopyWith: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "Copied") {
		t.Errorf("expected 'Copied' in output, got %q", out)
	}

	cfg, _, err := d.store.Load("blog2")
	if err != nil {
		t.Fatalf("destination config should exist: %v", err)
	}
	if len(cfg.Windows) != 1 || cfg.Windows[0].Name != "editor" {
		t.Errorf("unexpected config: %+v", cfg)
	}
}

func TestRunCopyWith_SourceNotFound(t *testing.T) {
	d := testDeps(t)

	err := runCopyWith(d, []string{"nonexistent", "dest"})
	if err == nil {
		t.Fatal("expected error for missing source")
	}
	if !strings.Contains(err.Error(), "config not found") {
		t.Errorf("error = %q, expected 'config not found'", err.Error())
	}
}

func TestRunCopyWith_DestExists(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("blog", &config.ProjectConfig{Windows: []config.Window{{Name: "editor", Panes: []config.Pane{{Command: "vim"}}}}})
	_ = d.store.Save("blog2", &config.ProjectConfig{Windows: []config.Window{{Name: "main", Panes: []config.Pane{{Command: "bash"}}}}})

	err := runCopyWith(d, []string{"blog", "blog2"})
	if err == nil {
		t.Fatal("expected error for existing destination")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("error = %q, expected 'already exists'", err.Error())
	}
}

func TestRunCopyWith_DestExists_Force(t *testing.T) {
	d := testDeps(t)
	d.copyForce = true
	d.openEditor = func(string) error { return nil }
	_ = d.store.Save("blog", &config.ProjectConfig{Windows: []config.Window{{Name: "editor", Panes: []config.Pane{{Command: "vim"}}}}})
	_ = d.store.Save("blog2", &config.ProjectConfig{Windows: []config.Window{{Name: "main", Panes: []config.Pane{{Command: "bash"}}}}})

	if err := runCopyWith(d, []string{"blog", "blog2"}); err != nil {
		t.Fatalf("runCopyWith with --force: %v", err)
	}

	cfg, _, err := d.store.Load("blog2")
	if err != nil {
		t.Fatalf("destination config should exist: %v", err)
	}
	if cfg.Windows[0].Name != "editor" {
		t.Errorf("expected overwritten config to have window 'editor', got %q", cfg.Windows[0].Name)
	}
}

func TestRunCopyWith_OpensEditor(t *testing.T) {
	d := testDeps(t)
	editorCalled := false
	d.openEditor = func(path string) error {
		editorCalled = true
		return nil
	}
	_ = d.store.Save("blog", &config.ProjectConfig{Windows: []config.Window{{Name: "editor", Panes: []config.Pane{{Command: "vim"}}}}})

	if err := runCopyWith(d, []string{"blog", "blog2"}); err != nil {
		t.Fatalf("runCopyWith: %v", err)
	}

	if !editorCalled {
		t.Error("expected editor to be called")
	}
}

func TestRunCopyWith_PreservesContent(t *testing.T) {
	d := testDeps(t)
	d.openEditor = func(string) error { return nil }

	// Write a config with a comment that would be lost through marshal/unmarshal
	raw := []byte("# yaml-language-server: $schema=https://raw.githubusercontent.com/Drew-Daniels/nux/main/schemas/project.schema.json\n# My custom comment\nwindows:\n  - name: editor\n    panes:\n      - command: vim\n")
	_ = d.store.SaveRaw("blog", raw)

	if err := runCopyWith(d, []string{"blog", "blog2"}); err != nil {
		t.Fatalf("runCopyWith: %v", err)
	}

	destData, err := os.ReadFile(d.store.Path("blog2"))
	if err != nil {
		t.Fatalf("reading destination file: %v", err)
	}
	if !strings.Contains(string(destData), "# My custom comment") {
		t.Error("expected comment to be preserved in copied file")
	}
}
