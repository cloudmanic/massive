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
	"time"

	"github.com/cloudmanic/massive-cli/internal/api"
	"github.com/spf13/cobra"
)

// stocksTradesCmd retrieves tick-level trade data for a specific stock ticker
// with optional timestamp filtering, sorting, and pagination. Each trade
// includes price, size, exchange, trade conditions, and nanosecond timestamps.
// Usage: massive stocks trades AAPL --timestamp 2025-01-06 --limit 10
var stocksTradesCmd = &cobra.Command{
	Use:   "trades [ticker]",
	Short: "Get tick-level trade data for a stock ticker",
	Long:  "Retrieve tick-level trade data for a stock ticker including price, size, exchange, conditions, and precise timestamps.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		timestamp, _ := cmd.Flags().GetString("timestamp")
		timestampGte, _ := cmd.Flags().GetString("timestamp-gte")
		timestampGt, _ := cmd.Flags().GetString("timestamp-gt")
		timestampLte, _ := cmd.Flags().GetString("timestamp-lte")
		timestampLt, _ := cmd.Flags().GetString("timestamp-lt")
		order, _ := cmd.Flags().GetString("order")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.TradesParams{
			Timestamp:    timestamp,
			TimestampGte: timestampGte,
			TimestampGt:  timestampGt,
			TimestampLte: timestampLte,
			TimestampLt:  timestampLt,
			Order:        order,
			Limit:        limit,
			Sort:         sort,
		}

		result, err := client.GetTrades(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Ticker: %s | Trades: %d\n\n", ticker, len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TIMESTAMP\tPRICE\tSIZE\tEXCHANGE\tTAPE\tID")
		fmt.Fprintln(w, "---------\t-----\t----\t--------\t----\t--")

		for _, trade := range result.Results {
			t := time.Unix(0, trade.SipTimestamp)
			fmt.Fprintf(w, "%s\t%.4f\t%.0f\t%d\t%d\t%s\n",
				t.Format("2006-01-02 15:04:05.000"),
				trade.Price, trade.Size, trade.Exchange, trade.Tape, trade.ID)
		}
		w.Flush()

		return nil
	},
}

// stocksLastTradeCmd retrieves the most recent trade for a specific stock
// ticker. Returns price, size, exchange, and timestamp information useful
// for monitoring current market activity.
// Usage: massive stocks last-trade AAPL
var stocksLastTradeCmd = &cobra.Command{
	Use:   "last-trade [ticker]",
	Short: "Get the most recent trade for a stock ticker",
	Long:  "Retrieve the last available trade for a stock ticker including price, size, exchange, and timestamp.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])

		result, err := client.GetLastTrade(ticker)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		trade := result.Results
		t := time.Unix(0, trade.SipTimestamp)

		fmt.Printf("Ticker:    %s\n", trade.Ticker)
		fmt.Printf("Price:     $%.4f\n", trade.Price)
		fmt.Printf("Size:      %.0f\n", trade.Size)
		fmt.Printf("Exchange:  %d\n", trade.Exchange)
		fmt.Printf("Tape:      %d\n", trade.Tape)
		fmt.Printf("Trade ID:  %s\n", trade.ID)
		fmt.Printf("Timestamp: %s\n", t.Format("2006-01-02 15:04:05.000"))

		return nil
	},
}

// stocksQuotesCmd retrieves tick-level NBBO quote data for a specific stock
// ticker with optional timestamp filtering, sorting, and pagination. Each
// quote includes bid/ask prices, sizes, exchange IDs, and nanosecond timestamps.
// Usage: massive stocks quotes AAPL --timestamp 2025-01-06 --limit 10
var stocksQuotesCmd = &cobra.Command{
	Use:   "quotes [ticker]",
	Short: "Get tick-level NBBO quote data for a stock ticker",
	Long:  "Retrieve tick-level NBBO (National Best Bid and Offer) quote data for a stock ticker including bid/ask prices, sizes, and exchange information.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		timestamp, _ := cmd.Flags().GetString("timestamp")
		timestampGte, _ := cmd.Flags().GetString("timestamp-gte")
		timestampGt, _ := cmd.Flags().GetString("timestamp-gt")
		timestampLte, _ := cmd.Flags().GetString("timestamp-lte")
		timestampLt, _ := cmd.Flags().GetString("timestamp-lt")
		order, _ := cmd.Flags().GetString("order")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.QuotesParams{
			Timestamp:    timestamp,
			TimestampGte: timestampGte,
			TimestampGt:  timestampGt,
			TimestampLte: timestampLte,
			TimestampLt:  timestampLt,
			Order:        order,
			Limit:        limit,
			Sort:         sort,
		}

		result, err := client.GetQuotes(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Ticker: %s | Quotes: %d\n\n", ticker, len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TIMESTAMP\tBID PRICE\tBID SIZE\tASK PRICE\tASK SIZE\tBID EX\tASK EX")
		fmt.Fprintln(w, "---------\t---------\t--------\t---------\t--------\t------\t------")

		for _, quote := range result.Results {
			t := time.Unix(0, quote.SipTimestamp)
			fmt.Fprintf(w, "%s\t%.4f\t%.0f\t%.4f\t%.0f\t%d\t%d\n",
				t.Format("2006-01-02 15:04:05.000"),
				quote.BidPrice, quote.BidSize,
				quote.AskPrice, quote.AskSize,
				quote.BidExchange, quote.AskExchange)
		}
		w.Flush()

		return nil
	},
}

// stocksLastQuoteCmd retrieves the most recent NBBO quote for a specific
// stock ticker. Returns the last available bid/ask prices, sizes, and
// exchange information for real-time market monitoring.
// Usage: massive stocks last-quote AAPL
var stocksLastQuoteCmd = &cobra.Command{
	Use:   "last-quote [ticker]",
	Short: "Get the most recent NBBO quote for a stock ticker",
	Long:  "Retrieve the last available NBBO (National Best Bid and Offer) quote for a stock ticker including bid/ask prices, sizes, and exchange information.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])

		result, err := client.GetLastQuote(ticker)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		quote := result.Results
		t := time.Unix(0, quote.SipTimestamp)

		fmt.Printf("Ticker:      %s\n", quote.Ticker)
		fmt.Printf("Bid Price:   $%.4f\n", quote.BidPrice)
		fmt.Printf("Bid Size:    %d\n", quote.BidSize)
		fmt.Printf("Bid Exchange: %d\n", quote.BidExchange)
		fmt.Printf("Ask Price:   $%.4f\n", quote.AskPrice)
		fmt.Printf("Ask Size:    %d\n", quote.AskSize)
		fmt.Printf("Ask Exchange: %d\n", quote.AskExchange)
		fmt.Printf("Tape:        %d\n", quote.Tape)
		fmt.Printf("Timestamp:   %s\n", t.Format("2006-01-02 15:04:05.000"))

		return nil
	},
}

// init registers the trades, last-trade, quotes, and last-quote commands
// and their flags under the stocks parent command.
func init() {
	// Trades command flags
	stocksTradesCmd.Flags().String("timestamp", "", "Filter by date (YYYY-MM-DD) or nanosecond timestamp")
	stocksTradesCmd.Flags().String("timestamp-gte", "", "Timestamp greater than or equal to")
	stocksTradesCmd.Flags().String("timestamp-gt", "", "Timestamp greater than")
	stocksTradesCmd.Flags().String("timestamp-lte", "", "Timestamp less than or equal to")
	stocksTradesCmd.Flags().String("timestamp-lt", "", "Timestamp less than")
	stocksTradesCmd.Flags().String("order", "", "Sort order (asc/desc)")
	stocksTradesCmd.Flags().String("limit", "1000", "Max number of results (max 50000)")
	stocksTradesCmd.Flags().String("sort", "", "Sort field (e.g., timestamp)")

	// Quotes command flags
	stocksQuotesCmd.Flags().String("timestamp", "", "Filter by date (YYYY-MM-DD) or nanosecond timestamp")
	stocksQuotesCmd.Flags().String("timestamp-gte", "", "Timestamp greater than or equal to")
	stocksQuotesCmd.Flags().String("timestamp-gt", "", "Timestamp greater than")
	stocksQuotesCmd.Flags().String("timestamp-lte", "", "Timestamp less than or equal to")
	stocksQuotesCmd.Flags().String("timestamp-lt", "", "Timestamp less than")
	stocksQuotesCmd.Flags().String("order", "", "Sort order (asc/desc)")
	stocksQuotesCmd.Flags().String("limit", "1000", "Max number of results (max 50000)")
	stocksQuotesCmd.Flags().String("sort", "", "Sort field (e.g., timestamp)")

	// Register all four commands under the stocks parent
	stocksCmd.AddCommand(stocksTradesCmd)
	stocksCmd.AddCommand(stocksLastTradeCmd)
	stocksCmd.AddCommand(stocksQuotesCmd)
	stocksCmd.AddCommand(stocksLastQuoteCmd)
}
