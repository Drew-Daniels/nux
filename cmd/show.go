package cmd

import (
	"fmt"

	"github.com/Drew-Daniels/nux/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var showRaw bool

var showCmd = &cobra.Command{
	Use:   "show <target> [target ...]",
	Short: "Print resolved config(s) for project(s)",
	Long: `Print the fully resolved config for one or more projects after interpolation,
variable expansion, and root resolution. Useful for debugging configs.

Supports glob patterns with + and group expansion with @.

Use --raw to print the config before interpolation (no variable or
environment expansion). This avoids exposing secrets in terminal output.

Multiple targets are written as a YAML stream (documents separated by ---).`,
	Example: `  nux show blog
  nux show web+
  nux show @work
  nux show blog --var port=8080
  nux show blog --raw`,
	Args: cobra.MinimumNArgs(1),
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
	showCmd.Flags().BoolVar(&showRaw, "raw", false, "print the config before interpolation (no variable or env expansion)")
	rootCmd.AddCommand(showCmd)
}

type showOutput struct {
	Name   string                `yaml:"name"`
	Root   string                `yaml:"root"`
	Source string                `yaml:"source"`
	Config *config.ProjectConfig `yaml:"config,omitempty"`
}

func runShow(_ *cobra.Command, args []string) error {
	d, err := setup()
	if err != nil {
		return err
	}
	return runShowWith(d, args)
}

func runShowWith(d *deps, args []string) error {
	targets, err := expandArgs(d, args)
	if err != nil {
		return err
	}

	for i, t := range targets {
		var data []byte
		if showRaw {
			data, err = marshalShowRawYAML(d, t.Project)
		} else {
			data, err = marshalShowResolvedYAML(d, t.Project)
		}
		if err != nil {
			return err
		}
		if i > 0 {
			_, _ = fmt.Fprint(d.stdout, "---\n")
		}
		_, _ = fmt.Fprint(d.stdout, string(data))
	}
	return nil
}

func marshalShowResolvedYAML(d *deps, name string) ([]byte, error) {
	result, err := d.resolver.Resolve(name)
	if err != nil {
		return nil, err
	}

	if err := applyVarOverrides(d, result.Config); err != nil {
		return nil, err
	}

	out := showOutput{
		Name:   result.Name,
		Root:   result.Root,
		Source: result.ConfigSource,
		Config: result.Config,
	}

	data, err := yaml.Marshal(out)
	if err != nil {
		return nil, fmt.Errorf("marshalling output: %w", err)
	}
	return data, nil
}

func marshalShowRawYAML(d *deps, name string) ([]byte, error) {
	cfg, cfgPath, err := d.store.Load(name)
	if err != nil {
		return nil, fmt.Errorf("config not found: %s", d.store.Path(name))
	}

	out := showOutput{
		Name:   config.NormalizeSessionName(name),
		Root:   cfg.Root,
		Source: cfgPath,
		Config: cfg,
	}

	data, err := yaml.Marshal(out)
	if err != nil {
		return nil, fmt.Errorf("marshalling output: %w", err)
	}
	return data, nil
}
