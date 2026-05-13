---
title: "reset"
weight: 9
---

## Usage

```text
nux reset
```

## Description

Removes the global nux config file so you can start fresh. Use `--projects` to also remove all project config files.

Before anything is deleted, a summary of what will be removed and what will be kept is printed.

## Flags

| Flag | Meaning |
|------|---------|
| `--force` | Skip the confirmation prompt |
| `--projects` | Also remove all project configs |

## Behavior

1. nux checks that the global `config.yaml` in the config directory exists and is accessible.
2. A preview of what will be removed and what will be kept is printed.
3. Unless `--force` is set, nux prompts for confirmation.
4. The global config file is deleted.
5. If `--projects` is set, the `projects/` directory and all its contents are also deleted.

Running tmux sessions are not affected by this command.

## Errors

- **`global config <path>: ...`** - the global config file is missing, unreadable, or not a regular file (the message includes the underlying error).

## Examples

```sh
# Remove global config with confirmation prompt
nux reset

# Remove global config without prompting
nux reset --force

# Remove global config and all project configs
nux reset --projects --force
```

## Notes

- To recreate the global config after a reset, run `nux config`.
- Project configs created with `nux new` are kept unless `--projects` is passed.
- This command never stops or modifies running tmux sessions.
