package cmd

import (
	"strings"
	"testing"

	"github.com/Drew-Daniels/nux/internal/config"
)

func TestRunDeleteWith_Force(t *testing.T) {
	d := testDeps(t)
	d.deleteForce = true
	_ = d.store.Save("blog", &config.ProjectConfig{Windows: []config.Window{{Name: "main", Panes: []config.Pane{{Command: "vim"}}}}})

	if err := runDeleteWith(d, []string{"blog"}); err != nil {
		t.Fatalf("runDeleteWith: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "Deleted") {
		t.Errorf("expected 'Deleted' in output, got %q", out)
	}

	_, _, err := d.store.Load("blog")
	if err == nil {
		t.Error("config should be deleted")
	}
}

func TestRunDeleteWith_Confirmed(t *testing.T) {
	d := testDeps(t)
	d.confirm = func(string) (bool, error) { return true, nil }
	_ = d.store.Save("blog", &config.ProjectConfig{Windows: []config.Window{{Name: "main", Panes: []config.Pane{{Command: "vim"}}}}})

	if err := runDeleteWith(d, []string{"blog"}); err != nil {
		t.Fatalf("runDeleteWith: %v", err)
	}

	_, _, err := d.store.Load("blog")
	if err == nil {
		t.Error("config should be deleted")
	}
}

func TestRunDeleteWith_Cancelled(t *testing.T) {
	d := testDeps(t)
	d.confirm = func(string) (bool, error) { return false, nil }
	_ = d.store.Save("blog", &config.ProjectConfig{Windows: []config.Window{{Name: "main", Panes: []config.Pane{{Command: "vim"}}}}})

	if err := runDeleteWith(d, []string{"blog"}); err != nil {
		t.Fatalf("runDeleteWith: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "Cancelled") {
		t.Errorf("expected 'Cancelled' in output, got %q", out)
	}

	_, _, err := d.store.Load("blog")
	if err != nil {
		t.Error("config should still exist after cancel")
	}
}

func TestRunDeleteWith_NotFound(t *testing.T) {
	d := testDeps(t)
	d.deleteForce = true

	err := runDeleteWith(d, []string{"missing"})
	if err == nil {
		t.Fatal("expected error for missing config")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error = %q, expected 'not found'", err.Error())
	}
}
