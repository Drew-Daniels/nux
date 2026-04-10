package cmd

import (
	"fmt"

	"github.com/Drew-Daniels/nux/internal/config"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop <session> [session ...]",
	Short: "Stop one or more tmux sessions",
	Long: `Stop one or more tmux sessions managed by nux.

Supports glob patterns with + and group expansion with @.`,
	Example: `  nux stop blog
  nux stop blog api docs
  nux stop web+
  nux stop @work`,
	Args: cobra.MinimumNArgs(1),
	RunE: runStop,
}

var stopAllCmd = &cobra.Command{
	Use:   "stop-all",
	Short: "Stop all running tmux sessions",
	RunE:  runStopAll,
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func runStop(_ *cobra.Command, args []string) error {
	d, err := setup()
	if err != nil {
		return err
	}
	return runStopWith(d, args)
}

func runStopWith(d *deps, args []string) error {
	targets, err := expandArgs(d, args)
	if err != nil {
		return err
	}

	for _, arg := range targets {
		normalized := config.NormalizeSessionName(arg.Project)
		if !d.client.HasSession(normalized) {
			return fmt.Errorf("session %q is not running", arg.Project)
		}
		if err := d.builder.StopSession(normalized); err != nil {
			return fmt.Errorf("stopping session %q: %w", arg.Project, err)
		}
	}
	return nil
}

func runStopAll(_ *cobra.Command, _ []string) error {
	d, err := setup()
	if err != nil {
		return err
	}
	return runStopAllWith(d)
}

func runStopAllWith(d *deps) error {
	return d.builder.StopAll()
}
