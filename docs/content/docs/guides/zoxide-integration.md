---
title: "Zoxide Integration"
weight: 5
---

[zoxide](https://github.com/ajeetdsouza/zoxide) tracks directories you visit frequently. When zoxide integration is enabled, nux can resolve a project by name even when there is no matching config file, by running `zoxide query <name>` as a fallback.

## Behavior

If you have visited a directory before, `nux my-obscure-project` can succeed when that name resolves via zoxide - even without a config file at `~/.config/nux/projects/my-obscure-project.yaml`.

When a project resolves through zoxide, no project config is loaded. The session is built using the `default_session` template (if configured) or as a bare session with a single window.

## Enabling

In global config (`~/.config/nux/config.yaml`):

```yaml
zoxide: true
```

When enabled, `nux doctor` verifies that the `zoxide` binary is on your PATH.

## Resolution order

When you pass a name to nux, resolution proceeds in this order:

1. **Explicit config** - a project file matching the name under `~/.config/nux/projects/`
2. **Zoxide** - `zoxide query <name>`, if enabled
3. **Directory** - `<project_dirs>/<name>` under each configured `project_dirs` root (direct path lookup)
4. **Error** - `"project not found: <name>"`

Zoxide is checked **before** directory scanning. This means a zoxide match takes priority over a directory that happens to exist under any configured `project_dirs` path but has no config file. An explicit config file always wins.

## Limitations

- Zoxide-resolved projects do **not** appear in `nux list`. That command only shows projects with config files.
- `nux doctor` checks for the zoxide binary but does not verify that your zoxide database has entries.
- If zoxide returns an empty path, resolution falls through to the next step.
