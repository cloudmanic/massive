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

// cryptoCmd is the parent command for all crypto market data subcommands
// including bars, snapshots, indicators, tickers, and trades.
var cryptoCmd = &cobra.Command{
	Use:   "crypto",
	Short: "Crypto market data commands",
}

// -------------------------------------------------------------------
// Aggregates Commands
// -------------------------------------------------------------------

// cryptoBarsCmd retrieves custom OHLC aggregate bars for a crypto ticker
// over a specified time range. Supports configurable timespan, multiplier,
// sort order, and result limit.
// Usage: massive crypto bars X:BTCUSD --from 2024-01-01 --to 2024-01-31
var cryptoBarsCmd = &cobra.Command{
	Use:   "bars [ticker]",
	Short: "Get OHLC aggregate bars for a crypto ticker",
	Long:  "Retrieve custom OHLC (Open, High, Low, Close) aggregate bar data for a crypto ticker over a specified time range.",
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

		params := api.BarsParams{
			Multiplier: multiplier,
			Timespan:   timespan,
			From:       from,
			To:         to,
			Adjusted:   adjusted,
			Sort:       sort,
			Limit:      limit,
		}

		result, err := client.GetCryptoBars(ticker, params)
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
			fmt.Fprintf(w, "%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\t%.4f\t%d\n",
				t.Format("2006-01-02"),
				bar.Open, bar.High, bar.Low, bar.Close,
				bar.Volume, bar.VWAP, bar.NumTrades)
		}
		w.Flush()

		return nil
	},
}

// cryptoDailyMarketSummaryCmd retrieves the grouped daily OHLC summary
// for all crypto tickers on a specified date. Useful for broad crypto
// market analysis and screening.
// Usage: massive crypto daily-market-summary 2024-01-09
var cryptoDailyMarketSummaryCmd = &cobra.Command{
	Use:   "daily-market-summary [date]",
	Short: "Get daily market summary for all crypto tickers",
	Long:  "Retrieve the daily OHLC, volume, and VWAP data for all crypto tickers on a specified trading date.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		date := args[0]
		adjusted, _ := cmd.Flags().GetString("adjusted")

		result, err := client.GetCryptoDailyMarketSummary(date, adjusted)
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
			fmt.Fprintf(w, "%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\t%.4f\t%d\n",
				s.Ticker, s.Open, s.High, s.Low, s.Close,
				s.Volume, s.VWAP, s.NumTrades)
		}
		w.Flush()

		return nil
	},
}

// cryptoDailyTickerSummaryCmd retrieves the daily open/close data for
// a specific crypto pair (from/to) on a given date. The response includes
// opening and closing trades along with aggregate prices.
// Usage: massive crypto daily-ticker-summary BTC USD 2024-01-09
var cryptoDailyTickerSummaryCmd = &cobra.Command{
	Use:   "daily-ticker-summary [from] [to] [date]",
	Short: "Get daily open/close for a crypto pair",
	Long:  "Retrieve the daily opening and closing prices and trades for a crypto currency pair on a given date.",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		from := strings.ToUpper(args[0])
		to := strings.ToUpper(args[1])
		date := args[2]
		adjusted, _ := cmd.Flags().GetString("adjusted")

		result, err := client.GetCryptoDailyTickerSummary(from, to, date, adjusted)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Symbol: %s | Date: %s | UTC: %v\n", result.Symbol, result.Day, result.IsUTC)
		fmt.Printf("Open:   %.4f\n", result.Open)
		fmt.Printf("Close:  %.4f\n\n", result.Close)

		if len(result.OpenTrades) > 0 {
			fmt.Printf("Open Trades: %d\n", len(result.OpenTrades))
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tPRICE\tSIZE\tEXCHANGE\tTIMESTAMP")
			fmt.Fprintln(w, "--\t-----\t----\t--------\t---------")
			for _, trade := range result.OpenTrades {
				t := time.UnixMilli(trade.Timestamp)
				fmt.Fprintf(w, "%s\t%.4f\t%.4f\t%d\t%s\n",
					trade.ID, trade.Price, trade.Size, trade.Exchange,
					t.Format("2006-01-02 15:04:05"))
			}
			w.Flush()
			fmt.Println()
		}

		if len(result.ClosingTrades) > 0 {
			fmt.Printf("Closing Trades: %d\n", len(result.ClosingTrades))
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tPRICE\tSIZE\tEXCHANGE\tTIMESTAMP")
			fmt.Fprintln(w, "--\t-----\t----\t--------\t---------")
			for _, trade := range result.ClosingTrades {
				t := time.UnixMilli(trade.Timestamp)
				fmt.Fprintf(w, "%s\t%.4f\t%.4f\t%d\t%s\n",
					trade.ID, trade.Price, trade.Size, trade.Exchange,
					t.Format("2006-01-02 15:04:05"))
			}
			w.Flush()
		}

		return nil
	},
}

// cryptoPreviousDayBarCmd retrieves the previous day's OHLC bar data
// for a specific crypto ticker. Useful for quick comparisons with
// current trading activity.
// Usage: massive crypto previous-day-bar X:BTCUSD
var cryptoPreviousDayBarCmd = &cobra.Command{
	Use:   "previous-day-bar [ticker]",
	Short: "Get previous day's bar for a crypto ticker",
	Long:  "Retrieve the previous trading day's OHLC bar data for a specific crypto ticker.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		adjusted, _ := cmd.Flags().GetString("adjusted")

		result, err := client.GetCryptoPreviousDayBar(ticker, adjusted)
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
			fmt.Fprintf(w, "%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\t%.4f\t%d\n",
				t.Format("2006-01-02"),
				bar.Open, bar.High, bar.Low, bar.Close,
				bar.Volume, bar.VWAP, bar.NumTrades)
		}
		w.Flush()

		return nil
	},
}

// -------------------------------------------------------------------
// Market Operations Commands
// -------------------------------------------------------------------

// cryptoConditionsCmd retrieves the list of condition codes that apply
// to crypto trade data. These codes describe the nature of trades
// (e.g., regular sale, block trade).
// Usage: massive crypto conditions
var cryptoConditionsCmd = &cobra.Command{
	Use:   "conditions",
	Short: "List crypto trade condition codes",
	Long:  "Retrieve the list of condition codes used for crypto trade data classification.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		result, err := client.GetCryptoConditions()
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Conditions: %d\n\n", result.Count)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTYPE\tASSET CLASS\tDATA TYPES")
		fmt.Fprintln(w, "--\t----\t----\t-----------\t----------")

		for _, c := range result.Results {
			dataTypes := strings.Join(c.DataTypes, ", ")
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
				c.ID, c.Name, c.Type, c.AssetClass, dataTypes)
		}
		w.Flush()

		return nil
	},
}

// cryptoExchangesCmd retrieves the list of known crypto exchanges
// with their identifiers and metadata.
// Usage: massive crypto exchanges
var cryptoExchangesCmd = &cobra.Command{
	Use:   "exchanges",
	Short: "List known crypto exchanges",
	Long:  "Retrieve a list of known cryptocurrency exchanges including their identifiers, names, and metadata.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		result, err := client.GetCryptoExchanges()
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Exchanges: %d\n\n", result.Count)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tACRONYM\tTYPE\tLOCALE")
		fmt.Fprintln(w, "--\t----\t-------\t----\t------")

		for _, e := range result.Results {
			acronym := e.Acronym
			if acronym == "" {
				acronym = "-"
			}
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
				e.ID, e.Name, acronym, e.Type, e.Locale)
		}
		w.Flush()

		return nil
	},
}

// cryptoMarketHolidaysCmd retrieves the list of upcoming market holidays
// and early-close days. Reuses the same endpoint as stocks since market
// holidays apply globally.
// Usage: massive crypto market-holidays
var cryptoMarketHolidaysCmd = &cobra.Command{
	Use:   "market-holidays",
	Short: "Get upcoming market holidays",
	Long:  "Retrieve upcoming market holidays and early-close days for all exchanges.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		result, err := client.GetMarketHolidays()
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		if len(result) == 0 {
			fmt.Println("No upcoming market holidays found.")
			return nil
		}

		fmt.Printf("Upcoming Market Holidays: %d\n\n", len(result))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "DATE\tEXCHANGE\tNAME\tSTATUS\tOPEN\tCLOSE")
		fmt.Fprintln(w, "----\t--------\t----\t------\t----\t-----")

		for _, h := range result {
			openTime := "-"
			closeTime := "-"
			if h.Open != "" {
				openTime = h.Open
			}
			if h.Close != "" {
				closeTime = h.Close
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				h.Date, h.Exchange, h.Name, h.Status, openTime, closeTime)
		}
		w.Flush()

		return nil
	},
}

// cryptoMarketStatusCmd retrieves the current real-time status of all
// markets including crypto exchanges and currency markets.
// Usage: massive crypto market-status
var cryptoMarketStatusCmd = &cobra.Command{
	Use:   "market-status",
	Short: "Get current market status",
	Long:  "Retrieve the real-time open/closed status of all markets including crypto exchanges, currency markets, and stock exchanges.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		result, err := client.GetMarketStatus()
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Market: %s | Server Time: %s\n", result.Market, result.ServerTime)
		fmt.Printf("After Hours: %v | Early Hours: %v\n\n", result.AfterHours, result.EarlyHours)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		fmt.Fprintln(w, "CURRENCIES")
		fmt.Fprintln(w, "----------")
		fmt.Fprintf(w, "Crypto\t%s\n", result.Currencies.Crypto)
		fmt.Fprintf(w, "Forex\t%s\n", result.Currencies.FX)
		fmt.Fprintln(w)

		fmt.Fprintln(w, "EXCHANGES")
		fmt.Fprintln(w, "--------")
		fmt.Fprintf(w, "NYSE\t%s\n", result.Exchanges.NYSE)
		fmt.Fprintf(w, "NASDAQ\t%s\n", result.Exchanges.Nasdaq)
		fmt.Fprintf(w, "OTC\t%s\n", result.Exchanges.OTC)

		w.Flush()

		return nil
	},
}

// -------------------------------------------------------------------
// Snapshot Commands
// -------------------------------------------------------------------

// cryptoSnapshotCmd retrieves the most recent snapshot for a single
// crypto ticker including the current day's bar, previous day's bar,
// latest minute bar, last trade, and fair market value.
// Usage: massive crypto snapshot X:BTCUSD
var cryptoSnapshotCmd = &cobra.Command{
	Use:   "snapshot [ticker]",
	Short: "Get snapshot for a single crypto ticker",
	Long:  "Retrieve the most recent snapshot for a single crypto ticker including current day, previous day, minute bar, last trade, and fair market value.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])

		result, err := client.GetCryptoSnapshotSingleTicker(ticker)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		t := result.Ticker
		fmt.Printf("Ticker: %s | Change: %.4f (%.2f%%) | FMV: %.4f\n\n",
			t.Ticker, t.TodaysChange, t.TodaysChangePct, t.FMV)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "PERIOD\tOPEN\tHIGH\tLOW\tCLOSE\tVOLUME\tVWAP")
		fmt.Fprintln(w, "------\t----\t----\t---\t-----\t------\t----")

		fmt.Fprintf(w, "Day\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\t%.4f\n",
			t.Day.Open, t.Day.High, t.Day.Low, t.Day.Close,
			t.Day.Volume, t.Day.VWAP)

		fmt.Fprintf(w, "Prev Day\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\t%.4f\n",
			t.PrevDay.Open, t.PrevDay.High, t.PrevDay.Low, t.PrevDay.Close,
			t.PrevDay.Volume, t.PrevDay.VWAP)

		fmt.Fprintf(w, "Minute\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\t%.4f\n",
			t.Min.Open, t.Min.High, t.Min.Low, t.Min.Close,
			t.Min.Volume, t.Min.VWAP)

		w.Flush()

		fmt.Printf("\nLast Trade: Price=%.4f Size=%.4f Exchange=%d\n",
			t.LastTrade.Price, t.LastTrade.Size, t.LastTrade.Exchange)

		return nil
	},
}

// cryptoSnapshotMarketCmd retrieves snapshot data for all crypto tickers
// or a filtered subset. Supports filtering by a comma-separated list of
// ticker symbols.
// Usage: massive crypto snapshot-market --tickers X:BTCUSD,X:ETHUSD
var cryptoSnapshotMarketCmd = &cobra.Command{
	Use:   "snapshot-market",
	Short: "Get snapshots for all or selected crypto tickers",
	Long:  "Retrieve snapshot data for all crypto tickers or a filtered subset specified by a comma-separated list of symbols.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		tickers, _ := cmd.Flags().GetString("tickers")

		params := api.CryptoSnapshotParams{
			Tickers: tickers,
		}

		result, err := client.GetCryptoSnapshotFullMarket(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Tickers: %d\n\n", len(result.Tickers))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tDAY OPEN\tDAY HIGH\tDAY LOW\tDAY CLOSE\tVOLUME\tCHANGE\tCHANGE %\tFMV")
		fmt.Fprintln(w, "------\t--------\t--------\t-------\t---------\t------\t------\t--------\t---")

		for _, t := range result.Tickers {
			fmt.Fprintf(w, "%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\t%.4f\t%.2f%%\t%.4f\n",
				t.Ticker, t.Day.Open, t.Day.High, t.Day.Low, t.Day.Close,
				t.Day.Volume, t.TodaysChange, t.TodaysChangePct, t.FMV)
		}
		w.Flush()

		return nil
	},
}

// cryptoGainersCmd retrieves the current top crypto gainers with snapshot
// data including day bar, previous day bar, and percentage change values.
// Usage: massive crypto gainers
var cryptoGainersCmd = &cobra.Command{
	Use:   "gainers",
	Short: "Get top gaining crypto tickers",
	Long:  "Retrieve the current top gainers in the crypto market with snapshot data including day bar, previous day bar, and change percentages.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		result, err := client.GetCryptoSnapshotTopMovers("gainers")
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		return printCryptoMoversTable("Gainers", result)
	},
}

// cryptoLosersCmd retrieves the current top crypto losers with snapshot
// data including day bar, previous day bar, and percentage change values.
// Usage: massive crypto losers
var cryptoLosersCmd = &cobra.Command{
	Use:   "losers",
	Short: "Get top losing crypto tickers",
	Long:  "Retrieve the current top losers in the crypto market with snapshot data including day bar, previous day bar, and change percentages.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		result, err := client.GetCryptoSnapshotTopMovers("losers")
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		return printCryptoMoversTable("Losers", result)
	},
}

// printCryptoMoversTable formats and prints a table of crypto gainers or
// losers snapshot data to stdout. The title parameter labels the output
// as either "Gainers" or "Losers" for display clarity.
func printCryptoMoversTable(title string, result *api.CryptoSnapshotResponse) error {
	fmt.Printf("Top %s: %d tickers\n\n", title, len(result.Tickers))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TICKER\tDAY OPEN\tDAY HIGH\tDAY LOW\tDAY CLOSE\tVOLUME\tCHANGE\tCHANGE %\tFMV")
	fmt.Fprintln(w, "------\t--------\t--------\t-------\t---------\t------\t------\t--------\t---")

	for _, t := range result.Tickers {
		fmt.Fprintf(w, "%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\t%.4f\t%.2f%%\t%.4f\n",
			t.Ticker, t.Day.Open, t.Day.High, t.Day.Low, t.Day.Close,
			t.Day.Volume, t.TodaysChange, t.TodaysChangePct, t.FMV)
	}
	w.Flush()

	return nil
}

// -------------------------------------------------------------------
// Technical Indicator Commands
// -------------------------------------------------------------------

// cryptoSMACmd retrieves Simple Moving Average (SMA) data for a crypto
// ticker over a specified time range. SMA smooths price data by calculating
// the arithmetic mean over a rolling window period.
// Usage: massive crypto sma X:BTCUSD --from 2025-01-06 --to 2025-01-10
var cryptoSMACmd = &cobra.Command{
	Use:   "sma [ticker]",
	Short: "Get Simple Moving Average (SMA) for a crypto ticker",
	Long:  "Retrieve Simple Moving Average (SMA) indicator data for a crypto ticker. SMA calculates the arithmetic mean of closing prices over a given window period.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		params := buildCryptoIndicatorParams(cmd)

		result, err := client.GetCryptoSMA(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		printIndicatorTable(ticker, "SMA", result)
		return nil
	},
}

// cryptoEMACmd retrieves Exponential Moving Average (EMA) data for a
// crypto ticker over a specified time range. EMA places greater weight
// on recent prices compared to SMA for quicker trend detection.
// Usage: massive crypto ema X:BTCUSD --from 2025-01-06 --to 2025-01-10
var cryptoEMACmd = &cobra.Command{
	Use:   "ema [ticker]",
	Short: "Get Exponential Moving Average (EMA) for a crypto ticker",
	Long:  "Retrieve Exponential Moving Average (EMA) indicator data for a crypto ticker. EMA places greater weight on recent prices for more responsive trend signals.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		params := buildCryptoIndicatorParams(cmd)

		result, err := client.GetCryptoEMA(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		printIndicatorTable(ticker, "EMA", result)
		return nil
	},
}

// cryptoRSICmd retrieves Relative Strength Index (RSI) data for a crypto
// ticker over a specified time range. RSI oscillates between 0 and 100
// to identify overbought or oversold conditions.
// Usage: massive crypto rsi X:BTCUSD --from 2025-01-06 --to 2025-01-10
var cryptoRSICmd = &cobra.Command{
	Use:   "rsi [ticker]",
	Short: "Get Relative Strength Index (RSI) for a crypto ticker",
	Long:  "Retrieve Relative Strength Index (RSI) indicator data for a crypto ticker. RSI measures the speed and magnitude of price changes, oscillating between 0 and 100.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		params := buildCryptoIndicatorParams(cmd)

		result, err := client.GetCryptoRSI(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		printIndicatorTable(ticker, "RSI", result)
		return nil
	},
}

// cryptoMACDCmd retrieves Moving Average Convergence/Divergence (MACD) data
// for a crypto ticker over a specified time range. MACD is a momentum
// indicator with three components: the MACD line, signal line, and histogram.
// Usage: massive crypto macd X:BTCUSD --from 2025-01-06 --to 2025-01-10
var cryptoMACDCmd = &cobra.Command{
	Use:   "macd [ticker]",
	Short: "Get Moving Average Convergence/Divergence (MACD) for a crypto ticker",
	Long:  "Retrieve MACD indicator data for a crypto ticker. MACD is a momentum indicator showing the relationship between two EMAs, with signal line and histogram.",
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

		result, err := client.GetCryptoMACD(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		printMACDTable(ticker, result)
		return nil
	},
}

// buildCryptoIndicatorParams extracts the common indicator flags from the
// cobra command and returns a populated IndicatorParams struct. This is
// shared by the crypto SMA, EMA, and RSI commands.
func buildCryptoIndicatorParams(cmd *cobra.Command) api.IndicatorParams {
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

// addCryptoIndicatorFlags registers the common flags shared by the crypto
// SMA, EMA, and RSI indicator subcommands. These include date range,
// window, timespan, series type, and pagination controls.
func addCryptoIndicatorFlags(cmd *cobra.Command, defaultWindow string) {
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

// -------------------------------------------------------------------
// Tickers Commands
// -------------------------------------------------------------------

// cryptoTickersCmd lists and searches crypto tickers from the Massive
// reference data. Supports filtering by name search, active status,
// sort field, sort order, and result limit.
// Usage: massive crypto tickers --search "Bitcoin"
var cryptoTickersCmd = &cobra.Command{
	Use:   "tickers",
	Short: "List and search crypto tickers",
	Long:  "Retrieve a list of crypto tickers with optional filtering by name, active status, and pagination controls.",
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

		params := api.CryptoTickersParams{
			Search: search,
			Active: active,
			Sort:   sort,
			Order:  order,
			Limit:  limit,
		}

		result, err := client.GetCryptoTickers(params)
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

// cryptoTickerOverviewCmd retrieves detailed reference information for
// a specific crypto ticker including currency details and active status.
// Usage: massive crypto ticker-overview X:BTCUSD
var cryptoTickerOverviewCmd = &cobra.Command{
	Use:   "ticker-overview [ticker]",
	Short: "Get detailed overview for a crypto ticker",
	Long:  "Retrieve detailed reference information for a specific crypto ticker including currency details, base currency, and active status.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])

		result, err := client.GetCryptoTickerOverview(ticker)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		r := result.Results
		fmt.Printf("Ticker:              %s\n", r.Ticker)
		fmt.Printf("Name:                %s\n", r.Name)
		fmt.Printf("Market:              %s\n", r.Market)
		fmt.Printf("Locale:              %s\n", r.Locale)
		fmt.Printf("Active:              %v\n", r.Active)
		fmt.Printf("Currency Symbol:     %s\n", r.CurrencySymbol)
		fmt.Printf("Currency Name:       %s\n", r.CurrencyName)
		fmt.Printf("Base Currency Symbol: %s\n", r.BaseCurrencySymbol)
		fmt.Printf("Base Currency Name:  %s\n", r.BaseCurrencyName)
		fmt.Printf("Last Updated:        %s\n", r.LastUpdatedUTC)

		return nil
	},
}

// -------------------------------------------------------------------
// Trades Commands
// -------------------------------------------------------------------

// cryptoTradesCmd retrieves tick-level trade data for a specific crypto
// ticker with optional timestamp filtering, sorting, and pagination.
// Usage: massive crypto trades X:BTCUSD --limit 10
var cryptoTradesCmd = &cobra.Command{
	Use:   "trades [ticker]",
	Short: "Get tick-level trade data for a crypto ticker",
	Long:  "Retrieve tick-level trade data for a crypto ticker including price, size, exchange, conditions, and timestamps.",
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

		params := api.CryptoTradesParams{
			Timestamp:    timestamp,
			TimestampGte: timestampGte,
			TimestampGt:  timestampGt,
			TimestampLte: timestampLte,
			TimestampLt:  timestampLt,
			Order:        order,
			Limit:        limit,
			Sort:         sort,
		}

		result, err := client.GetCryptoTrades(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Ticker: %s | Trades: %d\n\n", ticker, len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TIMESTAMP\tPRICE\tSIZE\tEXCHANGE\tID")
		fmt.Fprintln(w, "---------\t-----\t----\t--------\t--")

		for _, trade := range result.Results {
			t := time.Unix(0, trade.ParticipantTimestamp)
			fmt.Fprintf(w, "%s\t%.4f\t%.4f\t%d\t%s\n",
				t.Format("2006-01-02 15:04:05.000"),
				trade.Price, trade.Size, trade.Exchange, trade.ID)
		}
		w.Flush()

		return nil
	},
}

// cryptoLastTradeCmd retrieves the most recent trade for a specific crypto
// pair. Returns price, size, exchange, and timestamp information useful
// for monitoring current market activity.
// Usage: massive crypto last-trade BTC USD
var cryptoLastTradeCmd = &cobra.Command{
	Use:   "last-trade [from] [to]",
	Short: "Get the most recent trade for a crypto pair",
	Long:  "Retrieve the last available trade for a crypto currency pair including price, size, exchange, conditions, and timestamp.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		from := strings.ToUpper(args[0])
		to := strings.ToUpper(args[1])

		result, err := client.GetCryptoLastTrade(from, to)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		last := result.Last
		t := time.UnixMilli(last.Timestamp)

		fmt.Printf("Symbol:    %s\n", result.Symbol)
		fmt.Printf("Price:     %.4f\n", last.Price)
		fmt.Printf("Size:      %.4f\n", last.Size)
		fmt.Printf("Exchange:  %d\n", last.Exchange)
		fmt.Printf("Timestamp: %s\n", t.Format("2006-01-02 15:04:05.000"))

		if len(last.Conditions) > 0 {
			condStrs := make([]string, len(last.Conditions))
			for i, c := range last.Conditions {
				condStrs[i] = fmt.Sprintf("%d", c)
			}
			fmt.Printf("Conditions: %s\n", strings.Join(condStrs, ", "))
		}

		return nil
	},
}

// init registers the crypto parent command and all subcommands with
// their respective flags under the root command.
func init() {
	// Register crypto parent command under root
	rootCmd.AddCommand(cryptoCmd)

	// Bars command flags
	cryptoBarsCmd.Flags().String("multiplier", "1", "Size of the timespan multiplier")
	cryptoBarsCmd.Flags().String("timespan", "day", "Timespan (minute, hour, day, week, month, quarter, year)")
	cryptoBarsCmd.Flags().String("from", "", "Start date (YYYY-MM-DD) [required]")
	cryptoBarsCmd.Flags().String("to", "", "End date (YYYY-MM-DD) [required]")
	cryptoBarsCmd.Flags().String("adjusted", "true", "Adjust for splits (true/false)")
	cryptoBarsCmd.Flags().String("sort", "asc", "Sort order (asc/desc)")
	cryptoBarsCmd.Flags().String("limit", "5000", "Max number of results (max 50000)")
	cryptoBarsCmd.MarkFlagRequired("from")
	cryptoBarsCmd.MarkFlagRequired("to")
	cryptoCmd.AddCommand(cryptoBarsCmd)

	// Daily market summary command flags
	cryptoDailyMarketSummaryCmd.Flags().String("adjusted", "true", "Adjust for splits (true/false)")
	cryptoCmd.AddCommand(cryptoDailyMarketSummaryCmd)

	// Daily ticker summary command flags
	cryptoDailyTickerSummaryCmd.Flags().String("adjusted", "true", "Adjust for splits (true/false)")
	cryptoCmd.AddCommand(cryptoDailyTickerSummaryCmd)

	// Previous day bar command flags
	cryptoPreviousDayBarCmd.Flags().String("adjusted", "true", "Adjust for splits (true/false)")
	cryptoCmd.AddCommand(cryptoPreviousDayBarCmd)

	// Market operations commands
	cryptoCmd.AddCommand(cryptoConditionsCmd)
	cryptoCmd.AddCommand(cryptoExchangesCmd)
	cryptoCmd.AddCommand(cryptoMarketHolidaysCmd)
	cryptoCmd.AddCommand(cryptoMarketStatusCmd)

	// Snapshot commands
	cryptoCmd.AddCommand(cryptoSnapshotCmd)

	cryptoSnapshotMarketCmd.Flags().String("tickers", "", "Comma-separated list of ticker symbols (default: all)")
	cryptoCmd.AddCommand(cryptoSnapshotMarketCmd)

	cryptoCmd.AddCommand(cryptoGainersCmd)
	cryptoCmd.AddCommand(cryptoLosersCmd)

	// Technical indicator commands
	addCryptoIndicatorFlags(cryptoSMACmd, "10")
	cryptoCmd.AddCommand(cryptoSMACmd)

	addCryptoIndicatorFlags(cryptoEMACmd, "10")
	cryptoCmd.AddCommand(cryptoEMACmd)

	addCryptoIndicatorFlags(cryptoRSICmd, "14")
	cryptoCmd.AddCommand(cryptoRSICmd)

	// MACD flags
	cryptoMACDCmd.Flags().String("from", "", "Start date (YYYY-MM-DD) [required]")
	cryptoMACDCmd.Flags().String("to", "", "End date (YYYY-MM-DD) [required]")
	cryptoMACDCmd.Flags().String("timespan", "day", "Aggregate time window (minute, hour, day, week, month, quarter, year)")
	cryptoMACDCmd.Flags().String("adjusted", "true", "Adjust for splits (true/false)")
	cryptoMACDCmd.Flags().String("short-window", "12", "Short EMA period for MACD line")
	cryptoMACDCmd.Flags().String("long-window", "26", "Long EMA period for MACD line")
	cryptoMACDCmd.Flags().String("signal-window", "9", "Signal line EMA period")
	cryptoMACDCmd.Flags().String("series-type", "close", "Price type for calculation (open, high, low, close)")
	cryptoMACDCmd.Flags().String("order", "desc", "Sort order by timestamp (asc/desc)")
	cryptoMACDCmd.Flags().String("limit", "10", "Max number of results (max 5000)")
	cryptoMACDCmd.MarkFlagRequired("from")
	cryptoMACDCmd.MarkFlagRequired("to")
	cryptoCmd.AddCommand(cryptoMACDCmd)

	// Tickers command flags
	cryptoTickersCmd.Flags().String("search", "", "Search by name or symbol")
	cryptoTickersCmd.Flags().String("active", "", "Filter by active status (true/false)")
	cryptoTickersCmd.Flags().String("sort", "ticker", "Sort field (ticker, name)")
	cryptoTickersCmd.Flags().String("order", "asc", "Sort order (asc/desc)")
	cryptoTickersCmd.Flags().String("limit", "20", "Number of results to return (max 1000)")
	cryptoCmd.AddCommand(cryptoTickersCmd)

	// Ticker overview command
	cryptoCmd.AddCommand(cryptoTickerOverviewCmd)

	// Trades command flags
	cryptoTradesCmd.Flags().String("timestamp", "", "Filter by date (YYYY-MM-DD) or nanosecond timestamp")
	cryptoTradesCmd.Flags().String("timestamp-gte", "", "Timestamp greater than or equal to")
	cryptoTradesCmd.Flags().String("timestamp-gt", "", "Timestamp greater than")
	cryptoTradesCmd.Flags().String("timestamp-lte", "", "Timestamp less than or equal to")
	cryptoTradesCmd.Flags().String("timestamp-lt", "", "Timestamp less than")
	cryptoTradesCmd.Flags().String("order", "", "Sort order (asc/desc)")
	cryptoTradesCmd.Flags().String("limit", "1000", "Max number of results (max 50000)")
	cryptoTradesCmd.Flags().String("sort", "", "Sort field (e.g., timestamp)")
	cryptoCmd.AddCommand(cryptoTradesCmd)

	// Last trade command
	cryptoCmd.AddCommand(cryptoLastTradeCmd)
}
