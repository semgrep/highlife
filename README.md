# 🍾 highlife

Keep Homebrew packages installed and up-to-date from remote Brewfiles.

`highlife` pulls Brewfiles from Git repositories and runs `brew bundle` to
install and update the packages they declare. It can run on-demand or as a
background launchd service.

## Install

Via Homebrew:

    brew tap semgrep/highlife https://github.com/semgrep/highlife
    brew install semgrep/highlife/highlife

Or from source:

    go install github.com/semgrep/highlife@latest

## Quick start

Add a Brewfile source:

    highlife add https://github.com/example/dotfiles Brewfile

Run a sync:

    highlife sync

That will shallow-clone the repository (if not already cached), check out
only the requested path, and run `brew bundle --file=<path>`.

Set up automatic syncing in the background:

    highlife install

Add a shell hook to your `.zshrc` or `.bashrc` to get notified of issues:

    highlife shellhook

## Global flags

    --debug                Enable debug logging.
    --brew-timeout         Maximum time to allow brew commands to run (default: 15m).
    --skip-connectivity-check  Skip the internet connectivity check before syncing.

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

### sync [--dry-run] [--brew-min-interval MINUTES]

Pull all configured sources, then run `brew bundle` for each.

    highlife sync
    highlife sync --debug
    highlife sync --dry-run
    highlife sync --brew-min-interval 1440

`--dry-run` skips running `brew bundle`.

`--brew-min-interval` skips `brew bundle` for Brewfiles that have not
changed since the last successful sync within the given number of minutes.
Git pull always runs regardless. Changed Brewfiles are always processed
immediately, even within the interval window.

### status [--quiet]

Show the result of the last sync, including the last successful sync time
and per-source duration.

    highlife status

With `--quiet`, only failed sources are printed and the command exits
silently on success. The exit code is non-zero if any source failed.

### shellhook [--days-since-last-run DAYS] [--days-since-last-success DAYS]

Print a warning if the last sync had issues. Designed to be called from
`.zshrc` or `.bashrc`. Always exits zero so terminal startup is not
affected.

    highlife shellhook

Checks, in priority order:

1. Last sync had failures — warns and suggests `highlife status`.
2. No sync attempted in `--days-since-last-run` days (default: 7) — warns.
3. No successful sync in `--days-since-last-success` days (default: 14) — warns.

Warnings are colorized when stdout is a terminal, plain text otherwise.

### install [--interval MINUTES] [--brew-min-interval MINUTES]

Install a launchd user agent that runs `highlife sync` periodically.

    highlife install
    highlife install --interval 60 --brew-min-interval 1440

`--interval` controls how often launchd triggers a sync (default: 30
minutes). `--brew-min-interval` is passed to `sync --brew-min-interval`
to skip redundant `brew bundle` runs (default: 1440 minutes / 24 hours).

With `--run-at-load`, the agent also runs once at load (e.g. at login).

The plist is written to `~/Library/LaunchAgents/com.semgrep.highlife.sync.plist`.

### uninstall

Unload and remove the launchd agent.

    highlife uninstall

## Files

| Path | Purpose |
|------|---------|
| `~/.config/highlife/config.json` | Source list |
| `~/.local/state/highlife/state.json` | Last sync results and file hashes |
| `~/.local/state/highlife/repos/` | Cached shallow clones |
| `~/.local/state/highlife/highlife.log` | Log output from the launchd service |
| `~/.local/state/highlife/highlife.lock` | File lock to prevent concurrent syncs |

Paths respect `XDG_CONFIG_HOME` and `XDG_STATE_HOME` if set.

## Platform support

`sync`, `add`, `remove`, `list`, `status`, `clean`, and `shellhook` work on
any platform with Git and Homebrew. The `install` and `uninstall` commands
are macOS-only (launchd).
