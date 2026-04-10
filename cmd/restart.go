package cmd

import "github.com/spf13/cobra"

var restartCmd = &cobra.Command{
	Use:   "restart <target> [target ...]",
	Short: "Restart one or more tmux sessions",
	Long: `Stop and start tmux session(s), picking up any config changes.

Supports glob patterns with +, group expansion with @, and project:window
syntax to restart one or more windows (comma-separated) inside a session
without tearing down the rest.`,
	Example: `  nux restart blog
  nux restart blog api
  nux restart web+
  nux restart @work
  nux restart blog --var port=9090
  nux restart blog:editor
  nux restart blog:editor,server`,
	Args: cobra.MinimumNArgs(1),
	RunE: runRestart,
}

func init() {
	rootCmd.AddCommand(restartCmd)
}

func runRestart(_ *cobra.Command, args []string) error {
	d, err := setup()
	if err != nil {
		return err
	}
	return runRestartWith(d, args)
}

func runRestartWith(d *deps, args []string) error {
	targets, err := expandArgs(d, args)
	if err != nil {
		return err
	}

	for i, t := range targets {
		result, err := d.resolver.Resolve(t.Project)
		if err != nil {
			return err
		}

		if err := applyVarOverrides(d, result.Config); err != nil {
			return err
		}

		if t.Windows != nil {
			for _, w := range t.Windows {
				if err := d.builder.RestartWindow(result.Name, w, result.Config, result.Root); err != nil {
					return err
				}
			}
		} else {
			if err := d.builder.RestartSession(result.Name, result.Config, result.Root); err != nil {
				return err
			}
		}

		isLast := i == len(targets)-1
		if !d.noAttach && isLast {
			return d.client.AttachSession(result.Name)
		}
	}
	return nil
}
