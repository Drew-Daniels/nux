---
title: "completions"
weight: 11
---

## Usage

```text
nux completions <bash|zsh|fish>
```

## Description

Prints a **shell completion script** to standard output for the given shell. Redirect the output to the location your shell loads for completions.

## Install examples

**bash**

```sh
nux completions bash > /etc/bash_completion.d/nux
```

**zsh**

```sh
nux completions zsh > "${fpath[1]}/_nux"
```

**fish**

```sh
nux completions fish > ~/.config/fish/completions/nux.fish
```

## Notes

- If you install nux via **Homebrew** or **Nix**, completions are often installed automatically; you may not need the steps above.
- You may need elevated permissions for system-wide paths (for example `/etc/bash_completion.d`).
