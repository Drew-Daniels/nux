package picker

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type GumPicker struct {
	Build  CommandBuilder
	Stderr io.Writer
}

func (g *GumPicker) Pick(items []string, prompt string) (string, error) {
	args := []string{"filter", "--placeholder", prompt}
	input := strings.NewReader(strings.Join(items, "\n"))

	result, err := runExternal(g.Build, "gum", args, input, g.Stderr)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			return "", nil
		}
		return "", fmt.Errorf("gum failed: %w", err)
	}
	return result, nil
}
