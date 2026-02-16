//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/cloudmanic/massive-cli/internal/config"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

const (
	// wsDelayedURL is the base WebSocket endpoint for delayed (15-min) streaming data.
	wsDelayedURL = "wss://delayed.massive.com"

	// wsRealtimeURL is the base WebSocket endpoint for real-time streaming data.
	wsRealtimeURL = "wss://socket.massive.com"
)

// wsRealtime controls whether to connect to the real-time or delayed WebSocket endpoint.
var wsRealtime bool

// wsCmd is the parent command for all WebSocket streaming subcommands.
// It groups real-time streaming commands for stocks, crypto, forex, etc.
var wsCmd = &cobra.Command{
	Use:   "ws",
	Short: "WebSocket streaming commands",
}

// wsStocksCmd is the parent command for all stock-related WebSocket streaming
// subcommands. It groups real-time stock data streams like trades, quotes,
// aggregates, LULD bands, and fair market value.
var wsStocksCmd = &cobra.Command{
	Use:   "stocks",
	Short: "Stream real-time stocks data",
}

// wsStocksTradesCmd streams real-time stock trades from the WebSocket API.
// Accepts one or more ticker symbols as arguments, or use --all to subscribe
// to all tickers. Outputs trade events with symbol, price, size, and exchange.
// Usage: massive ws stocks trades AAPL MSFT --output table
var wsStocksTradesCmd = &cobra.Command{
	Use:   "trades [tickers...]",
	Short: "Stream real-time stock trades",
	Long:  "Stream real-time stock trade data via WebSocket. Provide ticker symbols as arguments or use --all for all tickers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		tickers, err := buildTickerParams("T", args, all)
		if err != nil {
			return err
		}
		return connectAndStream(cmd.Context(), "T", tickers, formatTrade)
	},
}

// wsStocksQuotesCmd streams real-time NBBO quotes from the WebSocket API.
// Accepts one or more ticker symbols as arguments, or use --all to subscribe
// to all tickers. Outputs quote events with symbol, bid/ask prices, and sizes.
// Usage: massive ws stocks quotes AAPL MSFT --output table
var wsStocksQuotesCmd = &cobra.Command{
	Use:   "quotes [tickers...]",
	Short: "Stream real-time stock quotes (NBBO)",
	Long:  "Stream real-time NBBO (National Best Bid and Offer) quote data via WebSocket. Provide ticker symbols as arguments or use --all for all tickers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		tickers, err := buildTickerParams("Q", args, all)
		if err != nil {
			return err
		}
		return connectAndStream(cmd.Context(), "Q", tickers, formatQuote)
	},
}

// wsStocksAggMinuteCmd streams per-minute aggregate bars from the WebSocket API.
// Accepts one or more ticker symbols as arguments, or use --all to subscribe
// to all tickers. Outputs OHLCV data aggregated at one-minute intervals.
// Usage: massive ws stocks agg-minute AAPL --output table
var wsStocksAggMinuteCmd = &cobra.Command{
	Use:   "agg-minute [tickers...]",
	Short: "Stream per-minute aggregate bars",
	Long:  "Stream per-minute OHLCV aggregate bar data via WebSocket. Provide ticker symbols as arguments or use --all for all tickers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		tickers, err := buildTickerParams("AM", args, all)
		if err != nil {
			return err
		}
		return connectAndStream(cmd.Context(), "AM", tickers, formatAggregate)
	},
}

// wsStocksAggSecondCmd streams per-second aggregate bars from the WebSocket API.
// Accepts one or more ticker symbols as arguments, or use --all to subscribe
// to all tickers. Outputs OHLCV data aggregated at one-second intervals.
// Usage: massive ws stocks agg-second AAPL --output table
var wsStocksAggSecondCmd = &cobra.Command{
	Use:   "agg-second [tickers...]",
	Short: "Stream per-second aggregate bars",
	Long:  "Stream per-second OHLCV aggregate bar data via WebSocket. Provide ticker symbols as arguments or use --all for all tickers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		tickers, err := buildTickerParams("A", args, all)
		if err != nil {
			return err
		}
		return connectAndStream(cmd.Context(), "A", tickers, formatAggregate)
	},
}

// wsStocksLULDCmd streams Limit Up/Limit Down band data from the WebSocket API.
// Accepts one or more ticker symbols as arguments, or use --all to subscribe
// to all tickers. Outputs LULD events with symbol, upper limit, and lower limit.
// Usage: massive ws stocks luld AAPL --output table
var wsStocksLULDCmd = &cobra.Command{
	Use:   "luld [tickers...]",
	Short: "Stream Limit Up/Limit Down bands",
	Long:  "Stream real-time LULD (Limit Up/Limit Down) band data via WebSocket. Provide ticker symbols as arguments or use --all for all tickers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		tickers, err := buildTickerParams("LULD", args, all)
		if err != nil {
			return err
		}
		return connectAndStream(cmd.Context(), "LULD", tickers, formatLULD)
	},
}

// wsStocksFMVCmd streams Fair Market Value data from the WebSocket API.
// Accepts one or more ticker symbols as arguments, or use --all to subscribe
// to all tickers. Uses a special path (/business/stocks/FMV) different from
// other stock channels.
// Usage: massive ws stocks fmv AAPL --output table
var wsStocksFMVCmd = &cobra.Command{
	Use:   "fmv [tickers...]",
	Short: "Stream Fair Market Value data",
	Long:  "Stream real-time Fair Market Value (FMV) data via WebSocket. Provide ticker symbols as arguments or use --all for all tickers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		tickers, err := buildTickerParams("FMV", args, all)
		if err != nil {
			return err
		}
		return connectAndStream(cmd.Context(), "FMV", tickers, formatFMV)
	},
}

// tableFormatter is a function type that formats a single WebSocket event as a
// table row and writes it to the provided tabwriter. Each channel type (trades,
// quotes, aggregates, etc.) provides its own formatter implementation.
type tableFormatter func(w *tabwriter.Writer, event map[string]interface{})

// connectAndStream is a convenience wrapper that connects to the stocks WebSocket
// endpoint, authenticates, subscribes, and streams data. Used by all ws stocks subcommands.
func connectAndStream(parentCtx context.Context, channel, tickerParams string, formatter tableFormatter) error {
	return connectAndStreamAsset(parentCtx, "stocks", channel, tickerParams, formatter)
}

// buildTickerParams constructs the subscription parameter string for the given
// channel and ticker symbols. If the all flag is set, it returns a wildcard
// subscription (e.g., "T.*"). Otherwise it prefixes each ticker with the channel
// name (e.g., "T.AAPL,T.MSFT"). Returns an error if no tickers are provided
// and the --all flag is not set.
func buildTickerParams(channel string, args []string, all bool) (string, error) {
	if all {
		return channel + ".*", nil
	}
	if len(args) == 0 {
		return "", fmt.Errorf("provide at least one ticker symbol or use --all")
	}
	parts := make([]string, len(args))
	for i, t := range args {
		parts[i] = channel + "." + strings.ToUpper(t)
	}
	return strings.Join(parts, ","), nil
}

// getWSBaseURL returns the appropriate WebSocket base URL based on whether
// the --realtime flag is set. Returns the delayed endpoint by default.
func getWSBaseURL() string {
	if wsRealtime {
		return wsRealtimeURL
	}
	return wsDelayedURL
}

// buildWSURL constructs the full WebSocket connection URL for the given asset
// class. The channel is NOT included in the URL path — it is specified in
// the subscribe message after connecting and authenticating.
func buildWSURL(assetClass string) string {
	return getWSBaseURL() + "/" + assetClass
}

// connectAndStreamAsset establishes a WebSocket connection to the Massive streaming
// API for any asset class, authenticates, subscribes to the given tickers, and
// reads messages in a loop until the context is cancelled (e.g., via Ctrl+C).
// The assetClass parameter determines the WebSocket path (e.g., "stocks", "crypto").
func connectAndStreamAsset(parentCtx context.Context, assetClass, channel, tickerParams string, formatter tableFormatter) error {
	apiKey, err := config.GetAPIKey()
	if err != nil {
		return err
	}

	url := buildWSURL(assetClass)

	// Set up signal handling for clean shutdown on Ctrl+C.
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		select {
		case <-sigCh:
			cancel()
		case <-ctx.Done():
		}
	}()

	// Connect to the WebSocket endpoint.
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}
	defer conn.Close()

	// Read the initial "connected" status message from the server.
	_, _, err = conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("failed to read connection message: %w", err)
	}

	// Authenticate by sending the API key in an auth action message.
	authMsg := map[string]string{
		"action": "auth",
		"params": apiKey,
	}
	if err := conn.WriteJSON(authMsg); err != nil {
		return fmt.Errorf("failed to send auth message: %w", err)
	}

	// Read the auth response and verify authentication succeeded.
	_, authResp, err := conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("failed to read auth response: %w", err)
	}

	var authEvents []map[string]interface{}
	if err := json.Unmarshal(authResp, &authEvents); err == nil {
		for _, ev := range authEvents {
			if status, ok := ev["status"].(string); ok && status != "auth_success" {
				msg, _ := ev["message"].(string)
				return fmt.Errorf("authentication failed: %s", msg)
			}
		}
	}

	// Send the subscribe message to start receiving events for the requested tickers.
	subscribeMsg := map[string]string{
		"action": "subscribe",
		"params": tickerParams,
	}
	if err := conn.WriteJSON(subscribeMsg); err != nil {
		return fmt.Errorf("failed to send subscribe message: %w", err)
	}

	// Read the subscribe response to check for errors.
	_, subResp, err := conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("failed to read subscribe response: %w", err)
	}

	var subEvents []map[string]interface{}
	if err := json.Unmarshal(subResp, &subEvents); err == nil {
		for _, ev := range subEvents {
			if status, ok := ev["status"].(string); ok && status == "error" {
				msg, _ := ev["message"].(string)
				return fmt.Errorf("subscription failed: %s", msg)
			}
		}
	}

	fmt.Fprintf(os.Stderr, "Connected to %s/%s, subscribed to: %s\n", assetClass, channel, tickerParams)

	// Set up table writer for table output mode. The header is printed once
	// and flushed after each batch of events to keep output streaming smoothly.
	var w *tabwriter.Writer
	if outputFormat == "table" {
		w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		printTableHeader(w, channel)
		w.Flush()
	}

	// Close the connection when the context is cancelled (e.g., Ctrl+C).
	// This causes ReadMessage in the main loop to return an error and exit cleanly.
	go func() {
		<-ctx.Done()
		_ = conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		)
		conn.Close()
	}()

	// Read messages in a loop until the connection is closed.
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			// If context was cancelled, exit cleanly.
			if ctx.Err() != nil {
				fmt.Fprintln(os.Stderr, "\nDisconnected.")
				return nil
			}
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				fmt.Fprintln(os.Stderr, "\nServer closed connection.")
				return nil
			}
			return fmt.Errorf("error reading message: %w", err)
		}

		// Parse the message as a JSON array of events.
		var events []map[string]interface{}
		if err := json.Unmarshal(message, &events); err != nil {
			// Try parsing as a single object (e.g., status messages).
			var single map[string]interface{}
			if err2 := json.Unmarshal(message, &single); err2 == nil {
				events = []map[string]interface{}{single}
			} else {
				fmt.Fprintf(os.Stderr, "Warning: failed to parse message: %s\n", string(message))
				continue
			}
		}

		// Output each event based on the configured output format.
		for _, event := range events {
			if outputFormat == "json" {
				line, err := json.Marshal(event)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to marshal event: %v\n", err)
					continue
				}
				fmt.Println(string(line))
			} else {
				formatter(w, event)
			}
		}

		// Flush the table writer after each batch to stream output to the terminal.
		if outputFormat == "table" && w != nil {
			w.Flush()
		}
	}
}

// printTableHeader writes the column header row for the specified channel's
// table output format to the given tabwriter. Supports stock channels (T, Q, AM, A, LULD, FMV, V)
// and crypto channels (XT, XQ, XA, XAS, FMV).
func printTableHeader(w *tabwriter.Writer, channel string) {
	switch channel {
	// Stock channels
	case "T":
		fmt.Fprintln(w, "TIME\tSYMBOL\tPRICE\tSIZE\tEXCHANGE")
		fmt.Fprintln(w, "----\t------\t-----\t----\t--------")
	case "Q":
		fmt.Fprintln(w, "TIME\tSYMBOL\tBID\tBID_SIZE\tASK\tASK_SIZE")
		fmt.Fprintln(w, "----\t------\t---\t--------\t---\t--------")
	case "AM", "A":
		fmt.Fprintln(w, "TIME\tSYMBOL\tOPEN\tHIGH\tLOW\tCLOSE\tVOLUME")
		fmt.Fprintln(w, "----\t------\t----\t----\t---\t-----\t------")
	case "LULD":
		fmt.Fprintln(w, "TIME\tSYMBOL\tHIGH\tLOW")
		fmt.Fprintln(w, "----\t------\t----\t---")
	case "V":
		fmt.Fprintln(w, "TIME\tSYMBOL\tVALUE")
		fmt.Fprintln(w, "----\t------\t-----")
	// Crypto channels
	case "XT":
		fmt.Fprintln(w, "TIME\tPAIR\tPRICE\tSIZE\tEXCHANGE")
		fmt.Fprintln(w, "----\t----\t-----\t----\t--------")
	case "XQ":
		fmt.Fprintln(w, "TIME\tPAIR\tBID\tBID_SIZE\tASK\tASK_SIZE")
		fmt.Fprintln(w, "----\t----\t---\t--------\t---\t--------")
	case "XA", "XAS":
		fmt.Fprintln(w, "TIME\tPAIR\tOPEN\tHIGH\tLOW\tCLOSE\tVOLUME")
		fmt.Fprintln(w, "----\t----\t----\t----\t---\t-----\t------")
	// Forex channels
	case "C":
		fmt.Fprintln(w, "TIME\tPAIR\tBID\tASK")
		fmt.Fprintln(w, "----\t----\t---\t---")
	case "CA", "CAS":
		fmt.Fprintln(w, "TIME\tPAIR\tOPEN\tHIGH\tLOW\tCLOSE\tVOLUME")
		fmt.Fprintln(w, "----\t----\t----\t----\t---\t-----\t------")
	// FMV channel (used for stocks, crypto, and forex)
	case "FMV":
		fmt.Fprintln(w, "TIME\tSYMBOL\tFMV")
		fmt.Fprintln(w, "----\t------\t---")
	}
}

// formatTimestamp converts a millisecond Unix timestamp from a JSON number
// to a human-readable time string. Returns "N/A" if the value cannot be
// converted to a valid timestamp.
func formatTimestamp(v interface{}) string {
	switch val := v.(type) {
	case float64:
		t := time.Unix(0, int64(val)*int64(time.Millisecond))
		return t.Format("15:04:05.000")
	case json.Number:
		n, err := val.Int64()
		if err != nil {
			return "N/A"
		}
		t := time.Unix(0, n*int64(time.Millisecond))
		return t.Format("15:04:05.000")
	default:
		return "N/A"
	}
}

// getStr safely extracts a string value from a map by key. Returns an empty
// string if the key is missing or the value is not a string.
func getStr(event map[string]interface{}, key string) string {
	if v, ok := event[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// getFloat safely extracts a float64 value from a map by key. Returns 0.0
// if the key is missing or the value cannot be converted to a float64.
func getFloat(event map[string]interface{}, key string) float64 {
	if v, ok := event[key]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return 0.0
}

// formatTrade formats a single trade event as a table row showing time,
// symbol, price, size, and exchange. Writes the formatted row to the
// provided tabwriter.
func formatTrade(w *tabwriter.Writer, event map[string]interface{}) {
	ts := formatTimestamp(event["t"])
	sym := getStr(event, "sym")
	price := getFloat(event, "p")
	size := getFloat(event, "s")
	exchange := getFloat(event, "x")
	fmt.Fprintf(w, "%s\t%s\t%.4f\t%.0f\t%.0f\n", ts, sym, price, size, exchange)
}

// formatQuote formats a single NBBO quote event as a table row showing time,
// symbol, bid price, bid size, ask price, and ask size. Writes the formatted
// row to the provided tabwriter.
func formatQuote(w *tabwriter.Writer, event map[string]interface{}) {
	ts := formatTimestamp(event["t"])
	sym := getStr(event, "sym")
	bid := getFloat(event, "bp")
	bidSize := getFloat(event, "bs")
	ask := getFloat(event, "ap")
	askSize := getFloat(event, "as")
	fmt.Fprintf(w, "%s\t%s\t%.4f\t%.0f\t%.4f\t%.0f\n", ts, sym, bid, bidSize, ask, askSize)
}

// formatAggregate formats a single aggregate bar event (per-minute or
// per-second) as a table row showing time, symbol, open, high, low, close,
// and volume. Writes the formatted row to the provided tabwriter.
func formatAggregate(w *tabwriter.Writer, event map[string]interface{}) {
	ts := formatTimestamp(event["s"])
	sym := getStr(event, "sym")
	open := getFloat(event, "o")
	high := getFloat(event, "h")
	low := getFloat(event, "l")
	close_ := getFloat(event, "c")
	volume := getFloat(event, "v")
	fmt.Fprintf(w, "%s\t%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\n", ts, sym, open, high, low, close_, volume)
}

// formatLULD formats a single Limit Up/Limit Down event as a table row
// showing time, symbol, upper limit (high), and lower limit (low). Writes
// the formatted row to the provided tabwriter.
func formatLULD(w *tabwriter.Writer, event map[string]interface{}) {
	ts := formatTimestamp(event["t"])
	sym := getStr(event, "T")
	high := getFloat(event, "h")
	low := getFloat(event, "l")
	fmt.Fprintf(w, "%s\t%s\t%.4f\t%.4f\n", ts, sym, high, low)
}

// formatFMV formats a single Fair Market Value event as a table row showing
// time, symbol, and FMV price. Writes the formatted row to the provided
// tabwriter.
func formatFMV(w *tabwriter.Writer, event map[string]interface{}) {
	ts := formatTimestamp(event["t"])
	sym := getStr(event, "sym")
	fmv := getFloat(event, "fmv")
	fmt.Fprintf(w, "%s\t%s\t%.4f\n", ts, sym, fmv)
}

// init registers the ws command under rootCmd, the stocks command under ws,
// and all stock streaming subcommands under ws stocks. Each subcommand gets
// an --all flag to subscribe to all available tickers.
func init() {
	// Register parent commands.
	rootCmd.AddCommand(wsCmd)
	wsCmd.PersistentFlags().BoolVar(&wsRealtime, "realtime", false, "Connect to real-time endpoint instead of delayed (15-min)")
	wsCmd.AddCommand(wsStocksCmd)

	// Add --all flag to each subcommand.
	wsStocksTradesCmd.Flags().Bool("all", false, "Subscribe to all tickers")
	wsStocksQuotesCmd.Flags().Bool("all", false, "Subscribe to all tickers")
	wsStocksAggMinuteCmd.Flags().Bool("all", false, "Subscribe to all tickers")
	wsStocksAggSecondCmd.Flags().Bool("all", false, "Subscribe to all tickers")
	wsStocksLULDCmd.Flags().Bool("all", false, "Subscribe to all tickers")
	wsStocksFMVCmd.Flags().Bool("all", false, "Subscribe to all tickers")

	// Register all subcommands under ws stocks.
	wsStocksCmd.AddCommand(wsStocksTradesCmd)
	wsStocksCmd.AddCommand(wsStocksQuotesCmd)
	wsStocksCmd.AddCommand(wsStocksAggMinuteCmd)
	wsStocksCmd.AddCommand(wsStocksAggSecondCmd)
	wsStocksCmd.AddCommand(wsStocksLULDCmd)
	wsStocksCmd.AddCommand(wsStocksFMVCmd)
}
