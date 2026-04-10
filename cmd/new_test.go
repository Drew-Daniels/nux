package cmd

import (
	"strings"
	"testing"

	"github.com/Drew-Daniels/nux/internal/config"
)

func TestRunNewWith_Creates(t *testing.T) {
	d := testDeps(t)
	editorCalled := false
	d.openEditor = func(path string) error {
		editorCalled = true
		return nil
	}

	if err := runNewWith(d, []string{"blog"}); err != nil {
		t.Fatalf("runNewWith: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "Created") {
		t.Errorf("expected 'Created' in output, got %q", out)
	}

	cfg, _, err := d.store.Load("blog")
	if err != nil {
		t.Fatalf("config should exist after creation: %v", err)
	}
	if len(cfg.Windows) != 1 || cfg.Windows[0].Name != "editor" {
		t.Errorf("unexpected config: %+v", cfg)
	}

	if !editorCalled {
		t.Error("expected editor to be called")
	}
}

func TestRunNewWith_NoEditorHint(t *testing.T) {
	d := testDeps(t)
	d.editor = ""
	d.openEditor = func(path string) error {
		return openInEditor(d, path)
	}

	if err := runNewWith(d, []string{"blog"}); err != nil {
		t.Fatalf("runNewWith: %v", err)
	}

	if _, _, err := d.store.Load("blog"); err != nil {
		t.Fatalf("config should exist after creation: %v", err)
	}

	stderr := stderrStr(d)
	if !strings.Contains(stderr, "$EDITOR") {
		t.Errorf("expected hint about $EDITOR in stderr, got %q", stderr)
	}
}

func TestRunNewWith_AlreadyExists(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("blog", &config.ProjectConfig{Command: "vim"})

	err := runNewWith(d, []string{"blog"})
	if err == nil {
		t.Fatal("expected error for existing config")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("error = %q, expected 'already exists'", err.Error())
	}
}
