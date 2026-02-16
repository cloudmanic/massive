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

// forexCmd is the parent command for all forex market data subcommands
// including bars, conversion, quotes, snapshots, and technical indicators.
var forexCmd = &cobra.Command{
	Use:   "forex",
	Short: "Forex market data commands",
}

// --- Aggregates ---

// forexBarsCmd retrieves custom OHLC aggregate bars for a forex ticker
// over a specified time range. Supports configurable timespan, multiplier,
// sort order, and result limit.
// Usage: massive forex bars C:EURUSD --from 2024-01-01 --to 2024-01-31
var forexBarsCmd = &cobra.Command{
	Use:   "bars [ticker]",
	Short: "Get OHLC aggregate bars for a forex ticker",
	Long:  "Retrieve custom OHLC (Open, High, Low, Close) aggregate bar data for a forex ticker over a specified time range.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		multiplier, _ := cmd.Flags().GetString("multiplier")
		timespan, _ := cmd.Flags().GetString("timespan")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		adjusted, _ := cmd.Flags().GetString("adjusted")
		sort, _ := cmd.Flags().GetString("sort")
		limit, _ := cmd.Flags().GetString("limit")

		params := api.ForexBarsParams{
			Multiplier: multiplier,
			Timespan:   timespan,
			From:       from,
			To:         to,
			Adjusted:   adjusted,
			Sort:       sort,
			Limit:      limit,
		}

		result, err := client.GetForexBars(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Ticker: %s | Bars: %d | Adjusted: %v\n\n", result.Ticker, result.ResultsCount, result.Adjusted)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "DATE\tOPEN\tHIGH\tLOW\tCLOSE\tVOLUME\tVWAP\tTRADES")
		fmt.Fprintln(w, "----\t----\t----\t---\t-----\t------\t----\t------")

		for _, bar := range result.Results {
			t := time.UnixMilli(bar.Timestamp)
			fmt.Fprintf(w, "%s\t%.6f\t%.6f\t%.6f\t%.6f\t%.0f\t%.6f\t%d\n",
				t.Format("2006-01-02"),
				bar.Open, bar.High, bar.Low, bar.Close,
				bar.Volume, bar.VWAP, bar.NumTrades)
		}
		w.Flush()

		return nil
	},
}

// forexDailyMarketSummaryCmd retrieves the grouped daily OHLC summary for
// all forex tickers on a specified date. Useful for broad forex market
// analysis and screening.
// Usage: massive forex daily-market-summary 2024-01-09
var forexDailyMarketSummaryCmd = &cobra.Command{
	Use:   "daily-market-summary [date]",
	Short: "Get daily market summary for all forex tickers",
	Long:  "Retrieve the daily OHLC, volume, and VWAP data for all forex tickers on a specified trading date.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		date := args[0]
		adjusted, _ := cmd.Flags().GetString("adjusted")

		params := api.ForexMarketSummaryParams{
			Adjusted: adjusted,
		}

		result, err := client.GetForexDailyMarketSummary(date, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Date: %s | Tickers: %d | Adjusted: %v\n\n", date, result.ResultsCount, result.Adjusted)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tOPEN\tHIGH\tLOW\tCLOSE\tVOLUME\tVWAP\tTRADES")
		fmt.Fprintln(w, "------\t----\t----\t---\t-----\t------\t----\t------")

		for _, s := range result.Results {
			fmt.Fprintf(w, "%s\t%.6f\t%.6f\t%.6f\t%.6f\t%.0f\t%.6f\t%d\n",
				s.Ticker, s.Open, s.High, s.Low, s.Close,
				s.Volume, s.VWAP, s.NumTrades)
		}
		w.Flush()

		return nil
	},
}

// forexPreviousDayBarCmd retrieves the previous day's OHLC bar data for
// a specific forex ticker. Useful for comparing current prices against
// the most recent close.
// Usage: massive forex previous-day-bar C:EURUSD
var forexPreviousDayBarCmd = &cobra.Command{
	Use:   "previous-day-bar [ticker]",
	Short: "Get previous day bar for a forex ticker",
	Long:  "Retrieve the previous day's OHLC aggregate bar data for a specific forex ticker.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		adjusted, _ := cmd.Flags().GetString("adjusted")

		result, err := client.GetForexPreviousDayBar(ticker, adjusted)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Ticker: %s | Adjusted: %v\n\n", result.Ticker, result.Adjusted)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "DATE\tOPEN\tHIGH\tLOW\tCLOSE\tVOLUME\tVWAP\tTRADES")
		fmt.Fprintln(w, "----\t----\t----\t---\t-----\t------\t----\t------")

		for _, bar := range result.Results {
			t := time.UnixMilli(bar.Timestamp)
			fmt.Fprintf(w, "%s\t%.6f\t%.6f\t%.6f\t%.6f\t%.0f\t%.6f\t%d\n",
				t.Format("2006-01-02"),
				bar.Open, bar.High, bar.Low, bar.Close,
				bar.Volume, bar.VWAP, bar.NumTrades)
		}
		w.Flush()

		return nil
	},
}

// --- Currency Conversion ---

// forexConvertCmd converts a specified amount from one currency to another
// using the latest exchange rate. The from and to currency codes are
// provided as positional arguments.
// Usage: massive forex convert USD EUR --amount 100 --precision 2
var forexConvertCmd = &cobra.Command{
	Use:   "convert [from] [to]",
	Short: "Convert between currencies",
	Long:  "Convert a specified amount from one currency to another using the latest forex exchange rate.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		from := strings.ToUpper(args[0])
		to := strings.ToUpper(args[1])
		amount, _ := cmd.Flags().GetString("amount")
		precision, _ := cmd.Flags().GetString("precision")

		params := api.ForexConversionParams{
			Amount:    amount,
			Precision: precision,
		}

		result, err := client.GetForexConversion(from, to, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Conversion: %s -> %s\n", result.From, result.To)
		fmt.Printf("Symbol: %s\n", result.Symbol)
		fmt.Printf("Initial Amount: %.2f\n", result.InitialAmount)
		fmt.Printf("Converted: %.6f\n", result.Converted)
		fmt.Printf("Ask: %.6f | Bid: %.6f\n", result.Last.Ask, result.Last.Bid)
		fmt.Printf("Exchange: %d\n", result.Last.Exchange)

		return nil
	},
}

// --- Quotes ---

// forexQuotesCmd retrieves tick-level quote data for a specific forex
// ticker with optional timestamp filtering, sorting, and pagination.
// Usage: massive forex quotes C:EURUSD --limit 10
var forexQuotesCmd = &cobra.Command{
	Use:   "quotes [ticker]",
	Short: "Get quotes for a forex ticker",
	Long:  "Retrieve tick-level quote data for a specific forex ticker with bid/ask prices, exchange IDs, and timestamps.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")
		order, _ := cmd.Flags().GetString("order")

		params := api.ForexQuotesParams{
			Limit: limit,
			Sort:  sort,
			Order: order,
		}

		result, err := client.GetForexQuotes(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Ticker: %s | Quotes: %d\n\n", ticker, len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TIMESTAMP\tASK PRICE\tBID PRICE\tASK EXCHANGE\tBID EXCHANGE")
		fmt.Fprintln(w, "---------\t---------\t---------\t------------\t------------")

		for _, q := range result.Results {
			t := time.UnixMilli(q.ParticipantTimestamp)
			fmt.Fprintf(w, "%s\t%.6f\t%.6f\t%d\t%d\n",
				t.Format("2006-01-02 15:04:05"),
				q.AskPrice, q.BidPrice, q.AskExchange, q.BidExchange)
		}
		w.Flush()

		return nil
	},
}

// forexLastQuoteCmd retrieves the most recent forex quote for a currency
// pair specified by the from and to currency codes.
// Usage: massive forex last-quote EUR USD
var forexLastQuoteCmd = &cobra.Command{
	Use:   "last-quote [from] [to]",
	Short: "Get last quote for a currency pair",
	Long:  "Retrieve the most recent forex quote for a currency pair including ask/bid prices, exchange, and timestamp.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		from := strings.ToUpper(args[0])
		to := strings.ToUpper(args[1])

		result, err := client.GetForexLastQuote(from, to)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Symbol: %s\n", result.Symbol)
		fmt.Printf("Ask: %.6f\n", result.Last.Ask)
		fmt.Printf("Bid: %.6f\n", result.Last.Bid)
		fmt.Printf("Exchange: %d\n", result.Last.Exchange)
		ts := time.UnixMilli(result.Last.Timestamp)
		fmt.Printf("Timestamp: %s\n", ts.Format("2006-01-02 15:04:05"))

		return nil
	},
}

// --- Snapshots ---

// forexSnapshotCmd retrieves the most recent snapshot for a single forex
// ticker, including the current day's bar, previous day's bar, and the
// last available quote data.
// Usage: massive forex snapshot C:EURUSD
var forexSnapshotCmd = &cobra.Command{
	Use:   "snapshot [ticker]",
	Short: "Get snapshot for a single forex ticker",
	Long:  "Retrieve the most recent snapshot for a single forex ticker including current day, previous day, last quote, and price change data.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])

		result, err := client.GetForexSnapshotTicker(ticker)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		t := result.Ticker
		fmt.Printf("Ticker: %s | Change: %.6f (%.2f%%)\n\n", t.Ticker, t.TodaysChange, t.TodaysChangePct)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "PERIOD\tOPEN\tHIGH\tLOW\tCLOSE")
		fmt.Fprintln(w, "------\t----\t----\t---\t-----")

		fmt.Fprintf(w, "Day\t%.6f\t%.6f\t%.6f\t%.6f\n",
			t.Day.Open, t.Day.High, t.Day.Low, t.Day.Close)

		fmt.Fprintf(w, "Prev Day\t%.6f\t%.6f\t%.6f\t%.6f\n",
			t.PrevDay.Open, t.PrevDay.High, t.PrevDay.Low, t.PrevDay.Close)

		w.Flush()

		fmt.Printf("\nLast Quote: Ask: %.6f | Bid: %.6f | Exchange: %d\n",
			t.LastQuote.Ask, t.LastQuote.Bid, t.LastQuote.Exchange)

		return nil
	},
}

// forexSnapshotMarketCmd retrieves snapshot data for all forex tickers or
// a filtered subset. Supports filtering by a comma-separated list of
// ticker symbols.
// Usage: massive forex snapshot-market --tickers C:EURUSD,C:GBPUSD
var forexSnapshotMarketCmd = &cobra.Command{
	Use:   "snapshot-market",
	Short: "Get snapshots for all or selected forex tickers",
	Long:  "Retrieve snapshot data for all forex tickers or a filtered subset specified by a comma-separated list of symbols.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		tickers, _ := cmd.Flags().GetString("tickers")

		params := api.ForexSnapshotAllParams{
			Tickers: tickers,
		}

		result, err := client.GetForexSnapshotAll(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Tickers: %d\n\n", result.Count)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tDAY OPEN\tDAY HIGH\tDAY LOW\tDAY CLOSE\tCHANGE\tCHANGE %")
		fmt.Fprintln(w, "------\t--------\t--------\t-------\t---------\t------\t--------")

		for _, t := range result.Tickers {
			fmt.Fprintf(w, "%s\t%.6f\t%.6f\t%.6f\t%.6f\t%.6f\t%.2f%%\n",
				t.Ticker, t.Day.Open, t.Day.High, t.Day.Low, t.Day.Close,
				t.TodaysChange, t.TodaysChangePct)
		}
		w.Flush()

		return nil
	},
}

// forexGainersCmd retrieves the current top forex gainers. Each ticker
// includes the current day's bar, previous day's bar, and percentage
// change values.
// Usage: massive forex gainers
var forexGainersCmd = &cobra.Command{
	Use:   "gainers",
	Short: "Get top gaining forex tickers",
	Long:  "Retrieve the current top gainers in the forex market with snapshot data including day bar, previous day bar, and change percentages.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		result, err := client.GetForexGainersLosers("gainers")
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		return printForexGainersLosersTable("Gainers", result)
	},
}

// forexLosersCmd retrieves the current top forex losers. Each ticker
// includes the current day's bar, previous day's bar, and percentage
// change values.
// Usage: massive forex losers
var forexLosersCmd = &cobra.Command{
	Use:   "losers",
	Short: "Get top losing forex tickers",
	Long:  "Retrieve the current top losers in the forex market with snapshot data including day bar, previous day bar, and change percentages.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		result, err := client.GetForexGainersLosers("losers")
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		return printForexGainersLosersTable("Losers", result)
	},
}

// printForexGainersLosersTable formats and prints a table of forex gainers
// or losers snapshot data to stdout. The title parameter labels the output
// as either "Gainers" or "Losers" for display clarity.
func printForexGainersLosersTable(title string, result *api.ForexSnapshotGainersLosersResponse) error {
	fmt.Printf("Top %s: %d tickers\n\n", title, len(result.Tickers))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TICKER\tDAY OPEN\tDAY HIGH\tDAY LOW\tDAY CLOSE\tCHANGE\tCHANGE %")
	fmt.Fprintln(w, "------\t--------\t--------\t-------\t---------\t------\t--------")

	for _, t := range result.Tickers {
		fmt.Fprintf(w, "%s\t%.6f\t%.6f\t%.6f\t%.6f\t%.6f\t%.2f%%\n",
			t.Ticker, t.Day.Open, t.Day.High, t.Day.Low, t.Day.Close,
			t.TodaysChange, t.TodaysChangePct)
	}
	w.Flush()

	return nil
}

// --- Technical Indicators ---

// forexSMACmd retrieves Simple Moving Average (SMA) data for a forex ticker
// over a specified time range. SMA smooths price data by calculating the
// arithmetic mean over a rolling window period.
// Usage: massive forex sma C:EURUSD --from 2025-01-06 --to 2025-01-10
var forexSMACmd = &cobra.Command{
	Use:   "sma [ticker]",
	Short: "Get Simple Moving Average (SMA) for a forex ticker",
	Long:  "Retrieve Simple Moving Average (SMA) indicator data for a forex ticker. SMA calculates the arithmetic mean of closing prices over a given window period.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		params := buildForexIndicatorParams(cmd)

		result, err := client.GetForexSMA(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		printForexIndicatorTable(ticker, "SMA", result)
		return nil
	},
}

// forexEMACmd retrieves Exponential Moving Average (EMA) data for a forex
// ticker over a specified time range. EMA places greater weight on recent
// prices compared to SMA for quicker trend detection.
// Usage: massive forex ema C:EURUSD --from 2025-01-06 --to 2025-01-10
var forexEMACmd = &cobra.Command{
	Use:   "ema [ticker]",
	Short: "Get Exponential Moving Average (EMA) for a forex ticker",
	Long:  "Retrieve Exponential Moving Average (EMA) indicator data for a forex ticker. EMA places greater weight on recent prices for more responsive trend signals.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		params := buildForexIndicatorParams(cmd)

		result, err := client.GetForexEMA(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		printForexIndicatorTable(ticker, "EMA", result)
		return nil
	},
}

// forexRSICmd retrieves Relative Strength Index (RSI) data for a forex
// ticker over a specified time range. RSI oscillates between 0 and 100
// to identify overbought or oversold conditions.
// Usage: massive forex rsi C:EURUSD --from 2025-01-06 --to 2025-01-10
var forexRSICmd = &cobra.Command{
	Use:   "rsi [ticker]",
	Short: "Get Relative Strength Index (RSI) for a forex ticker",
	Long:  "Retrieve Relative Strength Index (RSI) indicator data for a forex ticker. RSI measures the speed and magnitude of price changes, oscillating between 0 and 100.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		params := buildForexIndicatorParams(cmd)

		result, err := client.GetForexRSI(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		printForexIndicatorTable(ticker, "RSI", result)
		return nil
	},
}

// forexMACDCmd retrieves Moving Average Convergence/Divergence (MACD) data
// for a forex ticker over a specified time range. MACD is a momentum indicator
// with three components: the MACD line, signal line, and histogram.
// Usage: massive forex macd C:EURUSD --from 2025-01-06 --to 2025-01-10
var forexMACDCmd = &cobra.Command{
	Use:   "macd [ticker]",
	Short: "Get Moving Average Convergence/Divergence (MACD) for a forex ticker",
	Long:  "Retrieve MACD indicator data for a forex ticker. MACD is a momentum indicator showing the relationship between two EMAs, with signal line and histogram.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		timespan, _ := cmd.Flags().GetString("timespan")
		adjusted, _ := cmd.Flags().GetString("adjusted")
		shortWindow, _ := cmd.Flags().GetString("short-window")
		longWindow, _ := cmd.Flags().GetString("long-window")
		signalWindow, _ := cmd.Flags().GetString("signal-window")
		seriesType, _ := cmd.Flags().GetString("series-type")
		order, _ := cmd.Flags().GetString("order")
		limit, _ := cmd.Flags().GetString("limit")

		params := api.MACDParams{
			TimestampGTE: from,
			TimestampLTE: to,
			Timespan:     timespan,
			Adjusted:     adjusted,
			ShortWindow:  shortWindow,
			LongWindow:   longWindow,
			SignalWindow: signalWindow,
			SeriesType:   seriesType,
			Order:        order,
			Limit:        limit,
		}

		result, err := client.GetForexMACD(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		printForexMACDTable(ticker, result)
		return nil
	},
}

// buildForexIndicatorParams extracts the common indicator flags from the cobra
// command and returns a populated IndicatorParams struct. This is shared
// by the forex SMA, EMA, and RSI commands which all use the same parameters.
func buildForexIndicatorParams(cmd *cobra.Command) api.IndicatorParams {
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	timespan, _ := cmd.Flags().GetString("timespan")
	adjusted, _ := cmd.Flags().GetString("adjusted")
	window, _ := cmd.Flags().GetString("window")
	seriesType, _ := cmd.Flags().GetString("series-type")
	order, _ := cmd.Flags().GetString("order")
	limit, _ := cmd.Flags().GetString("limit")

	return api.IndicatorParams{
		TimestampGTE: from,
		TimestampLTE: to,
		Timespan:     timespan,
		Adjusted:     adjusted,
		Window:       window,
		SeriesType:   seriesType,
		Order:        order,
		Limit:        limit,
	}
}

// printForexIndicatorTable renders a formatted table of indicator values for
// the forex SMA, EMA, or RSI commands. Each row displays the date and value.
func printForexIndicatorTable(ticker, indicator string, result *api.IndicatorResponse) {
	fmt.Printf("Ticker: %s | Indicator: %s | Values: %d\n\n", ticker, indicator, len(result.Results.Values))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "DATE\tVALUE")
	fmt.Fprintln(w, "----\t-----")

	for _, v := range result.Results.Values {
		t := time.UnixMilli(v.Timestamp)
		fmt.Fprintf(w, "%s\t%.6f\n", t.Format("2006-01-02"), v.Value)
	}
	w.Flush()
}

// printForexMACDTable renders a formatted table of MACD indicator values
// including the MACD line, signal line, and histogram for each data point.
func printForexMACDTable(ticker string, result *api.MACDResponse) {
	fmt.Printf("Ticker: %s | Indicator: MACD | Values: %d\n\n", ticker, len(result.Results.Values))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "DATE\tMACD\tSIGNAL\tHISTOGRAM")
	fmt.Fprintln(w, "----\t----\t------\t---------")

	for _, v := range result.Results.Values {
		t := time.UnixMilli(v.Timestamp)
		fmt.Fprintf(w, "%s\t%.6f\t%.6f\t%.6f\n",
			t.Format("2006-01-02"), v.Value, v.Signal, v.Histogram)
	}
	w.Flush()
}

// addForexIndicatorFlags registers the common flags shared by the forex SMA,
// EMA, and RSI indicator subcommands. These include date range, window,
// timespan, series type, and pagination controls.
func addForexIndicatorFlags(cmd *cobra.Command, defaultWindow string) {
	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD) [required]")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD) [required]")
	cmd.Flags().String("timespan", "day", "Aggregate time window (minute, hour, day, week, month, quarter, year)")
	cmd.Flags().String("adjusted", "true", "Adjust for splits (true/false)")
	cmd.Flags().String("window", defaultWindow, "Number of periods for the indicator calculation")
	cmd.Flags().String("series-type", "close", "Price type for calculation (open, high, low, close)")
	cmd.Flags().String("order", "desc", "Sort order by timestamp (asc/desc)")
	cmd.Flags().String("limit", "10", "Max number of results (max 5000)")

	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")
}

// --- Tickers ---

// forexTickersCmd lists and searches forex tickers from the Massive
// reference data. Supports filtering by name search, active status,
// and pagination.
// Usage: massive forex tickers --search "EUR"
var forexTickersCmd = &cobra.Command{
	Use:   "tickers",
	Short: "List and search forex tickers",
	Long:  "Retrieve a list of forex tickers with optional filtering by name, active status, and pagination controls.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		search, _ := cmd.Flags().GetString("search")
		active, _ := cmd.Flags().GetString("active")
		sort, _ := cmd.Flags().GetString("sort")
		order, _ := cmd.Flags().GetString("order")
		limit, _ := cmd.Flags().GetString("limit")

		params := api.ForexTickerParams{
			Search: search,
			Active: active,
			Sort:   sort,
			Order:  order,
			Limit:  limit,
		}

		result, err := client.GetForexTickers(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Results: %d\n\n", result.Count)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tNAME\tMARKET\tACTIVE")
		fmt.Fprintln(w, "------\t----\t------\t------")

		for _, t := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%v\n",
				t.Ticker, t.Name, t.Market, t.Active)
		}
		w.Flush()

		return nil
	},
}

// forexTickerOverviewCmd retrieves detailed reference data for a specific
// forex ticker including market, locale, currency names, and active status.
// Usage: massive forex ticker-overview C:EURUSD
var forexTickerOverviewCmd = &cobra.Command{
	Use:   "ticker-overview [ticker]",
	Short: "Get detailed info for a forex ticker",
	Long:  "Retrieve detailed reference data for a specific forex ticker including market, locale, currency names, and active status.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])

		result, err := client.GetForexTickerOverview(ticker)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		r := result.Results
		fmt.Printf("Ticker: %s\n", r.Ticker)
		fmt.Printf("Name: %s\n", r.Name)
		fmt.Printf("Market: %s\n", r.Market)
		fmt.Printf("Locale: %s\n", r.Locale)
		fmt.Printf("Active: %v\n", r.Active)
		fmt.Printf("Currency: %s (%s)\n", r.CurrencyName, r.CurrencySymbol)
		fmt.Printf("Base Currency: %s (%s)\n", r.BaseCurrencyName, r.BaseCurrencySymbol)

		return nil
	},
}

// init registers the forex parent command and all forex subcommands with
// their respective flags under the root command.
func init() {
	// Bars flags
	forexBarsCmd.Flags().String("multiplier", "1", "Size of the timespan multiplier")
	forexBarsCmd.Flags().String("timespan", "day", "Timespan (minute, hour, day, week, month, quarter, year)")
	forexBarsCmd.Flags().String("from", "", "Start date (YYYY-MM-DD) [required]")
	forexBarsCmd.Flags().String("to", "", "End date (YYYY-MM-DD) [required]")
	forexBarsCmd.Flags().String("adjusted", "true", "Adjust for splits (true/false)")
	forexBarsCmd.Flags().String("sort", "asc", "Sort order (asc/desc)")
	forexBarsCmd.Flags().String("limit", "5000", "Max number of results (max 50000)")
	forexBarsCmd.MarkFlagRequired("from")
	forexBarsCmd.MarkFlagRequired("to")

	// Daily market summary flags
	forexDailyMarketSummaryCmd.Flags().String("adjusted", "true", "Adjust for splits (true/false)")

	// Previous day bar flags
	forexPreviousDayBarCmd.Flags().String("adjusted", "true", "Adjust for splits (true/false)")

	// Convert flags
	forexConvertCmd.Flags().String("amount", "1", "Amount to convert")
	forexConvertCmd.Flags().String("precision", "2", "Decimal precision for the converted amount")

	// Quotes flags
	forexQuotesCmd.Flags().String("limit", "10", "Max number of results")
	forexQuotesCmd.Flags().String("sort", "timestamp", "Sort field")
	forexQuotesCmd.Flags().String("order", "desc", "Sort order (asc/desc)")

	// Snapshot market flags
	forexSnapshotMarketCmd.Flags().String("tickers", "", "Comma-separated list of ticker symbols (default: all)")

	// SMA flags
	addForexIndicatorFlags(forexSMACmd, "10")

	// EMA flags
	addForexIndicatorFlags(forexEMACmd, "10")

	// RSI flags
	addForexIndicatorFlags(forexRSICmd, "14")

	// MACD flags
	forexMACDCmd.Flags().String("from", "", "Start date (YYYY-MM-DD) [required]")
	forexMACDCmd.Flags().String("to", "", "End date (YYYY-MM-DD) [required]")
	forexMACDCmd.Flags().String("timespan", "day", "Aggregate time window (minute, hour, day, week, month, quarter, year)")
	forexMACDCmd.Flags().String("adjusted", "true", "Adjust for splits (true/false)")
	forexMACDCmd.Flags().String("short-window", "12", "Short EMA period for MACD line")
	forexMACDCmd.Flags().String("long-window", "26", "Long EMA period for MACD line")
	forexMACDCmd.Flags().String("signal-window", "9", "Signal line EMA period")
	forexMACDCmd.Flags().String("series-type", "close", "Price type for calculation (open, high, low, close)")
	forexMACDCmd.Flags().String("order", "desc", "Sort order by timestamp (asc/desc)")
	forexMACDCmd.Flags().String("limit", "10", "Max number of results (max 5000)")
	forexMACDCmd.MarkFlagRequired("from")
	forexMACDCmd.MarkFlagRequired("to")

	// Tickers flags
	forexTickersCmd.Flags().String("search", "", "Search by currency pair name or symbol")
	forexTickersCmd.Flags().String("active", "", "Filter by active status (true/false)")
	forexTickersCmd.Flags().String("sort", "ticker", "Sort field (ticker, name)")
	forexTickersCmd.Flags().String("order", "asc", "Sort order (asc/desc)")
	forexTickersCmd.Flags().String("limit", "20", "Number of results to return (max 1000)")

	// Register all subcommands under forex
	forexCmd.AddCommand(forexBarsCmd)
	forexCmd.AddCommand(forexDailyMarketSummaryCmd)
	forexCmd.AddCommand(forexPreviousDayBarCmd)
	forexCmd.AddCommand(forexConvertCmd)
	forexCmd.AddCommand(forexQuotesCmd)
	forexCmd.AddCommand(forexLastQuoteCmd)
	forexCmd.AddCommand(forexSnapshotCmd)
	forexCmd.AddCommand(forexSnapshotMarketCmd)
	forexCmd.AddCommand(forexGainersCmd)
	forexCmd.AddCommand(forexLosersCmd)
	forexCmd.AddCommand(forexSMACmd)
	forexCmd.AddCommand(forexEMACmd)
	forexCmd.AddCommand(forexRSICmd)
	forexCmd.AddCommand(forexMACDCmd)
	forexCmd.AddCommand(forexTickersCmd)
	forexCmd.AddCommand(forexTickerOverviewCmd)

	// Register forex under root
	rootCmd.AddCommand(forexCmd)
}
