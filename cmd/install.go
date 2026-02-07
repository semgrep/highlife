package cmd

import "github.com/semgrep/highlife/internal/launchd"

type InstallCmd struct {
	Interval  int  `help:"How often launchd runs sync, in minutes." default:"30"`
	BrewMinInterval int `help:"Value passed to sync --brew-min-interval (skip brew bundle for unchanged Brewfiles within this many minutes)." default:"1440" name:"brew-min-interval"`
	RunAtLoad bool `help:"Run sync immediately when the launch agent is loaded (e.g. at login)." default:"false"`
}

func (c *InstallCmd) Run(g *Globals) error {
	return launchd.Install(launchd.InstallOptions{
		IntervalMinutes: c.Interval,
		BrewMinIntervalMins: c.BrewMinInterval,
		RunAtLoad:       c.RunAtLoad,
	})
}
