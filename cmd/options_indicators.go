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

// optionsSMACmd retrieves Simple Moving Average (SMA) data for an options
// contract ticker over a specified time range. SMA smooths price data by
// calculating the arithmetic mean over a rolling window period.
// Usage: massive options sma O:AAPL250117C00150000 --from 2025-01-06 --to 2025-01-10
var optionsSMACmd = &cobra.Command{
	Use:   "sma [ticker]",
	Short: "Get Simple Moving Average (SMA) for an options contract",
	Long:  "Retrieve Simple Moving Average (SMA) indicator data for an options contract ticker (e.g., O:AAPL250117C00150000). SMA calculates the arithmetic mean of values over a given window period.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		params := buildOptionsIndicatorParams(cmd)

		result, err := client.GetOptionsSMA(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		printOptionsIndicatorTable(ticker, "SMA", result)
		return nil
	},
}

// optionsEMACmd retrieves Exponential Moving Average (EMA) data for an options
// contract ticker over a specified time range. EMA places greater weight on
// recent values compared to SMA for quicker trend detection.
// Usage: massive options ema O:AAPL250117C00150000 --from 2025-01-06 --to 2025-01-10
var optionsEMACmd = &cobra.Command{
	Use:   "ema [ticker]",
	Short: "Get Exponential Moving Average (EMA) for an options contract",
	Long:  "Retrieve Exponential Moving Average (EMA) indicator data for an options contract ticker (e.g., O:AAPL250117C00150000). EMA places greater weight on recent values for more responsive trend signals.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		params := buildOptionsIndicatorParams(cmd)

		result, err := client.GetOptionsEMA(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		printOptionsIndicatorTable(ticker, "EMA", result)
		return nil
	},
}

// optionsRSICmd retrieves Relative Strength Index (RSI) data for an options
// contract ticker over a specified time range. RSI oscillates between 0 and
// 100 to identify overbought or oversold conditions.
// Usage: massive options rsi O:AAPL250117C00150000 --from 2025-01-06 --to 2025-01-10
var optionsRSICmd = &cobra.Command{
	Use:   "rsi [ticker]",
	Short: "Get Relative Strength Index (RSI) for an options contract",
	Long:  "Retrieve Relative Strength Index (RSI) indicator data for an options contract ticker (e.g., O:AAPL250117C00150000). RSI measures the speed and magnitude of price changes, oscillating between 0 and 100.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker := strings.ToUpper(args[0])
		params := buildOptionsIndicatorParams(cmd)

		result, err := client.GetOptionsRSI(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		printOptionsIndicatorTable(ticker, "RSI", result)
		return nil
	},
}

// optionsMACDCmd retrieves Moving Average Convergence/Divergence (MACD) data
// for an options contract ticker over a specified time range. MACD is a
// momentum indicator with three components: the MACD line, signal line, and
// histogram.
// Usage: massive options macd O:AAPL250117C00150000 --from 2025-01-06 --to 2025-01-10
var optionsMACDCmd = &cobra.Command{
	Use:   "macd [ticker]",
	Short: "Get Moving Average Convergence/Divergence (MACD) for an options contract",
	Long:  "Retrieve MACD indicator data for an options contract ticker (e.g., O:AAPL250117C00150000). MACD is a momentum indicator showing the relationship between two EMAs, with signal line and histogram.",
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

		result, err := client.GetOptionsMACD(ticker, params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		printOptionsMACDTable(ticker, result)
		return nil
	},
}

// buildOptionsIndicatorParams extracts the common indicator flags from the
// cobra command and returns a populated IndicatorParams struct. This is shared
// by the options SMA, EMA, and RSI commands which all use the same parameters.
func buildOptionsIndicatorParams(cmd *cobra.Command) api.IndicatorParams {
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

// printOptionsIndicatorTable renders a formatted table of indicator values for
// the options SMA, EMA, or RSI commands. Each row displays the date and
// computed value.
func printOptionsIndicatorTable(ticker, indicator string, result *api.IndicatorResponse) {
	fmt.Printf("Ticker: %s | Indicator: %s | Values: %d\n\n", ticker, indicator, len(result.Results.Values))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "DATE\tVALUE")
	fmt.Fprintln(w, "----\t-----")

	for _, v := range result.Results.Values {
		t := time.UnixMilli(v.Timestamp)
		fmt.Fprintf(w, "%s\t%.4f\n", t.Format("2006-01-02"), v.Value)
	}
	w.Flush()
}

// printOptionsMACDTable renders a formatted table of MACD indicator values
// including the MACD line, signal line, and histogram for each data point
// of an options contract ticker.
func printOptionsMACDTable(ticker string, result *api.MACDResponse) {
	fmt.Printf("Ticker: %s | Indicator: MACD | Values: %d\n\n", ticker, len(result.Results.Values))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "DATE\tMACD\tSIGNAL\tHISTOGRAM")
	fmt.Fprintln(w, "----\t----\t------\t---------")

	for _, v := range result.Results.Values {
		t := time.UnixMilli(v.Timestamp)
		fmt.Fprintf(w, "%s\t%.4f\t%.4f\t%.4f\n",
			t.Format("2006-01-02"), v.Value, v.Signal, v.Histogram)
	}
	w.Flush()
}

// addOptionsIndicatorFlags registers the common flags shared by the options
// SMA, EMA, and RSI indicator subcommands. These include date range, window,
// timespan, series type, and pagination controls.
func addOptionsIndicatorFlags(cmd *cobra.Command, defaultWindow string) {
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

// init registers the SMA, EMA, RSI, and MACD indicator subcommands and their
// flags under the options parent command.
func init() {
	// SMA flags
	addOptionsIndicatorFlags(optionsSMACmd, "10")
	optionsCmd.AddCommand(optionsSMACmd)

	// EMA flags
	addOptionsIndicatorFlags(optionsEMACmd, "10")
	optionsCmd.AddCommand(optionsEMACmd)

	// RSI flags
	addOptionsIndicatorFlags(optionsRSICmd, "14")
	optionsCmd.AddCommand(optionsRSICmd)

	// MACD flags
	optionsMACDCmd.Flags().String("from", "", "Start date (YYYY-MM-DD) [required]")
	optionsMACDCmd.Flags().String("to", "", "End date (YYYY-MM-DD) [required]")
	optionsMACDCmd.Flags().String("timespan", "day", "Aggregate time window (minute, hour, day, week, month, quarter, year)")
	optionsMACDCmd.Flags().String("adjusted", "true", "Adjust for splits (true/false)")
	optionsMACDCmd.Flags().String("short-window", "12", "Short EMA period for MACD line")
	optionsMACDCmd.Flags().String("long-window", "26", "Long EMA period for MACD line")
	optionsMACDCmd.Flags().String("signal-window", "9", "Signal line EMA period")
	optionsMACDCmd.Flags().String("series-type", "close", "Price type for calculation (open, high, low, close)")
	optionsMACDCmd.Flags().String("order", "desc", "Sort order by timestamp (asc/desc)")
	optionsMACDCmd.Flags().String("limit", "10", "Max number of results (max 5000)")

	optionsMACDCmd.MarkFlagRequired("from")
	optionsMACDCmd.MarkFlagRequired("to")

	optionsCmd.AddCommand(optionsMACDCmd)
}
