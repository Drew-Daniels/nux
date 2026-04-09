package cmd

import (
	"strings"
	"testing"
)

func TestRunVersionWith(t *testing.T) {
	old := Version
	oldCommit := Commit
	oldDate := Date
	Version = "1.2.3"
	Commit = "abc123"
	Date = "2025-01-01"
	defer func() { Version = old; Commit = oldCommit; Date = oldDate }()

	d := testDeps(t)
	if err := runVersionWith(d); err != nil {
		t.Fatalf("runVersionWith: %v", err)
	}
	out := stdoutStr(d)
	if !strings.Contains(out, "1.2.3") {
		t.Errorf("expected version in output, got %q", out)
	}
	if !strings.Contains(out, "abc123") {
		t.Errorf("expected commit in output, got %q", out)
	}
	if !strings.Contains(out, "2025-01-01") {
		t.Errorf("expected date in output, got %q", out)
	}
}

func TestFormatDuration_Minutes(t *testing.T) {
	tests := []struct {
		minutes int
		want    string
	}{
		{0, "0m"},
		{5, "5m"},
		{59, "59m"},
	}
	for _, tt := range tests {
		got := formatDuration(durationMinutes(tt.minutes))
		if got != tt.want {
			t.Errorf("formatDuration(%dm) = %q, want %q", tt.minutes, got, tt.want)
		}
	}
}

func TestFormatDuration_Hours(t *testing.T) {
	tests := []struct {
		minutes int
		want    string
	}{
		{60, "1h 0m"},
		{90, "1h 30m"},
		{150, "2h 30m"},
	}
	for _, tt := range tests {
		got := formatDuration(durationMinutes(tt.minutes))
		if got != tt.want {
			t.Errorf("formatDuration(%dm) = %q, want %q", tt.minutes, got, tt.want)
		}
	}
}
