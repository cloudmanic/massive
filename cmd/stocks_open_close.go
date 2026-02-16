//
// Date: 2026-02-14
// Copyright (c) 2026. All rights reserved.
//

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// stocksOpenCloseCmd retrieves the daily open, close, high, low, volume,
// pre-market, and after-hours prices for a specific stock ticker on a
// given date. Usage: massive stocks open-close AAPL 2024-01-09
var stocksOpenCloseCmd = &cobra.Command{
	Use:   "open-close [ticker] [date]",
	Short: "Get daily open/close data for a stock ticker",
	Long:  "Retrieve the opening and closing prices for a specific stock ticker on a given date, along with pre-market and after-hours trade prices.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		adjusted, _ := cmd.Flags().GetString("adjusted")
		ticker := strings.ToUpper(args[0])
		date := args[1]

		result, err := client.GetOpenClose(ticker, date, adjusted)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Symbol:      %s\n", result.Symbol)
		fmt.Printf("Date:        %s\n", result.From)
		fmt.Printf("Open:        $%.4f\n", result.Open)
		fmt.Printf("High:        $%.4f\n", result.High)
		fmt.Printf("Low:         $%.4f\n", result.Low)
		fmt.Printf("Close:       $%.4f\n", result.Close)
		fmt.Printf("Volume:      %d\n", result.Volume)
		fmt.Printf("Pre-Market:  $%.4f\n", result.PreMarket)
		fmt.Printf("After Hours: $%.4f\n", result.AfterHours)

		return nil
	},
}

// init registers the open-close command and its flags under the stocks parent command.
func init() {
	stocksOpenCloseCmd.Flags().String("adjusted", "true", "Adjust for splits (true/false)")
	stocksCmd.AddCommand(stocksOpenCloseCmd)
}
