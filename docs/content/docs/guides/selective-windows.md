---
title: "Selective Windows"
weight: 4
---

You can target **individual windows** with the `:window` suffix on the project name. This lets you restart a single window without tearing down the whole session.

## Syntax

```text
nux restart project:window
```

The window name matches the `name` field under `windows` in your project config.

## Restart

```sh
nux restart myapp:editor
```

Recreates just the `editor` window from the project config. The rest of the session stays untouched. This is useful when a process is stuck or you changed window commands and want a clean pane.

If the specified window name does not exist in the config, nux reports `"window not found in config"`.

## When to use selective windows

- Large projects where a single window's process is stuck.
- Iterating on one pane's commands while leaving others untouched.
- Recovering from a hung process without losing the rest of the session.

## Notes

- Selective windows are supported by `restart` only. `nux start` and `nux stop` always operate on the full session.
- To kill a single window manually without restart, use `tmux kill-window -t session:window`.
- `nux restart` accepts exactly one argument, so you cannot batch-restart multiple sessions with patterns or groups.
