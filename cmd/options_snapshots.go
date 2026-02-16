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

// optionsSnapshotsCmd is the parent command for all options snapshot
// subcommands including chain and contract snapshots.
var optionsSnapshotsCmd = &cobra.Command{
	Use:   "snapshots",
	Short: "Options market snapshot commands",
	Long:  "Retrieve real-time snapshot data for options contracts including day bar, Greeks, implied volatility, quotes, trades, and open interest.",
}

// optionsSnapshotsChainCmd retrieves snapshot data for all options contracts
// associated with a given underlying asset. The response includes day bar,
// contract details, Greeks, implied volatility, last quote, last trade,
// open interest, and underlying asset data for each contract. Supports
// filtering by strike price, expiration date, contract type, and pagination.
// Usage: massive options snapshots chain AAPL --strike-price 250 --expiration-date 2026-03-20
var optionsSnapshotsChainCmd = &cobra.Command{
	Use:   "chain [underlying]",
	Short: "Get options chain snapshot for an underlying asset",
	Long:  "Retrieve snapshot data for all options contracts associated with a given underlying asset ticker, with optional filters for strike price, expiration date, and contract type.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		underlying := strings.ToUpper(args[0])

		strikePrice, _ := cmd.Flags().GetString("strike-price")
		expirationDate, _ := cmd.Flags().GetString("expiration-date")
		contractType, _ := cmd.Flags().GetString("contract-type")
		strikePriceGTE, _ := cmd.Flags().GetString("strike-price-gte")
		strikePriceGT, _ := cmd.Flags().GetString("strike-price-gt")
		strikePriceLTE, _ := cmd.Flags().GetString("strike-price-lte")
		strikePriceLT, _ := cmd.Flags().GetString("strike-price-lt")
		expirationDateGTE, _ := cmd.Flags().GetString("expiration-date-gte")
		expirationDateGT, _ := cmd.Flags().GetString("expiration-date-gt")
		expirationDateLTE, _ := cmd.Flags().GetString("expiration-date-lte")
		expirationDateLT, _ := cmd.Flags().GetString("expiration-date-lt")
		order, _ := cmd.Flags().GetString("order")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.OptionsChainSnapshotParams{
			StrikePrice:       strikePrice,
			ExpirationDate:    expirationDate,
			ContractType:      contractType,
			StrikePriceGTE:    strikePriceGTE,
			StrikePriceGT:     strikePriceGT,
			StrikePriceLTE:    strikePriceLTE,
			StrikePriceLT:     strikePriceLT,
			ExpirationDateGTE: expirationDateGTE,
			ExpirationDateGT:  expirationDateGT,
			ExpirationDateLTE: expirationDateLTE,
			ExpirationDateLT:  expirationDateLT,
			Order:             order,
			Limit:             limit,
			Sort:              sort,
		}

		result, err := client.GetOptionsChainSnapshot(underlying, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		if len(result.Results) == 0 {
			fmt.Println("No options contracts found for:", underlying)
			return nil
		}

		fmt.Printf("Options Chain: %s (%d contracts)\n\n", underlying, len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "CONTRACT\tTYPE\tSTRIKE\tEXPIRATION\tCLOSE\tVOLUME\tOI\tIV\tDELTA\tGAMMA\tTHETA\tVEGA")
		fmt.Fprintln(w, "--------\t----\t------\t----------\t-----\t------\t--\t--\t-----\t-----\t-----\t----")

		for _, r := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t%.2f\t%s\t%.2f\t%.0f\t%.0f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\n",
				r.Details.Ticker, r.Details.ContractType, r.Details.StrikePrice,
				r.Details.ExpirationDate, r.Day.Close, r.Day.Volume,
				r.OpenInterest, r.ImpliedVolatility,
				r.Greeks.Delta, r.Greeks.Gamma, r.Greeks.Theta, r.Greeks.Vega)
		}
		w.Flush()

		return nil
	},
}

// optionsSnapshotsContractCmd retrieves the most recent snapshot for a
// single option contract identified by the underlying asset ticker and
// the option contract ticker. The snapshot includes the day bar, contract
// details, Greeks, implied volatility, last quote, last trade, open
// interest, and underlying asset information.
// Usage: massive options snapshots contract AAPL O:AAPL260320C00250000
var optionsSnapshotsContractCmd = &cobra.Command{
	Use:   "contract [underlying] [optionTicker]",
	Short: "Get snapshot for a single option contract",
	Long:  "Retrieve the most recent snapshot for a single option contract including day bar, contract details, Greeks, implied volatility, quotes, trades, and open interest.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		underlying := strings.ToUpper(args[0])
		optionTicker := strings.ToUpper(args[1])

		result, err := client.GetOptionContractSnapshot(underlying, optionTicker)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		r := result.Results

		fmt.Printf("Contract: %s (%s %s)\n", r.Details.Ticker, r.Details.ContractType, r.Details.ExerciseStyle)
		fmt.Printf("Underlying: %s | Strike: %.2f | Expiration: %s\n", r.UnderlyingAsset.Ticker, r.Details.StrikePrice, r.Details.ExpirationDate)
		fmt.Printf("Break Even: %.2f | IV: %.4f | Open Interest: %.0f\n\n", r.BreakEvenPrice, r.ImpliedVolatility, r.OpenInterest)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		fmt.Fprintln(w, "-- Day Bar --")
		fmt.Fprintln(w, "OPEN\tHIGH\tLOW\tCLOSE\tVOLUME\tVWAP\tCHANGE\tCHANGE %")
		fmt.Fprintln(w, "----\t----\t---\t-----\t------\t----\t------\t--------")
		fmt.Fprintf(w, "%.4f\t%.4f\t%.4f\t%.4f\t%.0f\t%.4f\t%.4f\t%.2f%%\n\n",
			r.Day.Open, r.Day.High, r.Day.Low, r.Day.Close,
			r.Day.Volume, r.Day.VWAP, r.Day.Change, r.Day.ChangePercent)

		fmt.Fprintln(w, "-- Greeks --")
		fmt.Fprintln(w, "DELTA\tGAMMA\tTHETA\tVEGA")
		fmt.Fprintln(w, "-----\t-----\t-----\t----")
		fmt.Fprintf(w, "%.4f\t%.4f\t%.4f\t%.4f\n",
			r.Greeks.Delta, r.Greeks.Gamma, r.Greeks.Theta, r.Greeks.Vega)

		w.Flush()

		return nil
	},
}

// init registers the options snapshots parent command and all snapshot
// subcommands with their respective flags under the options parent command.
func init() {
	optionsSnapshotsChainCmd.Flags().String("strike-price", "", "Filter by exact strike price")
	optionsSnapshotsChainCmd.Flags().String("expiration-date", "", "Filter by exact expiration date (YYYY-MM-DD)")
	optionsSnapshotsChainCmd.Flags().String("contract-type", "", "Filter by contract type (call or put)")
	optionsSnapshotsChainCmd.Flags().String("strike-price-gte", "", "Strike price greater than or equal to")
	optionsSnapshotsChainCmd.Flags().String("strike-price-gt", "", "Strike price greater than")
	optionsSnapshotsChainCmd.Flags().String("strike-price-lte", "", "Strike price less than or equal to")
	optionsSnapshotsChainCmd.Flags().String("strike-price-lt", "", "Strike price less than")
	optionsSnapshotsChainCmd.Flags().String("expiration-date-gte", "", "Expiration date greater than or equal to (YYYY-MM-DD)")
	optionsSnapshotsChainCmd.Flags().String("expiration-date-gt", "", "Expiration date greater than (YYYY-MM-DD)")
	optionsSnapshotsChainCmd.Flags().String("expiration-date-lte", "", "Expiration date less than or equal to (YYYY-MM-DD)")
	optionsSnapshotsChainCmd.Flags().String("expiration-date-lt", "", "Expiration date less than (YYYY-MM-DD)")
	optionsSnapshotsChainCmd.Flags().String("order", "", "Sort direction for results (asc or desc)")
	optionsSnapshotsChainCmd.Flags().String("limit", "", "Maximum number of results (default: 10, max: 250)")
	optionsSnapshotsChainCmd.Flags().String("sort", "", "Field to sort results by")

	optionsSnapshotsCmd.AddCommand(optionsSnapshotsChainCmd)
	optionsSnapshotsCmd.AddCommand(optionsSnapshotsContractCmd)

	optionsCmd.AddCommand(optionsSnapshotsCmd)
}
