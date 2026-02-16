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

// stocksMarketStatusCmd retrieves the current real-time status of all US
// stock exchanges, currency markets, and index groups. This shows whether
// markets are open, closed, in after-hours, or early-hours trading.
// Usage: massive stocks market-status
var stocksMarketStatusCmd = &cobra.Command{
	Use:   "market-status",
	Short: "Get current market status for all exchanges",
	Long:  "Retrieve the real-time open/closed status of US stock exchanges (NYSE, NASDAQ, OTC), currency markets (crypto, forex), and major index groups.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		result, err := client.GetMarketStatus()
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Market: %s | Server Time: %s\n", result.Market, result.ServerTime)
		fmt.Printf("After Hours: %v | Early Hours: %v\n\n", result.AfterHours, result.EarlyHours)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

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
		fmt.Fprintln(w)

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

		w.Flush()

		return nil
	},
}

// stocksMarketHolidaysCmd retrieves the list of upcoming market holidays
// and early-close days for NYSE, NASDAQ, and OTC exchanges. Useful for
// planning around market closures and shortened trading sessions.
// Usage: massive stocks market-holidays
var stocksMarketHolidaysCmd = &cobra.Command{
	Use:   "market-holidays",
	Short: "Get upcoming market holidays",
	Long:  "Retrieve upcoming market holidays and early-close days for NYSE, NASDAQ, and OTC exchanges, including open/close times for shortened sessions.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		result, err := client.GetMarketHolidays()
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

// stocksExchangesCmd retrieves a list of known stock exchanges with
// their identifiers, names, and metadata. Supports optional filtering
// by asset class and locale.
// Usage: massive stocks exchanges --asset-class stocks --locale us
var stocksExchangesCmd = &cobra.Command{
	Use:   "exchanges",
	Short: "List known exchanges",
	Long:  "Retrieve a list of known exchanges including their identifiers (MIC codes), names, asset classes, and other reference attributes. Supports filtering by asset class and locale.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		assetClass, _ := cmd.Flags().GetString("asset-class")
		locale, _ := cmd.Flags().GetString("locale")

		params := api.ExchangesParams{
			AssetClass: assetClass,
			Locale:     locale,
		}

		result, err := client.GetExchanges(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Exchanges: %d\n\n", result.Count)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tACRONYM\tMIC\tTYPE\tASSET CLASS\tLOCALE")
		fmt.Fprintln(w, "--\t----\t-------\t---\t----\t-----------\t------")

		for _, e := range result.Results {
			acronym := e.Acronym
			if acronym == "" {
				acronym = "-"
			}
			mic := e.MIC
			if mic == "" {
				mic = "-"
			}
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%s\n",
				e.ID, e.Name, acronym, mic, e.Type, e.AssetClass, e.Locale)
		}
		w.Flush()

		return nil
	},
}

// init registers the market-status, market-holidays, and exchanges
// commands and their flags under the stocks parent command.
func init() {
	stocksExchangesCmd.Flags().String("asset-class", "", "Filter by asset class (stocks, options, crypto, fx)")
	stocksExchangesCmd.Flags().String("locale", "", "Filter by locale (us, global)")

	stocksCmd.AddCommand(stocksMarketStatusCmd)
	stocksCmd.AddCommand(stocksMarketHolidaysCmd)
	stocksCmd.AddCommand(stocksExchangesCmd)
}
