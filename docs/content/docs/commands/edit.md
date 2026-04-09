---
title: "edit"
weight: 7
---

## Usage

```text
nux edit <name>
```

## Description

Opens the project config file for `<name>` in your `$EDITOR`. The file is located at `~/.config/nux/projects/<name>.yaml`.

This is a convenience shortcut - equivalent to running `$EDITOR ~/.config/nux/projects/<name>.yaml` yourself.

## Errors

- **`$EDITOR is not set`** - set the `EDITOR` environment variable in your shell profile (e.g. `export EDITOR=nvim`).
- **`config not found: <path>`** - no config file exists for that project name. Use `nux new <name>` to create one first.

## Examples

```sh
# Edit an existing project config
nux edit blog

# Common EDITOR values
export EDITOR=nvim    # Neovim
export EDITOR=vim     # Vim
export EDITOR=code    # VS Code
export EDITOR=nano    # Nano
```

## Notes

- After editing, run `nux validate <name>` to check for syntax errors before starting the session.
- To pick up config changes in a running session, use `nux restart <name>`.
