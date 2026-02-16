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

// etfGlobalCmd is the parent command for all ETF Global partner data
// subcommands including analytics and constituents.
var etfGlobalCmd = &cobra.Command{
	Use:   "etf-global",
	Short: "ETF Global partner data commands",
	Long:  "Access ETF Global partner data including analytics scores, risk metrics, and constituent holdings for exchange-traded funds.",
}

// etfGlobalAnalyticsCmd retrieves ETF Global analytics data including
// quantitative scores, risk assessments, reward metrics, and letter grades
// for exchange-traded funds. Supports filtering by ticker, date, score
// thresholds, and grade. Usage: massive etf-global analytics --ticker SPY
var etfGlobalAnalyticsCmd = &cobra.Command{
	Use:   "analytics",
	Short: "Get ETF Global analytics scores and ratings",
	Long:  "Retrieve quantitative analytics data from ETF Global including composite scores for technical, sentiment, behavioral, fundamental, global, and quality factors, plus risk and reward assessments.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker, _ := cmd.Flags().GetString("ticker")
		processedDate, _ := cmd.Flags().GetString("processed-date")
		effectiveDate, _ := cmd.Flags().GetString("effective-date")
		riskTotalScore, _ := cmd.Flags().GetString("risk-total-score")
		rewardScore, _ := cmd.Flags().GetString("reward-score")
		quantTotalScore, _ := cmd.Flags().GetString("quant-total-score")
		quantGrade, _ := cmd.Flags().GetString("quant-grade")
		sort, _ := cmd.Flags().GetString("sort")
		limit, _ := cmd.Flags().GetString("limit")

		params := api.ETFGlobalAnalyticsParams{
			CompositeTicker: ticker,
			ProcessedDate:   processedDate,
			EffectiveDate:   effectiveDate,
			RiskTotalScore:  riskTotalScore,
			RewardScore:     rewardScore,
			QuantTotalScore: quantTotalScore,
			QuantGrade:      quantGrade,
			Sort:            sort,
			Limit:           limit,
		}

		result, err := client.GetETFGlobalAnalytics(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("ETF Global Analytics | Results: %d\n\n", result.Count)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tDATE\tGRADE\tQUANT\tREWARD\tRISK\tTECH\tSENT\tFUND\tQUAL\tGLOBAL\tBEHAV")
		fmt.Fprintln(w, "------\t----\t-----\t-----\t------\t----\t----\t----\t----\t----\t------\t-----")

		for _, a := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%.1f\t%.1f\t%.1f\t%.1f\t%.1f\t%.1f\t%.1f\t%.1f\t%.1f\n",
				a.CompositeTicker,
				a.EffectiveDate,
				a.QuantGrade,
				a.QuantTotalScore,
				a.RewardScore,
				a.RiskTotalScore,
				a.QuantCompositeTechnical,
				a.QuantCompositeSentiment,
				a.QuantCompositeFundamental,
				a.QuantCompositeQuality,
				a.QuantCompositeGlobal,
				a.QuantCompositeBehavioral,
			)
		}
		w.Flush()

		return nil
	},
}

// etfGlobalConstituentsCmd retrieves the underlying constituent holdings
// of an ETF from ETF Global, including weights, market values, share counts,
// and security identifiers. Supports filtering by composite ticker,
// constituent ticker, date, and security IDs.
// Usage: massive etf-global constituents --ticker SPY
var etfGlobalConstituentsCmd = &cobra.Command{
	Use:   "constituents",
	Short: "Get ETF constituent holdings",
	Long:  "Retrieve the underlying constituent holdings of exchange-traded funds from ETF Global, including weights, market values, shares held, and security identifiers (ISIN, FIGI, SEDOL).",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker, _ := cmd.Flags().GetString("ticker")
		constituentTicker, _ := cmd.Flags().GetString("constituent-ticker")
		effectiveDate, _ := cmd.Flags().GetString("effective-date")
		processedDate, _ := cmd.Flags().GetString("processed-date")
		usCode, _ := cmd.Flags().GetString("us-code")
		isin, _ := cmd.Flags().GetString("isin")
		figi, _ := cmd.Flags().GetString("figi")
		sedol, _ := cmd.Flags().GetString("sedol")
		sort, _ := cmd.Flags().GetString("sort")
		limit, _ := cmd.Flags().GetString("limit")

		params := api.ETFGlobalConstituentsParams{
			CompositeTicker:   ticker,
			ConstituentTicker: constituentTicker,
			EffectiveDate:     effectiveDate,
			ProcessedDate:     processedDate,
			USCode:            usCode,
			ISIN:              isin,
			FIGI:              figi,
			SEDOL:             sedol,
			Sort:              sort,
			Limit:             limit,
		}

		result, err := client.GetETFGlobalConstituents(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("ETF Global Constituents | Results: %d\n\n", result.Count)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "RANK\tETF\tTICKER\tNAME\tWEIGHT\tSHARES\tMKT VALUE\tASSET CLASS\tEXCHANGE")
		fmt.Fprintln(w, "----\t---\t------\t----\t------\t------\t---------\t-----------\t--------")

		for _, c := range result.Results {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%.2f%%\t%.0f\t%.2f\t%s\t%s\n",
				c.ConstituentRank,
				c.CompositeTicker,
				c.ConstituentTicker,
				c.ConstituentName,
				c.Weight,
				c.SharesHeld,
				c.MarketValue,
				c.AssetClass,
				c.Exchange,
			)
		}
		w.Flush()

		return nil
	},
}

// init registers the etf-global parent command with rootCmd and registers
// the analytics and constituents subcommands with their respective flags.
func init() {
	// Register parent command with root
	rootCmd.AddCommand(etfGlobalCmd)

	// Analytics subcommand flags
	etfGlobalAnalyticsCmd.Flags().String("ticker", "", "Filter by ETF ticker symbol (e.g., SPY)")
	etfGlobalAnalyticsCmd.Flags().String("processed-date", "", "Filter by data processing date (YYYY-MM-DD)")
	etfGlobalAnalyticsCmd.Flags().String("effective-date", "", "Filter by data effective date (YYYY-MM-DD)")
	etfGlobalAnalyticsCmd.Flags().String("risk-total-score", "", "Filter by total risk score")
	etfGlobalAnalyticsCmd.Flags().String("reward-score", "", "Filter by reward score")
	etfGlobalAnalyticsCmd.Flags().String("quant-total-score", "", "Filter by total quant score")
	etfGlobalAnalyticsCmd.Flags().String("quant-grade", "", "Filter by quant grade (A, B, C, D, F)")
	etfGlobalAnalyticsCmd.Flags().String("sort", "", "Sort field with direction (e.g., quant_total_score.desc)")
	etfGlobalAnalyticsCmd.Flags().String("limit", "20", "Number of results to return (max 5000)")
	etfGlobalCmd.AddCommand(etfGlobalAnalyticsCmd)

	// Constituents subcommand flags
	etfGlobalConstituentsCmd.Flags().String("ticker", "", "Filter by ETF ticker symbol (e.g., SPY)")
	etfGlobalConstituentsCmd.Flags().String("constituent-ticker", "", "Filter by constituent ticker (e.g., AAPL)")
	etfGlobalConstituentsCmd.Flags().String("effective-date", "", "Filter by data effective date (YYYY-MM-DD)")
	etfGlobalConstituentsCmd.Flags().String("processed-date", "", "Filter by data processing date (YYYY-MM-DD)")
	etfGlobalConstituentsCmd.Flags().String("us-code", "", "Filter by US market identifier code")
	etfGlobalConstituentsCmd.Flags().String("isin", "", "Filter by ISIN (International Securities ID)")
	etfGlobalConstituentsCmd.Flags().String("figi", "", "Filter by FIGI (Financial Instrument Global Identifier)")
	etfGlobalConstituentsCmd.Flags().String("sedol", "", "Filter by SEDOL (UK trading code)")
	etfGlobalConstituentsCmd.Flags().String("sort", "", "Sort field with direction (e.g., weight.desc)")
	etfGlobalConstituentsCmd.Flags().String("limit", "20", "Number of results to return (max 5000)")
	etfGlobalCmd.AddCommand(etfGlobalConstituentsCmd)
}
