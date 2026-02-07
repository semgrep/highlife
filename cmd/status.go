package cmd

import (
	"fmt"
	"os"

	"github.com/semgrep/highlife/internal/state"
)

type StatusCmd struct {
	Quiet bool `optional:"" help:"Only print failures; silent if all OK."`
}

func (c *StatusCmd) Run() error {
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

	var anyFailed bool
	for _, r := range st.Results {
		if !r.Success {
			anyFailed = true
		}
	}

	if !c.Quiet {
		fmt.Printf("last sync: %s\n", st.LastSync.Format("2006-01-02 15:04:05"))
		for _, r := range st.Results {
			status := "ok"
			if !r.Success {
				status = "FAIL"
			}
			fmt.Printf("  %s\t%s\t%s\n", status, r.URL, r.Path)
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
