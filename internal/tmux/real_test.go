package tmux

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func dryClient() (*RealClient, *bytes.Buffer) {
	var buf bytes.Buffer
	c := NewRealClient()
	c.DryRun = true
	c.DryRunOut = &buf
	c.Stdin = strings.NewReader("")
	c.Stderr = &bytes.Buffer{}
	c.ExecCmd = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("true")
	}
	return c, &buf
}

func fakeCmd(stdout string) func(string, ...string) *exec.Cmd {
	return func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", "-n", stdout)
	}
}

func failCmd() func(string, ...string) *exec.Cmd {
	return func(name string, arg ...string) *exec.Cmd {
		return exec.Command("false")
	}
}

func TestNewRealClient(t *testing.T) {
	c := NewRealClient()
	if c.DryRun {
		t.Error("DryRun should default to false")
	}
	if c.LookupEnv == nil {
		t.Error("LookupEnv should be set")
	}
	if c.ExecCmd == nil {
		t.Error("ExecCmd should be set")
	}
}

func TestNewSession_DryRun(t *testing.T) {
	c, buf := dryClient()
	err := c.NewSession(NewSessionOpts{Name: "test", Root: "/tmp", Window: "editor", Detach: true})
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "new-session") {
		t.Errorf("expected new-session, got %q", out)
	}
	if !strings.Contains(out, "-d") {
		t.Errorf("expected -d flag, got %q", out)
	}
	if !strings.Contains(out, "-s test") {
		t.Errorf("expected -s test, got %q", out)
	}
	if !strings.Contains(out, "-c /tmp") {
		t.Errorf("expected -c /tmp, got %q", out)
	}
	if !strings.Contains(out, "-n editor") {
		t.Errorf("expected -n editor, got %q", out)
	}
}

func TestNewSession_DryRun_Minimal(t *testing.T) {
	c, buf := dryClient()
	err := c.NewSession(NewSessionOpts{Name: "bare"})
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if strings.Contains(out, "-d") {
		t.Error("should not include -d when Detach=false")
	}
	if strings.Contains(out, "-c") {
		t.Error("should not include -c when Root is empty")
	}
	if strings.Contains(out, "-n") {
		t.Error("should not include -n when Window is empty")
	}
}

func TestKillSession_DryRun(t *testing.T) {
	c, buf := dryClient()
	_ = c.KillSession("myproj")
	if !strings.Contains(buf.String(), "kill-session -t myproj") {
		t.Errorf("got %q", buf.String())
	}
}

func TestNewWindow_DryRun(t *testing.T) {
	c, buf := dryClient()
	_ = c.NewWindow("sess", NewWindowOpts{Name: "editor", Root: "/tmp"})
	out := buf.String()
	if !strings.Contains(out, "new-window -t sess") {
		t.Errorf("got %q", out)
	}
	if !strings.Contains(out, "-n editor") {
		t.Errorf("expected -n editor, got %q", out)
	}
	if !strings.Contains(out, "-c /tmp") {
		t.Errorf("expected -c /tmp, got %q", out)
	}
}

func TestNewWindow_DryRun_Minimal(t *testing.T) {
	c, buf := dryClient()
	_ = c.NewWindow("sess", NewWindowOpts{})
	out := buf.String()
	if strings.Contains(out, "-n") {
		t.Error("should not include -n when Name is empty")
	}
}

func TestKillWindow_DryRun(t *testing.T) {
	c, buf := dryClient()
	_ = c.KillWindow("sess", "editor")
	if !strings.Contains(buf.String(), "kill-window -t sess:editor") {
		t.Errorf("got %q", buf.String())
	}
}

func TestSplitWindow_DryRun(t *testing.T) {
	c, buf := dryClient()
	_ = c.SplitWindow("sess", "editor", SplitWindowOpts{Root: "/tmp"})
	out := buf.String()
	if !strings.Contains(out, "split-window -t sess:editor") {
		t.Errorf("got %q", out)
	}
	if !strings.Contains(out, "-v") {
		t.Errorf("expected -v, got %q", out)
	}
}

func TestSplitWindow_DryRun_NoRoot(t *testing.T) {
	c, buf := dryClient()
	_ = c.SplitWindow("sess", "editor", SplitWindowOpts{})
	out := buf.String()
	if !strings.Contains(out, "-v") {
		t.Errorf("expected -v, got %q", out)
	}
	if strings.Contains(out, "-c") {
		t.Error("should not include -c when Root is empty")
	}
}

func TestSelectLayout_DryRun(t *testing.T) {
	c, buf := dryClient()
	_ = c.SelectLayout("sess", "editor", "tiled")
	if !strings.Contains(buf.String(), "select-layout -t sess:editor tiled") {
		t.Errorf("got %q", buf.String())
	}
}

func TestSelectWindow_DryRun(t *testing.T) {
	c, buf := dryClient()
	_ = c.SelectWindow("sess", "editor")
	if !strings.Contains(buf.String(), "select-window -t sess:editor") {
		t.Errorf("got %q", buf.String())
	}
}

func TestSelectPane_DryRun(t *testing.T) {
	c, buf := dryClient()
	_ = c.SelectPane("sess", "editor", 2)
	if !strings.Contains(buf.String(), "select-pane -t sess:editor.2") {
		t.Errorf("got %q", buf.String())
	}
}

func TestSendKeys_DryRun(t *testing.T) {
	c, buf := dryClient()
	_ = c.SendKeys("sess:editor", "vim")
	out := buf.String()
	if !strings.Contains(out, "send-keys -t sess:editor vim Enter") {
		t.Errorf("got %q", out)
	}
}

func TestAttachSession_DryRun_OutsideTmux(t *testing.T) {
	c, buf := dryClient()
	c.LookupEnv = func(string) string { return "" }
	_ = c.AttachSession("myproj")
	if !strings.Contains(buf.String(), "attach-session -t myproj") {
		t.Errorf("got %q", buf.String())
	}
}

func TestAttachSession_DryRun_InsideTmux(t *testing.T) {
	c, buf := dryClient()
	c.LookupEnv = func(key string) string {
		if key == "TMUX" {
			return "/tmp/tmux-1000/default,12345,0"
		}
		return ""
	}
	_ = c.AttachSession("myproj")
	if !strings.Contains(buf.String(), "switch-client -t myproj") {
		t.Errorf("got %q", buf.String())
	}
}

func TestSetEnv_DryRun(t *testing.T) {
	c, buf := dryClient()
	_ = c.SetEnv("sess", "FOO", "bar")
	if !strings.Contains(buf.String(), "set-environment -t sess FOO bar") {
		t.Errorf("got %q", buf.String())
	}
}

func TestSetOption_DryRun(t *testing.T) {
	c, buf := dryClient()
	_ = c.SetOption("sess", "default-command", "/bin/zsh")
	if !strings.Contains(buf.String(), "set-option -t sess default-command /bin/zsh") {
		t.Errorf("got %q", buf.String())
	}
}

func TestSetHook_DryRun(t *testing.T) {
	c, buf := dryClient()
	_ = c.SetHook("sess", "session-closed[0]", "echo bye")
	out := buf.String()
	if !strings.Contains(out, "set-hook -t sess session-closed[0]") {
		t.Errorf("got %q", out)
	}
	if !strings.Contains(out, "run-shell 'echo bye'") {
		t.Errorf("expected run-shell wrapper, got %q", out)
	}
}

func TestSetHook_DryRun_QuoteEscaping(t *testing.T) {
	c, buf := dryClient()
	_ = c.SetHook("sess", "session-closed[0]", "echo 'hello'")
	out := buf.String()
	expected := `run-shell 'echo '\''hello'\'''`
	if !strings.Contains(out, expected) {
		t.Errorf("expected %q in output, got %q", expected, out)
	}
}

func TestHasSession_DryRun(t *testing.T) {
	c, buf := dryClient()
	c.ExecCmd = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("true")
	}
	got := c.HasSession("myproj")
	if !got {
		t.Error("expected true when command succeeds")
	}
	if !strings.Contains(buf.String(), "has-session") {
		t.Errorf("expected dry-run output, got %q", buf.String())
	}
}

func TestHasSession_NotFound(t *testing.T) {
	c, _ := dryClient()
	c.ExecCmd = failCmd()
	got := c.HasSession("missing")
	if got {
		t.Error("expected false when command fails")
	}
}

func TestIsInsideTmux_True(t *testing.T) {
	c, _ := dryClient()
	c.LookupEnv = func(key string) string {
		if key == "TMUX" {
			return "/tmp/tmux-1000/default,12345,0"
		}
		return ""
	}
	if !c.IsInsideTmux() {
		t.Error("expected true when TMUX is set")
	}
}

func TestIsInsideTmux_False(t *testing.T) {
	c, _ := dryClient()
	c.LookupEnv = func(string) string { return "" }
	if c.IsInsideTmux() {
		t.Error("expected false when TMUX is empty")
	}
}

func TestListSessions_DryRun(t *testing.T) {
	c, buf := dryClient()
	sessions, err := c.ListSessions()
	if err != nil {
		t.Fatal(err)
	}
	if sessions != nil {
		t.Errorf("expected nil sessions in dry-run, got %v", sessions)
	}
	if !strings.Contains(buf.String(), "list-sessions") {
		t.Errorf("expected list-sessions in dry-run output, got %q", buf.String())
	}
}

func TestListSessions_ParsesOutput(t *testing.T) {
	c := NewRealClient()
	c.ExecCmd = fakeCmd("blog|2|1700000000|1\napi|1|1700000000|0")
	c.Stdin = strings.NewReader("")
	c.Stderr = &bytes.Buffer{}

	sessions, err := c.ListSessions()
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}
	if sessions[0].Name != "blog" {
		t.Errorf("Name = %q, want blog", sessions[0].Name)
	}
	if sessions[0].Windows != 2 {
		t.Errorf("Windows = %d, want 2", sessions[0].Windows)
	}
	if !sessions[0].Attached {
		t.Error("expected Attached=true for blog")
	}
	if sessions[1].Attached {
		t.Error("expected Attached=false for api")
	}
}

func TestListSessions_Empty(t *testing.T) {
	c := NewRealClient()
	c.ExecCmd = fakeCmd("")
	sessions, err := c.ListSessions()
	if err != nil {
		t.Fatal(err)
	}
	if sessions != nil {
		t.Errorf("expected nil for empty output, got %v", sessions)
	}
}

func TestListSessions_BadWindowCount(t *testing.T) {
	c := NewRealClient()
	c.ExecCmd = fakeCmd("blog|notanum|1700000000|1")
	_, err := c.ListSessions()
	if err == nil {
		t.Fatal("expected error for bad window count")
	}
}

func TestListSessions_BadTimestamp(t *testing.T) {
	c := NewRealClient()
	c.ExecCmd = fakeCmd("blog|2|notanum|1")
	_, err := c.ListSessions()
	if err == nil {
		t.Fatal("expected error for bad timestamp")
	}
}

func TestListSessions_MalformedLine(t *testing.T) {
	c := NewRealClient()
	c.ExecCmd = fakeCmd("blog|2|1700000000|1\nbadline\napi|1|1700000000|0")
	sessions, err := c.ListSessions()
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions (skipping bad line), got %d", len(sessions))
	}
}

func TestListSessions_Error(t *testing.T) {
	c := NewRealClient()
	c.ExecCmd = failCmd()
	_, err := c.ListSessions()
	if err == nil {
		t.Fatal("expected error when command fails")
	}
}

func TestBaseIndex_Fallback(t *testing.T) {
	c := NewRealClient()
	c.ExecCmd = failCmd()
	if got := c.BaseIndex(); got != 0 {
		t.Errorf("BaseIndex = %d, want 0 (fallback)", got)
	}
}

func TestBaseIndex_Parsed(t *testing.T) {
	c := NewRealClient()
	c.ExecCmd = fakeCmd("1")
	if got := c.BaseIndex(); got != 1 {
		t.Errorf("BaseIndex = %d, want 1", got)
	}
}

func TestPaneBaseIndex_Fallback(t *testing.T) {
	c := NewRealClient()
	c.ExecCmd = failCmd()
	if got := c.PaneBaseIndex(); got != 0 {
		t.Errorf("PaneBaseIndex = %d, want 0 (fallback)", got)
	}
}

func TestPaneBaseIndex_Parsed(t *testing.T) {
	c := NewRealClient()
	c.ExecCmd = fakeCmd("1")
	if got := c.PaneBaseIndex(); got != 1 {
		t.Errorf("PaneBaseIndex = %d, want 1", got)
	}
}

func TestGlobalIntOption_BadOutput(t *testing.T) {
	c := NewRealClient()
	c.ExecCmd = fakeCmd("notanumber")
	if got := c.globalIntOption("base-index", 42); got != 42 {
		t.Errorf("got %d, want 42 (fallback)", got)
	}
}

func TestRun_NonDryRun(t *testing.T) {
	c := NewRealClient()
	c.ExecCmd = func(name string, arg ...string) *exec.Cmd { return exec.Command("true") }
	c.Stdin = strings.NewReader("")
	c.Stderr = &bytes.Buffer{}

	if err := c.run("some-command"); err != nil {
		t.Fatalf("run: %v", err)
	}
}

func TestRun_NonDryRun_Error(t *testing.T) {
	c := NewRealClient()
	c.ExecCmd = failCmd()
	c.Stdin = strings.NewReader("")
	c.Stderr = &bytes.Buffer{}

	if err := c.run("some-command"); err == nil {
		t.Fatal("expected error from failing command")
	}
}

func TestRunOutput_NonDryRun(t *testing.T) {
	c := NewRealClient()
	c.ExecCmd = fakeCmd("hello")

	out, err := c.runOutput("some-command")
	if err != nil {
		t.Fatalf("runOutput: %v", err)
	}
	if out != "hello" {
		t.Errorf("output = %q, want hello", out)
	}
}

func TestRunOutput_NonDryRun_Error(t *testing.T) {
	c := NewRealClient()
	c.ExecCmd = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("sh", "-c", "echo err >&2; exit 1")
	}

	_, err := c.runOutput("some-command")
	if err == nil {
		t.Fatal("expected error from failing command")
	}
	if !strings.Contains(err.Error(), "err") {
		t.Errorf("error should include stderr: %v", err)
	}
}
