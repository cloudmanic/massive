//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/cloudmanic/massive-cli/internal/api"
	"github.com/spf13/cobra"
)

// indicesBarsCmd retrieves custom OHLC aggregate bars for an index ticker
// over a specified time range. Supports configurable timespan, multiplier,
// sort order, and result limit.
// Usage: massive indices bars I:SPX --from 2025-01-06 --to 2025-01-08
var indicesBarsCmd = &cobra.Command{
	Use:   "bars [ticker]",
	Short: "Get OHLC aggregate bars for an index ticker",
	Long:  "Retrieve custom OHLC (Open, High, Low, Close) aggregate bar data for an index ticker over a specified time range.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		multiplier, _ := cmd.Flags().GetString("multiplier")
		timespan, _ := cmd.Flags().GetString("timespan")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		sort, _ := cmd.Flags().GetString("sort")
		limit, _ := cmd.Flags().GetString("limit")

		params := api.IndicesBarsParams{
			Multiplier: multiplier,
			Timespan:   timespan,
			From:       from,
			To:         to,
			Sort:       sort,
			Limit:      limit,
		}

		result, err := client.GetIndicesBars(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Ticker: %s | Bars: %d\n\n", result.Ticker, result.ResultsCount)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "DATE\tOPEN\tHIGH\tLOW\tCLOSE")
		fmt.Fprintln(w, "----\t----\t----\t---\t-----")

		for _, bar := range result.Results {
			t := time.UnixMilli(bar.Timestamp)
			fmt.Fprintf(w, "%s\t%.4f\t%.4f\t%.4f\t%.4f\n",
				t.Format("2006-01-02"),
				bar.Open, bar.High, bar.Low, bar.Close)
		}
		w.Flush()

		return nil
	},
}

// indicesDailyTickerSummaryCmd retrieves the daily open, close, high, low,
// and extended hours prices for a specific index ticker on a given date.
// Usage: massive indices daily-ticker-summary I:SPX 2025-01-06
var indicesDailyTickerSummaryCmd = &cobra.Command{
	Use:   "daily-ticker-summary [ticker] [date]",
	Short: "Get daily open/close data for an index ticker",
	Long:  "Retrieve the opening and closing prices for a specific index on a given date, along with pre-market and after-hours prices.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		date := args[1]

		result, err := client.GetIndicesDailyTickerSummary(ticker, date)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Index: %s | Date: %s\n\n", result.Symbol, result.From)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "FIELD\tVALUE")
		fmt.Fprintln(w, "-----\t-----")
		fmt.Fprintf(w, "Open\t%.4f\n", result.Open)
		fmt.Fprintf(w, "High\t%.4f\n", result.High)
		fmt.Fprintf(w, "Low\t%.4f\n", result.Low)
		fmt.Fprintf(w, "Close\t%.4f\n", result.Close)
		fmt.Fprintf(w, "After Hours\t%.4f\n", result.AfterHours)
		fmt.Fprintf(w, "Pre-Market\t%.4f\n", result.PreMarket)
		w.Flush()

		return nil
	},
}

// indicesPreviousDayBarCmd retrieves the previous trading day's open, high,
// low, and close data for a specified index ticker. Useful for quickly
// checking the most recent completed session's price data.
// Usage: massive indices previous-day-bar I:SPX
var indicesPreviousDayBarCmd = &cobra.Command{
	Use:   "previous-day-bar [ticker]",
	Short: "Get previous day OHLC data for an index ticker",
	Long:  "Retrieve the previous trading day's open, high, low, and close (OHLC) data for a specified index ticker.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])

		result, err := client.GetIndicesPreviousDayBar(ticker)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Ticker: %s | Results: %d\n\n", result.Ticker, result.ResultsCount)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tDATE\tOPEN\tHIGH\tLOW\tCLOSE")
		fmt.Fprintln(w, "------\t----\t----\t----\t---\t-----")

		for _, bar := range result.Results {
			t := time.UnixMilli(bar.Timestamp)
			fmt.Fprintf(w, "%s\t%s\t%.4f\t%.4f\t%.4f\t%.4f\n",
				bar.Ticker,
				t.Format("2006-01-02"),
				bar.Open, bar.High, bar.Low, bar.Close)
		}
		w.Flush()

		return nil
	},
}

// init registers all indices aggregates subcommands under the indices
// parent command with their respective flags.
func init() {
	indicesBarsCmd.Flags().String("multiplier", "1", "Size of the timespan multiplier")
	indicesBarsCmd.Flags().String("timespan", "day", "Timespan (minute, hour, day, week, month, quarter, year)")
	indicesBarsCmd.Flags().String("from", "", "Start date (YYYY-MM-DD) [required]")
	indicesBarsCmd.Flags().String("to", "", "End date (YYYY-MM-DD) [required]")
	indicesBarsCmd.Flags().String("sort", "asc", "Sort order (asc/desc)")
	indicesBarsCmd.Flags().String("limit", "5000", "Max number of results (max 50000)")

	indicesBarsCmd.MarkFlagRequired("from")
	indicesBarsCmd.MarkFlagRequired("to")

	indicesCmd.AddCommand(indicesBarsCmd)
	indicesCmd.AddCommand(indicesDailyTickerSummaryCmd)
	indicesCmd.AddCommand(indicesPreviousDayBarCmd)
}
