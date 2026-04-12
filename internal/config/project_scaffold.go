package config

// NewProjectFileContents returns the default bytes written by nux new: schema
// modeline, one short comment, and a minimal valid windows definition. Optional
// project fields are omitted entirely; authors add them from the docs as needed.
func NewProjectFileContents() []byte {
	return []byte(ProjectSchemaModeline + `# Optional: root, env, vars, on_start / on_ready / on_detach / on_stop — see project-config in the docs.

windows:
  - name: editor
    panes:
      - ""  # interactive shell (tmux default); replace with a command to run on create

`)
}
