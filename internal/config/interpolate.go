package config

import (
	"os"
	"os/exec"
	"strings"
)

// CommandRunner executes a shell command and returns its output.
type CommandRunner func(command string) (string, error)

// EnvLookup expands environment variable references in a string.
type EnvLookup func(s string) string

// DefaultCommandRunner runs a command via sh -c and returns trimmed output.
func DefaultCommandRunner(command string) (string, error) {
	out, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(out), "\n"), nil
}

// Interpolator resolves {{var}} placeholders and environment variables in configs.
type Interpolator struct {
	RunCommand CommandRunner
	ExpandEnv  EnvLookup
}

// NewInterpolator returns an Interpolator with default shell and env backends.
func NewInterpolator() *Interpolator {
	return &Interpolator{
		RunCommand: DefaultCommandRunner,
		ExpandEnv:  os.ExpandEnv,
	}
}

// Interpolate resolves all variable and environment references in a ProjectConfig.
func (ip *Interpolator) Interpolate(cfg *ProjectConfig) error {
	vars, err := ip.resolveVars(cfg.Vars)
	if err != nil {
		return err
	}

	return ip.applyTransform(cfg, func(s string) string {
		for k, v := range vars {
			s = strings.ReplaceAll(s, "{{"+k+"}}", v)
		}
		return ip.ExpandEnv(s)
	})
}

// InterpolateVars resolves only {{var}} placeholders without re-expanding
// environment variables. Use this when the config has already been through a
// full Interpolate pass and only custom variable overrides need reapplying.
func (ip *Interpolator) InterpolateVars(cfg *ProjectConfig) error {
	vars, err := ip.resolveVars(cfg.Vars)
	if err != nil {
		return err
	}

	return ip.applyTransform(cfg, func(s string) string {
		for k, v := range vars {
			s = strings.ReplaceAll(s, "{{"+k+"}}", v)
		}
		return s
	})
}

func (ip *Interpolator) applyTransform(cfg *ProjectConfig, fn func(string) string) error {
	cfg.Root = fn(cfg.Root)
	cfg.Command = fn(cfg.Command)

	applySlice := func(ss []string) {
		for i, s := range ss {
			ss[i] = fn(s)
		}
	}
	applySlice(cfg.OnStart)
	applySlice(cfg.OnReady)
	applySlice(cfg.OnDetach)
	applySlice(cfg.OnStop)

	for k, v := range cfg.Env {
		cfg.Env[k] = fn(v)
	}

	for i := range cfg.Windows {
		w := &cfg.Windows[i]
		w.Root = fn(w.Root)
		w.Command = fn(w.Command)
		for k, v := range w.Env {
			w.Env[k] = fn(v)
		}
		for j := range w.Panes {
			p := &w.Panes[j]
			p.Root = fn(p.Root)
			p.Command = fn(p.Command)
		}
	}

	return nil
}

func (ip *Interpolator) resolveVars(vars map[string]string) (map[string]string, error) {
	resolved := make(map[string]string, len(vars))
	for k, v := range vars {
		val, err := ip.resolveValue(v)
		if err != nil {
			return nil, err
		}
		resolved[k] = val
	}
	return resolved, nil
}

func (ip *Interpolator) resolveValue(v string) (string, error) {
	if strings.HasPrefix(v, "`") && strings.HasSuffix(v, "`") && len(v) > 2 {
		cmd := v[1 : len(v)-1]
		return ip.RunCommand(cmd)
	}
	return v, nil
}
