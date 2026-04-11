---
title: "Project config"
weight: 2
---

Each project is defined in a YAML file under the nux config directory:

```text
$XDG_CONFIG_HOME/nux/projects/<name>.yaml
```

On most systems this is `~/.config/nux/projects/<name>.yaml`. The file name (without `.yaml`) becomes the project name on the CLI.

## Session name normalization

The tmux session name is derived from the project name with these transformations:

- Dots (`.`), colons (`:`), and spaces are replaced with underscores (`_`)
- Leading dashes are stripped
- Consecutive underscores are collapsed to one
- Trailing underscores are stripped

For example, a file named `my.cool-project.yaml` produces the session name `my_cool-project`.

## Fields

### `root` (string)

Project root directory. Supports `~` expansion and `{{var}}` interpolation (see [Custom variables]({{< relref "custom-variables" >}})). If omitted, defaults to `<projects_dir>/<name>` from [global config]({{< relref "global-config" >}}).

### `command` (string)

Single-command shorthand for a one-window, one-pane session. **Mutually exclusive** with `windows`.

### `windows` (list)

Window definitions. **Mutually exclusive** with top-level `command`.

Each window object supports:

| Field | Required | Description |
|-------|----------|-------------|
| `name` | yes | Window name shown in the tmux status bar |
| `panes` | yes | List of panes (at least one). Use `panes: [""]` for a bare shell. |
| `root` | no | Working directory for this window. Relative paths resolve against the project `root`. Absolute paths and `~` paths are used as-is. |
| `layout` | no | `even-horizontal`, `even-vertical`, `main-horizontal`, `main-vertical`, `tiled`, or a custom tmux layout string |
| `env` | no | Environment variables for all panes in this window. Merged with project-level `env`; window values take precedence. See [Environment variables]({{< relref "environment-variables" >}}). |

### Panes

Each pane may be:

- A **string** (command shorthand) - e.g. `- npm run dev`
- An **object** with optional fields:

| Field | Description |
|-------|-------------|
| `root` | Working directory override for this pane. Relative to the window root. |
| `command` | Command to run in this pane. |

### `env` (map)

Environment variables applied to the tmux session via `tmux set-environment`. See [Environment variables]({{< relref "environment-variables" >}}).

### `vars` (map)

Custom variables for `{{var}}` interpolation. See [Custom variables]({{< relref "custom-variables" >}}).

### `on_start`, `on_stop`, `on_ready`, `on_detach` (lists of strings)

Shell commands run at lifecycle points for this project session. See [Lifecycle hooks]({{< relref "/docs/guides/lifecycle-hooks" >}}).

## Example

```yaml
root: ~/projects/my-api

env:
  NODE_ENV: development
  PORT: "4000"

vars:
  log_level: debug

on_start:
  - docker compose up -d

on_stop:
  - docker compose stop

windows:
  - name: api
    panes:
      - npm run dev

  - name: workers
    layout: even-vertical
    panes:
      - npm run worker:email
      - npm run worker:jobs

  - name: db
    root: docker
    panes:
      - docker compose up postgres
      - command: docker compose logs -f postgres
```

## Minimal configs

For a single-command project, the simplest config is:

```yaml
command: nvim
```

This creates one window with one pane running `nvim`, rooted at `<projects_dir>/<name>`.
