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

Use an object with a `windows` array for the same structure as a full [project config]({{< relref "project-config" >}}):

```yaml
default_session:
  windows:
    - name: editor
      command: nvim
    - name: term
```

This creates a two-window session for any project that lacks its own config file.

## Interaction with `pane_init`

[`pane_init`]({{< relref "global-config" >}}) commands run in every pane of every session, including default sessions. They execute before the `default_session` command or window commands.

## When default sessions are used

Default sessions apply when all of these are true:

1. You run `nux <name>` (or a name resolves through a group, pattern, or picker)
2. No config file exists at `~/.config/nux/projects/<name>.yaml`
3. The project directory is found via `projects_dir` or zoxide

If the project directory cannot be found at all, nux reports an error instead of creating a default session.
