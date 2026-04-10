package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Drew-Daniels/nux/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Open or create the global config file",
	Long: `Open the global nux config in $EDITOR. If the config does not exist yet,
a scaffold with commented examples is created first.`,
	Args: cobra.NoArgs,
	RunE: runConfig,
	Example: `  nux config
  nux config --config-dir ~/dotfiles/nux`,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfig(_ *cobra.Command, _ []string) error {
	d, err := setup()
	if err != nil {
		cfgDir := config.DefaultConfigDir()
		if opts.configDir != "" {
			cfgDir = opts.configDir
		}
		cfgPath := filepath.Join(cfgDir, "config.yaml")
		editor := opts.editor()
		if editor == "" {
			return fmt.Errorf("loading global config failed and $EDITOR is not set: %w", err)
		}
		_, _ = fmt.Fprintf(os.Stderr, "warning: %v\n", err)
		cmd := exec.Command(editor, cfgPath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	cfgDir := config.DefaultConfigDir()
	if opts.configDir != "" {
		cfgDir = opts.configDir
	}

	return runConfigWith(d, cfgDir)
}

func runConfigWith(d *deps, cfgDir string) error {
	cfgPath := filepath.Join(cfgDir, "config.yaml")

	if _, err := d.checkStat(cfgPath); err == nil {
		return d.openEditor(cfgPath)
	}

	projectsDir := filepath.Join(cfgDir, "projects")
	if err := os.MkdirAll(projectsDir, 0o755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	if err := os.WriteFile(cfgPath, config.ScaffoldGlobalConfig(), 0o600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	_, _ = fmt.Fprintf(d.stdout, "Created %s\n", cfgPath)

	return d.openEditor(cfgPath)
}
