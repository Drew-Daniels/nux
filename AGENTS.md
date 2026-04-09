# Agent Guidelines for nux

## Project overview

nux is a tmux session manager written in Go. It builds tmux sessions declaratively from YAML configs, with support for batch operations, ad-hoc layouts, zoxide integration, and interactive pickers.

## Architecture

```
cmd/           CLI commands (cobra). Each file is one subcommand.
               root.go owns flag parsing, setup(), and the deps struct.
internal/
  config/      YAML config types, loading, validation, interpolation.
  tmux/        Client interface, RealClient (shelling out to tmux),
               MockClient, and Builder (session construction logic).
  resolver/    Project name resolution: config file -> zoxide -> directory.
  picker/      Fuzzy finder backends (fzf, gum).
  ui/          Terminal UI helpers (tables, prompts).
docs/          Hugo documentation site.
schemas/       JSON schemas for config files.
```

## Key patterns

- **Dependency injection via `deps`**: All command handlers receive a `*deps` struct built by `setup()`. This holds the tmux client, builder, resolver, store, and all runtime config. Tests inject mocks through this seam.
- **`Client` interface**: All tmux interaction goes through `tmux.Client`. `RealClient` shells out to tmux. `MockClient` records calls for assertions. Never call tmux directly outside the client.
- **`Builder` caches tmux options**: `BaseIndex()` and `PaneBaseIndex()` are queried once at `NewBuilder` construction time and stored as fields. Do not call `client.BaseIndex()` or `client.PaneBaseIndex()` per-build.
- **Error accumulation**: Builder methods collect errors into `[]error` and return `errors.Join(errs...)`. Individual tmux command failures don't abort the build.
- **Config normalization**: Session names are normalized via `NormalizeSessionName` (dots, colons, spaces become underscores). The project name comes from the config filename, not a field inside the YAML.

## tmux compatibility

nux must work with any user tmux config. Do not hardcode tmux option values:

- Use `b.firstWindow()` for the first window index (respects `base-index`).
- Use `b.paneBase()` for the first pane index (respects `pane-base-index`).
- Never hardcode `"0"` as a window or pane index.

## Testing

- Run `go test -race ./...` or `just test`.
- Run `just check` (fmt + lint + test) before committing.
- Builder tests use `MockClient`. Set `BaseIndexReturn` and `PaneBaseIndexReturn` to test non-default tmux configurations.
- Test helpers: `assertCalled`, `assertCalledWith`, `callsFor`.

## Conventions

- **Commits**: Conventional Commits format. Lefthook enforces this.
- **Lint**: golangci-lint with errcheck, govet, staticcheck, unused, ineffassign, misspell, revive.
- **No exported-symbol doc comments enforced** (revive `exported` rule is disabled).
- **Formatting**: `gofmt`.

## Adding a new command

1. Create `cmd/<name>.go` with a cobra command.
2. Call `setup()` at the top of the handler to get `*deps`.
3. Register the command in `init()` with `rootCmd.AddCommand`.
4. Add a doc page at `docs/content/docs/commands/<name>.md`.
5. Add the command to the table in `docs/content/docs/commands/_index.md`.

## Adding a new config field

1. Add the field to the struct in `internal/config/types.go` with yaml, json, and jsonschema tags.
2. If it needs validation, add checks in `internal/config/validate.go`.
3. If it needs interpolation, add it to `applyTransform` in `internal/config/interpolate.go`.
4. Wire it into the builder in `internal/tmux/builder.go`.
5. Regenerate JSON schemas with `just schemas`.
6. Update `docs/content/docs/configuration/project-config.md`.

## Documentation

- Docs live under `docs/` as a Hugo site.
- When adding features, update the relevant doc pages and `README.md`.
- Run `just docs` to preview locally.
