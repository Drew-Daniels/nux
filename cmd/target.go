package cmd

import (
	"fmt"
	"strings"
)

// sessionArg is one CLI token: a project name, or a project plus one or more
// window names (subset start / restart).
type sessionArg struct {
	Project string
	// Windows is nil for a full session; non-nil means only these windows (in order).
	Windows []string
}

// parseSessionToken parses a single argument like "blog", "blog:editor", or
// "blog:editor,server" (comma-separated window names).
func parseSessionToken(s string) (sessionArg, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return sessionArg{}, fmt.Errorf("empty argument")
	}
	if !strings.Contains(s, ":") {
		return sessionArg{Project: s}, nil
	}
	i := strings.IndexByte(s, ':')
	project := strings.TrimSpace(s[:i])
	rest := strings.TrimSpace(s[i+1:])
	if project == "" {
		return sessionArg{}, fmt.Errorf("invalid target %q: missing project name before ':'", s)
	}
	if rest == "" {
		return sessionArg{}, fmt.Errorf("invalid target %q: missing window name after ':'", s)
	}
	parts := strings.Split(rest, ",")
	windows := make([]string, 0, len(parts))
	for _, p := range parts {
		w := strings.TrimSpace(p)
		if w == "" {
			return sessionArg{}, fmt.Errorf("invalid target %q: empty window name in list", s)
		}
		windows = append(windows, w)
	}
	return sessionArg{Project: project, Windows: windows}, nil
}

// ParseTarget splits a "project:window" argument into project and the first
// window name only. For multi-window targets (comma-separated), it returns
// the first window. If there is no ":" separator, window is empty.
func ParseTarget(arg string) (project, window string) {
	sa, err := parseSessionToken(arg)
	if err != nil {
		return arg, ""
	}
	if sa.Windows == nil {
		return sa.Project, ""
	}
	return sa.Project, sa.Windows[0]
}
