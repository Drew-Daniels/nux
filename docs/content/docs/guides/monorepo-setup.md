---
title: "Monorepo Setup"
weight: 1
---

In a monorepo you often want each tmux window to start in a different subdirectory while still belonging to one nux project. Use **per-window `root`** values: they are paths **relative to the project `root`**, so one config can place windows under `apps/frontend`, `apps/backend`, `packages/shared`, and so on.

## Example: project "acme-mono"

Assume the repo on disk is `~/code/acme-mono` with sub-apps under `apps/` and shared code under `packages/`. A single project config can set the session `root` to the monorepo root and give each window its own subdirectory.

```yaml
# ~/.config/nux/projects/acme-mono.yaml
root: ~/code/acme-mono

windows:
  - name: root
    # optional: omit root to stay at project root, or set "." explicitly
    root: "."
    panes:
      - echo "repo root"

  - name: frontend
    root: apps/frontend
    panes:
      - npm run dev

  - name: backend
    root: apps/backend
    panes:
      - go run ./cmd/api

  - name: shared
    root: packages/shared
    panes:
      - npm test -- --watch
```

- The session and default working directory for the project are anchored at `~/code/acme-mono`.
- Each window’s `root` is resolved under that project root, so `apps/frontend` means `~/code/acme-mono/apps/frontend`.

Use this pattern whenever you want one `nux acme-mono` (or `nux acme-mono:frontend,backend`) to spin up the whole tree without maintaining separate nux projects per package.
