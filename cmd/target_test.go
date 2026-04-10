package cmd

import (
	"reflect"
	"testing"
)

func TestParseTarget_WithColon(t *testing.T) {
	proj, win := ParseTarget("blog:editor")
	if proj != "blog" {
		t.Errorf("project = %q, want blog", proj)
	}
	if win != "editor" {
		t.Errorf("window = %q, want editor", win)
	}
}

func TestParseTarget_MultiReturnsFirstWindow(t *testing.T) {
	proj, win := ParseTarget("blog:server,editor")
	if proj != "blog" || win != "server" {
		t.Errorf("got %q, %q want blog, server", proj, win)
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

func TestParseSessionToken_MultiWindow(t *testing.T) {
	sa, err := parseSessionToken("blog:a,b")
	if err != nil {
		t.Fatal(err)
	}
	if sa.Project != "blog" || !reflect.DeepEqual(sa.Windows, []string{"a", "b"}) {
		t.Errorf("got %+v", sa)
	}
}

func TestParseSessionToken_Invalid(t *testing.T) {
	for _, s := range []string{":only", "blog:", "blog:a,,b", ""} {
		if _, err := parseSessionToken(s); err == nil {
			t.Errorf("expected error for %q", s)
		}
	}
}
