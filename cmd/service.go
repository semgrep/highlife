package cmd

type ServiceCmd struct {
	Install   ServiceInstallCmd   `cmd:"" help:"Install the launchd service for automatic syncing."`
	Uninstall ServiceUninstallCmd `cmd:"" help:"Uninstall the launchd service."`
}
