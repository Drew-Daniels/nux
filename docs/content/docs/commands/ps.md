---
title: "ps"
weight: 5
---

## Usage

```text
nux ps
```

## Description

Shows currently running tmux sessions in a table.

| Column | Meaning |
|--------|---------|
| **NAME** | Session name |
| **WINDOWS** | Number of windows in the session |
| **ATTACHED** | `yes` if a client is attached, `no` otherwise |
| **UPTIME** | How long the session has been running (e.g. `5m`, `1h 32m`) |

If no sessions are running, nux prints `No running sessions.`

This command shows all tmux sessions, not just those started by nux.

## Example output

```text
NAME     WINDOWS   ATTACHED   UPTIME
blog     2         yes        5m
api      3         no         1h 32m
docs     1         no         45m
```

## Notes

- For the full project catalog (including projects without a running session), use [`nux list`]({{< relref "/docs/commands/list" >}}).
- Uptime is rounded to the nearest minute.
