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

// tmxCmd is the parent command for all TMX partner data subcommands
// including corporate events from the Wall Street Horizon data feed.
var tmxCmd = &cobra.Command{
	Use:   "tmx",
	Short: "TMX partner data commands",
	Long:  "Access TMX/Wall Street Horizon partner data including corporate events such as earnings announcements, dividend dates, investor conferences, and stock splits.",
}

// tmxCorporateEventsCmd retrieves corporate events from the TMX/Wall
// Street Horizon data feed. Supports filtering by ticker, date range,
// event type, status, ISIN, trading venue, and TMX-specific identifiers.
// Output can be formatted as a table or JSON.
// Usage: massive tmx corporate-events --ticker AAPL --type earnings_announcement_date
var tmxCorporateEventsCmd = &cobra.Command{
	Use:   "corporate-events",
	Short: "Get corporate events from TMX/Wall Street Horizon",
	Long:  "Retrieve structured corporate event data from Wall Street Horizon's comprehensive global events calendar, including earnings announcements, dividend dates, investor conferences, and stock splits.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker, _ := cmd.Flags().GetString("ticker")
		ticker = strings.ToUpper(ticker)
		tickerAnyOf, _ := cmd.Flags().GetString("ticker-any-of")
		tickerAnyOf = strings.ToUpper(tickerAnyOf)
		date, _ := cmd.Flags().GetString("date")
		dateAnyOf, _ := cmd.Flags().GetString("date-any-of")
		dateGT, _ := cmd.Flags().GetString("date-gt")
		dateGTE, _ := cmd.Flags().GetString("date-gte")
		dateLT, _ := cmd.Flags().GetString("date-lt")
		dateLTE, _ := cmd.Flags().GetString("date-lte")
		eventType, _ := cmd.Flags().GetString("type")
		typeAnyOf, _ := cmd.Flags().GetString("type-any-of")
		status, _ := cmd.Flags().GetString("status")
		statusAnyOf, _ := cmd.Flags().GetString("status-any-of")
		isin, _ := cmd.Flags().GetString("isin")
		isinAnyOf, _ := cmd.Flags().GetString("isin-any-of")
		tradingVenue, _ := cmd.Flags().GetString("trading-venue")
		tradingVenueAnyOf, _ := cmd.Flags().GetString("trading-venue-any-of")
		tmxCompanyID, _ := cmd.Flags().GetString("tmx-company-id")
		tmxRecordID, _ := cmd.Flags().GetString("tmx-record-id")
		tmxRecordIDAnyOf, _ := cmd.Flags().GetString("tmx-record-id-any-of")
		sort, _ := cmd.Flags().GetString("sort")
		limit, _ := cmd.Flags().GetString("limit")

		params := api.TMXCorporateEventsParams{
			Ticker:           ticker,
			TickerAnyOf:      tickerAnyOf,
			Date:             date,
			DateAnyOf:        dateAnyOf,
			DateGT:           dateGT,
			DateGTE:          dateGTE,
			DateLT:           dateLT,
			DateLTE:          dateLTE,
			Type:             eventType,
			TypeAnyOf:        typeAnyOf,
			Status:           status,
			StatusAnyOf:      statusAnyOf,
			ISIN:             isin,
			ISINAnyOf:        isinAnyOf,
			TradingVenue:     tradingVenue,
			TradingVenueAnyOf: tradingVenueAnyOf,
			TMXCompanyID:     tmxCompanyID,
			TMXRecordID:      tmxRecordID,
			TMXRecordIDAnyOf: tmxRecordIDAnyOf,
			Sort:             sort,
			Limit:            limit,
		}

		result, err := client.GetTMXCorporateEvents(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Corporate Events: %d result(s)\n\n", len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tDATE\tTYPE\tNAME\tSTATUS\tCOMPANY\tVENUE")
		fmt.Fprintln(w, "------\t----\t----\t----\t------\t-------\t-----")

		for _, e := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				e.Ticker, e.Date, e.Type, e.Name,
				e.Status, e.CompanyName, e.TradingVenue)
		}
		w.Flush()

		return nil
	},
}

// init registers the tmx command as a subcommand of the root command
// and registers all TMX subcommands with their flags.
func init() {
	// Corporate events command flags
	tmxCorporateEventsCmd.Flags().String("ticker", "", "Stock ticker symbol (e.g. AAPL)")
	tmxCorporateEventsCmd.Flags().String("ticker-any-of", "", "Match any of the specified tickers (comma-separated, e.g. AAPL,MSFT,TSLA)")
	tmxCorporateEventsCmd.Flags().String("date", "", "Exact event date (YYYY-MM-DD)")
	tmxCorporateEventsCmd.Flags().String("date-any-of", "", "Match any of the specified dates (comma-separated)")
	tmxCorporateEventsCmd.Flags().String("date-gt", "", "Event date greater than (YYYY-MM-DD)")
	tmxCorporateEventsCmd.Flags().String("date-gte", "", "Event date greater than or equal (YYYY-MM-DD)")
	tmxCorporateEventsCmd.Flags().String("date-lt", "", "Event date less than (YYYY-MM-DD)")
	tmxCorporateEventsCmd.Flags().String("date-lte", "", "Event date less than or equal (YYYY-MM-DD)")
	tmxCorporateEventsCmd.Flags().String("type", "", "Event type (e.g. earnings_announcement_date, dividend, stock_split, investor_conference)")
	tmxCorporateEventsCmd.Flags().String("type-any-of", "", "Match any of the specified event types (comma-separated)")
	tmxCorporateEventsCmd.Flags().String("status", "", "Event status (approved, canceled, confirmed, historical, pending_approval, postponed, unconfirmed)")
	tmxCorporateEventsCmd.Flags().String("status-any-of", "", "Match any of the specified statuses (comma-separated)")
	tmxCorporateEventsCmd.Flags().String("isin", "", "International Securities Identification Number")
	tmxCorporateEventsCmd.Flags().String("isin-any-of", "", "Match any of the specified ISINs (comma-separated)")
	tmxCorporateEventsCmd.Flags().String("trading-venue", "", "Market Identifier Code / MIC (e.g. XNAS, XNYS)")
	tmxCorporateEventsCmd.Flags().String("trading-venue-any-of", "", "Match any of the specified trading venues (comma-separated)")
	tmxCorporateEventsCmd.Flags().String("tmx-company-id", "", "TMX company numeric ID")
	tmxCorporateEventsCmd.Flags().String("tmx-record-id", "", "TMX event record identifier")
	tmxCorporateEventsCmd.Flags().String("tmx-record-id-any-of", "", "Match any of the specified TMX record IDs (comma-separated)")
	tmxCorporateEventsCmd.Flags().String("sort", "", "Sort field with direction (e.g. date.asc, date.desc)")
	tmxCorporateEventsCmd.Flags().String("limit", "100", "Max number of results (default 100, max 50000)")

	// Register subcommands under tmx parent
	tmxCmd.AddCommand(tmxCorporateEventsCmd)

	// Register tmx under root command
	rootCmd.AddCommand(tmxCmd)
}
