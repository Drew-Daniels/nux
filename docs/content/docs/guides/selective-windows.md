---
title: "Selective Windows"
weight: 4
---

You can target **individual windows** with the `:window` suffix on the project name, optionally with **multiple windows** separated by commas. Window names match the `name` field under `windows` in your project config.

## Syntax

```text
nux project:window
nux project:window1,window2,...
nux restart project:window
nux restart project:window1,window2,...
```

## Starting a subset (root command)

```sh
nux myapp:editor
nux myapp:editor,server
```

Creates the session with **only** the listed windows, in **the order you list** (not necessarily YAML order). Session-level options (`env`, `on_start`, `on_ready`, hooks) still apply; `on_start` and `on_ready` run against the **first window in your list** as the anchor.

If the session **already exists**, nux selects the **first** window you listed and attaches (it does not rebuild the session or replay hooks).

You cannot combine `:window` targets with `--run`, `--layout`, or `--panes`. Projects that use the single `command` form (no `windows` array) cannot use `:window`.

## Restart

```sh
nux restart myapp:editor
nux restart myapp:editor,server
```

Recreates the listed window(s) from the project config. The rest of the session stays untouched. With multiple windows, restart runs **in the order listed**.

If a window name does not exist in the config, nux reports `"window not found in config"`.

## When to use selective windows

- Large projects where you only need a few panes right now.
- **Restart** when a single window's process is stuck or you changed its pane commands.
- Iterating on one pane while leaving others untouched.

## Notes

- `nux stop` and `nux delete` operate on the **project** name only; `blog:editor` stops or deletes the `blog` project/session config, not a single window.
- To kill a single window manually without restart, use `tmux kill-window -t session:window`.
- `nux restart` supports multiple targets, glob patterns (`+`), and `@group` expansion.
