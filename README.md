# nux

[![CI](https://github.com/Drew-Daniels/nux/actions/workflows/ci.yml/badge.svg)](https://github.com/Drew-Daniels/nux/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Drew-Daniels/nux)](https://goreportcard.com/report/github.com/Drew-Daniels/nux)
[![Go Reference](https://pkg.go.dev/badge/github.com/Drew-Daniels/nux.svg)](https://pkg.go.dev/github.com/Drew-Daniels/nux)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A modern tmux session manager written in Go.

nux manages tmux sessions declaratively through project configs, with a unified interface that handles both configured and ad-hoc projects. It replaces tmuxinator with a single binary, zero runtime dependencies beyond tmux, and first-class support for batch session management.

## Features

- **Batch operations** - start, stop, or restart multiple sessions at once with session groups (`nux @work`) and glob patterns (`nux web+`)
- **Convention over configuration** - projects in `~/projects/` work without any config file, using a default session template
- **Declarative YAML configs** - familiar format for tmuxinator users, with lifecycle hooks, environment variables, and custom variables
- **Interactive picker** - fuzzy finder integration (fzf/gum) for session selection
- **Zoxide integration** - smart directory discovery as a resolver fallback
- **Selective windows** - start or restart only the windows you name (`nux blog:editor`, `nux blog:editor,server`, `nux restart blog:editor`)
- **Pane split direction** - control whether panes split horizontally or vertically per-pane in config
- **Config inspector** - `nux show` prints the fully resolved config after interpolation and variable expansion (multiple targets, globs, and groups supported)
- **Dry-run mode** - preview tmux commands without executing
- **JSON schemas** - editor intellisense for config files via yaml-language-server
- **Doctor command** - environment validation with actionable fix suggestions

## Installation

### Homebrew (macOS, Linux)

```sh
brew install Drew-Daniels/tap/nux
```

### Go install

```sh
go install github.com/Drew-Daniels/nux@latest
```

### Nix flake

```sh
nix profile install github:Drew-Daniels/nux
```

### From source

```sh
git clone https://github.com/Drew-Daniels/nux.git
cd nux
go build -o nux .
```

## Quick Start

```sh
# Start a session for a project (uses default template if no config exists)
nux blog

# Start multiple sessions
nux blog api docs

# Start all sessions in a group
nux @work

# Run a command instead of the project config
nux -x "just dev"
nux -x "just dev" blog

# Ad-hoc layout with a command in every pane
nux -x "fish" -l tiled -p 4 blog

# Open interactive picker
nux

# Stop sessions
nux stop blog
nux stop web+
nux stop-all

# Restart a session (picks up config changes)
nux restart blog
nux restart web+
nux restart @work

# Restart just one window
nux restart blog:editor

# Start only certain windows (comma-separated, in that order)
nux blog:editor,server

# List available projects
nux list

# Show running sessions
nux ps

# Print resolved config for a project (or several: nux show web+)
nux show blog
```

## Configuration

Config files live in `~/.config/nux/` (XDG-aware):

```
~/.config/nux/
  config.yaml              # global settings and session groups
  projects/
    blog.yaml              # project-specific config
    api-server.yaml
```

### Global Config

```yaml
projects_dir: ~/projects
default_shell: /opt/homebrew/bin/fish
pane_init:
  - cls

default_session:
  windows:
    - name: editor
      layout: tiled
      panes:
        - nvim
        - ""
    - name: shell

groups:
  work:
    - web-app
    - web-api
    - web-admin

picker: fzf
picker_on_bare: true
zoxide: true
```

### Project Config

```yaml
root: ~/projects/api-server

env:
  NODE_ENV: development

on_start:
  - docker compose up -d
on_stop:
  - docker compose stop

windows:
  - name: editor
    layout: main-horizontal
    panes:
      - nvim
      - ""
  - name: stack
```

### Config Management

```sh
nux new blog           # create config from template, open in $EDITOR
nux edit blog          # edit existing config
nux delete blog        # delete config (with confirmation)
nux validate           # validate all configs
nux validate blog      # validate a specific config
nux validate web+      # validate every matching config
nux show blog          # print resolved config after interpolation
```

## Diagnostics

```sh
nux doctor             # check tmux, config, tools, completions
nux version            # print version info
```

## Shell Completions

```sh
nux completions bash > /etc/bash_completion.d/nux
nux completions zsh > "${fpath[1]}/_nux"
nux completions fish > ~/.config/fish/completions/nux.fish
```

## Editor Intellisense

Add a modeline to your config files for autocomplete and validation:

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/Drew-Daniels/nux/main/schemas/project.schema.json
root: ~/projects/my-project
windows: ...
```

## Community

- [Contributing](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md)
- [Security policy](SECURITY.md)
- [Questions](https://github.com/Drew-Daniels/nux/discussions) · [Issues](https://github.com/Drew-Daniels/nux/issues)

## License

[MIT](LICENSE)
