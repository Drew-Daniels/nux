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

func TestRunShowWith_Raw(t *testing.T) {
	d := testDeps(t)
	showRaw = true
	t.Cleanup(func() { showRaw = false })

	_ = d.store.Save("api", &config.ProjectConfig{
		Root:    "~/projects/api",
		Command: "echo {{greeting}} $HOME",
		Vars:    map[string]string{"greeting": "hi"},
		Env:     map[string]string{"SECRET": "${API_KEY}"},
	})

	if err := runShowWith(d, []string{"api"}); err != nil {
		t.Fatalf("runShowWith --raw: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "{{greeting}}") {
		t.Errorf("expected raw {{greeting}} placeholder preserved, got %q", out)
	}
	if !strings.Contains(out, "${API_KEY}") {
		t.Errorf("expected raw ${API_KEY} preserved, got %q", out)
	}
	if !strings.Contains(out, "$HOME") {
		t.Errorf("expected raw $HOME preserved, got %q", out)
	}
}

func TestRunShowWith_Raw_NotFound(t *testing.T) {
	d := testDeps(t)
	showRaw = true
	t.Cleanup(func() { showRaw = false })

	err := runShowWith(d, []string{"missing"})
	if err == nil {
		t.Fatal("expected error for missing project with --raw")
	}
}

func TestRunShowWith_GlobMultiDoc(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("web-api", &config.ProjectConfig{
		Root:    d.global.ProjectsDir,
		Command: "a",
	})
	_ = d.store.Save("web-ui", &config.ProjectConfig{
		Root:    d.global.ProjectsDir,
		Command: "b",
	})

	if err := runShowWith(d, []string{"web+"}); err != nil {
		t.Fatalf("runShowWith: %v", err)
	}

	out := stdoutStr(d)
	if strings.Count(out, "---\n") != 1 {
		t.Errorf("expected exactly one YAML document separator for 2 projects, got %q", out)
	}
	if !strings.Contains(out, "web-api") || !strings.Contains(out, "web-ui") {
		t.Errorf("expected both project names in output, got %q", out)
	}
}

func TestRunShowWith_Group(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("alpha", &config.ProjectConfig{
		Root:    d.global.ProjectsDir,
		Command: "x",
	})
	_ = d.store.Save("bravo", &config.ProjectConfig{
		Root:    d.global.ProjectsDir,
		Command: "y",
	})
	d.global.Groups = map[string][]string{"batch": {"alpha", "bravo"}}

	if err := runShowWith(d, []string{"@batch"}); err != nil {
		t.Fatalf("runShowWith: %v", err)
	}

	out := stdoutStr(d)
	if strings.Count(out, "---\n") != 1 {
		t.Errorf("expected one separator for 2 group members, got %q", out)
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
