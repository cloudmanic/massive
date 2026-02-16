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

// optionsBarsCmd retrieves custom OHLC aggregate bars for an options contract
// ticker over a specified time range. Supports configurable timespan, multiplier,
// adjusted, sort order, and result limit.
// Usage: massive options bars O:AAPL250221C00230000 --from 2025-02-10 --to 2025-02-14
var optionsBarsCmd = &cobra.Command{
	Use:   "bars [ticker]",
	Short: "Get OHLC aggregate bars for an options contract",
	Long:  "Retrieve custom OHLC (Open, High, Low, Close) aggregate bar data for an options contract over a specified time range.",
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
		adjusted, _ := cmd.Flags().GetString("adjusted")
		sort, _ := cmd.Flags().GetString("sort")
		limit, _ := cmd.Flags().GetString("limit")

		params := api.OptionsBarsParams{
			Multiplier: multiplier,
			Timespan:   timespan,
			From:       from,
			To:         to,
			Adjusted:   adjusted,
			Sort:       sort,
			Limit:      limit,
		}

		result, err := client.GetOptionsBars(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Ticker: %s | Bars: %d | Adjusted: %v\n\n", result.Ticker, result.ResultsCount, result.Adjusted)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "DATE\tOPEN\tHIGH\tLOW\tCLOSE\tVOLUME\tVWAP\tTRADES")
		fmt.Fprintln(w, "----\t----\t----\t---\t-----\t------\t----\t------")

		for _, bar := range result.Results {
			t := time.UnixMilli(bar.Timestamp)
			fmt.Fprintf(w, "%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\t%.4f\t%d\n",
				t.Format("2006-01-02"),
				bar.Open, bar.High, bar.Low, bar.Close,
				bar.Volume, bar.VWAP, bar.NumTrades)
		}
		w.Flush()

		return nil
	},
}

// optionsDailyTickerSummaryCmd retrieves the daily open, close, high, low,
// volume, and extended hours prices for a specific options contract on a
// given date.
// Usage: massive options daily-ticker-summary O:AAPL250221C00230000 2025-02-10
var optionsDailyTickerSummaryCmd = &cobra.Command{
	Use:   "daily-ticker-summary [ticker] [date]",
	Short: "Get daily open/close data for an options contract",
	Long:  "Retrieve the opening and closing prices for a specific options contract on a given date, along with pre-market and after-hours trade prices.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		date := args[1]
		adjusted, _ := cmd.Flags().GetString("adjusted")

		result, err := client.GetOptionsDailyTickerSummary(ticker, date, adjusted)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Contract: %s | Date: %s\n\n", result.Symbol, result.From)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "FIELD\tVALUE")
		fmt.Fprintln(w, "-----\t-----")
		fmt.Fprintf(w, "Open\t%.4f\n", result.Open)
		fmt.Fprintf(w, "High\t%.4f\n", result.High)
		fmt.Fprintf(w, "Low\t%.4f\n", result.Low)
		fmt.Fprintf(w, "Close\t%.4f\n", result.Close)
		fmt.Fprintf(w, "Volume\t%.0f\n", result.Volume)
		fmt.Fprintf(w, "After Hours\t%.4f\n", result.AfterHours)
		fmt.Fprintf(w, "Pre-Market\t%.4f\n", result.PreMarket)
		w.Flush()

		return nil
	},
}

// optionsPreviousDayBarCmd retrieves the previous trading day's open, high,
// low, close, volume, VWAP, and trade count for a specified options contract
// ticker. Useful for quickly checking the most recent completed session's
// price data.
// Usage: massive options previous-day-bar O:AAPL250221C00230000
var optionsPreviousDayBarCmd = &cobra.Command{
	Use:   "previous-day-bar [ticker]",
	Short: "Get previous day OHLC data for an options contract",
	Long:  "Retrieve the previous trading day's open, high, low, close, volume, VWAP, and trade count for a specified options contract.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		adjusted, _ := cmd.Flags().GetString("adjusted")

		result, err := client.GetOptionsPreviousDayBar(ticker, adjusted)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Ticker: %s | Results: %d | Adjusted: %v\n\n", result.Ticker, result.ResultsCount, result.Adjusted)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tDATE\tOPEN\tHIGH\tLOW\tCLOSE\tVOLUME\tVWAP\tTRADES")
		fmt.Fprintln(w, "------\t----\t----\t----\t---\t-----\t------\t----\t------")

		for _, bar := range result.Results {
			t := time.UnixMilli(bar.Timestamp)
			fmt.Fprintf(w, "%s\t%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\t%.4f\t%d\n",
				bar.Ticker,
				t.Format("2006-01-02"),
				bar.Open, bar.High, bar.Low, bar.Close,
				bar.Volume, bar.VWAP, bar.NumTrades)
		}
		w.Flush()

		return nil
	},
}

// init registers all options aggregates subcommands under the options
// parent command with their respective flags.
func init() {
	optionsBarsCmd.Flags().String("multiplier", "1", "Size of the timespan multiplier")
	optionsBarsCmd.Flags().String("timespan", "day", "Timespan (minute, hour, day, week, month, quarter, year)")
	optionsBarsCmd.Flags().String("from", "", "Start date (YYYY-MM-DD) [required]")
	optionsBarsCmd.Flags().String("to", "", "End date (YYYY-MM-DD) [required]")
	optionsBarsCmd.Flags().String("adjusted", "true", "Adjust for splits (true/false)")
	optionsBarsCmd.Flags().String("sort", "asc", "Sort order (asc/desc)")
	optionsBarsCmd.Flags().String("limit", "5000", "Max number of results (max 50000)")

	optionsBarsCmd.MarkFlagRequired("from")
	optionsBarsCmd.MarkFlagRequired("to")

	optionsDailyTickerSummaryCmd.Flags().String("adjusted", "true", "Adjust for splits (true/false)")

	optionsPreviousDayBarCmd.Flags().String("adjusted", "true", "Adjust for splits (true/false)")

	optionsCmd.AddCommand(optionsBarsCmd)
	optionsCmd.AddCommand(optionsDailyTickerSummaryCmd)
	optionsCmd.AddCommand(optionsPreviousDayBarCmd)
}
