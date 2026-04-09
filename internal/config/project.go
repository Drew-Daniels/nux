package config

import "path/filepath"

// ProjectConfigDir returns the directory containing project config files.
func ProjectConfigDir() string {
	return filepath.Join(DefaultConfigDir(), "projects")
}
