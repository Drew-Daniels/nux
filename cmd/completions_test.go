package cmd

import (
	"testing"
)

func TestRunCompletionsWith_Bash(t *testing.T) {
	d := testDeps(t)
	if err := runCompletionsWith(d, []string{"bash"}); err != nil {
		t.Fatalf("runCompletionsWith bash: %v", err)
	}
	out := stdoutStr(d)
	if len(out) == 0 {
		t.Error("expected bash completion output")
	}
}

func TestRunCompletionsWith_Zsh(t *testing.T) {
	d := testDeps(t)
	if err := runCompletionsWith(d, []string{"zsh"}); err != nil {
		t.Fatalf("runCompletionsWith zsh: %v", err)
	}
	out := stdoutStr(d)
	if len(out) == 0 {
		t.Error("expected zsh completion output")
	}
}

func TestRunCompletionsWith_Fish(t *testing.T) {
	d := testDeps(t)
	if err := runCompletionsWith(d, []string{"fish"}); err != nil {
		t.Fatalf("runCompletionsWith fish: %v", err)
	}
	out := stdoutStr(d)
	if len(out) == 0 {
		t.Error("expected fish completion output")
	}
}

func TestRunCompletionsWith_Unknown(t *testing.T) {
	d := testDeps(t)
	if err := runCompletionsWith(d, []string{"powershell"}); err != nil {
		t.Fatalf("runCompletionsWith unknown: %v", err)
	}
	out := stdoutStr(d)
	if len(out) != 0 {
		t.Error("expected no output for unknown shell")
	}
}
