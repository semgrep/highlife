package cmd

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/semgrep/highlife/internal/state"
)

type StatusCmd struct {
	Quiet bool `optional:"" help:"Only print failures; silent if all OK."`
}

func (c *StatusCmd) Run(g *Globals) error {
	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	if st.LastSync.IsZero() {
		if !c.Quiet {
			fmt.Println("no sync has been run yet")
		}
		return nil
	}

	anyFailed := slices.ContainsFunc(st.Results, func(r state.SourceResult) bool {
		return !r.Success
	})

	if !c.Quiet {
		fmt.Printf("last sync:            %s (%s ago)\n", st.LastSync.Format("2006-01-02 15:04:05"), timeAgo(st.LastSync))
		if st.LastSuccessfulSync.IsZero() {
			fmt.Println("last successful sync: never")
		} else {
			fmt.Printf("last successful sync: %s (%s ago)\n", st.LastSuccessfulSync.Format("2006-01-02 15:04:05"), timeAgo(st.LastSuccessfulSync))
		}
		for _, r := range st.Results {
			status := "ok"
			if !r.Success {
				status = "FAIL"
			}
			fmt.Printf("  %s\t%s\t%s\t%s\n", status, r.Duration.Round(time.Millisecond), r.URL, r.Path)
			if !r.Success && r.Error != "" {
				fmt.Printf("    error: %s\n", r.Error)
			}
		}
	} else {
		// Quiet mode: only print failures.
		for _, r := range st.Results {
			if !r.Success {
				fmt.Printf("FAIL\t%s\t%s\t%s\n", r.URL, r.Path, r.Error)
			}
		}
	}

	if anyFailed {
		os.Exit(1)
	}
	return nil
}

func timeAgo(t time.Time) string {
	d := time.Since(t).Truncate(time.Second)

	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	secs := int(d.Seconds()) % 60

	switch {
	case days > 0:
		return fmt.Sprintf("%dd %dh", days, hours)
	case hours > 0:
		return fmt.Sprintf("%dh %dm", hours, mins)
	case mins > 0:
		return fmt.Sprintf("%dm %ds", mins, secs)
	default:
		return fmt.Sprintf("%ds", secs)
	}
}
