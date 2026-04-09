package resolver

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Drew-Daniels/nux/internal/config"
)

func TestResolveRoot_Tilde(t *testing.T) {
	home, _ := os.UserHomeDir()
	got := ResolveRoot("~/code/myproject", "/ignored")
	want := filepath.Join(home, "code/myproject")
	if got != want {
		t.Errorf("ResolveRoot(~/code/myproject) = %q, want %q", got, want)
	}
}

func TestResolveRoot_Absolute(t *testing.T) {
	got := ResolveRoot("/opt/projects/foo", "/ignored")
	if got != "/opt/projects/foo" {
		t.Errorf("ResolveRoot(/opt/projects/foo) = %q, want /opt/projects/foo", got)
	}
}

func TestResolveRoot_Relative(t *testing.T) {
	got := ResolveRoot("myproject", "/home/user/projects")
	want := "/home/user/projects/myproject"
	if got != want {
		t.Errorf("ResolveRoot(myproject, /home/user/projects) = %q, want %q", got, want)
	}
}

func TestResolveRoot_RelativeWithTildeBase(t *testing.T) {
	home, _ := os.UserHomeDir()
	got := ResolveRoot("myproject", "~/projects")
	want := filepath.Join(home, "projects", "myproject")
	if got != want {
		t.Errorf("ResolveRoot(myproject, ~/projects) = %q, want %q", got, want)
	}
}

func TestExpandGroup_Valid(t *testing.T) {
	r := NewResolverWithStore(&config.GlobalConfig{
		Groups: map[string][]string{
			"work": {"alpha", "bravo", "charlie"},
		},
	}, config.NewProjectStore(t.TempDir()))

	members, err := r.ExpandGroup("work")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(members) != 3 {
		t.Fatalf("expected 3 members, got %d", len(members))
	}
	if members[0] != "alpha" || members[1] != "bravo" || members[2] != "charlie" {
		t.Errorf("unexpected members: %v", members)
	}
}

func TestExpandGroup_NotFound(t *testing.T) {
	r := NewResolverWithStore(&config.GlobalConfig{
		Groups: map[string][]string{},
	}, config.NewProjectStore(t.TempDir()))

	_, err := r.ExpandGroup("missing")
	if err == nil {
		t.Fatal("expected error for missing group")
	}
}

func TestExpandGlob_NoMatches(t *testing.T) {
	projects := []config.ProjectInfo{
		{Name: "alpha"},
		{Name: "bravo"},
	}
	_, err := ExpandGlobFrom("zzz+", projects, nil)
	if err == nil {
		t.Fatal("expected error for no matches")
	}
}

func TestExpandGlob_Empty(t *testing.T) {
	_, err := ExpandGlobFrom("foo+", nil, nil)
	if err == nil {
		t.Fatal("expected error for no matches")
	}
}

func TestExpandGlob_Matches(t *testing.T) {
	projects := []config.ProjectInfo{
		{Name: "web-api"},
		{Name: "web-router"},
		{Name: "web-util"},
		{Name: "other-thing"},
	}

	tests := []struct {
		pattern string
		want    []string
	}{
		{"web+", []string{"web-api", "web-router", "web-util"}},
		{"+util", []string{"web-util"}},
		{"+rout+", []string{"web-router"}},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			got, err := ExpandGlobFrom(tt.pattern, projects, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("got[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestExpandGlob_MatchesSessions(t *testing.T) {
	sessions := []string{"test_api", "test_frontend", "other"}

	got, err := ExpandGlobFrom("test_+", nil, sessions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"test_api", "test_frontend"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestExpandGlob_MergesProjectsAndSessions(t *testing.T) {
	projects := []config.ProjectInfo{
		{Name: "test_api"},
		{Name: "test_cli"},
	}
	sessions := []string{"test_api", "test_worker"}

	got, err := ExpandGlobFrom("test_+", projects, sessions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"test_api", "test_cli", "test_worker"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestExpandGlob_SessionsOnly_NoMatches(t *testing.T) {
	sessions := []string{"alpha", "bravo"}
	_, err := ExpandGlobFrom("zzz+", nil, sessions)
	if err == nil {
		t.Fatal("expected error for no matches")
	}
}

func testResolver(t *testing.T, cfgs map[string]*config.ProjectConfig) *Resolver {
	t.Helper()
	dir := t.TempDir()
	store := config.NewProjectStore(dir)
	for name, cfg := range cfgs {
		if err := store.Save(name, cfg); err != nil {
			t.Fatalf("saving %s: %v", name, err)
		}
	}

	projectsDir := t.TempDir()

	r := NewResolverWithStore(&config.GlobalConfig{
		ProjectsDir: projectsDir,
	}, store)
	r = r.WithHomeDir(func() (string, error) { return "/home/test", nil })
	r = r.WithDirChecker(func(path string) (os.FileInfo, error) {
		return os.Stat(path)
	})
	return r
}

func TestResolve_FromConfig(t *testing.T) {
	r := testResolver(t, map[string]*config.ProjectConfig{
		"blog": {Root: "/home/test/blog", Windows: []config.Window{{Name: "editor"}}},
	})

	result, err := r.Resolve("blog")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if result.Name != "blog" {
		t.Errorf("Name = %q, want blog", result.Name)
	}
	if result.Root != "/home/test/blog" {
		t.Errorf("Root = %q, want /home/test/blog", result.Root)
	}
	if result.ConfigSource != "project" {
		t.Errorf("ConfigSource = %q, want project", result.ConfigSource)
	}
	if result.Config == nil {
		t.Error("Config should not be nil")
	}
}

func TestResolve_FromConfig_ValidationError(t *testing.T) {
	r := testResolver(t, map[string]*config.ProjectConfig{
		"bad": {Command: "vim", Windows: []config.Window{{Name: "editor"}}},
	})

	_, err := r.Resolve("bad")
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestResolve_FromZoxide(t *testing.T) {
	r := testResolver(t, nil)
	r.global.Zoxide = true
	r = r.WithZoxideQuerier(func(name string) (string, error) {
		if name == "blog" {
			return "/home/test/blog", nil
		}
		return "", fmt.Errorf("not found")
	})

	result, err := r.Resolve("blog")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if result.ConfigSource != "zoxide" {
		t.Errorf("ConfigSource = %q, want zoxide", result.ConfigSource)
	}
	if result.Root != "/home/test/blog" {
		t.Errorf("Root = %q, want /home/test/blog", result.Root)
	}
}

func TestResolve_FromDirectory(t *testing.T) {
	projectsDir := t.TempDir()
	blogDir := filepath.Join(projectsDir, "blog")
	if err := os.Mkdir(blogDir, 0o755); err != nil {
		t.Fatal(err)
	}

	store := config.NewProjectStore(t.TempDir())
	r := NewResolverWithStore(&config.GlobalConfig{
		ProjectsDir: projectsDir,
	}, store)
	r = r.WithHomeDir(func() (string, error) { return "/home/test", nil })

	result, err := r.Resolve("blog")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if result.ConfigSource != "directory" {
		t.Errorf("ConfigSource = %q, want directory", result.ConfigSource)
	}
	if result.Root != blogDir {
		t.Errorf("Root = %q, want %q", result.Root, blogDir)
	}
	if result.Config != nil {
		t.Error("Config should be nil for directory source")
	}
}

func TestResolve_NotFound(t *testing.T) {
	r := testResolver(t, nil)

	_, err := r.Resolve("nonexistent")
	if err == nil {
		t.Fatal("expected error for not found")
	}
}

func TestResolve_ZoxideDisabled(t *testing.T) {
	r := testResolver(t, nil)
	r.global.Zoxide = false
	called := false
	r = r.WithZoxideQuerier(func(_ string) (string, error) {
		called = true
		return "/tmp/x", nil
	})

	_, _ = r.Resolve("anything")
	if called {
		t.Error("zoxide querier should not be called when disabled")
	}
}

func TestResolve_ZoxideFallsThrough(t *testing.T) {
	projectsDir := t.TempDir()
	blogDir := filepath.Join(projectsDir, "blog")
	if err := os.Mkdir(blogDir, 0o755); err != nil {
		t.Fatal(err)
	}

	store := config.NewProjectStore(t.TempDir())
	r := NewResolverWithStore(&config.GlobalConfig{
		ProjectsDir: projectsDir,
		Zoxide:      true,
	}, store)
	r = r.WithHomeDir(func() (string, error) { return "/home/test", nil })
	r = r.WithZoxideQuerier(func(_ string) (string, error) {
		return "", fmt.Errorf("not found")
	})

	result, err := r.Resolve("blog")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if result.ConfigSource != "directory" {
		t.Errorf("should fall through to directory, got %q", result.ConfigSource)
	}
}

func TestExpandGlob_Method(t *testing.T) {
	dir := t.TempDir()
	store := config.NewProjectStore(dir)
	_ = store.Save("web-api", &config.ProjectConfig{Command: "a"})
	_ = store.Save("web-ui", &config.ProjectConfig{Command: "b"})
	_ = store.Save("other", &config.ProjectConfig{Command: "c"})

	r := NewResolverWithStore(&config.GlobalConfig{}, store)
	got, err := r.ExpandGlob("web+", nil)
	if err != nil {
		t.Fatalf("ExpandGlob: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 matches, got %d: %v", len(got), got)
	}
}

func TestExpandTilde_NoTilde(t *testing.T) {
	r := NewResolverWithStore(&config.GlobalConfig{}, config.NewProjectStore(t.TempDir()))
	got := r.expandTilde("/absolute/path")
	if got != "/absolute/path" {
		t.Errorf("expandTilde = %q, want /absolute/path", got)
	}
}

func TestWithSetters(t *testing.T) {
	r := NewResolverWithStore(&config.GlobalConfig{}, config.NewProjectStore(t.TempDir()))

	r2 := r.WithZoxideQuerier(func(_ string) (string, error) { return "", nil })
	if r2 == nil {
		t.Error("WithZoxideQuerier returned nil")
	}

	r3 := r.WithInterpolator(config.NewInterpolator())
	if r3 == nil {
		t.Error("WithInterpolator returned nil")
	}

	r4 := r.WithDirChecker(os.Stat)
	if r4 == nil {
		t.Error("WithDirChecker returned nil")
	}

	r5 := r.WithHomeDir(os.UserHomeDir)
	if r5 == nil {
		t.Error("WithHomeDir returned nil")
	}
}

func TestNormalizeSessionName_Delegation(t *testing.T) {
	got := config.NormalizeSessionName("my.project:name")
	want := "my_project_name"
	if got != want {
		t.Errorf("NormalizeSessionName(my.project:name) = %q, want %q", got, want)
	}
}

func TestExpandTilde_HomeDirError(t *testing.T) {
	r := NewResolverWithStore(&config.GlobalConfig{}, config.NewProjectStore(t.TempDir()))
	r = r.WithHomeDir(func() (string, error) { return "", fmt.Errorf("no home") })
	got := r.expandTilde("~/projects")
	if got != "~/projects" {
		t.Errorf("expandTilde should return path unchanged on error, got %q", got)
	}
}

func TestExpandGlob_ListError(t *testing.T) {
	store := config.NewProjectStore("/nonexistent/path/that/does/not/exist")
	r := NewResolverWithStore(&config.GlobalConfig{}, store)
	_, err := r.ExpandGlob("web+", nil)
	if err == nil {
		t.Fatal("expected error from List failure")
	}
}

func TestResolve_SessionNameNormalized(t *testing.T) {
	r := testResolver(t, map[string]*config.ProjectConfig{
		"my.project": {Root: "/tmp", Windows: []config.Window{{Name: "editor"}}},
	})

	result, err := r.Resolve("my.project")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if result.Name != "my_project" {
		t.Errorf("Name = %q, want my_project", result.Name)
	}
}
