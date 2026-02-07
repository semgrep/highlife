package cmd

import "github.com/semgrep/highlife/internal/launchd"

type ServiceUninstallCmd struct{}

func (c *ServiceUninstallCmd) Run(g *Globals) error {
	return launchd.Uninstall()
}
