package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestValidateLayoutFlags_None(t *testing.T) {
	d := testDeps(t)
	if err := validateLayoutFlags(d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateLayoutFlags_ValidLayout(t *testing.T) {
	d := testDeps(t)
	d.layout = "tiled"
	d.panes = 4
	if err := validateLayoutFlags(d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateLayoutFlags_InvalidLayout(t *testing.T) {
	d := testDeps(t)
	d.layout = "bogus"
	d.panes = 2
	if err := validateLayoutFlags(d); err == nil {
		t.Fatal("expected error for invalid layout")
	}
}

func TestValidateLayoutFlags_NegativePanes(t *testing.T) {
	d := testDeps(t)
	d.layout = "tiled"
	d.panes = -1
	if err := validateLayoutFlags(d); err == nil {
		t.Fatal("expected error for negative panes")
	}
}

func TestValidateLayoutFlags_DefaultsLayout(t *testing.T) {
	d := testDeps(t)
	d.panes = 4
	if err := validateLayoutFlags(d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.layout != "tiled" {
		t.Errorf("layout = %q, want tiled (default)", d.layout)
	}
}

func TestValidateLayoutFlags_DefaultsPanes(t *testing.T) {
	d := testDeps(t)
	d.layout = "tiled"
	if err := validateLayoutFlags(d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.panes != 2 {
		t.Errorf("panes = %d, want 2 (default)", d.panes)
	}
}

func TestAdHocLayoutFromDeps_Nil(t *testing.T) {
	d := testDeps(t)
	if got := adHocLayoutFromDeps(d); got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
}

func TestAdHocLayoutFromDeps_Set(t *testing.T) {
	d := testDeps(t)
	d.layout = "tiled"
	d.panes = 4
	got := adHocLayoutFromDeps(d)
	if got == nil {
		t.Fatal("expected non-nil")
	}
	if got.Layout != "tiled" || got.Panes != 4 {
		t.Errorf("got %+v, want {tiled, 4}", got)
	}
}

func TestMatchSubcommand_ExactMatch(t *testing.T) {
	cmd, ok := matchSubcommand("list")
	if !ok {
		t.Fatal("expected match for 'list'")
	}
	if cmd.Name() != "list" {
		t.Errorf("matched %q, want list", cmd.Name())
	}
}

func TestMatchSubcommand_AliasMatch(t *testing.T) {
	cmd, ok := matchSubcommand("ls")
	if !ok {
		t.Fatal("expected match for 'ls' alias")
	}
	if cmd.Name() != "list" {
		t.Errorf("matched %q, want list", cmd.Name())
	}
}

func TestMatchSubcommand_UnambiguousPrefix(t *testing.T) {
	cmd, ok := matchSubcommand("li")
	if !ok {
		t.Fatal("expected match for 'li'")
	}
	if cmd.Name() != "list" {
		t.Errorf("matched %q, want list", cmd.Name())
	}
}

func TestMatchSubcommand_AmbiguousPrefix(t *testing.T) {
	_, ok := matchSubcommand("v")
	if ok {
		t.Error("expected no match for ambiguous prefix 'v'")
	}
}

func TestMatchSubcommand_NoMatch(t *testing.T) {
	_, ok := matchSubcommand("zzz")
	if ok {
		t.Error("expected no match for 'zzz'")
	}
}

func TestOptionsEditor(t *testing.T) {
	o := options{editorFunc: func() string { return "/usr/bin/vim" }}
	if got := o.editor(); got != "/usr/bin/vim" {
		t.Errorf("editor() = %q, want /usr/bin/vim", got)
	}
}

func TestOptionsEditor_Fallback(t *testing.T) {
	t.Setenv("EDITOR", "nano")
	o := options{}
	if got := o.editor(); got != "nano" {
		t.Errorf("editor() = %q, want nano", got)
	}
}

func TestParseVars_Valid(t *testing.T) {
	vars := parseVars([]string{"key=value", "foo=bar"}, &bytes.Buffer{})
	if vars["key"] != "value" {
		t.Errorf("key = %q, want value", vars["key"])
	}
	if vars["foo"] != "bar" {
		t.Errorf("foo = %q, want bar", vars["foo"])
	}
}

func TestParseVars_Malformed(t *testing.T) {
	var stderr bytes.Buffer
	vars := parseVars([]string{"noequals", "good=val"}, &stderr)
	if _, ok := vars["noequals"]; ok {
		t.Error("malformed var should be skipped")
	}
	if vars["good"] != "val" {
		t.Errorf("good = %q, want val", vars["good"])
	}
	if !strings.Contains(stderr.String(), "malformed") {
		t.Errorf("expected warning in stderr, got %q", stderr.String())
	}
}

func TestParseVars_Empty(t *testing.T) {
	vars := parseVars(nil, &bytes.Buffer{})
	if len(vars) != 0 {
		t.Errorf("expected empty map, got %v", vars)
	}
}
