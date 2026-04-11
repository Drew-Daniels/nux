---
title: Comparison with alternatives
weight: 4
description: "How nux compares to tmuxinator, tmuxp, smug, sesh, tms, tmux-sessionx, and other tmux session managers."
---

Several tools manage tmux sessions. This page compares nux with the most popular alternatives so you can decide which fits your workflow.

All of these tools are good. They solve overlapping problems - starting, switching, and managing tmux sessions - and if one already works for you, there may be no reason to switch. This page focuses on where they differ.

## Config-driven session managers

These tools build tmux sessions from declarative config files (YAML/JSON).

| Feature | nux | tmuxinator | tmuxp | smug |
|---------|-----|------------|-------|------|
| Language | Go | Ruby | Python | Go |
| Config format | YAML | YAML + ERB | YAML / JSON | YAML |
| Runtime deps | tmux only | Ruby, tmux | Python, tmux | tmux only |
| Install | single binary | gem | pip / pipx | single binary |
| Batch start (multiple sessions) | `nux blog api docs` | one at a time | `tmuxp load a b` | one at a time |
| Session groups | `nux @work` | no | no | no |
| Pattern matching | `nux web+` | no | no | no |
| Zero-config sessions | yes (`project_dirs`) | no | no | no |
| Selective windows | `project:window` | no | no | `project:window` |
| Zoxide integration | built-in | no | no | no |
| Interactive picker | fzf / gum | no | no | no |
| Dry-run mode | `--dry-run` | no | no | `--debug` (log only) |
| Ephemeral sessions | `--run "cmd"` | no | no | no |
| Custom variables | `{{var}}` + `--var` | ERB (`<%= %>`) | no (use env vars) | `${var}` syntax |
| Lifecycle hooks | on_start, on_stop, on_ready, on_detach | on_project_start, on_project_stop, etc. | `before_script` | before_start, stop, attach_hook, detach_hook |
| Session freezing | no | no | `tmuxp freeze` | no |
| JSON config | no | no | yes | no |
| Plugin system | no | no | yes | no |
| Python scripting | no | no | yes (libtmux) | no |
| JSON schemas | yes | no | no | no |

## Session switchers

These tools focus on quickly creating and switching between sessions, rather than defining layouts in config files.

### sesh

[sesh](https://github.com/joshmedeski/sesh) is a Go-based session manager centered around zoxide integration and fast session switching. It auto-names sessions based on the directory, connects to existing sessions, and integrates tightly with tmux keybindings and fzf/zoxide for fuzzy finding.

**Overlap with nux:** Both integrate zoxide and support interactive pickers. Both can start sessions from directories without config files.

**Where sesh differs:** sesh is purely a session switcher - it creates bare sessions from directories but does not define window layouts, panes, or commands in config files. If you need multi-window layouts, hooks, or custom variables, sesh does not cover that.

**Where nux differs:** nux is config-driven. It defines full session layouts (windows, panes, commands, hooks, env vars) in YAML, plus batch operations, groups, and pattern matching. nux is more tool than workflow - it doesn't try to integrate into your tmux keybindings.

### tmux-sessionizer (tms)

[tmux-sessionizer](https://github.com/jrmoulton/tmux-sessionizer) is a Rust-based tool inspired by ThePrimeagen's tmux-sessionizer script. It scans configured directories for git repos, presents them in a built-in TUI picker, and opens them as sessions. It also has git worktree support, session switching with previews, and a `tms kill` command that auto-jumps to another session.

**Overlap with nux:** Both scan directories for projects and offer fuzzy session selection.

**Where tms differs:** tms is git-centric - it discovers projects by scanning for git repos and opens worktrees as separate windows. It has a custom TUI picker (not fzf/gum) with session preview. It's designed to be bound to tmux keys for fast switching.

**Where nux differs:** nux doesn't care whether a project is a git repo. It's config-driven with full YAML layouts, batch operations, and groups. tms has no equivalent of session groups, pattern matching, or declarative window/pane definitions.

### tmux-sessionx

[tmux-sessionx](https://github.com/omerxx/tmux-sessionx) is a tmux plugin (installed via TPM) that adds a fuzzy session manager popup inside tmux. It uses fzf-tmux with preview windows, and supports creating, renaming, and deleting sessions without leaving tmux. It also integrates with zoxide and tmuxinator.

**Overlap with nux:** Both offer fuzzy session selection and zoxide integration.

**Where sessionx differs:** It's a tmux plugin, not a standalone CLI. It lives inside tmux and is triggered by a keybinding (`prefix + O`). It can optionally delegate session creation to tmuxinator. The focus is on switching and managing existing sessions with a rich fzf interface - not on defining layouts.

**Where nux differs:** nux is a standalone binary that defines sessions declaratively. sessionx has no concept of config files, window layouts, batch operations, or session groups.

## tmux plugins (resurrect / continuum)

[tmux-resurrect](https://github.com/tmux-plugins/tmux-resurrect) and [tmux-continuum](https://github.com/tmux-plugins/tmux-continuum) take a different approach entirely. They save and restore your tmux environment (sessions, windows, pane layouts, running programs) automatically after system restarts.

**These solve a different problem.** Resurrect/continuum preserve state across reboots; nux (and the tools above) build sessions from scratch based on a desired state. You can use both - nux to define your ideal layouts, resurrect to recover if your machine reboots unexpectedly.

## Honorable mentions

**[tmuxifier](https://github.com/jimeh/tmuxifier)** - Shell-based (1.5k stars). Instead of YAML, you write shell scripts that call tmux commands with helper functions. More flexible than YAML if you want arbitrary shell logic, but more verbose. Been around since 2012.

**[laio](https://laio.sh/)** - Rust-based. Uses a flexbox-inspired DSL for defining pane layouts (rows, columns, proportional sizing). Supports session serialization (like tmuxp freeze) and project-local `.laio.yaml` files. Newer entry worth watching.

**[teamocil](https://github.com/remiprev/teamocil)** - Ruby-based, YAML config, similar to tmuxinator but simpler. Largely inactive (last significant update years ago). If you're considering teamocil, tmuxinator or nux are more actively maintained alternatives.

**[dmux](https://github.com/zdcthomas/dmux)** - Rust-based workspace manager with fzf integration and directory-agnostic profiles. Smaller community but a clean design.

## Where nux stands out

**Batch-first design.** nux is built for people who run many sessions at once. Session groups (`@work`), glob patterns (`web+`), and multi-argument start (`nux blog api docs`) are first-class. The alternatives are designed around one session at a time.

**Zero runtime dependencies.** nux is a single static binary. tmuxinator requires Ruby; tmuxp requires Python. smug and sesh share this advantage as Go/Rust binaries.

**Convention over configuration.** Any directory under a configured `project_dirs` path can become a session with no YAML. The config-driven alternatives require a config file for every project. sesh and tms also support configless sessions, but without the ability to define layouts when you need them.

**Integrated discovery.** The interactive picker (fzf/gum), zoxide fallback, auto-detect from the current directory, and JSON schemas for editor IntelliSense mean nux adapts to how you navigate.

## Where alternatives stand out

**tmuxinator** is the most mature (10+ years) with the largest community. ERB templating is more powerful than nux's `{{var}}` system if you need conditional logic in configs.

**tmuxp** offers a Python API (libtmux) for programmatic session control, a plugin system, and session freezing (`tmuxp freeze` snapshots a running session into a config file).

**smug** is the closest to nux in philosophy - single Go binary, YAML config, minimal surface area. It supports selective window start/stop with the same `project:window` syntax.

**sesh** is ideal if you want fast directory-based session switching integrated into tmux keybindings, without the overhead of config files.

**tms** is the best choice for git-heavy workflows with worktree support and a built-in TUI picker with session previews.

## When to choose nux

- You manage many sessions simultaneously and want batch start/stop.
- You prefer zero-config convention for most projects, with YAML only where needed.
- You want a single binary with no runtime dependencies.
- You value integrated tooling (picker, zoxide, dry-run, JSON schemas).

## When to choose something else

- You need ERB templating with conditional logic in configs - use **tmuxinator**.
- You want to script tmux from Python or need a plugin system - use **tmuxp**.
- You want the simplest possible config-driven tool and don't need batch features - use **smug**.
- You want fast directory-based session switching without config files - use **sesh**.
- You work heavily with git worktrees and want a TUI picker with previews - use **tms**.
- You want a session switcher that lives inside tmux as a plugin - use **tmux-sessionx**.
- You want flexbox-style layout DSL with session serialization - use **laio**.
- You want automatic session persistence across reboots - use **tmux-resurrect** + **tmux-continuum** (works alongside any session manager).
