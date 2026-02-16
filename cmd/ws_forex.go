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

// wsForexCmd is the parent command for all forex WebSocket streaming subcommands.
// It groups quotes, aggregates, and fair market value streams under "massive ws forex".
var wsForexCmd = &cobra.Command{
	Use:   "forex",
	Short: "Stream real-time forex data via WebSocket",
}

// wsForexQuotesCmd streams real-time forex quote data for the specified currency
// pairs. Each event includes the pair name, bid price, and ask price. Supports
// both table and JSON output formats and an --all flag to subscribe to all pairs.
// Usage: massive ws forex quotes C:EURUSD C:USD/CNH
var wsForexQuotesCmd = &cobra.Command{
	Use:   "quotes [tickers...]",
	Short: "Stream real-time forex quotes",
	Long:  "Stream real-time forex quote data including bid and ask prices for specified currency pairs.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}
		tickers := args
		if all {
			tickers = []string{"*"}
		}
		params := buildForexSubscriptionParams("C", tickers)
		return connectAndStreamAsset(cmd.Context(), "forex", "C", params, formatForexQuote)
	},
}

// wsForexAggMinuteCmd streams per-minute aggregate bars for the specified forex
// currency pairs. Each event includes open, high, low, close prices and volume.
// Usage: massive ws forex agg-minute C:EURUSD
var wsForexAggMinuteCmd = &cobra.Command{
	Use:   "agg-minute [tickers...]",
	Short: "Stream per-minute forex aggregates",
	Long:  "Stream per-minute aggregate bar data (OHLCV) for specified forex currency pairs.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}
		tickers := args
		if all {
			tickers = []string{"*"}
		}
		params := buildForexSubscriptionParams("CA", tickers)
		return connectAndStreamAsset(cmd.Context(), "forex", "CA", params, formatForexAgg)
	},
}

// wsForexAggSecondCmd streams per-second aggregate bars for the specified forex
// currency pairs. Each event includes open, high, low, close prices and volume.
// Usage: massive ws forex agg-second C:EURUSD
var wsForexAggSecondCmd = &cobra.Command{
	Use:   "agg-second [tickers...]",
	Short: "Stream per-second forex aggregates",
	Long:  "Stream per-second aggregate bar data (OHLCV) for specified forex currency pairs.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}
		tickers := args
		if all {
			tickers = []string{"*"}
		}
		params := buildForexSubscriptionParams("CAS", tickers)
		return connectAndStreamAsset(cmd.Context(), "forex", "CAS", params, formatForexAgg)
	},
}

// wsForexFMVCmd streams Fair Market Value data for the specified forex currency
// pairs. Each event includes the symbol and its computed fair market value.
// Usage: massive ws forex fmv C:EURUSD
var wsForexFMVCmd = &cobra.Command{
	Use:   "fmv [tickers...]",
	Short: "Stream forex Fair Market Value",
	Long:  "Stream Fair Market Value (FMV) data for specified forex currency pairs.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}
		tickers := args
		if all {
			tickers = []string{"*"}
		}
		params := buildForexSubscriptionParams("FMV", tickers)
		return connectAndStreamAsset(cmd.Context(), "forex", "FMV", params, formatFMV)
	},
}

// buildForexSubscriptionParams constructs the subscription parameter string
// for a WebSocket subscribe message. Each ticker is prefixed with the channel
// name and a dot separator (e.g., "C.C:EURUSD"). When all is true via the tickers
// slice containing "*", it returns the wildcard pattern (e.g., "C.*").
func buildForexSubscriptionParams(channel string, tickers []string) string {
	parts := make([]string, len(tickers))
	for i, t := range tickers {
		parts[i] = channel + "." + strings.ToUpper(t)
	}
	return strings.Join(parts, ",")
}

// formatForexQuote formats a single forex quote event as a table row showing
// time, pair, bid price, and ask price. Writes the formatted row to the
// provided tabwriter.
func formatForexQuote(w *tabwriter.Writer, event map[string]interface{}) {
	ts := formatTimestamp(event["t"])
	pair := getStr(event, "p")
	bid := getFloat(event, "b")
	ask := getFloat(event, "a")
	fmt.Fprintf(w, "%s\t%s\t%.6f\t%.6f\n", ts, pair, bid, ask)
}

// formatForexAgg formats a single forex aggregate event (per-minute or per-second)
// as a table row showing time, pair, open, high, low, close, and volume.
// Writes the formatted row to the provided tabwriter.
func formatForexAgg(w *tabwriter.Writer, event map[string]interface{}) {
	ts := formatTimestamp(event["s"])
	pair := getStr(event, "pair")
	o := getFloat(event, "o")
	h := getFloat(event, "h")
	l := getFloat(event, "l")
	c := getFloat(event, "c")
	v := getFloat(event, "v")
	fmt.Fprintf(w, "%s\t%s\t%.6f\t%.6f\t%.6f\t%.6f\t%.0f\n", ts, pair, o, h, l, c, v)
}

// init registers all forex WebSocket subcommands under the wsForexCmd parent,
// which is itself registered under the wsCmd parent. Each subcommand receives
// an --all flag for subscribing to all available tickers.
func init() {
	wsForexQuotesCmd.Flags().Bool("all", false, "Subscribe to all forex pairs")
	wsForexAggMinuteCmd.Flags().Bool("all", false, "Subscribe to all forex pairs")
	wsForexAggSecondCmd.Flags().Bool("all", false, "Subscribe to all forex pairs")
	wsForexFMVCmd.Flags().Bool("all", false, "Subscribe to all forex pairs")

	wsForexCmd.AddCommand(wsForexQuotesCmd)
	wsForexCmd.AddCommand(wsForexAggMinuteCmd)
	wsForexCmd.AddCommand(wsForexAggSecondCmd)
	wsForexCmd.AddCommand(wsForexFMVCmd)

	wsCmd.AddCommand(wsForexCmd)
}
