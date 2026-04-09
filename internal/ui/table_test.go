package ui

import (
	"strings"
	"testing"
)

func TestRender_BasicTable(t *testing.T) {
	tbl := &Table{
		Headers: []string{"NAME", "STATUS", "ROOT"},
		Rows: [][]string{
			{"my-app", "running", "~/projects/my-app"},
			{"blog", "stopped", "~/projects/blog"},
		},
	}

	got := tbl.Render()

	lines := strings.Split(got, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d:\n%s", len(lines), got)
	}

	if !strings.Contains(lines[0], "NAME") || !strings.Contains(lines[0], "STATUS") {
		t.Errorf("headers missing: %s", lines[0])
	}

	if !strings.Contains(lines[1], "my-app") {
		t.Errorf("first row missing data: %s", lines[1])
	}
}

func TestRender_EmptyRows(t *testing.T) {
	tbl := &Table{
		Headers: []string{"NAME", "VALUE"},
		Rows:    [][]string{},
	}

	got := tbl.Render()
	lines := strings.Split(got, "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line (header only), got %d:\n%s", len(lines), got)
	}
}

func TestRender_SingleColumn(t *testing.T) {
	tbl := &Table{
		Headers: []string{"ITEM"},
		Rows: [][]string{
			{"alpha"},
			{"beta"},
		},
	}

	got := tbl.Render()
	lines := strings.Split(got, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d:\n%s", len(lines), got)
	}

	if strings.TrimSpace(lines[0]) != "ITEM" {
		t.Errorf("expected header ITEM, got %q", strings.TrimSpace(lines[0]))
	}
}

func TestRender_NoHeaders(t *testing.T) {
	tbl := &Table{
		Headers: []string{},
		Rows:    [][]string{{"a", "b"}},
	}

	got := tbl.Render()
	if got != "" {
		t.Errorf("expected empty string for no headers, got %q", got)
	}
}

func TestRender_ColumnWidthFromData(t *testing.T) {
	tbl := &Table{
		Headers: []string{"X", "Y"},
		Rows: [][]string{
			{"short", "a-much-longer-value"},
		},
	}

	got := tbl.Render()
	lines := strings.Split(got, "\n")

	headerX := strings.Index(lines[0], "Y")
	dataX := strings.Index(lines[1], "a-much-longer-value")

	if headerX != dataX {
		t.Errorf("columns not aligned: header Y at %d, data at %d", headerX, dataX)
	}
}

func TestRender_ShortRow(t *testing.T) {
	tbl := &Table{
		Headers: []string{"A", "B", "C"},
		Rows: [][]string{
			{"only-one"},
		},
	}

	got := tbl.Render()
	lines := strings.Split(got, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d:\n%s", len(lines), got)
	}
}
