---
title: "Default session"
weight: 4
---

When a project has **no** config file under `~/.config/nux/projects/`, nux builds a session from `default_session` in [global config]({{< relref "global-config" >}}). This is why commands like `nux blog` can work with zero project file - the template provides the layout.

## No default_session set

If `default_session` is omitted from global config entirely, nux still creates a session - it just gets a single empty window with no commands. This is a bare tmux session rooted at the project directory.

## String shorthand

A single string is treated as one window running that command:

```yaml
default_session: nvim
```

This creates a session with one window. The `pane_init` commands (if configured) run first, then `nvim` is sent as a keystroke.

## Object form

Use an object with a `windows` array. Each window uses the **same fields** as in a [project config]({{< relref "project-config" >}}) (`name`, `panes`, `layout`, `env`, and so on). Every window must have at least one pane.

`default_session` itself only accepts **`command`** (single-pane fallback for the whole template) and **`windows`**. It does **not** support project-level options such as `root`, `vars`, top-level `env`, or lifecycle hooks (`on_start`, `on_ready`, …). For those, add a real project file under `projects/`.

```yaml
default_session:
  windows:
    - name: editor
      panes:
        - nvim
    - name: term
      panes:
        - ""
```

This creates a two-window session for any project that lacks its own config file.

## Panes and `split`

Multi-pane windows are defined with a window-level `panes` list. The field reference and examples are in [Project config]({{< relref "project-config" >}}) under **Panes**. Two details usually cause confusion:

### `split` only applies from the second pane onward

The **first** pane in the list is the window’s initial pane. nux does not run `split-window` for it, so a `split` field on **only** the first entry has no effect. For each **additional** pane, nux splits from the current layout and `split` chooses the direction:

| Value | tmux flag | Result |
|-------|-----------|--------|
| `vertical` (default) | `split-window -v` | New pane **below** the active one (stacked) |
| `horizontal` | `split-window -h` | New pane **beside** the active one (side by side) |

### `layout` vs per-pane `split`

- **`split`** on a pane: how **that** pane is created when it is added (one split at a time).
- **`layout`** on the window: optional **tmux layout** applied **after** all panes exist (for example `tiled`, `even-vertical`, or a custom layout string). You can use both; think of splits as building blocks and `layout` as a final arrangement pass.

### Same rules everywhere

`default_session` windows use the same pane logic as project configs. If anything is unclear, read the **Panes** table and example under [Project config]({{< relref "project-config" >}}).

## Interaction with `pane_init`

[`pane_init`]({{< relref "global-config" >}}) commands run in every pane of every session, including default sessions. They execute before the `default_session` command or window commands.

## When default sessions are used

Default sessions apply when all of these are true:

1. You run `nux <name>` (or a name resolves through a group, pattern, or picker)
2. No config file exists at `~/.config/nux/projects/<name>.yaml`
3. The project directory is found via `projects_dir` or zoxide

If the project directory cannot be found at all, nux reports an error instead of creating a default session.
