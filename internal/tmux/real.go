package tmux

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type RealClient struct {
	DryRun    bool
	DryRunOut io.Writer
	Stdin     io.Reader
	Stderr    io.Writer
	LookupEnv func(string) string
	ExecCmd   func(name string, arg ...string) *exec.Cmd
}

func NewRealClient() *RealClient {
	return &RealClient{
		DryRunOut: os.Stdout,
		Stdin:     os.Stdin,
		Stderr:    os.Stderr,
		LookupEnv: os.Getenv,
		ExecCmd:   exec.Command,
	}
}

func (c *RealClient) command(args ...string) *exec.Cmd {
	return c.ExecCmd("tmux", args...)
}

func (c *RealClient) run(args ...string) error {
	if c.DryRun {
		_, _ = fmt.Fprintln(c.DryRunOut, "tmux "+strings.Join(args, " "))
		return nil
	}
	cmd := c.command(args...)
	cmd.Stdin = c.Stdin
	cmd.Stderr = c.Stderr
	return cmd.Run()
}

func (c *RealClient) runOutput(args ...string) (string, error) {
	if c.DryRun {
		_, _ = fmt.Fprintln(c.DryRunOut, "tmux "+strings.Join(args, " "))
		return "", nil
	}
	out, err := c.command(args...).Output()
	return strings.TrimSpace(string(out)), err
}

func (c *RealClient) HasSession(name string) bool {
	if c.DryRun {
		_, _ = fmt.Fprintln(c.DryRunOut, "# check: tmux has-session -t "+name)
	}
	// Always query the live server so dry-run output accurately reflects
	// whether the session would be created or just attached.
	cmd := c.command("has-session", "-t", name)
	return cmd.Run() == nil
}

func (c *RealClient) NewSession(opts NewSessionOpts) error {
	args := []string{"new-session"}
	if opts.Detach {
		args = append(args, "-d")
	}
	args = append(args, "-s", opts.Name)
	if opts.Root != "" {
		args = append(args, "-c", opts.Root)
	}
	if opts.Window != "" {
		args = append(args, "-n", opts.Window)
	}
	return c.run(args...)
}

func (c *RealClient) KillSession(name string) error {
	return c.run("kill-session", "-t", name)
}

func (c *RealClient) NewWindow(session string, opts NewWindowOpts) error {
	args := []string{"new-window", "-t", session}
	if opts.Name != "" {
		args = append(args, "-n", opts.Name)
	}
	if opts.Root != "" {
		args = append(args, "-c", opts.Root)
	}
	return c.run(args...)
}

func (c *RealClient) KillWindow(session, window string) error {
	return c.run("kill-window", "-t", session+":"+window)
}

func (c *RealClient) SplitWindow(session, window string, opts SplitWindowOpts) error {
	args := []string{"split-window", "-t", session + ":" + window}
	if opts.Horizontal {
		args = append(args, "-h")
	} else {
		args = append(args, "-v")
	}
	if opts.Root != "" {
		args = append(args, "-c", opts.Root)
	}
	return c.run(args...)
}

func (c *RealClient) SelectLayout(session, window, layout string) error {
	return c.run("select-layout", "-t", session+":"+window, layout)
}

func (c *RealClient) SelectWindow(session, window string) error {
	return c.run("select-window", "-t", session+":"+window)
}

func (c *RealClient) SelectPane(session, window string, pane int) error {
	target := fmt.Sprintf("%s:%s.%d", session, window, pane)
	return c.run("select-pane", "-t", target)
}

func (c *RealClient) SendKeys(target, keys string) error {
	return c.run("send-keys", "-t", target, keys, "Enter")
}

func (c *RealClient) AttachSession(name string) error {
	if c.IsInsideTmux() {
		return c.run("switch-client", "-t", name)
	}
	return c.run("attach-session", "-t", name)
}

func (c *RealClient) SetEnv(session, key, value string) error {
	return c.run("set-environment", "-t", session, key, value)
}

func (c *RealClient) SetOption(session, key, value string) error {
	return c.run("set-option", "-t", session, key, value)
}

func (c *RealClient) SetHook(session, hookName, command string) error {
	escaped := strings.ReplaceAll(command, "'", "'\"'\"'")
	hook := fmt.Sprintf("run-shell '%s'", escaped)
	return c.run("set-hook", "-t", session, hookName, hook)
}

func (c *RealClient) ListSessions() ([]SessionInfo, error) {
	format := "#{session_name}|#{session_windows}|#{session_created}|#{session_attached}"
	out, err := c.runOutput("list-sessions", "-F", format)
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}

	var sessions []SessionInfo
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue
		}

		windows, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("parsing window count %q: %w", parts[1], err)
		}

		created, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing created timestamp %q: %w", parts[2], err)
		}

		attached := parts[3] == "1"

		sessions = append(sessions, SessionInfo{
			Name:     parts[0],
			Windows:  windows,
			Created:  time.Unix(created, 0),
			Attached: attached,
		})
	}
	return sessions, nil
}

func (c *RealClient) IsInsideTmux() bool {
	return c.LookupEnv("TMUX") != ""
}

func (c *RealClient) BaseIndex() int {
	return c.globalIntOption("base-index", 0)
}

func (c *RealClient) PaneBaseIndex() int {
	return c.globalIntOption("pane-base-index", 0)
}

func (c *RealClient) globalIntOption(name string, fallback int) int {
	// Always query the live tmux server - this is read-only and must
	// return accurate values even during dry-run so that generated
	// targets use the correct base indices.
	out, err := c.command("show-options", "-gv", name).Output()
	if err != nil {
		return fallback
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return fallback
	}
	return n
}
