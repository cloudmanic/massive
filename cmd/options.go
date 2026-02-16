//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package cmd

import (
	"github.com/spf13/cobra"
)

// optionsCmd is the parent command for all options market data subcommands
// including aggregates, contracts, and snapshots. All options-related
// subcommands are registered as children of this command.
var optionsCmd = &cobra.Command{
	Use:   "options",
	Short: "Options market data commands",
}

// init registers the options command as a subcommand of the root command.
func init() {
	rootCmd.AddCommand(optionsCmd)
}
