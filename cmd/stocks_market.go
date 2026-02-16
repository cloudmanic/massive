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

// stocksMarketCmd retrieves the grouped daily OHLC summary for all US
// stocks on a specified date. Useful for broad market analysis and
// screening. Usage: massive stocks market 2024-01-09
var stocksMarketCmd = &cobra.Command{
	Use:   "market [date]",
	Short: "Get daily market summary for all stocks",
	Long:  "Retrieve the daily OHLC, volume, and VWAP data for all US stocks on a specified trading date.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		date := args[0]
		adjusted, _ := cmd.Flags().GetString("adjusted")
		includeOTC, _ := cmd.Flags().GetString("include-otc")

		params := api.MarketSummaryParams{
			Adjusted:   adjusted,
			IncludeOTC: includeOTC,
		}

		result, err := client.GetMarketSummary(date, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Date: %s | Tickers: %d | Adjusted: %v\n\n", date, result.ResultsCount, result.Adjusted)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tOPEN\tHIGH\tLOW\tCLOSE\tVOLUME\tVWAP\tTRADES")
		fmt.Fprintln(w, "------\t----\t----\t---\t-----\t------\t----\t------")

		for _, s := range result.Results {
			fmt.Fprintf(w, "%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\t%.4f\t%d\n",
				s.Ticker, s.Open, s.High, s.Low, s.Close,
				s.Volume, s.VWAP, s.NumTrades)
		}
		w.Flush()

		return nil
	},
}

// init registers the market command and its flags under the stocks parent command.
func init() {
	stocksMarketCmd.Flags().String("adjusted", "true", "Adjust for splits (true/false)")
	stocksMarketCmd.Flags().String("include-otc", "false", "Include OTC securities (true/false)")
	stocksCmd.AddCommand(stocksMarketCmd)
}
