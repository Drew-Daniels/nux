---
title: "list"
weight: 4
---

## Usage

```text
nux list
```

**Alias:** `nux ls`

## Description

Lists all projects that have a config file under `~/.config/nux/projects/`, along with their current status.

| Column | Meaning |
|--------|---------|
| **NAME** | Project name (derived from the config filename) |
| **STATUS** | `running` if a tmux session is active, `-` otherwise |
| **WINDOWS** | Number of windows in the running session (blank if not running) |
| **UPTIME** | How long the session has been running (blank if not running) |
| **CONFIG** | `project` (all listed entries come from config files) |
| **ROOT** | The `root` value from the project config, if set |

Only projects with explicit config files appear here. Convention-based projects (discovered via `project_dirs` or zoxide) are not included unless they also have a config file.

## Examples

```sh
nux list
nux ls
```

Example output:

```text
NAME     STATUS    WINDOWS   UPTIME    CONFIG    ROOT
blog     running   2         5m        project   ~/projects/blog
api      -                             project   ~/code/api
docs     running   1         1h 2m     project
```

## Notes

- For only **running** tmux sessions (not the full project catalog), use [`nux ps`]({{< relref "/docs/commands/ps" >}}).
- The ROOT column shows the raw value from config before resolution, so it may contain `~` or `{{var}}` placeholders.
