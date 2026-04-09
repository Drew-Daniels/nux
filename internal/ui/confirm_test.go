package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestConfirm_Yes(t *testing.T) {
	p := &Prompter{In: strings.NewReader("y\n"), Out: &bytes.Buffer{}}
	ok, err := p.Confirm("Continue?")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("expected true for 'y'")
	}
}

func TestConfirm_YesFullWord(t *testing.T) {
	p := &Prompter{In: strings.NewReader("yes\n"), Out: &bytes.Buffer{}}
	ok, err := p.Confirm("Continue?")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("expected true for 'yes'")
	}
}

func TestConfirm_No(t *testing.T) {
	p := &Prompter{In: strings.NewReader("n\n"), Out: &bytes.Buffer{}}
	ok, err := p.Confirm("Continue?")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("expected false for 'n'")
	}
}

func TestConfirm_Empty(t *testing.T) {
	p := &Prompter{In: strings.NewReader("\n"), Out: &bytes.Buffer{}}
	ok, err := p.Confirm("Continue?")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("expected false for empty input")
	}
}

func TestConfirm_EOF(t *testing.T) {
	p := &Prompter{In: strings.NewReader(""), Out: &bytes.Buffer{}}
	ok, err := p.Confirm("Continue?")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("expected false on EOF")
	}
}

func TestConfirm_Prompt(t *testing.T) {
	var buf bytes.Buffer
	p := &Prompter{In: strings.NewReader("n\n"), Out: &buf}
	_, _ = p.Confirm("Delete it?")
	if !strings.Contains(buf.String(), "Delete it?") {
		t.Errorf("expected prompt in output, got %q", buf.String())
	}
}
