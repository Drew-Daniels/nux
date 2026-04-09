package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var showCmd = &cobra.Command{
	Use:   "show <project>",
	Short: "Print the resolved config for a project",
	Long: `Print the fully resolved config for a project after interpolation,
variable expansion, and root resolution. Useful for debugging configs.`,
	Example: `  nux show blog
  nux show blog --var port=8080`,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		d, err := setup()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		projects, err := d.store.List()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		names := make([]string, len(projects))
		for i, p := range projects {
			names[i] = p.Name
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: runShow,
}

func init() {
	showCmd.Flags().StringArrayVar(&opts.vars, "var", nil, "override a custom variable (key=value, repeatable)")
	rootCmd.AddCommand(showCmd)
}

type showOutput struct {
	Name   string            `yaml:"name"`
	Root   string            `yaml:"root"`
	Source string            `yaml:"source"`
	Config *showOutputConfig `yaml:"config,omitempty"`
}

type showOutputConfig struct {
	Command  string             `yaml:"command,omitempty"`
	Env      map[string]string  `yaml:"env,omitempty"`
	Vars     map[string]string  `yaml:"vars,omitempty"`
	OnStart  []string           `yaml:"on_start,omitempty"`
	OnReady  []string           `yaml:"on_ready,omitempty"`
	OnDetach []string           `yaml:"on_detach,omitempty"`
	OnStop   []string           `yaml:"on_stop,omitempty"`
	Windows  []showOutputWindow `yaml:"windows,omitempty"`
}

type showOutputWindow struct {
	Name    string           `yaml:"name"`
	Root    string           `yaml:"root,omitempty"`
	Layout  string           `yaml:"layout,omitempty"`
	Command string           `yaml:"command,omitempty"`
	Panes   []showOutputPane `yaml:"panes,omitempty"`
}

type showOutputPane struct {
	Root    string `yaml:"root,omitempty"`
	Command string `yaml:"command,omitempty"`
}

func runShow(_ *cobra.Command, args []string) error {
	d, err := setup()
	if err != nil {
		return err
	}
	return runShowWith(d, args)
}

func runShowWith(d *deps, args []string) error {
	name := args[0]
	result, err := d.resolver.Resolve(name)
	if err != nil {
		return err
	}

	if len(d.vars) > 0 && result.Config != nil {
		if result.Config.Vars == nil {
			result.Config.Vars = make(map[string]string)
		}
		for k, v := range d.vars {
			result.Config.Vars[k] = v
		}
		if err := d.resolver.Interpolator.InterpolateVars(result.Config); err != nil {
			return fmt.Errorf("interpolation failed: %w", err)
		}
	}

	out := showOutput{
		Name:   result.Name,
		Root:   result.Root,
		Source: result.ConfigSource,
	}

	if result.Config != nil {
		cfg := result.Config
		sc := &showOutputConfig{
			Command:  cfg.Command,
			Env:      cfg.Env,
			Vars:     cfg.Vars,
			OnStart:  cfg.OnStart,
			OnReady:  cfg.OnReady,
			OnDetach: cfg.OnDetach,
			OnStop:   cfg.OnStop,
		}
		for _, w := range cfg.Windows {
			sw := showOutputWindow{
				Name:    w.Name,
				Root:    w.Root,
				Layout:  w.Layout,
				Command: w.Command,
			}
			for _, p := range w.Panes {
				sw.Panes = append(sw.Panes, showOutputPane{
					Root:    p.Root,
					Command: p.Command,
				})
			}
			sc.Windows = append(sc.Windows, sw)
		}
		out.Config = sc
	}

	data, err := yaml.Marshal(out)
	if err != nil {
		return fmt.Errorf("marshalling output: %w", err)
	}
	_, _ = fmt.Fprint(d.stdout, string(data))
	return nil
}
