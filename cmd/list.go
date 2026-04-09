package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Drew-Daniels/nux/internal/config"
	"github.com/Drew-Daniels/nux/internal/tmux"
	"github.com/Drew-Daniels/nux/internal/ui"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List available projects",
	Long:    `Show all available projects with their config source and session status.`,
	RunE:    runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(_ *cobra.Command, _ []string) error {
	d, err := setup()
	if err != nil {
		return err
	}
	return runListWith(d)
}

func runListWith(d *deps) error {
	projects, err := d.store.List()
	if err != nil {
		return fmt.Errorf("listing projects: %w", err)
	}

	sessions, err := d.client.ListSessions()
	if err != nil {
		sessions = nil
	}

	sessionMap := make(map[string]tmux.SessionInfo, len(sessions))
	for _, s := range sessions {
		sessionMap[s.Name] = s
	}

	configs := make(map[string]*config.ProjectConfig, len(projects))
	for _, p := range projects {
		cfg, _, err := d.store.Load(p.Name)
		if err != nil {
			_, _ = fmt.Fprintf(d.stderr, "warning: loading %s: %v\n", p.Name, err)
			continue
		}
		configs[p.Name] = cfg
	}

	tbl := &ui.Table{Headers: []string{"NAME", "STATUS", "WINDOWS", "UPTIME", "CONFIG", "ROOT"}}

	for _, p := range projects {
		status := "-"
		windows := ""
		uptime := ""
		sessionName := config.NormalizeSessionName(p.Name)
		if s, ok := sessionMap[sessionName]; ok {
			status = "running"
			windows = strconv.Itoa(s.Windows)
			uptime = formatDuration(time.Since(s.Created))
		}

		root := ""
		if cfg := configs[p.Name]; cfg != nil {
			root = cfg.Root
		}

		tbl.Rows = append(tbl.Rows, []string{p.Name, status, windows, uptime, "project", root})
	}

	_, _ = fmt.Fprintln(d.stdout, tbl.Render())
	return nil
}
