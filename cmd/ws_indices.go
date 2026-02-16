//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package cmd

import (
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

// wsIndicesCmd is the parent command for all WebSocket-based indices streaming
// subcommands. It groups real-time index data streams such as per-minute
// aggregates, per-second aggregates, and live index values.
var wsIndicesCmd = &cobra.Command{
	Use:   "indices",
	Short: "Stream real-time indices data",
}

// wsIndicesAggMinuteCmd streams per-minute aggregate bars for one or more
// index tickers via WebSocket. Each message contains the symbol, open, close,
// high, low, and start/end timestamps for the aggregate window.
// Usage: massive ws indices agg-minute I:SPX I:DJI
var wsIndicesAggMinuteCmd = &cobra.Command{
	Use:   "agg-minute [tickers...]",
	Short: "Stream per-minute aggregate bars for indices",
	Long:  "Connect to the Massive WebSocket API and stream real-time per-minute OHLC aggregate data for one or more index tickers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")

		tickers, err := buildWsIndicesTickers(args, all, "AM")
		if err != nil {
			return err
		}

		return streamIndicesWebSocket("AM", tickers, formatIndicesAggEvent)
	},
}

// wsIndicesAggSecondCmd streams per-second aggregate bars for one or more
// index tickers via WebSocket. Each message contains the symbol, open, close,
// high, low, and start/end timestamps for the aggregate window.
// Usage: massive ws indices agg-second I:SPX
var wsIndicesAggSecondCmd = &cobra.Command{
	Use:   "agg-second [tickers...]",
	Short: "Stream per-second aggregate bars for indices",
	Long:  "Connect to the Massive WebSocket API and stream real-time per-second OHLC aggregate data for one or more index tickers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")

		tickers, err := buildWsIndicesTickers(args, all, "A")
		if err != nil {
			return err
		}

		return streamIndicesWebSocket("A", tickers, formatIndicesAggEvent)
	},
}

// wsIndicesValueCmd streams real-time index values for one or more index
// tickers via WebSocket. Each message contains the ticker symbol, the
// current index value, and a timestamp.
// Usage: massive ws indices value I:SPX I:COMP
var wsIndicesValueCmd = &cobra.Command{
	Use:   "value [tickers...]",
	Short: "Stream real-time index values",
	Long:  "Connect to the Massive WebSocket API and stream real-time index value updates for one or more index tickers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")

		tickers, err := buildWsIndicesTickers(args, all, "V")
		if err != nil {
			return err
		}

		return streamIndicesWebSocket("V", tickers, formatIndicesValueEvent)
	},
}

// indicesAggEvent represents a per-minute or per-second aggregate event
// received from the indices WebSocket stream. Fields match the wire format
// with short JSON keys (ev, sym, op, o, c, h, l, s, e).
type indicesAggEvent struct {
	Ev     string  `json:"ev"`
	Symbol string  `json:"sym"`
	OpenP  float64 `json:"op"`
	Open   float64 `json:"o"`
	Close  float64 `json:"c"`
	High   float64 `json:"h"`
	Low    float64 `json:"l"`
	Start  int64   `json:"s"`
	End    int64   `json:"e"`
}

// indicesValueEvent represents a real-time index value event received from
// the indices WebSocket stream. Contains the event type, index value,
// ticker symbol, and a timestamp.
type indicesValueEvent struct {
	Ev     string  `json:"ev"`
	Value  float64 `json:"val"`
	Ticker string  `json:"T"`
	Time   int64   `json:"t"`
}

// indicesEventFormatter is a function type that formats a raw JSON event
// message for display. It handles both JSON and table output modes depending
// on the global outputFormat setting.
type indicesEventFormatter func(w *tabwriter.Writer, raw json.RawMessage) error

// buildWsIndicesTickers validates and constructs the subscription parameter
// string for the indices WebSocket connection. When the --all flag is set,
// it subscribes to all tickers using a wildcard. Otherwise, it requires at
// least one ticker symbol to be provided as a positional argument. The
// channel prefix (e.g., "AM", "A", "V") is prepended to each ticker.
func buildWsIndicesTickers(args []string, all bool, channel string) (string, error) {
	if all {
		return channel + ".*", nil
	}

	if len(args) == 0 {
		return "", fmt.Errorf("at least one ticker symbol is required (e.g., I:SPX) or use --all")
	}

	params := make([]string, len(args))
	for i, t := range args {
		params[i] = channel + "." + strings.ToUpper(t)
	}

	return strings.Join(params, ","), nil
}

// streamIndicesWebSocket establishes a WebSocket connection to the Massive
// indices streaming endpoint, subscribes to the specified tickers, and reads
// messages in a loop. Each received message is parsed as a JSON array of
// events and formatted using the provided formatter function. The connection
// is cleanly shut down when the user presses Ctrl+C (SIGINT/SIGTERM).
func streamIndicesWebSocket(channel string, subscribeParams string, formatter indicesEventFormatter) error {
	apiKey, err := config.GetAPIKey()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("wss://socket.massive.com/indices/%s?apiKey=%s", channel, apiKey)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}
	defer conn.Close()

	// Send the subscribe message to begin receiving events for the
	// specified tickers.
	subscribeMsg := map[string]string{
		"action": "subscribe",
		"params": subscribeParams,
	}

	if err := conn.WriteJSON(subscribeMsg); err != nil {
		return fmt.Errorf("failed to send subscribe message: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Connected to indices/%s stream. Subscribed to: %s\n", channel, subscribeParams)
	fmt.Fprintf(os.Stderr, "Press Ctrl+C to stop.\n\n")

	// Set up signal handling for clean shutdown on Ctrl+C.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Create a tabwriter for table-formatted output. The writer is flushed
	// after each batch of events to ensure timely display.
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Print the table header if output format is table.
	if outputFormat != "json" {
		printIndicesStreamHeader(w, channel)
		w.Flush()
	}

	// done channel signals when the read loop has finished, either due to
	// an error or a clean connection close.
	done := make(chan struct{})

	go func() {
		defer close(done)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				// Check if the error is due to a normal close; if so, exit
				// silently. Otherwise, report the error.
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					return
				}
				fmt.Fprintf(os.Stderr, "\nWebSocket read error: %v\n", err)
				return
			}

			// Messages arrive as a JSON array of event objects.
			var events []json.RawMessage
			if err := json.Unmarshal(message, &events); err != nil {
				// If the message is not an array, try to print it as a
				// single raw JSON line (e.g., status messages).
				if outputFormat == "json" {
					fmt.Println(string(message))
				}
				continue
			}

			for _, raw := range events {
				if err := formatter(w, raw); err != nil {
					fmt.Fprintf(os.Stderr, "Format error: %v\n", err)
				}
			}

			if outputFormat != "json" {
				w.Flush()
			}
		}
	}()

	// Block until either the interrupt signal is received or the read loop
	// finishes on its own.
	select {
	case <-interrupt:
		fmt.Fprintf(os.Stderr, "\nShutting down...\n")

		// Send a close message to the server and wait briefly for a
		// clean close handshake.
		closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
		if err := conn.WriteMessage(websocket.CloseMessage, closeMsg); err != nil {
			return fmt.Errorf("failed to send close message: %w", err)
		}

		// Wait for the read loop to finish or timeout after 3 seconds.
		select {
		case <-done:
		case <-time.After(3 * time.Second):
		}

	case <-done:
	}

	return nil
}

// printIndicesStreamHeader writes the column header row for table-formatted
// output. The header varies depending on the channel: aggregate channels
// (AM, A) show OHLC data while the value channel (V) shows the index value.
func printIndicesStreamHeader(w *tabwriter.Writer, channel string) {
	if channel == "V" {
		fmt.Fprintln(w, "TIME\tSYMBOL\tVALUE")
		fmt.Fprintln(w, "----\t------\t-----")
	} else {
		fmt.Fprintln(w, "TIME\tSYMBOL\tOPEN\tHIGH\tLOW\tCLOSE")
		fmt.Fprintln(w, "----\t------\t----\t----\t---\t-----")
	}
}

// formatIndicesAggEvent formats a single aggregate event (AM or A) for
// display. When output format is JSON, the raw event is printed as-is.
// When output format is table, the event is parsed into an indicesAggEvent
// struct and displayed as a formatted table row with timestamp, symbol,
// and OHLC values.
func formatIndicesAggEvent(w *tabwriter.Writer, raw json.RawMessage) error {
	if outputFormat == "json" {
		fmt.Println(string(raw))
		return nil
	}

	var evt indicesAggEvent
	if err := json.Unmarshal(raw, &evt); err != nil {
		return fmt.Errorf("failed to parse aggregate event: %w", err)
	}

	// Skip status messages that don't have a symbol.
	if evt.Symbol == "" {
		return nil
	}

	t := time.UnixMilli(evt.Start)
	fmt.Fprintf(w, "%s\t%s\t%.4f\t%.4f\t%.4f\t%.4f\n",
		t.Format("15:04:05"),
		evt.Symbol,
		evt.Open,
		evt.High,
		evt.Low,
		evt.Close,
	)

	return nil
}

// formatIndicesValueEvent formats a single index value event (V) for
// display. When output format is JSON, the raw event is printed as-is.
// When output format is table, the event is parsed into an indicesValueEvent
// struct and displayed as a formatted table row with timestamp, ticker
// symbol, and the current index value.
func formatIndicesValueEvent(w *tabwriter.Writer, raw json.RawMessage) error {
	if outputFormat == "json" {
		fmt.Println(string(raw))
		return nil
	}

	var evt indicesValueEvent
	if err := json.Unmarshal(raw, &evt); err != nil {
		return fmt.Errorf("failed to parse value event: %w", err)
	}

	// Skip status messages that don't have a ticker.
	if evt.Ticker == "" {
		return nil
	}

	t := time.UnixMilli(evt.Time)
	fmt.Fprintf(w, "%s\t%s\t%.4f\n",
		t.Format("15:04:05"),
		evt.Ticker,
		evt.Value,
	)

	return nil
}

// init registers the indices WebSocket streaming subcommands under the
// wsIndicesCmd parent, adds the --all flag to each subcommand, and registers
// wsIndicesCmd under the shared wsCmd parent command.
func init() {
	wsIndicesAggMinuteCmd.Flags().Bool("all", false, "Subscribe to all index tickers")
	wsIndicesAggSecondCmd.Flags().Bool("all", false, "Subscribe to all index tickers")
	wsIndicesValueCmd.Flags().Bool("all", false, "Subscribe to all index tickers")

	wsIndicesCmd.AddCommand(wsIndicesAggMinuteCmd)
	wsIndicesCmd.AddCommand(wsIndicesAggSecondCmd)
	wsIndicesCmd.AddCommand(wsIndicesValueCmd)

	wsCmd.AddCommand(wsIndicesCmd)
}
