package tmux

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/Drew-Daniels/nux/internal/config"
)

// AdHocLayout holds layout overrides from CLI flags (--layout, --panes).
type AdHocLayout struct {
	Layout string
	Panes  int
}

type Builder struct {
	client        Client
	global        *config.GlobalConfig
	adHocLayout   *AdHocLayout
	baseIndex     int
	paneBaseIndex int
}

func NewBuilder(client Client, global *config.GlobalConfig) *Builder {
	return &Builder{
		client:        client,
		global:        global,
		baseIndex:     client.BaseIndex(),
		paneBaseIndex: client.PaneBaseIndex(),
	}
}

// SetAdHocLayout configures a layout override for the next Build or
// BuildEphemeral call. Pass nil to clear.
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
	return nil
}

func (b *Builder) Build(name string, cfg *config.ProjectConfig, root string) error {
	if cfg == nil {
		return b.buildDefault(name, root)
	}
	if len(cfg.Windows) > 0 {
		return b.buildWindowed(name, cfg, root)
	}
	return b.buildCommand(name, cfg, root)
}

func (b *Builder) buildDefault(name, root string) error {
	ds := b.global.DefaultSession

	if b.adHocLayout != nil && (ds == nil || len(ds.Windows) == 0) {
		return b.buildAdHoc(name, root, ds)
	}

	if ds == nil {
		return b.createDetachedSession(name, root, "")
	}

	if len(ds.Windows) > 0 {
		cfg := &config.ProjectConfig{Windows: ds.Windows}
		return b.buildWindowed(name, cfg, root)
	}

	if err := b.createDetachedSession(name, root, ""); err != nil {
		return err
	}

	fw := b.firstWindow()
	firstWindow := name + ":" + fw
	var errs []error
	errs = append(errs, b.sendPaneInit(firstWindow)...)
	if ds.Command != "" {
		errs = append(errs, b.client.SendKeys(firstWindow, ds.Command))
	}

	return errors.Join(errs...)
}

func (b *Builder) buildAdHoc(name, root string, ds *config.DefaultSession) error {
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
	for i := 0; i < b.adHocLayout.Panes; i++ {
		paneTarget := fmt.Sprintf("%s.%d", target, pb+i)
		errs = append(errs, b.sendPaneInit(paneTarget)...)
	}

	if ds != nil && ds.Command != "" {
		errs = append(errs, b.client.SendKeys(fmt.Sprintf("%s.%d", target, pb), ds.Command))
	}

	errs = append(errs, b.client.SelectPane(name, fw, pb))

	return errors.Join(errs...)
}

func (b *Builder) buildCommand(name string, cfg *config.ProjectConfig, root string) error {
	if err := b.createDetachedSession(name, root, ""); err != nil {
		return err
	}

	firstWindow := name + ":" + b.firstWindow()
	var errs []error

	errs = append(errs, b.applySessionSettings(name, cfg, firstWindow)...)
	errs = append(errs, b.sendPaneInit(firstWindow)...)
	errs = append(errs, b.client.SendKeys(firstWindow, cfg.Command))
	errs = append(errs, b.sendOnReady(cfg, firstWindow)...)

	return errors.Join(errs...)
}

func (b *Builder) buildWindowed(name string, cfg *config.ProjectConfig, root string) error {
	firstWin := cfg.Windows[0]
	winRoot := windowRoot(firstWin.Root, root)

	if err := b.createDetachedSession(name, winRoot, firstWin.Name); err != nil {
		return err
	}

	firstWindow := name + ":" + firstWin.Name
	var errs []error

	errs = append(errs, b.applySessionSettings(name, cfg, firstWindow)...)
	errs = append(errs, b.startWindow(name, firstWin, root))

	for _, w := range cfg.Windows[1:] {
		wr := windowRoot(w.Root, root)
		errs = append(errs, b.client.NewWindow(name, NewWindowOpts{Name: w.Name, Root: wr}))
		errs = append(errs, b.startWindow(name, w, root))
	}

	errs = append(errs, b.client.SelectWindow(name, firstWin.Name))
	errs = append(errs, b.sendOnReady(cfg, firstWindow)...)

	return errors.Join(errs...)
}

func (b *Builder) applySessionSettings(name string, cfg *config.ProjectConfig, firstWindow string) []error {
	var errs []error

	for k, v := range cfg.Env {
		errs = append(errs, b.client.SetEnv(name, k, v))
	}

	if b.global.DefaultShell != "" {
		errs = append(errs, b.client.SetOption(name, "default-command", b.global.DefaultShell))
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

func (b *Builder) sendPaneInit(target string) []error {
	var errs []error
	for _, cmd := range b.global.PaneInit {
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
		errs = append(errs, b.client.SendKeys(target, fmt.Sprintf("export %s=%s", k, env[k])))
	}
	return errs
}

func (b *Builder) startWindow(session string, w config.Window, projectRoot string) error {
	var errs []error
	target := session + ":" + w.Name

	if w.Command != "" {
		errs = append(errs, b.sendPaneInit(target)...)
		errs = append(errs, b.sendWindowEnv(target, w.Env)...)
		errs = append(errs, b.client.SendKeys(target, w.Command))
		return errors.Join(errs...)
	}

	wr := windowRoot(w.Root, projectRoot)
	pb := b.paneBase()

	for i, p := range w.Panes {
		paneTarget := fmt.Sprintf("%s.%d", target, pb+i)

		if i > 0 {
			pr := windowRoot(p.Root, wr)
			errs = append(errs, b.client.SplitWindow(session, w.Name, SplitWindowOpts{
				Root:       pr,
				Horizontal: p.Split == "horizontal",
			}))
		}

		errs = append(errs, b.sendPaneInit(paneTarget)...)
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

func (b *Builder) BuildEphemeral(name string, command string, root string) error {
	if err := b.createDetachedSession(name, root, ""); err != nil {
		return err
	}

	fw := b.firstWindow()
	target := name + ":" + fw

	if b.adHocLayout != nil {
		var errs []error
		for i := 1; i < b.adHocLayout.Panes; i++ {
			errs = append(errs, b.client.SplitWindow(name, fw, SplitWindowOpts{Root: root}))
		}
		if b.adHocLayout.Layout != "" {
			errs = append(errs, b.client.SelectLayout(name, fw, b.adHocLayout.Layout))
		}
		pb := b.paneBase()
		errs = append(errs, b.client.SendKeys(fmt.Sprintf("%s.%d", target, pb), command))
		errs = append(errs, b.client.SelectPane(name, fw, pb))
		return errors.Join(errs...)
	}

	return b.client.SendKeys(target, command)
}

func (b *Builder) StopSession(name string) error {
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

	return b.startWindow(session, w, root)
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

func windowRoot(winRoot, projectRoot string) string {
	if winRoot != "" {
		if strings.HasPrefix(winRoot, "/") || strings.HasPrefix(winRoot, "~") {
			return winRoot
		}
		return filepath.Join(projectRoot, winRoot)
	}
	return projectRoot
}
