package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestStore_Path(t *testing.T) {
	s := NewProjectStore("/tmp/configs")
	got := s.Path("blog")
	want := "/tmp/configs/blog.yaml"
	if got != want {
		t.Errorf("Path(blog) = %q, want %q", got, want)
	}
}

func TestStore_SaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	s := NewProjectStore(dir)

	cfg := &ProjectConfig{
		Root: "~/projects/blog",
		Env:  map[string]string{"FOO": "bar"},
		Windows: []Window{
			{Name: "editor", Panes: []Pane{{Command: "nvim"}}},
		},
	}

	if err := s.Save("blog", cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, path, err := s.Load("blog")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if path != filepath.Join(dir, "blog.yaml") {
		t.Errorf("path = %q, want %q", path, filepath.Join(dir, "blog.yaml"))
	}
	if loaded.Root != "~/projects/blog" {
		t.Errorf("Root = %q, want ~/projects/blog", loaded.Root)
	}
	if loaded.Windows[0].Panes[0].Command != "nvim" {
		t.Errorf("pane command = %q, want nvim", loaded.Windows[0].Panes[0].Command)
	}
	if loaded.Env["FOO"] != "bar" {
		t.Errorf("Env[FOO] = %q, want bar", loaded.Env["FOO"])
	}
}

func TestStore_Load_NotFound(t *testing.T) {
	s := NewProjectStore(t.TempDir())
	_, _, err := s.Load("missing")
	if err == nil {
		t.Fatal("expected error for missing config")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected os.ErrNotExist, got %v", err)
	}
}

func TestStore_Load_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "bad.yaml"), []byte(":\n  :\n  - :\n    bad: ["), 0o644); err != nil {
		t.Fatal(err)
	}
	s := NewProjectStore(dir)
	_, _, err := s.Load("bad")
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestStore_List_Empty(t *testing.T) {
	s := NewProjectStore(t.TempDir())
	projects, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(projects) != 0 {
		t.Errorf("expected 0 projects, got %d", len(projects))
	}
}

func TestStore_List_MissingDir(t *testing.T) {
	s := NewProjectStore("/nonexistent/path/abc123")
	projects, err := s.List()
	if err != nil {
		t.Fatalf("List on missing dir should return nil error, got %v", err)
	}
	if projects != nil {
		t.Errorf("expected nil projects, got %v", projects)
	}
}

func TestStore_List_FiltersNonYAML(t *testing.T) {
	dir := t.TempDir()
	s := NewProjectStore(dir)

	_ = s.Save("alpha", &ProjectConfig{Windows: []Window{{Name: "main", Panes: []Pane{{Command: "a"}}}}})
	_ = os.WriteFile(filepath.Join(dir, "readme.md"), []byte("# hi"), 0o644)
	_ = os.Mkdir(filepath.Join(dir, "subdir"), 0o755)

	projects, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}
	if projects[0].Name != "alpha" {
		t.Errorf("Name = %q, want alpha", projects[0].Name)
	}
}

func TestStore_Delete(t *testing.T) {
	dir := t.TempDir()
	s := NewProjectStore(dir)

	_ = s.Save("deleteme", &ProjectConfig{Windows: []Window{{Name: "main", Panes: []Pane{{Command: "echo"}}}}})
	if err := s.Delete("deleteme"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, _, err := s.Load("deleteme")
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected file to be gone, got %v", err)
	}
}

func TestNewProjectFileContents_UnmarshalAndValidate(t *testing.T) {
	var cfg ProjectConfig
	if err := yaml.Unmarshal(NewProjectFileContents(), &cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if errs := Validate(&cfg); len(errs) > 0 {
		t.Fatalf("validate: %v", errs)
	}
	if len(cfg.Windows) != 1 || cfg.Windows[0].Name != "editor" {
		t.Fatalf("windows = %+v", cfg.Windows)
	}
	if len(cfg.Windows[0].Panes) != 1 || cfg.Windows[0].Panes[0].Command != "" {
		t.Fatalf("panes = %+v", cfg.Windows[0].Panes)
	}
}

func TestStore_SaveRaw_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	s := NewProjectStore(dir)
	if err := s.SaveRaw("blog", NewProjectFileContents()); err != nil {
		t.Fatalf("SaveRaw: %v", err)
	}
	cfg, _, err := s.Load("blog")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Windows) != 1 || cfg.Windows[0].Name != "editor" {
		t.Fatalf("loaded = %+v", cfg)
	}
}

func TestStore_Save_OmitsEmptyOptionalYAMLKeys(t *testing.T) {
	dir := t.TempDir()
	s := NewProjectStore(dir)
	cfg := &ProjectConfig{
		Windows: []Window{{Name: "editor", Panes: []Pane{{Command: ""}}}},
	}
	if err := s.Save("x", cfg); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(s.Path("x"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	for _, key := range []string{"on_stop:", "on_ready:", "on_detach:", "on_start:", "env:", "vars:", "layout:"} {
		if strings.Contains(content, key) {
			t.Errorf("did not expect %q in marshaled file:\n%s", key, content)
		}
	}
}

func TestStore_Save_CreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "dir")
	s := NewProjectStore(dir)

	err := s.Save("proj", &ProjectConfig{Windows: []Window{{Name: "main", Panes: []Pane{{Command: "echo"}}}}})
	if err != nil {
		t.Fatalf("Save should create dirs: %v", err)
	}

	_, _, err = s.Load("proj")
	if err != nil {
		t.Fatalf("Load after Save: %v", err)
	}
}
