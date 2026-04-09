package ui

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Prompter struct {
	In  io.Reader
	Out io.Writer
}

func (p *Prompter) Confirm(prompt string) (bool, error) {
	_, _ = fmt.Fprintf(p.Out, "%s [y/N]: ", prompt)
	scanner := bufio.NewScanner(p.In)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, err
		}
		return false, nil
	}
	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return answer == "y" || answer == "yes", nil
}
