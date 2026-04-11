package cmd

import (
	"testing"

	"github.com/Drew-Daniels/nux/internal/config"
	"github.com/Drew-Daniels/nux/internal/tmux"
)

func TestRunRestartWith_VarOverrides(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	d.vars = map[string]string{"greeting": "hello"}
	// No vars in file so Resolve leaves {{greeting}} in Command; applyVarOverrides
	// merges CLI vars before RestartSession builds.
	_ = d.store.Save("blog", &config.ProjectConfig{
		Root:    d.global.ProjectDirs[0],
		Command: "echo {{greeting}}",
	})

	if err := runRestartWith(d, []string{"blog"}); err != nil {
		t.Fatalf("runRestartWith: %v", err)
	}

	mock := d.client.(*tmux.MockClient)
	found := false
	for _, c := range mock.Calls {
		if c.Method == "SendKeys" && len(c.Args) >= 2 && c.Args[1] == "echo hello" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected SendKeys with CLI --var override (echo hello)")
	}
}

func TestRunRestartWith_FullSession(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	mock := d.client.(*tmux.MockClient)
	mock.HasSessionReturn = true
	_ = d.store.Save("blog", &config.ProjectConfig{
		Root: d.global.ProjectDirs[0],
		Windows: []config.Window{
			{Name: "editor", Panes: []config.Pane{{Command: "vim"}}},
		},
	})

	if err := runRestartWith(d, []string{"blog"}); err != nil {
		t.Fatalf("runRestartWith: %v", err)
	}

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
		Root: d.global.ProjectDirs[0],
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
	mock := d.client.(*tmux.MockClient)
	mock.HasSessionReturn = true
	_ = d.store.Save("blog", &config.ProjectConfig{
		Root: d.global.ProjectDirs[0],
		Windows: []config.Window{
			{Name: "editor", Panes: []config.Pane{{Command: "vim"}}},
		},
	})

	if err := runRestartWith(d, []string{"blog"}); err != nil {
		t.Fatalf("runRestartWith: %v", err)
	}

	if !mock.Called("AttachSession") {
		t.Error("expected AttachSession when noAttach=false")
	}
}

func TestRunRestartWith_Window_Attaches(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("blog", &config.ProjectConfig{
		Root: d.global.ProjectDirs[0],
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

func TestRunRestartWith_GlobMulti(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	mock := d.client.(*tmux.MockClient)
	mock.HasSessionReturn = true
	_ = d.store.Save("web-api", &config.ProjectConfig{
		Root: d.global.ProjectDirs[0],
		Windows: []config.Window{
			{Name: "editor", Panes: []config.Pane{{Command: "vim"}}},
		},
	})
	_ = d.store.Save("web-ui", &config.ProjectConfig{
		Root: d.global.ProjectDirs[0],
		Windows: []config.Window{
			{Name: "editor", Panes: []config.Pane{{Command: "vim"}}},
		},
	})

	if err := runRestartWith(d, []string{"web+"}); err != nil {
		t.Fatalf("runRestartWith: %v", err)
	}

	n := 0
	for _, c := range mock.Calls {
		if c.Method == "KillSession" {
			n++
		}
	}
	if n != 2 {
		t.Errorf("expected 2 KillSession calls for web+ expansion, got %d", n)
	}
}

func TestRunRestartWith_Group(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	mock := d.client.(*tmux.MockClient)
	mock.HasSessionReturn = true
	_ = d.store.Save("alpha", &config.ProjectConfig{
		Root:    d.global.ProjectDirs[0],
		Windows: []config.Window{{Name: "editor", Panes: []config.Pane{{Command: "vim"}}}},
	})
	_ = d.store.Save("bravo", &config.ProjectConfig{
		Root:    d.global.ProjectDirs[0],
		Windows: []config.Window{{Name: "editor", Panes: []config.Pane{{Command: "vim"}}}},
	})
	d.global.Groups = map[string][]string{"batch": {"alpha", "bravo"}}

	if err := runRestartWith(d, []string{"@batch"}); err != nil {
		t.Fatalf("runRestartWith: %v", err)
	}

	n := 0
	for _, c := range mock.Calls {
		if c.Method == "KillSession" {
			n++
		}
	}
	if n != 2 {
		t.Errorf("expected 2 KillSession calls for @batch, got %d", n)
	}
}

func TestRunRestartWith_Glob_AttachesLast(t *testing.T) {
	d := testDeps(t)
	mock := d.client.(*tmux.MockClient)
	mock.HasSessionReturn = true
	_ = d.store.Save("web-api", &config.ProjectConfig{
		Root: d.global.ProjectDirs[0],
		Windows: []config.Window{
			{Name: "editor", Panes: []config.Pane{{Command: "vim"}}},
		},
	})
	_ = d.store.Save("web-ui", &config.ProjectConfig{
		Root: d.global.ProjectDirs[0],
		Windows: []config.Window{
			{Name: "editor", Panes: []config.Pane{{Command: "vim"}}},
		},
	})

	if err := runRestartWith(d, []string{"web+"}); err != nil {
		t.Fatalf("runRestartWith: %v", err)
	}

	var lastAttach string
	for _, c := range mock.Calls {
		if c.Method == "AttachSession" && len(c.Args) > 0 {
			lastAttach = c.Args[0]
		}
	}
	if lastAttach != "web-ui" {
		t.Errorf("expected final AttachSession for web-ui (sorted glob order), got %q", lastAttach)
	}
}

func TestRunRestartWith_MultiWindow(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	_ = d.store.Save("blog", &config.ProjectConfig{
		Root: d.global.ProjectDirs[0],
		Windows: []config.Window{
			{Name: "editor", Panes: []config.Pane{{Command: "vim"}}},
			{Name: "server", Panes: []config.Pane{{Command: "go run ."}}},
		},
	})

	if err := runRestartWith(d, []string{"blog:editor,server"}); err != nil {
		t.Fatalf("runRestartWith: %v", err)
	}

	mock := d.client.(*tmux.MockClient)
	n := 0
	for _, c := range mock.Calls {
		if c.Method == "KillWindow" {
			n++
		}
	}
	if n != 2 {
		t.Errorf("expected 2 KillWindow calls, got %d", n)
	}
}
