package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{
	Use:     "copy <source> <destination>",
	Aliases: []string{"cp"},
	Short:   "Copy an existing project config to a new name",
	Long:    `Copy an existing project config file to a new name and open it in $EDITOR.`,
	Args:    cobra.ExactArgs(2),
	RunE:    runCopy,
}

func init() {
	copyCmd.Flags().BoolVar(&opts.copyForce, "force", false, "overwrite destination if it already exists")
	rootCmd.AddCommand(copyCmd)
}

func runCopy(_ *cobra.Command, args []string) error {
	d, err := setup()
	if err != nil {
		return err
	}
	return runCopyWith(d, args)
}

func runCopyWith(d *deps, args []string) error {
	source := args[0]
	dest := args[1]

	data, sourcePath, err := d.store.LoadRaw(source)
	if err != nil {
		return fmt.Errorf("config not found: %s", d.store.Path(source))
	}

	destPath := d.store.Path(dest)
	if _, _, err := d.store.Load(dest); err == nil && !d.copyForce {
		return fmt.Errorf("config already exists: %s (use --force to overwrite)", destPath)
	}

	if err := d.store.SaveRaw(dest, data); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	_, _ = fmt.Fprintf(d.stdout, "Copied %s -> %s\n", sourcePath, destPath)

	if err := d.openEditor(destPath); err != nil {
		return err
	}

	return validateProjectAfterEdit(d, dest)
}
