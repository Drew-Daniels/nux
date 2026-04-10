package config

import "gopkg.in/yaml.v3"

// GlobalConfig holds top-level settings from ~/.config/nux/config.yaml.
type GlobalConfig struct {
	DefaultShell   string              `yaml:"default_shell" json:"default_shell,omitempty" jsonschema:"description=Shell to set as tmux default-command for new sessions."`
	PaneInit       []string            `yaml:"pane_init" json:"pane_init,omitempty" jsonschema:"description=Commands run in every pane before pane-specific commands."`
	DefaultSession *DefaultSession     `yaml:"default_session" json:"default_session,omitempty" jsonschema:"description=Template used for projects without a config file. Accepts a layout string or full session definition."`
	ProjectsDir    string              `yaml:"projects_dir" json:"projects_dir,omitempty" jsonschema:"description=Base directory for project discovery. Supports ~ expansion."`
	Picker         string              `yaml:"picker" json:"picker,omitempty" jsonschema:"enum=fzf,enum=gum,description=Fuzzy finder backend for interactive session selection."`
	PickerOnBare   bool                `yaml:"picker_on_bare" json:"picker_on_bare,omitempty" jsonschema:"description=Open picker when nux is run with no arguments outside a project directory."`
	Zoxide         bool                `yaml:"zoxide" json:"zoxide,omitempty" jsonschema:"description=Use zoxide for directory discovery as a resolver fallback."`
	Groups         map[string][]string `yaml:"groups" json:"groups,omitempty" jsonschema:"description=Named groups of projects for batch operations (e.g. nux @work)."`
}

// DefaultSession is the template applied to projects without a config file.
type DefaultSession struct {
	Command string   `yaml:"command" json:"command,omitempty" jsonschema:"description=Command to run in the first pane (string shorthand form)."`
	Windows []Window `yaml:"windows" json:"windows,omitempty" jsonschema:"description=Window definitions for the default session template."`
}

// UnmarshalYAML allows DefaultSession to be a plain string or a full object.
func (ds *DefaultSession) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		ds.Command = value.Value
		return nil
	}
	type raw DefaultSession
	return value.Decode((*raw)(ds))
}

// ProjectConfig represents a single project YAML file.
type ProjectConfig struct {
	Root     string            `yaml:"root" json:"root,omitempty" jsonschema:"description=Project root directory. Supports ~ expansion and variable interpolation."`
	Command  string            `yaml:"command" json:"command,omitempty" jsonschema:"description=Command for a single-window session. Mutually exclusive with windows."`
	OnStart  []string          `yaml:"on_start" json:"on_start,omitempty" jsonschema:"description=Commands sent to the first pane after the session is created."`
	OnReady  []string          `yaml:"on_ready" json:"on_ready,omitempty" jsonschema:"description=Commands sent to the first pane once at the end of the initial session build."`
	OnDetach []string          `yaml:"on_detach" json:"on_detach,omitempty" jsonschema:"description=Commands run each time a client detaches."`
	OnStop   []string          `yaml:"on_stop" json:"on_stop,omitempty" jsonschema:"description=Commands run when the session is closed."`
	Env      map[string]string `yaml:"env" json:"env,omitempty" jsonschema:"description=Environment variables set for all panes via tmux set-environment."`
	Vars     map[string]string `yaml:"vars" json:"vars,omitempty" jsonschema:"description=Custom variables for {{var}} interpolation in config values."`
	Windows  []Window          `yaml:"windows" json:"windows,omitempty" jsonschema:"description=Window definitions. Mutually exclusive with command."`
}

// Window defines a tmux window inside a project session.
type Window struct {
	Name    string            `yaml:"name" json:"name" jsonschema:"required,description=Window name shown in the tmux status bar."`
	Root    string            `yaml:"root" json:"root,omitempty" jsonschema:"description=Working directory override. Relative paths resolve against project root."`
	Layout  string            `yaml:"layout" json:"layout,omitempty" jsonschema:"enum=even-horizontal,enum=even-vertical,enum=main-horizontal,enum=main-vertical,enum=tiled,description=Tmux pane layout. Also accepts custom tmux layout strings."`
	Command string            `yaml:"command" json:"command,omitempty" jsonschema:"description=Command for a single-pane window. Mutually exclusive with panes."`
	Env     map[string]string `yaml:"env" json:"env,omitempty" jsonschema:"description=Environment variables set in all panes of this window. Merged with project-level env; window values take precedence."`
	Panes   []Pane            `yaml:"panes" json:"panes,omitempty" jsonschema:"description=Pane definitions. Mutually exclusive with command."`
}

// Pane defines a tmux pane inside a window.
type Pane struct {
	Root    string `yaml:"root" json:"root,omitempty" jsonschema:"description=Working directory override for this pane."`
	Command string `yaml:"command" json:"command,omitempty" jsonschema:"description=Command to run in this pane."`
	Split   string `yaml:"split" json:"split,omitempty" jsonschema:"enum=horizontal,enum=vertical,description=Split direction when creating this pane. Horizontal splits side-by-side; vertical splits top-bottom. Default is vertical."`
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
