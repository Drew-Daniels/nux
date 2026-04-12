package tmux

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Drew-Daniels/nux/internal/config"
	"github.com/Drew-Daniels/nux/internal/resolver"
)

// AdHocLayout holds overrides from CLI flags (--layout, --panes, --run).
type AdHocLayout struct {
	Layout  string
	Panes   int
	Command string
}

type Builder struct {
	client        Client
	global        *config.GlobalConfig
	adHocLayout   *AdHocLayout
	baseIndex     int
	paneBaseIndex int
	indexResolved bool
}

func NewBuilder(client Client, global *config.GlobalConfig) *Builder {
	return &Builder{
		client: client,
		global: global,
	}
}

// resolveIndices queries base-index and pane-base-index from the running tmux
// server exactly once. This is deferred from construction time so that the
// server has been started (by createDetachedSession) before we query it.
func (b *Builder) resolveIndices() {
	if b.indexResolved {
		return
	}
	b.baseIndex = b.client.BaseIndex()
	b.paneBaseIndex = b.client.PaneBaseIndex()
	b.indexResolved = true
}

// SetAdHocLayout configures a layout override for the next Build call.
// Pass nil to clear.
func (b *Builder) SetAdHocLayout(layout *AdHocLayout) {
	b.adHocLayout = layout
}

func (b *Builder) firstWindow() string {
	return strconv.Itoa(b.baseIndex)
}

func (b *Builder) paneBase() int {
	return b.paneBaseIndex
}

func (b *Builder) createDetachedSession(name, root, window string) error {
	if err := b.client.NewSession(NewSessionOpts{
		Name:   name,
		Root:   root,
		Window: window,
		Detach: true,
	}); err != nil {
		return fmt.Errorf("creating session: %w", err)
	}
	b.resolveIndices()
	return nil
}

func (b *Builder) Build(name string, cfg *config.ProjectConfig, root string) error {
	if cfg == nil {
		return b.buildDefault(name, root)
	}
	return b.buildWindowed(name, cfg, root)
}

func (b *Builder) buildDefault(name, root string) error {
	ds := b.global.DefaultSession

	// -x: ignore the default session template, use ad-hoc or bare.
	if b.adHocCommand() != "" {
		if b.hasAdHocPanes() {
			return b.buildAdHoc(name, root)
		}
		return b.buildBare(name, root)
	}

	// --layout/--panes without -x: ad-hoc layout.
	if b.hasAdHocPanes() {
		return b.buildAdHoc(name, root)
	}

	// No flags and no default session: bare shell.
	if ds == nil || len(ds.Windows) == 0 {
		return b.buildBare(name, root)
	}

	// Default session with windows: build them.
	return b.buildWindowed(name, &config.ProjectConfig{Windows: ds.Windows}, root)
}

func (b *Builder) buildBare(name, root string) error {
	if err := b.createDetachedSession(name, root, ""); err != nil {
		return err
	}

	cmd := b.adHocCommand()
	if cmd == "" {
		return nil
	}

	firstWindow := name + ":" + b.firstWindow()
	var errs []error
	errs = append(errs, b.sendPaneInit(firstWindow, nil)...)
	errs = append(errs, b.client.SendKeys(firstWindow, cmd))
	return errors.Join(errs...)
}

func (b *Builder) buildAdHoc(name, root string) error {
	if err := b.createDetachedSession(name, root, ""); err != nil {
		return err
	}

	fw := b.firstWindow()
	target := name + ":" + fw
	var errs []error

	for i := 1; i < b.adHocLayout.Panes; i++ {
		errs = append(errs, b.client.SplitWindow(name, fw, SplitWindowOpts{Root: root}))
	}

	if b.adHocLayout.Layout != "" {
		errs = append(errs, b.client.SelectLayout(name, fw, b.adHocLayout.Layout))
	}

	pb := b.paneBase()
	cmd := b.adHocCommand()
	for i := 0; i < b.adHocLayout.Panes; i++ {
		paneTarget := fmt.Sprintf("%s.%d", target, pb+i)
		errs = append(errs, b.sendPaneInit(paneTarget, nil)...)
		if cmd != "" {
			errs = append(errs, b.client.SendKeys(paneTarget, cmd))
		}
	}

	errs = append(errs, b.client.SelectPane(name, fw, pb))

	return errors.Join(errs...)
}

func (b *Builder) buildWindowed(name string, cfg *config.ProjectConfig, root string) error {
	return b.buildWindowList(name, cfg, root, cfg.Windows)
}

// BuildWindows creates a session containing only the named windows, in the
// given order. Session-level settings and hooks use the first window in names
// as the anchor (on_start / on_ready targets).
func (b *Builder) BuildWindows(name string, cfg *config.ProjectConfig, root string, names []string) error {
	if cfg == nil {
		return fmt.Errorf("window selection requires a project config")
	}
	if len(cfg.Windows) == 0 {
		return fmt.Errorf("project %q has no windows defined", name)
	}
	if len(names) == 0 {
		return fmt.Errorf("no windows specified")
	}

	windows := make([]config.Window, 0, len(names))
	for _, n := range names {
		w, ok := findWindow(cfg, n)
		if !ok {
			return fmt.Errorf("window %q not found in config", n)
		}
		windows = append(windows, w)
	}

	return b.buildWindowList(name, cfg, root, windows)
}

func (b *Builder) buildWindowList(name string, cfg *config.ProjectConfig, root string, windows []config.Window) error {
	firstWin := windows[0]
	winRoot := windowRoot(firstWin.Root, root)

	if err := b.createDetachedSession(name, winRoot, firstWin.Name); err != nil {
		return err
	}

	firstTarget := name + ":" + firstWin.Name
	var errs []error

	errs = append(errs, b.applySessionSettings(name, cfg, firstTarget)...)
	errs = append(errs, b.startWindow(name, firstWin, root, cfg.PaneInit))

	for _, w := range windows[1:] {
		wr := windowRoot(w.Root, root)
		errs = append(errs, b.client.NewWindow(name, NewWindowOpts{Name: w.Name, Root: wr}))
		errs = append(errs, b.startWindow(name, w, root, cfg.PaneInit))
	}

	errs = append(errs, b.client.SelectWindow(name, firstWin.Name))
	errs = append(errs, b.sendOnReady(cfg, firstTarget)...)

	return errors.Join(errs...)
}

func (b *Builder) applySessionSettings(name string, cfg *config.ProjectConfig, firstWindow string) []error {
	var errs []error

	for k, v := range cfg.Env {
		errs = append(errs, b.client.SetEnv(name, k, v))
	}

	shell := cfg.DefaultShell
	if shell == "" {
		shell = b.global.DefaultShell
	}
	if shell != "" {
		errs = append(errs, b.client.SetOption(name, "default-command", shell))
	}

	for _, cmd := range cfg.OnStart {
		errs = append(errs, b.client.SendKeys(firstWindow, cmd))
	}

	for i, cmd := range cfg.OnStop {
		hookName := fmt.Sprintf("session-closed[%d]", i)
		errs = append(errs, b.client.SetHook(name, hookName, cmd))
	}

	for i, cmd := range cfg.OnDetach {
		hookName := fmt.Sprintf("client-detached[%d]", i)
		errs = append(errs, b.client.SetHook(name, hookName, cmd))
	}

	return errs
}

func (b *Builder) sendOnReady(cfg *config.ProjectConfig, firstWindow string) []error {
	var errs []error
	for _, cmd := range cfg.OnReady {
		errs = append(errs, b.client.SendKeys(firstWindow, cmd))
	}
	return errs
}

func (b *Builder) hasAdHocPanes() bool {
	return b.adHocLayout != nil && (b.adHocLayout.Layout != "" || b.adHocLayout.Panes > 0)
}

func (b *Builder) adHocCommand() string {
	if b.adHocLayout != nil {
		return b.adHocLayout.Command
	}
	return ""
}

func (b *Builder) sendPaneInit(target string, projectPaneInit []string) []error {
	var errs []error
	for _, cmd := range b.global.PaneInit {
		errs = append(errs, b.client.SendKeys(target, cmd))
	}
	for _, cmd := range projectPaneInit {
		errs = append(errs, b.client.SendKeys(target, cmd))
	}
	return errs
}

func (b *Builder) sendWindowEnv(target string, env map[string]string) []error {
	if len(env) == 0 {
		return nil
	}
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var errs []error
	for _, k := range keys {
		errs = append(errs, b.client.SendKeys(target, fmt.Sprintf("export %s='%s'", k, shellEscape(env[k]))))
	}
	return errs
}

func shellEscape(s string) string {
	return strings.ReplaceAll(s, "'", `'\''`)
}

func (b *Builder) startWindow(session string, w config.Window, projectRoot string, projectPaneInit []string) error {
	target := session + ":" + w.Name
	wr := windowRoot(w.Root, projectRoot)
	pb := b.paneBase()
	var errs []error

	for i, p := range w.Panes {
		paneTarget := fmt.Sprintf("%s.%d", target, pb+i)

		if i > 0 {
			pr := windowRoot(p.Root, wr)
			errs = append(errs, b.client.SplitWindow(session, w.Name, SplitWindowOpts{
				Root: pr,
			}))
		}

		errs = append(errs, b.sendPaneInit(paneTarget, projectPaneInit)...)
		errs = append(errs, b.sendWindowEnv(paneTarget, w.Env)...)

		if p.Command != "" {
			errs = append(errs, b.client.SendKeys(paneTarget, p.Command))
		}
	}

	layout := w.Layout
	if layout == "" && b.adHocLayout != nil {
		layout = b.adHocLayout.Layout
	}
	if layout != "" {
		errs = append(errs, b.client.SelectLayout(session, w.Name, layout))
	}

	errs = append(errs, b.client.SelectPane(session, w.Name, pb))

	return errors.Join(errs...)
}

func (b *Builder) StopSession(name string) error {
	if !b.client.HasSession(name) {
		return nil
	}
	return b.client.KillSession(name)
}

func (b *Builder) StopAll() error {
	sessions, err := b.client.ListSessions()
	if err != nil {
		return err
	}
	var errs []error
	for _, s := range sessions {
		errs = append(errs, b.client.KillSession(s.Name))
	}
	return errors.Join(errs...)
}

func (b *Builder) RestartSession(name string, cfg *config.ProjectConfig, root string) error {
	if err := b.StopSession(name); err != nil {
		return fmt.Errorf("stopping session %q: %w", name, err)
	}
	return b.Build(name, cfg, root)
}

func (b *Builder) RestartWindow(session, windowName string, cfg *config.ProjectConfig, root string) error {
	b.resolveIndices()

	w, ok := findWindow(cfg, windowName)
	if !ok {
		return fmt.Errorf("window %q not found in config", windowName)
	}

	if err := b.client.KillWindow(session, windowName); err != nil {
		return fmt.Errorf("killing window %q: %w", windowName, err)
	}

	wr := windowRoot(w.Root, root)
	if err := b.client.NewWindow(session, NewWindowOpts{Name: w.Name, Root: wr}); err != nil {
		return fmt.Errorf("creating window %q: %w", windowName, err)
	}

	return b.startWindow(session, w, root, cfg.PaneInit)
}

func findWindow(cfg *config.ProjectConfig, name string) (config.Window, bool) {
	if cfg == nil {
		return config.Window{}, false
	}
	for _, w := range cfg.Windows {
		if w.Name == name {
			return w, true
		}
	}
	return config.Window{}, false
}

// windowRoot resolves a window or pane working directory for tmux -c. Tilde
// must be expanded here: tmux does not treat ~ as home in -c, which would
// leave every pane in $HOME.
func windowRoot(winRoot, projectRoot string) string {
	return resolver.ResolveRoot(winRoot, projectRoot)
}
