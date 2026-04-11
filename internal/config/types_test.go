package config

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDefaultSession_UnmarshalYAML_String(t *testing.T) {
	input := `"htop"`
	var ds DefaultSession
	if err := yaml.Unmarshal([]byte(input), &ds); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if ds.Command != "htop" {
		t.Errorf("Command = %q, want htop", ds.Command)
	}
	if len(ds.Windows) != 0 {
		t.Errorf("expected no windows, got %d", len(ds.Windows))
	}
}

func TestDefaultSession_UnmarshalYAML_Object(t *testing.T) {
	input := `
windows:
  - name: editor
    panes:
      - vim
`
	var ds DefaultSession
	if err := yaml.Unmarshal([]byte(input), &ds); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if ds.Command != "" {
		t.Errorf("Command = %q, want empty", ds.Command)
	}
	if len(ds.Windows) != 1 {
		t.Fatalf("expected 1 window, got %d", len(ds.Windows))
	}
	if ds.Windows[0].Name != "editor" {
		t.Errorf("window name = %q, want editor", ds.Windows[0].Name)
	}
}

func TestWindow_UnmarshalYAML_RejectsCommand(t *testing.T) {
	input := `
name: editor
command: vim
`
	var w Window
	err := yaml.Unmarshal([]byte(input), &w)
	if err == nil {
		t.Fatal("expected error for window-level command")
	}
	if !strings.Contains(err.Error(), "not a valid window field") {
		t.Errorf("expected actionable error, got %q", err.Error())
	}
	if !strings.Contains(err.Error(), "panes: [vim]") {
		t.Errorf("expected panes suggestion, got %q", err.Error())
	}
}

func TestWindow_UnmarshalYAML_ValidPanes(t *testing.T) {
	input := `
name: editor
panes:
  - vim
`
	var w Window
	if err := yaml.Unmarshal([]byte(input), &w); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if w.Name != "editor" {
		t.Errorf("Name = %q, want editor", w.Name)
	}
	if len(w.Panes) != 1 || w.Panes[0].Command != "vim" {
		t.Errorf("Panes = %+v, want [{Command: vim}]", w.Panes)
	}
}

func TestPane_UnmarshalYAML_String(t *testing.T) {
	input := `"vim ."`
	var p Pane
	if err := yaml.Unmarshal([]byte(input), &p); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if p.Command != "vim ." {
		t.Errorf("Command = %q, want 'vim .'", p.Command)
	}
}

func TestPane_UnmarshalYAML_Object(t *testing.T) {
	input := `
command: make watch
root: src
`
	var p Pane
	if err := yaml.Unmarshal([]byte(input), &p); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if p.Command != "make watch" {
		t.Errorf("Command = %q, want 'make watch'", p.Command)
	}
	if p.Root != "src" {
		t.Errorf("Root = %q, want src", p.Root)
	}
}

func TestProjectConfig_UnmarshalYAML_Full(t *testing.T) {
	input := `
root: ~/projects/myapp
env:
  NODE_ENV: development
windows:
  - name: editor
    env:
      PORT: "3000"
    panes:
      - vim
      - command: make watch
        root: backend
`
	var cfg ProjectConfig
	if err := yaml.Unmarshal([]byte(input), &cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if cfg.Root != "~/projects/myapp" {
		t.Errorf("Root = %q", cfg.Root)
	}
	if cfg.Env["NODE_ENV"] != "development" {
		t.Errorf("Env[NODE_ENV] = %q", cfg.Env["NODE_ENV"])
	}
	if len(cfg.Windows) != 1 {
		t.Fatalf("expected 1 window, got %d", len(cfg.Windows))
	}
	w := cfg.Windows[0]
	if w.Env["PORT"] != "3000" {
		t.Errorf("window Env[PORT] = %q", w.Env["PORT"])
	}
	if len(w.Panes) != 2 {
		t.Fatalf("expected 2 panes, got %d", len(w.Panes))
	}
	if w.Panes[0].Command != "vim" {
		t.Errorf("pane 0 command = %q, want vim", w.Panes[0].Command)
	}
	if w.Panes[1].Root != "backend" {
		t.Errorf("pane 1 root = %q, want backend", w.Panes[1].Root)
	}
}
