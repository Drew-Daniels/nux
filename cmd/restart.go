package cmd

import "github.com/spf13/cobra"

var restartCmd = &cobra.Command{
	Use:   "restart <session>",
	Short: "Restart a tmux session",
	Long: `Stop and start a tmux session, picking up any config changes.

Supports :window syntax to restart individual windows without touching the
rest of the session.`,
	Example: `  nux restart blog
  nux restart blog:editor`,
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
	arg := args[0]

	projectName, windowName := ParseTarget(arg)

	result, err := d.resolver.Resolve(projectName)
	if err != nil {
		return err
	}

	if windowName != "" {
		if err := d.builder.RestartWindow(result.Name, windowName, result.Config, result.Root); err != nil {
			return err
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
