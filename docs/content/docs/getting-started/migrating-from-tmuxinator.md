---
title: Migrating from tmuxinator
weight: 3
---

nux replaces tmuxinator-style YAML with its own schema and CLI. This page maps concepts and highlights what does not carry over one-to-one.

## Templating and variables

tmuxinator uses ERB (`<%= ... %>`). nux does not embed Ruby.

- **Custom template variables** use `{{var}}` syntax.
- **Environment variables** use `${VAR}` syntax (shell-style).

## Project root

In tmuxinator, `root` is often explicit. In nux, if you omit `root`, it defaults to `<first_project_dirs>/<project_name>`: the first `project_dirs` entry when global config lists multiple paths, otherwise that single path (from global config and the project name).

## Hooks

| tmuxinator        | nux        |
|-------------------|------------|
| `on_project_stop` | `on_stop`  |

nux also supports `on_start`, `on_ready`, and `on_detach`.

`on_stop` is registered as a tmux hook (`session-closed`), so it runs when the session ends even if the process exits abnormally.

## Windows and panes

- Each entry under `panes` can be a **string** (a single command) or an **object** with `root` and/or `command`.
- **Pane names** are for humans reading the config; they are not used as functional identifiers like some tmux scripts assume.

## Startup window and pane

tmuxinator has `startup_window` and `startup_pane`. nux does not: it always selects the **first window** after creation. Reorder windows in YAML if you need a different default.

## tmux integration flags

These tmuxinator fields are **not** in nux:

- `tmux_options`
- `tmux_command`
- `socket_name`

Configure tmux itself via `tmux.conf` and standard environment (for example `TMUX_TMPDIR` or socket behavior outside nux).

## Pre-window commands

tmuxinator `pre_window` runs before each pane command. The nux equivalent is **`pane_init`** in the **global** `~/.config/nux/config.yaml`, which applies to new panes according to your setup.

## Attach behavior

tmuxinator allows `attach: false` per project. nux has no per-project attach flag; use the CLI flag **`--no-attach`** when you start a session if you do not want to attach immediately.

## Session name normalization

nux normalizes session names from the config filename: dots, colons, and spaces become underscores, leading dashes are stripped, and consecutive underscores are collapsed. For example, `my.project.yaml` becomes session `my_project`. tmuxinator uses the `name` field directly.

## Config file location

| tmuxinator                    | nux                                      |
|-------------------------------|------------------------------------------|
| `~/.tmuxinator/<name>.yml`     | `~/.config/nux/projects/<name>.yaml`    |

Global options live in `~/.config/nux/config.yaml`. nux follows the XDG Base Directory Specification, so if `$XDG_CONFIG_HOME` is set, it uses that instead of `~/.config`.

## Simple single-process sessions

nux requires at least one window with `panes`, like tmuxinator. For a single TUI or long-running command, use one window and one pane string (or a pane object with `command`).

## Side-by-side example

**tmuxinator-style project (concept)**

```yaml
name: myapp
root: ~/code/myapp
pre_window: export FOO=bar
windows:
  - editor:
      layout: main-vertical
      panes:
        - vim
        - cargo watch
  - logs: tail -f log/development.log
```

**Rough nux equivalent**

```yaml
# root defaults to <first_project_dirs>/myapp if omitted (first entry when project_dirs is a list)
# put "export FOO=bar" in ~/.config/nux/config.yaml under pane_init
# if you need the same prelude in every pane (tmuxinator pre_window)
windows:
  - name: editor
    layout: main-vertical
    panes:
      - vim
      - cargo watch
  - name: logs
    panes:
      - tail -f log/development.log
```

Use `pane_init` in global config for tmuxinator-style `pre_window` when the same commands should run before every pane command. nux uses `layout` on each window where you need a tmux layout. Align window and pane lists with how you want tmux to arrange them.

After migrating files, run `nux doctor` and start with `nux <project>` (or `nux --no-attach <project>`) to validate behavior.
