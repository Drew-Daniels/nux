---
title: "Lifecycle Hooks"
weight: 2
---

Project configs can define four hooks that run shell commands at well-defined points in the session lifecycle.

| Hook | When it runs | Mechanism |
|------|-------------|-----------|
| `on_start` | After the session and first window are created, before remaining windows | `SendKeys` to the first pane |
| `on_stop` | When the session ends | tmux hook (`session-closed`) |
| `on_ready` | Once, at the end of the full session build | `SendKeys` to the first pane |
| `on_detach` | Each time a client detaches from the session | tmux hook (`client-detached`) |

## `on_start`

Runs **once** after the session and its first window are created. Commands are sent as keystrokes to the first pane of the first window. Use it for one-time setup like starting background services.

```yaml
on_start:
  - docker compose up -d
```

Because `on_start` uses `SendKeys`, the commands run inside the pane's shell (inheriting the session root, environment, and `pane_init` state). They are not detached background processes.

Note: `on_start` fires before windows 2..N are created. If your startup command needs all windows to exist first, use `on_ready` instead.

## `on_stop`

Registered as a tmux hook via `set-hook session-closed`. It fires when the session is torn down - whether via `nux stop`, `tmux kill-session`, or terminal crash. Commands run through tmux `run-shell` in a detached context.

```yaml
on_stop:
  - docker compose stop
```

Multiple commands are supported - each is registered as an indexed tmux hook (`session-closed[0]`, `session-closed[1]`, etc.) so all of them run.

## `on_ready`

Runs **once** at the end of session creation, after all windows, panes, and `on_start` commands have been processed. Commands are delivered via `SendKeys` to the first pane.

Use it for commands that should execute after the full layout is ready.

```yaml
on_ready:
  - echo "Session ready"
```

`on_ready` does **not** re-run on subsequent re-attaches to the session. It runs only during the initial build.

## `on_detach`

Registered as a tmux hook via `set-hook client-detached`. It fires **every** time a client detaches from the session.

```yaml
on_detach:
  - notify-send "Detached from session"
```

Multiple commands are supported - each is registered as an indexed tmux hook, just like `on_stop`.

## Hook values go through interpolation

All hook command strings support `{{var}}` custom variable expansion and `${VAR}` environment variable expansion, just like other config fields.

## Summary

| Hook | Runs once or every time? | Runs in a pane? |
|------|-------------------------|-----------------|
| `on_start` | Once (before windows 2..N) | Yes (first pane) |
| `on_stop` | Every session close | No (tmux run-shell) |
| `on_ready` | Once (after all windows) | Yes (first pane) |
| `on_detach` | Every detach | No (tmux run-shell) |

Pair `on_start` / `on_stop` for long-lived processes tied to the session lifetime. Use `on_detach` for per-client behavior like notifications.
