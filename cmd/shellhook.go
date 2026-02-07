package cmd

import (
	"fmt"
	"slices"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/semgrep/highlife/internal/state"
)

type ShellHookCmd struct {
	DaysSinceLastSuccess int `help:"Warn if last successful sync was more than this many days ago." default:"14" name:"days-since-last-success"`
	DaysSinceLastRun     int `help:"Warn if no sync has been attempted in this many days." default:"7" name:"days-since-last-run"`
}

var warnStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)

const suggestion = "ensure highlife is running in the background via `highlife install` or run manually via `highlife sync`"

func (c *ShellHookCmd) Run(g *Globals) error {
	st, err := state.Load()
	if err != nil {
		return nil // swallow errors silently
	}

	if st.LastSync.IsZero() {
		return nil
	}

	// Check if the most recent sync had any failures.
	anyFailed := slices.ContainsFunc(st.Results, func(r state.SourceResult) bool {
		return !r.Success
	})
	if anyFailed {
		fmt.Println(warnStyle.Render("highlife: last sync had failures") + " — run `highlife status` for details")
		return nil
	}

	// Check if no sync has been attempted recently.
	if time.Since(st.LastSync) > time.Duration(c.DaysSinceLastRun)*24*time.Hour {
		fmt.Println(warnStyle.Render(fmt.Sprintf("highlife: no sync attempted in over %d days", c.DaysSinceLastRun)) + " — " + suggestion)
		return nil
	}

	// Check if last successful sync is stale.
	if !st.LastSuccessfulSync.IsZero() && time.Since(st.LastSuccessfulSync) > time.Duration(c.DaysSinceLastSuccess)*24*time.Hour {
		fmt.Println(warnStyle.Render(fmt.Sprintf("highlife: no successful sync in over %d days", c.DaysSinceLastSuccess)) + " — " + suggestion)
		return nil
	}

	return nil
}
