//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// wsFuturesCmd is the parent command for all futures WebSocket streaming
// subcommands. It groups trades, quotes, and aggregate streams under
// "massive ws futures".
var wsFuturesCmd = &cobra.Command{
	Use:   "futures",
	Short: "Stream real-time futures data via WebSocket",
}

// wsFuturesTradesCmd streams real-time trade data for the specified futures
// contracts. Each event includes the symbol, price, and size of the trade.
// Usage: massive ws futures trades ESZ4 NQZ4
var wsFuturesTradesCmd = &cobra.Command{
	Use:   "trades [tickers...]",
	Short: "Stream real-time futures trades",
	Long:  "Stream real-time trade data for specified futures contracts including price, size, and timestamps.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}
		tickers := args
		if all {
			tickers = []string{"*"}
		}
		params := buildFuturesSubscriptionParams("T", tickers)
		return connectAndStreamAsset(cmd.Context(), "futures", "T", params, formatFuturesTrade)
	},
}

// wsFuturesQuotesCmd streams real-time quote data for the specified futures
// contracts. Each event includes bid/ask prices and sizes.
// Usage: massive ws futures quotes ESZ4
var wsFuturesQuotesCmd = &cobra.Command{
	Use:   "quotes [tickers...]",
	Short: "Stream real-time futures quotes",
	Long:  "Stream real-time quote data for specified futures contracts including bid/ask prices and sizes.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}
		tickers := args
		if all {
			tickers = []string{"*"}
		}
		params := buildFuturesSubscriptionParams("Q", tickers)
		return connectAndStreamAsset(cmd.Context(), "futures", "Q", params, formatFuturesQuote)
	},
}

// wsFuturesAggMinuteCmd streams per-minute aggregate bars for the specified
// futures contracts. Each event includes OHLCV data for the completed bar.
// Usage: massive ws futures agg-minute ESZ4
var wsFuturesAggMinuteCmd = &cobra.Command{
	Use:   "agg-minute [tickers...]",
	Short: "Stream per-minute futures aggregates",
	Long:  "Stream per-minute aggregate bar data (OHLCV) for specified futures contracts.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}
		tickers := args
		if all {
			tickers = []string{"*"}
		}
		params := buildFuturesSubscriptionParams("AM", tickers)
		return connectAndStreamAsset(cmd.Context(), "futures", "AM", params, formatAggregate)
	},
}

// wsFuturesAggSecondCmd streams per-second aggregate bars for the specified
// futures contracts. Each event includes OHLCV data for the completed bar.
// Usage: massive ws futures agg-second ESZ4
var wsFuturesAggSecondCmd = &cobra.Command{
	Use:   "agg-second [tickers...]",
	Short: "Stream per-second futures aggregates",
	Long:  "Stream per-second aggregate bar data (OHLCV) for specified futures contracts.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}
		tickers := args
		if all {
			tickers = []string{"*"}
		}
		params := buildFuturesSubscriptionParams("A", tickers)
		return connectAndStreamAsset(cmd.Context(), "futures", "A", params, formatAggregate)
	},
}

// buildFuturesSubscriptionParams constructs the subscription parameter string
// for a WebSocket subscribe message. Each ticker is prefixed with the channel
// name and a dot separator (e.g., "T.ESZ4").
func buildFuturesSubscriptionParams(channel string, tickers []string) string {
	parts := make([]string, len(tickers))
	for i, t := range tickers {
		parts[i] = channel + "." + strings.ToUpper(t)
	}
	return strings.Join(parts, ",")
}

// formatFuturesTrade formats a single futures trade event as a table row
// showing time, symbol, price, and size. Writes the formatted row to the
// provided tabwriter.
func formatFuturesTrade(w *tabwriter.Writer, event map[string]interface{}) {
	ts := formatTimestamp(event["t"])
	sym := getStr(event, "sym")
	price := getFloat(event, "p")
	size := getFloat(event, "s")
	fmt.Fprintf(w, "%s\t%s\t%.4f\t%.0f\n", ts, sym, price, size)
}

// formatFuturesQuote formats a single futures quote event as a table row
// showing time, symbol, bid price, bid size, ask price, and ask size. Writes
// the formatted row to the provided tabwriter.
func formatFuturesQuote(w *tabwriter.Writer, event map[string]interface{}) {
	ts := formatTimestamp(event["t"])
	sym := getStr(event, "sym")
	bid := getFloat(event, "bp")
	bidSize := getFloat(event, "bs")
	ask := getFloat(event, "ap")
	askSize := getFloat(event, "as")
	fmt.Fprintf(w, "%s\t%s\t%.4f\t%.0f\t%.4f\t%.0f\n", ts, sym, bid, bidSize, ask, askSize)
}

// init registers the futures WebSocket command and all its subcommands under
// the ws parent command.
func init() {
	wsFuturesTradesCmd.Flags().Bool("all", false, "Subscribe to all futures trades")
	wsFuturesQuotesCmd.Flags().Bool("all", false, "Subscribe to all futures quotes")
	wsFuturesAggMinuteCmd.Flags().Bool("all", false, "Subscribe to all futures per-minute aggregates")
	wsFuturesAggSecondCmd.Flags().Bool("all", false, "Subscribe to all futures per-second aggregates")

	wsFuturesCmd.AddCommand(wsFuturesTradesCmd)
	wsFuturesCmd.AddCommand(wsFuturesQuotesCmd)
	wsFuturesCmd.AddCommand(wsFuturesAggMinuteCmd)
	wsFuturesCmd.AddCommand(wsFuturesAggSecondCmd)

	wsCmd.AddCommand(wsFuturesCmd)
}
