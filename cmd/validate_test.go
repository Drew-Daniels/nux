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
		Windows: []config.Window{{Name: "editor"}},
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
		Command: "vim",
		Windows: []config.Window{{Name: "editor"}},
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
		Windows: []config.Window{{Name: "editor"}},
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
		Command: "vim",
		Windows: []config.Window{{Name: "editor"}},
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
		Windows: []config.Window{{Name: "editor"}},
	})

	if err := runValidateWith(d, []string{"blog"}); err != nil {
		t.Fatalf("runValidateWith single: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "[ok]") {
		t.Errorf("expected [ok], got %q", out)
	}
}

func TestRunValidateWith_All(t *testing.T) {
	d := testDeps(t)

	if err := runValidateWith(d, nil); err != nil {
		t.Fatalf("runValidateWith all: %v", err)
	}
}
