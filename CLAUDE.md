# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What is nux?

nux is a tmux session manager written in Go. It builds tmux sessions declaratively from YAML configs, with support for batch operations, ad-hoc layouts, zoxide integration, and interactive pickers (fzf/gum).

## Common Commands

```bash
just check          # fmt + lint + test (run before committing)
just test           # go test -race ./...
just lint           # golangci-lint run
just fmt            # gofmt -w .
just build          # go build -o bin/nux .
just install        # go install .
just schemas        # regenerate JSON schemas from config types
just setup          # install all dev tools and git hooks
just docs           # start Hugo dev server for docs
```

Run a single test:
```bash
go test -race -run TestFunctionName ./cmd/...
go test -race -run TestFunctionName ./internal/config/...
```

## Architecture

See [AGENTS.md](AGENTS.md) for detailed architecture, key patterns, and conventions.

Quick summary:
- **`cmd/`** - Cobra CLI commands. Each file is one subcommand. `root.go` owns flag parsing, `setup()`, and the `deps` struct.
- **`internal/config/`** - YAML config types, loading, validation, interpolation, JSON schema generation (`gen/`).
- **`internal/tmux/`** - `Client` interface, `RealClient` (shells out to tmux), `MockClient`, and `Builder` (session construction).
- **`internal/resolver/`** - Project name resolution: config file -> zoxide -> directory.
- **`internal/picker/`** - Fuzzy finder backends (fzf, gum).
- **`internal/ui/`** - Terminal UI helpers (tables, prompts).
- **`docs/`** - Hugo documentation site.
- **`schemas/`** - Auto-generated JSON schemas for config files.

## Critical Patterns

- **Dependency injection via `deps` struct**: All command handlers call `setup()` to get `*deps` containing the tmux client, builder, resolver, store, and function pointers. Tests inject mocks through this seam using `testDeps()`.
- **`Client` interface**: All tmux interaction goes through `tmux.Client`. Never call tmux directly outside the client.
- **Builder caches tmux options**: `BaseIndex()` and `PaneBaseIndex()` are queried once at `NewBuilder` time. Never hardcode `"0"` as a window or pane index - use `b.firstWindow()` and `b.paneBase()`.
- **Error accumulation**: Builder collects errors into `[]error` and returns `errors.Join(errs...)`. Individual failures don't abort the build.
- **Session name normalization**: Names come from config filenames, not YAML fields. `NormalizeSessionName` converts dots/colons/spaces to underscores.

## Commits

Conventional Commits format is enforced by lefthook:
```
feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert(scope)?: description
```

## Pre-commit Hooks (lefthook)

- **pre-commit**: gofmt, golangci-lint, schema regeneration (if types.go changed), gitleaks
- **commit-msg**: conventional commit format enforced
- **pre-push**: tests with race detector, 60% coverage threshold
