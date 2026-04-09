---
title: "Custom variables"
weight: 6
---

The `vars` field in [project config]({{< relref "project-config" >}}) defines named values you can reuse across the file.

## `{{var}}` interpolation

Use double-brace names in any string field - `root`, `command`, hooks, `env` values, and window or pane fields:

```yaml
vars:
  app_name: myapp
  deploy_env: staging

root: ~/projects/{{app_name}}

env:
  APP_NAME: "{{app_name}}"
  DEPLOY_ENV: "{{deploy_env}}"

windows:
  - name: api
    command: ./bin/{{app_name}}-api
```

## Resolution order

Interpolation happens in this order:

1. **Backtick evaluation** - values wrapped in backticks run as shell commands; the output becomes the variable value
2. **Custom variable substitution** - `{{var}}` placeholders are replaced with values from the `vars` block (or CLI overrides)
3. **Environment expansion** - `${VAR}` references are expanded via Go's `os.ExpandEnv`

If a `{{var}}` placeholder has no matching key, it remains as a literal string (no error is raised).

## CLI overrides

Override variables at runtime with `--var`:

```sh
nux --var port=4000 blog
nux --var port=4000 --var env=staging blog
```

The `--var` flag is repeatable. CLI values take precedence over the YAML `vars` block for the same key.

## Dynamic values with backticks

A value wrapped in backticks runs as a shell command; stdout is used as the variable value:

```yaml
vars:
  port: "`echo 3000`"
  branch: "`git branch --show-current`"
```

Use sparingly - this runs a command every time nux resolves the config.

{{< hint warning >}}
**Trust model:** Backtick values execute arbitrary shell commands as your user. Only load config files you wrote or trust. If you receive a project config from someone else, review it before running `nux` against it.
{{< /hint >}}

## Template-style reuse

One pattern for several similar services: keep one project file per service and drive differences through `vars`:

```yaml
vars:
  service: payments
  region: us-east-1

root: ~/work/services/{{service}}

env:
  SERVICE_NAME: "{{service}}"
  AWS_REGION: "{{region}}"

windows:
  - name: dev
    command: make dev SERVICE={{service}}
```

Adjust `service` and `region` per project file or override via `--var` without duplicating layout structure.
