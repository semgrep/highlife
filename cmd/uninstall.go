package cmd

import "github.com/semgrep/highlife/internal/launchd"

type UninstallCmd struct{}

func (c *UninstallCmd) Run(g *Globals) error {
	return launchd.Uninstall()
}
