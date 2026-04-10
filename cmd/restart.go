package cmd

import "github.com/spf13/cobra"

var restartCmd = &cobra.Command{
	Use:   "restart <session>",
	Short: "Restart a tmux session",
	Long: `Stop and start a tmux session, picking up any config changes.

Supports project:window syntax to restart one or more windows (comma-separated)
without tearing down the rest of the session.`,
	Example: `  nux restart blog
  nux restart blog --var port=9090
  nux restart blog:editor
  nux restart blog:editor,server`,
	Args: cobra.ExactArgs(1),
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
	sa, err := parseSessionToken(args[0])
	if err != nil {
		return err
	}

	result, err := d.resolver.Resolve(sa.Project)
	if err != nil {
		return err
	}

	if err := applyVarOverrides(d, result.Config); err != nil {
		return err
	}

	if sa.Windows != nil {
		for _, w := range sa.Windows {
			if err := d.builder.RestartWindow(result.Name, w, result.Config, result.Root); err != nil {
				return err
			}
		}
	} else {
		if err := d.builder.RestartSession(result.Name, result.Config, result.Root); err != nil {
			return err
		}
	}

	if !d.noAttach {
		return d.client.AttachSession(result.Name)
	}
	return nil
}
