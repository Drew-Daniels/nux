package cmd

import (
	"strings"
	"testing"

	"github.com/Drew-Daniels/nux/internal/config"
)

func TestValidateAll_NoConfigs(t *testing.T) {
	d := testDeps(t)

	if err := validateAll(d); err != nil {
		t.Fatalf("validateAll: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "No project configs") {
		t.Errorf("expected 'No project configs', got %q", out)
	}
}

func TestValidateAll_AllValid(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("blog", &config.ProjectConfig{
		Windows: []config.Window{{Name: "editor", Panes: []config.Pane{{Command: "vim"}}}},
	})

	if err := validateAll(d); err != nil {
		t.Fatalf("validateAll: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "[ok]") {
		t.Errorf("expected [ok], got %q", out)
	}
}

func TestValidateAll_WithErrors(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("bad", &config.ProjectConfig{
		Windows: []config.Window{{Name: "", Layout: "bogus"}},
	})

	err := validateAll(d)
	if err == nil {
		t.Fatal("expected error for invalid config")
	}

	errOut := stderrStr(d)
	if !strings.Contains(errOut, "[error]") {
		t.Errorf("expected [error] in stderr, got %q", errOut)
	}
}

func TestValidateProject_Valid(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("blog", &config.ProjectConfig{
		Windows: []config.Window{{Name: "editor", Panes: []config.Pane{{Command: "vim"}}}},
	})

	if err := validateProject(d, "blog"); err != nil {
		t.Fatalf("validateProject: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "[ok]") {
		t.Errorf("expected [ok], got %q", out)
	}
}

func TestValidateProject_Invalid(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("bad", &config.ProjectConfig{
		Windows: []config.Window{{Name: "", Layout: "bogus"}},
	})

	err := validateProject(d, "bad")
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}

func TestValidateProject_NotFound(t *testing.T) {
	d := testDeps(t)

	err := validateProject(d, "missing")
	if err == nil {
		t.Fatal("expected error for missing config")
	}
}

func TestRunValidateWith_Delegating(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("blog", &config.ProjectConfig{
		Windows: []config.Window{{Name: "editor", Panes: []config.Pane{{Command: "vim"}}}},
	})

	if err := runValidateWith(d, []string{"blog"}); err != nil {
		t.Fatalf("runValidateWith single: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "[ok]") {
		t.Errorf("expected [ok], got %q", out)
	}
}

func TestRunValidateWith_Glob(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("web-api", &config.ProjectConfig{
		Windows: []config.Window{{Name: "editor", Panes: []config.Pane{{Command: "vim"}}}},
	})
	_ = d.store.Save("web-ui", &config.ProjectConfig{
		Windows: []config.Window{{Name: "editor", Panes: []config.Pane{{Command: "vim"}}}},
	})

	if err := runValidateWith(d, []string{"web+"}); err != nil {
		t.Fatalf("runValidateWith web+: %v", err)
	}

	out := stdoutStr(d)
	if strings.Count(out, "[ok]") != 2 {
		t.Errorf("expected two [ok] lines, got %q", out)
	}
}

func TestRunValidateWith_MultipleNames(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("blog", &config.ProjectConfig{
		Windows: []config.Window{{Name: "editor", Panes: []config.Pane{{Command: "vim"}}}},
	})
	_ = d.store.Save("api", &config.ProjectConfig{
		Windows: []config.Window{{Name: "server", Panes: []config.Pane{{Command: "go run ."}}}},
	})

	if err := runValidateWith(d, []string{"blog", "api"}); err != nil {
		t.Fatalf("runValidateWith: %v", err)
	}

	out := stdoutStr(d)
	if strings.Count(out, "[ok]") != 2 {
		t.Errorf("expected two [ok] lines, got %q", out)
	}
}

func TestRunValidateWith_Group(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("alpha", &config.ProjectConfig{
		Windows: []config.Window{{Name: "editor", Panes: []config.Pane{{Command: "vim"}}}},
	})
	_ = d.store.Save("bravo", &config.ProjectConfig{
		Windows: []config.Window{{Name: "editor", Panes: []config.Pane{{Command: "vim"}}}},
	})
	d.global.Groups = map[string][]string{"batch": {"alpha", "bravo"}}

	if err := runValidateWith(d, []string{"@batch"}); err != nil {
		t.Fatalf("runValidateWith @batch: %v", err)
	}

	out := stdoutStr(d)
	if strings.Count(out, "[ok]") != 2 {
		t.Errorf("expected two [ok] lines, got %q", out)
	}
}

func TestRunValidateWith_All(t *testing.T) {
	d := testDeps(t)

	if err := runValidateWith(d, nil); err != nil {
		t.Fatalf("runValidateWith all: %v", err)
	}
}
