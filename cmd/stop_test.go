package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Drew-Daniels/nux/internal/config"
	"github.com/Drew-Daniels/nux/internal/tmux"
)

func TestRunStopWith(t *testing.T) {
	d := testDeps(t)
	mock := d.client.(*tmux.MockClient)
	mock.HasSessionReturn = true
	_ = d.store.Save("blog", &config.ProjectConfig{Windows: []config.Window{{Name: "main", Panes: []config.Pane{{Command: "vim"}}}}})

	if err := runStopWith(d, []string{"blog"}); err != nil {
		t.Fatalf("runStopWith: %v", err)
	}

	found := false
	for _, c := range mock.Calls {
		if c.Method == "KillSession" && len(c.Args) > 0 && c.Args[0] == "blog" {
			found = true
		}
	}
	if !found {
		t.Error("expected KillSession with 'blog'")
	}
}

func TestRunStopWith_NormalizesSessionName(t *testing.T) {
	d := testDeps(t)
	mock := d.client.(*tmux.MockClient)
	mock.HasSessionReturn = true

	if err := runStopWith(d, []string{"my.project"}); err != nil {
		t.Fatalf("runStopWith: %v", err)
	}

	found := false
	for _, c := range mock.Calls {
		if c.Method == "KillSession" && len(c.Args) > 0 && c.Args[0] == "my_project" {
			found = true
		}
	}
	if !found {
		t.Error("expected KillSession with normalized name 'my_project'")
	}
}

func TestRunStopWith_NotRunning(t *testing.T) {
	d := testDeps(t)

	err := runStopWith(d, []string{"studios"})
	if err == nil {
		t.Fatal("expected error for session not running")
	}
	if !strings.Contains(err.Error(), "is not running") {
		t.Errorf("error = %q, expected 'is not running'", err.Error())
	}
}

func TestRunStopWith_ExpandError(t *testing.T) {
	d := testDeps(t)
	err := runStopWith(d, []string{"@missing"})
	if err == nil {
		t.Fatal("expected error for missing group")
	}
}

func TestRunStopWith_KillError(t *testing.T) {
	d := testDeps(t)
	mock := d.client.(*tmux.MockClient)
	mock.HasSessionReturn = true
	mock.DefaultError = fmt.Errorf("kill failed")
	_ = d.store.Save("blog", &config.ProjectConfig{Windows: []config.Window{{Name: "main", Panes: []config.Pane{{Command: "vim"}}}}})

	err := runStopWith(d, []string{"blog"})
	if err == nil {
		t.Fatal("expected error from KillSession failure")
	}
	if !strings.Contains(err.Error(), "stopping session") {
		t.Errorf("error = %q, expected 'stopping session'", err.Error())
	}
}

func TestRunStopAllWith(t *testing.T) {
	d := testDeps(t)
	mock := d.client.(*tmux.MockClient)
	mock.ListSessionsReturn = []tmux.SessionInfo{
		{Name: "a"},
		{Name: "b"},
	}

	if err := runStopAllWith(d); err != nil {
		t.Fatalf("runStopAllWith: %v", err)
	}

	kills := 0
	for _, c := range mock.Calls {
		if c.Method == "KillSession" {
			kills++
		}
	}
	if kills != 2 {
		t.Errorf("expected 2 KillSession calls, got %d", kills)
	}

	out := stderrStr(d)
	if !strings.Contains(out, "Stopping a (1/2)") || !strings.Contains(out, "Stopping b (2/2)") {
		t.Errorf("stderr = %q, want progress lines for a and b", out)
	}
}

func TestRunStopAllWith_NoSessionsMessage(t *testing.T) {
	d := testDeps(t)
	mock := d.client.(*tmux.MockClient)
	mock.ListSessionsReturn = []tmux.SessionInfo{}

	if err := runStopAllWith(d); err != nil {
		t.Fatalf("runStopAllWith: %v", err)
	}

	out := stderrStr(d)
	if !strings.Contains(out, "No tmux sessions running") {
		t.Errorf("stderr = %q, want empty-session message", out)
	}
}
