package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var resetForce bool
var resetProjects bool

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Remove the global config and optionally all project configs",
	Long: `Remove nux's global config file so you can start fresh. Use --projects to
also remove all project config files.

Before anything is deleted, a summary of what will be removed is printed.
Use --force to skip the confirmation prompt.

Running sessions are not affected. To recreate the global config after a
reset, run nux config.`,
	Args: cobra.NoArgs,
	RunE: runReset,
	Example: `  nux reset
  nux reset --force
  nux reset --projects --force`,
}

func init() {
	resetCmd.Flags().BoolVar(&resetForce, "force", false, "skip confirmation prompt")
	resetCmd.Flags().BoolVar(&resetProjects, "projects", false, "also remove all project configs")
	rootCmd.AddCommand(resetCmd)
}

func runReset(_ *cobra.Command, _ []string) error {
	d, err := setup()
	if err != nil {
		return err
	}
	return runResetWith(d)
}

func runResetWith(d *deps) error {
	cfgDir := d.projectCfgDir
	cfgPath := filepath.Join(filepath.Dir(cfgDir), "config.yaml")

	if _, err := d.checkStat(cfgPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config not found: %s", cfgPath)
		}
		return fmt.Errorf("config: %w", err)
	}

	projectsExist := false
	if info, err := d.checkStat(cfgDir); err == nil && info.IsDir() {
		entries, _ := os.ReadDir(cfgDir)
		for _, e := range entries {
			if !e.IsDir() {
				projectsExist = true
				break
			}
		}
	}

	_, _ = fmt.Fprintln(d.stdout, "Will remove:")
	_, _ = fmt.Fprintf(d.stdout, "  config:   %s\n", cfgPath)
	if resetProjects && projectsExist {
		_, _ = fmt.Fprintf(d.stdout, "  projects: %s\n", cfgDir)
	}
	_, _ = fmt.Fprintln(d.stdout)
	_, _ = fmt.Fprintln(d.stdout, "Will keep:")
	if !resetProjects && projectsExist {
		_, _ = fmt.Fprintf(d.stdout, "  projects: %s (use --projects to remove)\n", cfgDir)
	}
	_, _ = fmt.Fprintln(d.stdout, "  Running tmux sessions are not affected.")
	_, _ = fmt.Fprintln(d.stdout)

	if !resetForce {
		ok, err := d.confirm("Proceed?")
		if err != nil {
			return err
		}
		if !ok {
			_, _ = fmt.Fprintln(d.stdout, "Cancelled.")
			return nil
		}
	}

	if err := os.Remove(cfgPath); err != nil {
		return fmt.Errorf("removing config: %w", err)
	}

	if resetProjects {
		if err := os.RemoveAll(cfgDir); err != nil {
			return fmt.Errorf("removing project configs: %w", err)
		}
	}

	_, _ = fmt.Fprintln(d.stdout, "Done.")
	return nil
}
