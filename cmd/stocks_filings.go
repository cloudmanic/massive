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

// stocksFilingsCmd is the parent command for all SEC filing data subcommands
// including 10-K sections, risk factors, and risk categories.
var stocksFilingsCmd = &cobra.Command{
	Use:   "filings",
	Short: "SEC filing data commands",
	Long:  "Retrieve SEC filing data including 10-K sections, risk factor disclosures, and risk factor taxonomy categories.",
}

// stocksFilingsSectionsCmd retrieves plain-text content of specific sections
// from SEC 10-K filings. Supports filtering by ticker, CIK, section type,
// filing date, and period end date.
// Usage: massive stocks filings sections --ticker AAPL --section risk_factors
var stocksFilingsSectionsCmd = &cobra.Command{
	Use:   "sections",
	Short: "Get 10-K section content from SEC filings",
	Long:  "Retrieve plain-text content of specific sections from SEC 10-K filings, such as business descriptions, risk factors, and management discussion.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker, _ := cmd.Flags().GetString("ticker")
		cik, _ := cmd.Flags().GetString("cik")
		section, _ := cmd.Flags().GetString("section")
		filingDate, _ := cmd.Flags().GetString("filing-date")
		filingDateGt, _ := cmd.Flags().GetString("filing-date-gt")
		filingDateLt, _ := cmd.Flags().GetString("filing-date-lt")
		periodEnd, _ := cmd.Flags().GetString("period-end")
		periodEndGt, _ := cmd.Flags().GetString("period-end-gt")
		periodEndLt, _ := cmd.Flags().GetString("period-end-lt")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.SECFilingSectionsParams{
			Ticker:       strings.ToUpper(ticker),
			CIK:          cik,
			Section:      section,
			FilingDate:   filingDate,
			FilingDateGt: filingDateGt,
			FilingDateLt: filingDateLt,
			PeriodEnd:    periodEnd,
			PeriodEndGt:  periodEndGt,
			PeriodEndLt:  periodEndLt,
			Limit:        limit,
			Sort:         sort,
		}

		result, err := client.GetSECFilingSections(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Results: %d\n\n", len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tSECTION\tFILING DATE\tPERIOD END\tTEXT PREVIEW")
		fmt.Fprintln(w, "------\t-------\t-----------\t----------\t------------")

		for _, s := range result.Results {
			preview := s.Text
			if len(preview) > 80 {
				preview = preview[:80] + "..."
			}
			// Replace newlines so the preview stays on one table row
			preview = strings.ReplaceAll(preview, "\n", " ")
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				s.Ticker, s.Section, s.FilingDate, s.PeriodEnd, preview)
		}
		w.Flush()

		return nil
	},
}

// stocksFilingsRiskFactorsCmd retrieves standardized, machine-readable risk
// factor disclosures from SEC filings. Each risk is classified into a
// three-level taxonomy with supporting text.
// Usage: massive stocks filings risk-factors --ticker AAPL
var stocksFilingsRiskFactorsCmd = &cobra.Command{
	Use:   "risk-factors",
	Short: "Get categorized risk factor disclosures from SEC filings",
	Long:  "Retrieve standardized, machine-readable risk factor disclosures from SEC filings, classified into a three-level taxonomy with supporting text.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker, _ := cmd.Flags().GetString("ticker")
		cik, _ := cmd.Flags().GetString("cik")
		filingDate, _ := cmd.Flags().GetString("filing-date")
		filingDateGt, _ := cmd.Flags().GetString("filing-date-gt")
		filingDateLt, _ := cmd.Flags().GetString("filing-date-lt")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.RiskFactorsParams{
			Ticker:       strings.ToUpper(ticker),
			CIK:          cik,
			FilingDate:   filingDate,
			FilingDateGt: filingDateGt,
			FilingDateLt: filingDateLt,
			Limit:        limit,
			Sort:         sort,
		}

		result, err := client.GetRiskFactors(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Results: %d\n\n", len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tFILING DATE\tPRIMARY\tSECONDARY\tTERTIARY")
		fmt.Fprintln(w, "------\t-----------\t-------\t---------\t--------")

		for _, rf := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				rf.Ticker, rf.FilingDate,
				rf.PrimaryCategory, rf.SecondaryCategory, rf.TertiaryCategory)
		}
		w.Flush()

		return nil
	},
}

// stocksFilingsRiskCategoriesCmd retrieves the taxonomy used to classify
// risk factors in SEC filing disclosures. Each entry includes a three-level
// classification with a description and taxonomy version.
// Usage: massive stocks filings risk-categories --primary-category financial_and_market
var stocksFilingsRiskCategoriesCmd = &cobra.Command{
	Use:   "risk-categories",
	Short: "Get the risk factor classification taxonomy",
	Long:  "Retrieve the taxonomy used to classify risk factors in SEC filing disclosures, including primary, secondary, and tertiary categories with descriptions.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		primaryCategory, _ := cmd.Flags().GetString("primary-category")
		secondaryCategory, _ := cmd.Flags().GetString("secondary-category")
		tertiaryCategory, _ := cmd.Flags().GetString("tertiary-category")
		taxonomy, _ := cmd.Flags().GetString("taxonomy")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.RiskCategoriesParams{
			PrimaryCategory:   primaryCategory,
			SecondaryCategory: secondaryCategory,
			TertiaryCategory:  tertiaryCategory,
			Taxonomy:          taxonomy,
			Limit:             limit,
			Sort:              sort,
		}

		result, err := client.GetRiskCategories(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Results: %d\n\n", len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "PRIMARY\tSECONDARY\tTERTIARY\tTAXONOMY\tDESCRIPTION")
		fmt.Fprintln(w, "-------\t---------\t--------\t--------\t-----------")

		for _, rc := range result.Results {
			desc := rc.Description
			if len(desc) > 60 {
				desc = desc[:60] + "..."
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%.1f\t%s\n",
				rc.PrimaryCategory, rc.SecondaryCategory, rc.TertiaryCategory,
				rc.Taxonomy, desc)
		}
		w.Flush()

		return nil
	},
}

// init registers the filings parent command and its subcommands under the
// stocks parent command, along with all their respective flags.
func init() {
	// Register sections command flags
	stocksFilingsSectionsCmd.Flags().String("ticker", "", "Filter by stock ticker symbol")
	stocksFilingsSectionsCmd.Flags().String("cik", "", "Filter by SEC Central Index Key (10-digit zero-padded)")
	stocksFilingsSectionsCmd.Flags().String("section", "", "Filter by section type (business, risk_factors, etc.)")
	stocksFilingsSectionsCmd.Flags().String("filing-date", "", "Filter by exact filing date (YYYY-MM-DD)")
	stocksFilingsSectionsCmd.Flags().String("filing-date-gt", "", "Filing date greater than (YYYY-MM-DD)")
	stocksFilingsSectionsCmd.Flags().String("filing-date-lt", "", "Filing date less than (YYYY-MM-DD)")
	stocksFilingsSectionsCmd.Flags().String("period-end", "", "Filter by exact period end date (YYYY-MM-DD)")
	stocksFilingsSectionsCmd.Flags().String("period-end-gt", "", "Period end date greater than (YYYY-MM-DD)")
	stocksFilingsSectionsCmd.Flags().String("period-end-lt", "", "Period end date less than (YYYY-MM-DD)")
	stocksFilingsSectionsCmd.Flags().String("limit", "10", "Number of results to return (max 9999)")
	stocksFilingsSectionsCmd.Flags().String("sort", "period_end.desc", "Sort order (e.g., period_end.desc, filing_date.asc)")

	// Register risk factors command flags
	stocksFilingsRiskFactorsCmd.Flags().String("ticker", "", "Filter by stock ticker symbol")
	stocksFilingsRiskFactorsCmd.Flags().String("cik", "", "Filter by SEC Central Index Key (10-digit zero-padded)")
	stocksFilingsRiskFactorsCmd.Flags().String("filing-date", "", "Filter by exact filing date (YYYY-MM-DD)")
	stocksFilingsRiskFactorsCmd.Flags().String("filing-date-gt", "", "Filing date greater than (YYYY-MM-DD)")
	stocksFilingsRiskFactorsCmd.Flags().String("filing-date-lt", "", "Filing date less than (YYYY-MM-DD)")
	stocksFilingsRiskFactorsCmd.Flags().String("limit", "10", "Number of results to return (max 49999)")
	stocksFilingsRiskFactorsCmd.Flags().String("sort", "filing_date.desc", "Sort order (e.g., filing_date.desc, filing_date.asc)")

	// Register risk categories command flags
	stocksFilingsRiskCategoriesCmd.Flags().String("primary-category", "", "Filter by primary risk category")
	stocksFilingsRiskCategoriesCmd.Flags().String("secondary-category", "", "Filter by secondary risk category")
	stocksFilingsRiskCategoriesCmd.Flags().String("tertiary-category", "", "Filter by tertiary risk category")
	stocksFilingsRiskCategoriesCmd.Flags().String("taxonomy", "", "Filter by taxonomy version (e.g., 1.0)")
	stocksFilingsRiskCategoriesCmd.Flags().String("limit", "20", "Number of results to return (max 999)")
	stocksFilingsRiskCategoriesCmd.Flags().String("sort", "taxonomy.desc", "Sort order (e.g., taxonomy.desc, primary_category.asc)")

	// Register subcommands under the filings parent
	stocksFilingsCmd.AddCommand(stocksFilingsSectionsCmd)
	stocksFilingsCmd.AddCommand(stocksFilingsRiskFactorsCmd)
	stocksFilingsCmd.AddCommand(stocksFilingsRiskCategoriesCmd)

	// Register filings under the stocks parent command
	stocksCmd.AddCommand(stocksFilingsCmd)
}
