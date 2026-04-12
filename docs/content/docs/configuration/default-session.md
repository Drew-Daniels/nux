---
title: "Default session"
weight: 4
---

When a project has **no** config file under `~/.config/nux/projects/`, nux builds a session from `default_session` in [global config]({{< relref "global-config" >}}). This is why commands like `nux blog` can work with zero project file - the template provides the layout.

## No default_session set

If `default_session` is omitted from global config entirely, nux still creates a session - it just gets a single empty window with no commands. This is a bare tmux session rooted at the project directory.

## Object form

Use an object with a `windows` array. Each window uses the **same fields** as in a [project config]({{< relref "project-config" >}}) (`name`, `panes`, `layout`, `env`, and so on). Every window must have at least one pane.

`default_session` only accepts **`windows`**. It does **not** support project-level options such as `root`, `vars`, top-level `env`, or lifecycle hooks (`on_start`, `on_ready`, …). For those, add a real project file under `projects/`.

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

## Panes and layout

Multi-pane windows are defined with a window-level `panes` list. The first pane occupies the initial window; each additional entry creates a new pane via `split-window`. Use `layout` on the window to rearrange all panes after they are created (for example `tiled`, `even-vertical`, or a custom tmux layout string).

`default_session` windows follow the same rules as project configs. See the **Panes** table under [Project config]({{< relref "project-config" >}}) for the full field reference.

## Interaction with `pane_init`

[`pane_init`]({{< relref "global-config" >}}) in the global config runs in every pane of every session, including default sessions, before each pane’s command.

Project-level [`pane_init`]({{< relref "project-config" >}}) applies only when a project YAML is loaded; default sessions (no project file) use global `pane_init` only.

## When default sessions are used

Default sessions apply when all of these are true:

1. You run `nux <name>` (or a name resolves through a group, pattern, or picker)
2. No config file exists at `~/.config/nux/projects/<name>.yaml`
3. The project directory is found via `project_dirs` path(s) or zoxide

If the project directory cannot be found at all, nux reports an error instead of creating a default session.
