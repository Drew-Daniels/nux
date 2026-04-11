package config

import (
	"os"
	"strings"
	"testing"
)

func TestValidate_Valid(t *testing.T) {
	cfg := &ProjectConfig{
		Root: "~/projects/test",
		Windows: []Window{
			{Name: "editor", Layout: "tiled", Panes: []Pane{{Command: "vim"}}},
			{Name: "shell", Layout: "even-horizontal", Panes: []Pane{{Command: ""}}},
		},
	}
	if errs := Validate(cfg); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidate_CommandAndWindowsMutuallyExclusive(t *testing.T) {
	cfg := &ProjectConfig{
		Command: "vim",
		Windows: []Window{{Name: "editor", Panes: []Pane{{Command: "vim"}}}},
	}
	errs := Validate(cfg)
	if len(errs) == 0 {
		t.Fatal("expected error for command + windows")
	}
	assertContains(t, errs[0].Error(), "mutually exclusive")
}

func TestValidate_WindowRequiresPanes(t *testing.T) {
	cfg := &ProjectConfig{
		Windows: []Window{
			{Name: "editor"},
		},
	}
	errs := Validate(cfg)
	if len(errs) == 0 {
		t.Fatal("expected error for window without panes")
	}
	assertContains(t, errs[0].Error(), "at least one pane is required")
}

func TestValidate_WindowNameRequired(t *testing.T) {
	cfg := &ProjectConfig{
		Windows: []Window{{Layout: "tiled", Panes: []Pane{{Command: ""}}}},
	}
	errs := Validate(cfg)
	if len(errs) == 0 {
		t.Fatal("expected error for missing window name")
	}
	assertContains(t, errs[0].Error(), "name is required")
}

func TestValidate_InvalidLayout(t *testing.T) {
	cfg := &ProjectConfig{
		Windows: []Window{{Name: "editor", Layout: "bogus", Panes: []Pane{{Command: ""}}}},
	}
	errs := Validate(cfg)
	if len(errs) == 0 {
		t.Fatal("expected error for invalid layout")
	}
	assertContains(t, errs[0].Error(), "invalid layout")
}

func TestValidate_ValidLayouts(t *testing.T) {
	layouts := []string{"even-horizontal", "even-vertical", "main-horizontal", "main-vertical", "tiled", ""}
	for _, l := range layouts {
		cfg := &ProjectConfig{
			Windows: []Window{{Name: "w", Layout: l, Panes: []Pane{{Command: ""}}}},
		}
		if errs := Validate(cfg); len(errs) != 0 {
			t.Errorf("layout %q should be valid, got %v", l, errs)
		}
	}
}

func TestValidate_CustomLayout(t *testing.T) {
	cfg := &ProjectConfig{
		Windows: []Window{{Name: "w", Layout: "b]cd,159x43,0,0{79x43,0,0,0,79x43,80,0,1}", Panes: []Pane{{Command: ""}}}},
	}
	// Custom tmux layout strings start with a hex dimension and contain commas
	// Our heuristic checks for a comma at position 4
	if errs := Validate(cfg); len(errs) != 0 {
		t.Errorf("custom layout should be valid, got %v", errs)
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	cfg := &ProjectConfig{
		Command: "vim",
		Windows: []Window{
			{Layout: "bogus"},
		},
	}
	errs := Validate(cfg)
	if len(errs) < 4 {
		t.Fatalf("expected at least 4 errors (mutually exclusive, name, panes, layout), got %d: %v", len(errs), errs)
	}
}

func TestValidate_CommandOnlyIsValid(t *testing.T) {
	cfg := &ProjectConfig{
		Root:    "~/projects/test",
		Command: "just dev",
	}
	if errs := Validate(cfg); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidateAllWith(t *testing.T) {
	dir := t.TempDir()
	store := NewProjectStore(dir)

	_ = store.Save("valid", &ProjectConfig{
		Windows: []Window{{Name: "editor", Panes: []Pane{{Command: "vim"}}}},
	})
	_ = store.Save("invalid", &ProjectConfig{
		Command: "vim",
		Windows: []Window{{Name: "editor", Panes: []Pane{{Command: "vim"}}}},
	})

	results, err := ValidateAllWith(store)
	if err != nil {
		t.Fatalf("ValidateAllWith: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	var validResult, invalidResult *ValidationResult
	for i := range results {
		switch results[i].Name {
		case "valid":
			validResult = &results[i]
		case "invalid":
			invalidResult = &results[i]
		}
	}
	if validResult == nil || invalidResult == nil {
		t.Fatal("expected both valid and invalid results")
	}
	if len(validResult.Errors) != 0 {
		t.Errorf("valid config should have no errors, got %v", validResult.Errors)
	}
	if len(invalidResult.Errors) == 0 {
		t.Error("invalid config should have errors")
	}
}

func TestValidateAllWith_LoadError(t *testing.T) {
	dir := t.TempDir()
	store := NewProjectStore(dir)

	if err := os.WriteFile(dir+"/broken.yaml", []byte(":\n  :\n  - :\n    bad: ["), 0o644); err != nil {
		t.Fatal(err)
	}

	results, err := ValidateAllWith(store)
	if err != nil {
		t.Fatalf("ValidateAllWith: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if len(results[0].Errors) == 0 {
		t.Error("expected load error to be recorded")
	}
}

func assertContains(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Errorf("expected %q to contain %q", s, substr)
	}
}
