//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package cmd

import (
	"github.com/spf13/cobra"
)

// indicesCmd is the parent command for all index market data subcommands
// including tickers and aggregates. All indices-related subcommands are
// registered as children of this command.
var indicesCmd = &cobra.Command{
	Use:   "indices",
	Short: "Index market data commands",
}

// init registers the indices command as a subcommand of the root command.
func init() {
	rootCmd.AddCommand(indicesCmd)
}
