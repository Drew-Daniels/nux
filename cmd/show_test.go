package cmd

import (
	"strings"
	"testing"

	"github.com/Drew-Daniels/nux/internal/config"
)

func TestRunShowWith_WithConfig(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("blog", &config.ProjectConfig{
		Root:    d.global.ProjectsDir,
		Command: "vim",
		Env:     map[string]string{"NODE_ENV": "dev"},
	})

	if err := runShowWith(d, []string{"blog"}); err != nil {
		t.Fatalf("runShowWith: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "blog") {
		t.Error("expected project name in output")
	}
	if !strings.Contains(out, "NODE_ENV") {
		t.Error("expected env var in output")
	}
}

func TestRunShowWith_WithWindows(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("app", &config.ProjectConfig{
		Root: d.global.ProjectsDir,
		Windows: []config.Window{
			{
				Name:   "editor",
				Layout: "tiled",
				Panes: []config.Pane{
					{Command: "vim", Root: "src"},
				},
			},
		},
	})

	if err := runShowWith(d, []string{"app"}); err != nil {
		t.Fatalf("runShowWith: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "editor") {
		t.Error("expected window name in output")
	}
	if !strings.Contains(out, "vim") {
		t.Error("expected pane command in output")
	}
}

func TestRunShowWith_NotFound(t *testing.T) {
	d := testDeps(t)

	err := runShowWith(d, []string{"missing"})
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestRunShowWith_WithVarOverrides(t *testing.T) {
	d := testDeps(t)
	d.vars = map[string]string{"port": "9090"}
	_ = d.store.Save("api", &config.ProjectConfig{
		Root:    d.global.ProjectsDir,
		Command: "serve --port={{port}}",
		Vars:    map[string]string{"port": "8080"},
	})

	if err := runShowWith(d, []string{"api"}); err != nil {
		t.Fatalf("runShowWith: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "9090") {
		t.Errorf("expected overridden var value 9090 in output, got %q", out)
	}
}
