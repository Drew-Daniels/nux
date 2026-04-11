package resolver

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Drew-Daniels/nux/internal/config"
)

type Result struct {
	Name         string
	Config       *config.ProjectConfig
	Root         string
	ConfigSource string
	ConfigPath   string
}

type ZoxideQuerier func(name string) (string, error)

func DefaultZoxideQuerier(name string) (string, error) {
	out, err := exec.Command("zoxide", "query", name).Output()
	if err != nil {
		return "", err
	}
	path := strings.TrimSpace(string(out))
	if path == "" {
		return "", fmt.Errorf("zoxide returned empty path")
	}
	return path, nil
}

type DirChecker func(path string) (os.FileInfo, error)

type HomeDirFunc func() (string, error)

type Resolver struct {
	global       *config.GlobalConfig
	store        config.ProjectStore
	zoxideQuery  ZoxideQuerier
	Interpolator *config.Interpolator
	checkDir     DirChecker
	homeDir      HomeDirFunc
}

func NewResolverWithStore(global *config.GlobalConfig, store config.ProjectStore) *Resolver {
	return &Resolver{
		global:       global,
		store:        store,
		zoxideQuery:  DefaultZoxideQuerier,
		Interpolator: config.NewInterpolator(),
		checkDir:     os.Stat,
		homeDir:      os.UserHomeDir,
	}
}

func (r *Resolver) WithZoxideQuerier(q ZoxideQuerier) *Resolver {
	r.zoxideQuery = q
	return r
}

func (r *Resolver) WithInterpolator(ip *config.Interpolator) *Resolver {
	r.Interpolator = ip
	return r
}

func (r *Resolver) WithDirChecker(dc DirChecker) *Resolver {
	r.checkDir = dc
	return r
}

func (r *Resolver) WithHomeDir(fn HomeDirFunc) *Resolver {
	r.homeDir = fn
	return r
}

func (r *Resolver) Resolve(name string) (*Result, error) {
	cfg, cfgPath, err := r.store.Load(name)
	if err == nil {
		return r.resolveFromConfig(name, cfg, cfgPath)
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("loading project config: %w", err)
	}

	if r.global.Zoxide {
		if result, err := r.resolveFromZoxide(name); err == nil {
			return result, nil
		}
	}

	if result, err := r.resolveFromDirectory(name); err == nil {
		return result, nil
	}

	return nil, fmt.Errorf("project not found: %s", name)
}

func (r *Resolver) resolveFromConfig(name string, cfg *config.ProjectConfig, cfgPath string) (*Result, error) {
	if errs := config.Validate(cfg); len(errs) > 0 {
		msgs := make([]string, len(errs))
		for i, e := range errs {
			msgs[i] = e.Error()
		}
		return nil, fmt.Errorf("validation failed: %s", strings.Join(msgs, "; "))
	}

	if err := r.Interpolator.Interpolate(cfg); err != nil {
		return nil, fmt.Errorf("interpolation failed: %w", err)
	}

	root := resolveRootWith(cfg.Root, r.global.FirstProjectDir(), r.homeDir)

	return &Result{
		Name:         config.NormalizeSessionName(name),
		Config:       cfg,
		Root:         root,
		ConfigSource: "project",
		ConfigPath:   cfgPath,
	}, nil
}

func (r *Resolver) resolveFromZoxide(name string) (*Result, error) {
	path, err := r.zoxideQuery(name)
	if err != nil {
		return nil, err
	}

	return &Result{
		Name:         config.NormalizeSessionName(name),
		Config:       nil,
		Root:         path,
		ConfigSource: "zoxide",
	}, nil
}

func (r *Resolver) resolveFromDirectory(name string) (*Result, error) {
	for _, dir := range r.global.ProjectDirs {
		candidate := filepath.Join(r.expandTilde(dir), name)
		info, err := r.checkDir(candidate)
		if err != nil || !info.IsDir() {
			continue
		}
		return &Result{
			Name:         config.NormalizeSessionName(name),
			Config:       nil,
			Root:         candidate,
			ConfigSource: "directory",
		}, nil
	}
	return nil, fmt.Errorf("not a directory in any project_dirs: %s", name)
}

func (r *Resolver) ExpandGlob(pattern string, sessionNames []string) ([]string, error) {
	projects, err := r.store.List()
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool, len(projects))
	for _, p := range projects {
		seen[p.Name] = true
	}
	for _, name := range r.listProjectDirNames() {
		if !seen[name] {
			projects = append(projects, config.ProjectInfo{Name: name})
		}
	}

	return ExpandGlobFrom(pattern, projects, sessionNames)
}

func (r *Resolver) listProjectDirNames() []string {
	seen := make(map[string]bool)
	var names []string
	for _, dir := range r.global.ProjectDirs {
		expanded := r.expandTilde(dir)
		if expanded == "" {
			continue
		}
		entries, err := os.ReadDir(expanded)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() && !seen[e.Name()] {
				seen[e.Name()] = true
				names = append(names, e.Name())
			}
		}
	}
	return names
}

func ExpandGlobFrom(pattern string, projects []config.ProjectInfo, sessionNames []string) ([]string, error) {
	re, err := regexp.Compile("^" + strings.ReplaceAll(regexp.QuoteMeta(pattern), `\+`, ".*") + "$")
	if err != nil {
		return nil, fmt.Errorf("invalid glob pattern: %w", err)
	}

	seen := make(map[string]bool)
	var matches []string
	for _, p := range projects {
		if re.MatchString(p.Name) && !seen[p.Name] {
			seen[p.Name] = true
			matches = append(matches, p.Name)
		}
	}
	for _, name := range sessionNames {
		if re.MatchString(name) && !seen[name] {
			seen[name] = true
			matches = append(matches, name)
		}
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no projects or sessions matched pattern: %s", pattern)
	}

	sort.Strings(matches)
	return matches, nil
}

func (r *Resolver) ExpandGroup(groupName string) ([]string, error) {
	members, ok := r.global.Groups[groupName]
	if !ok {
		return nil, fmt.Errorf("group not found: %s", groupName)
	}
	return members, nil
}

func ResolveRoot(root string, projectsDir string) string {
	return resolveRootWith(root, projectsDir, os.UserHomeDir)
}

// ResolveRoots expands each dir with ~ and returns the absolute paths.
func ResolveRoots(dirs config.StringOrList) []string {
	out := make([]string, len(dirs))
	for i, d := range dirs {
		out[i] = ResolveRoot(d, "")
	}
	return out
}

func resolveRootWith(root string, projectsDir string, homeDir HomeDirFunc) string {
	expand := func(p string) string { return expandTildeWith(p, homeDir) }
	root = expand(root)
	if filepath.IsAbs(root) {
		return root
	}
	return filepath.Join(expand(projectsDir), root)
}

func (r *Resolver) expandTilde(path string) string {
	return expandTildeWith(path, r.homeDir)
}

func expandTildeWith(path string, homeDir HomeDirFunc) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}
	home, err := homeDir()
	if err != nil {
		return path
	}
	return filepath.Join(home, path[1:])
}
