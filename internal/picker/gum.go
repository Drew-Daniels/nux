package picker

import (
	"fmt"
	"io"
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
		return "", fmt.Errorf("gum failed: %w", err)
	}
	if result == "" {
		return "", fmt.Errorf("selection cancelled")
	}
	return result, nil
}
