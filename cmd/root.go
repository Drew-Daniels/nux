package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Drew-Daniels/nux/internal/config"
	"github.com/Drew-Daniels/nux/internal/picker"
	"github.com/Drew-Daniels/nux/internal/resolver"
	"github.com/Drew-Daniels/nux/internal/tmux"
	"github.com/Drew-Daniels/nux/internal/ui"
	"github.com/spf13/cobra"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

type options struct {
	run         string
	layout      string
	panes       int
	vars        []string
	noAttach    bool
	dryRun      bool
	force       bool
	deleteForce bool
	copyForce   bool
	configDir   string
	projectDirs string
	dir         string
	editorFunc  func() string
}

func (o *options) editor() string {
	if o.editorFunc != nil {
		return o.editorFunc()
	}
	return os.Getenv("EDITOR")
}

// opts is package-level because cobra binds flags to persistent variables.
// All flag reads happen in setup(); command handlers read from deps only.
var opts options

type deps struct {
	global        *config.GlobalConfig
	client        tmux.Client
	builder       *tmux.Builder
	resolver      *resolver.Resolver
	store         config.ProjectStore
	projectCfgDir string

	noAttach    bool
	force       bool
	deleteForce bool
	copyForce   bool
	run         string
	layout      string
	panes       int
	dir         string
	editor      string
	vars        map[string]string
	stdin       io.Reader
	stdout      io.Writer
	stderr      io.Writer

	getwd        func() (string, error)
	confirm      func(prompt string) (bool, error)
	openEditor   func(path string) error
	newPicker    func(backend string, stderr io.Writer) (picker.Picker, error)
	execCmd      func(name string, arg ...string) *exec.Cmd
	help         func() error
	checkBin     func(name string) (path string, ok bool)
	probeVersion func() (string, error)
	checkStat    func(path string) (os.FileInfo, error)
}

var rootCmd = &cobra.Command{
	Use:   "nux [project ...] [flags]",
	Short: "A modern tmux session manager",
	Long: `nux manages tmux sessions declaratively through project configs.

Start sessions by name, attach to running sessions, or batch-start entire
groups. Projects without explicit configs use the default session template.

When run with no arguments inside a project directory, nux starts or attaches
to a session for that directory. When run with no arguments outside a project
directory, nux opens an interactive picker (if configured).`,
	Example: `  nux                           # picker or current directory session
  nux blog                      # start/attach the "blog" session
  nux blog api docs             # start multiple sessions
  nux @work                     # start all sessions in the "work" group
  nux web+                      # start all projects matching "web*"
  nux blog:editor               # start only the "editor" window
  nux blog:editor,server        # start only those windows (in this order)
  nux -x "just dev"             # run a command in the current directory
  nux -x "fish" blog            # run fish in the blog session
  nux -l tiled -p 4             # 4 equal panes in the current directory
  nux myproject -l main-vertical -p 3
  nux -x "fish" -l tiled -p 4 blog  # fish in each of 4 tiled panes`,
	Args:              cobra.ArbitraryArgs,
	DisableAutoGenTag: true,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	SilenceErrors:     true,
	SilenceUsage:      true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.RunE = runRoot
	rootCmd.Flags().StringVarP(&opts.run, "run", "x", "", "run a command instead of the project config")
	rootCmd.Flags().StringVarP(&opts.layout, "layout", "l", "", "ad-hoc tmux layout (e.g. tiled, even-horizontal)")
	_ = rootCmd.RegisterFlagCompletionFunc("layout", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{
			"even-horizontal\thorizontal split",
			"even-vertical\tvertical split",
			"main-horizontal\tlarge top pane",
			"main-vertical\tlarge left pane",
			"tiled\tequal grid",
		}, cobra.ShellCompDirectiveNoFileComp
	})
	rootCmd.Flags().IntVarP(&opts.panes, "panes", "p", 0, "number of panes for the ad-hoc layout")
	rootCmd.Flags().StringArrayVar(&opts.vars, "var", nil, "override a custom variable (key=value, repeatable)")
	rootCmd.Flags().BoolVar(&opts.noAttach, "no-attach", false, "start session(s) without attaching")
	rootCmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "print tmux commands without executing (still queries tmux for session state)")
	rootCmd.Flags().BoolVar(&opts.force, "force", false, "override nested session prevention")
	rootCmd.Flags().StringVar(&opts.configDir, "config-dir", "", "override config directory path (global config and project configs)")
	rootCmd.Flags().StringVar(&opts.projectDirs, "project-dirs", "", "override project directories path")
	rootCmd.Flags().StringVarP(&opts.dir, "dir", "C", "", "session root directory (skips name-based resolution)")

	restartCmd.Flags().AddFlag(rootCmd.Flag("no-attach"))
	restartCmd.Flags().StringArrayVar(&opts.vars, "var", nil, "override a custom variable (key=value, repeatable)")

	rootCmd.AddCommand(stopAllCmd)
}

func loadGlobalConfig() (*config.GlobalConfig, string, error) {
	if opts.configDir != "" {
		path := filepath.Join(opts.configDir, "config.yaml")
		cfg, err := config.LoadGlobalFrom(path)
		return cfg, path, err
	}
	cfg, err := config.LoadGlobal()
	return cfg, config.GlobalConfigPath(), err
}

func setup() (*deps, error) {
	global, _, err := loadGlobalConfig()
	if err != nil {
		return nil, fmt.Errorf("loading global config: %w", err)
	}
	if opts.projectDirs != "" {
		global.ProjectDirs = config.StringOrList{opts.projectDirs}
	}

	client := tmux.NewRealClient()
	client.DryRun = opts.dryRun

	var projectCfgDir string
	if opts.configDir != "" {
		projectCfgDir = filepath.Join(opts.configDir, "projects")
	} else {
		projectCfgDir = config.ProjectConfigDir()
	}
	store := config.NewProjectStore(projectCfgDir)
	builder := tmux.NewBuilder(client, global)
	res := resolver.NewResolverWithStore(global, store)

	editor := opts.editor()
	stdin := os.Stdin
	stdout := rootCmd.OutOrStdout()
	stderr := rootCmd.ErrOrStderr()

	prompter := &ui.Prompter{In: stdin, Out: stdout}

	d := &deps{
		global:        global,
		client:        client,
		builder:       builder,
		resolver:      res,
		store:         store,
		projectCfgDir: projectCfgDir,
		noAttach:      opts.noAttach,
		force:         opts.force,
		deleteForce:   opts.deleteForce,
		copyForce:     opts.copyForce,
		run:           opts.run,
		layout:        opts.layout,
		panes:         opts.panes,
		dir:           opts.dir,
		editor:        editor,
		vars:          parseVars(opts.vars, stderr),
		stdin:         stdin,
		stdout:        stdout,
		stderr:        stderr,
		getwd:         os.Getwd,
		confirm:       prompter.Confirm,
		newPicker:     picker.New,
		execCmd:       exec.Command,
		help:          rootCmd.Help,
		checkBin:      defaultBinaryChecker,
		probeVersion:  defaultVersionProber,
		checkStat:     os.Stat,
	}
	d.openEditor = func(path string) error {
		return openInEditor(d, path)
	}
	return d, nil
}

func runRoot(_ *cobra.Command, args []string) error {
	d, err := setup()
	if err != nil {
		return err
	}

	if err := validateLayoutFlags(d); err != nil {
		return err
	}
	d.builder.SetAdHocLayout(adHocLayoutFromDeps(d))

	if d.dir != "" {
		if len(args) == 0 {
			return fmt.Errorf("--dir requires a session name")
		}
		if _, ok := matchSubcommand(args[0]); ok {
			return fmt.Errorf("--dir cannot be combined with subcommands")
		}
	} else if len(args) > 0 {
		if sub, ok := matchSubcommand(args[0]); ok {
			return sub.RunE(sub, args[1:])
		}
	}

	if d.client.IsInsideTmux() && !d.force {
		return fmt.Errorf("already inside a tmux session (use --force to override)")
	}

	if len(args) == 0 {
		return runBareNux(d)
	}

	return runSessions(d, args)
}

func matchSubcommand(prefix string) (*cobra.Command, bool) {
	var matches []*cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == prefix {
			return cmd, true
		}
		for _, alias := range cmd.Aliases {
			if alias == prefix {
				return cmd, true
			}
		}
	}
	for _, cmd := range rootCmd.Commands() {
		if strings.HasPrefix(cmd.Name(), prefix) {
			matches = append(matches, cmd)
		}
	}
	if len(matches) == 1 {
		return matches[0], true
	}
	return nil, false
}

func validateLayoutFlags(d *deps) error {
	if d.layout == "" && d.panes == 0 {
		return nil
	}

	if d.panes < 0 {
		return fmt.Errorf("--panes must be a positive number")
	}

	if d.layout != "" && !config.IsValidLayout(d.layout) {
		return fmt.Errorf("invalid --layout %q (must be even-horizontal, even-vertical, main-horizontal, main-vertical, tiled, or a custom layout string)", d.layout)
	}

	if d.layout == "" {
		d.layout = "tiled"
	}
	if d.panes == 0 {
		d.panes = 2
	}

	return nil
}

func adHocLayoutFromDeps(d *deps) *tmux.AdHocLayout {
	if d.layout == "" && d.panes == 0 && d.run == "" {
		return nil
	}
	return &tmux.AdHocLayout{Layout: d.layout, Panes: d.panes, Command: d.run}
}
