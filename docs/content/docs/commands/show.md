---
title: "show"
weight: 3
---

## Usage

```text
nux show <project>
```

## Description

Prints the fully resolved config for a project as YAML. The output reflects the state after interpolation (`{{var}}` substitution), environment expansion (`${VAR}`), and root resolution - exactly what nux would use to build the session.

If the project has no config file (resolved via `projects_dir` or zoxide), nux prints the session name, resolved root, and config source without a `config` section.

## Flags

| Flag | Meaning |
|------|---------|
| `--var key=value` | Override a custom variable. Repeatable. |

## Output fields

| Field | Description |
|-------|-------------|
| `name` | Normalized tmux session name |
| `root` | Resolved project root directory |
| `source` | How the project was resolved: `project`, `directory`, or `zoxide` |
| `config` | Full resolved config (omitted when no config file exists) |

## Examples

```sh
# Show resolved config for a project
nux show blog
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
          command: nvim
        - name: dev
          command: npm run dev
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

## Notes

- This command never starts or attaches to a session. It is read-only.
- Useful for debugging `{{var}}` interpolation, backtick command expansion, and root resolution.
- Shell completions suggest project names from config files.
