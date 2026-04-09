package cmd

import (
	"testing"

	"github.com/Drew-Daniels/nux/internal/config"
	"github.com/Drew-Daniels/nux/internal/tmux"
)

func TestRunRestartWith_FullSession(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	_ = d.store.Save("blog", &config.ProjectConfig{
		Root: d.global.ProjectsDir,
		Windows: []config.Window{
			{Name: "editor", Panes: []config.Pane{{Command: "vim"}}},
		},
	})

	if err := runRestartWith(d, []string{"blog"}); err != nil {
		t.Fatalf("runRestartWith: %v", err)
	}

	mock := d.client.(*tmux.MockClient)
	if !mock.Called("KillSession") {
		t.Error("expected KillSession")
	}
	if !mock.Called("NewSession") {
		t.Error("expected NewSession")
	}
}

func TestRunRestartWith_SingleWindow(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	_ = d.store.Save("blog", &config.ProjectConfig{
		Root: d.global.ProjectsDir,
		Windows: []config.Window{
			{Name: "editor", Panes: []config.Pane{{Command: "vim"}}},
			{Name: "server", Panes: []config.Pane{{Command: "go run ."}}},
		},
	})

	if err := runRestartWith(d, []string{"blog:server"}); err != nil {
		t.Fatalf("runRestartWith: %v", err)
	}

	mock := d.client.(*tmux.MockClient)
	if !mock.Called("KillWindow") {
		t.Error("expected KillWindow")
	}
	if !mock.Called("NewWindow") {
		t.Error("expected NewWindow")
	}
	if mock.Called("KillSession") {
		t.Error("should not kill entire session for window restart")
	}
}

func TestRunRestartWith_NotFound(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true

	err := runRestartWith(d, []string{"missing"})
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestRunRestartWith_FullSession_Attaches(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("blog", &config.ProjectConfig{
		Root: d.global.ProjectsDir,
		Windows: []config.Window{
			{Name: "editor", Panes: []config.Pane{{Command: "vim"}}},
		},
	})

	if err := runRestartWith(d, []string{"blog"}); err != nil {
		t.Fatalf("runRestartWith: %v", err)
	}

	mock := d.client.(*tmux.MockClient)
	if !mock.Called("AttachSession") {
		t.Error("expected AttachSession when noAttach=false")
	}
}

func TestRunRestartWith_Window_Attaches(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("blog", &config.ProjectConfig{
		Root: d.global.ProjectsDir,
		Windows: []config.Window{
			{Name: "editor", Panes: []config.Pane{{Command: "vim"}}},
			{Name: "server", Panes: []config.Pane{{Command: "go run ."}}},
		},
	})

	if err := runRestartWith(d, []string{"blog:editor"}); err != nil {
		t.Fatalf("runRestartWith: %v", err)
	}

	mock := d.client.(*tmux.MockClient)
	if !mock.Called("AttachSession") {
		t.Error("expected AttachSession when noAttach=false")
	}
}
