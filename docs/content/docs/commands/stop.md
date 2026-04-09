---
title: "stop"
weight: 1
---

## Usage

```text
nux stop <session> [session ...]
```

## Description

Stops one or more tmux sessions by name. You can pass multiple session arguments in a single invocation.

Session names accept **glob patterns** (for example `web+` matches related session names) and **group expansion** with the `@` prefix (for example `@work` expands to every session in that group).

## stop-all

```text
nux stop-all
```

Kills **every** running tmux session that nux manages. This is a blunt reset when you want no sessions left running. It is separate from `nux stop` (which requires at least one session argument).

## Examples

```sh
# Stop a single session
nux stop blog

# Stop all sessions matching a pattern
nux stop web+

# Stop every session in a named group
nux stop @work

# Stop several sessions at once
nux stop api worker frontend

# Stop every running session
nux stop-all
```

## Notes

- Pattern and group syntax matches how you start sessions, so `nux stop` mirrors `nux` start behavior for batch targets.
