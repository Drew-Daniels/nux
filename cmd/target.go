package cmd

import "strings"

// ParseTarget splits a "project:window" argument into its components.
// If there is no ":" separator, window is empty.
func ParseTarget(arg string) (project, window string) {
	if i := strings.IndexByte(arg, ':'); i >= 0 {
		return arg[:i], arg[i+1:]
	}
	return arg, ""
}
