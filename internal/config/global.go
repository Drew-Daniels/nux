package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v3"
)

func DefaultConfigDir() string {
	return filepath.Join(xdg.ConfigHome, "nux")
}

func GlobalConfigPath() string {
	return filepath.Join(DefaultConfigDir(), "config.yaml")
}

func LoadGlobalFrom(path string) (*GlobalConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return GlobalDefaults(), nil
		}
		return nil, err
	}

	cfg := GlobalDefaults()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func LoadGlobal() (*GlobalConfig, error) {
	return LoadGlobalFrom(GlobalConfigPath())
}

func GlobalDefaults() *GlobalConfig {
	return &GlobalConfig{
		ProjectsDir:  "~/projects",
		Picker:       "fzf",
		PickerOnBare: false,
		Zoxide:       false,
	}
}

const GlobalSchemaModeline = "# yaml-language-server: $schema=https://raw.githubusercontent.com/Drew-Daniels/nux/main/schemas/global.schema.json\n"

func ScaffoldGlobalConfig() []byte {
	return []byte(GlobalSchemaModeline + `# nux global configuration

# Base directory for project discovery (supports ~ expansion).
projects_dir: ~/projects

# Fuzzy finder backend: fzf or gum.
picker: fzf

# Open the picker when nux is run with no arguments outside a project.
picker_on_bare: false

# Use zoxide for directory lookup when no config file matches.
zoxide: false

# Shell to set as tmux default-command for new sessions.
# default_shell: /bin/zsh

# Commands run in every pane before pane-specific commands.
# pane_init:
#   - eval "$(direnv hook zsh)"

# Template for projects without a config file.
# default_session:
#   windows:
#     - name: editor
#       command: vim

# Named groups for batch operations (e.g. nux @work).
# groups:
#   work:
#     - api
#     - frontend
`)
}
