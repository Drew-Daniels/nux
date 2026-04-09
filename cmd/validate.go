package cmd

import (
	"fmt"

	"github.com/Drew-Daniels/nux/internal/config"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [name]",
	Short: "Validate project config syntax",
	Long:  `Validate project config files for structural errors. Validates all configs if no name is given.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func runValidate(_ *cobra.Command, args []string) error {
	d, err := setup()
	if err != nil {
		return err
	}
	return runValidateWith(d, args)
}

func runValidateWith(d *deps, args []string) error {
	if len(args) > 0 {
		return validateProject(d, args[0])
	}

	return validateAll(d)
}

func validateAll(d *deps) error {
	projects, err := d.store.List()
	if err != nil {
		return fmt.Errorf("listing projects: %w", err)
	}
	if len(projects) == 0 {
		_, _ = fmt.Fprintln(d.stdout, "No project configs found.")
		return nil
	}

	hasErrors := false
	for _, p := range projects {
		cfg, _, err := d.store.Load(p.Name)
		if err != nil {
			_, _ = fmt.Fprintf(d.stderr, "  [error] %s: %v\n", p.Name, err)
			hasErrors = true
			continue
		}
		if errs := config.Validate(cfg); len(errs) > 0 {
			for _, e := range errs {
				_, _ = fmt.Fprintf(d.stderr, "  [error] %s: %v\n", p.Name, e)
			}
			hasErrors = true
		} else {
			_, _ = fmt.Fprintf(d.stdout, "  [ok]    %s\n", p.Name)
		}
	}
	if hasErrors {
		return fmt.Errorf("one or more configs have errors")
	}
	return nil
}

func validateProject(d *deps, name string) error {
	cfg, _, err := d.store.Load(name)
	if err != nil {
		return fmt.Errorf("loading config %q: %w", name, err)
	}

	errs := config.Validate(cfg)
	if len(errs) == 0 {
		_, _ = fmt.Fprintf(d.stdout, "  [ok]    %s\n", name)
		return nil
	}

	for _, e := range errs {
		_, _ = fmt.Fprintf(d.stderr, "  [error] %s: %v\n", name, e)
	}
	return fmt.Errorf("config %q has errors", name)
}
