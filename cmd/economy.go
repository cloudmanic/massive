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

// economyCmd is the parent command for all economic data subcommands
// including inflation, labor market, and treasury yield indicators.
var economyCmd = &cobra.Command{
	Use:   "economy",
	Short: "Economic data commands",
	Long:  "Access economic indicators from the Federal Reserve including inflation (CPI/PCE), labor market data, and treasury yields.",
}

// economyInflationCmd retrieves inflation indicator data from the Federal
// Reserve including headline and core CPI and PCE measures. Supports date
// range filtering, sorting, and limiting results.
// Usage: massive economy inflation --date-gte 2025-01-01 --date-lte 2025-12-31
var economyInflationCmd = &cobra.Command{
	Use:   "inflation",
	Short: "Get inflation indicators (CPI and PCE)",
	Long:  "Retrieve inflation indicators including headline and core inflation measures from the CPI and PCE indexes for tracking economic price trends.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		date, _ := cmd.Flags().GetString("date")
		dateGT, _ := cmd.Flags().GetString("date-gt")
		dateGTE, _ := cmd.Flags().GetString("date-gte")
		dateLT, _ := cmd.Flags().GetString("date-lt")
		dateLTE, _ := cmd.Flags().GetString("date-lte")
		sort, _ := cmd.Flags().GetString("sort")
		limit, _ := cmd.Flags().GetString("limit")

		params := api.InflationParams{
			Date:    date,
			DateGT:  dateGT,
			DateGTE: dateGTE,
			DateLT:  dateLT,
			DateLTE: dateLTE,
			Sort:    sort,
			Limit:   limit,
		}

		result, err := client.GetInflation(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		if len(result.Results) == 0 {
			fmt.Println("No inflation data found for the given parameters.")
			return nil
		}

		fmt.Printf("Inflation Data | Results: %d\n\n", len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "DATE\tCPI\tCPI CORE\tPCE\tPCE CORE\tPCE SPENDING")
		fmt.Fprintln(w, "----\t---\t--------\t---\t--------\t------------")

		for _, r := range result.Results {
			fmt.Fprintf(w, "%s\t%.3f\t%.3f\t%.3f\t%.3f\t%.1f\n",
				r.Date, r.CPI, r.CPICore, r.PCE, r.PCECore, r.PCESpending)
		}
		w.Flush()

		return nil
	},
}

// economyLaborMarketCmd retrieves labor market indicator data from the
// Federal Reserve including unemployment rate, labor force participation,
// average hourly earnings, and job openings.
// Usage: massive economy labor-market --date-gte 2025-01-01 --limit 12
var economyLaborMarketCmd = &cobra.Command{
	Use:   "labor-market",
	Short: "Get labor market indicators",
	Long:  "Retrieve key labor market indicators from the Federal Reserve including unemployment rate, labor force participation rate, average hourly earnings, and job openings.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		date, _ := cmd.Flags().GetString("date")
		dateGT, _ := cmd.Flags().GetString("date-gt")
		dateGTE, _ := cmd.Flags().GetString("date-gte")
		dateLT, _ := cmd.Flags().GetString("date-lt")
		dateLTE, _ := cmd.Flags().GetString("date-lte")
		sort, _ := cmd.Flags().GetString("sort")
		limit, _ := cmd.Flags().GetString("limit")

		params := api.LaborMarketParams{
			Date:    date,
			DateGT:  dateGT,
			DateGTE: dateGTE,
			DateLT:  dateLT,
			DateLTE: dateLTE,
			Sort:    sort,
			Limit:   limit,
		}

		result, err := client.GetLaborMarket(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		if len(result.Results) == 0 {
			fmt.Println("No labor market data found for the given parameters.")
			return nil
		}

		fmt.Printf("Labor Market Data | Results: %d\n\n", len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "DATE\tUNEMPLOYMENT\tPARTICIPATION\tHOURLY EARNINGS\tJOB OPENINGS")
		fmt.Fprintln(w, "----\t------------\t-------------\t---------------\t------------")

		for _, r := range result.Results {
			jobOpenings := "-"
			if r.JobOpenings > 0 {
				jobOpenings = fmt.Sprintf("%.0f", r.JobOpenings)
			}
			fmt.Fprintf(w, "%s\t%.1f%%\t%.1f%%\t$%.2f\t%s\n",
				r.Date, r.UnemploymentRate, r.LaborForceParticipationRate,
				r.AvgHourlyEarnings, jobOpenings)
		}
		w.Flush()

		return nil
	},
}

// economyTreasuryYieldsCmd retrieves daily treasury yield curve data from
// the Federal Reserve across multiple maturities from 1-month to 30-year.
// Usage: massive economy treasury-yields --date-gte 2026-01-01 --sort date.desc
var economyTreasuryYieldsCmd = &cobra.Command{
	Use:   "treasury-yields",
	Short: "Get treasury yield curve data",
	Long:  "Retrieve daily treasury yield curve data from the Federal Reserve across multiple maturities from 1-month to 30-year durations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		date, _ := cmd.Flags().GetString("date")
		dateGT, _ := cmd.Flags().GetString("date-gt")
		dateGTE, _ := cmd.Flags().GetString("date-gte")
		dateLT, _ := cmd.Flags().GetString("date-lt")
		dateLTE, _ := cmd.Flags().GetString("date-lte")
		sort, _ := cmd.Flags().GetString("sort")
		limit, _ := cmd.Flags().GetString("limit")

		params := api.TreasuryYieldParams{
			Date:    date,
			DateGT:  dateGT,
			DateGTE: dateGTE,
			DateLT:  dateLT,
			DateLTE: dateLTE,
			Sort:    sort,
			Limit:   limit,
		}

		result, err := client.GetTreasuryYields(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		if len(result.Results) == 0 {
			fmt.Println("No treasury yield data found for the given parameters.")
			return nil
		}

		fmt.Printf("Treasury Yields | Results: %d\n\n", len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "DATE\t1M\t3M\t6M\t1Y\t2Y\t3Y\t5Y\t7Y\t10Y\t20Y\t30Y")
		fmt.Fprintln(w, "----\t--\t--\t--\t--\t--\t--\t--\t--\t---\t---\t---")

		for _, r := range result.Results {
			fmt.Fprintf(w, "%s\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\n",
				r.Date,
				r.Yield1Month, r.Yield3Month, r.Yield6Month,
				r.Yield1Year, r.Yield2Year, r.Yield3Year,
				r.Yield5Year, r.Yield7Year, r.Yield10Year,
				r.Yield20Year, r.Yield30Year)
		}
		w.Flush()

		return nil
	},
}

// addEconomyDateFlags registers the common date filtering flags used by all
// economy subcommands including exact date, greater than, greater than or
// equal, less than, and less than or equal date operators.
func addEconomyDateFlags(cmd *cobra.Command) {
	cmd.Flags().String("date", "", "Exact observation date (YYYY-MM-DD)")
	cmd.Flags().String("date-gt", "", "Date greater than (YYYY-MM-DD)")
	cmd.Flags().String("date-gte", "", "Date greater than or equal (YYYY-MM-DD)")
	cmd.Flags().String("date-lt", "", "Date less than (YYYY-MM-DD)")
	cmd.Flags().String("date-lte", "", "Date less than or equal (YYYY-MM-DD)")
	cmd.Flags().String("sort", "date.desc", "Sort order (date.asc, date.desc)")
	cmd.Flags().String("limit", "100", "Max number of results (max 50000)")
}

// init registers the economy parent command with the root command and
// adds all economy subcommands with their respective flags.
func init() {
	rootCmd.AddCommand(economyCmd)

	addEconomyDateFlags(economyInflationCmd)
	addEconomyDateFlags(economyLaborMarketCmd)
	addEconomyDateFlags(economyTreasuryYieldsCmd)

	economyCmd.AddCommand(economyInflationCmd)
	economyCmd.AddCommand(economyLaborMarketCmd)
	economyCmd.AddCommand(economyTreasuryYieldsCmd)
}
