# highlife

Keep Homebrew packages installed and up-to-date from remote Brewfiles.

`highlife` pulls Brewfiles from Git repositories and runs `brew bundle` to
install and update the packages they declare. It can run on-demand or as a
background launchd service.

## Install

    go install github.com/tpetr/highlife@latest

## Quick start

Add a Brewfile source:

    highlife source add https://github.com/example/dotfiles Brewfile

Run a sync:

    highlife sync

That will shallow-clone the repository (if not already cached), check out
only the requested path, and run `brew bundle --file=<path>`.

## Usage

### source add URL [PATH]

Register a Brewfile source. PATH defaults to `Brewfile`.

    highlife source add https://github.com/example/dotfiles
    highlife source add https://github.com/example/dotfiles packages/Brewfile

Multiple sources can point to the same repository with different paths.

### source remove URL

Remove all sources for the given URL. If no other sources reference the
repository, its cached clone is deleted.

### source list

List all configured sources.

### source clean [--debug]

Delete all cached repository clones. The next `sync` will re-clone as
needed.

    highlife source clean

With `--debug`, prints each directory as it is removed.

### sync [--debug] [--dry-run]

Clone or pull all configured sources, then run `brew bundle` for each.

    highlife sync
    highlife sync --debug
    highlife sync --dry-run

`--debug` prints each shell command before it runs (git and brew).
`--dry-run` implies `--debug` but skips running `brew bundle`.

### status [--quiet]

Show the result of the last sync.

    highlife status

With `--quiet`, only failed sources are printed and the command exits
silently on success. The exit code is non-zero if any source failed.

### service install

Install a launchd user agent that runs `highlife sync` daily at 09:00 and
once at load.

    highlife service install

The plist is written to `~/Library/LaunchAgents/com.highlife.sync.plist`.

### service uninstall

Unload and remove the launchd agent.

    highlife service uninstall

## Files

| Path | Purpose |
|------|---------|
| `~/.config/highlife/config.json` | Source list |
| `~/.local/state/highlife/state.json` | Last sync results |
| `~/.local/state/highlife/repos/` | Cached shallow clones |
| `~/.local/state/highlife/highlife.log` | Log output from the launchd service |

Paths respect `XDG_CONFIG_HOME` and `XDG_STATE_HOME` if set.

## Platform support

`sync`, `source`, and `status` work on any platform with Git and Homebrew.
The `service` command is macOS-only (launchd).
