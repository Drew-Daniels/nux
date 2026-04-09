package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/Drew-Daniels/nux/internal/config"
)

func TestRunDoctorWith_AllPass(t *testing.T) {
	d := testDeps(t)

	checkBin := func(name string) (string, bool) { return "/usr/bin/" + name, true }
	probeVersion := func() (string, error) { return "tmux 3.6", nil }
	checkStat := func(path string) (os.FileInfo, error) { return os.Stat(d.global.ProjectsDir) }

	err := runDoctorWith(d, checkBin, probeVersion, checkStat)
	if err != nil {
		t.Fatalf("runDoctorWith: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "All checks passed") {
		t.Errorf("expected 'All checks passed', got %q", out)
	}
	if !strings.Contains(out, "tmux 3.6") {
		t.Errorf("expected version in output, got %q", out)
	}
}

func TestRunDoctorWith_TmuxMissing(t *testing.T) {
	d := testDeps(t)

	checkBin := func(name string) (string, bool) { return "", false }
	probeVersion := func() (string, error) { return "", nil }
	checkStat := func(path string) (os.FileInfo, error) { return os.Stat(d.global.ProjectsDir) }

	err := runDoctorWith(d, checkBin, probeVersion, checkStat)
	if err == nil {
		t.Fatal("expected error when tmux missing")
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "[missing]") {
		t.Errorf("expected '[missing]' in output, got %q", out)
	}
}

func TestRunDoctorChecks_ZoxideMissing(t *testing.T) {
	d := testDeps(t)
	d.global.Zoxide = true

	checkBin := func(name string) (string, bool) {
		if name == "zoxide" {
			return "", false
		}
		return "/usr/bin/" + name, true
	}
	checkStat := func(path string) (os.FileInfo, error) { return os.Stat(d.global.ProjectsDir) }

	ok := runDoctorChecks(d, checkBin, checkStat)
	if ok {
		t.Error("expected false when zoxide missing")
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "[missing]") {
		t.Errorf("expected '[missing]' in output, got %q", out)
	}
}

func TestRunDoctorChecks_PickerMissing(t *testing.T) {
	d := testDeps(t)
	d.global.Picker = "fzf"

	checkBin := func(name string) (string, bool) {
		if name == "fzf" {
			return "", false
		}
		return "/usr/bin/" + name, true
	}
	checkStat := func(path string) (os.FileInfo, error) { return os.Stat(d.global.ProjectsDir) }

	_ = runDoctorChecks(d, checkBin, checkStat)

	out := stdoutStr(d)
	if !strings.Contains(out, "[warn]") {
		t.Errorf("expected '[warn]' in output, got %q", out)
	}
}

func TestRunDoctorChecks_InvalidConfig(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("bad", &config.ProjectConfig{
		Command: "vim",
		Windows: []config.Window{{Name: "editor"}},
	})

	checkBin := func(name string) (string, bool) { return "/usr/bin/" + name, true }
	checkStat := func(path string) (os.FileInfo, error) { return os.Stat(d.global.ProjectsDir) }

	_ = runDoctorChecks(d, checkBin, checkStat)

	errOut := stderrStr(d)
	if !strings.Contains(errOut, "[fail]") {
		t.Errorf("expected '[fail]' in stderr, got %q", errOut)
	}
}

func TestRunDoctorChecks_DirsMissing(t *testing.T) {
	d := testDeps(t)

	checkBin := func(name string) (string, bool) { return "/usr/bin/" + name, true }
	checkStat := func(path string) (os.FileInfo, error) {
		return nil, os.ErrNotExist
	}

	_ = runDoctorChecks(d, checkBin, checkStat)

	out := stdoutStr(d)
	if !strings.Contains(out, "[warn]") {
		t.Errorf("expected '[warn]' for missing dirs, got %q", out)
	}
}

func TestRunDoctorWith_WithValidConfigs(t *testing.T) {
	d := testDeps(t)
	_ = d.store.Save("good", &config.ProjectConfig{
		Windows: []config.Window{{Name: "editor"}},
	})

	checkBin := func(name string) (string, bool) { return "/usr/bin/" + name, true }
	probeVersion := func() (string, error) { return "tmux 3.6", nil }
	checkStat := func(path string) (os.FileInfo, error) { return os.Stat(d.global.ProjectsDir) }

	err := runDoctorWith(d, checkBin, probeVersion, checkStat)
	if err != nil {
		t.Fatalf("runDoctorWith: %v", err)
	}

	out := stdoutStr(d)
	if !strings.Contains(out, "all configs valid") {
		t.Errorf("expected 'all configs valid', got %q", out)
	}
}
