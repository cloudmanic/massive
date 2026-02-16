//
// Date: 2026-02-14
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

// stocksTickersCmd lists and searches stock tickers from the Massive
// reference data. Supports filtering by ticker symbol, name search,
// type, market, exchange, and active status.
// Usage: massive stocks tickers --search "Apple"
var stocksTickersCmd = &cobra.Command{
	Use:   "tickers",
	Short: "List and search stock tickers",
	Long:  "Retrieve a list of stock tickers with optional filtering by symbol, name, type, market, exchange, and active status.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker, _ := cmd.Flags().GetString("ticker")
		search, _ := cmd.Flags().GetString("search")
		tickerType, _ := cmd.Flags().GetString("type")
		market, _ := cmd.Flags().GetString("market")
		exchange, _ := cmd.Flags().GetString("exchange")
		active, _ := cmd.Flags().GetString("active")
		sort, _ := cmd.Flags().GetString("sort")
		order, _ := cmd.Flags().GetString("order")
		limit, _ := cmd.Flags().GetString("limit")

		params := api.TickerParams{
			Ticker:   ticker,
			Search:   search,
			Type:     tickerType,
			Market:   market,
			Exchange: exchange,
			Active:   active,
			Sort:     sort,
			Order:    order,
			Limit:    limit,
		}

		result, err := client.GetTickers(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Results: %d\n\n", result.Count)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tNAME\tTYPE\tEXCHANGE\tACTIVE")
		fmt.Fprintln(w, "------\t----\t----\t--------\t------")

		for _, t := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\n",
				t.Ticker, t.Name, t.Type, t.PrimaryExchange, t.Active)
		}
		w.Flush()

		return nil
	},
}

// init registers the tickers command and its flags under the stocks parent command.
func init() {
	stocksTickersCmd.Flags().String("ticker", "", "Filter by specific ticker symbol")
	stocksTickersCmd.Flags().String("search", "", "Search by company name or symbol")
	stocksTickersCmd.Flags().String("type", "", "Filter by ticker type (CS, ETF, etc.)")
	stocksTickersCmd.Flags().String("market", "", "Filter by market (stocks, crypto, fx)")
	stocksTickersCmd.Flags().String("exchange", "", "Filter by primary exchange")
	stocksTickersCmd.Flags().String("active", "", "Filter by active status (true/false)")
	stocksTickersCmd.Flags().String("sort", "ticker", "Sort field (ticker, name, market, type)")
	stocksTickersCmd.Flags().String("order", "asc", "Sort order (asc/desc)")
	stocksTickersCmd.Flags().String("limit", "20", "Number of results to return (max 1000)")
	stocksCmd.AddCommand(stocksTickersCmd)
}
