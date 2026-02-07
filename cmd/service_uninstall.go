package cmd

import "github.com/tpetr/highlife/internal/launchd"

type ServiceUninstallCmd struct{}

func (c *ServiceUninstallCmd) Run() error {
	return launchd.Uninstall()
}
