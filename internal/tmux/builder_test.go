package tmux

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Drew-Daniels/nux/internal/config"
	"github.com/Drew-Daniels/nux/internal/resolver"
)

func newTestBuilder(client *MockClient, global *config.GlobalConfig) *Builder {
	if global == nil {
		global = &config.GlobalConfig{}
	}
	return NewBuilder(client, global)
}

func TestBuild_WindowsWithPanesAndLayout(t *testing.T) {
	mock := &MockClient{}
	global := &config.GlobalConfig{
		PaneInit: []string{"source ~/.bashrc"},
	}
	builder := newTestBuilder(mock, global)

	cfg := &config.ProjectConfig{
		Env: map[string]string{"APP_ENV": "dev"},
		Windows: []config.Window{
			{
				Name:   "editor",
				Layout: "main-vertical",
				Panes: []config.Pane{
					{Command: "vim"},
					{Command: "make watch"},
				},
			},
			{
				Name: "server",
				Root: "backend",
				Panes: []config.Pane{
					{Command: "go run ."},
				},
			},
		},
		OnStart: []string{"echo starting"},
		OnStop:  []string{"echo bye"},
		OnReady: []string{"echo attached"},
	}

	err := builder.Build("myproj", cfg, "/home/user/myproj")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	assertCalled(t, mock, "NewSession")
	assertCalled(t, mock, "SetEnv")
	assertCalled(t, mock, "SetHook")
	assertCalled(t, mock, "SplitWindow")
	assertCalled(t, mock, "SelectLayout")
	assertCalled(t, mock, "SelectPane")
	assertCalled(t, mock, "SelectWindow")

	assertCalledWith(t, mock, "NewSession", "myproj")
	assertCalledWith(t, mock, "SetEnv", "APP_ENV")
	assertCalledWith(t, mock, "SelectLayout", "main-vertical")

	sendKeysTargets := callsFor(mock, "SendKeys")
	if len(sendKeysTargets) == 0 {
		t.Fatal("expected SendKeys calls")
	}

	foundVim := false
	foundInit := false
	for _, c := range sendKeysTargets {
		if len(c.Args) >= 2 && c.Args[1] == "vim" {
			foundVim = true
		}
		if len(c.Args) >= 2 && c.Args[1] == "source ~/.bashrc" {
			foundInit = true
		}
	}
	if !foundVim {
		t.Error("expected SendKeys with 'vim'")
	}
	if !foundInit {
		t.Error("expected SendKeys with pane_init command")
	}
}

func TestBuildWindows_UserOrder(t *testing.T) {
	mock := &MockClient{}
	global := &config.GlobalConfig{}
	builder := newTestBuilder(mock, global)

	cfg := &config.ProjectConfig{
		Windows: []config.Window{
			{Name: "editor", Panes: []config.Pane{{Command: "vim"}}},
			{Name: "server", Panes: []config.Pane{{Command: "go run ."}}},
		},
		OnStart: []string{"echo hi"},
		OnReady: []string{"echo ready"},
	}

	err := builder.BuildWindows("myproj", cfg, "/home/user/myproj", []string{"server", "editor"})
	if err != nil {
		t.Fatalf("BuildWindows: %v", err)
	}

	var ns *NewSessionOpts
	for _, c := range mock.Calls {
		if c.Method == "NewSession" && c.Opts != nil {
			if o, ok := c.Opts.(NewSessionOpts); ok {
				ns = &o
				break
			}
		}
	}
	if ns == nil || ns.Window != "server" {
		t.Errorf("first window should be server (user order), got %+v", ns)
	}

	assertCalledWith(t, mock, "SendKeys", "myproj:server")
}

func TestBuildWindows_UnknownWindow(t *testing.T) {
	mock := &MockClient{}
	builder := newTestBuilder(mock, nil)
	cfg := &config.ProjectConfig{
		Windows: []config.Window{{Name: "a", Panes: []config.Pane{{Command: "x"}}}},
	}
	err := builder.BuildWindows("p", cfg, "/r", []string{"missing"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBuild_NilConfig_DefaultSessionWithWindows(t *testing.T) {
	mock := &MockClient{}
	global := &config.GlobalConfig{
		DefaultSession: &config.DefaultSession{
			Windows: []config.Window{
				{Name: "main", Panes: []config.Pane{{Command: "htop"}}},
			},
		},
		PaneInit: []string{"export TERM=xterm"},
	}
	builder := newTestBuilder(mock, global)

	err := builder.Build("scratch", nil, "/tmp/scratch")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	assertCalled(t, mock, "NewSession")
	assertCalledWith(t, mock, "SendKeys", "htop")
	assertCalledWith(t, mock, "SendKeys", "export TERM=xterm")
}

func TestBuild_NilConfig_DefaultWindows(t *testing.T) {
	mock := &MockClient{}
	global := &config.GlobalConfig{
		DefaultSession: &config.DefaultSession{
			Windows: []config.Window{
				{
					Name:   "editor",
					Layout: "tiled",
					Panes: []config.Pane{
						{Command: "nvim"},
						{},
					},
				},
				{Name: "shell", Panes: []config.Pane{{Command: ""}}},
			},
		},
		PaneInit: []string{"cls"},
	}
	builder := newTestBuilder(mock, global)

	err := builder.Build("project", nil, "/tmp/project")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	assertCalled(t, mock, "NewSession")
	assertCalled(t, mock, "SplitWindow")
	assertCalledWith(t, mock, "SelectLayout", "tiled")
	assertCalledWith(t, mock, "SendKeys", "nvim")
	assertCalledWith(t, mock, "SendKeys", "cls")
	assertCalled(t, mock, "NewWindow")
	assertCalled(t, mock, "SelectWindow")
}

func TestBuild_ProjectPaneInit_AfterGlobal(t *testing.T) {
	mock := &MockClient{}
	global := &config.GlobalConfig{PaneInit: []string{"global-init"}}
	builder := newTestBuilder(mock, global)
	cfg := &config.ProjectConfig{
		PaneInit: []string{"project-init"},
		Windows:  []config.Window{{Name: "w", Panes: []config.Pane{{Command: "vim"}}}},
	}
	if err := builder.Build("proj", cfg, "/tmp/p"); err != nil {
		t.Fatalf("Build: %v", err)
	}
	var paneKeys []string
	for _, c := range mock.Calls {
		if c.Method != "SendKeys" || len(c.Args) < 2 {
			continue
		}
		if strings.HasPrefix(c.Args[0], "proj:w.") {
			paneKeys = append(paneKeys, c.Args[1])
		}
	}
	if len(paneKeys) < 3 || paneKeys[0] != "global-init" || paneKeys[1] != "project-init" || paneKeys[2] != "vim" {
		t.Fatalf("first pane SendKeys want [global-init project-init vim], got %v", paneKeys)
	}
}

func TestBuild_NilConfig_NoDefaults(t *testing.T) {
	mock := &MockClient{}
	builder := newTestBuilder(mock, nil)

	err := builder.Build("bare", nil, "/tmp/bare")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	assertCalled(t, mock, "NewSession")
	if mock.Called("SendKeys") {
		t.Error("bare session with no defaults should not send keys")
	}
}

func TestBuild_RunCommand_NilConfig(t *testing.T) {
	mock := &MockClient{}
	builder := newTestBuilder(mock, nil)
	builder.SetAdHocLayout(&AdHocLayout{Command: "go test ./..."})

	err := builder.Build("run-test", nil, "/home/user/project")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	assertCalled(t, mock, "NewSession")
	assertCalledWith(t, mock, "SendKeys", "go test ./...")

	if mock.Called("SetEnv") || mock.Called("SetHook") {
		t.Error("bare session with --run should not set env or hooks")
	}
	if mock.Called("SplitWindow") {
		t.Error("--run without --panes should not split windows")
	}
}

func TestBuild_RunCommand_OverridesDefaultSessionWindows(t *testing.T) {
	mock := &MockClient{}
	global := &config.GlobalConfig{
		DefaultSession: &config.DefaultSession{
			Windows: []config.Window{
				{Name: "editor", Panes: []config.Pane{{Command: "nvim"}}},
			},
		},
	}
	builder := newTestBuilder(mock, global)
	builder.SetAdHocLayout(&AdHocLayout{Command: "fish"})

	err := builder.Build("proj", nil, "/tmp/proj")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	assertCalledWith(t, mock, "SendKeys", "fish")

	if mock.Called("NewWindow") {
		t.Error("--run should skip default_session windows")
	}
	for _, c := range callsFor(mock, "SendKeys") {
		if len(c.Args) >= 2 && c.Args[1] == "nvim" {
			t.Error("--run should not send default_session window commands")
		}
	}
}

func TestBuild_RunCommand_WithAdHocLayout_OverridesDefaultSessionWindows(t *testing.T) {
	mock := &MockClient{}
	global := &config.GlobalConfig{
		DefaultSession: &config.DefaultSession{
			Windows: []config.Window{
				{Name: "editor", Panes: []config.Pane{{Command: "nvim"}}},
			},
		},
	}
	builder := newTestBuilder(mock, global)
	builder.SetAdHocLayout(&AdHocLayout{Layout: "tiled", Panes: 3, Command: "fish"})

	err := builder.Build("proj", nil, "/tmp/proj")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	splits := callsFor(mock, "SplitWindow")
	if len(splits) != 2 {
		t.Fatalf("expected 2 SplitWindow calls for 3 panes, got %d", len(splits))
	}
	assertCalledWith(t, mock, "SelectLayout", "tiled")

	cmdCalls := 0
	for _, c := range callsFor(mock, "SendKeys") {
		if len(c.Args) >= 2 && c.Args[1] == "fish" {
			cmdCalls++
		}
	}
	if cmdCalls != 3 {
		t.Errorf("expected fish sent to 3 panes, got %d", cmdCalls)
	}

	if mock.Called("NewWindow") {
		t.Error("--run with --layout should skip default_session windows")
	}
}

func TestStopSession(t *testing.T) {
	mock := &MockClient{HasSessionReturn: true}
	builder := newTestBuilder(mock, nil)

	err := builder.StopSession("myproj")
	if err != nil {
		t.Fatalf("StopSession returned error: %v", err)
	}

	assertCalledWith(t, mock, "KillSession", "myproj")
}

func TestStopSession_NotRunning(t *testing.T) {
	mock := &MockClient{HasSessionReturn: false}
	builder := newTestBuilder(mock, nil)

	err := builder.StopSession("myproj")
	if err != nil {
		t.Fatalf("StopSession returned error: %v", err)
	}

	if mock.Called("KillSession") {
		t.Error("should not call KillSession when session is not running")
	}
}

func TestStopAll(t *testing.T) {
	mock := &MockClient{
		ListSessionsReturn: []SessionInfo{
			{Name: "a"},
			{Name: "b"},
		},
	}
	builder := newTestBuilder(mock, nil)

	err := builder.StopAll(StopAllOpts{})
	if err != nil {
		t.Fatalf("StopAll returned error: %v", err)
	}

	kills := callsFor(mock, "KillSession")
	if len(kills) != 2 {
		t.Fatalf("expected 2 KillSession calls, got %d", len(kills))
	}
}

func TestStopAll_OnSessionOrderAndIndices(t *testing.T) {
	mock := &MockClient{
		ListSessionsReturn: []SessionInfo{
			{Name: "a"},
			{Name: "b"},
		},
	}
	builder := newTestBuilder(mock, nil)

	var seen []string
	err := builder.StopAll(StopAllOpts{
		OnSession: func(name string, index, total int) {
			seen = append(seen, fmt.Sprintf("%s:%d/%d", name, index, total))
		},
	})
	if err != nil {
		t.Fatalf("StopAll: %v", err)
	}
	want := []string{"a:1/2", "b:2/2"}
	if len(seen) != len(want) {
		t.Fatalf("OnSession calls = %d, want %d: %v", len(seen), len(want), seen)
	}
	for i := range want {
		if seen[i] != want[i] {
			t.Errorf("OnSession[%d] = %q, want %q", i, seen[i], want[i])
		}
	}
}

func TestStopAll_OnEmpty(t *testing.T) {
	mock := &MockClient{ListSessionsReturn: []SessionInfo{}}
	builder := newTestBuilder(mock, nil)

	var emptyCalls int
	err := builder.StopAll(StopAllOpts{
		OnEmpty: func() { emptyCalls++ },
	})
	if err != nil {
		t.Fatalf("StopAll: %v", err)
	}
	if emptyCalls != 1 {
		t.Errorf("OnEmpty calls = %d, want 1", emptyCalls)
	}
	if mock.Called("KillSession") {
		t.Error("should not KillSession when list is empty")
	}
}

func TestRestartWindow(t *testing.T) {
	mock := &MockClient{}
	builder := newTestBuilder(mock, nil)

	cfg := &config.ProjectConfig{
		Windows: []config.Window{
			{Name: "editor", Panes: []config.Pane{{Command: "vim"}}},
			{Name: "server", Panes: []config.Pane{{Command: "go run ."}}},
		},
	}

	err := builder.RestartWindow("myproj", "server", cfg, "/home/user/myproj")
	if err != nil {
		t.Fatalf("RestartWindow returned error: %v", err)
	}

	assertCalledWith(t, mock, "KillWindow", "server")
	assertCalledWith(t, mock, "NewWindow", "server")
	assertCalledWith(t, mock, "SendKeys", "go run .")
}

func TestRestartWindow_NotFound(t *testing.T) {
	mock := &MockClient{}
	builder := newTestBuilder(mock, nil)

	cfg := &config.ProjectConfig{
		Windows: []config.Window{
			{Name: "editor", Panes: []config.Pane{{Command: ""}}},
		},
	}

	err := builder.RestartWindow("myproj", "missing", cfg, "/root")
	if err == nil {
		t.Fatal("expected error for missing window")
	}
}

func TestBuild_DefaultShell(t *testing.T) {
	mock := &MockClient{}
	global := &config.GlobalConfig{
		DefaultShell: "/usr/bin/fish",
	}
	builder := newTestBuilder(mock, global)

	cfg := &config.ProjectConfig{
		Windows: []config.Window{
			{Name: "main", Panes: []config.Pane{{Command: "echo hi"}}},
		},
	}

	err := builder.Build("test", cfg, "/tmp")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	assertCalledWith(t, mock, "SetOption", "default-command")
	found := false
	for _, c := range callsFor(mock, "SetOption") {
		if len(c.Args) >= 3 && c.Args[2] == "/usr/bin/fish" {
			found = true
		}
	}
	if !found {
		t.Error("expected SetOption with /usr/bin/fish")
	}
}

func TestBuild_ProjectDefaultShell_OverridesGlobal(t *testing.T) {
	mock := &MockClient{}
	global := &config.GlobalConfig{DefaultShell: "/usr/bin/fish"}
	builder := newTestBuilder(mock, global)
	cfg := &config.ProjectConfig{
		DefaultShell: "/bin/bash",
		Windows:      []config.Window{{Name: "main", Panes: []config.Pane{{Command: "echo"}}}},
	}
	if err := builder.Build("p", cfg, "/tmp"); err != nil {
		t.Fatalf("Build: %v", err)
	}
	found := false
	for _, c := range callsFor(mock, "SetOption") {
		if len(c.Args) >= 3 && c.Args[1] == "default-command" && c.Args[2] == "/bin/bash" {
			found = true
		}
	}
	if !found {
		t.Error("expected project default_shell to override global (want /bin/bash)")
	}
}

func TestBuild_ProjectDefaultShell_WhenGlobalUnset(t *testing.T) {
	mock := &MockClient{}
	builder := newTestBuilder(mock, &config.GlobalConfig{})
	cfg := &config.ProjectConfig{
		DefaultShell: "/usr/local/bin/fish",
		Windows:      []config.Window{{Name: "main", Panes: []config.Pane{{Command: "echo"}}}},
	}
	if err := builder.Build("p", cfg, "/tmp"); err != nil {
		t.Fatalf("Build: %v", err)
	}
	found := false
	for _, c := range callsFor(mock, "SetOption") {
		if len(c.Args) >= 3 && c.Args[2] == "/usr/local/bin/fish" {
			found = true
		}
	}
	if !found {
		t.Error("expected SetOption with project default_shell when global is empty")
	}
}

func TestBuild_WindowRootRelative(t *testing.T) {
	mock := &MockClient{}
	builder := newTestBuilder(mock, nil)

	cfg := &config.ProjectConfig{
		Windows: []config.Window{
			{
				Name: "api",
				Root: "services/api",
				Panes: []config.Pane{
					{Command: "go run ."},
				},
			},
		},
	}

	err := builder.Build("proj", cfg, "/home/user/proj")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	sessions := callsFor(mock, "NewSession")
	if len(sessions) == 0 {
		t.Fatal("expected NewSession call")
	}
	if sessions[0].Args[1] != "/home/user/proj/services/api" {
		t.Errorf("expected root /home/user/proj/services/api, got %s", sessions[0].Args[1])
	}
}

func TestBuild_AdHocLayout_NilConfig(t *testing.T) {
	mock := &MockClient{}
	builder := newTestBuilder(mock, nil)
	builder.SetAdHocLayout(&AdHocLayout{Layout: "tiled", Panes: 4})

	err := builder.Build("scratch", nil, "/tmp/scratch")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	assertCalled(t, mock, "NewSession")

	splits := callsFor(mock, "SplitWindow")
	if len(splits) != 3 {
		t.Fatalf("expected 3 SplitWindow calls for 4 panes, got %d", len(splits))
	}

	assertCalledWith(t, mock, "SelectLayout", "tiled")
	assertCalled(t, mock, "SelectPane")
}

func TestBuild_AdHocLayout_NilConfig_WithPaneInit(t *testing.T) {
	mock := &MockClient{}
	global := &config.GlobalConfig{
		PaneInit: []string{"clear"},
	}
	builder := newTestBuilder(mock, global)
	builder.SetAdHocLayout(&AdHocLayout{Layout: "even-horizontal", Panes: 2})

	err := builder.Build("proj", nil, "/tmp/proj")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	splits := callsFor(mock, "SplitWindow")
	if len(splits) != 1 {
		t.Fatalf("expected 1 SplitWindow call for 2 panes, got %d", len(splits))
	}

	assertCalledWith(t, mock, "SelectLayout", "even-horizontal")

	initCalls := 0
	for _, c := range callsFor(mock, "SendKeys") {
		if len(c.Args) >= 2 && c.Args[1] == "clear" {
			initCalls++
		}
	}
	if initCalls != 2 {
		t.Errorf("expected pane_init sent to 2 panes, got %d", initCalls)
	}
}

func TestBuild_AdHocLayout_WithDefaultWindows(t *testing.T) {
	mock := &MockClient{}
	global := &config.GlobalConfig{
		DefaultSession: &config.DefaultSession{
			Windows: []config.Window{
				{Name: "main", Panes: []config.Pane{{Command: "htop"}}},
			},
		},
	}
	builder := newTestBuilder(mock, global)
	builder.SetAdHocLayout(&AdHocLayout{Layout: "tiled", Panes: 2})

	err := builder.Build("scratch", nil, "/tmp/scratch")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	splits := callsFor(mock, "SplitWindow")
	if len(splits) != 1 {
		t.Fatalf("expected 1 SplitWindow, got %d", len(splits))
	}

	assertCalledWith(t, mock, "SelectLayout", "tiled")
}

func TestBuild_AdHocLayout_DoesNotOverrideConfigWindows(t *testing.T) {
	mock := &MockClient{}
	builder := newTestBuilder(mock, nil)
	builder.SetAdHocLayout(&AdHocLayout{Layout: "tiled", Panes: 4})

	cfg := &config.ProjectConfig{
		Windows: []config.Window{
			{
				Name:   "editor",
				Layout: "main-vertical",
				Panes:  []config.Pane{{Command: "vim"}},
			},
		},
	}

	err := builder.Build("proj", cfg, "/tmp/proj")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	layouts := callsFor(mock, "SelectLayout")
	if len(layouts) != 1 {
		t.Fatalf("expected 1 SelectLayout call, got %d", len(layouts))
	}
	if layouts[0].Args[2] != "main-vertical" {
		t.Errorf("expected config layout main-vertical to be preserved, got %s", layouts[0].Args[2])
	}
}

func TestBuild_AdHocLayout_FillsEmptyWindowLayout(t *testing.T) {
	mock := &MockClient{}
	builder := newTestBuilder(mock, nil)
	builder.SetAdHocLayout(&AdHocLayout{Layout: "tiled", Panes: 4})

	cfg := &config.ProjectConfig{
		Windows: []config.Window{
			{
				Name:  "editor",
				Panes: []config.Pane{{Command: "vim"}},
			},
		},
	}

	err := builder.Build("proj", cfg, "/tmp/proj")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	assertCalledWith(t, mock, "SelectLayout", "tiled")
}

func TestBuild_RunCommand_AdHocLayout(t *testing.T) {
	mock := &MockClient{}
	builder := newTestBuilder(mock, nil)
	builder.SetAdHocLayout(&AdHocLayout{Layout: "even-vertical", Panes: 3, Command: "go test ./..."})

	err := builder.Build("run-test", nil, "/home/user/project")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	assertCalled(t, mock, "NewSession")

	splits := callsFor(mock, "SplitWindow")
	if len(splits) != 2 {
		t.Fatalf("expected 2 SplitWindow calls for 3 panes, got %d", len(splits))
	}

	assertCalledWith(t, mock, "SelectLayout", "even-vertical")
	assertCalled(t, mock, "SelectPane")

	cmdCalls := 0
	for _, c := range callsFor(mock, "SendKeys") {
		if len(c.Args) >= 2 && c.Args[1] == "go test ./..." {
			cmdCalls++
		}
	}
	if cmdCalls != 3 {
		t.Errorf("expected command sent to 3 panes, got %d", cmdCalls)
	}
}

func TestBuild_RunCommand_AdHocLayout_OverridesDefaultWindows(t *testing.T) {
	mock := &MockClient{}
	global := &config.GlobalConfig{
		DefaultSession: &config.DefaultSession{
			Windows: []config.Window{
				{Name: "main", Panes: []config.Pane{{Command: "htop"}}},
			},
		},
	}
	builder := newTestBuilder(mock, global)
	builder.SetAdHocLayout(&AdHocLayout{Layout: "tiled", Panes: 2, Command: "fish"})

	err := builder.Build("scratch", nil, "/tmp/scratch")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	for _, c := range callsFor(mock, "SendKeys") {
		if len(c.Args) >= 2 && c.Args[1] == "htop" {
			t.Error("--run command should override default_session windows")
		}
	}
	assertCalledWith(t, mock, "SendKeys", "fish")
}

func TestBuild_AdHocLayout_BaseIndex1(t *testing.T) {
	mock := &MockClient{BaseIndexReturn: 1, PaneBaseIndexReturn: 1}
	builder := newTestBuilder(mock, nil)
	builder.SetAdHocLayout(&AdHocLayout{Layout: "tiled", Panes: 4})

	err := builder.Build("scratch", nil, "/tmp/scratch")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	for _, c := range callsFor(mock, "SplitWindow") {
		if c.Args[1] != "1" {
			t.Errorf("SplitWindow should target window 1, got %q", c.Args[1])
		}
	}
	for _, c := range callsFor(mock, "SelectLayout") {
		if c.Args[1] != "1" {
			t.Errorf("SelectLayout should target window 1, got %q", c.Args[1])
		}
	}
	for _, c := range callsFor(mock, "SelectPane") {
		if c.Args[1] != "1" {
			t.Errorf("SelectPane should target window 1, got %q", c.Args[1])
		}
		if c.Args[2] != "1" {
			t.Errorf("SelectPane should target pane 1, got %q", c.Args[2])
		}
	}
	for _, c := range callsFor(mock, "SendKeys") {
		if !strings.HasPrefix(c.Args[0], "scratch:1.1") {
			t.Errorf("SendKeys target should start with scratch:1.1, got %q", c.Args[0])
		}
	}
}

func TestBuild_RunCommand_BaseIndex1(t *testing.T) {
	mock := &MockClient{BaseIndexReturn: 1, PaneBaseIndexReturn: 1}
	builder := newTestBuilder(mock, nil)
	builder.SetAdHocLayout(&AdHocLayout{Layout: "tiled", Panes: 2, Command: "make test"})

	err := builder.Build("run", nil, "/tmp")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	for _, c := range callsFor(mock, "SplitWindow") {
		if c.Args[1] != "1" {
			t.Errorf("SplitWindow should target window 1, got %q", c.Args[1])
		}
	}
	for _, c := range callsFor(mock, "SendKeys") {
		if !strings.HasPrefix(c.Args[0], "run:1") {
			t.Errorf("SendKeys target should start with run:1, got %q", c.Args[0])
		}
	}
}

func TestBuild_Windowed_PaneBaseIndex1(t *testing.T) {
	mock := &MockClient{PaneBaseIndexReturn: 1}
	builder := newTestBuilder(mock, nil)

	cfg := &config.ProjectConfig{
		Windows: []config.Window{
			{
				Name: "editor",
				Panes: []config.Pane{
					{Command: "vim"},
					{Command: "make watch"},
				},
			},
		},
	}

	err := builder.Build("proj", cfg, "/tmp/proj")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	sendKeys := callsFor(mock, "SendKeys")
	for _, c := range sendKeys {
		target := c.Args[0]
		if !strings.Contains(target, ".") {
			continue
		}
		parts := strings.SplitN(target, ".", 2)
		paneIdx := parts[1]
		if paneIdx == "0" {
			t.Errorf("SendKeys targets pane 0 but pane-base-index is 1: %q", target)
		}
	}

	selectPanes := callsFor(mock, "SelectPane")
	for _, c := range selectPanes {
		if c.Args[2] != "1" {
			t.Errorf("SelectPane should target pane 1, got %q", c.Args[2])
		}
	}
}

func TestBuild_WindowEnv(t *testing.T) {
	mock := &MockClient{}
	builder := newTestBuilder(mock, nil)

	cfg := &config.ProjectConfig{
		Windows: []config.Window{
			{
				Name: "api",
				Env:  map[string]string{"PORT": "3000", "DEBUG": "true"},
				Panes: []config.Pane{
					{Command: "go run ."},
					{Command: "make watch"},
				},
			},
		},
	}

	err := builder.Build("proj", cfg, "/tmp/proj")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	sends := callsFor(mock, "SendKeys")
	var exports []string
	for _, c := range sends {
		if len(c.Args) >= 2 && strings.HasPrefix(c.Args[1], "export ") {
			exports = append(exports, c.Args[1])
		}
	}

	if len(exports) != 4 {
		t.Fatalf("expected 4 export commands (2 vars x 2 panes), got %d: %v", len(exports), exports)
	}

	assertCalledWith(t, mock, "SendKeys", "export DEBUG='true'")
	assertCalledWith(t, mock, "SendKeys", "export PORT='3000'")
}

func TestBuild_WindowEnv_MultipleWindows(t *testing.T) {
	mock := &MockClient{}
	builder := newTestBuilder(mock, nil)

	cfg := &config.ProjectConfig{
		Windows: []config.Window{
			{
				Name:  "api",
				Env:   map[string]string{"PORT": "3000"},
				Panes: []config.Pane{{Command: "go run ."}},
			},
			{
				Name:  "frontend",
				Env:   map[string]string{"PORT": "5173"},
				Panes: []config.Pane{{Command: "npm run dev"}},
			},
		},
	}

	err := builder.Build("proj", cfg, "/tmp/proj")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	sends := callsFor(mock, "SendKeys")

	var apiExport, feExport bool
	for _, c := range sends {
		if len(c.Args) < 2 {
			continue
		}
		if strings.Contains(c.Args[0], ":api") && c.Args[1] == "export PORT='3000'" {
			apiExport = true
		}
		if strings.Contains(c.Args[0], ":frontend") && c.Args[1] == "export PORT='5173'" {
			feExport = true
		}
	}

	if !apiExport {
		t.Error("expected export PORT='3000' sent to api window")
	}
	if !feExport {
		t.Error("expected export PORT='5173' sent to frontend window")
	}
}

func TestBuild_WindowEnv_SinglePane(t *testing.T) {
	mock := &MockClient{}
	builder := newTestBuilder(mock, nil)

	cfg := &config.ProjectConfig{
		Windows: []config.Window{
			{
				Name:  "server",
				Env:   map[string]string{"PORT": "8080"},
				Panes: []config.Pane{{Command: "npm start"}},
			},
		},
	}

	err := builder.Build("proj", cfg, "/tmp/proj")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	sends := callsFor(mock, "SendKeys")
	exportIdx := -1
	cmdIdx := -1
	for i, c := range sends {
		if len(c.Args) >= 2 && c.Args[1] == "export PORT='8080'" {
			exportIdx = i
		}
		if len(c.Args) >= 2 && c.Args[1] == "npm start" {
			cmdIdx = i
		}
	}

	if exportIdx == -1 {
		t.Fatal("expected export PORT='8080'")
	}
	if cmdIdx == -1 {
		t.Fatal("expected npm start")
	}
	if exportIdx >= cmdIdx {
		t.Error("export should come before the pane command")
	}
}

func TestBuild_WindowEnv_WithProjectEnv(t *testing.T) {
	mock := &MockClient{}
	builder := newTestBuilder(mock, nil)

	cfg := &config.ProjectConfig{
		Env: map[string]string{"NODE_ENV": "development"},
		Windows: []config.Window{
			{
				Name:  "api",
				Env:   map[string]string{"PORT": "3000"},
				Panes: []config.Pane{{Command: "npm start"}},
			},
		},
	}

	err := builder.Build("proj", cfg, "/tmp/proj")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	assertCalledWith(t, mock, "SetEnv", "NODE_ENV")
	assertCalledWith(t, mock, "SendKeys", "export PORT='3000'")

	if !mock.Called("SetEnv") {
		t.Error("project-level env should use SetEnv")
	}
}

func TestBuild_WindowEnv_ShellEscaped(t *testing.T) {
	mock := &MockClient{}
	builder := newTestBuilder(mock, nil)

	cfg := &config.ProjectConfig{
		Windows: []config.Window{
			{
				Name: "app",
				Env: map[string]string{
					"MSG":    "hello world",
					"QUOTES": "it's fine",
					"DOLLAR": "price is $5",
				},
				Panes: []config.Pane{{Command: "node server.js"}},
			},
		},
	}

	err := builder.Build("test", cfg, "/tmp/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertCalledWith(t, mock, "SendKeys", "export DOLLAR='price is $5'")
	assertCalledWith(t, mock, "SendKeys", "export MSG='hello world'")
	assertCalledWith(t, mock, "SendKeys", "export QUOTES='it'\\''s fine'")
}

func TestRestartSession(t *testing.T) {
	mock := &MockClient{HasSessionReturn: true}
	builder := newTestBuilder(mock, nil)

	cfg := &config.ProjectConfig{
		Windows: []config.Window{
			{Name: "editor", Panes: []config.Pane{{Command: "vim"}}},
		},
	}

	err := builder.RestartSession("myproj", cfg, "/home/user/myproj")
	if err != nil {
		t.Fatalf("RestartSession returned error: %v", err)
	}

	assertCalledWith(t, mock, "KillSession", "myproj")
	assertCalled(t, mock, "NewSession")
	assertCalledWith(t, mock, "SendKeys", "vim")
}

func TestBuild_OnDetachHooks(t *testing.T) {
	mock := &MockClient{}
	builder := newTestBuilder(mock, nil)

	cfg := &config.ProjectConfig{
		OnDetach: []string{"echo detached"},
		Windows:  []config.Window{{Name: "main", Panes: []config.Pane{{Command: "echo hi"}}}},
	}

	err := builder.Build("proj", cfg, "/tmp/proj")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	assertCalledWith(t, mock, "SetHook", "client-detached[0]")
}

func TestBuild_ErrorFromCreateSession(t *testing.T) {
	mock := &MockClient{DefaultError: fmt.Errorf("tmux not found")}
	builder := newTestBuilder(mock, nil)

	err := builder.Build("proj", nil, "/tmp/proj")
	if err == nil {
		t.Fatal("expected error when session creation fails")
	}
}

func TestBuild_WindowRoot_Absolute(t *testing.T) {
	mock := &MockClient{}
	builder := newTestBuilder(mock, nil)

	cfg := &config.ProjectConfig{
		Windows: []config.Window{
			{Name: "api", Root: "/opt/api", Panes: []config.Pane{{Command: "go run ."}}},
		},
	}

	err := builder.Build("proj", cfg, "/home/user/proj")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	sessions := callsFor(mock, "NewSession")
	if len(sessions) == 0 {
		t.Fatal("expected NewSession call")
	}
	if sessions[0].Args[1] != "/opt/api" {
		t.Errorf("expected root /opt/api, got %s", sessions[0].Args[1])
	}
}

func TestBuild_WindowRoot_Tilde(t *testing.T) {
	mock := &MockClient{}
	builder := newTestBuilder(mock, nil)

	cfg := &config.ProjectConfig{
		Windows: []config.Window{
			{Name: "api", Root: "~/code/api", Panes: []config.Pane{{Command: "go run ."}}},
		},
	}

	err := builder.Build("proj", cfg, "/home/user/proj")
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	sessions := callsFor(mock, "NewSession")
	if len(sessions) == 0 {
		t.Fatal("expected NewSession call")
	}
	want := resolver.ResolveRoot("~/code/api", "/home/user/proj")
	if sessions[0].Args[1] != want {
		t.Errorf("expected expanded tilde root %q, got %q", want, sessions[0].Args[1])
	}
}

func TestFindWindow_NilConfig(t *testing.T) {
	_, ok := findWindow(nil, "editor")
	if ok {
		t.Error("expected false for nil config")
	}
}

func TestWindowRoot_Empty(t *testing.T) {
	got := windowRoot("", "/home/user/proj")
	want := resolver.ResolveRoot("", "/home/user/proj")
	if got != want {
		t.Errorf("windowRoot('', ...) = %q, want %q", got, want)
	}
}

// --- helpers ---

func assertCalled(t *testing.T, mock *MockClient, method string) {
	t.Helper()
	if !mock.Called(method) {
		t.Errorf("expected %s to be called", method)
	}
}

func assertCalledWith(t *testing.T, mock *MockClient, method, argSubstr string) {
	t.Helper()
	for _, c := range mock.Calls {
		if c.Method != method {
			continue
		}
		for _, a := range c.Args {
			if a == argSubstr {
				return
			}
		}
	}
	t.Errorf("expected %s to be called with arg %q", method, argSubstr)
}

func callsFor(mock *MockClient, method string) []Call {
	var result []Call
	for _, c := range mock.Calls {
		if c.Method == method {
			result = append(result, c)
		}
	}
	return result
}
