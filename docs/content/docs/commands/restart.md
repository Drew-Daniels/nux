---
title: "restart"
weight: 2
---

## Usage

```text
nux restart <target> [target ...]
```

## Description

Stops each matching session, then starts it again so changes in the project config are picked up. By default, nux attaches to the **last** restarted session.

Targets support **glob patterns** (`+` matches any suffix, same as `nux` and `nux stop`), **`@group`** expansion from the global config, and multiple space-separated names.

Use the **`:window`** suffix to restart one or more windows (comma-separated) inside a session without tearing down the rest.

## Flags

| Flag | Meaning |
|------|---------|
| `--no-attach` | Do not attach to the session after restart |
| `--var key=value` | Override a custom template variable (repeatable); same as `nux` and `nux show` |

## Examples

```sh
# Full session restart (picks up config changes)
nux restart blog

# Several projects, or every project matching a pattern
nux restart blog api
nux restart web+
nux restart @work

# Restart only the "editor" window in the "blog" session
nux restart blog:editor

# Restart several windows in order
nux restart blog:editor,server

# Override {{var}} placeholders for this restart only
nux restart blog --var port=9090

# Restart without attaching
nux restart api --no-attach
```

## Behavior

**Full restart** (`nux restart blog`):

1. Kills the existing session (`tmux kill-session`)
2. Re-resolves the project config from disk
3. Rebuilds the session from the fresh config
4. Attaches (unless `--no-attach` is set)

**Window restart** (`nux restart blog:editor` or `nux restart blog:editor,server`):

1. For each listed window in order: looks up the window definition in the project config
2. Kills that window (`tmux kill-window`)
3. Creates a new window with the same name and config
4. Sets up panes, commands, and layout from the config

If a window name does not exist in the project config, nux reports `"window not found in config"`.

## Notes

- With multiple targets, **`--no-attach`** skips attaching; otherwise nux attaches only to the last session in the expanded list (sorted order for globs).
- A token like **`project+:window`** is still treated as a single glob pattern (not “expand projects then restart that window”). List targets explicitly if you need window restarts for several projects.
- Window restarts are useful when a process is stuck or you changed commands for a single window.
- **`--var`** updates the config `vars` map and re-runs `{{var}}` substitution on the loaded project. If a placeholder was already fully replaced using values from the YAML `vars` block, use **`--var`** to override those map entries before the next substitution, or edit the config on disk; see [Custom variables]({{< relref "/docs/configuration/custom-variables" >}}).
