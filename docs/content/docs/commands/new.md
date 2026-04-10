---
title: "new"
weight: 6
---

## Usage

```text
nux new <name>
```

## Description

Creates a new project config file at `~/.config/nux/projects/<name>.yaml` and opens it in `$EDITOR` if set.

The generated config includes a [schema modeline]({{< relref "/docs/configuration/editor-intellisense" >}}) for editor IntelliSense and a minimal template:

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/Drew-Daniels/nux/main/schemas/project.schema.json
windows:
  - name: editor
```

Edit this file to define your windows, panes, commands, and other project settings.

## Behavior

1. If a config for `<name>` already exists, nux exits with an error (it will not overwrite).
2. The config file is written to disk.
3. nux prints `Created <path>`.
4. If `$EDITOR` is set, nux opens the new file in your editor.
5. If `$EDITOR` is not set, nux prints `hint: set $EDITOR to open new configs automatically`.

## Errors

- **`config already exists: <path>`** - a file with that name already exists. Use `nux edit <name>` to modify it, or `nux delete <name>` to remove it first.

## Examples

```sh
# Create and edit a new project config
nux new blog

# Create without opening an editor (when $EDITOR is unset)
unset EDITOR && nux new api
```

## Notes

- The project name comes from the argument, not from a field inside the YAML. The filename `blog.yaml` means the project is called `blog`.
- After creating the config, start the session with `nux blog`.
