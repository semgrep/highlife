# highlife

Keep Homebrew packages installed and up-to-date from remote Brewfiles.

`highlife` pulls Brewfiles from Git repositories and runs `brew bundle` to
install and update the packages they declare. It can run on-demand or as a
background launchd service.

## Install

    go install github.com/semgrep/highlife@latest

## Quick start

Add a Brewfile source:

    highlife add https://github.com/example/dotfiles Brewfile

Run a sync:

    highlife sync

That will shallow-clone the repository (if not already cached), check out
only the requested path, and run `brew bundle --file=<path>`.

## Global flags

    --debug          Enable debug logging.
    --brew-timeout   Maximum time to allow brew commands to run (default: 15m).

## Usage

### add URL [PATH...]

Register one or more Brewfile sources. PATH defaults to `Brewfile`.

    highlife add https://github.com/example/dotfiles
    highlife add https://github.com/example/dotfiles packages/Brewfile
    highlife add https://github.com/example/dotfiles Brewfile packages/Brewfile

Multiple sources can point to the same repository with different paths.
The repository is cloned eagerly to verify the files exist.

### remove URL [PATH]

Remove sources for the given URL. If PATH is specified, only that source is
removed; otherwise all sources for the URL are removed. If no other sources
reference the repository, its cached clone is deleted.

### list

List all configured sources.

### clean

Delete all cached repository clones. The next `sync` will re-clone as
needed.

    highlife clean

With `--debug`, prints each directory as it is removed.

### sync [--dry-run] [--delay MINUTES]

Clone or pull all configured sources, then run `brew bundle` for each.

    highlife sync
    highlife sync --debug
    highlife sync --dry-run
    highlife sync --delay 1440

`--dry-run` skips running `brew bundle`.
`--delay` skips the sync entirely if the last successful sync was less than
the given number of minutes ago. Useful for the launchd service to avoid
redundant runs.

### status [--quiet]

Show the result of the last sync, including the last successful sync time
and per-source duration.

    highlife status

With `--quiet`, only failed sources are printed and the command exits
silently on success. The exit code is non-zero if any source failed.

### install [--interval MINUTES] [--delay MINUTES]

Install a launchd user agent that runs `highlife sync` periodically.

    highlife install
    highlife install --interval 60 --delay 1440

`--interval` controls how often launchd triggers a sync (default: 30
minutes). `--delay` is passed to `sync --delay` to skip redundant runs
(default: 1440 minutes / 24 hours). The agent also runs once at load.

The plist is written to `~/Library/LaunchAgents/com.highlife.sync.plist`.

### uninstall

Unload and remove the launchd agent.

    highlife uninstall

## Files

| Path | Purpose |
|------|---------|
| `~/.config/highlife/config.json` | Source list |
| `~/.local/state/highlife/state.json` | Last sync results |
| `~/.local/state/highlife/repos/` | Cached shallow clones |
| `~/.local/state/highlife/highlife.log` | Log output from the launchd service |

Paths respect `XDG_CONFIG_HOME` and `XDG_STATE_HOME` if set.

## Platform support

`sync`, `add`, `remove`, `list`, `status`, and `clean` work on any platform
with Git and Homebrew. The `install` and `uninstall` commands are macOS-only
(launchd).
