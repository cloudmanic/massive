//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// indicesMarketStatusCmd retrieves the current real-time trading status
// of all US exchanges, currency markets, and major index groups. This
// shows whether markets are open, closed, in after-hours, or early-hours
// trading. The indices groups section is especially relevant for index
// tracking, covering S&P, Dow Jones, NASDAQ, MSCI, FTSE Russell, and others.
// Usage: massive indices market-status
var indicesMarketStatusCmd = &cobra.Command{
	Use:   "market-status",
	Short: "Get current market status for indices and exchanges",
	Long:  "Retrieve the real-time open/closed status of US stock exchanges (NYSE, NASDAQ, OTC), currency markets (crypto, forex), and major index groups (S&P, Dow Jones, NASDAQ, MSCI, FTSE Russell, and others).",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		result, err := client.GetIndicesMarketStatus()
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Market: %s | Server Time: %s\n", result.Market, result.ServerTime)
		fmt.Printf("After Hours: %v | Early Hours: %v\n\n", result.AfterHours, result.EarlyHours)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		fmt.Fprintln(w, "INDICES GROUPS")
		fmt.Fprintln(w, "--------------")
		fmt.Fprintf(w, "S&P\t%s\n", result.IndicesGroups.SAndP)
		fmt.Fprintf(w, "Dow Jones\t%s\n", result.IndicesGroups.DowJones)
		fmt.Fprintf(w, "NASDAQ\t%s\n", result.IndicesGroups.Nasdaq)
		fmt.Fprintf(w, "MSCI\t%s\n", result.IndicesGroups.MSCI)
		fmt.Fprintf(w, "FTSE Russell\t%s\n", result.IndicesGroups.FTSERussell)
		fmt.Fprintf(w, "Societe Generale\t%s\n", result.IndicesGroups.SocieteGenerale)
		fmt.Fprintf(w, "MStar\t%s\n", result.IndicesGroups.MStar)
		fmt.Fprintf(w, "MStarC\t%s\n", result.IndicesGroups.MStarC)
		fmt.Fprintf(w, "CCCY\t%s\n", result.IndicesGroups.CCCY)
		fmt.Fprintf(w, "CGI\t%s\n", result.IndicesGroups.CGI)
		fmt.Fprintln(w)

		fmt.Fprintln(w, "EXCHANGES")
		fmt.Fprintln(w, "--------")
		fmt.Fprintf(w, "NYSE\t%s\n", result.Exchanges.NYSE)
		fmt.Fprintf(w, "NASDAQ\t%s\n", result.Exchanges.Nasdaq)
		fmt.Fprintf(w, "OTC\t%s\n", result.Exchanges.OTC)
		fmt.Fprintln(w)

		fmt.Fprintln(w, "CURRENCIES")
		fmt.Fprintln(w, "----------")
		fmt.Fprintf(w, "Crypto\t%s\n", result.Currencies.Crypto)
		fmt.Fprintf(w, "Forex\t%s\n", result.Currencies.FX)

		w.Flush()

		return nil
	},
}

// indicesMarketHolidaysCmd retrieves the list of upcoming market holidays
// and early-close days for NYSE, NASDAQ, and OTC exchanges. This is useful
// for planning around market closures and shortened trading sessions that
// affect index calculations and trading.
// Usage: massive indices market-holidays
var indicesMarketHolidaysCmd = &cobra.Command{
	Use:   "market-holidays",
	Short: "Get upcoming market holidays",
	Long:  "Retrieve upcoming market holidays and early-close days for NYSE, NASDAQ, and OTC exchanges, including open/close times for shortened sessions that affect index trading.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		result, err := client.GetIndicesMarketHolidays()
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		if len(result) == 0 {
			fmt.Println("No upcoming market holidays found.")
			return nil
		}

		fmt.Printf("Upcoming Market Holidays: %d\n\n", len(result))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "DATE\tEXCHANGE\tNAME\tSTATUS\tOPEN\tCLOSE")
		fmt.Fprintln(w, "----\t--------\t----\t------\t----\t-----")

		for _, h := range result {
			openTime := "-"
			closeTime := "-"
			if h.Open != "" {
				openTime = h.Open
			}
			if h.Close != "" {
				closeTime = h.Close
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				h.Date, h.Exchange, h.Name, h.Status, openTime, closeTime)
		}
		w.Flush()

		return nil
	},
}

// init registers the market-status and market-holidays commands as
// subcommands of the indices parent command.
func init() {
	indicesCmd.AddCommand(indicesMarketStatusCmd)
	indicesCmd.AddCommand(indicesMarketHolidaysCmd)
}
