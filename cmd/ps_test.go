package cmd

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Drew-Daniels/nux/internal/tmux"
)

func TestRunPsWith_Sessions(t *testing.T) {
	d := testDeps(t)
	mock := d.client.(*tmux.MockClient)
	mock.ListSessionsReturn = []tmux.SessionInfo{
		{Name: "blog", Windows: 2, Attached: true, Created: time.Now().Add(-30 * time.Minute)},
		{Name: "api", Windows: 1, Attached: false, Created: time.Now().Add(-2 * time.Hour)},
	}

	if err := runPsWith(d); err != nil {
		t.Fatalf("runPsWith: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "blog") {
		t.Error("expected blog in output")
	}
	if !strings.Contains(out, "api") {
		t.Error("expected api in output")
	}
	if !strings.Contains(out, "yes") {
		t.Error("expected 'yes' for attached session")
	}
}

func TestRunPsWith_Empty(t *testing.T) {
	d := testDeps(t)

	if err := runPsWith(d); err != nil {
		t.Fatalf("runPsWith: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "No running sessions") {
		t.Errorf("expected 'No running sessions', got %q", out)
	}
}

func TestRunPsWith_ListError(t *testing.T) {
	d := testDeps(t)
	mock := d.client.(*tmux.MockClient)
	mock.ListSessionsError = fmt.Errorf("tmux not running")

	err := runPsWith(d)
	if err == nil {
		t.Fatal("expected error from ListSessions failure")
	}
	if !strings.Contains(err.Error(), "listing sessions") {
		t.Errorf("error = %q, expected 'listing sessions'", err.Error())
	}
}
