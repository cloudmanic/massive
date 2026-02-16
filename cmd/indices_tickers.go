//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/cloudmanic/massive-cli/internal/api"
	"github.com/spf13/cobra"
)

// indicesTickersCmd lists and searches index tickers from the Massive
// reference data. Supports filtering by ticker symbol, name search,
// and active status. The market parameter is automatically set to
// "indices" so only index tickers are returned.
// Usage: massive indices tickers --search "S&P"
var indicesTickersCmd = &cobra.Command{
	Use:   "tickers",
	Short: "List and search index tickers",
	Long:  "Retrieve a list of index tickers with optional filtering by symbol, name, and active status.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker, _ := cmd.Flags().GetString("ticker")
		search, _ := cmd.Flags().GetString("search")
		active, _ := cmd.Flags().GetString("active")
		sort, _ := cmd.Flags().GetString("sort")
		order, _ := cmd.Flags().GetString("order")
		limit, _ := cmd.Flags().GetString("limit")

		params := api.IndicesTickerParams{
			Ticker: ticker,
			Search: search,
			Active: active,
			Sort:   sort,
			Order:  order,
			Limit:  limit,
		}

		result, err := client.GetIndicesTickers(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Results: %d\n\n", result.Count)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tNAME\tSOURCE FEED\tACTIVE")
		fmt.Fprintln(w, "------\t----\t-----------\t------")

		for _, t := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%v\n",
				t.Ticker, t.Name, t.SourceFeed, t.Active)
		}
		w.Flush()

		return nil
	},
}

// init registers the indices tickers command and its flags under the
// indices parent command. Flags support filtering by ticker symbol,
// name search, active status, sort field, sort order, and result limit.
func init() {
	indicesTickersCmd.Flags().String("ticker", "", "Filter by specific index ticker symbol (e.g., I:SPX)")
	indicesTickersCmd.Flags().String("search", "", "Search by index name or symbol")
	indicesTickersCmd.Flags().String("active", "", "Filter by active status (true/false)")
	indicesTickersCmd.Flags().String("sort", "ticker", "Sort field (ticker, name)")
	indicesTickersCmd.Flags().String("order", "asc", "Sort order (asc/desc)")
	indicesTickersCmd.Flags().String("limit", "20", "Number of results to return (max 1000)")
	indicesCmd.AddCommand(indicesTickersCmd)
}
