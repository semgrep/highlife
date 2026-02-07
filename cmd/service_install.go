package cmd

import "github.com/tpetr/highlife/internal/launchd"

type ServiceInstallCmd struct{}

func (c *ServiceInstallCmd) Run() error {
	return launchd.Install()
}
