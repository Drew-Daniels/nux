# Getting help

## Before you open an issue

1. Run **`nux doctor`** to verify tmux, config paths, and optional tools.
2. Run **`nux validate`** on your project config if the problem is YAML-related.
3. Check the **[documentation](https://drew-daniels.github.io/nux/)** (quickstart, configuration, commands).

Many problems are configuration or environment issues rather than bugs in nux.

## Where to ask

| Need | Where |
|------|--------|
| Questions, how-tos, troubleshooting ideas | [GitHub Discussions](https://github.com/Drew-Daniels/nux/discussions) |
| Bug reports or feature requests | [Issues](https://github.com/Drew-Daniels/nux/issues) (use the templates) |
| Security vulnerabilities | See [SECURITY.md](SECURITY.md) - **do not** file a public issue |

## What to include in a bug report

- Output of `nux version` and `tmux -V`
- OS and architecture
- Relevant config with secrets redacted
- What you expected vs what happened

Using the bug report template speeds up triage.
