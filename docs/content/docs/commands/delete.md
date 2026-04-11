---
title: "delete"
weight: 8
---

## Usage

```text
nux delete <name> [name ...]
```

**Alias:** `nux del`

## Description

Deletes one or more project config files. By default, nux prompts for confirmation before each deletion.

This only removes config files at `~/.config/nux/projects/<name>.yaml`. It does **not** delete project directories, source code, or stop any running sessions.

Supports [glob patterns]({{< relref "/docs/guides/pattern-matching" >}}) with `+` and [group expansion]({{< relref "/docs/configuration/session-groups" >}}) with `@`.

## Flags

| Flag | Meaning |
|------|---------|
| `--force` | Skip confirmation prompts |

## Behavior

1. Arguments are expanded (patterns and groups resolved to individual names).
2. For each target, nux checks that the config file exists.
3. Unless `--force` is set, nux prompts: `Delete config for "<name>"? [y/N]`
4. If confirmed, the file is deleted and nux prints `Deleted config for "<name>"`.
5. If declined, nux prints `Cancelled.` and moves on to the next target.

## Errors

- **`config not found: <path>`** - no config file exists for that project name.

## Examples

```sh
# Delete with confirmation prompt
nux delete old-project

# Delete without prompting
nux delete old-project --force

# Delete multiple projects
nux delete blog api docs

# Delete all projects matching a pattern
nux del web+ --force

# Delete all projects in a group
nux del @deprecated --force
```

## Notes

- If a session for the deleted project is still running, it continues until stopped. Use `nux stop <name>` first if you want to tear down the session.
- After deletion, the project can still start via convention (directory under any configured `project_dirs` path) or zoxide, but without any custom config.
