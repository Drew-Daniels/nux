package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Drew-Daniels/nux/internal/config"
	"github.com/Drew-Daniels/nux/internal/picker"
	"github.com/Drew-Daniels/nux/internal/tmux"
)

func TestExpandArgs_PlainNames(t *testing.T) {
	d := testDeps(t)
	names, err := expandArgs(d, []string{"blog", "api"})
	if err != nil {
		t.Fatalf("expandArgs: %v", err)
	}
	if len(names) != 2 || names[0].Project != "blog" || names[1].Project != "api" ||
		names[0].Windows != nil || names[1].Windows != nil {
		t.Errorf("got %+v, want [blog api] without windows", names)
	}
}

func TestExpandArgs_Group(t *testing.T) {
	d := testDeps(t)
	d.global.Groups = map[string][]string{
		"work": {"alpha", "bravo"},
	}

	names, err := expandArgs(d, []string{"@work"})
	if err != nil {
		t.Fatalf("expandArgs: %v", err)
	}
	if len(names) != 2 || names[0].Project != "alpha" || names[1].Project != "bravo" {
		t.Errorf("got %+v, want [alpha bravo]", names)
	}
}

func TestExpandArgs_GroupNotFound(t *testing.T) {
	d := testDeps(t)
	_, err := expandArgs(d, []string{"@missing"})
	if err == nil {
		t.Fatal("expected error for missing group")
	}
}

func TestExpandArgs_ColonTarget(t *testing.T) {
	d := testDeps(t)
	names, err := expandArgs(d, []string{"blog:editor"})
	if err != nil {
		t.Fatalf("expandArgs: %v", err)
	}
	if len(names) != 1 || names[0].Project != "blog" || len(names[0].Windows) != 1 || names[0].Windows[0] != "editor" {
		t.Errorf("got %+v, want blog + [editor]", names)
	}
}

func TestExpandArgs_MultiWindowTarget(t *testing.T) {
	d := testDeps(t)
	names, err := expandArgs(d, []string{"blog: editor , server "})
	if err != nil {
		t.Fatalf("expandArgs: %v", err)
	}
	if len(names) != 1 || names[0].Project != "blog" || len(names[0].Windows) != 2 ||
		names[0].Windows[0] != "editor" || names[0].Windows[1] != "server" {
		t.Errorf("got %+v, want blog + [editor server]", names)
	}
}

func TestExpandArgs_Glob(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("web-api", &config.ProjectConfig{Command: "a"})
	_ = d.store.Save("web-ui", &config.ProjectConfig{Command: "b"})
	_ = d.store.Save("other", &config.ProjectConfig{Command: "c"})

	names, err := expandArgs(d, []string{"web+"})
	if err != nil {
		t.Fatalf("expandArgs: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("got %v, want 2 web-* matches", names)
	}
}

func TestRunSessions_Single(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	_ = d.store.Save("blog", &config.ProjectConfig{
		Root:    d.global.ProjectsDir,
		Command: "vim",
	})

	err := runSessions(d, []string{"blog"})
	if err != nil {
		t.Fatalf("runSessions: %v", err)
	}

	mock := d.client.(*tmux.MockClient)
	if !mock.Called("NewSession") {
		t.Error("expected NewSession to be called")
	}
}

func TestRunSessions_SkipsExisting(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	mock := d.client.(*tmux.MockClient)
	mock.HasSessionReturn = true
	_ = d.store.Save("blog", &config.ProjectConfig{Root: d.global.ProjectsDir, Command: "vim"})

	err := runSessions(d, []string{"blog"})
	if err != nil {
		t.Fatalf("runSessions: %v", err)
	}

	if mock.Called("NewSession") {
		t.Error("should skip NewSession for existing session")
	}
}

func TestRunSessions_WithVarOverrides(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	d.vars = map[string]string{"greeting": "hello"}
	_ = d.store.Save("api", &config.ProjectConfig{
		Root:    d.global.ProjectsDir,
		Command: "echo {{greeting}}",
		Vars:    map[string]string{"greeting": "hi"},
	})

	err := runSessions(d, []string{"api"})
	if err != nil {
		t.Fatalf("runSessions: %v", err)
	}

	mock := d.client.(*tmux.MockClient)
	if !mock.Called("NewSession") {
		t.Error("expected session to be built")
	}
}

func TestRunSessions_WithRunCommand(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	d.run = "go test ./..."
	d.builder.SetAdHocLayout(&tmux.AdHocLayout{Command: "go test ./..."})

	blogDir := filepath.Join(d.global.ProjectsDir, "blog")
	if err := os.Mkdir(blogDir, 0o755); err != nil {
		t.Fatal(err)
	}

	err := runSessions(d, []string{"blog"})
	if err != nil {
		t.Fatalf("runSessions: %v", err)
	}

	mock := d.client.(*tmux.MockClient)
	if !mock.Called("NewSession") {
		t.Error("expected NewSession")
	}
	found := false
	for _, c := range mock.Calls {
		if c.Method == "SendKeys" && len(c.Args) >= 2 && c.Args[1] == "go test ./..." {
			found = true
		}
	}
	if !found {
		t.Error("expected SendKeys with run command")
	}
}

func TestRunSessions_RunCommand_SkipsProjectConfig(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	d.run = "fish"
	d.builder.SetAdHocLayout(&tmux.AdHocLayout{Command: "fish"})

	_ = d.store.Save("blog", &config.ProjectConfig{
		Root: d.global.ProjectsDir,
		Windows: []config.Window{
			{Name: "editor", Panes: []config.Pane{{Command: "nvim"}}},
			{Name: "server", Panes: []config.Pane{{Command: "go run ."}}},
		},
	})

	err := runSessions(d, []string{"blog"})
	if err != nil {
		t.Fatalf("runSessions: %v", err)
	}

	mock := d.client.(*tmux.MockClient)
	if !mock.Called("NewSession") {
		t.Error("expected NewSession")
	}

	foundFish := false
	for _, c := range mock.Calls {
		if c.Method == "SendKeys" && len(c.Args) >= 2 && c.Args[1] == "fish" {
			foundFish = true
		}
	}
	if !foundFish {
		t.Error("expected SendKeys with fish")
	}

	if mock.Called("NewWindow") {
		t.Error("--run should skip project config windows")
	}
	for _, c := range mock.Calls {
		if c.Method == "SendKeys" && len(c.Args) >= 2 && c.Args[1] == "nvim" {
			t.Error("--run should not send project config commands")
		}
	}
}

func TestTryAutoDetect_InsideProjectsDir(t *testing.T) {
	d := testDeps(t)
	blogDir := filepath.Join(d.global.ProjectsDir, "blog")
	if err := os.Mkdir(blogDir, 0o755); err != nil {
		t.Fatal(err)
	}
	d.getwd = func() (string, error) { return blogDir, nil }

	result, ok := tryAutoDetect(d)
	if !ok {
		t.Fatal("expected auto-detect to succeed")
	}
	if result.Name != "blog" {
		t.Errorf("Name = %q, want blog", result.Name)
	}
}

func TestTryAutoDetect_OutsideProjectsDir(t *testing.T) {
	d := testDeps(t)
	d.getwd = func() (string, error) { return "/some/other/dir", nil }

	_, ok := tryAutoDetect(d)
	if ok {
		t.Error("expected auto-detect to fail outside projects dir")
	}
}

func TestCollectPickerItems(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("blog", &config.ProjectConfig{Command: "a"})
	_ = d.store.Save("api", &config.ProjectConfig{Command: "b"})
	mock := d.client.(*tmux.MockClient)
	mock.ListSessionsReturn = []tmux.SessionInfo{
		{Name: "blog"},
		{Name: "scratch"},
	}

	items := collectPickerItems(d)
	if len(items) != 3 {
		t.Fatalf("expected 3 items (blog, api, scratch), got %d: %v", len(items), items)
	}

	seen := make(map[string]bool)
	for _, item := range items {
		seen[item] = true
	}
	if !seen["blog"] || !seen["api"] || !seen["scratch"] {
		t.Errorf("missing expected items: %v", items)
	}
}

func TestCollectPickerItems_DedupesNormalizedNames(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("my.project", &config.ProjectConfig{Command: "a"})
	mock := d.client.(*tmux.MockClient)
	mock.ListSessionsReturn = []tmux.SessionInfo{
		{Name: "my_project"},
	}

	items := collectPickerItems(d)
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %v", items)
	}
	if items[0] != "my.project" {
		t.Fatalf("got %q, want my.project (project name wins over tmux session name)", items[0])
	}
}

func TestRunBareNux_AutoDetect(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	blogDir := filepath.Join(d.global.ProjectsDir, "blog")
	if err := os.Mkdir(blogDir, 0o755); err != nil {
		t.Fatal(err)
	}
	d.getwd = func() (string, error) { return blogDir, nil }

	err := runBareNux(d)
	if err != nil {
		t.Fatalf("runBareNux: %v", err)
	}

	mock := d.client.(*tmux.MockClient)
	if !mock.Called("NewSession") {
		t.Error("expected NewSession for auto-detected project")
	}
}

func TestRunSessions_AttachesLast(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("blog", &config.ProjectConfig{Root: d.global.ProjectsDir, Command: "vim"})

	err := runSessions(d, []string{"blog"})
	if err != nil {
		t.Fatalf("runSessions: %v", err)
	}

	mock := d.client.(*tmux.MockClient)
	if !mock.Called("AttachSession") {
		t.Error("expected AttachSession for last session")
	}
}

func TestRunBareNux_WithRunCommand(t *testing.T) {
	d := testDeps(t)
	d.run = "echo hi"
	d.builder.SetAdHocLayout(&tmux.AdHocLayout{Command: "echo hi"})

	blogDir := filepath.Join(d.global.ProjectsDir, "blog")
	if err := os.Mkdir(blogDir, 0o755); err != nil {
		t.Fatal(err)
	}
	d.getwd = func() (string, error) { return blogDir, nil }

	err := runBareNux(d)
	if err != nil {
		t.Fatalf("runBareNux: %v", err)
	}

	mock := d.client.(*tmux.MockClient)
	if !mock.Called("AttachSession") {
		t.Error("expected AttachSession")
	}
	found := false
	for _, c := range mock.Calls {
		if c.Method == "SendKeys" && len(c.Args) >= 2 && c.Args[1] == "echo hi" {
			found = true
		}
	}
	if !found {
		t.Error("expected SendKeys with run command")
	}
}

func TestOpenInEditor(t *testing.T) {
	d := testDeps(t)
	d.editor = "echo"
	err := openInEditor(d, "/tmp/test.yaml")
	if err != nil {
		t.Fatalf("openInEditor: %v", err)
	}
}

func TestOpenInEditor_NoEditor(t *testing.T) {
	d := testDeps(t)
	d.editor = ""
	err := openInEditor(d, "/tmp/test.yaml")
	if err != nil {
		t.Error("expected nil when no editor set")
	}
	stderr := stderrStr(d)
	if !strings.Contains(stderr, "$EDITOR") {
		t.Errorf("expected hint about $EDITOR in stderr, got %q", stderr)
	}
}

func TestOpenInEditor_EditorFailure(t *testing.T) {
	d := testDeps(t)
	d.editor = "false"
	d.execCmd = exec.Command

	err := openInEditor(d, "/tmp/test.yaml")
	if err == nil {
		t.Fatal("expected error from failing editor")
	}
	if !strings.Contains(err.Error(), "editor failed") {
		t.Errorf("error = %q, expected 'editor failed'", err.Error())
	}
}

type fakePicker struct {
	choice string
	called bool
}

func (f *fakePicker) Pick([]string, string) (string, error) {
	f.called = true
	return f.choice, nil
}

func TestRunBareNux_Picker(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	d.global.PickerOnBare = true
	d.getwd = func() (string, error) { return "/some/other/dir", nil }
	_ = d.store.Save("blog", &config.ProjectConfig{Root: d.global.ProjectsDir, Command: "a"})

	fp := &fakePicker{choice: "blog"}
	d.newPicker = func(_ string, _ io.Writer) (picker.Picker, error) {
		return fp, nil
	}

	err := runBareNux(d)
	if err != nil {
		t.Fatalf("runBareNux: %v", err)
	}
	if !fp.called {
		t.Error("expected picker to be called")
	}
}

func TestRunBareNux_NoProjectsNoPicker(t *testing.T) {
	d := testDeps(t)
	d.getwd = func() (string, error) { return "/some/other/dir", nil }
	d.global.PickerOnBare = true

	err := runBareNux(d)
	if err == nil {
		t.Fatal("expected error when no projects found")
	}
	if !strings.Contains(err.Error(), "no projects") {
		t.Errorf("error = %q, expected 'no projects'", err.Error())
	}
}

func TestRunBareNux_Help(t *testing.T) {
	d := testDeps(t)
	d.getwd = func() (string, error) { return "/some/other/dir", nil }

	helpCalled := false
	d.help = func() error { helpCalled = true; return nil }

	err := runBareNux(d)
	if err != nil {
		t.Fatalf("runBareNux: %v", err)
	}
	if !helpCalled {
		t.Error("expected help to be called")
	}
}

func TestRunBareNux_PickerDismissed(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	d.global.PickerOnBare = true
	d.getwd = func() (string, error) { return "/some/other/dir", nil }
	_ = d.store.Save("blog", &config.ProjectConfig{Root: d.global.ProjectsDir, Command: "a"})

	fp := &fakePicker{choice: ""}
	d.newPicker = func(_ string, _ io.Writer) (picker.Picker, error) {
		return fp, nil
	}

	err := runBareNux(d)
	if err != nil {
		t.Fatalf("runBareNux: %v", err)
	}
	mock := d.client.(*tmux.MockClient)
	if mock.Called("NewSession") {
		t.Error("should not create session when picker is dismissed")
	}
}

func TestRunBareNux_PickerError(t *testing.T) {
	d := testDeps(t)
	d.global.PickerOnBare = true
	d.getwd = func() (string, error) { return "/some/other/dir", nil }
	_ = d.store.Save("blog", &config.ProjectConfig{Root: d.global.ProjectsDir, Command: "a"})

	d.newPicker = func(_ string, _ io.Writer) (picker.Picker, error) {
		return nil, fmt.Errorf("no picker binary")
	}

	err := runBareNux(d)
	if err == nil {
		t.Fatal("expected error from picker creation failure")
	}
}

func TestRunBareNux_AutoDetect_Attaches(t *testing.T) {
	d := testDeps(t)
	blogDir := filepath.Join(d.global.ProjectsDir, "blog")
	if err := os.Mkdir(blogDir, 0o755); err != nil {
		t.Fatal(err)
	}
	d.getwd = func() (string, error) { return blogDir, nil }

	err := runBareNux(d)
	if err != nil {
		t.Fatalf("runBareNux: %v", err)
	}

	mock := d.client.(*tmux.MockClient)
	if !mock.Called("AttachSession") {
		t.Error("expected AttachSession when noAttach=false")
	}
}

func TestRunSessions_VarWithRunWarning(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	d.run = "go test ./..."
	d.vars = map[string]string{"port": "8080"}
	d.builder.SetAdHocLayout(&tmux.AdHocLayout{Command: "go test ./..."})

	blogDir := filepath.Join(d.global.ProjectsDir, "blog")
	if err := os.Mkdir(blogDir, 0o755); err != nil {
		t.Fatal(err)
	}

	err := runSessions(d, []string{"blog"})
	if err != nil {
		t.Fatalf("runSessions: %v", err)
	}

	stderr := stderrStr(d)
	if !strings.Contains(stderr, "--var is ignored") {
		t.Errorf("expected warning in stderr, got %q", stderr)
	}
}

func TestRunSessions_BuildError(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	mock := d.client.(*tmux.MockClient)
	mock.DefaultError = fmt.Errorf("session failed")
	_ = d.store.Save("blog", &config.ProjectConfig{
		Root:    d.global.ProjectsDir,
		Command: "vim",
	})

	err := runSessions(d, []string{"blog"})
	if err == nil {
		t.Fatal("expected error from build failure")
	}
	if !strings.Contains(err.Error(), "building session") {
		t.Errorf("error = %q, expected 'building session'", err.Error())
	}
}

func TestRunSessions_ResolveError(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true

	err := runSessions(d, []string{"nonexistent"})
	if err == nil {
		t.Fatal("expected error from resolve failure")
	}
}

func TestRunSessions_Subset_NewSession_UserOrder(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	_ = d.store.Save("blog", &config.ProjectConfig{
		Root: d.global.ProjectsDir,
		Windows: []config.Window{
			{Name: "editor", Panes: []config.Pane{{Command: "nvim"}}},
			{Name: "server", Panes: []config.Pane{{Command: "go run ."}}},
		},
	})

	err := runSessions(d, []string{"blog:server,editor"})
	if err != nil {
		t.Fatalf("runSessions: %v", err)
	}

	mock := d.client.(*tmux.MockClient)
	var first *tmux.NewSessionOpts
	for _, c := range mock.Calls {
		if c.Method == "NewSession" && c.Opts != nil {
			if o, ok := c.Opts.(tmux.NewSessionOpts); ok {
				first = &o
				break
			}
		}
	}
	if first == nil || first.Window != "server" {
		t.Errorf("first NewSession window = %v, want server", first)
	}
}

func TestRunSessions_Subset_Existing_SelectsWindow(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	mock := d.client.(*tmux.MockClient)
	mock.HasSessionReturn = true
	_ = d.store.Save("blog", &config.ProjectConfig{
		Root: d.global.ProjectsDir,
		Windows: []config.Window{
			{Name: "editor", Panes: []config.Pane{{Command: "nvim"}}},
		},
	})

	err := runSessions(d, []string{"blog:editor"})
	if err != nil {
		t.Fatalf("runSessions: %v", err)
	}

	found := false
	for _, c := range mock.Calls {
		if c.Method == "SelectWindow" && len(c.Args) >= 2 && c.Args[0] == "blog" && c.Args[1] == "editor" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected SelectWindow for blog:editor")
	}
	if mock.Called("NewSession") {
		t.Error("should not create session when it already exists")
	}
}

func TestRunSessions_Subset_AdhocFlagsError(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	d.run = "fish"
	_ = d.store.Save("blog", &config.ProjectConfig{
		Root: d.global.ProjectsDir,
		Windows: []config.Window{
			{Name: "editor", Panes: []config.Pane{{Command: "nvim"}}},
		},
	})

	err := runSessions(d, []string{"blog:editor"})
	if err == nil {
		t.Fatal("expected error when combining --run with :window")
	}
	if !strings.Contains(err.Error(), "cannot combine") {
		t.Errorf("error = %q", err.Error())
	}
}

func TestRunSessions_Subset_CommandOnlyConfigError(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	_ = d.store.Save("blog", &config.ProjectConfig{
		Root:    d.global.ProjectsDir,
		Command: "vim",
	})

	err := runSessions(d, []string{"blog:editor"})
	if err == nil {
		t.Fatal("expected error for command-only project with :window")
	}
}

func TestRunBareNux_AutoDetect_BuildError(t *testing.T) {
	d := testDeps(t)
	d.noAttach = true
	mock := d.client.(*tmux.MockClient)
	mock.DefaultError = fmt.Errorf("create failed")
	blogDir := filepath.Join(d.global.ProjectsDir, "blog")
	if err := os.Mkdir(blogDir, 0o755); err != nil {
		t.Fatal(err)
	}
	d.getwd = func() (string, error) { return blogDir, nil }

	err := runBareNux(d)
	if err == nil {
		t.Fatal("expected error from build failure")
	}
	if !strings.Contains(err.Error(), "building session") {
		t.Errorf("error = %q, expected 'building session'", err.Error())
	}
}
