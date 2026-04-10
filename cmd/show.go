package cmd

import (
	"fmt"

	"github.com/Drew-Daniels/nux/internal/config"
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
	name := args[0]
	result, err := d.resolver.Resolve(name)
	if err != nil {
		return err
	}

	if err := applyVarOverrides(d, result.Config); err != nil {
		return err
	}

	out := showOutput{
		Name:   result.Name,
		Root:   result.Root,
		Source: result.ConfigSource,
		Config: result.Config,
	}

	data, err := yaml.Marshal(out)
	if err != nil {
		return fmt.Errorf("marshalling output: %w", err)
	}
	_, _ = fmt.Fprint(d.stdout, string(data))
	return nil
}
