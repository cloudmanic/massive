//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/cloudmanic/massive-cli/internal/config"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

// wsOptionsCmd is the parent command for all options WebSocket streaming
// subcommands including trades, quotes, aggregates, and fair market value.
var wsOptionsCmd = &cobra.Command{
	Use:   "options",
	Short: "Stream real-time options data",
}

// wsOptionsTradesCmd streams real-time options trade events over a WebSocket
// connection. Accepts one or more option contract tickers as positional
// arguments or the --all flag to subscribe to all trades. Each trade event
// includes the symbol, price, size, exchange, conditions, timestamp, and
// sequence number.
// Usage: massive ws options trades O:SPY241220P00720000
var wsOptionsTradesCmd = &cobra.Command{
	Use:   "trades [tickers...]",
	Short: "Stream real-time options trades",
	Long:  "Stream real-time options trade data via WebSocket. Provide one or more option contract tickers (e.g., O:SPY241220P00720000) or use --all to subscribe to all trades.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		tickers := args
		if all {
			tickers = []string{"*"}
		}

		// Build subscription params with T. prefix for trades channel
		params := buildOptionsSubscriptionParams("T", tickers)

		return connectAndStreamOptions("T", params, func(msg map[string]interface{}) {
			if outputFormat == "json" {
				printJSON(msg)
				return
			}
			printOptionsTradeRow(msg)
		})
	},
}

// wsOptionsQuotesCmd streams real-time options NBBO quote events over a
// WebSocket connection. Maximum of 1000 contracts per connection. Accepts
// one or more option contract tickers as positional arguments or the --all
// flag. Each quote event includes the symbol, bid/ask exchanges, bid/ask
// prices, bid/ask sizes, timestamp, and sequence number.
// Usage: massive ws options quotes O:SPY241220P00720000
var wsOptionsQuotesCmd = &cobra.Command{
	Use:   "quotes [tickers...]",
	Short: "Stream real-time options quotes (max 1000 contracts/connection)",
	Long:  "Stream real-time options NBBO quote data via WebSocket. Maximum 1000 contracts per connection. Provide option contract tickers or use --all.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		tickers := args
		if all {
			tickers = []string{"*"}
		}

		// Build subscription params with Q. prefix for quotes channel
		params := buildOptionsSubscriptionParams("Q", tickers)

		return connectAndStreamOptions("Q", params, func(msg map[string]interface{}) {
			if outputFormat == "json" {
				printJSON(msg)
				return
			}
			printOptionsQuoteRow(msg)
		})
	},
}

// wsOptionsAggMinuteCmd streams per-minute aggregate bar events for options
// contracts over a WebSocket connection. Accepts one or more option contract
// tickers as positional arguments or the --all flag. Each aggregate event
// includes OHLCV data, VWAP, average price, and the bar time window.
// Usage: massive ws options agg-minute O:SPY241220P00720000
var wsOptionsAggMinuteCmd = &cobra.Command{
	Use:   "agg-minute [tickers...]",
	Short: "Stream per-minute options aggregates",
	Long:  "Stream per-minute aggregate bar data for options contracts via WebSocket. Includes open, high, low, close, volume, and VWAP.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		tickers := args
		if all {
			tickers = []string{"*"}
		}

		// Build subscription params with AM. prefix for per-minute aggregates
		params := buildOptionsSubscriptionParams("AM", tickers)

		return connectAndStreamOptions("AM", params, func(msg map[string]interface{}) {
			if outputFormat == "json" {
				printJSON(msg)
				return
			}
			printOptionsAggRow(msg)
		})
	},
}

// wsOptionsAggSecondCmd streams per-second aggregate bar events for options
// contracts over a WebSocket connection. Accepts one or more option contract
// tickers as positional arguments or the --all flag. Each aggregate event
// includes OHLCV data, VWAP, average price, and the bar time window.
// Usage: massive ws options agg-second O:SPY241220P00720000
var wsOptionsAggSecondCmd = &cobra.Command{
	Use:   "agg-second [tickers...]",
	Short: "Stream per-second options aggregates",
	Long:  "Stream per-second aggregate bar data for options contracts via WebSocket. Includes open, high, low, close, volume, and VWAP.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		tickers := args
		if all {
			tickers = []string{"*"}
		}

		// Build subscription params with A. prefix for per-second aggregates
		params := buildOptionsSubscriptionParams("A", tickers)

		return connectAndStreamOptions("A", params, func(msg map[string]interface{}) {
			if outputFormat == "json" {
				printJSON(msg)
				return
			}
			printOptionsAggRow(msg)
		})
	},
}

// wsOptionsFMVCmd streams Fair Market Value events for options contracts
// over a WebSocket connection. Uses the /business/options/FMV path. Accepts
// one or more option contract tickers as positional arguments or the --all
// flag. Each FMV event includes the fair market value, symbol, and timestamp.
// Usage: massive ws options fmv O:SPY241220P00720000
var wsOptionsFMVCmd = &cobra.Command{
	Use:   "fmv [tickers...]",
	Short: "Stream options Fair Market Value",
	Long:  "Stream Fair Market Value (FMV) data for options contracts via WebSocket. Uses the /business/options/FMV endpoint.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		tickers := args
		if all {
			tickers = []string{"*"}
		}

		// Build subscription params with FMV. prefix for fair market value
		params := buildOptionsSubscriptionParams("FMV", tickers)

		return connectAndStreamOptions("FMV", params, func(msg map[string]interface{}) {
			if outputFormat == "json" {
				printJSON(msg)
				return
			}
			printOptionsFMVRow(msg)
		})
	},
}

// buildOptionsSubscriptionParams constructs the subscription parameter string
// for a WebSocket subscribe message. Each ticker is prefixed with the channel
// name and a dot separator (e.g., "T.O:SPY241220P00720000"). Multiple tickers
// are joined with commas.
func buildOptionsSubscriptionParams(channel string, tickers []string) string {
	parts := make([]string, len(tickers))
	for i, t := range tickers {
		parts[i] = channel + "." + t
	}
	return strings.Join(parts, ",")
}

// connectAndStreamOptions establishes a WebSocket connection to the Massive
// options streaming endpoint, sends a subscription message, and reads incoming
// messages in a loop. The handler function is called for each received event.
// The connection is cleanly shut down when the user presses Ctrl+C. The channel
// parameter determines the WebSocket path: FMV uses /business/options/FMV while
// all other channels use /options/{channel}.
func connectAndStreamOptions(channel string, subscriptionParams string, handler func(map[string]interface{})) error {
	// Get the API key for authentication
	apiKey, err := config.GetAPIKey()
	if err != nil {
		return err
	}

	// Build the WebSocket URL based on the channel type
	var wsPath string
	if channel == "FMV" {
		wsPath = "/business/options/FMV"
	} else {
		wsPath = "/options/" + channel
	}

	wsURL := url.URL{
		Scheme:   "wss",
		Host:     "socket.massive.com",
		Path:     wsPath,
		RawQuery: "apiKey=" + url.QueryEscape(apiKey),
	}

	// Set up context with cancellation for clean shutdown on Ctrl+C
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		fmt.Println("\nShutting down...")
		cancel()
	}()

	// Establish the WebSocket connection
	fmt.Printf("Connecting to %s...\n", wsURL.String())
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, wsURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}
	defer conn.Close()

	fmt.Println("Connected. Subscribing...")

	// Send the subscription message to start receiving events
	subscribeMsg := map[string]string{
		"action": "subscribe",
		"params": subscriptionParams,
	}

	if err := conn.WriteJSON(subscribeMsg); err != nil {
		return fmt.Errorf("failed to send subscribe message: %w", err)
	}

	fmt.Printf("Subscribed to: %s\n\n", subscriptionParams)

	// Print table header if using table output format
	if outputFormat != "json" {
		printOptionsTableHeader(channel)
	}

	// Start a goroutine to close the connection when context is cancelled
	go func() {
		<-ctx.Done()
		conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		)
	}()

	// Read messages in a loop until the connection closes or context is cancelled
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				if ctx.Err() != nil {
					return nil
				}
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					return nil
				}
				log.Printf("Read error: %v", err)
				return fmt.Errorf("WebSocket read error: %w", err)
			}

			// Parse the received message as a JSON array of events
			var events []map[string]interface{}
			if err := json.Unmarshal(message, &events); err != nil {
				// Try parsing as a single object if array parsing fails
				var single map[string]interface{}
				if err2 := json.Unmarshal(message, &single); err2 != nil {
					log.Printf("Failed to parse message: %s", string(message))
					continue
				}
				events = []map[string]interface{}{single}
			}

			// Process each event through the handler function
			for _, event := range events {
				handler(event)
			}
		}
	}
}

// printOptionsTableHeader prints the column header row for the table output
// format based on the options channel type. Each channel has different fields
// that are displayed as columns.
func printOptionsTableHeader(channel string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	switch channel {
	case "T":
		fmt.Fprintln(w, "TIME\tSYMBOL\tPRICE\tSIZE\tEXCHANGE")
		fmt.Fprintln(w, "----\t------\t-----\t----\t--------")
	case "Q":
		fmt.Fprintln(w, "TIME\tSYMBOL\tBID\tBID_SIZE\tASK\tASK_SIZE")
		fmt.Fprintln(w, "----\t------\t---\t--------\t---\t--------")
	case "AM", "A":
		fmt.Fprintln(w, "TIME\tSYMBOL\tOPEN\tHIGH\tLOW\tCLOSE\tVOLUME")
		fmt.Fprintln(w, "----\t------\t----\t----\t---\t-----\t------")
	case "FMV":
		fmt.Fprintln(w, "TIME\tSYMBOL\tFMV")
		fmt.Fprintln(w, "----\t------\t---")
	}
	w.Flush()
}

// printOptionsTradeRow formats and prints a single options trade event as
// a tab-separated table row. Extracts the timestamp, symbol, price, size,
// and exchange fields from the event map.
func printOptionsTradeRow(msg map[string]interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	ts := formatOptionsTimestamp(msg["t"])
	sym := safeOptionsString(msg["sym"])
	price := safeOptionsFloat(msg["p"])
	size := safeOptionsFloat(msg["s"])
	exchange := safeOptionsString(msg["x"])
	fmt.Fprintf(w, "%s\t%s\t%.4f\t%.0f\t%s\n", ts, sym, price, size, exchange)
	w.Flush()
}

// printOptionsQuoteRow formats and prints a single options quote event as
// a tab-separated table row. Extracts the timestamp, symbol, bid price,
// bid size, ask price, and ask size fields from the event map.
func printOptionsQuoteRow(msg map[string]interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	ts := formatOptionsTimestamp(msg["t"])
	sym := safeOptionsString(msg["sym"])
	bid := safeOptionsFloat(msg["bp"])
	bidSize := safeOptionsFloat(msg["bs"])
	ask := safeOptionsFloat(msg["ap"])
	askSize := safeOptionsFloat(msg["as"])
	fmt.Fprintf(w, "%s\t%s\t%.4f\t%.0f\t%.4f\t%.0f\n", ts, sym, bid, bidSize, ask, askSize)
	w.Flush()
}

// printOptionsAggRow formats and prints a single options aggregate event as
// a tab-separated table row. Used for both per-minute (AM) and per-second (A)
// aggregate channels. Extracts the start time, symbol, open, high, low, close,
// and volume fields from the event map.
func printOptionsAggRow(msg map[string]interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	ts := formatOptionsTimestamp(msg["s"])
	sym := safeOptionsString(msg["sym"])
	open := safeOptionsFloat(msg["o"])
	high := safeOptionsFloat(msg["h"])
	low := safeOptionsFloat(msg["l"])
	close := safeOptionsFloat(msg["c"])
	volume := safeOptionsFloat(msg["v"])
	fmt.Fprintf(w, "%s\t%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.0f\n", ts, sym, open, high, low, close, volume)
	w.Flush()
}

// printOptionsFMVRow formats and prints a single options Fair Market Value
// event as a tab-separated table row. Extracts the timestamp, symbol, and
// FMV fields from the event map.
func printOptionsFMVRow(msg map[string]interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	ts := formatOptionsTimestamp(msg["t"])
	sym := safeOptionsString(msg["sym"])
	fmv := safeOptionsFloat(msg["fmv"])
	fmt.Fprintf(w, "%s\t%s\t%.4f\n", ts, sym, fmv)
	w.Flush()
}

// formatOptionsTimestamp converts a timestamp value from a WebSocket event
// into a human-readable time string. Handles both millisecond and nanosecond
// epoch timestamps by checking the magnitude. Returns "N/A" if the value
// cannot be converted to a number.
func formatOptionsTimestamp(v interface{}) string {
	if v == nil {
		return "N/A"
	}

	var ms int64
	switch val := v.(type) {
	case float64:
		// Determine if the timestamp is in milliseconds or nanoseconds
		if val > 1e15 {
			// Nanosecond timestamp - convert to milliseconds
			ms = int64(val) / 1e6
		} else if val > 1e12 {
			// Microsecond timestamp - convert to milliseconds
			ms = int64(val) / 1e3
		} else {
			// Already in milliseconds
			ms = int64(val)
		}
	case json.Number:
		n, err := val.Int64()
		if err != nil {
			return "N/A"
		}
		if n > 1e15 {
			ms = n / 1e6
		} else if n > 1e12 {
			ms = n / 1e3
		} else {
			ms = n
		}
	default:
		return "N/A"
	}

	t := time.Unix(0, ms*int64(time.Millisecond))
	return t.Format("15:04:05.000")
}

// safeOptionsString extracts a string value from an interface{}, returning
// an empty string if the value is nil or not a string type. Used to safely
// access string fields from parsed JSON event maps.
func safeOptionsString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

// safeOptionsFloat extracts a float64 value from an interface{}, returning
// 0 if the value is nil or not a numeric type. Used to safely access numeric
// fields from parsed JSON event maps.
func safeOptionsFloat(v interface{}) float64 {
	if v == nil {
		return 0
	}
	if f, ok := v.(float64); ok {
		return f
	}
	return 0
}

// init registers the options WebSocket command and all its subcommands under
// the ws parent command. Each subcommand gets an --all flag for subscribing
// to all available option contract events.
func init() {
	// Add the --all flag to each subcommand
	wsOptionsTradesCmd.Flags().Bool("all", false, "Subscribe to all options trades")
	wsOptionsQuotesCmd.Flags().Bool("all", false, "Subscribe to all options quotes")
	wsOptionsAggMinuteCmd.Flags().Bool("all", false, "Subscribe to all options per-minute aggregates")
	wsOptionsAggSecondCmd.Flags().Bool("all", false, "Subscribe to all options per-second aggregates")
	wsOptionsFMVCmd.Flags().Bool("all", false, "Subscribe to all options FMV events")

	// Register subcommands under the options WebSocket parent
	wsOptionsCmd.AddCommand(wsOptionsTradesCmd)
	wsOptionsCmd.AddCommand(wsOptionsQuotesCmd)
	wsOptionsCmd.AddCommand(wsOptionsAggMinuteCmd)
	wsOptionsCmd.AddCommand(wsOptionsAggSecondCmd)
	wsOptionsCmd.AddCommand(wsOptionsFMVCmd)

	// Register the options command under the ws parent command
	wsCmd.AddCommand(wsOptionsCmd)
}
