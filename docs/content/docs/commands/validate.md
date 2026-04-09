---
title: "validate"
weight: 9
---

## Usage

```text
nux validate [name]
```

## Description

Validates project configuration files and reports errors. Each config is printed with **`[ok]`** or **`[error]`**.

- **With no arguments:** validates all project configs under `~/.config/nux/projects/`.
- **With `name`:** validates only that project's config.

### What gets checked

- `command` and `windows` are mutually exclusive at the project level.
- Every window must have a `name`.
- `command` and `panes` are mutually exclusive at the window level.
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

# Check one project
nux validate blog
```

## Notes

- Run this after editing configs by hand, generating them from templates, or as part of a CI pipeline.
- Validation checks structural correctness but does not verify that commands, paths, or directories actually exist.
