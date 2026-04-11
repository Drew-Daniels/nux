package config

import (
	"fmt"

	"github.com/invopop/jsonschema"
	"gopkg.in/yaml.v3"
)

// GlobalConfig holds top-level settings from ~/.config/nux/config.yaml.
type GlobalConfig struct {
	DefaultShell   string              `yaml:"default_shell" json:"default_shell,omitempty" jsonschema:"description=Shell to set as tmux default-command for new sessions."`
	PaneInit       []string            `yaml:"pane_init" json:"pane_init,omitempty" jsonschema:"description=Commands run in every pane before pane-specific commands."`
	DefaultSession *DefaultSession     `yaml:"default_session" json:"default_session,omitempty" jsonschema:"description=Template used for projects without a config file. Accepts a layout string or full session definition."`
	ProjectDirs    StringOrList        `yaml:"project_dirs" json:"project_dirs,omitempty" jsonschema:"description=Directories for project discovery. A single string or a list of paths. Supports ~ expansion."`
	Picker         string              `yaml:"picker" json:"picker,omitempty" jsonschema:"enum=fzf,enum=gum,description=Fuzzy finder backend for interactive session selection."`
	PickerOnBare   bool                `yaml:"picker_on_bare" json:"picker_on_bare,omitempty" jsonschema:"description=Open picker when nux is run with no arguments outside a project directory."`
	Zoxide         bool                `yaml:"zoxide" json:"zoxide,omitempty" jsonschema:"description=Use zoxide for directory discovery as a resolver fallback."`
	Groups         map[string][]string `yaml:"groups" json:"groups,omitempty" jsonschema:"description=Named groups of projects for batch operations (e.g. nux @work)."`
}

// StringOrList holds one or more strings. It unmarshals from either a single
// YAML string or a list of strings, so config authors can write:
//
//	project_dirs: ~/projects
//	# or
//	project_dirs:
//	  - ~/projects
//	  - ~/work
type StringOrList []string

func (s *StringOrList) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		*s = StringOrList{value.Value}
		return nil
	}
	var list []string
	if err := value.Decode(&list); err != nil {
		return err
	}
	*s = list
	return nil
}

// FirstProjectDir returns the first configured project directory, or an
// empty string if none are configured. Used as the base for relative root
// resolution in project configs.
func (g *GlobalConfig) FirstProjectDir() string {
	if len(g.ProjectDirs) > 0 {
		return g.ProjectDirs[0]
	}
	return ""
}

func (StringOrList) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{Type: "string", Description: "A single directory path."},
			{Type: "array", Items: &jsonschema.Schema{Type: "string"}, Description: "A list of directory paths."},
		},
		Description: "Directories for project discovery. A single string or a list of paths. Supports ~ expansion.",
	}
}

// DefaultSession is the template applied to projects without a config file.
type DefaultSession struct {
	Windows []Window `yaml:"windows" json:"windows,omitempty" jsonschema:"description=Window definitions for the default session template."`
}

// UnmarshalYAML rejects the removed string shorthand with an actionable message.
func (ds *DefaultSession) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		return fmt.Errorf("default_session must be an object with a windows array, not a string; use default_session: {windows: [{name: main, panes: [%s]}]}", value.Value)
	}
	type raw DefaultSession
	return value.Decode((*raw)(ds))
}

// ProjectConfig represents a single project YAML file.
type ProjectConfig struct {
	Root     string            `yaml:"root" json:"root,omitempty" jsonschema:"description=Project root directory. Supports ~ expansion and variable interpolation."`
	OnStart  []string          `yaml:"on_start" json:"on_start,omitempty" jsonschema:"description=Commands sent to the first pane after the session is created."`
	OnReady  []string          `yaml:"on_ready" json:"on_ready,omitempty" jsonschema:"description=Commands sent to the first pane once at the end of the initial session build."`
	OnDetach []string          `yaml:"on_detach" json:"on_detach,omitempty" jsonschema:"description=Commands run each time a client detaches."`
	OnStop   []string          `yaml:"on_stop" json:"on_stop,omitempty" jsonschema:"description=Commands run when the session is closed."`
	Env      map[string]string `yaml:"env" json:"env,omitempty" jsonschema:"description=Environment variables set for all panes via tmux set-environment."`
	Vars     map[string]string `yaml:"vars" json:"vars,omitempty" jsonschema:"description=Custom variables for {{var}} interpolation in config values."`
	Windows  []Window          `yaml:"windows" json:"windows" jsonschema:"required,minItems=1,description=Window definitions for the session."`
}

// UnmarshalYAML rejects the removed command field with an actionable message.
func (pc *ProjectConfig) UnmarshalYAML(value *yaml.Node) error {
	type raw ProjectConfig
	if err := value.Decode((*raw)(pc)); err != nil {
		return err
	}
	if value.Kind == yaml.MappingNode {
		for i := 0; i < len(value.Content)-1; i += 2 {
			if value.Content[i].Value == "command" {
				return fmt.Errorf("\"command\" is not a valid project field; use windows instead (e.g. windows: [{name: main, panes: [%s]}])", value.Content[i+1].Value)
			}
		}
	}
	return nil
}

// Window defines a tmux window inside a project session.
type Window struct {
	Name   string            `yaml:"name" json:"name" jsonschema:"required,description=Window name shown in the tmux status bar."`
	Root   string            `yaml:"root" json:"root,omitempty" jsonschema:"description=Working directory override. Relative paths resolve against project root."`
	Layout string            `yaml:"layout" json:"layout,omitempty" jsonschema:"description=Tmux pane layout. Named layouts or custom tmux layout strings."`
	Env    map[string]string `yaml:"env" json:"env,omitempty" jsonschema:"description=Environment variables set in all panes of this window. Merged with project-level env; window values take precedence."`
	Panes  []Pane            `yaml:"panes" json:"panes" jsonschema:"required,minItems=1,description=Pane definitions. Every window must have at least one pane."`
}

// Window.JSONSchemaExtend replaces the layout enum with a oneOf that accepts
// both named layouts and custom tmux layout strings (hex dimension prefix).
func (Window) JSONSchemaExtend(s *jsonschema.Schema) {
	if layout, ok := s.Properties.Get("layout"); ok {
		layout.Enum = nil
		layout.OneOf = []*jsonschema.Schema{
			{Type: "string", Enum: []any{"even-horizontal", "even-vertical", "main-horizontal", "main-vertical", "tiled"}},
			{Type: "string", Pattern: `^[0-9a-f]{4},`, Description: "Custom tmux layout string."},
		}
	}
}

// UnmarshalYAML rejects the removed window-level command field with an
// actionable error message directing users to panes.
func (w *Window) UnmarshalYAML(value *yaml.Node) error {
	type raw Window
	if err := value.Decode((*raw)(w)); err != nil {
		return err
	}
	if value.Kind == yaml.MappingNode {
		for i := 0; i < len(value.Content)-1; i += 2 {
			if value.Content[i].Value == "command" {
				return fmt.Errorf("window %q: \"command\" is not a valid window field; use panes instead (e.g. panes: [%s])", w.Name, value.Content[i+1].Value)
			}
		}
	}
	return nil
}

// Pane defines a tmux pane inside a window.
type Pane struct {
	Root    string `yaml:"root" json:"root,omitempty" jsonschema:"description=Working directory override for this pane."`
	Command string `yaml:"command" json:"command,omitempty" jsonschema:"description=Command to run in this pane."`
}

func (Pane) JSONSchema() *jsonschema.Schema {
	objProps := jsonschema.NewProperties()
	objProps.Set("root", &jsonschema.Schema{
		Type:        "string",
		Description: "Working directory override for this pane.",
	})
	objProps.Set("command", &jsonschema.Schema{
		Type:        "string",
		Description: "Command to run in this pane.",
	})
	return &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{Type: "string", Description: "Command shorthand — equivalent to {command: <value>}."},
			{
				Type:                 "object",
				Properties:           objProps,
				AdditionalProperties: jsonschema.FalseSchema,
			},
		},
		Description: "Pane definition. A plain string is shorthand for {command: <string>}.",
	}
}

// UnmarshalYAML allows a Pane to be a plain command string or a full object.
func (p *Pane) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		p.Command = value.Value
		return nil
	}
	type raw Pane
	return value.Decode((*raw)(p))
}

// ProjectInfo is a lightweight reference to a discovered project config.
type ProjectInfo struct {
	Name string
	Path string
}
