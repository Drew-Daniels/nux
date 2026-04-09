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
