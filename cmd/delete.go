package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a project config",
	Long:  `Delete a project config file. Prompts for confirmation unless --force is set.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDelete,
}

func init() {
	deleteCmd.Flags().BoolVar(&opts.deleteForce, "force", false, "skip confirmation prompt")
	rootCmd.AddCommand(deleteCmd)
}

func runDelete(_ *cobra.Command, args []string) error {
	d, err := setup()
	if err != nil {
		return err
	}
	return runDeleteWith(d, args)
}

func runDeleteWith(d *deps, args []string) error {
	name := args[0]
	path := d.store.Path(name)

	if _, _, err := d.store.Load(name); err != nil {
		return fmt.Errorf("config not found: %s", path)
	}

	if !d.deleteForce {
		ok, err := d.confirm(fmt.Sprintf("Delete config for %q?", name))
		if err != nil {
			return err
		}
		if !ok {
			_, _ = fmt.Fprintln(d.stdout, "Cancelled.")
			return nil
		}
	}

	if err := d.store.Delete(name); err != nil {
		return fmt.Errorf("deleting config: %w", err)
	}
	_, _ = fmt.Fprintf(d.stdout, "Deleted config for %q\n", name)
	return nil
}
