---
title: Quickstart
weight: 2
description: "Set up nux in minutes - configure global defaults, create project configs, start sessions, and use groups and patterns."
---

This walkthrough assumes nux is installed and tmux 3.0+ is available. See [Installation]({{< relref "/docs/getting-started/installation" >}}) if you have not set that up yet.

## Global config

Create `~/.config/nux/config.yaml` to configure defaults:

```yaml
project_dirs: ~/projects
default_shell: /bin/zsh
pane_init:
  - eval "$(direnv hook zsh)"
default_session:
  windows:
    - name: main
      panes:
        - dev
```

- `project_dirs` - base directory or list of directories for convention-based project discovery; a string or YAML list of strings (default: `~/projects`). When multiple paths are set, nux scans all of them; the first entry is the base for relative `root` in project configs.
- `default_shell` - shell used for new panes (optional, tmux default if omitted)
- `pane_init` - commands run in each pane before pane-specific commands (optional)
- `default_session` - template for projects without a config file; an object with a `windows` array (same shape as project config), for example one window with a single pane command

If you skip the global config entirely, nux uses built-in defaults (`project_dirs: ~/projects`, `picker: fzf`).

## Convention over configuration

If `~/projects/blog` exists, you can start a session with no project file:

```sh
nux blog
```

nux resolves the directory under one of your `project_dirs` paths and applies the `default_session` layout (or creates a bare session if no `default_session` is configured).

## Auto-detect

If you are already inside a project directory, a bare `nux` resolves the current directory automatically:

```sh
cd ~/projects/blog
nux
```

This starts or attaches to the `blog` session without naming it explicitly.

## Project-specific config

Scaffold a config for more control:

```sh
nux new blog
```

This creates `~/.config/nux/projects/blog.yaml` and opens it in `$EDITOR`. A minimal example:

```yaml
windows:
  - name: editor
    panes:
      - nvim
  - name: shell
    panes:
      - ""
```

Adjust `windows` and `panes` to match how you work.

## Listing projects and sessions

See configured projects and their status:

```sh
nux list    # or: nux ls
```

```text
NAME     STATUS    WINDOWS   UPTIME    CONFIG    ROOT
blog     running   2         5m        project   ~/projects/blog
api      -                             project   ~/code/api
```

See running tmux sessions:

```sh
nux ps
```

```text
NAME     WINDOWS   ATTACHED   UPTIME
blog     2         yes        5m
api      3         no         1h 2m
```

## Starting multiple sessions

Pass multiple project names:

```sh
nux blog api docs
```

## Selective windows (start)

With a multi-window project config, start only the windows you need (comma-separated, in that order):

```sh
nux blog:editor
nux blog:editor,shell
```

See [Selective windows]({{< relref "/docs/guides/selective-windows" >}}) for details and how this interacts with hooks.

## Session groups

In `~/.config/nux/config.yaml`, define groups:

```yaml
groups:
  work:
    - blog
    - api
    - docs
```

Start everything in the group:

```sh
nux @work
```

## Stopping sessions

```sh
nux stop blog              # stop one
nux stop web+              # stop by pattern
nux stop @work             # stop a group
nux stop-all               # stop everything
```

## Restarting after config changes

```sh
nux restart blog              # full session restart
nux restart blog:editor       # restart just one window
nux restart blog:editor,shell # restart several windows in order
```

## Dry-run and diagnostics

Preview what nux would do without executing:

```sh
nux --dry-run blog
```

Check your environment:

```sh
nux doctor
```

## Running commands

Run a command instead of using the project config:

```sh
nux -x "just dev"                     # in the current directory
nux -x "just dev" blog                # in the blog project directory
nux -x "fish" -l tiled -p 4 blog     # in every pane of an ad-hoc layout
```
