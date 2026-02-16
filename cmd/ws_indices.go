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

// wsIndicesCmd is the parent command for all WebSocket-based indices streaming
// subcommands. It groups real-time index data streams such as per-minute
// aggregates, per-second aggregates, and live index values.
var wsIndicesCmd = &cobra.Command{
	Use:   "indices",
	Short: "Stream real-time indices data",
}

// wsIndicesAggMinuteCmd streams per-minute aggregate bars for one or more
// index tickers via WebSocket. Each message contains the symbol, open, close,
// high, low, and start/end timestamps for the aggregate window.
// Usage: massive ws indices agg-minute I:SPX I:DJI
var wsIndicesAggMinuteCmd = &cobra.Command{
	Use:   "agg-minute [tickers...]",
	Short: "Stream per-minute aggregate bars for indices",
	Long:  "Connect to the Massive WebSocket API and stream real-time per-minute OHLC aggregate data for one or more index tickers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}
		tickers := args
		if all {
			tickers = []string{"*"}
		}
		params := buildIndicesSubscriptionParams("AM", tickers)
		return connectAndStreamAsset(cmd.Context(), "indices", "AM", params, formatAggregate)
	},
}

// wsIndicesAggSecondCmd streams per-second aggregate bars for one or more
// index tickers via WebSocket. Each message contains the symbol, open, close,
// high, low, and start/end timestamps for the aggregate window.
// Usage: massive ws indices agg-second I:SPX
var wsIndicesAggSecondCmd = &cobra.Command{
	Use:   "agg-second [tickers...]",
	Short: "Stream per-second aggregate bars for indices",
	Long:  "Connect to the Massive WebSocket API and stream real-time per-second OHLC aggregate data for one or more index tickers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}
		tickers := args
		if all {
			tickers = []string{"*"}
		}
		params := buildIndicesSubscriptionParams("A", tickers)
		return connectAndStreamAsset(cmd.Context(), "indices", "A", params, formatAggregate)
	},
}

// wsIndicesValueCmd streams real-time index values for one or more index
// tickers via WebSocket. Each message contains the ticker symbol, the
// current index value, and a timestamp.
// Usage: massive ws indices value I:SPX I:COMP
var wsIndicesValueCmd = &cobra.Command{
	Use:   "value [tickers...]",
	Short: "Stream real-time index values",
	Long:  "Connect to the Massive WebSocket API and stream real-time index value updates for one or more index tickers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}
		tickers := args
		if all {
			tickers = []string{"*"}
		}
		params := buildIndicesSubscriptionParams("V", tickers)
		return connectAndStreamAsset(cmd.Context(), "indices", "V", params, formatIndicesValue)
	},
}

// buildIndicesSubscriptionParams constructs the subscription parameter string
// for a WebSocket subscribe message. Each ticker is prefixed with the channel
// name and a dot separator (e.g., "AM.I:SPX").
func buildIndicesSubscriptionParams(channel string, tickers []string) string {
	parts := make([]string, len(tickers))
	for i, t := range tickers {
		if t == "*" {
			parts[i] = channel + ".*"
		} else {
			parts[i] = channel + "." + strings.ToUpper(t)
		}
	}
	return strings.Join(parts, ",")
}

// formatIndicesValue formats a single index value event (V) for display.
// Uses the shared helper functions from ws_stocks.go to extract and format
// the event data as a table row with timestamp, ticker symbol, and value.
func formatIndicesValue(w *tabwriter.Writer, event map[string]interface{}) {
	ts := formatTimestamp(event["t"])
	sym := getStr(event, "T")
	val := getFloat(event, "val")
	fmt.Fprintf(w, "%s\t%s\t%.4f\n", ts, sym, val)
}

// init registers the indices WebSocket streaming subcommands under the
// wsIndicesCmd parent, adds the --all flag to each subcommand, and registers
// wsIndicesCmd under the shared wsCmd parent command.
func init() {
	wsIndicesAggMinuteCmd.Flags().Bool("all", false, "Subscribe to all index tickers")
	wsIndicesAggSecondCmd.Flags().Bool("all", false, "Subscribe to all index tickers")
	wsIndicesValueCmd.Flags().Bool("all", false, "Subscribe to all index tickers")

	wsIndicesCmd.AddCommand(wsIndicesAggMinuteCmd)
	wsIndicesCmd.AddCommand(wsIndicesAggSecondCmd)
	wsIndicesCmd.AddCommand(wsIndicesValueCmd)

	wsCmd.AddCommand(wsIndicesCmd)
}
