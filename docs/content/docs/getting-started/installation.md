---
title: Installation
weight: 1
---

Install nux using one of the methods below. You need **tmux 3.0 or newer** on your PATH before using nux.

## Homebrew

```sh
brew install Drew-Daniels/tap/nux
```

To upgrade later:

```sh
brew upgrade nux
```

## Go install

Requires Go 1.26 or newer:

```sh
go install github.com/Drew-Daniels/nux@latest
```

Ensure `$(go env GOPATH)/bin` is on your `PATH`.

To upgrade, run the same command again - it always installs the latest version.

## Nix

Install into your profile:

```sh
nix profile install github:Drew-Daniels/nux
```

Try without installing:

```sh
nix run github:Drew-Daniels/nux -- doctor
```

## From source

```sh
git clone https://github.com/Drew-Daniels/nux.git
cd nux
go build -o nux .
```

Move the binary somewhere on your `PATH`, or run it from the build directory.

## Platform support

nux works on **macOS** and **Linux** - anywhere tmux runs. Windows is not supported (tmux is not available on Windows without WSL).

## Verify the install

After installation, run:

```sh
nux doctor
```

This checks tmux, config paths, and other expectations so you can fix issues before starting sessions.

```sh
nux version
```

This prints the exact version, commit, and build date of your installed binary.

## Shell completions

Generate completion scripts for your shell:

**Bash**

```sh
nux completions bash > /usr/local/etc/bash_completion.d/nux
```

**Zsh**

```sh
nux completions zsh > "${fpath[1]}/_nux"
```

**Fish**

```sh
nux completions fish > ~/.config/fish/completions/nux.fish
```

Adjust paths to match your system. If you installed via Homebrew or Nix, completions may already be configured automatically.

## Uninstall

**Homebrew**

```sh
brew uninstall nux
brew untap Drew-Daniels/tap
```

**Go install**

Remove the binary from your Go bin directory:

```sh
rm "$(go env GOPATH)/bin/nux"
```

**Nix**

```sh
nix profile remove nux
```

**From source**

Delete the binary you built and, optionally, the cloned repository.

**Config cleanup**

nux stores config files under `~/.config/nux/` (or `$XDG_CONFIG_HOME/nux/`). Remove this directory to clean up all project configs and global settings:

```sh
rm -rf ~/.config/nux
```
