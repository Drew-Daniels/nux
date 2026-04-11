package config

import (
	"testing"
)

func TestInterpolate_CustomVars(t *testing.T) {
	cfg := &ProjectConfig{
		Root: "~/projects/{{name}}",
		Vars: map[string]string{
			"name": "my-project",
		},
		Windows: []Window{
			{Name: "editor", Panes: []Pane{{Command: "{{name}}"}}},
		},
	}
	if err := NewInterpolator().Interpolate(cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Root != "~/projects/my-project" {
		t.Errorf("root = %q, want ~/projects/my-project", cfg.Root)
	}
	if cfg.Windows[0].Panes[0].Command != "my-project" {
		t.Errorf("pane command = %q, want my-project", cfg.Windows[0].Panes[0].Command)
	}
}

func TestInterpolate_EnvVars(t *testing.T) {
	t.Setenv("NUX_TEST_PORT", "3000")
	cfg := &ProjectConfig{
		Env: map[string]string{
			"PORT": "${NUX_TEST_PORT}",
		},
	}
	if err := NewInterpolator().Interpolate(cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Env["PORT"] != "3000" {
		t.Errorf("env PORT = %q, want 3000", cfg.Env["PORT"])
	}
}

func TestInterpolate_VarsBeforeEnv(t *testing.T) {
	t.Setenv("NUX_TEST_DIR", "/opt")

	cfg := &ProjectConfig{
		Root: "${NUX_TEST_DIR}/{{sub}}",
		Vars: map[string]string{"sub": "app"},
	}
	if err := NewInterpolator().Interpolate(cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Root != "/opt/app" {
		t.Errorf("root = %q, want /opt/app", cfg.Root)
	}
}

func TestInterpolate_Hooks(t *testing.T) {
	cfg := &ProjectConfig{
		OnStart:  []string{"echo {{name}}"},
		OnStop:   []string{"echo stop"},
		OnReady:  []string{"echo ready"},
		OnDetach: []string{"echo detach"},
		Vars:     map[string]string{"name": "test"},
	}
	if err := NewInterpolator().Interpolate(cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.OnStart[0] != "echo test" {
		t.Errorf("on_start = %q, want echo test", cfg.OnStart[0])
	}
}

func TestInterpolate_WindowsAndPanes(t *testing.T) {
	cfg := &ProjectConfig{
		Vars: map[string]string{"dir": "src", "editor": "nvim"},
		Windows: []Window{
			{
				Root: "{{dir}}",
				Panes: []Pane{
					{Root: "{{dir}}/sub", Command: "{{editor}} ."},
				},
			},
		},
	}
	if err := NewInterpolator().Interpolate(cfg); err != nil {
		t.Fatal(err)
	}
	w := cfg.Windows[0]
	if w.Root != "src" {
		t.Errorf("window root = %q, want src", w.Root)
	}
	p := w.Panes[0]
	if p.Root != "src/sub" {
		t.Errorf("pane root = %q, want src/sub", p.Root)
	}
	if p.Command != "nvim ." {
		t.Errorf("pane command = %q, want nvim .", p.Command)
	}
}

func TestInterpolate_BacktickCommand(t *testing.T) {
	cfg := &ProjectConfig{
		Vars: map[string]string{"out": "`echo hello`"},
		Windows: []Window{
			{Name: "main", Panes: []Pane{{Command: "{{out}}"}}},
		},
	}
	if err := NewInterpolator().Interpolate(cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Windows[0].Panes[0].Command != "hello" {
		t.Errorf("pane command = %q, want hello", cfg.Windows[0].Panes[0].Command)
	}
}

func TestInterpolate_WindowEnv(t *testing.T) {
	t.Setenv("NUX_TEST_HOST", "localhost")
	cfg := &ProjectConfig{
		Vars: map[string]string{"port": "3000"},
		Windows: []Window{
			{
				Name: "api",
				Env: map[string]string{
					"PORT":         "{{port}}",
					"DATABASE_URL": "postgres://${NUX_TEST_HOST}/mydb",
				},
				Panes: []Pane{{Command: "npm start"}},
			},
		},
	}
	if err := NewInterpolator().Interpolate(cfg); err != nil {
		t.Fatal(err)
	}
	w := cfg.Windows[0]
	if w.Env["PORT"] != "3000" {
		t.Errorf("window env PORT = %q, want 3000", w.Env["PORT"])
	}
	if w.Env["DATABASE_URL"] != "postgres://localhost/mydb" {
		t.Errorf("window env DATABASE_URL = %q, want postgres://localhost/mydb", w.Env["DATABASE_URL"])
	}
}

func TestInterpolateVars_OnlyExpandsVars(t *testing.T) {
	t.Setenv("NUX_TEST_SHOULD_NOT_EXPAND", "expanded")
	cfg := &ProjectConfig{
		Root: "$NUX_TEST_SHOULD_NOT_EXPAND/{{name}}",
		Vars: map[string]string{
			"name": "my-project",
		},
		Windows: []Window{
			{Name: "main", Panes: []Pane{{Command: "{{name}}"}}},
		},
	}
	if err := NewInterpolator().InterpolateVars(cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Windows[0].Panes[0].Command != "my-project" {
		t.Errorf("pane command = %q, want my-project", cfg.Windows[0].Panes[0].Command)
	}
	if cfg.Root != "$NUX_TEST_SHOULD_NOT_EXPAND/my-project" {
		t.Errorf("root = %q, want $NUX_TEST_SHOULD_NOT_EXPAND/my-project (env should not expand)", cfg.Root)
	}
}

func TestInterpolate_NoVars(t *testing.T) {
	cfg := &ProjectConfig{
		Root: "~/projects/plain",
		Windows: []Window{
			{Name: "main", Panes: []Pane{{Command: "vim"}}},
		},
	}
	if err := NewInterpolator().Interpolate(cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Root != "~/projects/plain" {
		t.Errorf("root = %q, want ~/projects/plain", cfg.Root)
	}
}
