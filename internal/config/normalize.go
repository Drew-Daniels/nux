package config

import (
	"regexp"
	"strings"
)

var (
	multiUnderscore = regexp.MustCompile(`_+`)
	sessionReplacer = strings.NewReplacer(".", "_", ":", "_", " ", "_")
)

// NormalizeSessionName replaces characters that tmux disallows in session names.
func NormalizeSessionName(name string) string {
	name = sessionReplacer.Replace(name)
	name = strings.TrimLeft(name, "-")
	name = multiUnderscore.ReplaceAllString(name, "_")
	name = strings.TrimRight(name, "_")
	return name
}
