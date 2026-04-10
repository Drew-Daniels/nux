package cmd

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/Drew-Daniels/nux/internal/config"
	"github.com/Drew-Daniels/nux/internal/resolver"
)

func runSessions(d *deps, args []string) error {
	if d.run != "" && len(d.vars) > 0 {
		_, _ = fmt.Fprintln(d.stderr, "warning: --var is ignored when using --run/-x")
	}

	targets, err := expandArgs(d, args)
	if err != nil {
		return err
	}

	if err := validateAdhocSubset(d, targets); err != nil {
		return err
	}

	for i, arg := range targets {
		result, err := d.resolver.Resolve(arg.Project)
		if err != nil {
			return err
		}

		cfg := effectiveConfig(d, result)

		if err := applyVarOverrides(d, cfg); err != nil {
			return err
		}

		if arg.Windows != nil {
			if err := validateSubsetConfig(cfg, result.Name); err != nil {
				return err
			}
			if !d.client.HasSession(result.Name) {
				if err := d.builder.BuildWindows(result.Name, cfg, result.Root, arg.Windows); err != nil {
					return fmt.Errorf("building session %q: %w", result.Name, err)
				}
			} else {
				focus := arg.Windows[0]
				if err := d.client.SelectWindow(result.Name, focus); err != nil {
					return fmt.Errorf("selecting window %q: %w", focus, err)
				}
			}
		} else {
			if err := buildIfAbsent(d, result); err != nil {
				return err
			}
		}

		isLast := i == len(targets)-1
		if !d.noAttach && isLast {
			return d.client.AttachSession(result.Name)
		}
	}
	return nil
}

func effectiveConfig(d *deps, result *resolver.Result) *config.ProjectConfig {
	if d.run != "" {
		return nil
	}
	return result.Config
}

func buildIfAbsent(d *deps, result *resolver.Result) error {
	cfg := effectiveConfig(d, result)
	if d.client.HasSession(result.Name) {
		return nil
	}
	if err := d.builder.Build(result.Name, cfg, result.Root); err != nil {
		return fmt.Errorf("building session %q: %w", result.Name, err)
	}
	return nil
}

func applyVarOverrides(d *deps, cfg *config.ProjectConfig) error {
	if len(d.vars) == 0 || cfg == nil {
		return nil
	}
	if cfg.Vars == nil {
		cfg.Vars = make(map[string]string)
	}
	for k, v := range d.vars {
		cfg.Vars[k] = v
	}
	if err := d.resolver.Interpolator.InterpolateVars(cfg); err != nil {
		return fmt.Errorf("interpolation failed: %w", err)
	}
	return nil
}

func validateSubsetConfig(cfg *config.ProjectConfig, sessionName string) error {
	if cfg == nil {
		return fmt.Errorf("window selection requires a project config with windows for %q", sessionName)
	}
	if len(cfg.Windows) == 0 {
		return fmt.Errorf("project %q has no windows defined (use `nux` without :window)", sessionName)
	}
	return nil
}

func validateAdhocSubset(d *deps, targets []sessionArg) error {
	var hasSubset bool
	for _, t := range targets {
		if t.Windows != nil {
			hasSubset = true
			break
		}
	}
	if !hasSubset {
		return nil
	}
	if d.run != "" || d.layout != "" || d.panes != 0 {
		return fmt.Errorf("cannot combine --run, --layout, or --panes with project:window targets")
	}
	return nil
}

func expandArgs(d *deps, args []string) ([]sessionArg, error) {
	var sessionNames []string
	for _, arg := range args {
		if strings.Contains(arg, "+") {
			sessions, _ := d.client.ListSessions()
			for _, s := range sessions {
				sessionNames = append(sessionNames, s.Name)
			}
			break
		}
	}

	var targets []sessionArg
	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "@"):
			group := strings.TrimPrefix(arg, "@")
			members, err := d.resolver.ExpandGroup(group)
			if err != nil {
				return nil, err
			}
			for _, m := range members {
				targets = append(targets, sessionArg{Project: m})
			}

		case strings.Contains(arg, "+"):
			matches, err := d.resolver.ExpandGlob(arg, sessionNames)
			if err != nil {
				return nil, err
			}
			for _, m := range matches {
				targets = append(targets, sessionArg{Project: m})
			}

		case strings.Contains(arg, ":"):
			sa, err := parseSessionToken(arg)
			if err != nil {
				return nil, err
			}
			targets = append(targets, sa)

		default:
			targets = append(targets, sessionArg{Project: arg})
		}
	}
	return targets, nil
}

func tryAutoDetect(d *deps) (*resolver.Result, bool) {
	cwd, _ := d.getwd()
	projectsDir := resolver.ResolveRoot(d.global.ProjectsDir, "")

	rel, err := filepath.Rel(projectsDir, cwd)
	if err != nil || strings.HasPrefix(rel, "..") {
		return nil, false
	}

	name := strings.SplitN(rel, string(filepath.Separator), 2)[0]
	result, err := d.resolver.Resolve(name)
	if err != nil {
		return nil, false
	}
	return result, true
}

func runBareNux(d *deps) error {
	if result, ok := tryAutoDetect(d); ok {
		if err := buildIfAbsent(d, result); err != nil {
			return err
		}
		if !d.noAttach {
			return d.client.AttachSession(result.Name)
		}
		return nil
	}

	if d.global.PickerOnBare {
		items := collectPickerItems(d)
		if len(items) == 0 {
			return fmt.Errorf("no projects or sessions found")
		}
		p, err := d.newPicker(d.global.Picker, d.stderr)
		if err != nil {
			return err
		}
		chosen, err := p.Pick(items, "project")
		if err != nil {
			return err
		}
		if chosen == "" {
			return nil
		}
		return runSessions(d, []string{chosen})
	}

	return d.help()
}

func collectPickerItems(d *deps) []string {
	// Dedupe by normalized session name so e.g. config "my.project" and tmux
	// session "my_project" count as one entry (project name preferred).
	seen := make(map[string]bool)
	var items []string

	projects, _ := d.store.List()
	for _, p := range projects {
		k := config.NormalizeSessionName(p.Name)
		if !seen[k] {
			seen[k] = true
			items = append(items, p.Name)
		}
	}

	sessions, _ := d.client.ListSessions()
	for _, s := range sessions {
		k := config.NormalizeSessionName(s.Name)
		if !seen[k] {
			seen[k] = true
			items = append(items, s.Name)
		}
	}

	return items
}

func parseVars(raw []string, w io.Writer) map[string]string {
	vars := make(map[string]string, len(raw))
	for _, v := range raw {
		k, val, ok := strings.Cut(v, "=")
		if !ok {
			_, _ = fmt.Fprintf(w, "warning: ignoring malformed --var %q (expected key=value)\n", v)
			continue
		}
		vars[k] = val
	}
	return vars
}
