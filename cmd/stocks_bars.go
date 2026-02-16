//
// Date: 2026-02-14
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

// stocksBarsCmd retrieves custom OHLC aggregate bars for a stock ticker
// over a specified time range. Supports configurable timespan, multiplier,
// sort order, and result limit. Usage: massive stocks bars AAPL --from 2024-01-01 --to 2024-01-31
var stocksBarsCmd = &cobra.Command{
	Use:   "bars [ticker]",
	Short: "Get OHLC aggregate bars for a stock ticker",
	Long:  "Retrieve custom OHLC (Open, High, Low, Close) aggregate bar data for a stock ticker over a specified time range.",
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

		params := api.BarsParams{
			Multiplier: multiplier,
			Timespan:   timespan,
			From:       from,
			To:         to,
			Adjusted:   adjusted,
			Sort:       sort,
			Limit:      limit,
		}

		result, err := client.GetBars(ticker, params)
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

// init registers the bars command and its flags under the stocks parent command.
func init() {
	stocksBarsCmd.Flags().String("multiplier", "1", "Size of the timespan multiplier")
	stocksBarsCmd.Flags().String("timespan", "day", "Timespan (minute, hour, day, week, month, quarter, year)")
	stocksBarsCmd.Flags().String("from", "", "Start date (YYYY-MM-DD) [required]")
	stocksBarsCmd.Flags().String("to", "", "End date (YYYY-MM-DD) [required]")
	stocksBarsCmd.Flags().String("adjusted", "true", "Adjust for splits (true/false)")
	stocksBarsCmd.Flags().String("sort", "asc", "Sort order (asc/desc)")
	stocksBarsCmd.Flags().String("limit", "5000", "Max number of results (max 50000)")

	stocksBarsCmd.MarkFlagRequired("from")
	stocksBarsCmd.MarkFlagRequired("to")

	stocksCmd.AddCommand(stocksBarsCmd)
}
