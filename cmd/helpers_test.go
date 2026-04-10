package cmd

import (
	"bytes"
	"io"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/Drew-Daniels/nux/internal/config"
	"github.com/Drew-Daniels/nux/internal/picker"
	"github.com/Drew-Daniels/nux/internal/resolver"
	"github.com/Drew-Daniels/nux/internal/tmux"
)

func durationMinutes(m int) time.Duration {
	return time.Duration(m) * time.Minute
}

type noopPicker struct{}

func (noopPicker) Pick([]string, string) (string, error) { return "", nil }

func testDeps(t *testing.T) *deps {
	t.Helper()

	projectCfgDir := t.TempDir()
	store := config.NewProjectStore(projectCfgDir)
	global := &config.GlobalConfig{
		ProjectsDir: t.TempDir(),
	}
	client := &tmux.MockClient{}
	builder := tmux.NewBuilder(client, global)
	res := resolver.NewResolverWithStore(global, store)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	return &deps{
		global:        global,
		client:        client,
		builder:       builder,
		resolver:      res,
		store:         store,
		projectCfgDir: projectCfgDir,
		noAttach:      false,
		force:         false,
		deleteForce:   false,
		run:           "",
		layout:        "",
		panes:         0,
		editor:        "echo",
		vars:          map[string]string{},
		stdin:         strings.NewReader(""),
		stdout:        stdout,
		stderr:        stderr,
		getwd:         func() (string, error) { return "/tmp/test", nil },
		confirm:       func(string) (bool, error) { return true, nil },
		openEditor:    func(string) error { return nil },
		newPicker:     func(string, io.Writer) (picker.Picker, error) { return noopPicker{}, nil },
		execCmd:       exec.Command,
		help:          func() error { return nil },
	}
}

func stdoutStr(d *deps) string {
	return d.stdout.(*bytes.Buffer).String()
}

func stderrStr(d *deps) string {
	return d.stderr.(*bytes.Buffer).String()
}
