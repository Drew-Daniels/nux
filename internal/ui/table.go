package ui

import (
	"fmt"
	"strings"
)

type Table struct {
	Headers []string
	Rows    [][]string
}

func (t *Table) Render() string {
	if len(t.Headers) == 0 {
		return ""
	}

	widths := make([]int, len(t.Headers))
	for i, h := range t.Headers {
		widths[i] = len(h)
	}
	for _, row := range t.Rows {
		for i := range widths {
			if i < len(row) && len(row[i]) > widths[i] {
				widths[i] = len(row[i])
			}
		}
	}

	var b strings.Builder
	for i, h := range t.Headers {
		if i > 0 {
			b.WriteString("  ")
		}
		fmt.Fprintf(&b, "%-*s", widths[i], h)
	}

	for _, row := range t.Rows {
		b.WriteByte('\n')
		for i := range widths {
			if i > 0 {
				b.WriteString("  ")
			}
			val := ""
			if i < len(row) {
				val = row[i]
			}
			fmt.Fprintf(&b, "%-*s", widths[i], val)
		}
	}

	return b.String()
}
