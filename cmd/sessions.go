package cmd

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/Drew-Daniels/nux/internal/config"
	"github.com/Drew-Daniels/nux/internal/resolver"
)

func runEphemeral(d *deps) error {
	cwd, _ := d.getwd()
	name := config.NormalizeSessionName(filepath.Base(cwd))
	if err := d.builder.BuildEphemeral(name, d.run, cwd); err != nil {
		return err
	}
	if !d.noAttach {
		return d.client.AttachSession(name)
	}
	return nil
}

func runSessions(d *deps, args []string) error {
	names, err := expandArgs(d, args)
	if err != nil {
		return err
	}

	for i, name := range names {
		result, err := d.resolver.Resolve(name)
		if err != nil {
			return err
		}

		if len(d.vars) > 0 && result.Config != nil {
			if result.Config.Vars == nil {
				result.Config.Vars = make(map[string]string)
			}
			for k, v := range d.vars {
				result.Config.Vars[k] = v
			}
			if err := d.resolver.Interpolator.InterpolateVars(result.Config); err != nil {
				return fmt.Errorf("interpolation failed: %w", err)
			}
		}

		if !d.client.HasSession(result.Name) {
			if err := d.builder.Build(result.Name, result.Config, result.Root); err != nil {
				return fmt.Errorf("building session %q: %w", result.Name, err)
			}
		}

		isLast := i == len(names)-1
		if !d.noAttach && isLast {
			return d.client.AttachSession(result.Name)
		}
	}
	return nil
}

func expandArgs(d *deps, args []string) ([]string, error) {
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

	var names []string
	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "@"):
			group := strings.TrimPrefix(arg, "@")
			members, err := d.resolver.ExpandGroup(group)
			if err != nil {
				return nil, err
			}
			names = append(names, members...)

		case strings.Contains(arg, "+"):
			matches, err := d.resolver.ExpandGlob(arg, sessionNames)
			if err != nil {
				return nil, err
			}
			names = append(names, matches...)

		case strings.Contains(arg, ":"):
			project, _ := ParseTarget(arg)
			names = append(names, project)

		default:
			names = append(names, arg)
		}
	}
	return names, nil
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
		if !d.client.HasSession(result.Name) {
			if err := d.builder.Build(result.Name, result.Config, result.Root); err != nil {
				return fmt.Errorf("building session: %w", err)
			}
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
	seen := make(map[string]bool)
	var items []string

	projects, _ := d.store.List()
	for _, p := range projects {
		if !seen[p.Name] {
			seen[p.Name] = true
			items = append(items, p.Name)
		}
	}

	sessions, _ := d.client.ListSessions()
	for _, s := range sessions {
		if !seen[s.Name] {
			seen[s.Name] = true
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
