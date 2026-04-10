# Contributing to nux

- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Security policy](SECURITY.md) (report vulnerabilities privately)
- [Support / where to ask](SUPPORT.md)

## Prerequisites

- Go 1.26+
- tmux 3.0+
- [just](https://github.com/casey/just)

## Getting Started

### Option A: just setup (universal)

```sh
just setup
```

This installs all required tools via `go install` and sets up git hooks with lefthook.

### Option B: nix develop (Nix users)

```sh
nix develop
```

Uses `flake.nix` to provide all dependencies.

## Development Workflow

```sh
just check    # fmt + lint + test
just test     # tests only
just lint     # golangci-lint only
just fmt      # format code
just cover    # tests with coverage report
just build    # build binary to bin/nux
```

## Commit Conventions

This project uses [Conventional Commits](https://www.conventionalcommits.org/). Lefthook enforces this on every commit.

```
feat: add session restore
fix: handle missing config file
docs: update CLI usage examples
```

## Pull Requests

1. Branch from `main`.
2. Keep commits atomic and focused.
3. Run `just check` before pushing.
4. Open a PR and describe what changed and why.

## Releasing

Releases are fully automated via GitHub Actions and [goreleaser](https://goreleaser.com/).

1. Ensure `main` is green.
2. Tag the release: `git tag v1.2.0`
3. Push the tag: `git push origin v1.2.0`

CI builds binaries for linux/darwin (amd64/arm64), creates a GitHub release with a changelog, publishes to the Homebrew tap, and produces deb/rpm/apk packages. No manual steps after the tag push.
