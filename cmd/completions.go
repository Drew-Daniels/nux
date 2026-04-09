package cmd

import (
	"github.com/spf13/cobra"
)

var completionsCmd = &cobra.Command{
	Use:       "completions <bash|zsh|fish>",
	Short:     "Generate shell completions",
	Long:      `Generate shell completion scripts for bash, zsh, or fish.`,
	ValidArgs: []string{"bash", "zsh", "fish"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE:      runCompletions,
	Example: `  nux completions bash > /etc/bash_completion.d/nux
  nux completions zsh > "${fpath[1]}/_nux"
  nux completions fish > ~/.config/fish/completions/nux.fish`,
}

func init() {
	rootCmd.AddCommand(completionsCmd)
}

func runCompletions(_ *cobra.Command, args []string) error {
	d, err := setup()
	if err != nil {
		return err
	}
	return runCompletionsWith(d, args)
}

func runCompletionsWith(d *deps, args []string) error {
	switch args[0] {
	case "bash":
		return rootCmd.GenBashCompletion(d.stdout)
	case "zsh":
		return rootCmd.GenZshCompletion(d.stdout)
	case "fish":
		return rootCmd.GenFishCompletion(d.stdout, true)
	}
	return nil
}
