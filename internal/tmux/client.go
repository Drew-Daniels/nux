package tmux

import "time"

type Client interface {
	HasSession(name string) bool
	NewSession(opts NewSessionOpts) error
	KillSession(name string) error
	NewWindow(session string, opts NewWindowOpts) error
	KillWindow(session, window string) error
	SplitWindow(session, window string, opts SplitWindowOpts) error
	SelectLayout(session, window, layout string) error
	SelectWindow(session, window string) error
	SelectPane(session, window string, pane int) error
	SendKeys(target, keys string) error
	AttachSession(name string) error
	SetEnv(session, key, value string) error
	SetOption(session, key, value string) error
	SetHook(session, hookName, command string) error
	ListSessions() ([]SessionInfo, error)
	IsInsideTmux() bool
	BaseIndex() int
	PaneBaseIndex() int
}

type NewSessionOpts struct {
	Name   string
	Root   string
	Window string
	Detach bool
}

type NewWindowOpts struct {
	Name string
	Root string
}

type SplitWindowOpts struct {
	Root string
}

type SessionInfo struct {
	Name     string
	Windows  int
	Created  time.Time
	Attached bool
}
