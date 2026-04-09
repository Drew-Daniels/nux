---
title: "Environment variables"
weight: 5
---

The `env` field in a [project config]({{< relref "project-config" >}}) sets session-scoped environment variables using `tmux set-environment`. Every pane in that session inherits them.

## Syntax

Values are strings. You can reference existing process environment variables with `${VAR}` and use default fallbacks with `${VAR:-default}`:

```yaml
env:
  NODE_ENV: development
  PORT: "3000"
  DATABASE_URL: "${DATABASE_URL_DEV}"
  LOG_LEVEL: "${LOG_LEVEL:-info}"
```

## Interpolation

Environment variable values go through the same interpolation pipeline as other config fields:

1. **Custom variables** - `{{var}}` placeholders are resolved from the `vars` block
2. **Environment expansion** - `${VAR}` references are expanded using Go's `os.ExpandEnv`

This means you can combine both:

```yaml
vars:
  service: api

env:
  SERVICE_NAME: "{{service}}"
  SERVICE_URL: "http://localhost:${PORT:-3000}/{{service}}"
```

## Window-level env

Windows can also define their own `env` field. Window-level variables are sent to each pane in that window via `send-keys`, complementing the session-level variables set by the project `env`.

```yaml
env:
  APP_ENV: development

windows:
  - name: api
    env:
      PORT: "4000"
    panes:
      - npm run dev
  - name: worker
    env:
      PORT: "5000"
    panes:
      - npm run worker
```

Use window-level `env` when different windows need different values for the same variable. Project-level `env` applies to the whole session; window-level `env` narrows the scope.

## Behavior

- Project-level variables are attached to the tmux **session** via `tmux set-environment`, not just the first pane.
- New panes, split panes, and reattached clients all see session-level values.
- Window-level variables are sent as `export` commands into each pane of that window.
- These are tmux session-level variables, separate from your shell's exported environment.

## Typical uses

- `NODE_ENV`, `RAILS_ENV`, `APP_ENV` for framework mode
- `PORT`, `HOST` for local servers
- `DATABASE_URL` for database connections
- Paths or URLs that differ per machine but should stay out of shell profiles
