package cmd

import (
	"fmt"

	"github.com/Drew-Daniels/nux/internal/config"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "Open a project config in $EDITOR",
	Args:  cobra.ExactArgs(1),
	RunE:  runEdit,
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func runEdit(_ *cobra.Command, args []string) error {
	d, err := setup()
	if err != nil {
		return err
	}
	return runEditWith(d, args)
}

func runEditWith(d *deps, args []string) error {
	name := args[0]
	path := d.store.Path(name)

	if _, _, err := d.store.Load(name); err != nil {
		return fmt.Errorf("config not found: %s", path)
	}

	if d.editor == "" {
		return fmt.Errorf("$EDITOR is not set")
	}

	if err := d.openEditor(path); err != nil {
		return err
	}

	return validateProjectAfterEdit(d, name)
}

func validateProjectAfterEdit(d *deps, name string) error {
	cfg, _, err := d.store.Load(name)
	if err != nil {
		_, _ = fmt.Fprintf(d.stderr, "warning: config has syntax errors: %v\n", err)
		return nil
	}

	errs := config.Validate(cfg)
	for _, e := range errs {
		_, _ = fmt.Fprintf(d.stderr, "  [error] %v\n", e)
	}
	if len(errs) == 0 {
		_, _ = fmt.Fprintln(d.stdout, "Config valid.")
	}
	return nil
}
