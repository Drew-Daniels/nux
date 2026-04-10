package cmd

import (
	"fmt"

	"github.com/Drew-Daniels/nux/internal/config"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create a new project config",
	Long:  `Create a new project config from the default template and open it in $EDITOR.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func runNew(_ *cobra.Command, args []string) error {
	d, err := setup()
	if err != nil {
		return err
	}
	return runNewWith(d, args)
}

func runNewWith(d *deps, args []string) error {
	name := args[0]
	path := d.store.Path(name)

	if _, _, err := d.store.Load(name); err == nil {
		return fmt.Errorf("config already exists: %s", path)
	}

	cfg := &config.ProjectConfig{
		Windows: []config.Window{
			{Name: "editor", Panes: []config.Pane{{Command: ""}}},
		},
	}

	if err := d.store.Save(name, cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	_, _ = fmt.Fprintf(d.stdout, "Created %s\n", path)

	return d.openEditor(path)
}

func openInEditor(d *deps, path string) error {
	if d.editor == "" {
		_, _ = fmt.Fprintln(d.stderr, "hint: set $EDITOR to open new configs automatically")
		return nil
	}
	cmd := d.execCmd(d.editor, path)
	cmd.Stdin = d.stdin
	cmd.Stdout = d.stdout
	cmd.Stderr = d.stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editor failed: %w", err)
	}
	return nil
}
