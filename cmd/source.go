package cmd

type SourceCmd struct {
	Add    SourceAddCmd    `cmd:"" help:"Add a Brewfile source."`
	Remove SourceRemoveCmd `cmd:"" help:"Remove a Brewfile source."`
	List   SourceListCmd   `cmd:"" help:"List Brewfile sources."`
	Clean  SourceCleanCmd  `cmd:"" help:"Remove all cached repo clones."`
}
