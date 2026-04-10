package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Drew-Daniels/nux/internal/ui"
	"github.com/spf13/cobra"
)

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Show running tmux sessions",
	Long:  `Show all currently running tmux sessions.`,
	RunE:  runPs,
}

func init() {
	rootCmd.AddCommand(psCmd)
}

func runPs(_ *cobra.Command, _ []string) error {
	d, err := setup()
	if err != nil {
		return err
	}
	return runPsWith(d)
}

func runPsWith(d *deps) error {
	sessions, err := d.client.ListSessions()
	if err != nil {
		return fmt.Errorf("listing sessions: %w", err)
	}
	if len(sessions) == 0 {
		_, _ = fmt.Fprintln(d.stdout, "No running sessions.")
		return nil
	}

	tbl := &ui.Table{Headers: []string{"NAME", "WINDOWS", "ATTACHED", "UPTIME"}}

	for _, s := range sessions {
		attached := "no"
		if s.Attached {
			attached = "yes"
		}
		tbl.Rows = append(tbl.Rows, []string{
			s.Name,
			strconv.Itoa(s.Windows),
			attached,
			formatDuration(time.Since(s.Created)),
		})
	}

	_, _ = fmt.Fprintln(d.stdout, tbl.Render())
	return nil
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}
