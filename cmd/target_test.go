package cmd

import "testing"

func TestParseTarget_WithColon(t *testing.T) {
	proj, win := ParseTarget("blog:editor")
	if proj != "blog" {
		t.Errorf("project = %q, want blog", proj)
	}
	if win != "editor" {
		t.Errorf("window = %q, want editor", win)
	}
}

func TestParseTarget_WithoutColon(t *testing.T) {
	proj, win := ParseTarget("blog")
	if proj != "blog" {
		t.Errorf("project = %q, want blog", proj)
	}
	if win != "" {
		t.Errorf("window = %q, want empty", win)
	}
}
