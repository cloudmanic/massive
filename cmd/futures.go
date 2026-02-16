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

// futuresCmd is the parent command for all futures market data subcommands
// including bars, contracts, products, schedules, exchanges, snapshot,
// trades, and quotes.
var futuresCmd = &cobra.Command{
	Use:   "futures",
	Short: "Futures market data commands",
}

// futuresBarsCmd retrieves aggregate bar data for a specific futures ticker
// with configurable resolution, time window, sorting, and result limits.
// Usage: massive futures bars ESM5 --resolution 1day --window-start 2025-03-01 --limit 10
var futuresBarsCmd = &cobra.Command{
	Use:   "bars [ticker]",
	Short: "Get aggregate bars for a futures ticker",
	Long:  "Retrieve aggregate bar data for a futures ticker including open, high, low, close, volume, settlement price, and dollar volume over configurable time windows.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		resolution, _ := cmd.Flags().GetString("resolution")
		windowStart, _ := cmd.Flags().GetString("window-start")
		windowStartGte, _ := cmd.Flags().GetString("window-start-gte")
		windowStartGt, _ := cmd.Flags().GetString("window-start-gt")
		windowStartLte, _ := cmd.Flags().GetString("window-start-lte")
		windowStartLt, _ := cmd.Flags().GetString("window-start-lt")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.FuturesAggParams{
			Resolution:     resolution,
			WindowStart:    windowStart,
			WindowStartGte: windowStartGte,
			WindowStartGt:  windowStartGt,
			WindowStartLte: windowStartLte,
			WindowStartLt:  windowStartLt,
			Limit:          limit,
			Sort:           sort,
		}

		result, err := client.GetFuturesAggs(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Ticker: %s | Bars: %d\n\n", ticker, len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "WINDOW START\tOPEN\tHIGH\tLOW\tCLOSE\tVOLUME\tSETTLEMENT\tTRANSACTIONS")
		fmt.Fprintln(w, "------------\t----\t----\t---\t-----\t------\t----------\t------------")

		for _, bar := range result.Results {
			t := time.Unix(0, bar.WindowStart)
			fmt.Fprintf(w, "%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\t%.4f\t%d\n",
				t.Format("2006-01-02 15:04:05"),
				bar.Open, bar.High, bar.Low, bar.Close,
				bar.Volume, bar.SettlementPrice, bar.Transactions)
		}
		w.Flush()

		return nil
	},
}

// futuresContractsCmd retrieves a list of futures contracts matching the
// provided filter criteria. Supports filtering by product code, ticker,
// active status, type, and date ranges.
// Usage: massive futures contracts --product-code ES --active true --limit 10
var futuresContractsCmd = &cobra.Command{
	Use:   "contracts",
	Short: "List futures contracts",
	Long:  "Retrieve a list of futures contracts with optional filtering by product code, ticker, active status, type, and trade date ranges.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		productCode, _ := cmd.Flags().GetString("product-code")
		ticker, _ := cmd.Flags().GetString("ticker")
		active, _ := cmd.Flags().GetString("active")
		contractType, _ := cmd.Flags().GetString("type")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.FuturesContractsParams{
			ProductCode: productCode,
			Ticker:      ticker,
			Active:      active,
			Type:        contractType,
			Limit:       limit,
			Sort:        sort,
		}

		result, err := client.GetFuturesContracts(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Contracts: %d\n\n", len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tNAME\tPRODUCT\tVENUE\tTYPE\tACTIVE\tDAYS TO MAT\tSETTLEMENT DATE")
		fmt.Fprintln(w, "------\t----\t-------\t-----\t----\t------\t-----------\t---------------")

		for _, c := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%v\t%d\t%s\n",
				c.Ticker, c.Name, c.ProductCode, c.TradingVenue,
				c.Type, c.Active, c.DaysToMaturity, c.SettlementDate)
		}
		w.Flush()

		return nil
	},
}

// futuresProductsCmd retrieves a list of futures products matching the
// provided filter criteria. Supports filtering by name, product code,
// sector, asset class, trading venue, and type.
// Usage: massive futures products --sector index --asset-class equity_index
var futuresProductsCmd = &cobra.Command{
	Use:   "products",
	Short: "List futures products",
	Long:  "Retrieve a list of futures products with optional filtering by name, product code, sector, asset class, trading venue, and type.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		productCode, _ := cmd.Flags().GetString("product-code")
		name, _ := cmd.Flags().GetString("name")
		sector, _ := cmd.Flags().GetString("sector")
		assetClass, _ := cmd.Flags().GetString("asset-class")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.FuturesProductsParams{
			ProductCode: productCode,
			Name:        name,
			Sector:      sector,
			AssetClass:  assetClass,
			Limit:       limit,
			Sort:        sort,
		}

		result, err := client.GetFuturesProducts(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Products: %d\n\n", len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "CODE\tNAME\tSECTOR\tASSET CLASS\tVENUE\tTYPE\tSETTLEMENT")
		fmt.Fprintln(w, "----\t----\t------\t-----------\t-----\t----\t----------")

		for _, p := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				p.ProductCode, p.Name, p.Sector, p.AssetClass,
				p.TradingVenue, p.Type, p.SettlementMethod)
		}
		w.Flush()

		return nil
	},
}

// futuresSchedulesCmd retrieves a list of futures schedule events matching
// the provided filters. Supports filtering by product code, session end
// date, and trading venue.
// Usage: massive futures schedules --product-code ES --session-end-date 2025-03-15
var futuresSchedulesCmd = &cobra.Command{
	Use:   "schedules",
	Short: "List futures schedule events",
	Long:  "Retrieve a list of futures schedule events such as settlements, last trade dates, and session boundaries with optional filtering by product code, session end date, and trading venue.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		productCode, _ := cmd.Flags().GetString("product-code")
		sessionEndDate, _ := cmd.Flags().GetString("session-end-date")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.FuturesSchedulesParams{
			ProductCode:    productCode,
			SessionEndDate: sessionEndDate,
			Limit:          limit,
			Sort:           sort,
		}

		result, err := client.GetFuturesSchedules(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Schedules: %d\n\n", len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "EVENT\tPRODUCT CODE\tPRODUCT NAME\tSESSION END\tTIMESTAMP\tVENUE")
		fmt.Fprintln(w, "-----\t------------\t------------\t-----------\t---------\t-----")

		for _, s := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				s.Event, s.ProductCode, s.ProductName,
				s.SessionEndDate, s.Timestamp, s.TradingVenue)
		}
		w.Flush()

		return nil
	},
}

// futuresExchangesCmd retrieves a list of known futures exchanges with
// their identifiers, names, MIC codes, and metadata.
// Usage: massive futures exchanges --limit 10
var futuresExchangesCmd = &cobra.Command{
	Use:   "exchanges",
	Short: "List futures exchanges",
	Long:  "Retrieve a list of known futures exchanges including their identifiers, MIC codes, names, and other reference attributes.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		limit, _ := cmd.Flags().GetString("limit")

		params := api.FuturesExchangesParams{
			Limit: limit,
		}

		result, err := client.GetFuturesExchanges(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Exchanges: %d\n\n", result.Count)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tACRONYM\tMIC\tTYPE\tLOCALE\tURL")
		fmt.Fprintln(w, "--\t----\t-------\t---\t----\t------\t---")

		for _, e := range result.Results {
			acronym := e.Acronym
			if acronym == "" {
				acronym = "-"
			}
			mic := e.MIC
			if mic == "" {
				mic = "-"
			}
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%s\n",
				e.ID, e.Name, acronym, mic, e.Type, e.Locale, e.URL)
		}
		w.Flush()

		return nil
	},
}

// futuresSnapshotCmd retrieves snapshot data for futures contracts including
// nested details, last minute bar, last quote, last trade, and session data.
// Supports filtering by product code and ticker.
// Usage: massive futures snapshot --product-code ES --ticker ESM5
var futuresSnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Get futures contract snapshots",
	Long:  "Retrieve snapshot data for futures contracts including open interest, last minute bar, last quote, last trade, and current session OHLC data.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		productCode, _ := cmd.Flags().GetString("product-code")
		ticker, _ := cmd.Flags().GetString("ticker")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.FuturesSnapshotParams{
			ProductCode: productCode,
			Ticker:      ticker,
			Limit:       limit,
			Sort:        sort,
		}

		result, err := client.GetFuturesSnapshot(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Snapshots: %d\n\n", result.Count)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tPRODUCT\tLAST PRICE\tBID\tASK\tSESS OPEN\tSESS HIGH\tSESS LOW\tSESS CLOSE\tCHANGE\tVOLUME")
		fmt.Fprintln(w, "------\t-------\t----------\t---\t---\t---------\t---------\t--------\t----------\t------\t------")

		for _, snap := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\n",
				snap.Ticker, snap.ProductCode,
				snap.LastTrade.Price,
				snap.LastQuote.BidPrice, snap.LastQuote.AskPrice,
				snap.Session.Open, snap.Session.High, snap.Session.Low,
				snap.Session.Close, snap.Session.Change, snap.Session.Volume)
		}
		w.Flush()

		return nil
	},
}

// futuresTradesCmd retrieves tick-level trade data for a specific futures
// ticker with optional session date filtering, sorting, and pagination.
// Usage: massive futures trades ESM5 --session-end-date 2025-03-15 --limit 10
var futuresTradesCmd = &cobra.Command{
	Use:   "trades [ticker]",
	Short: "Get tick-level trade data for a futures ticker",
	Long:  "Retrieve tick-level trade data for a futures ticker including price, size, sequence numbers, and nanosecond timestamps.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		sessionEndDate, _ := cmd.Flags().GetString("session-end-date")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.FuturesTradesParams{
			SessionEndDate: sessionEndDate,
			Limit:          limit,
			Sort:           sort,
		}

		result, err := client.GetFuturesTrades(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Ticker: %s | Trades: %d\n\n", ticker, len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TIMESTAMP\tPRICE\tSIZE\tSESSION END\tSEQUENCE")
		fmt.Fprintln(w, "---------\t-----\t----\t-----------\t--------")

		for _, trade := range result.Results {
			t := time.Unix(0, trade.Timestamp)
			fmt.Fprintf(w, "%s\t%.4f\t%.0f\t%s\t%d\n",
				t.Format("2006-01-02 15:04:05.000"),
				trade.Price, trade.Size, trade.SessionEndDate, trade.SequenceNumber)
		}
		w.Flush()

		return nil
	},
}

// futuresQuotesCmd retrieves tick-level quote data for a specific futures
// ticker with optional session date filtering, sorting, and pagination.
// Each quote includes bid/ask prices, sizes, and nanosecond timestamps.
// Usage: massive futures quotes ESM5 --session-end-date 2025-03-15 --limit 10
var futuresQuotesCmd = &cobra.Command{
	Use:   "quotes [ticker]",
	Short: "Get tick-level quote data for a futures ticker",
	Long:  "Retrieve tick-level quote data for a futures ticker including bid/ask prices, sizes, and nanosecond timestamps.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		sessionEndDate, _ := cmd.Flags().GetString("session-end-date")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.FuturesQuotesParams{
			SessionEndDate: sessionEndDate,
			Limit:          limit,
			Sort:           sort,
		}

		result, err := client.GetFuturesQuotes(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Ticker: %s | Quotes: %d\n\n", ticker, len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TIMESTAMP\tBID PRICE\tBID SIZE\tASK PRICE\tASK SIZE\tSESSION END")
		fmt.Fprintln(w, "---------\t---------\t--------\t---------\t--------\t-----------")

		for _, quote := range result.Results {
			t := time.Unix(0, quote.Timestamp)
			fmt.Fprintf(w, "%s\t%.4f\t%.0f\t%.4f\t%.0f\t%s\n",
				t.Format("2006-01-02 15:04:05.000"),
				quote.BidPrice, quote.BidSize,
				quote.AskPrice, quote.AskSize,
				quote.SessionEndDate)
		}
		w.Flush()

		return nil
	},
}

// init registers the futures parent command and all subcommands with their
// respective flags under the root command.
func init() {
	// Bars command flags
	futuresBarsCmd.Flags().String("resolution", "1day", "Bar resolution (1min, 15mins, 1hr, 1day)")
	futuresBarsCmd.Flags().String("window-start", "", "Filter by window start date or timestamp")
	futuresBarsCmd.Flags().String("window-start-gte", "", "Window start greater than or equal to")
	futuresBarsCmd.Flags().String("window-start-gt", "", "Window start greater than")
	futuresBarsCmd.Flags().String("window-start-lte", "", "Window start less than or equal to")
	futuresBarsCmd.Flags().String("window-start-lt", "", "Window start less than")
	futuresBarsCmd.Flags().String("limit", "5000", "Max number of results")
	futuresBarsCmd.Flags().String("sort", "asc", "Sort order (asc/desc)")

	// Contracts command flags
	futuresContractsCmd.Flags().String("product-code", "", "Filter by product code (e.g., ES, NQ, CL)")
	futuresContractsCmd.Flags().String("ticker", "", "Filter by specific ticker symbol")
	futuresContractsCmd.Flags().String("active", "", "Filter by active status (true/false)")
	futuresContractsCmd.Flags().String("type", "", "Filter by contract type")
	futuresContractsCmd.Flags().String("limit", "20", "Max number of results")
	futuresContractsCmd.Flags().String("sort", "", "Sort field")

	// Products command flags
	futuresProductsCmd.Flags().String("product-code", "", "Filter by product code (e.g., ES, NQ, CL)")
	futuresProductsCmd.Flags().String("name", "", "Filter by product name")
	futuresProductsCmd.Flags().String("sector", "", "Filter by sector")
	futuresProductsCmd.Flags().String("asset-class", "", "Filter by asset class")
	futuresProductsCmd.Flags().String("limit", "20", "Max number of results")
	futuresProductsCmd.Flags().String("sort", "", "Sort field")

	// Schedules command flags
	futuresSchedulesCmd.Flags().String("product-code", "", "Filter by product code (e.g., ES, NQ, CL)")
	futuresSchedulesCmd.Flags().String("session-end-date", "", "Filter by session end date (YYYY-MM-DD)")
	futuresSchedulesCmd.Flags().String("limit", "20", "Max number of results")
	futuresSchedulesCmd.Flags().String("sort", "", "Sort field")

	// Exchanges command flags
	futuresExchangesCmd.Flags().String("limit", "20", "Max number of results")

	// Snapshot command flags
	futuresSnapshotCmd.Flags().String("product-code", "", "Filter by product code (e.g., ES, NQ, CL)")
	futuresSnapshotCmd.Flags().String("ticker", "", "Filter by specific ticker symbol")
	futuresSnapshotCmd.Flags().String("limit", "20", "Max number of results")
	futuresSnapshotCmd.Flags().String("sort", "", "Sort field")

	// Trades command flags
	futuresTradesCmd.Flags().String("session-end-date", "", "Filter by session end date (YYYY-MM-DD)")
	futuresTradesCmd.Flags().String("limit", "1000", "Max number of results")
	futuresTradesCmd.Flags().String("sort", "", "Sort field (e.g., timestamp)")

	// Quotes command flags
	futuresQuotesCmd.Flags().String("session-end-date", "", "Filter by session end date (YYYY-MM-DD)")
	futuresQuotesCmd.Flags().String("limit", "1000", "Max number of results")
	futuresQuotesCmd.Flags().String("sort", "", "Sort field (e.g., timestamp)")

	// Register all subcommands under the futures parent
	futuresCmd.AddCommand(futuresBarsCmd)
	futuresCmd.AddCommand(futuresContractsCmd)
	futuresCmd.AddCommand(futuresProductsCmd)
	futuresCmd.AddCommand(futuresSchedulesCmd)
	futuresCmd.AddCommand(futuresExchangesCmd)
	futuresCmd.AddCommand(futuresSnapshotCmd)
	futuresCmd.AddCommand(futuresTradesCmd)
	futuresCmd.AddCommand(futuresQuotesCmd)

	// Register the futures parent under root
	rootCmd.AddCommand(futuresCmd)
}
