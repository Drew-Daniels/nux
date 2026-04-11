---
title: "Tool Sessions"
weight: 7
---

**Tool sessions** are nux projects that exist mainly to launch a **single utility** (database TUI, log tailer, container dashboard, and so on). They are not special-cased in nux: they are normal project configs that often define **one window** and sometimes omit a meaningful repo layout.

## Typical shape

- **`root`** can be **omitted** or set to any directory (home, a log path, a repo, or `/tmp`).
- For a **single command** and one pane, use one window with a single `panes` entry.

## Example configs

**Database client (single TUI)**

```yaml
root: ~/projects/myapp
windows:
  - name: main
    panes:
      - pgcli postgresql://localhost/myapp
```

**Log follower**

```yaml
root: /var/log/myapp
windows:
  - name: main
    panes:
      - tail -f api.log
```

**Monitoring dashboard**

```yaml
windows:
  - name: main
    panes:
      - lazydocker
```

Adjust names and commands to match the tools you use.

## Discovery and commands

Tool sessions **show up in `nux list`** like any other project and work with **`nux`**, **`nux stop`**, **`nux restart`**, and patterns, so you can treat “infrastructure” panes as first-class sessions without a separate tool.
