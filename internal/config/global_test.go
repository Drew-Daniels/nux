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
	if len(cfg.ProjectDirs) != 1 || cfg.ProjectDirs[0] != "~/projects" {
		t.Errorf("ProjectDirs = %v, want [~/projects]", cfg.ProjectDirs)
	}
	if cfg.Picker != "fzf" {
		t.Errorf("Picker = %q, want fzf", cfg.Picker)
	}
}

func TestLoadGlobalFrom_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := `
project_dirs: ~/code
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
	if len(cfg.ProjectDirs) != 1 || cfg.ProjectDirs[0] != "~/code" {
		t.Errorf("ProjectDirs = %v, want [~/code]", cfg.ProjectDirs)
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
	if len(cfg.ProjectDirs) != 1 || cfg.ProjectDirs[0] != "~/projects" {
		t.Errorf("ProjectDirs = %v, want [~/projects]", cfg.ProjectDirs)
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

func TestLoadGlobalFrom_ProjectDirsList(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := `
project_dirs:
  - ~/projects
  - ~/work
  - ~/docs
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadGlobalFrom(path)
	if err != nil {
		t.Fatalf("LoadGlobalFrom: %v", err)
	}
	want := StringOrList{"~/projects", "~/work", "~/docs"}
	if len(cfg.ProjectDirs) != len(want) {
		t.Fatalf("ProjectDirs = %v, want %v", cfg.ProjectDirs, want)
	}
	for i := range want {
		if cfg.ProjectDirs[i] != want[i] {
			t.Errorf("ProjectDirs[%d] = %q, want %q", i, cfg.ProjectDirs[i], want[i])
		}
	}
}

func TestFirstProjectDir(t *testing.T) {
	cfg := &GlobalConfig{ProjectDirs: StringOrList{"~/a", "~/b"}}
	if got := cfg.FirstProjectDir(); got != "~/a" {
		t.Errorf("FirstProjectDir = %q, want ~/a", got)
	}

	empty := &GlobalConfig{}
	if got := empty.FirstProjectDir(); got != "" {
		t.Errorf("FirstProjectDir on empty = %q, want empty", got)
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
	if len(cfg.ProjectDirs) != 1 || cfg.ProjectDirs[0] != "~/projects" {
		t.Errorf("ProjectDirs should keep default, got %v", cfg.ProjectDirs)
	}
}
