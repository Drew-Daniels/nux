---
title: "delete"
weight: 8
---

## Usage

```text
nux delete <name>
```

## Description

Deletes the project config file for `<name>`. By default, nux prompts for confirmation before deleting.

This only removes the config file at `~/.config/nux/projects/<name>.yaml`. It does **not** delete the project directory, source code, or stop any running sessions.

## Flags

| Flag | Meaning |
|------|---------|
| `--force` | Skip the confirmation prompt |

## Behavior

1. nux checks that `~/.config/nux/projects/<name>.yaml` exists.
2. Unless `--force` is set, nux prompts: `Delete config for "<name>"? [y/N]`
3. If confirmed, the file is deleted and nux prints `Deleted config for "<name>"`.
4. If declined, nux prints `Cancelled.` and exits cleanly.

## Errors

- **`config not found: <path>`** - no config file exists for that project name.

## Examples

```sh
# Delete with confirmation prompt
nux delete old-project

# Delete without prompting
nux delete old-project --force
```

## Notes

- If a session for the deleted project is still running, it continues until stopped. Use `nux stop <name>` first if you want to tear down the session.
- After deletion, the project can still start via convention (directory under `projects_dir`) or zoxide, but without any custom config.
