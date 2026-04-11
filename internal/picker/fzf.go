package picker

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type FzfPicker struct {
	Build  CommandBuilder
	Stderr io.Writer
}

func (f *FzfPicker) Pick(items []string, prompt string) (string, error) {
	args := []string{"--prompt", prompt + " ", "--reverse"}
	input := strings.NewReader(strings.Join(items, "\n"))

	result, err := runExternal(f.Build, "fzf", args, input, f.Stderr)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && (exitErr.ExitCode() == 1 || exitErr.ExitCode() == 130) {
			return "", nil
		}
		return "", fmt.Errorf("fzf failed: %w", err)
	}
	return result, nil
}
