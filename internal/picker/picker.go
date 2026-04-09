package picker

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Picker interface {
	Pick(items []string, prompt string) (string, error)
}

func New(backend string, stderr io.Writer) (Picker, error) {
	switch backend {
	case "fzf":
		return &FzfPicker{Build: exec.Command, Stderr: stderr}, nil
	case "gum":
		return &GumPicker{Build: exec.Command, Stderr: stderr}, nil
	default:
		return nil, fmt.Errorf("unknown picker backend: %s", backend)
	}
}

type CommandBuilder func(name string, args ...string) *exec.Cmd

func runExternal(build CommandBuilder, bin string, args []string, stdin io.Reader, stderr io.Writer) (string, error) {
	cmd := build(bin, args...)
	cmd.Stdin = stdin
	cmd.Stderr = stderr

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
