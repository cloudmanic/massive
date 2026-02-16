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

// optionsContractsListCmd lists and searches options contracts from the
// Massive reference data. Supports filtering by underlying ticker, contract
// type, expiration date, strike price, and various range filters.
// Usage: massive options contracts list --underlying-ticker AAPL --contract-type call
var optionsContractsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List and search options contracts",
	Long:  "Retrieve a list of options contracts with optional filtering by underlying ticker, contract type, expiration date, strike price, and range filters.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		underlyingTicker, _ := cmd.Flags().GetString("underlying-ticker")
		contractType, _ := cmd.Flags().GetString("contract-type")
		expirationDate, _ := cmd.Flags().GetString("expiration-date")
		asOf, _ := cmd.Flags().GetString("as-of")
		strikePrice, _ := cmd.Flags().GetString("strike-price")
		expired, _ := cmd.Flags().GetString("expired")
		expirationDateGte, _ := cmd.Flags().GetString("expiration-date-gte")
		expirationDateGt, _ := cmd.Flags().GetString("expiration-date-gt")
		expirationDateLte, _ := cmd.Flags().GetString("expiration-date-lte")
		expirationDateLt, _ := cmd.Flags().GetString("expiration-date-lt")
		strikePriceGte, _ := cmd.Flags().GetString("strike-price-gte")
		strikePriceGt, _ := cmd.Flags().GetString("strike-price-gt")
		strikePriceLte, _ := cmd.Flags().GetString("strike-price-lte")
		strikePriceLt, _ := cmd.Flags().GetString("strike-price-lt")
		order, _ := cmd.Flags().GetString("order")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.OptionsContractsParams{
			UnderlyingTicker:  underlyingTicker,
			ContractType:      contractType,
			ExpirationDate:    expirationDate,
			AsOf:              asOf,
			StrikePrice:       strikePrice,
			Expired:           expired,
			ExpirationDateGte: expirationDateGte,
			ExpirationDateGt:  expirationDateGt,
			ExpirationDateLte: expirationDateLte,
			ExpirationDateLt:  expirationDateLt,
			StrikePriceGte:    strikePriceGte,
			StrikePriceGt:     strikePriceGt,
			StrikePriceLte:    strikePriceLte,
			StrikePriceLt:     strikePriceLt,
			Order:             order,
			Limit:             limit,
			Sort:              sort,
		}

		result, err := client.GetOptionsContracts(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Results: %d\n\n", len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tUNDERLYING\tTYPE\tSTRIKE\tEXPIRATION\tSTYLE\tEXCHANGE")
		fmt.Fprintln(w, "------\t----------\t----\t------\t----------\t-----\t--------")

		for _, c := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%.2f\t%s\t%s\t%s\n",
				c.Ticker, c.UnderlyingTicker, c.ContractType,
				c.StrikePrice, c.ExpirationDate, c.ExerciseStyle,
				c.PrimaryExchange)
		}
		w.Flush()

		if result.NextURL != "" {
			fmt.Println("\nMore results available. Increase --limit or use pagination.")
		}

		return nil
	},
}

// optionsContractsGetCmd retrieves detailed information about a single
// options contract identified by its options ticker symbol.
// Usage: massive options contracts get O:AAPL260218C00190000
var optionsContractsGetCmd = &cobra.Command{
	Use:   "get [options_ticker]",
	Short: "Get details for a specific options contract",
	Long:  "Retrieve detailed information about a single options contract by its ticker (e.g., O:AAPL260218C00190000).",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		optionsTicker := args[0]
		asOf, _ := cmd.Flags().GetString("as-of")

		result, err := client.GetOptionsContract(optionsTicker, asOf)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		c := result.Results

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "Ticker:\t%s\n", c.Ticker)
		fmt.Fprintf(w, "Underlying:\t%s\n", c.UnderlyingTicker)
		fmt.Fprintf(w, "Type:\t%s\n", c.ContractType)
		fmt.Fprintf(w, "Exercise Style:\t%s\n", c.ExerciseStyle)
		fmt.Fprintf(w, "Strike Price:\t%.2f\n", c.StrikePrice)
		fmt.Fprintf(w, "Expiration:\t%s\n", c.ExpirationDate)
		fmt.Fprintf(w, "Shares/Contract:\t%d\n", c.SharesPerContract)
		fmt.Fprintf(w, "Primary Exchange:\t%s\n", c.PrimaryExchange)
		fmt.Fprintf(w, "CFI:\t%s\n", c.CFI)

		if c.Correction != 0 {
			fmt.Fprintf(w, "Correction:\t%d\n", c.Correction)
		}

		if len(c.AdditionalUnderlyings) > 0 {
			fmt.Fprintf(w, "\nAdditional Underlyings:\n")
			for _, au := range c.AdditionalUnderlyings {
				fmt.Fprintf(w, "  %s\t%s\t%.4f\n", au.Underlying, au.Type, au.Amount)
			}
		}

		w.Flush()

		return nil
	},
}

// optionsContractsCmd is the parent command for options contract subcommands
// including list and get. It is registered under the optionsCmd parent.
var optionsContractsCmd = &cobra.Command{
	Use:   "contracts",
	Short: "Options contract reference data commands",
}

// init registers the options contracts commands and their flags under the
// options parent command. The list command supports filtering by underlying
// ticker, contract type, expiration date, strike price, range filters,
// sort, order, and limit. The get command accepts an as-of date for
// historical snapshots.
func init() {
	// List command flags
	optionsContractsListCmd.Flags().String("underlying-ticker", "", "Filter by underlying stock ticker (e.g., AAPL)")
	optionsContractsListCmd.Flags().String("contract-type", "", "Filter by contract type (call, put)")
	optionsContractsListCmd.Flags().String("expiration-date", "", "Filter by exact expiration date (YYYY-MM-DD)")
	optionsContractsListCmd.Flags().String("as-of", "", "Historical snapshot date (YYYY-MM-DD, default: today)")
	optionsContractsListCmd.Flags().String("strike-price", "", "Filter by exact strike price")
	optionsContractsListCmd.Flags().String("expired", "", "Include expired contracts (true/false)")
	optionsContractsListCmd.Flags().String("expiration-date-gte", "", "Expiration date greater than or equal to (YYYY-MM-DD)")
	optionsContractsListCmd.Flags().String("expiration-date-gt", "", "Expiration date greater than (YYYY-MM-DD)")
	optionsContractsListCmd.Flags().String("expiration-date-lte", "", "Expiration date less than or equal to (YYYY-MM-DD)")
	optionsContractsListCmd.Flags().String("expiration-date-lt", "", "Expiration date less than (YYYY-MM-DD)")
	optionsContractsListCmd.Flags().String("strike-price-gte", "", "Strike price greater than or equal to")
	optionsContractsListCmd.Flags().String("strike-price-gt", "", "Strike price greater than")
	optionsContractsListCmd.Flags().String("strike-price-lte", "", "Strike price less than or equal to")
	optionsContractsListCmd.Flags().String("strike-price-lt", "", "Strike price less than")
	optionsContractsListCmd.Flags().String("order", "asc", "Sort order (asc/desc)")
	optionsContractsListCmd.Flags().String("limit", "20", "Number of results to return (max 1000)")
	optionsContractsListCmd.Flags().String("sort", "ticker", "Sort field (ticker, underlying_ticker, expiration_date, strike_price)")

	// Get command flags
	optionsContractsGetCmd.Flags().String("as-of", "", "Historical snapshot date (YYYY-MM-DD, default: today)")

	// Register subcommands under the contracts parent
	optionsContractsCmd.AddCommand(optionsContractsListCmd)
	optionsContractsCmd.AddCommand(optionsContractsGetCmd)

	// Register contracts under the options parent command
	optionsCmd.AddCommand(optionsContractsCmd)
}
