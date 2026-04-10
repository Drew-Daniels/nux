---
title: "Global config"
weight: 1
---

The global config file controls defaults, discovery, pickers, and session groups for all projects.

## Location

nux follows the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/latest/). The config file is located at:

```text
$XDG_CONFIG_HOME/nux/config.yaml
```

On most systems this resolves to `~/.config/nux/config.yaml`. If `$XDG_CONFIG_HOME` is set to a custom path, nux uses that instead.

You can override the config directory at runtime with `--config-dir <path>`. This also changes where nux looks for project configs (`<path>/projects/`).

If the config file does not exist, nux applies built-in defaults silently (no error).

## Fields

### `projects_dir` (string)

Base directory used when discovering projects by convention. Tilde (`~`) is expanded. Default: `~/projects`.

Override at runtime with `--projects-dir <path>`.

### `default_shell` (string)

Shell passed to tmux as `default-command`, so new panes start in that shell. If omitted, tmux uses its own `default-shell` setting.

### `pane_init` (list of strings)

Commands run in **every** pane of **every** session before any pane-specific command. Useful for shell setup that should always happen, like initializing direnv or sourcing an environment file.

### `default_session`

Template used when a project has no config file (see [Default session]({{< relref "default-session" >}}), including **Panes and `split`**).

- **String shorthand:** a single value treated as a one-window command (for example `nvim`).
- **Object form:** an object with optional `command` and a `windows` array (same window shape as [project config]({{< relref "project-config" >}}), not a full project file).

If omitted, projects without config files get a bare session with a single empty window.

### `picker` (string)

Interactive picker backend: `fzf` or `gum`. Default: `fzf`.

### `picker_on_bare` (bool)

When `true`, running `nux` with no arguments **outside** a project directory opens the picker. Default: `false`.

### `zoxide` (bool)

Use [zoxide](https://github.com/ajeetdsouza/zoxide) for directory discovery as a resolver fallback. Default: `false`.

### `groups` (map of string to list of strings)

Named session groups: each key is a group name, each value is a list of project names. Used for batch start/stop (see [Session groups]({{< relref "session-groups" >}})).

## Built-in defaults

When no config file exists, nux uses:

```yaml
projects_dir: ~/projects
picker: fzf
picker_on_bare: false
zoxide: false
```

## Example configs

### Minimal

Enough to get started — everything else uses built-in defaults:

```yaml
projects_dir: ~/projects
```

### Editor-centric workflow

Two windows for every unconfigured project: an editor and a shell. Direnv hooks ensure `.envrc` files are loaded automatically.

```yaml
projects_dir: ~/code

default_shell: /bin/zsh

pane_init:
  - eval "$(direnv hook zsh)"

default_session:
  windows:
    - name: editor
      panes:
        - nvim
    - name: shell
      panes:
        - ""
```

### Split-pane development

A tiled editor window with a side terminal, plus a dedicated shell window:

```yaml
projects_dir: ~/dev

default_session:
  windows:
    - name: editor
      layout: main-vertical
      panes:
        - nvim
        - command: ""
          split: horizontal
    - name: shell
      panes:
        - ""
```

### Team with groups and zoxide

Batch-start related projects and let zoxide find directories outside `projects_dir`:

```yaml
projects_dir: ~/work

default_session: nvim

picker: fzf
picker_on_bare: true
zoxide: true

groups:
  platform:
    - api
    - web
    - admin
  infra:
    - terraform
    - k8s-configs
```

### Full kitchen-sink

```yaml
projects_dir: ~/code

default_shell: /bin/zsh

pane_init:
  - eval "$(direnv hook zsh)"

default_session:
  windows:
    - name: editor
      panes:
        - nvim
    - name: shell
      panes:
        - ""

picker: fzf
picker_on_bare: true
zoxide: true

groups:
  work:
    - api
    - web
  personal:
    - blog
    - dotfiles
```
