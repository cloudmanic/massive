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

	"github.com/cloudmanic/massive-cli/internal/api"
	"github.com/spf13/cobra"
)

// stocksSnapshotsCmd is the parent command for all snapshot subcommands
// including ticker, all, gainers, and losers snapshots.
var stocksSnapshotsCmd = &cobra.Command{
	Use:   "snapshots",
	Short: "Stock market snapshot commands",
	Long:  "Retrieve real-time snapshot data for stock tickers including current day, previous day, and minute-level market data.",
}

// stocksSnapshotsTickerCmd retrieves the most recent snapshot for a single
// stock ticker. The snapshot includes the current day's bar, previous day's
// bar, latest minute bar, and the day's price change values.
// Usage: massive stocks snapshots ticker AAPL
var stocksSnapshotsTickerCmd = &cobra.Command{
	Use:   "ticker [symbol]",
	Short: "Get snapshot for a single stock ticker",
	Long:  "Retrieve the most recent snapshot for a single stock ticker including current day, previous day, minute bar, and price change data.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])

		result, err := client.GetSnapshotTicker(ticker)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		t := result.Ticker
		fmt.Printf("Ticker: %s | Change: %.4f (%.2f%%)\n\n", t.Ticker, t.TodaysChange, t.TodaysChangePct)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "PERIOD\tOPEN\tHIGH\tLOW\tCLOSE\tVOLUME\tVWAP")
		fmt.Fprintln(w, "------\t----\t----\t---\t-----\t------\t----")

		fmt.Fprintf(w, "Day\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\t%.4f\n",
			t.Day.Open, t.Day.High, t.Day.Low, t.Day.Close,
			t.Day.Volume, t.Day.VWAP)

		fmt.Fprintf(w, "Prev Day\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\t%.4f\n",
			t.PrevDay.Open, t.PrevDay.High, t.PrevDay.Low, t.PrevDay.Close,
			t.PrevDay.Volume, t.PrevDay.VWAP)

		fmt.Fprintf(w, "Minute\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\t%.4f\n",
			t.Min.Open, t.Min.High, t.Min.Low, t.Min.Close,
			t.Min.Volume, t.Min.VWAP)

		w.Flush()

		return nil
	},
}

// stocksSnapshotsAllCmd retrieves snapshot data for all US stock tickers
// or a filtered subset. Supports filtering by a comma-separated list of
// ticker symbols and optional OTC inclusion.
// Usage: massive stocks snapshots all --tickers AAPL,MSFT,TSLA
var stocksSnapshotsAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Get snapshots for all or selected stock tickers",
	Long:  "Retrieve snapshot data for all US stock tickers or a filtered subset specified by a comma-separated list of symbols.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		tickers, _ := cmd.Flags().GetString("tickers")
		includeOTC, _ := cmd.Flags().GetString("include-otc")

		params := api.AllTickersSnapshotParams{
			Tickers:    tickers,
			IncludeOTC: includeOTC,
		}

		result, err := client.GetSnapshotAllTickers(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Tickers: %d\n\n", result.Count)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tDAY OPEN\tDAY HIGH\tDAY LOW\tDAY CLOSE\tVOLUME\tCHANGE\tCHANGE %")
		fmt.Fprintln(w, "------\t--------\t--------\t-------\t---------\t------\t------\t--------")

		for _, t := range result.Tickers {
			fmt.Fprintf(w, "%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\t%.4f\t%.2f%%\n",
				t.Ticker, t.Day.Open, t.Day.High, t.Day.Low, t.Day.Close,
				t.Day.Volume, t.TodaysChange, t.TodaysChangePct)
		}
		w.Flush()

		return nil
	},
}

// stocksSnapshotsGainersCmd retrieves the current top 20 gainers in
// the US stock market. Each ticker includes the current day's bar,
// previous day's bar, and percentage change values.
// Usage: massive stocks snapshots gainers
var stocksSnapshotsGainersCmd = &cobra.Command{
	Use:   "gainers",
	Short: "Get top gaining stock tickers",
	Long:  "Retrieve the current top 20 gainers in the US stock market with snapshot data including day bar, previous day bar, and change percentages.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		includeOTC, _ := cmd.Flags().GetString("include-otc")

		params := api.GainersLosersParams{
			IncludeOTC: includeOTC,
		}

		result, err := client.GetSnapshotGainersLosers("gainers", params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		return printGainersLosersTable("Gainers", result)
	},
}

// stocksSnapshotsLosersCmd retrieves the current top 20 losers in
// the US stock market. Each ticker includes the current day's bar,
// previous day's bar, and percentage change values.
// Usage: massive stocks snapshots losers
var stocksSnapshotsLosersCmd = &cobra.Command{
	Use:   "losers",
	Short: "Get top losing stock tickers",
	Long:  "Retrieve the current top 20 losers in the US stock market with snapshot data including day bar, previous day bar, and change percentages.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		includeOTC, _ := cmd.Flags().GetString("include-otc")

		params := api.GainersLosersParams{
			IncludeOTC: includeOTC,
		}

		result, err := client.GetSnapshotGainersLosers("losers", params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		return printGainersLosersTable("Losers", result)
	},
}

// printGainersLosersTable formats and prints a table of gainers or losers
// snapshot data to stdout. The title parameter labels the output as either
// "Gainers" or "Losers" for display clarity.
func printGainersLosersTable(title string, result *api.GainersLosersSnapshotResponse) error {
	fmt.Printf("Top %s: %d tickers\n\n", title, len(result.Tickers))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TICKER\tDAY OPEN\tDAY HIGH\tDAY LOW\tDAY CLOSE\tVOLUME\tCHANGE\tCHANGE %")
	fmt.Fprintln(w, "------\t--------\t--------\t-------\t---------\t------\t------\t--------")

	for _, t := range result.Tickers {
		fmt.Fprintf(w, "%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\t%.4f\t%.2f%%\n",
			t.Ticker, t.Day.Open, t.Day.High, t.Day.Low, t.Day.Close,
			t.Day.Volume, t.TodaysChange, t.TodaysChangePct)
	}
	w.Flush()

	return nil
}

// init registers the snapshots parent command and all snapshot subcommands
// with their respective flags under the stocks parent command.
func init() {
	stocksSnapshotsAllCmd.Flags().String("tickers", "", "Comma-separated list of ticker symbols (default: all)")
	stocksSnapshotsAllCmd.Flags().String("include-otc", "false", "Include OTC securities (true/false)")

	stocksSnapshotsGainersCmd.Flags().String("include-otc", "false", "Include OTC securities (true/false)")

	stocksSnapshotsLosersCmd.Flags().String("include-otc", "false", "Include OTC securities (true/false)")

	stocksSnapshotsCmd.AddCommand(stocksSnapshotsTickerCmd)
	stocksSnapshotsCmd.AddCommand(stocksSnapshotsAllCmd)
	stocksSnapshotsCmd.AddCommand(stocksSnapshotsGainersCmd)
	stocksSnapshotsCmd.AddCommand(stocksSnapshotsLosersCmd)

	stocksCmd.AddCommand(stocksSnapshotsCmd)
}
