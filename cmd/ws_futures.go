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
	"os"
	"os/signal"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/cloudmanic/massive-cli/internal/config"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

// wsFuturesCmd is the parent command for all futures WebSocket streaming
// subcommands. It groups trades, quotes, and aggregate streams under
// "massive ws futures".
var wsFuturesCmd = &cobra.Command{
	Use:   "futures",
	Short: "Stream real-time futures data via WebSocket",
}

// wsFuturesTradesCmd streams real-time trade data for the specified futures
// contracts. Each event includes the symbol, price, and size of the trade.
// Usage: massive ws futures trades ESZ4 NQZ4
var wsFuturesTradesCmd = &cobra.Command{
	Use:   "trades [tickers...]",
	Short: "Stream real-time futures trades",
	Long:  "Stream real-time trade data for specified futures contracts including price, size, and timestamps.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		tickers := args

		if !all && len(tickers) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		subscriptions := buildFuturesSubscriptions("T", tickers, all)
		return runFuturesWebSocket("T", subscriptions, renderFuturesTrade)
	},
}

// wsFuturesQuotesCmd streams real-time quote data for the specified futures
// contracts. Each event includes bid/ask prices and sizes.
// Usage: massive ws futures quotes ESZ4
var wsFuturesQuotesCmd = &cobra.Command{
	Use:   "quotes [tickers...]",
	Short: "Stream real-time futures quotes",
	Long:  "Stream real-time quote data for specified futures contracts including bid/ask prices and sizes.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		tickers := args

		if !all && len(tickers) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		subscriptions := buildFuturesSubscriptions("Q", tickers, all)
		return runFuturesWebSocket("Q", subscriptions, renderFuturesQuote)
	},
}

// wsFuturesAggMinuteCmd streams per-minute aggregate bars for the specified
// futures contracts. Each event includes OHLCV data for the completed bar.
// Usage: massive ws futures agg-minute ESZ4
var wsFuturesAggMinuteCmd = &cobra.Command{
	Use:   "agg-minute [tickers...]",
	Short: "Stream per-minute futures aggregates",
	Long:  "Stream per-minute aggregate bar data (OHLCV) for specified futures contracts.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		tickers := args

		if !all && len(tickers) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		subscriptions := buildFuturesSubscriptions("AM", tickers, all)
		return runFuturesWebSocket("AM", subscriptions, renderFuturesAggregate)
	},
}

// wsFuturesAggSecondCmd streams per-second aggregate bars for the specified
// futures contracts. Each event includes OHLCV data for the completed bar.
// Usage: massive ws futures agg-second ESZ4
var wsFuturesAggSecondCmd = &cobra.Command{
	Use:   "agg-second [tickers...]",
	Short: "Stream per-second futures aggregates",
	Long:  "Stream per-second aggregate bar data (OHLCV) for specified futures contracts.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		tickers := args

		if !all && len(tickers) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		subscriptions := buildFuturesSubscriptions("A", tickers, all)
		return runFuturesWebSocket("A", subscriptions, renderFuturesAggregate)
	},
}

// futuresTradeEvent represents a single futures trade message received from the
// WebSocket. Fields map to the API's event schema for the "T" channel.
type futuresTradeEvent struct {
	Ev        string  `json:"ev"`
	Symbol    string  `json:"sym"`
	Price     float64 `json:"p"`
	Size      float64 `json:"s"`
	Timestamp int64   `json:"t"`
	Sequence  int64   `json:"q"`
	Exchange  int     `json:"z"`
}

// futuresQuoteEvent represents a single futures quote message received from the
// WebSocket. Fields map to the API's event schema for the "Q" channel.
type futuresQuoteEvent struct {
	Ev           string  `json:"ev"`
	Symbol       string  `json:"sym"`
	BidPrice     float64 `json:"bp"`
	BidSize      float64 `json:"bs"`
	BidTimestamp int64   `json:"bt"`
	AskPrice     float64 `json:"ap"`
	AskSize      float64 `json:"as"`
	AskTimestamp int64   `json:"at"`
	Timestamp    int64   `json:"t"`
}

// futuresAggregateEvent represents a futures aggregate bar message received
// from the WebSocket. Used for both per-minute (AM) and per-second (A) channels.
type futuresAggregateEvent struct {
	Ev          string  `json:"ev"`
	Symbol      string  `json:"sym"`
	Open        float64 `json:"o"`
	High        float64 `json:"h"`
	Low         float64 `json:"l"`
	Close       float64 `json:"c"`
	Volume      float64 `json:"v"`
	DollarVol   float64 `json:"dv"`
	NumTrades   int64   `json:"n"`
	Start       int64   `json:"s"`
	End         int64   `json:"e"`
}

// buildFuturesSubscriptions constructs the subscription parameter string for the
// given channel and tickers. When all is true, it subscribes to the wildcard
// pattern for that channel. Otherwise, it prefixes each ticker with the channel
// name separated by a dot (e.g., "T.ESZ4").
func buildFuturesSubscriptions(channel string, tickers []string, all bool) string {
	if all {
		return channel + ".*"
	}

	parts := make([]string, len(tickers))
	for i, t := range tickers {
		parts[i] = channel + "." + strings.ToUpper(t)
	}
	return strings.Join(parts, ",")
}

// connectFuturesWebSocket establishes a WebSocket connection to the Massive
// futures streaming endpoint for the given channel. It loads the API key from
// the environment or config file and dials the WebSocket URL. Returns the open
// connection or an error if the connection fails.
func connectFuturesWebSocket(channel string) (*websocket.Conn, error) {
	apiKey, err := config.GetAPIKey()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("wss://socket.massive.com/futures/%s?apiKey=%s", channel, apiKey)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to futures WebSocket: %w", err)
	}

	return conn, nil
}

// sendFuturesSubscription sends a subscribe action over the WebSocket connection
// with the provided params string. This tells the server which tickers and
// channels to stream data for.
func sendFuturesSubscription(conn *websocket.Conn, params string) error {
	msg := map[string]string{
		"action": "subscribe",
		"params": params,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription message: %w", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("failed to send subscription: %w", err)
	}

	return nil
}

// runFuturesWebSocket is the main streaming loop for futures WebSocket commands.
// It connects to the WebSocket, sends the subscription, and continuously reads
// messages until the user presses Ctrl+C. Each raw message is passed to the
// provided render function for display. The render function receives the raw
// JSON bytes and a tabwriter for table output, and returns true if a row was
// written that should be flushed.
func runFuturesWebSocket(channel string, subscriptions string, render func([]byte, *tabwriter.Writer) bool) error {
	conn, err := connectFuturesWebSocket(channel)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := sendFuturesSubscription(conn, subscriptions); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Connected to futures %s stream. Press Ctrl+C to stop.\n", channel)

	// Set up signal handling for clean shutdown on Ctrl+C.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var w *tabwriter.Writer
	if outputFormat == "table" {
		w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	}

	// Read messages in a goroutine so we can respect context cancellation.
	msgChan := make(chan []byte)
	errChan := make(chan error)

	go futuresReadLoop(conn, msgChan, errChan)

	for {
		select {
		case <-ctx.Done():
			fmt.Fprintln(os.Stderr, "\nDisconnecting...")
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			return nil

		case err := <-errChan:
			return fmt.Errorf("WebSocket read error: %w", err)

		case msg := <-msgChan:
			if outputFormat == "json" {
				fmt.Println(string(msg))
			} else {
				if render(msg, w) {
					w.Flush()
				}
			}
		}
	}
}

// futuresReadLoop continuously reads messages from the WebSocket connection and
// sends them to the provided message channel. If a read error occurs, it is
// sent to the error channel and the loop exits.
func futuresReadLoop(conn *websocket.Conn, msgChan chan<- []byte, errChan chan<- error) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			errChan <- err
			return
		}
		msgChan <- msg
	}
}

// renderFuturesTrade parses a raw WebSocket message as a futures trade event
// array and writes each trade as a table row. Returns true if any rows were
// written so the caller knows to flush the tabwriter.
func renderFuturesTrade(msg []byte, w *tabwriter.Writer) bool {
	var events []futuresTradeEvent
	if err := json.Unmarshal(msg, &events); err != nil {
		// Try single event
		var single futuresTradeEvent
		if err := json.Unmarshal(msg, &single); err != nil {
			log.Printf("failed to parse futures trade: %s", string(msg))
			return false
		}
		events = []futuresTradeEvent{single}
	}

	wrote := false
	for _, ev := range events {
		if ev.Ev != "T" {
			continue
		}
		t := time.Unix(0, ev.Timestamp*int64(time.Millisecond))
		fmt.Fprintf(w, "%s\t%s\t%.4f\t%.0f\n",
			t.Format("15:04:05.000"),
			ev.Symbol,
			ev.Price,
			ev.Size)
		wrote = true
	}
	return wrote
}

// renderFuturesQuote parses a raw WebSocket message as a futures quote event
// array and writes each quote as a table row. Returns true if any rows were
// written so the caller knows to flush the tabwriter.
func renderFuturesQuote(msg []byte, w *tabwriter.Writer) bool {
	var events []futuresQuoteEvent
	if err := json.Unmarshal(msg, &events); err != nil {
		var single futuresQuoteEvent
		if err := json.Unmarshal(msg, &single); err != nil {
			log.Printf("failed to parse futures quote: %s", string(msg))
			return false
		}
		events = []futuresQuoteEvent{single}
	}

	wrote := false
	for _, ev := range events {
		if ev.Ev != "Q" {
			continue
		}
		t := time.Unix(0, ev.Timestamp*int64(time.Millisecond))
		fmt.Fprintf(w, "%s\t%s\t%.4f\t%.0f\t%.4f\t%.0f\n",
			t.Format("15:04:05.000"),
			ev.Symbol,
			ev.BidPrice,
			ev.BidSize,
			ev.AskPrice,
			ev.AskSize)
		wrote = true
	}
	return wrote
}

// renderFuturesAggregate parses a raw WebSocket message as a futures aggregate
// event array and writes each bar as a table row. Works for both per-minute (AM)
// and per-second (A) channels. Returns true if any rows were written.
func renderFuturesAggregate(msg []byte, w *tabwriter.Writer) bool {
	var events []futuresAggregateEvent
	if err := json.Unmarshal(msg, &events); err != nil {
		var single futuresAggregateEvent
		if err := json.Unmarshal(msg, &single); err != nil {
			log.Printf("failed to parse futures aggregate: %s", string(msg))
			return false
		}
		events = []futuresAggregateEvent{single}
	}

	wrote := false
	for _, ev := range events {
		if ev.Ev != "AM" && ev.Ev != "A" {
			continue
		}
		t := time.Unix(0, ev.Start*int64(time.Millisecond))
		fmt.Fprintf(w, "%s\t%s\t%.4f\t%.4f\t%.4f\t%.4f\t%.2f\n",
			t.Format("15:04:05.000"),
			ev.Symbol,
			ev.Open,
			ev.High,
			ev.Low,
			ev.Close,
			ev.Volume)
		wrote = true
	}
	return wrote
}

// init registers all futures WebSocket subcommands under the wsFuturesCmd
// parent, which is itself registered under the wsCmd parent. Each subcommand
// receives an --all flag for subscribing to all available contracts.
func init() {
	wsFuturesTradesCmd.Flags().Bool("all", false, "Subscribe to all futures contracts")
	wsFuturesQuotesCmd.Flags().Bool("all", false, "Subscribe to all futures contracts")
	wsFuturesAggMinuteCmd.Flags().Bool("all", false, "Subscribe to all futures contracts")
	wsFuturesAggSecondCmd.Flags().Bool("all", false, "Subscribe to all futures contracts")

	wsFuturesCmd.AddCommand(wsFuturesTradesCmd)
	wsFuturesCmd.AddCommand(wsFuturesQuotesCmd)
	wsFuturesCmd.AddCommand(wsFuturesAggMinuteCmd)
	wsFuturesCmd.AddCommand(wsFuturesAggSecondCmd)

	wsCmd.AddCommand(wsFuturesCmd)
}
