package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadGlobalFrom_Missing(t *testing.T) {
	cfg, err := LoadGlobalFrom("/nonexistent/config.yaml")
	if err != nil {
		t.Fatalf("missing file should return defaults, got error: %v", err)
	}
	if cfg.ProjectsDir != "~/projects" {
		t.Errorf("ProjectsDir = %q, want ~/projects", cfg.ProjectsDir)
	}
	if cfg.Picker != "fzf" {
		t.Errorf("Picker = %q, want fzf", cfg.Picker)
	}
}

func TestLoadGlobalFrom_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := `
projects_dir: ~/code
picker: gum
zoxide: true
picker_on_bare: true
default_shell: /bin/zsh
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadGlobalFrom(path)
	if err != nil {
		t.Fatalf("LoadGlobalFrom: %v", err)
	}
	if cfg.ProjectsDir != "~/code" {
		t.Errorf("ProjectsDir = %q, want ~/code", cfg.ProjectsDir)
	}
	if cfg.Picker != "gum" {
		t.Errorf("Picker = %q, want gum", cfg.Picker)
	}
	if !cfg.Zoxide {
		t.Error("Zoxide should be true")
	}
	if !cfg.PickerOnBare {
		t.Error("PickerOnBare should be true")
	}
	if cfg.DefaultShell != "/bin/zsh" {
		t.Errorf("DefaultShell = %q, want /bin/zsh", cfg.DefaultShell)
	}
}

func TestLoadGlobalFrom_Invalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(":\n  bad: ["), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadGlobalFrom(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestGlobalDefaults(t *testing.T) {
	cfg := GlobalDefaults()
	if cfg.ProjectsDir != "~/projects" {
		t.Errorf("ProjectsDir = %q, want ~/projects", cfg.ProjectsDir)
	}
	if cfg.Picker != "fzf" {
		t.Errorf("Picker = %q, want fzf", cfg.Picker)
	}
	if cfg.PickerOnBare {
		t.Error("PickerOnBare should be false by default")
	}
	if cfg.Zoxide {
		t.Error("Zoxide should be false by default")
	}
}

func TestLoadGlobalFrom_PartialOverride(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("picker: gum\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadGlobalFrom(path)
	if err != nil {
		t.Fatalf("LoadGlobalFrom: %v", err)
	}
	if cfg.Picker != "gum" {
		t.Errorf("Picker = %q, want gum", cfg.Picker)
	}
	if cfg.ProjectsDir != "~/projects" {
		t.Errorf("ProjectsDir should keep default, got %q", cfg.ProjectsDir)
	}
}
