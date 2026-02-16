//
// Date: 2026-02-14
// Copyright (c) 2026. All rights reserved.
//

package cmd

import (
	"github.com/spf13/cobra"
)

// stocksCmd is the parent command for all stock market data subcommands
// including open-close, bars, tickers, and market summary.
var stocksCmd = &cobra.Command{
	Use:   "stocks",
	Short: "Stock market data commands",
}

// init registers the stocks command as a subcommand of the root command.
func init() {
	rootCmd.AddCommand(stocksCmd)
}
