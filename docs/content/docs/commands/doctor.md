---
title: "doctor"
weight: 10
---

## Usage

```text
nux doctor
```

## Description

Runs a diagnostic suite for your nux and tmux environment. Each check is reported with **`[ok]`**, **`[warn]`**, **`[fail]`**, or **`[missing]`**.

Checks performed, in order:

1. **nux version** - prints the running version, OS, and architecture
2. **tmux binary** - verifies `tmux` is on PATH and prints its version
3. **Global config** - confirms the global config loaded successfully
4. **Zoxide binary** - checked only when `zoxide: true` in global config
5. **Picker binary** - checked only when a picker is configured (fzf, gum)
6. **Config directory** - verifies `~/.config/nux/projects/` exists
7. **Projects directory** - verifies the `projects_dir` path exists
8. **Project configs** - validates every project YAML file

If any check fails, `doctor` exits with a non-zero status and prints "some checks failed."

## Example output

```text
  nux 0.5.0 (darwin/arm64)

  [ok]      tmux (/opt/homebrew/bin/tmux)
            tmux 3.5
  [ok]      global config
  [ok]      config directory (~/.config/nux/projects)
  [ok]      projects directory (~/projects)
  [ok]      4 project config(s)
  [ok]      all configs valid

All checks passed.
```

## Notes

- Run `doctor` after installation, on a new machine, or when something behaves unexpectedly.
- If the global config file is missing or malformed, `doctor` will fail during setup before any checks run. Fix the config file first, then re-run.
- Use the output when reporting bugs.
