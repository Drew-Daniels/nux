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
	LoadRaw(name string) ([]byte, string, error)
	List() ([]ProjectInfo, error)
	Save(name string, cfg *ProjectConfig) error
	SaveRaw(name string, content []byte) error
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

func (s *DirProjectStore) LoadRaw(name string) ([]byte, string, error) {
	path := s.Path(name)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, path, err
	}
	return data, path, nil
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
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	content := append([]byte(ProjectSchemaModeline), data...)
	return s.SaveRaw(name, content)
}

// SaveRaw writes a project file as-is (caller supplies full bytes, including
// the schema modeline if desired). Used for hand-authored scaffolds.
func (s *DirProjectStore) SaveRaw(name string, content []byte) error {
	if err := os.MkdirAll(s.Dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(s.Path(name), content, 0o600)
}

func (s *DirProjectStore) Delete(name string) error {
	return os.Remove(s.Path(name))
}
