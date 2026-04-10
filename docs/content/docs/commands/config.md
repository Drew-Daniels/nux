---
title: "config"
weight: 5
---

## Usage

```text
nux config
```

## Description

Opens the global nux config in `$EDITOR`. If the config file does not exist yet, a scaffold with commented examples is created first.

The generated scaffold contains all available settings with sensible defaults. Optional settings are commented out with examples so you can enable them as needed.

## Behavior

**When no config exists:**

1. The config directory and `projects/` subdirectory are created.
2. A scaffold `config.yaml` is written with default values and commented examples.
3. nux prints `Created <path>`.
4. The file is opened in `$EDITOR`.

**When a config already exists:**

1. The file is opened in `$EDITOR`.

## Generated scaffold

The scaffold includes all global options:

```yaml
# Base directory for project discovery (supports ~ expansion).
projects_dir: ~/projects

# Fuzzy finder backend: fzf or gum.
picker: fzf

# Open the picker when nux is run with no arguments outside a project.
picker_on_bare: false

# Use zoxide for directory lookup when no config file matches.
zoxide: false
```

Optional settings like `default_shell`, `pane_init`, `default_session`, and `groups` are included as commented examples.

## Examples

```sh
# Open or create the global config
nux config

# Use a custom config directory
nux config --config-dir ~/dotfiles/nux
```
