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

	Source  cmd.SourceCmd  `cmd:"" help:"Manage Brewfile sources."`
	Sync    cmd.SyncCmd    `cmd:"" help:"Sync all sources (clone/pull + brew bundle)."`
	Status  cmd.StatusCmd  `cmd:"" help:"Show sync status."`
	Service cmd.ServiceCmd `cmd:"" help:"Manage the launchd background service."`
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
