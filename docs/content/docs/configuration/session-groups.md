---
title: "Session groups"
weight: 3
---

Groups are named lists of project names in [global config]({{< relref "global-config" >}}) under the `groups` key. They let you start or stop several projects at once without typing each name.

## Defining groups

```yaml
groups:
  work:
    - api
    - web
    - workers
  personal:
    - blog
    - dotfiles
```

Each entry is a project name - the basename of a file in `~/.config/nux/projects/` (without `.yaml`), or any name resolvable via `projects_dir` or zoxide.

## Usage

Start every session in a group:

```sh
nux @work
```

Stop every session in a group:

```sh
nux stop @work
```

Replace `work` with any group key you defined.

## How groups are resolved

When you pass `@work`, nux expands it to the list of member names (`api`, `web`, `workers`), then resolves each name individually through the normal resolution chain (config file, zoxide, directory). If any member fails to resolve, nux reports an error for that specific project.

## Combining with other arguments

You can mix groups, patterns, and plain names in a single invocation:

```sh
nux @work blog web+
```

nux expands everything first, then starts each resolved session.

## Errors

- **`group not found: <name>`** - the group key does not exist in global config. Check your `groups` definition.
