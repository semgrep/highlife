package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/log"
	"github.com/semgrep/highlife/cmd"
)

type CLI struct {
	cmd.Globals

	Add       cmd.AddCmd       `cmd:"" help:"Add a Brewfile source."`
	Remove    cmd.RemoveCmd    `cmd:"" help:"Remove a Brewfile source."`
	List      cmd.ListCmd      `cmd:"" help:"List Brewfile sources."`
	Clean     cmd.CleanCmd     `cmd:"" help:"Remove all cached repo clones."`
	Sync      cmd.SyncCmd      `cmd:"" help:"Sync all sources (clone/pull + brew bundle)."`
	Status    cmd.StatusCmd    `cmd:"" help:"Show sync status."`
	Install   cmd.InstallCmd   `cmd:"" help:"Install the launchd service for automatic syncing."`
	Uninstall cmd.UninstallCmd `cmd:"" help:"Uninstall the launchd service."`
}

func main() {
	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("highlife"),
		kong.Description("Keep Homebrew packages installed and up-to-date from remote Brewfiles."),
	)
	if cli.Debug {
		log.SetLevel(log.DebugLevel)
	}
	if err := ctx.Run(&cli.Globals); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
