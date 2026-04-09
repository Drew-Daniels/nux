---
title: "Interactive Picker"
weight: 6
---

When you run **`nux` with no arguments** from a directory that is not tied to a project, nux can open an interactive fuzzy finder instead of printing help.

## What the picker shows

The picker combines:

- **Running** tmux sessions
- **Configured** projects (files in `~/.config/nux/projects/`)

Pick an entry to start or attach to it as usual.

## Backends

| Backend | Notes |
|---------|--------|
| **fzf** | Default |
| **gum** | Alternative TUI picker |

Set the backend in global config:

```yaml
picker: fzf   # or gum
```

## Enabling the picker

The picker only activates when **both** conditions are met:

1. `picker_on_bare: true` is set in global config
2. You run bare `nux` (no subcommand, no project name) from **outside** a recognized project directory

```yaml
picker_on_bare: true
```

If `picker_on_bare` is `false` (the default), running bare `nux` outside a project directory prints the help text instead.

## Auto-detect behavior

Before the picker activates, nux checks whether your current working directory is **inside** the configured `projects_dir`. If it is, nux auto-detects the project name from the directory and starts or attaches to that session directly - the picker is skipped entirely.

For example, if `projects_dir` is `~/projects` and you run `nux` from `~/projects/blog`, nux starts the `blog` session without showing the picker.

## When to use the picker

The picker is most useful when you work across many repos and want muscle-memory `nux` from any directory. Combined with `picker_on_bare: true`, a bare `nux` invocation always gets you to a session quickly.
