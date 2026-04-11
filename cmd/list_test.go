package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Drew-Daniels/nux/internal/config"
	"github.com/Drew-Daniels/nux/internal/tmux"
)

func TestRunListWith_WithProjects(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("blog", &config.ProjectConfig{Root: "~/blog", Windows: []config.Window{{Name: "main", Panes: []config.Pane{{Command: "vim"}}}}})
	_ = d.store.Save("api", &config.ProjectConfig{Root: "~/api", Windows: []config.Window{{Name: "main", Panes: []config.Pane{{Command: "go run ."}}}}})

	mock := d.client.(*tmux.MockClient)
	mock.ListSessionsReturn = []tmux.SessionInfo{
		{Name: "blog", Windows: 2, Attached: true},
	}

	if err := runListWith(d); err != nil {
		t.Fatalf("runListWith: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "blog") {
		t.Error("expected blog in output")
	}
	if !strings.Contains(out, "api") {
		t.Error("expected api in output")
	}
	if !strings.Contains(out, "running") {
		t.Error("expected 'running' status for blog")
	}
}

func TestRunListWith_Empty(t *testing.T) {
	d := testDeps(t)

	if err := runListWith(d); err != nil {
		t.Fatalf("runListWith: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "NAME") {
		t.Error("expected header row even with no projects")
	}
}

func TestRunListWith_SessionsError(t *testing.T) {
	d := testDeps(t)
	mock := d.client.(*tmux.MockClient)
	mock.ListSessionsError = fmt.Errorf("tmux not running")
	_ = d.store.Save("blog", &config.ProjectConfig{Root: "~/blog", Windows: []config.Window{{Name: "main", Panes: []config.Pane{{Command: "vim"}}}}})

	if err := runListWith(d); err != nil {
		t.Fatalf("runListWith: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "blog") {
		t.Error("expected blog in output even when sessions fail")
	}
}
