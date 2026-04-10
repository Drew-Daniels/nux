package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Drew-Daniels/nux/internal/config"
	"github.com/Drew-Daniels/nux/internal/resolver"
	"github.com/spf13/cobra"
)

type binaryChecker func(name string) (path string, ok bool)

type versionProber func() (string, error)

type statFunc func(path string) (os.FileInfo, error)

func defaultBinaryChecker(name string) (string, bool) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", false
	}
	return path, true
}

func defaultVersionProber() (string, error) {
	out, err := exec.Command("tmux", "-V").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check environment and configuration",
	Long:  `Run diagnostic checks on the nux environment and report issues with suggested fixes.`,
	RunE:  runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func runDoctor(_ *cobra.Command, _ []string) error {
	d, err := setup()
	if err != nil {
		return err
	}
	return runDoctorWith(d, defaultBinaryChecker, defaultVersionProber, os.Stat)
}

func runDoctorWith(d *deps, checkBin binaryChecker, probeVersion versionProber, checkStat statFunc) error {
	out := d.stdout

	_, _ = fmt.Fprintf(out, "  nux %s (%s/%s)\n\n", Version, runtime.GOOS, runtime.GOARCH)
	ok := true

	if path, found := checkBin("tmux"); !found {
		_, _ = fmt.Fprintf(out, "  [missing] tmux\n")
		ok = false
	} else {
		_, _ = fmt.Fprintf(out, "  [ok]      tmux (%s)\n", path)
		if ver, err := probeVersion(); err == nil {
			_, _ = fmt.Fprintf(out, "            %s\n", ver)
		}
	}

	_, _ = fmt.Fprintf(out, "  [ok]      global config\n")

	ok = runDoctorChecks(d, checkBin, checkStat) && ok

	_, _ = fmt.Fprintln(out)
	if ok {
		_, _ = fmt.Fprintln(out, "All checks passed.")
		return nil
	}
	return fmt.Errorf("some checks failed")
}

func runDoctorChecks(d *deps, checkBin binaryChecker, checkStat statFunc) bool {
	out := d.stdout
	errOut := d.stderr
	global := d.global
	ok := true

	if global.Zoxide {
		if _, found := checkBin("zoxide"); !found {
			_, _ = fmt.Fprintf(out, "  [missing] zoxide\n")
			ok = false
		}
	}

	if global.Picker != "" {
		if _, found := checkBin(global.Picker); !found {
			_, _ = fmt.Fprintf(out, "  [warn]    picker %q not found, interactive selection will fail\n", global.Picker)
		}
	}

	if _, err := checkStat(d.projectCfgDir); err != nil {
		_, _ = fmt.Fprintf(out, "  [warn]    config directory missing: %s\n", d.projectCfgDir)
	} else {
		_, _ = fmt.Fprintf(out, "  [ok]      config directory (%s)\n", d.projectCfgDir)
	}

	projectsDir := resolver.ResolveRoot(global.ProjectsDir, "")
	if _, err := checkStat(projectsDir); err != nil {
		_, _ = fmt.Fprintf(out, "  [warn]    projects directory missing: %s\n", projectsDir)
	} else {
		_, _ = fmt.Fprintf(out, "  [ok]      projects directory (%s)\n", projectsDir)
	}

	results, err := config.ValidateAllWith(d.store)
	if err != nil {
		_, _ = fmt.Fprintf(errOut, "  [warn]    cannot list project configs: %v\n", err)
		return ok
	}

	_, _ = fmt.Fprintf(out, "  [ok]      %d project config(s)\n", len(results))

	invalid := 0
	for _, r := range results {
		if len(r.Errors) > 0 {
			for _, e := range r.Errors {
				_, _ = fmt.Fprintf(errOut, "  [fail]    %s: %v\n", r.Name, e)
			}
			invalid++
		}
	}
	if invalid > 0 {
		ok = false
	} else if len(results) > 0 {
		_, _ = fmt.Fprintf(out, "  [ok]      all configs valid\n")
	}

	return ok
}
