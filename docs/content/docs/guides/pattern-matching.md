---
title: "Pattern Matching"
weight: 3
---

Many nux commands accept a **project pattern** instead of an exact name. The special character **`+` means "zero or more characters"** and can appear anywhere in the pattern.

## Where `+` can appear

- **Suffix:** `web+` matches names starting with `web` (e.g. `web-app`, `web-admin`)
- **Prefix:** `+wiki` matches names ending with `wiki` (e.g. `dev-wiki`, `team-wiki`)
- **Infix:** `+api+` matches any name containing `api` (e.g. `my-api-server`)
- **Multiple:** `w+a+` matches names starting with `w` that also contain `a`

Internally, each `+` is translated to `.*` and the pattern is anchored as `^...$`, so it must match the full project name.

## Why `+` instead of `*`

Using `+` avoids **shell glob expansion**. You do not need to quote patterns; `nux web+` works as intended without special escaping.

## What gets matched

Patterns match against project names from **config files** (`~/.config/nux/projects/*.yaml`). Convention-based projects (directories under `projects_dir` without config files) and zoxide entries are not included in pattern matching.

Results are sorted alphabetically.

## Commands that support patterns

- **`nux <pattern>`** - start or attach to all matching projects
- **`nux stop <pattern>`** - stop all matching sessions

`nux restart` does **not** support patterns - it takes exactly one session argument.

## Errors

If no projects match the pattern, nux reports:

```text
no projects matched pattern: <pattern>
```

## Examples

```sh
# Start all projects starting with "web"
nux web+

# Stop all projects containing "api"
nux stop +api+

# Start all projects ending with "wiki"
nux +wiki
```
