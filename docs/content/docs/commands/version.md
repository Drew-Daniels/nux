---
title: "version"
weight: 12
---

## Usage

```text
nux version
```

## Description

Prints build metadata: version, git commit, and build date.

## Example output

```text
nux 0.5.0
  commit: a1b2c3d
  built:  2026-04-01T12:00:00Z
```

When running a development build (via `go install` or `go build` without ldflags), the output shows placeholder values:

```text
nux dev
  commit: none
  built:  unknown
```

## Changelog

See the [GitHub releases page](https://github.com/Drew-Daniels/nux/releases) for release notes and changelogs.

## Notes

- Useful when reporting bugs or verifying that your PATH points to the expected install.
- Release builds embed version info via `-ldflags` at compile time.
