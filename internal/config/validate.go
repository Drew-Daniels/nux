package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var validLayouts = map[string]bool{
	"even-horizontal": true,
	"even-vertical":   true,
	"main-horizontal": true,
	"main-vertical":   true,
	"tiled":           true,
}

func IsValidLayout(layout string) bool {
	if layout == "" {
		return true
	}
	if validLayouts[layout] {
		return true
	}
	if len(layout) > 4 && layout[4] == ',' {
		return true
	}
	return false
}

func layoutError(context, layout string) error {
	return fmt.Errorf("%s: invalid layout %q (must be even-horizontal, even-vertical, main-horizontal, main-vertical, tiled, or a custom layout string)", context, layout)
}

func Validate(cfg *ProjectConfig) []error {
	var errs []error

	if len(cfg.Windows) == 0 {
		errs = append(errs, fmt.Errorf("at least one window is required"))
	}

	errs = append(errs, validateWindows(cfg.Windows)...)

	return errs
}

func validateWindows(windows []Window) []error {
	var errs []error
	for i, w := range windows {
		label := fmt.Sprintf("windows[%d]", i)

		if w.Name == "" {
			errs = append(errs, fmt.Errorf("%s: name is required", label))
		}

		if len(w.Panes) == 0 {
			errs = append(errs, fmt.Errorf("%s: at least one pane is required; use panes: [\"\"] for a bare shell", label))
		}

		if !IsValidLayout(w.Layout) {
			errs = append(errs, layoutError(label, w.Layout))
		}
	}
	return errs
}

var validPickers = map[string]bool{
	"fzf": true,
	"gum": true,
	"":    true,
}

func ValidateGlobal(cfg *GlobalConfig) (errs []error, warnings []error) {
	if !validPickers[cfg.Picker] {
		errs = append(errs, fmt.Errorf("picker: invalid value %q (must be fzf or gum)", cfg.Picker))
	}

	for i, d := range cfg.ProjectDirs {
		if strings.TrimSpace(d) == "" {
			errs = append(errs, fmt.Errorf("project_dirs[%d]: empty path", i))
		}
	}

	if cfg.DefaultSession != nil {
		for _, e := range validateWindows(cfg.DefaultSession.Windows) {
			errs = append(errs, fmt.Errorf("default_session.%w", e))
		}
	}

	for name, members := range cfg.Groups {
		if len(members) == 0 {
			warnings = append(warnings, fmt.Errorf("group %q is empty", name))
		}
	}

	for i, d := range cfg.ProjectDirs {
		expanded := expandHome(d)
		if _, err := os.Stat(expanded); os.IsNotExist(err) {
			warnings = append(warnings, fmt.Errorf("project_dirs[%d]: directory does not exist: %s", i, d))
		}
	}

	return errs, warnings
}

func expandHome(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return filepath.Join(home, path[1:])
}

type ValidationResult struct {
	Name   string
	Errors []error
}

func ValidateAllWith(store ProjectStore) ([]ValidationResult, error) {
	projects, err := store.List()
	if err != nil {
		return nil, err
	}

	var results []ValidationResult
	for _, p := range projects {
		cfg, _, err := store.Load(p.Name)
		if err != nil {
			results = append(results, ValidationResult{Name: p.Name, Errors: []error{err}})
			continue
		}
		if errs := Validate(cfg); len(errs) > 0 {
			results = append(results, ValidationResult{Name: p.Name, Errors: errs})
		} else {
			results = append(results, ValidationResult{Name: p.Name})
		}
	}
	return results, nil
}
