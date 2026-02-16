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

// benzingaCmd is the parent command for all Benzinga partner data subcommands
// including news, ratings, earnings, guidance, and analysts.
var benzingaCmd = &cobra.Command{
	Use:   "benzinga",
	Short: "Benzinga partner data commands",
	Long:  "Access Benzinga financial data including news articles, analyst ratings, earnings reports, corporate guidance, and analyst details.",
}

// benzingaNewsCmd retrieves Benzinga news articles from the Massive API.
// Supports filtering by ticker symbols, publication date range, channels,
// tags, and author. Results can be displayed as a table or raw JSON.
// Usage: massive benzinga news --tickers AAPL --limit 5
var benzingaNewsCmd = &cobra.Command{
	Use:   "news",
	Short: "Get Benzinga news articles",
	Long:  "Retrieve Benzinga news articles with optional filtering by tickers, publication date range, channels, tags, author, and result limit.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		tickers, _ := cmd.Flags().GetString("tickers")
		tickersAnyOf, _ := cmd.Flags().GetString("tickers-any-of")
		published, _ := cmd.Flags().GetString("published")
		publishedGte, _ := cmd.Flags().GetString("published-from")
		publishedLte, _ := cmd.Flags().GetString("published-to")
		channels, _ := cmd.Flags().GetString("channels")
		tags, _ := cmd.Flags().GetString("tags")
		author, _ := cmd.Flags().GetString("author")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.BenzingaNewsParams{
			Tickers:      strings.ToUpper(tickers),
			TickersAnyOf: strings.ToUpper(tickersAnyOf),
			Published:    published,
			PublishedGte: publishedGte,
			PublishedLte: publishedLte,
			Channels:     channels,
			Tags:         tags,
			Author:       author,
			Limit:        limit,
			Sort:         sort,
		}

		result, err := client.GetBenzingaNews(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		// Display results count header
		fmt.Printf("Benzinga News Articles: %d\n\n", result.Count)

		if len(result.Results) == 0 {
			fmt.Println("No news articles found.")
			return nil
		}

		// Print each news article in a readable table format
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "DATE\tAUTHOR\tTICKERS\tTITLE")
		fmt.Fprintln(w, "----\t------\t-------\t-----")

		for _, article := range result.Results {
			date := formatBenzingaDate(article.Published)
			tickers := truncateBenzingaString(strings.Join(article.Tickers, ","), 20)
			title := truncateBenzingaString(article.Title, 60)

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				date, article.Author, tickers, title)
		}
		w.Flush()

		return nil
	},
}

// benzingaRatingsCmd retrieves Benzinga analyst ratings from the Massive API.
// Supports filtering by ticker, date range, rating action, price target action,
// and importance level. Results can be displayed as a table or raw JSON.
// Usage: massive benzinga ratings --ticker AAPL --limit 10
var benzingaRatingsCmd = &cobra.Command{
	Use:   "ratings",
	Short: "Get Benzinga analyst ratings",
	Long:  "Retrieve Benzinga analyst ratings with optional filtering by ticker, date range, rating action, price target action, importance, and result limit.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker, _ := cmd.Flags().GetString("ticker")
		tickerAnyOf, _ := cmd.Flags().GetString("ticker-any-of")
		date, _ := cmd.Flags().GetString("date")
		dateGte, _ := cmd.Flags().GetString("date-from")
		dateLte, _ := cmd.Flags().GetString("date-to")
		importance, _ := cmd.Flags().GetString("importance")
		ratingAction, _ := cmd.Flags().GetString("rating-action")
		priceTargetAction, _ := cmd.Flags().GetString("price-target-action")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.BenzingaRatingsParams{
			Ticker:            strings.ToUpper(ticker),
			TickerAnyOf:       strings.ToUpper(tickerAnyOf),
			Date:              date,
			DateGte:           dateGte,
			DateLte:           dateLte,
			Importance:        importance,
			RatingAction:      ratingAction,
			PriceTargetAction: priceTargetAction,
			Limit:             limit,
			Sort:              sort,
		}

		result, err := client.GetBenzingaRatings(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		// Display results count header
		fmt.Printf("Benzinga Analyst Ratings: %d\n\n", result.Count)

		if len(result.Results) == 0 {
			fmt.Println("No analyst ratings found.")
			return nil
		}

		// Print each rating in a readable table format
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "DATE\tTICKER\tFIRM\tACTION\tRATING\tPT\tPREV PT")
		fmt.Fprintln(w, "----\t------\t----\t------\t------\t--\t-------")

		for _, rating := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%.2f\t%.2f\n",
				rating.Date, rating.Ticker,
				truncateBenzingaString(rating.Firm, 20),
				rating.RatingAction, rating.Rating,
				rating.PriceTarget, rating.PreviousPriceTarget)
		}
		w.Flush()

		return nil
	},
}

// benzingaEarningsCmd retrieves Benzinga earnings data from the Massive API.
// Supports filtering by ticker, date range, fiscal period, date status,
// and importance level. Results can be displayed as a table or raw JSON.
// Usage: massive benzinga earnings --ticker AAPL --limit 10
var benzingaEarningsCmd = &cobra.Command{
	Use:   "earnings",
	Short: "Get Benzinga earnings reports",
	Long:  "Retrieve Benzinga earnings reports with optional filtering by ticker, date range, fiscal period, date status, importance, and result limit.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker, _ := cmd.Flags().GetString("ticker")
		tickerAnyOf, _ := cmd.Flags().GetString("ticker-any-of")
		date, _ := cmd.Flags().GetString("date")
		dateGte, _ := cmd.Flags().GetString("date-from")
		dateLte, _ := cmd.Flags().GetString("date-to")
		dateStatus, _ := cmd.Flags().GetString("date-status")
		fiscalYear, _ := cmd.Flags().GetString("fiscal-year")
		fiscalPeriod, _ := cmd.Flags().GetString("fiscal-period")
		importance, _ := cmd.Flags().GetString("importance")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.BenzingaEarningsParams{
			Ticker:       strings.ToUpper(ticker),
			TickerAnyOf:  strings.ToUpper(tickerAnyOf),
			Date:         date,
			DateGte:      dateGte,
			DateLte:      dateLte,
			DateStatus:   dateStatus,
			FiscalYear:   fiscalYear,
			FiscalPeriod: fiscalPeriod,
			Importance:   importance,
			Limit:        limit,
			Sort:         sort,
		}

		result, err := client.GetBenzingaEarnings(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		// Display results count header
		fmt.Printf("Benzinga Earnings Reports: %d\n\n", result.Count)

		if len(result.Results) == 0 {
			fmt.Println("No earnings reports found.")
			return nil
		}

		// Print each earnings record in a readable table format
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "DATE\tTICKER\tCOMPANY\tSTATUS\tACT EPS\tEST EPS\tEPS SURP%\tPERIOD")
		fmt.Fprintln(w, "----\t------\t-------\t------\t-------\t-------\t---------\t------")

		for _, earn := range result.Results {
			company := truncateBenzingaString(earn.CompanyName, 20)
			period := fmt.Sprintf("%s %d", earn.FiscalPeriod, earn.FiscalYear)

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%.2f\t%.2f\t%.2f%%\t%s\n",
				earn.Date, earn.Ticker, company, earn.DateStatus,
				earn.ActualEPS, earn.EstimatedEPS, earn.EPSSurprisePercent,
				period)
		}
		w.Flush()

		return nil
	},
}

// benzingaGuidanceCmd retrieves Benzinga corporate guidance data from the
// Massive API. Supports filtering by ticker, date range, fiscal period,
// positioning, and importance level. Results can be displayed as a table
// or raw JSON.
// Usage: massive benzinga guidance --ticker AAPL --limit 10
var benzingaGuidanceCmd = &cobra.Command{
	Use:   "guidance",
	Short: "Get Benzinga corporate guidance",
	Long:  "Retrieve Benzinga corporate guidance with optional filtering by ticker, date range, fiscal period, positioning, importance, and result limit.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker, _ := cmd.Flags().GetString("ticker")
		tickerAnyOf, _ := cmd.Flags().GetString("ticker-any-of")
		date, _ := cmd.Flags().GetString("date")
		dateGte, _ := cmd.Flags().GetString("date-from")
		dateLte, _ := cmd.Flags().GetString("date-to")
		positioning, _ := cmd.Flags().GetString("positioning")
		fiscalYear, _ := cmd.Flags().GetString("fiscal-year")
		fiscalPeriod, _ := cmd.Flags().GetString("fiscal-period")
		importance, _ := cmd.Flags().GetString("importance")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.BenzingaGuidanceParams{
			Ticker:       strings.ToUpper(ticker),
			TickerAnyOf:  strings.ToUpper(tickerAnyOf),
			Date:         date,
			DateGte:      dateGte,
			DateLte:      dateLte,
			Positioning:  positioning,
			FiscalYear:   fiscalYear,
			FiscalPeriod: fiscalPeriod,
			Importance:   importance,
			Limit:        limit,
			Sort:         sort,
		}

		result, err := client.GetBenzingaGuidance(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		// Display results count header
		fmt.Printf("Benzinga Corporate Guidance: %d\n\n", result.Count)

		if len(result.Results) == 0 {
			fmt.Println("No corporate guidance found.")
			return nil
		}

		// Print each guidance record in a readable table format
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "DATE\tTICKER\tCOMPANY\tPOS\tEPS RANGE\tREV RANGE\tPERIOD")
		fmt.Fprintln(w, "----\t------\t-------\t---\t---------\t---------\t------")

		for _, guide := range result.Results {
			company := truncateBenzingaString(guide.CompanyName, 20)
			period := fmt.Sprintf("%s %d", guide.FiscalPeriod, guide.FiscalYear)
			epsRange := fmt.Sprintf("%.2f-%.2f", guide.MinEPSGuidance, guide.MaxEPSGuidance)
			revRange := fmt.Sprintf("%.0f-%.0f", guide.MinRevenueGuidance, guide.MaxRevenueGuidance)

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				guide.Date, guide.Ticker, company, guide.Positioning,
				epsRange, revRange, period)
		}
		w.Flush()

		return nil
	},
}

// benzingaAnalystsCmd retrieves Benzinga analyst details from the Massive API.
// Supports filtering by analyst ID, firm ID, analyst name, and firm name.
// Results can be displayed as a table or raw JSON.
// Usage: massive benzinga analysts --firm-name "Goldman Sachs" --limit 10
var benzingaAnalystsCmd = &cobra.Command{
	Use:   "analysts",
	Short: "Get Benzinga analyst details",
	Long:  "Retrieve Benzinga analyst details with optional filtering by analyst ID, firm ID, analyst name, firm name, and result limit.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		benzingaID, _ := cmd.Flags().GetString("benzinga-id")
		benzingaFirmID, _ := cmd.Flags().GetString("benzinga-firm-id")
		fullName, _ := cmd.Flags().GetString("full-name")
		firmName, _ := cmd.Flags().GetString("firm-name")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.BenzingaAnalystsParams{
			BenzingaID:     benzingaID,
			BenzingaFirmID: benzingaFirmID,
			FullName:       fullName,
			FirmName:       firmName,
			Limit:          limit,
			Sort:           sort,
		}

		result, err := client.GetBenzingaAnalysts(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		// Display results count header
		fmt.Printf("Benzinga Analysts: %d\n\n", len(result.Results))

		if len(result.Results) == 0 {
			fmt.Println("No analysts found.")
			return nil
		}

		// Print each analyst in a readable table format
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tFIRM\tSCORE\tSUCCESS RATE\tAVG RETURN\tRATINGS")
		fmt.Fprintln(w, "----\t----\t-----\t------------\t----------\t-------")

		for _, analyst := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t%.1f\t%.0f%%\t%.1f%%\t%.0f\n",
				truncateBenzingaString(analyst.FullName, 25),
				truncateBenzingaString(analyst.FirmName, 20),
				analyst.SmartScore,
				analyst.OverallSuccessRate*100,
				analyst.OverallAvgReturn,
				analyst.TotalRatings)
		}
		w.Flush()

		return nil
	},
}

// formatBenzingaDate extracts the date portion from an ISO 8601 timestamp
// string. If the string is shorter than 10 characters, it returns the
// original string unchanged.
func formatBenzingaDate(ts string) string {
	if len(ts) >= 10 {
		return ts[:10]
	}
	return ts
}

// truncateBenzingaString shortens a string to the specified maximum length,
// appending "..." if truncation occurs. Returns the original string if it
// is within the limit.
func truncateBenzingaString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// init registers the benzinga parent command with the root command and
// registers all Benzinga subcommands (news, ratings, earnings, guidance,
// analysts) along with their flags.
func init() {
	rootCmd.AddCommand(benzingaCmd)

	// Register subcommands under benzinga parent
	benzingaCmd.AddCommand(benzingaNewsCmd)
	benzingaCmd.AddCommand(benzingaRatingsCmd)
	benzingaCmd.AddCommand(benzingaEarningsCmd)
	benzingaCmd.AddCommand(benzingaGuidanceCmd)
	benzingaCmd.AddCommand(benzingaAnalystsCmd)

	// Benzinga News flags
	benzingaNewsCmd.Flags().String("tickers", "", "Filter by ticker symbols (e.g., AAPL)")
	benzingaNewsCmd.Flags().String("tickers-any-of", "", "Filter by any of these tickers (comma-separated)")
	benzingaNewsCmd.Flags().String("published", "", "Filter by exact publication date (ISO 8601)")
	benzingaNewsCmd.Flags().String("published-from", "", "Filter articles published on or after this date (ISO 8601)")
	benzingaNewsCmd.Flags().String("published-to", "", "Filter articles published on or before this date (ISO 8601)")
	benzingaNewsCmd.Flags().String("channels", "", "Filter by news channels (e.g., Tech)")
	benzingaNewsCmd.Flags().String("tags", "", "Filter by content tags")
	benzingaNewsCmd.Flags().String("author", "", "Filter by author name")
	benzingaNewsCmd.Flags().String("limit", "10", "Number of results to return (max 50000)")
	benzingaNewsCmd.Flags().String("sort", "published.desc", "Sort order (e.g., published.asc, published.desc)")

	// Benzinga Ratings flags
	benzingaRatingsCmd.Flags().String("ticker", "", "Filter by ticker symbol (e.g., AAPL)")
	benzingaRatingsCmd.Flags().String("ticker-any-of", "", "Filter by any of these tickers (comma-separated)")
	benzingaRatingsCmd.Flags().String("date", "", "Filter by exact date (YYYY-MM-DD)")
	benzingaRatingsCmd.Flags().String("date-from", "", "Filter ratings on or after this date (YYYY-MM-DD)")
	benzingaRatingsCmd.Flags().String("date-to", "", "Filter ratings on or before this date (YYYY-MM-DD)")
	benzingaRatingsCmd.Flags().String("importance", "", "Filter by importance level (0-5)")
	benzingaRatingsCmd.Flags().String("rating-action", "", "Filter by rating action (upgrades, downgrades, maintains, initiates_coverage_on, etc.)")
	benzingaRatingsCmd.Flags().String("price-target-action", "", "Filter by price target action (raises, lowers, maintains, announces, sets)")
	benzingaRatingsCmd.Flags().String("limit", "10", "Number of results to return (max 50000)")
	benzingaRatingsCmd.Flags().String("sort", "date.desc", "Sort order (e.g., date.asc, date.desc)")

	// Benzinga Earnings flags
	benzingaEarningsCmd.Flags().String("ticker", "", "Filter by ticker symbol (e.g., AAPL)")
	benzingaEarningsCmd.Flags().String("ticker-any-of", "", "Filter by any of these tickers (comma-separated)")
	benzingaEarningsCmd.Flags().String("date", "", "Filter by exact date (YYYY-MM-DD)")
	benzingaEarningsCmd.Flags().String("date-from", "", "Filter earnings on or after this date (YYYY-MM-DD)")
	benzingaEarningsCmd.Flags().String("date-to", "", "Filter earnings on or before this date (YYYY-MM-DD)")
	benzingaEarningsCmd.Flags().String("date-status", "", "Filter by date status (projected, confirmed)")
	benzingaEarningsCmd.Flags().String("fiscal-year", "", "Filter by fiscal year (e.g., 2026)")
	benzingaEarningsCmd.Flags().String("fiscal-period", "", "Filter by fiscal period (Q1, Q2, Q3, Q4, H1, FY)")
	benzingaEarningsCmd.Flags().String("importance", "", "Filter by importance level (0-5)")
	benzingaEarningsCmd.Flags().String("limit", "10", "Number of results to return (max 50000)")
	benzingaEarningsCmd.Flags().String("sort", "last_updated.desc", "Sort order (e.g., date.asc, last_updated.desc)")

	// Benzinga Guidance flags
	benzingaGuidanceCmd.Flags().String("ticker", "", "Filter by ticker symbol (e.g., AAPL)")
	benzingaGuidanceCmd.Flags().String("ticker-any-of", "", "Filter by any of these tickers (comma-separated)")
	benzingaGuidanceCmd.Flags().String("date", "", "Filter by exact date (YYYY-MM-DD)")
	benzingaGuidanceCmd.Flags().String("date-from", "", "Filter guidance on or after this date (YYYY-MM-DD)")
	benzingaGuidanceCmd.Flags().String("date-to", "", "Filter guidance on or before this date (YYYY-MM-DD)")
	benzingaGuidanceCmd.Flags().String("positioning", "", "Filter by positioning (primary, secondary)")
	benzingaGuidanceCmd.Flags().String("fiscal-year", "", "Filter by fiscal year (e.g., 2026)")
	benzingaGuidanceCmd.Flags().String("fiscal-period", "", "Filter by fiscal period (Q1, Q2, Q3, Q4)")
	benzingaGuidanceCmd.Flags().String("importance", "", "Filter by importance level (0-5)")
	benzingaGuidanceCmd.Flags().String("limit", "10", "Number of results to return (max 50000)")
	benzingaGuidanceCmd.Flags().String("sort", "date.desc", "Sort order (e.g., date.asc, date.desc)")

	// Benzinga Analysts flags
	benzingaAnalystsCmd.Flags().String("benzinga-id", "", "Filter by Benzinga analyst ID")
	benzingaAnalystsCmd.Flags().String("benzinga-firm-id", "", "Filter by Benzinga firm ID")
	benzingaAnalystsCmd.Flags().String("full-name", "", "Filter by analyst full name")
	benzingaAnalystsCmd.Flags().String("firm-name", "", "Filter by firm name")
	benzingaAnalystsCmd.Flags().String("limit", "10", "Number of results to return (max 50000)")
	benzingaAnalystsCmd.Flags().String("sort", "", "Sort order (comma-separated columns with .asc/.desc)")
}
