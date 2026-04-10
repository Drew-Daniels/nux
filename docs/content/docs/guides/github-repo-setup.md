---
title: GitHub repository setup
weight: 99
---

This page documents maintainer-facing GitHub settings for **nux**. Most files live in the repo (templates, Dependabot, workflows); a few settings are only in the GitHub UI or API.

## Files in the repo

| Item | Location |
|------|----------|
| Security policy | [SECURITY.md](https://github.com/Drew-Daniels/nux/blob/main/SECURITY.md) |
| Code of Conduct | [CODE_OF_CONDUCT.md](https://github.com/Drew-Daniels/nux/blob/main/CODE_OF_CONDUCT.md) |
| Support / where to ask | [SUPPORT.md](https://github.com/Drew-Daniels/nux/blob/main/SUPPORT.md) |
| Issue forms | [.github/ISSUE_TEMPLATE/](https://github.com/Drew-Daniels/nux/tree/main/.github/ISSUE_TEMPLATE) |
| PR template | [.github/PULL_REQUEST_TEMPLATE.md](https://github.com/Drew-Daniels/nux/blob/main/.github/PULL_REQUEST_TEMPLATE.md) |
| Dependabot | [.github/dependabot.yml](https://github.com/Drew-Daniels/nux/blob/main/.github/dependabot.yml) |
| Sponsor button | [.github/FUNDING.yml](https://github.com/Drew-Daniels/nux/blob/main/.github/FUNDING.yml) |
| Stale issues (optional) | [.github/workflows/stale.yml](https://github.com/Drew-Daniels/nux/blob/main/.github/workflows/stale.yml) |

## Settings applied via GitHub API

These can be configured with `gh` (see project history). To reproduce or audit:

- **Discussions** - enable on the repository.
- **Dependabot security updates** - enable under security analysis.
- **Branch protection (`main`)** - required status checks should match current CI job names (see below); strict mode; **admins may bypass** (`enforce_admins` off); no force-push; no deletion.
- **Labels** - ensure `security` and `breaking` exist (defaults may already be present).

### Branch protection (re-apply)

List exact check names from a recent commit on `main`:

```sh
gh api "repos/Drew-Daniels/nux/commits/main/check-runs" --jq '.check_runs[].name' | sort -u
```

Then `PUT /repos/{owner}/{repo}/branches/main/protection` with `required_status_checks.checks` matching those names. nux CI includes **vulncheck** in addition to lint, test, vet, schemas, gitleaks, and matrix **build** jobs.

Do **not** require the **Docs** workflow `deploy` check unless every pull request runs it. On nux, docs deploy only runs on pushes to `main` when `docs/**` changes, so it is not a reliable PR gate and is omitted from required checks.

### Optional UI steps

- **Private vulnerability reporting** - Settings → Security → Code security → Private vulnerability reporting.
- **GitHub Sponsors** - complete [Sponsors](https://github.com/sponsors) onboarding so the Sponsor button from `FUNDING.yml` works.

## CI job names

The workflow file is named `CI`. Required checks use the **job id** from each job (e.g. `lint`, `test`, `vulncheck`), not the workflow display name.
