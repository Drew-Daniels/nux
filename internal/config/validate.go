package config

import "fmt"

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

	if cfg.Command != "" && len(cfg.Windows) > 0 {
		errs = append(errs, fmt.Errorf("command and windows are mutually exclusive"))
	}

	for i, w := range cfg.Windows {
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
