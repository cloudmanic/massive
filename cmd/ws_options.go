//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// wsOptionsCmd is the parent command for all options WebSocket streaming
// subcommands including trades, quotes, aggregates, and fair market value.
var wsOptionsCmd = &cobra.Command{
	Use:   "options",
	Short: "Stream real-time options data",
}

// wsOptionsTradesCmd streams real-time options trade events over a WebSocket
// connection. Accepts one or more option contract tickers as positional
// arguments or the --all flag to subscribe to all trades.
// Usage: massive ws options trades O:SPY241220P00720000
var wsOptionsTradesCmd = &cobra.Command{
	Use:   "trades [tickers...]",
	Short: "Stream real-time options trades",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}
		tickers := args
		if all {
			tickers = []string{"*"}
		}
		params := buildOptionsSubscriptionParams("T", tickers)
		return connectAndStreamAsset(cmd.Context(), "options", "T", params, formatTrade)
	},
}

// wsOptionsQuotesCmd streams real-time options NBBO quote events over a
// WebSocket connection. Maximum of 1000 contracts per connection.
// Usage: massive ws options quotes O:SPY241220P00720000
var wsOptionsQuotesCmd = &cobra.Command{
	Use:   "quotes [tickers...]",
	Short: "Stream real-time options quotes (max 1000 contracts/connection)",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}
		tickers := args
		if all {
			tickers = []string{"*"}
		}
		params := buildOptionsSubscriptionParams("Q", tickers)
		return connectAndStreamAsset(cmd.Context(), "options", "Q", params, formatQuote)
	},
}

// wsOptionsAggMinuteCmd streams per-minute aggregate bar events for options
// contracts over a WebSocket connection.
// Usage: massive ws options agg-minute O:SPY241220P00720000
var wsOptionsAggMinuteCmd = &cobra.Command{
	Use:   "agg-minute [tickers...]",
	Short: "Stream per-minute options aggregates",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}
		tickers := args
		if all {
			tickers = []string{"*"}
		}
		params := buildOptionsSubscriptionParams("AM", tickers)
		return connectAndStreamAsset(cmd.Context(), "options", "AM", params, formatAggregate)
	},
}

// wsOptionsAggSecondCmd streams per-second aggregate bar events for options
// contracts over a WebSocket connection.
// Usage: massive ws options agg-second O:SPY241220P00720000
var wsOptionsAggSecondCmd = &cobra.Command{
	Use:   "agg-second [tickers...]",
	Short: "Stream per-second options aggregates",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}
		tickers := args
		if all {
			tickers = []string{"*"}
		}
		params := buildOptionsSubscriptionParams("A", tickers)
		return connectAndStreamAsset(cmd.Context(), "options", "A", params, formatAggregate)
	},
}

// wsOptionsFMVCmd streams Fair Market Value events for options contracts
// over a WebSocket connection.
// Usage: massive ws options fmv O:SPY241220P00720000
var wsOptionsFMVCmd = &cobra.Command{
	Use:   "fmv [tickers...]",
	Short: "Stream options Fair Market Value",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}
		tickers := args
		if all {
			tickers = []string{"*"}
		}
		params := buildOptionsSubscriptionParams("FMV", tickers)
		return connectAndStreamAsset(cmd.Context(), "options", "FMV", params, formatFMV)
	},
}

// buildOptionsSubscriptionParams constructs the subscription parameter string
// for a WebSocket subscribe message. Each ticker is prefixed with the channel
// name and a dot separator (e.g., "T.O:SPY241220P00720000").
func buildOptionsSubscriptionParams(channel string, tickers []string) string {
	parts := make([]string, len(tickers))
	for i, t := range tickers {
		parts[i] = channel + "." + t
	}
	return strings.Join(parts, ",")
}

// init registers the options WebSocket command and all its subcommands under
// the ws parent command.
func init() {
	wsOptionsTradesCmd.Flags().Bool("all", false, "Subscribe to all options trades")
	wsOptionsQuotesCmd.Flags().Bool("all", false, "Subscribe to all options quotes")
	wsOptionsAggMinuteCmd.Flags().Bool("all", false, "Subscribe to all options per-minute aggregates")
	wsOptionsAggSecondCmd.Flags().Bool("all", false, "Subscribe to all options per-second aggregates")
	wsOptionsFMVCmd.Flags().Bool("all", false, "Subscribe to all options FMV events")

	wsOptionsCmd.AddCommand(wsOptionsTradesCmd)
	wsOptionsCmd.AddCommand(wsOptionsQuotesCmd)
	wsOptionsCmd.AddCommand(wsOptionsAggMinuteCmd)
	wsOptionsCmd.AddCommand(wsOptionsAggSecondCmd)
	wsOptionsCmd.AddCommand(wsOptionsFMVCmd)

	wsCmd.AddCommand(wsOptionsCmd)
}
