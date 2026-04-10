package config

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const ProjectSchemaModeline = "# yaml-language-server: $schema=https://raw.githubusercontent.com/Drew-Daniels/nux/main/schemas/project.schema.json\n"

// ProjectStore abstracts operations on project config files, providing
// an injection seam for testing without filesystem access.
type ProjectStore interface {
	Load(name string) (*ProjectConfig, string, error)
	List() ([]ProjectInfo, error)
	Save(name string, cfg *ProjectConfig) error
	Delete(name string) error
	Path(name string) string
}

// DirProjectStore implements ProjectStore backed by a directory on disk.
type DirProjectStore struct {
	Dir string
}

func NewProjectStore(dir string) *DirProjectStore {
	return &DirProjectStore{Dir: dir}
}

func DefaultProjectStore() *DirProjectStore {
	return &DirProjectStore{Dir: ProjectConfigDir()}
}

func (s *DirProjectStore) Path(name string) string {
	return filepath.Join(s.Dir, name+".yaml")
}

func (s *DirProjectStore) Load(name string) (*ProjectConfig, string, error) {
	path := s.Path(name)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}

	var cfg ProjectConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, path, err
	}
	return &cfg, path, nil
}

func (s *DirProjectStore) List() ([]ProjectInfo, error) {
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var projects []ProjectInfo
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".yaml") {
			continue
		}
		projects = append(projects, ProjectInfo{
			Name: strings.TrimSuffix(name, ".yaml"),
			Path: filepath.Join(s.Dir, name),
		})
	}
	return projects, nil
}

func (s *DirProjectStore) Save(name string, cfg *ProjectConfig) error {
	if err := os.MkdirAll(s.Dir, 0o755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	content := append([]byte(ProjectSchemaModeline), data...)
	return os.WriteFile(s.Path(name), content, 0o644)
}

func (s *DirProjectStore) Delete(name string) error {
	return os.Remove(s.Path(name))
}
