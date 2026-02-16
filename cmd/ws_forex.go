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

// wsForexCmd is the parent command for all forex WebSocket streaming subcommands.
// It groups quotes, aggregates, and fair market value streams under "massive ws forex".
var wsForexCmd = &cobra.Command{
	Use:   "forex",
	Short: "Stream real-time forex data via WebSocket",
}

// wsForexQuotesCmd streams real-time forex quote data for the specified currency
// pairs. Each event includes the pair name, bid price, and ask price. Supports
// both table and JSON output formats and an --all flag to subscribe to all pairs.
// Usage: massive ws forex quotes C:EURUSD C:USD/CNH
var wsForexQuotesCmd = &cobra.Command{
	Use:   "quotes [tickers...]",
	Short: "Stream real-time forex quotes",
	Long:  "Stream real-time forex quote data including bid and ask prices for specified currency pairs.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		tickers := args

		if !all && len(tickers) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		subscriptions := buildForexSubscriptions("C", tickers, all)
		return runForexWebSocket("C", subscriptions, renderForexQuote)
	},
}

// wsForexAggMinuteCmd streams per-minute aggregate bars for the specified forex
// currency pairs. Each event includes open, high, low, close prices and volume.
// Usage: massive ws forex agg-minute C:EURUSD
var wsForexAggMinuteCmd = &cobra.Command{
	Use:   "agg-minute [tickers...]",
	Short: "Stream per-minute forex aggregates",
	Long:  "Stream per-minute aggregate bar data (OHLCV) for specified forex currency pairs.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		tickers := args

		if !all && len(tickers) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		subscriptions := buildForexSubscriptions("CA", tickers, all)
		return runForexWebSocket("CA", subscriptions, renderForexAggregate)
	},
}

// wsForexAggSecondCmd streams per-second aggregate bars for the specified forex
// currency pairs. Each event includes open, high, low, close prices and volume.
// Usage: massive ws forex agg-second C:EURUSD
var wsForexAggSecondCmd = &cobra.Command{
	Use:   "agg-second [tickers...]",
	Short: "Stream per-second forex aggregates",
	Long:  "Stream per-second aggregate bar data (OHLCV) for specified forex currency pairs.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		tickers := args

		if !all && len(tickers) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		subscriptions := buildForexSubscriptions("CAS", tickers, all)
		return runForexWebSocket("CAS", subscriptions, renderForexAggregate)
	},
}

// wsForexFMVCmd streams Fair Market Value data for the specified forex currency
// pairs. Each event includes the symbol and its computed fair market value.
// Usage: massive ws forex fmv C:EURUSD
var wsForexFMVCmd = &cobra.Command{
	Use:   "fmv [tickers...]",
	Short: "Stream forex Fair Market Value",
	Long:  "Stream Fair Market Value (FMV) data for specified forex currency pairs.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		tickers := args

		if !all && len(tickers) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		subscriptions := buildForexSubscriptions("FMV", tickers, all)
		return runForexWebSocket("FMV", subscriptions, renderForexFMV)
	},
}

// forexQuoteEvent represents a single forex quote message received from the
// WebSocket. Fields map to the API's event schema for the "C" channel.
type forexQuoteEvent struct {
	Ev   string  `json:"ev"`
	Pair string  `json:"p"`
	Ex   int     `json:"x"`
	Ask  float64 `json:"a"`
	Bid  float64 `json:"b"`
	T    int64   `json:"t"`
}

// forexAggregateEvent represents a forex aggregate bar message received from
// the WebSocket. Used for both per-minute (CA) and per-second (CAS) channels.
type forexAggregateEvent struct {
	Ev     string  `json:"ev"`
	Pair   string  `json:"pair"`
	Open   float64 `json:"o"`
	Close  float64 `json:"c"`
	High   float64 `json:"h"`
	Low    float64 `json:"l"`
	Volume float64 `json:"v"`
	Start  int64   `json:"s"`
	End    int64   `json:"e"`
}

// forexFMVEvent represents a Fair Market Value message received from the
// WebSocket for the "FMV" channel.
type forexFMVEvent struct {
	Ev     string  `json:"ev"`
	FMV    float64 `json:"fmv"`
	Symbol string  `json:"sym"`
	T      int64   `json:"t"`
}

// buildForexSubscriptions constructs the subscription parameter string for the
// given channel and tickers. When all is true, it subscribes to the wildcard
// pattern for that channel. Otherwise, it prefixes each ticker with the channel
// name separated by a dot (e.g., "C.C:EURUSD").
func buildForexSubscriptions(channel string, tickers []string, all bool) string {
	if all {
		return channel + ".*"
	}

	parts := make([]string, len(tickers))
	for i, t := range tickers {
		parts[i] = channel + "." + strings.ToUpper(t)
	}
	return strings.Join(parts, ",")
}

// connectForexWebSocket establishes a WebSocket connection to the Massive forex
// streaming endpoint for the given channel. It loads the API key from the
// environment or config file and dials the WebSocket URL. Returns the open
// connection or an error if the connection fails.
func connectForexWebSocket(channel string) (*websocket.Conn, error) {
	apiKey, err := config.GetAPIKey()
	if err != nil {
		return nil, err
	}

	// FMV uses a different URL path than the other forex channels.
	wsPath := fmt.Sprintf("forex/%s", channel)
	if channel == "FMV" {
		wsPath = "business/forex/FMV"
	}

	url := fmt.Sprintf("wss://socket.massive.com/%s?apiKey=%s", wsPath, apiKey)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to forex WebSocket: %w", err)
	}

	return conn, nil
}

// sendForexSubscription sends a subscribe action over the WebSocket connection
// with the provided params string. This tells the server which tickers and
// channels to stream data for.
func sendForexSubscription(conn *websocket.Conn, params string) error {
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

// runForexWebSocket is the main streaming loop for forex WebSocket commands.
// It connects to the WebSocket, sends the subscription, and continuously reads
// messages until the user presses Ctrl+C. Each raw message is passed to the
// provided render function for display. The render function receives the raw
// JSON bytes and a tabwriter for table output, and returns true if a row was
// written that should be flushed.
func runForexWebSocket(channel string, subscriptions string, render func([]byte, *tabwriter.Writer) bool) error {
	conn, err := connectForexWebSocket(channel)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := sendForexSubscription(conn, subscriptions); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Connected to forex %s stream. Press Ctrl+C to stop.\n", channel)

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

	go forexReadLoop(conn, msgChan, errChan)

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

// forexReadLoop continuously reads messages from the WebSocket connection and
// sends them to the provided message channel. If a read error occurs, it is
// sent to the error channel and the loop exits.
func forexReadLoop(conn *websocket.Conn, msgChan chan<- []byte, errChan chan<- error) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			errChan <- err
			return
		}
		msgChan <- msg
	}
}

// renderForexQuote parses a raw WebSocket message as a forex quote event array
// and writes each quote as a table row. Returns true if any rows were written
// so the caller knows to flush the tabwriter.
func renderForexQuote(msg []byte, w *tabwriter.Writer) bool {
	var events []forexQuoteEvent
	if err := json.Unmarshal(msg, &events); err != nil {
		// Try single event
		var single forexQuoteEvent
		if err := json.Unmarshal(msg, &single); err != nil {
			log.Printf("failed to parse forex quote: %s", string(msg))
			return false
		}
		events = []forexQuoteEvent{single}
	}

	wrote := false
	for _, ev := range events {
		if ev.Ev != "C" {
			continue
		}
		t := time.Unix(0, ev.T*int64(time.Millisecond))
		fmt.Fprintf(w, "%s\t%s\t%.6f\t%.6f\n",
			t.Format("15:04:05.000"),
			ev.Pair,
			ev.Bid,
			ev.Ask)
		wrote = true
	}
	return wrote
}

// renderForexAggregate parses a raw WebSocket message as a forex aggregate event
// array and writes each bar as a table row. Works for both per-minute (CA) and
// per-second (CAS) channels. Returns true if any rows were written.
func renderForexAggregate(msg []byte, w *tabwriter.Writer) bool {
	var events []forexAggregateEvent
	if err := json.Unmarshal(msg, &events); err != nil {
		var single forexAggregateEvent
		if err := json.Unmarshal(msg, &single); err != nil {
			log.Printf("failed to parse forex aggregate: %s", string(msg))
			return false
		}
		events = []forexAggregateEvent{single}
	}

	wrote := false
	for _, ev := range events {
		if ev.Ev != "CA" && ev.Ev != "CAS" {
			continue
		}
		t := time.Unix(0, ev.Start*int64(time.Millisecond))
		fmt.Fprintf(w, "%s\t%s\t%.6f\t%.6f\t%.6f\t%.6f\t%.2f\n",
			t.Format("15:04:05.000"),
			ev.Pair,
			ev.Open,
			ev.High,
			ev.Low,
			ev.Close,
			ev.Volume)
		wrote = true
	}
	return wrote
}

// renderForexFMV parses a raw WebSocket message as a forex FMV event array and
// writes each value as a table row. Returns true if any rows were written.
func renderForexFMV(msg []byte, w *tabwriter.Writer) bool {
	var events []forexFMVEvent
	if err := json.Unmarshal(msg, &events); err != nil {
		var single forexFMVEvent
		if err := json.Unmarshal(msg, &single); err != nil {
			log.Printf("failed to parse forex FMV: %s", string(msg))
			return false
		}
		events = []forexFMVEvent{single}
	}

	wrote := false
	for _, ev := range events {
		if ev.Ev != "FMV" {
			continue
		}
		t := time.Unix(0, ev.T*int64(time.Millisecond))
		fmt.Fprintf(w, "%s\t%s\t%.6f\n",
			t.Format("15:04:05.000"),
			ev.Symbol,
			ev.FMV)
		wrote = true
	}
	return wrote
}

// init registers all forex WebSocket subcommands under the wsForexCmd parent,
// which is itself registered under the wsCmd parent. Each subcommand receives
// an --all flag for subscribing to all available tickers.
func init() {
	wsForexQuotesCmd.Flags().Bool("all", false, "Subscribe to all forex pairs")
	wsForexAggMinuteCmd.Flags().Bool("all", false, "Subscribe to all forex pairs")
	wsForexAggSecondCmd.Flags().Bool("all", false, "Subscribe to all forex pairs")
	wsForexFMVCmd.Flags().Bool("all", false, "Subscribe to all forex pairs")

	wsForexCmd.AddCommand(wsForexQuotesCmd)
	wsForexCmd.AddCommand(wsForexAggMinuteCmd)
	wsForexCmd.AddCommand(wsForexAggSecondCmd)
	wsForexCmd.AddCommand(wsForexFMVCmd)

	wsCmd.AddCommand(wsForexCmd)
}
