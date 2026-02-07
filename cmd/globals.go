package cmd

import "time"

// Globals holds flags shared across all commands.
type Globals struct {
	Debug       bool          `help:"Print commands before running them." default:"false"`
	BrewTimeout time.Duration `help:"Maximum time to allow brew commands to run." default:"15m" name:"brew-timeout"`
}
