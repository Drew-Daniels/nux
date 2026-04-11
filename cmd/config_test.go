package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunConfigWith_CreatesWhenMissing(t *testing.T) {
	d := testDeps(t)
	cfgDir := filepath.Join(t.TempDir(), "nux")

	editorPath := ""
	d.openEditor = func(path string) error {
		editorPath = path
		return nil
	}

	if err := runConfigWith(d, cfgDir); err != nil {
		t.Fatalf("runConfigWith: %v", err)
	}

	cfgPath := filepath.Join(cfgDir, "config.yaml")
	if _, err := os.Stat(cfgPath); err != nil {
		t.Fatalf("config file should exist: %v", err)
	}

	projectsDir := filepath.Join(cfgDir, "projects")
	info, err := os.Stat(projectsDir)
	if err != nil {
		t.Fatalf("projects dir should exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("projects path should be a directory")
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "Created") {
		t.Errorf("expected 'Created' in output, got %q", out)
	}

	if editorPath != cfgPath {
		t.Errorf("editor path = %q, want %q", editorPath, cfgPath)
	}
}

func TestRunConfigWith_OpensEditorWhenExists(t *testing.T) {
	d := testDeps(t)
	cfgDir := t.TempDir()

	cfgPath := filepath.Join(cfgDir, "config.yaml")
	_ = os.WriteFile(cfgPath, []byte("picker: gum\n"), 0o644)

	editorPath := ""
	d.openEditor = func(path string) error {
		editorPath = path
		return nil
	}

	if err := runConfigWith(d, cfgDir); err != nil {
		t.Fatalf("runConfigWith: %v", err)
	}

	if editorPath != cfgPath {
		t.Errorf("editor path = %q, want %q", editorPath, cfgPath)
	}

	out := stdoutStr(d)
	if strings.Contains(out, "Created") {
		t.Error("should not print 'Created' when config already exists")
	}
}

func TestRunConfigWith_DoesNotOverwriteExisting(t *testing.T) {
	d := testDeps(t)
	cfgDir := t.TempDir()

	cfgPath := filepath.Join(cfgDir, "config.yaml")
	original := []byte("picker: gum\n")
	_ = os.WriteFile(cfgPath, original, 0o644)

	if err := runConfigWith(d, cfgDir); err != nil {
		t.Fatalf("runConfigWith: %v", err)
	}

	data, _ := os.ReadFile(cfgPath)
	if string(data) != string(original) {
		t.Errorf("config was overwritten: got %q, want %q", string(data), string(original))
	}
}

func TestRunConfigWith_ScaffoldContainsDefaults(t *testing.T) {
	d := testDeps(t)
	cfgDir := filepath.Join(t.TempDir(), "nux")

	if err := runConfigWith(d, cfgDir); err != nil {
		t.Fatalf("runConfigWith: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(cfgDir, "config.yaml"))
	if err != nil {
		t.Fatalf("reading config: %v", err)
	}

	content := string(data)
	for _, want := range []string{"project_dirs", "picker", "picker_on_bare", "zoxide"} {
		if !strings.Contains(content, want) {
			t.Errorf("scaffold missing %q", want)
		}
	}
}
