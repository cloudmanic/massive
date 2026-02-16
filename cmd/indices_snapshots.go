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

// indicesSnapshotsCmd is the parent command for all indices snapshot
// subcommands including ticker and all-tickers snapshots.
var indicesSnapshotsCmd = &cobra.Command{
	Use:   "snapshots",
	Short: "Index market snapshot commands",
	Long:  "Retrieve real-time snapshot data for market indices including current value, session metrics, and market status.",
}

// indicesSnapshotsTickerCmd retrieves the most recent snapshot for a
// single index ticker. The snapshot includes the current value, trading
// session data (open, high, low, close), change values, market status,
// and last update timestamp.
// Usage: massive indices snapshots ticker I:SPX
var indicesSnapshotsTickerCmd = &cobra.Command{
	Use:   "ticker [symbol]",
	Short: "Get snapshot for a single index ticker",
	Long:  "Retrieve the most recent snapshot for a single index ticker including current value, session data, change values, and market status.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])

		params := api.IndicesSnapshotParams{
			TickerAnyOf: ticker,
		}

		result, err := client.GetIndicesSnapshot(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		if len(result.Results) == 0 {
			fmt.Println("No snapshot data found for ticker:", ticker)
			return nil
		}

		idx := result.Results[0]
		fmt.Printf("Index: %s (%s)\n", idx.Ticker, idx.Name)
		fmt.Printf("Value: %.2f | Change: %.2f (%.4f%%)\n", idx.Value, idx.Session.Change, idx.Session.ChangePercent)
		fmt.Printf("Market Status: %s | Timeframe: %s\n\n", idx.MarketStatus, idx.Timeframe)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "OPEN\tHIGH\tLOW\tCLOSE\tPREV CLOSE")
		fmt.Fprintln(w, "----\t----\t---\t-----\t----------")
		fmt.Fprintf(w, "%.2f\t%.2f\t%.2f\t%.2f\t%.2f\n",
			idx.Session.Open, idx.Session.High, idx.Session.Low,
			idx.Session.Close, idx.Session.PreviousClose)
		w.Flush()

		return nil
	},
}

// indicesSnapshotsAllCmd retrieves snapshot data for all indices or a
// filtered subset. Supports filtering by a comma-separated list of
// ticker symbols and optional limit, order, and sort parameters.
// Usage: massive indices snapshots all --tickers I:SPX,I:DJI,I:COMP
var indicesSnapshotsAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Get snapshots for all or selected index tickers",
	Long:  "Retrieve snapshot data for all indices or a filtered subset specified by a comma-separated list of ticker symbols.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		tickers, _ := cmd.Flags().GetString("tickers")
		limit, _ := cmd.Flags().GetString("limit")
		order, _ := cmd.Flags().GetString("order")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.IndicesSnapshotParams{
			TickerAnyOf: tickers,
			Limit:       limit,
			Order:       order,
			Sort:        sort,
		}

		result, err := client.GetIndicesSnapshot(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Indices: %d\n\n", len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tNAME\tVALUE\tOPEN\tHIGH\tLOW\tCLOSE\tCHANGE\tCHANGE %\tSTATUS")
		fmt.Fprintln(w, "------\t----\t-----\t----\t----\t---\t-----\t------\t--------\t------")

		for _, idx := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.4f%%\t%s\n",
				idx.Ticker, idx.Name, idx.Value,
				idx.Session.Open, idx.Session.High, idx.Session.Low,
				idx.Session.Close, idx.Session.Change, idx.Session.ChangePercent,
				idx.MarketStatus)
		}
		w.Flush()

		return nil
	},
}

// init registers the snapshots parent command and all snapshot subcommands
// with their respective flags under the indices parent command.
func init() {
	indicesSnapshotsAllCmd.Flags().String("tickers", "", "Comma-separated list of index ticker symbols (e.g. I:SPX,I:DJI)")
	indicesSnapshotsAllCmd.Flags().String("limit", "", "Maximum number of results (default: 10, max: 250)")
	indicesSnapshotsAllCmd.Flags().String("order", "", "Order results (asc or desc)")
	indicesSnapshotsAllCmd.Flags().String("sort", "", "Field to sort results by")

	indicesSnapshotsCmd.AddCommand(indicesSnapshotsTickerCmd)
	indicesSnapshotsCmd.AddCommand(indicesSnapshotsAllCmd)

	indicesCmd.AddCommand(indicesSnapshotsCmd)
}
