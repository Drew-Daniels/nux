---
title: "Troubleshooting"
weight: 10
---

Common issues and how to resolve them.

## "nux: unknown command" or "command not found"

The `nux` binary is not on your PATH. Verify the install:

```sh
which nux
nux version
```

If you installed via `go install`, make sure `$(go env GOPATH)/bin` is in your PATH. If you installed via Homebrew, run `brew link nux`.

## "cannot start session from inside tmux"

nux refuses to start new sessions from inside an existing tmux session to prevent nested session confusion. Options:

- Detach from the current session first (`prefix + d`), then run `nux`.
- Override the guard with `--force`:

```sh
nux --force blog
```

## Session name doesn't match the config filename

nux normalizes session names for tmux compatibility:

- Dots (`.`), colons (`:`), and spaces become underscores (`_`)
- Leading dashes are stripped
- Consecutive underscores are collapsed

For example, `my.api-server.yaml` becomes session `my_api-server`. Run `nux list` to see the mapping between config names and session names.

## Commands not running in panes

If `command` values in your config aren't executing:

1. Check that `pane_init` commands in global config aren't interfering. nux sends `pane_init` commands before each pane's command via `SendKeys`. If a `pane_init` command starts an interactive process, subsequent commands won't reach the shell.

2. Verify with `--dry-run` to see exactly what tmux commands nux sends:

```sh
nux --dry-run blog
```

3. Check that the `root` directory exists. If the directory doesn't exist, the pane's shell may start in a fallback directory and commands may fail.

## `on_start` runs before all windows exist

This is by design. `on_start` fires after the session and first window are created but before windows 2..N. If your startup command needs all windows to exist, use `on_ready` instead - it runs once at the end of the full session build.

## Picker doesn't appear

The interactive picker requires two conditions:

1. `picker_on_bare: true` in global config
2. You run bare `nux` (no arguments) from outside a recognized project directory

If you're inside `projects_dir`, nux auto-detects the project instead of showing the picker. Move to a directory outside `projects_dir` and try again.

Also verify your picker binary is installed:

```sh
which fzf   # or: which gum
nux doctor  # checks picker availability
```

## Zoxide-resolved projects have no layout

When a project resolves through zoxide (no config file), nux builds the session from the `default_session` template in global config. If no `default_session` is configured, you get a bare single-window session.

To add a layout, create a config file:

```sh
nux new my-project
```

## Config changes not picked up

Running sessions don't re-read configs automatically. After editing a config:

```sh
nux restart blog           # full session restart
nux restart blog:editor    # restart just one window
```

## "project not found"

nux resolves project names in this order:

1. Config file at `~/.config/nux/projects/<name>.yaml`
2. Zoxide query (if enabled)
3. Directory at `<projects_dir>/<name>`

If none match, you get this error. Check:

- The config filename matches (without `.yaml`)
- The directory exists under `projects_dir`
- Zoxide knows about the directory (`zoxide query <name>`)

Run `nux doctor` for a full environment check.

## `--dry-run` still talks to tmux

`--dry-run` prints the tmux commands nux would execute, but it still queries the live tmux server for session state (`has-session`) and option values (`base-index`, `pane-base-index`). This keeps the output accurate - it shows whether a session would be created or attached based on real state. If tmux is not running, the dry-run output will reflect that (all sessions treated as new).

## `nux show` displays interpolated values

`nux show` prints the fully resolved config after variable and environment expansion. If your config contains secrets in `env` or `vars`, the expanded values will appear in the output. Use `--raw` to print the config before interpolation:

```sh
nux show blog --raw
```

## tmux version too old

nux requires tmux 3.0 or newer. Check your version:

```sh
tmux -V
```

On macOS with Homebrew: `brew upgrade tmux`. On Linux, check your package manager or build from source.
