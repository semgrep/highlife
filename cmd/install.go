package cmd

import "github.com/semgrep/highlife/internal/launchd"

type InstallCmd struct {
	Interval int `help:"How often launchd runs sync, in minutes." default:"30"`
	Delay    int `help:"Delay value passed to sync (skip if last successful sync was less than this many minutes ago)." default:"1440"`
}

func (c *InstallCmd) Run(g *Globals) error {
	return launchd.Install(launchd.InstallOptions{
		IntervalMinutes: c.Interval,
		DelayMinutes:    c.Delay,
	})
}
