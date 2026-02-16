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

// wsCryptoCmd is the parent command for all crypto WebSocket streaming
// subcommands. It groups real-time crypto data streams under the "ws crypto"
// namespace including trades, quotes, aggregates, and fair market value.
var wsCryptoCmd = &cobra.Command{
	Use:   "crypto",
	Short: "Stream real-time crypto data",
}

// wsCryptoTradesCmd streams real-time crypto trade events via WebSocket.
// Each trade event includes the crypto pair, price, size, exchange, and
// timestamp. Supports subscribing to specific tickers or all tickers.
// Usage: massive ws crypto trades X:BTC-USD X:ETH-USD
var wsCryptoTradesCmd = &cobra.Command{
	Use:   "trades [tickers...]",
	Short: "Stream real-time crypto trades",
	Long:  "Stream real-time crypto trade data via WebSocket. Each event includes pair, price, size, exchange, and timestamp.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		tickers := args
		if all {
			tickers = []string{"*"}
		}

		// Build subscription parameters for XT (trades) channel with crypto tickers.
		params := buildCryptoSubscriptionParams("XT", tickers)
		return connectAndStreamAsset(cmd.Context(), "crypto", "XT", params, formatCryptoTrade)
	},
}

// wsCryptoQuotesCmd streams real-time crypto quote events via WebSocket.
// Each quote event includes the crypto pair, bid/ask prices, bid/ask sizes,
// exchange, and timestamp. Supports subscribing to specific tickers or all.
// Usage: massive ws crypto quotes X:BTC-USD
var wsCryptoQuotesCmd = &cobra.Command{
	Use:   "quotes [tickers...]",
	Short: "Stream real-time crypto quotes",
	Long:  "Stream real-time crypto quote data via WebSocket. Each event includes pair, bid/ask prices and sizes, exchange, and timestamp.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		tickers := args
		if all {
			tickers = []string{"*"}
		}

		// Build subscription parameters for XQ (quotes) channel with crypto tickers.
		params := buildCryptoSubscriptionParams("XQ", tickers)
		return connectAndStreamAsset(cmd.Context(), "crypto", "XQ", params, formatCryptoQuote)
	},
}

// wsCryptoAggMinuteCmd streams real-time per-minute aggregate bar data for
// crypto pairs via WebSocket. Each event includes open, high, low, close,
// volume, VWAP, and the start/end timestamps of the aggregate window.
// Usage: massive ws crypto agg-minute X:BTC-USD
var wsCryptoAggMinuteCmd = &cobra.Command{
	Use:   "agg-minute [tickers...]",
	Short: "Stream per-minute crypto aggregates",
	Long:  "Stream real-time per-minute aggregate bar data for crypto pairs via WebSocket. Each event includes OHLCV data and VWAP.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		tickers := args
		if all {
			tickers = []string{"*"}
		}

		// Build subscription parameters for XA (per-minute aggregate) channel with crypto tickers.
		params := buildCryptoSubscriptionParams("XA", tickers)
		return connectAndStreamAsset(cmd.Context(), "crypto", "XA", params, formatCryptoAggregate)
	},
}

// wsCryptoAggSecondCmd streams real-time per-second aggregate bar data for
// crypto pairs via WebSocket. Each event includes the same fields as the
// per-minute aggregate but at second-level granularity.
// Usage: massive ws crypto agg-second X:BTC-USD
var wsCryptoAggSecondCmd = &cobra.Command{
	Use:   "agg-second [tickers...]",
	Short: "Stream per-second crypto aggregates",
	Long:  "Stream real-time per-second aggregate bar data for crypto pairs via WebSocket. Each event includes OHLCV data and VWAP.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		tickers := args
		if all {
			tickers = []string{"*"}
		}

		// Build subscription parameters for XAS (per-second aggregate) channel with crypto tickers.
		params := buildCryptoSubscriptionParams("XAS", tickers)
		return connectAndStreamAsset(cmd.Context(), "crypto", "XAS", params, formatCryptoAggregate)
	},
}

// wsCryptoFMVCmd streams real-time Fair Market Value (FMV) data for crypto
// pairs via WebSocket. FMV represents a calculated fair price for a crypto
// asset across multiple exchanges. Each event includes symbol, FMV, and timestamp.
// Usage: massive ws crypto fmv X:BTC-USD
var wsCryptoFMVCmd = &cobra.Command{
	Use:   "fmv [tickers...]",
	Short: "Stream crypto Fair Market Value",
	Long:  "Stream real-time Fair Market Value (FMV) data for crypto pairs via WebSocket. FMV provides a calculated fair price across exchanges.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		tickers := args
		if all {
			tickers = []string{"*"}
		}

		// Build subscription parameters for FMV (Fair Market Value) channel with crypto tickers.
		params := buildCryptoSubscriptionParams("FMV", tickers)
		return connectAndStreamAsset(cmd.Context(), "crypto", "FMV", params, formatCryptoFMV)
	},
}

// buildCryptoSubscriptionParams constructs the subscription parameter string for
// the WebSocket subscribe message. It prefixes each ticker with the channel
// name (e.g., "XT.X:BTC-USD"). Multiple tickers are comma-separated.
func buildCryptoSubscriptionParams(channel string, tickers []string) string {
	parts := make([]string, len(tickers))
	for i, t := range tickers {
		parts[i] = channel + "." + t
	}
	return strings.Join(parts, ",")
}

// formatCryptoTrade formats a single crypto trade event as a table row showing
// time, pair, price, size, and exchange. Uses shared helpers formatTimestamp,
// getStr, and getFloat to extract and format values from the event map.
func formatCryptoTrade(w *tabwriter.Writer, event map[string]interface{}) {
	ts := formatTimestamp(event["t"])
	pair := getStr(event, "pair")
	price := getFloat(event, "p")
	size := getFloat(event, "s")
	exchange := getFloat(event, "x")
	fmt.Fprintf(w, "%s\t%s\t%.4f\t%.4f\t%.0f\n", ts, pair, price, size, exchange)
}

// formatCryptoQuote formats a single crypto quote event as a table row showing
// time, pair, bid price, bid size, ask price, and ask size. Uses shared helpers
// to extract and format values from the event map.
func formatCryptoQuote(w *tabwriter.Writer, event map[string]interface{}) {
	ts := formatTimestamp(event["t"])
	pair := getStr(event, "pair")
	bid := getFloat(event, "bp")
	bidSize := getFloat(event, "bs")
	ask := getFloat(event, "ap")
	askSize := getFloat(event, "as")
	fmt.Fprintf(w, "%s\t%s\t%.4f\t%.4f\t%.4f\t%.4f\n", ts, pair, bid, bidSize, ask, askSize)
}

// formatCryptoAggregate formats a single crypto aggregate bar event (per-minute or
// per-second) as a table row showing time, pair, open, high, low, close, and
// volume. Uses shared helpers to extract and format values from the event map.
func formatCryptoAggregate(w *tabwriter.Writer, event map[string]interface{}) {
	ts := formatTimestamp(event["s"])
	pair := getStr(event, "pair")
	open := getFloat(event, "o")
	high := getFloat(event, "h")
	low := getFloat(event, "l")
	close_ := getFloat(event, "c")
	volume := getFloat(event, "v")
	fmt.Fprintf(w, "%s\t%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\n", ts, pair, open, high, low, close_, volume)
}

// formatCryptoFMV formats a single crypto Fair Market Value event as a table row
// showing time, symbol, and FMV price. Uses shared helpers to extract and format
// values from the event map.
func formatCryptoFMV(w *tabwriter.Writer, event map[string]interface{}) {
	ts := formatTimestamp(event["t"])
	sym := getStr(event, "sym")
	fmv := getFloat(event, "fmv")
	fmt.Fprintf(w, "%s\t%s\t%.4f\n", ts, sym, fmv)
}

// init registers the crypto WebSocket command and all its subcommands under
// the ws parent command. Each subcommand gets an --all flag to subscribe
// to all available tickers for that channel.
func init() {
	// Add the --all flag to each subcommand for subscribing to all tickers.
	wsCryptoTradesCmd.Flags().Bool("all", false, "Subscribe to all crypto trade events")
	wsCryptoQuotesCmd.Flags().Bool("all", false, "Subscribe to all crypto quote events")
	wsCryptoAggMinuteCmd.Flags().Bool("all", false, "Subscribe to all crypto per-minute aggregates")
	wsCryptoAggSecondCmd.Flags().Bool("all", false, "Subscribe to all crypto per-second aggregates")
	wsCryptoFMVCmd.Flags().Bool("all", false, "Subscribe to all crypto FMV events")

	// Register subcommands under the crypto parent.
	wsCryptoCmd.AddCommand(wsCryptoTradesCmd)
	wsCryptoCmd.AddCommand(wsCryptoQuotesCmd)
	wsCryptoCmd.AddCommand(wsCryptoAggMinuteCmd)
	wsCryptoCmd.AddCommand(wsCryptoAggSecondCmd)
	wsCryptoCmd.AddCommand(wsCryptoFMVCmd)

	// Register the crypto command under the ws parent command.
	wsCmd.AddCommand(wsCryptoCmd)
}
