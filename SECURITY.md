# Security Policy

## Scope

nux reads project YAML, runs lifecycle hooks, and executes tmux commands on your behalf. Security issues in any of these areas are in scope:

- Arbitrary command execution via malicious project config, hooks, or interpolated values
- Path traversal or unintended file overwrite when resolving project paths or writing configs
- Secret or sensitive data leakage in logs, `nux show` output, or error messages
- Vulnerabilities in nux's use of external tools (tmux, fzf) when introduced by nux's integration

## Reporting a Vulnerability

**Do not open a public issue for security vulnerabilities.**

Use [GitHub's private vulnerability reporting](https://github.com/Drew-Daniels/nux/security/advisories/new) to submit a report. This keeps the details confidential until a fix is available.

If you prefer email, contact **drew@drewdaniels.dev** with the subject line "nux security report".

### What to include

- nux version (`nux version`)
- tmux version (`tmux -V`)
- OS and architecture
- Description of the vulnerability
- Steps to reproduce or a proof of concept
- Impact assessment (what an attacker could do)

### Response timeline

- **Acknowledge** within 48 hours
- **Triage and initial assessment** within 1 week
- **Fix or mitigation** targeting the next patch release

## Supported Versions

Only the latest release is supported with security updates. If you are on an older version, upgrade to the latest before reporting.

## Disclosure

After a fix is released, a security advisory will be published on the GitHub repository with credit to the reporter (unless anonymity is requested).
