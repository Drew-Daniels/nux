package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	RunE:  runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(_ *cobra.Command, _ []string) error {
	d, err := setup()
	if err != nil {
		return err
	}
	return runVersionWith(d)
}

func runVersionWith(d *deps) error {
	_, _ = fmt.Fprintf(d.stdout, "nux %s\n", Version)
	_, _ = fmt.Fprintf(d.stdout, "  commit: %s\n", Commit)
	_, _ = fmt.Fprintf(d.stdout, "  built:  %s\n", Date)
	return nil
}
