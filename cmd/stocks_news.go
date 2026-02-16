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

// stocksNewsCmd retrieves news articles for stocks from the Massive API.
// Supports filtering by ticker symbol, publication date range, and
// sorting. Results can be displayed as a table or raw JSON.
// Usage: massive stocks news --ticker AAPL --limit 5
var stocksNewsCmd = &cobra.Command{
	Use:   "news",
	Short: "Get stock market news articles",
	Long:  "Retrieve stock market news articles with optional filtering by ticker symbol, publication date range, sort order, and result limit.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker, _ := cmd.Flags().GetString("ticker")
		publishedUTC, _ := cmd.Flags().GetString("published-utc")
		publishedFrom, _ := cmd.Flags().GetString("published-from")
		publishedTo, _ := cmd.Flags().GetString("published-to")
		order, _ := cmd.Flags().GetString("order")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.NewsParams{
			Ticker:          strings.ToUpper(ticker),
			PublishedUTC:    publishedUTC,
			PublishedUTCGte: publishedFrom,
			PublishedUTCLte: publishedTo,
			Order:           order,
			Limit:           limit,
			Sort:            sort,
		}

		result, err := client.GetNews(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		// Display results count header
		fmt.Printf("News Articles: %d\n\n", result.Count)

		if len(result.Results) == 0 {
			fmt.Println("No news articles found.")
			return nil
		}

		// Print each news article in a readable table format
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "DATE\tSOURCE\tTICKERS\tTITLE")
		fmt.Fprintln(w, "----\t------\t-------\t-----")

		for _, article := range result.Results {
			// Format the published date to just the date portion
			date := formatPublishedDate(article.PublishedUTC)
			tickers := truncateString(strings.Join(article.Tickers, ","), 20)
			title := truncateString(article.Title, 60)

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				date, article.Publisher.Name, tickers, title)
		}
		w.Flush()

		return nil
	},
}

// formatPublishedDate extracts the date portion from an RFC3339
// timestamp string. If the string is shorter than 10 characters,
// it returns the original string unchanged.
func formatPublishedDate(utc string) string {
	if len(utc) >= 10 {
		return utc[:10]
	}
	return utc
}

// truncateString shortens a string to the specified maximum length,
// appending "..." if truncation occurs. Returns the original string
// if it is within the limit.
func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// init registers the news command and its flags under the stocks parent command.
func init() {
	stocksNewsCmd.Flags().String("ticker", "", "Filter by ticker symbol (e.g., AAPL)")
	stocksNewsCmd.Flags().String("published-utc", "", "Filter by exact publication date (YYYY-MM-DD)")
	stocksNewsCmd.Flags().String("published-from", "", "Filter articles published on or after this date (YYYY-MM-DD)")
	stocksNewsCmd.Flags().String("published-to", "", "Filter articles published on or before this date (YYYY-MM-DD)")
	stocksNewsCmd.Flags().String("order", "desc", "Sort order (asc/desc)")
	stocksNewsCmd.Flags().String("limit", "10", "Number of results to return (max 1000)")
	stocksNewsCmd.Flags().String("sort", "published_utc", "Sort field (published_utc)")
	stocksCmd.AddCommand(stocksNewsCmd)
}
