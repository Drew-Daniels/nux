package picker

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNew_Fzf(t *testing.T) {
	p, err := New("fzf", os.Stderr)
	if err != nil {
		t.Fatalf("New(fzf): %v", err)
	}
	if _, ok := p.(*FzfPicker); !ok {
		t.Errorf("expected *FzfPicker, got %T", p)
	}
}

func TestNew_Gum(t *testing.T) {
	p, err := New("gum", os.Stderr)
	if err != nil {
		t.Fatalf("New(gum): %v", err)
	}
	if _, ok := p.(*GumPicker); !ok {
		t.Errorf("expected *GumPicker, got %T", p)
	}
}

func TestNew_Unknown(t *testing.T) {
	_, err := New("unknown", os.Stderr)
	if err == nil {
		t.Fatal("expected error for unknown backend")
	}
}

func helperScript(t *testing.T, script string) string {
	t.Helper()
	dir := t.TempDir()
	ext := ".sh"
	if runtime.GOOS == "windows" {
		ext = ".bat"
	}
	path := filepath.Join(dir, "helper"+ext)
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRunExternal_Success(t *testing.T) {
	script := helperScript(t, "#!/bin/sh\necho selected-item\n")
	build := func(name string, args ...string) *exec.Cmd {
		return exec.Command(script)
	}

	result, err := runExternal(build, "fake", nil, nil, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("runExternal: %v", err)
	}
	if result != "selected-item" {
		t.Errorf("result = %q, want selected-item", result)
	}
}

func TestRunExternal_Error(t *testing.T) {
	build := func(name string, args ...string) *exec.Cmd {
		return exec.Command("false")
	}

	_, err := runExternal(build, "fake", nil, nil, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error from failing command")
	}
}

func TestFzfPicker_Success(t *testing.T) {
	script := helperScript(t, "#!/bin/sh\necho blog\n")
	p := &FzfPicker{
		Build:  func(name string, args ...string) *exec.Cmd { return exec.Command(script) },
		Stderr: &bytes.Buffer{},
	}
	result, err := p.Pick([]string{"blog", "api"}, "project")
	if err != nil {
		t.Fatalf("Pick: %v", err)
	}
	if result != "blog" {
		t.Errorf("result = %q, want blog", result)
	}
}

func TestFzfPicker_Error(t *testing.T) {
	p := &FzfPicker{
		Build:  func(name string, args ...string) *exec.Cmd { return exec.Command("false") },
		Stderr: &bytes.Buffer{},
	}
	_, err := p.Pick([]string{"blog"}, "project")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGumPicker_Success(t *testing.T) {
	script := helperScript(t, "#!/bin/sh\necho blog\n")
	p := &GumPicker{
		Build:  func(name string, args ...string) *exec.Cmd { return exec.Command(script) },
		Stderr: &bytes.Buffer{},
	}
	result, err := p.Pick([]string{"blog", "api"}, "project")
	if err != nil {
		t.Fatalf("Pick: %v", err)
	}
	if result != "blog" {
		t.Errorf("result = %q, want blog", result)
	}
}

func TestGumPicker_EmptyResult(t *testing.T) {
	script := helperScript(t, "#!/bin/sh\necho\n")
	p := &GumPicker{
		Build:  func(name string, args ...string) *exec.Cmd { return exec.Command(script) },
		Stderr: &bytes.Buffer{},
	}
	_, err := p.Pick([]string{"blog"}, "project")
	if err == nil {
		t.Fatal("expected error for empty selection")
	}
}

func TestGumPicker_Error(t *testing.T) {
	p := &GumPicker{
		Build:  func(name string, args ...string) *exec.Cmd { return exec.Command("false") },
		Stderr: &bytes.Buffer{},
	}
	_, err := p.Pick([]string{"blog"}, "project")
	if err == nil {
		t.Fatal("expected error")
	}
}
