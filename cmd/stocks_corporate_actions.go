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

// stocksDividendsCmd retrieves historical cash dividend distributions for
// a specified stock ticker. Supports filtering by ex-dividend date range,
// frequency, distribution type, and result limit. Output can be formatted
// as a table or JSON. Usage: massive stocks dividends --ticker AAPL
var stocksDividendsCmd = &cobra.Command{
	Use:   "dividends",
	Short: "Get historical dividend data for stocks",
	Long:  "Retrieve historical cash dividend distributions including declaration dates, ex-dividend dates, record dates, pay dates, cash amounts, frequencies, and split-adjusted values.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker, _ := cmd.Flags().GetString("ticker")
		ticker = strings.ToUpper(ticker)
		exDividendDate, _ := cmd.Flags().GetString("ex-dividend-date")
		exDividendDateGTE, _ := cmd.Flags().GetString("ex-dividend-date-gte")
		exDividendDateLTE, _ := cmd.Flags().GetString("ex-dividend-date-lte")
		exDividendDateGT, _ := cmd.Flags().GetString("ex-dividend-date-gt")
		exDividendDateLT, _ := cmd.Flags().GetString("ex-dividend-date-lt")
		frequency, _ := cmd.Flags().GetString("frequency")
		distributionType, _ := cmd.Flags().GetString("distribution-type")
		sort, _ := cmd.Flags().GetString("sort")
		limit, _ := cmd.Flags().GetString("limit")

		params := api.DividendsParams{
			Ticker:            ticker,
			ExDividendDate:    exDividendDate,
			ExDividendDateGT:  exDividendDateGT,
			ExDividendDateGTE: exDividendDateGTE,
			ExDividendDateLT:  exDividendDateLT,
			ExDividendDateLTE: exDividendDateLTE,
			Frequency:         frequency,
			DistributionType:  distributionType,
			Sort:              sort,
			Limit:             limit,
		}

		result, err := client.GetDividends(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Dividends: %d result(s)\n\n", len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tEX-DIV DATE\tPAY DATE\tCASH AMT\tCURRENCY\tFREQ\tTYPE\tSPLIT-ADJ AMT")
		fmt.Fprintln(w, "------\t-----------\t--------\t--------\t--------\t----\t----\t-------------")

		for _, d := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%.6f\t%s\t%d\t%s\t%.6f\n",
				d.Ticker, d.ExDividendDate, d.PayDate,
				d.CashAmount, d.Currency, d.Frequency,
				d.DistributionType, d.SplitAdjustedCashAmount)
		}
		w.Flush()

		return nil
	},
}

// stocksSplitsCmd retrieves historical stock split events for a specified
// stock ticker. Supports filtering by execution date range, adjustment
// type, and result limit. Output can be formatted as a table or JSON.
// Usage: massive stocks splits --ticker AAPL
var stocksSplitsCmd = &cobra.Command{
	Use:   "splits",
	Short: "Get historical stock split data",
	Long:  "Retrieve historical stock split events including execution dates, split ratios (split_from and split_to), adjustment types (forward_split, reverse_split, stock_dividend), and historical adjustment factors.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker, _ := cmd.Flags().GetString("ticker")
		ticker = strings.ToUpper(ticker)
		executionDate, _ := cmd.Flags().GetString("execution-date")
		executionDateGTE, _ := cmd.Flags().GetString("execution-date-gte")
		executionDateLTE, _ := cmd.Flags().GetString("execution-date-lte")
		executionDateGT, _ := cmd.Flags().GetString("execution-date-gt")
		executionDateLT, _ := cmd.Flags().GetString("execution-date-lt")
		adjustmentType, _ := cmd.Flags().GetString("adjustment-type")
		sort, _ := cmd.Flags().GetString("sort")
		limit, _ := cmd.Flags().GetString("limit")

		params := api.SplitsParams{
			Ticker:           ticker,
			ExecutionDate:    executionDate,
			ExecutionDateGT:  executionDateGT,
			ExecutionDateGTE: executionDateGTE,
			ExecutionDateLT:  executionDateLT,
			ExecutionDateLTE: executionDateLTE,
			AdjustmentType:   adjustmentType,
			Sort:             sort,
			Limit:            limit,
		}

		result, err := client.GetSplits(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Splits: %d result(s)\n\n", len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tEXECUTION DATE\tSPLIT FROM\tSPLIT TO\tTYPE\tADJ FACTOR")
		fmt.Fprintln(w, "------\t--------------\t----------\t--------\t----\t----------")

		for _, s := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t%.0f\t%.0f\t%s\t%.6f\n",
				s.Ticker, s.ExecutionDate,
				s.SplitFrom, s.SplitTo,
				s.AdjustmentType, s.HistoricalAdjustmentFactor)
		}
		w.Flush()

		return nil
	},
}

// init registers the dividends and splits commands and their flags under
// the stocks parent command.
func init() {
	// Dividends command flags
	stocksDividendsCmd.Flags().String("ticker", "", "Stock ticker symbol (e.g. AAPL)")
	stocksDividendsCmd.Flags().String("ex-dividend-date", "", "Exact ex-dividend date (YYYY-MM-DD)")
	stocksDividendsCmd.Flags().String("ex-dividend-date-gt", "", "Ex-dividend date greater than (YYYY-MM-DD)")
	stocksDividendsCmd.Flags().String("ex-dividend-date-gte", "", "Ex-dividend date greater than or equal (YYYY-MM-DD)")
	stocksDividendsCmd.Flags().String("ex-dividend-date-lt", "", "Ex-dividend date less than (YYYY-MM-DD)")
	stocksDividendsCmd.Flags().String("ex-dividend-date-lte", "", "Ex-dividend date less than or equal (YYYY-MM-DD)")
	stocksDividendsCmd.Flags().String("frequency", "", "Dividend frequency (0=one-time, 1=annual, 4=quarterly, 12=monthly)")
	stocksDividendsCmd.Flags().String("distribution-type", "", "Distribution type (recurring, special, supplemental, irregular, unknown)")
	stocksDividendsCmd.Flags().String("sort", "", "Sort field with direction (e.g. ex_dividend_date.desc)")
	stocksDividendsCmd.Flags().String("limit", "100", "Max number of results (max 5000)")

	// Splits command flags
	stocksSplitsCmd.Flags().String("ticker", "", "Stock ticker symbol (e.g. AAPL)")
	stocksSplitsCmd.Flags().String("execution-date", "", "Exact execution date (YYYY-MM-DD)")
	stocksSplitsCmd.Flags().String("execution-date-gt", "", "Execution date greater than (YYYY-MM-DD)")
	stocksSplitsCmd.Flags().String("execution-date-gte", "", "Execution date greater than or equal (YYYY-MM-DD)")
	stocksSplitsCmd.Flags().String("execution-date-lt", "", "Execution date less than (YYYY-MM-DD)")
	stocksSplitsCmd.Flags().String("execution-date-lte", "", "Execution date less than or equal (YYYY-MM-DD)")
	stocksSplitsCmd.Flags().String("adjustment-type", "", "Adjustment type (forward_split, reverse_split, stock_dividend)")
	stocksSplitsCmd.Flags().String("sort", "", "Sort field with direction (e.g. execution_date.desc)")
	stocksSplitsCmd.Flags().String("limit", "100", "Max number of results (max 5000)")

	// Register under stocks parent command
	stocksCmd.AddCommand(stocksDividendsCmd)
	stocksCmd.AddCommand(stocksSplitsCmd)
}
