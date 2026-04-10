---
title: "Commands"
weight: 2
bookCollapseSection: true
---

nux provides subcommands for managing sessions, projects, and your environment. Running **`nux`** with no subcommand starts or attaches to sessions, or opens the interactive picker when configured.

## Command overview

| Command | Description |
|---------|-------------|
| `nux [names...]` | Start or attach to sessions |
| [`nux stop`]({{< relref "stop" >}}) | Stop sessions by name, pattern, or group |
| `nux stop-all` | Stop every running tmux session |
| [`nux restart`]({{< relref "restart" >}}) | Restart a session or a single window |
| [`nux list`]({{< relref "list" >}}) / `nux ls` | List configured projects and their status |
| [`nux ps`]({{< relref "ps" >}}) | Show running tmux sessions |
| [`nux show`]({{< relref "show" >}}) | Print resolved config for a project |
| [`nux new`]({{< relref "new" >}}) | Create a new project config file |
| [`nux edit`]({{< relref "edit" >}}) | Open a project config in `$EDITOR` |
| [`nux delete`]({{< relref "delete" >}}) | Delete a project config file |
| [`nux validate`]({{< relref "validate" >}}) | Validate project configs for structural errors |
| [`nux doctor`]({{< relref "doctor" >}}) | Run environment diagnostics |
| [`nux completions`]({{< relref "completions" >}}) | Generate shell completions |
| [`nux version`]({{< relref "version" >}}) | Print version and build info |

## Global flags

These flags apply to the root `nux` command (starting/attaching sessions):

| Flag | Short | Description |
|------|-------|-------------|
| `--run <command>` | `-x` | Run a command instead of using the project config. Combines with `--layout`/`--panes` and project names. |
| `--layout <name>` | `-l` | Apply an ad-hoc tmux layout (`tiled`, `even-horizontal`, `even-vertical`, `main-horizontal`, `main-vertical`, or a custom layout string). |
| `--panes <n>` | `-p` | Number of panes for the ad-hoc layout. Defaults to 2 if only `--layout` is given. |
| `--no-attach` | | Start session(s) without attaching. Also available on `restart`. |
| `--dry-run` | | Print the tmux commands nux would execute without actually running them. |
| `--force` | | Override the nested session guard. By default, nux refuses to start sessions from inside tmux. |
| `--config-dir <path>` | | Override the config directory path. Both global config and project configs are read from this directory (default: `~/.config/nux`). |
| `--projects-dir <path>` | | Override the `projects_dir` value from global config at runtime. |
| `--var key=value` | | Override a custom variable. Repeatable: `--var port=4000 --var env=staging`. |

## Auto-detect

If you run bare `nux` from inside a directory under your configured `projects_dir`, nux automatically detects the project name from the directory and starts or attaches to that session. No arguments needed:

```sh
cd ~/projects/blog
nux
```

If the current directory is not under `projects_dir`, nux falls through to the picker (if enabled) or prints help.

## Ad-hoc layouts

The `--layout` (`-l`) and `--panes` (`-p`) flags let you create multi-pane sessions on the fly without writing a config file. This is useful when you want a quick layout for a one-off task.

```sh
# 4 equal panes in the current directory
nux -l tiled -p 4

# 3 panes with a large left pane
nux myproject -l main-vertical -p 3
```

When only `--layout` is given, `--panes` defaults to 2. When only `--panes` is given, `--layout` defaults to `tiled`.

For projects with a config file, the ad-hoc layout fills in as a fallback for windows that don't specify their own layout. It never overrides a layout already set in the config.

## Running commands

The `--run` (or `-x`) flag creates a session that runs the given command instead of using the project config. The project name is still used for directory resolution and session naming, but the config's windows, panes, and commands are bypassed.

```sh
# Run a command in the current directory
nux -x "just dev"

# Run a command in the blog project directory
nux -x "just dev" blog

# Run a command in every pane of an ad-hoc layout
nux -x "fish" -l tiled -p 4 blog

# Start in background
nux -x "just serve" --no-attach
```

When combined with `--layout`/`--panes`, the command runs in **every** pane. Without a layout, it runs in the first pane only.

When used without a project name, the session is derived from the current directory (auto-detect) or falls through to the picker.

If you want a command to run in every pane across all sessions (not just ad-hoc ones), use [`pane_init`]({{< relref "/docs/configuration/global-config" >}}) in your global config instead.

## Nested session guard

If you are already inside a tmux session, nux refuses to start new sessions to avoid confusion. Use `--force` to override this:

```sh
nux --force blog
```

## Dry-run mode

Preview what nux would do without executing anything:

```sh
nux --dry-run @work
```

This prints every `tmux` command nux would invoke, which is useful for debugging configs or understanding the build sequence.
