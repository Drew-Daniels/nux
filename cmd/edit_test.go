package cmd

import (
	"strings"
	"testing"

	"github.com/Drew-Daniels/nux/internal/config"
)

func TestRunEditWith_OpensEditor(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("blog", &config.ProjectConfig{Command: "vim"})

	editorPath := ""
	d.openEditor = func(path string) error {
		editorPath = path
		return nil
	}

	if err := runEditWith(d, []string{"blog"}); err != nil {
		t.Fatalf("runEditWith: %v", err)
	}

	if editorPath == "" {
		t.Error("expected editor to be called")
	}
	if !strings.HasSuffix(editorPath, "blog.yaml") {
		t.Errorf("editor path = %q, expected blog.yaml suffix", editorPath)
	}
}

func TestRunEditWith_NoEditor(t *testing.T) {
	d := testDeps(t)
	d.editor = ""
	_ = d.store.Save("blog", &config.ProjectConfig{Command: "vim"})

	err := runEditWith(d, []string{"blog"})
	if err == nil {
		t.Fatal("expected error when EDITOR not set")
	}
	if !strings.Contains(err.Error(), "EDITOR") {
		t.Errorf("error = %q, expected EDITOR message", err.Error())
	}
}

func TestRunEditWith_NotFound(t *testing.T) {
	d := testDeps(t)

	err := runEditWith(d, []string{"missing"})
	if err == nil {
		t.Fatal("expected error for missing config")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error = %q, expected 'not found'", err.Error())
	}
}
