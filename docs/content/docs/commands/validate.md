---
title: "validate"
weight: 9
---

## Usage

```text
nux validate [name ...]
```

## Description

Validates project configuration files and reports errors. Each config is printed with **`[ok]`** or **`[error]`**.

- **With no arguments:** validates all project configs under `~/.config/nux/projects/`.
- **With one or more targets:** validates each expanded project. Targets support **glob patterns** (`+`), **`@group`** expansion, and multiple space-separated names (same rules as `nux` / `nux stop`).

### What gets checked

- `command` and `windows` are mutually exclusive at the project level.
- Every window must have a `name`.
- Every window must have at least one pane.
- `layout` values must be a recognized tmux layout name or a valid custom layout string.
- `split` values on panes must be `horizontal` or `vertical` (if set).

## Example output

```text
  [ok]    blog
  [ok]    api
  [error] workers: "command" and "windows" are mutually exclusive
```

If any config has errors, nux exits with a non-zero status.

If no project configs are found, nux prints `No project configs found.`

## Examples

```sh
# Check everything
nux validate

# Check one or more projects
nux validate blog
nux validate blog api
nux validate web+
nux validate @work
```

## Notes

- Run this after editing configs by hand, generating them from templates, or as part of a CI pipeline.
- Validation checks structural correctness but does not verify that commands, paths, or directories actually exist.
