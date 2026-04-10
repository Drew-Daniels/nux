---
title: "show"
weight: 3
---

## Usage

```text
nux show <target> [target ...]
```

## Description

Prints the fully resolved config for one or more projects as YAML. The output reflects the state after interpolation (`{{var}}` substitution), environment expansion (`${VAR}`), and root resolution - exactly what nux would use to build the session.

Targets support **glob patterns** (`+`), **`@group`** expansion, and multiple space-separated names. Multiple projects are written as a **YAML stream**: one document per project, separated by a `---` line.

If a project has no config file (resolved via `projects_dir` or zoxide), nux prints the session name, resolved root, and config source without a `config` section.

## Flags

| Flag | Meaning |
|------|---------|
| `--var key=value` | Override a custom variable. Repeatable. Ignored with `--raw`. |
| `--raw` | Print the config before interpolation (no variable or env expansion). Useful for inspecting configs that contain secrets. |

## Output fields

| Field | Description |
|-------|-------------|
| `name` | Normalized tmux session name |
| `root` | Resolved project root directory |
| `source` | How the project was resolved: `project`, `directory`, or `zoxide`. With `--raw`, this shows the config file path. |
| `config` | Full resolved config (omitted when no config file exists) |

## Examples

```sh
# Show resolved config for a project
nux show blog

# Several projects (YAML stream)
nux show web+
nux show @work
```

```yaml
name: blog
root: /home/user/projects/blog
source: project
config:
    env:
        NODE_ENV: development
    windows:
        - name: editor
          panes:
            - nvim
        - name: dev
          panes:
            - npm run dev
```

```sh
# Show a project resolved by directory convention (no config file)
nux show scratch
```

```yaml
name: scratch
root: /home/user/projects/scratch
source: directory
```

```sh
# Override a variable to see how it resolves
nux show blog --var port=4000
```

```sh
# Show the raw config without interpolation (secrets stay unexpanded)
nux show blog --raw
```

## Notes

- This command never starts or attaches to a session. It is read-only.
- Useful for debugging `{{var}}` interpolation, backtick command expansion, and root resolution.
- Use `--raw` to inspect configs that contain secrets in `env` or `vars` without expanding them into terminal output.
- `--raw` only works for projects that have a config file. With multiple targets, each document is raw independently.
- Shell completions suggest project names from config files.
