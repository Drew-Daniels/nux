package tmux

import "fmt"

type Call struct {
	Method string
	Args   []string
	Opts   any
}

type MockClient struct {
	Calls []Call

	HasSessionReturn    bool
	HasSessionFunc      func(string) bool
	ListSessionsReturn  []SessionInfo
	ListSessionsError   error
	IsInsideTmuxReturn  bool
	BaseIndexReturn     int
	PaneBaseIndexReturn int
	DefaultError        error
}

func (m *MockClient) record(method string, opts any, args ...string) {
	m.Calls = append(m.Calls, Call{Method: method, Args: args, Opts: opts})
}

func (m *MockClient) Called(method string) bool {
	for _, c := range m.Calls {
		if c.Method == method {
			return true
		}
	}
	return false
}

func (m *MockClient) HasSession(name string) bool {
	m.record("HasSession", nil, name)
	if m.HasSessionFunc != nil {
		return m.HasSessionFunc(name)
	}
	return m.HasSessionReturn
}

func (m *MockClient) NewSession(opts NewSessionOpts) error {
	m.record("NewSession", opts, opts.Name, opts.Root, opts.Window)
	return m.DefaultError
}

func (m *MockClient) KillSession(name string) error {
	m.record("KillSession", nil, name)
	return m.DefaultError
}

func (m *MockClient) NewWindow(session string, opts NewWindowOpts) error {
	m.record("NewWindow", opts, session, opts.Name, opts.Root)
	return m.DefaultError
}

func (m *MockClient) KillWindow(session, window string) error {
	m.record("KillWindow", nil, session, window)
	return m.DefaultError
}

func (m *MockClient) SplitWindow(session, window string, opts SplitWindowOpts) error {
	m.record("SplitWindow", opts, session, window, opts.Root)
	return m.DefaultError
}

func (m *MockClient) SelectLayout(session, window, layout string) error {
	m.record("SelectLayout", nil, session, window, layout)
	return m.DefaultError
}

func (m *MockClient) SelectWindow(session, window string) error {
	m.record("SelectWindow", nil, session, window)
	return m.DefaultError
}

func (m *MockClient) SelectPane(session, window string, pane int) error {
	m.record("SelectPane", nil, session, window, fmt.Sprintf("%d", pane))
	return m.DefaultError
}

func (m *MockClient) SendKeys(target, keys string) error {
	m.record("SendKeys", nil, target, keys)
	return m.DefaultError
}

func (m *MockClient) AttachSession(name string) error {
	m.record("AttachSession", nil, name)
	return m.DefaultError
}

func (m *MockClient) SetEnv(session, key, value string) error {
	m.record("SetEnv", nil, session, key, value)
	return m.DefaultError
}

func (m *MockClient) SetOption(session, key, value string) error {
	m.record("SetOption", nil, session, key, value)
	return m.DefaultError
}

func (m *MockClient) SetHook(session, hookName, command string) error {
	m.record("SetHook", nil, session, hookName, command)
	return m.DefaultError
}

func (m *MockClient) ListSessions() ([]SessionInfo, error) {
	m.record("ListSessions", nil)
	return m.ListSessionsReturn, m.ListSessionsError
}

func (m *MockClient) IsInsideTmux() bool {
	m.record("IsInsideTmux", nil)
	return m.IsInsideTmuxReturn
}

func (m *MockClient) BaseIndex() int {
	m.record("BaseIndex", nil)
	return m.BaseIndexReturn
}

func (m *MockClient) PaneBaseIndex() int {
	m.record("PaneBaseIndex", nil)
	return m.PaneBaseIndexReturn
}

var _ Client = (*MockClient)(nil)
var _ Client = (*RealClient)(nil)
