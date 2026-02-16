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

// wsCryptoCmd is the parent command for all crypto WebSocket streaming
// subcommands. It groups real-time crypto data streams under the "ws crypto"
// namespace including trades, quotes, aggregates, and fair market value.
var wsCryptoCmd = &cobra.Command{
	Use:   "crypto",
	Short: "Stream real-time crypto data",
}

// wsCryptoTradesCmd streams real-time crypto trade events via WebSocket.
// Each trade event includes the crypto pair, price, size, exchange, and
// timestamp. Supports subscribing to specific tickers or all tickers.
// Usage: massive ws crypto trades X:BTC-USD X:ETH-USD
var wsCryptoTradesCmd = &cobra.Command{
	Use:   "trades [tickers...]",
	Short: "Stream real-time crypto trades",
	Long:  "Stream real-time crypto trade data via WebSocket. Each event includes pair, price, size, exchange, and timestamp.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		tickers := args
		if all {
			tickers = []string{"*"}
		}

		subscriptions := buildCryptoSubscriptions("XT", tickers)
		return connectCryptoWebSocket("XT", subscriptions, printCryptoTrade)
	},
}

// wsCryptoQuotesCmd streams real-time crypto quote events via WebSocket.
// Each quote event includes the crypto pair, bid/ask prices, bid/ask sizes,
// exchange, and timestamp. Supports subscribing to specific tickers or all.
// Usage: massive ws crypto quotes X:BTC-USD
var wsCryptoQuotesCmd = &cobra.Command{
	Use:   "quotes [tickers...]",
	Short: "Stream real-time crypto quotes",
	Long:  "Stream real-time crypto quote data via WebSocket. Each event includes pair, bid/ask prices and sizes, exchange, and timestamp.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		tickers := args
		if all {
			tickers = []string{"*"}
		}

		subscriptions := buildCryptoSubscriptions("XQ", tickers)
		return connectCryptoWebSocket("XQ", subscriptions, printCryptoQuote)
	},
}

// wsCryptoAggMinuteCmd streams real-time per-minute aggregate bar data for
// crypto pairs via WebSocket. Each event includes open, high, low, close,
// volume, VWAP, and the start/end timestamps of the aggregate window.
// Usage: massive ws crypto agg-minute X:BTC-USD
var wsCryptoAggMinuteCmd = &cobra.Command{
	Use:   "agg-minute [tickers...]",
	Short: "Stream per-minute crypto aggregates",
	Long:  "Stream real-time per-minute aggregate bar data for crypto pairs via WebSocket. Each event includes OHLCV data and VWAP.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		tickers := args
		if all {
			tickers = []string{"*"}
		}

		subscriptions := buildCryptoSubscriptions("XA", tickers)
		return connectCryptoWebSocket("XA", subscriptions, printCryptoAggregate)
	},
}

// wsCryptoAggSecondCmd streams real-time per-second aggregate bar data for
// crypto pairs via WebSocket. Each event includes the same fields as the
// per-minute aggregate but at second-level granularity.
// Usage: massive ws crypto agg-second X:BTC-USD
var wsCryptoAggSecondCmd = &cobra.Command{
	Use:   "agg-second [tickers...]",
	Short: "Stream per-second crypto aggregates",
	Long:  "Stream real-time per-second aggregate bar data for crypto pairs via WebSocket. Each event includes OHLCV data and VWAP.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		tickers := args
		if all {
			tickers = []string{"*"}
		}

		subscriptions := buildCryptoSubscriptions("XAS", tickers)
		return connectCryptoWebSocket("XAS", subscriptions, printCryptoAggregate)
	},
}

// wsCryptoFMVCmd streams real-time Fair Market Value (FMV) data for crypto
// pairs via WebSocket. FMV represents a calculated fair price for a crypto
// asset across multiple exchanges. Each event includes symbol, FMV, and timestamp.
// Usage: massive ws crypto fmv X:BTC-USD
var wsCryptoFMVCmd = &cobra.Command{
	Use:   "fmv [tickers...]",
	Short: "Stream crypto Fair Market Value",
	Long:  "Stream real-time Fair Market Value (FMV) data for crypto pairs via WebSocket. FMV provides a calculated fair price across exchanges.",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return fmt.Errorf("provide at least one ticker or use --all")
		}

		tickers := args
		if all {
			tickers = []string{"*"}
		}

		subscriptions := buildCryptoSubscriptions("FMV", tickers)
		return connectCryptoWebSocket("FMV", subscriptions, printCryptoFMV)
	},
}

// cryptoTradeEvent represents a single crypto trade event received from
// the WebSocket stream. Fields map to the XT channel data structure.
type cryptoTradeEvent struct {
	Ev        string  `json:"ev"`
	Pair      string  `json:"pair"`
	Price     float64 `json:"p"`
	Timestamp int64   `json:"t"`
	Size      float64 `json:"s"`
	Conditions []int  `json:"c"`
	TradeID   string  `json:"i"`
	Exchange  int     `json:"x"`
	Received  int64   `json:"r"`
}

// cryptoQuoteEvent represents a single crypto quote event received from
// the WebSocket stream. Fields map to the XQ channel data structure.
type cryptoQuoteEvent struct {
	Ev       string  `json:"ev"`
	Pair     string  `json:"pair"`
	BidPrice float64 `json:"bp"`
	BidSize  float64 `json:"bs"`
	AskPrice float64 `json:"ap"`
	AskSize  float64 `json:"as"`
	Timestamp int64  `json:"t"`
	Exchange int     `json:"x"`
	Received int64   `json:"r"`
}

// cryptoAggEvent represents a crypto aggregate bar event received from
// the WebSocket stream. Used for both XA (per-minute) and XAS (per-second)
// channels. Fields include OHLCV data and VWAP.
type cryptoAggEvent struct {
	Ev     string  `json:"ev"`
	Pair   string  `json:"pair"`
	Open   float64 `json:"o"`
	Close  float64 `json:"c"`
	High   float64 `json:"h"`
	Low    float64 `json:"l"`
	Volume float64 `json:"v"`
	VWAP   float64 `json:"vw"`
	AvgSize float64 `json:"z"`
	Start  int64   `json:"s"`
	End    int64   `json:"e"`
}

// cryptoFMVEvent represents a Fair Market Value event received from the
// WebSocket stream. Provides a calculated fair price for a crypto asset.
type cryptoFMVEvent struct {
	Ev     string  `json:"ev"`
	FMV    float64 `json:"fmv"`
	Symbol string  `json:"sym"`
	Timestamp int64 `json:"t"`
}

// buildCryptoSubscriptions constructs the subscription parameter string for
// the WebSocket subscribe message. It prefixes each ticker with the channel
// name (e.g., "XT.X:BTC-USD"). Multiple tickers are comma-separated.
func buildCryptoSubscriptions(channel string, tickers []string) string {
	parts := make([]string, len(tickers))
	for i, t := range tickers {
		parts[i] = channel + "." + t
	}
	return strings.Join(parts, ",")
}

// connectCryptoWebSocket establishes a WebSocket connection to the Massive
// crypto streaming endpoint, sends a subscribe message for the given channel
// and subscriptions, and reads messages in a loop until the user presses
// Ctrl+C. Each incoming message is passed to the provided handler function
// for display. The connection uses the API key from config for authentication.
func connectCryptoWebSocket(channel string, subscriptions string, handler func([]byte, *tabwriter.Writer) error) error {
	apiKey, err := config.GetAPIKey()
	if err != nil {
		return err
	}

	// Build the WebSocket URL. FMV uses a different path than other crypto channels.
	var wsURL string
	if channel == "FMV" {
		wsURL = fmt.Sprintf("wss://socket.massive.com/business/crypto/FMV?apiKey=%s", apiKey)
	} else {
		wsURL = fmt.Sprintf("wss://socket.massive.com/crypto/%s?apiKey=%s", channel, apiKey)
	}

	fmt.Printf("Connecting to %s stream...\n", channel)

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}
	defer conn.Close()

	// Send the subscribe message to begin receiving events for the requested tickers.
	subscribeMsg := map[string]string{
		"action": "subscribe",
		"params": subscriptions,
	}
	if err := conn.WriteJSON(subscribeMsg); err != nil {
		return fmt.Errorf("failed to send subscribe message: %w", err)
	}

	fmt.Printf("Subscribed to: %s\n\n", subscriptions)

	// Set up a signal handler so Ctrl+C gracefully closes the connection.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Create a tabwriter for table output formatting.
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Print the table header if we are in table output mode.
	if outputFormat == "table" {
		printCryptoTableHeader(channel, w)
	}

	// Read messages in a loop until interrupted.
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					return
				}
				fmt.Fprintf(os.Stderr, "\nWebSocket read error: %v\n", err)
				return
			}

			if outputFormat == "json" {
				fmt.Println(string(message))
				continue
			}

			if err := handler(message, w); err != nil {
				fmt.Fprintf(os.Stderr, "\nError processing message: %v\n", err)
			}
		}
	}()

	// Wait for either the read loop to finish or a Ctrl+C interrupt.
	select {
	case <-done:
		return nil
	case <-interrupt:
		fmt.Println("\nDisconnecting...")

		// Send a close message to the server for a clean shutdown.
		err := conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		)
		if err != nil {
			return fmt.Errorf("error sending close message: %w", err)
		}

		// Wait briefly for the read loop to finish after close.
		select {
		case <-done:
		case <-time.After(time.Second):
		}
		return nil
	}
}

// printCryptoTableHeader writes the appropriate column header to the
// tabwriter based on the WebSocket channel type (XT, XQ, XA, XAS, FMV).
func printCryptoTableHeader(channel string, w *tabwriter.Writer) {
	switch channel {
	case "XT":
		fmt.Fprintln(w, "TIME\tPAIR\tPRICE\tSIZE\tEXCHANGE")
		fmt.Fprintln(w, "----\t----\t-----\t----\t--------")
	case "XQ":
		fmt.Fprintln(w, "TIME\tPAIR\tBID\tBID_SIZE\tASK\tASK_SIZE")
		fmt.Fprintln(w, "----\t----\t---\t--------\t---\t--------")
	case "XA", "XAS":
		fmt.Fprintln(w, "TIME\tPAIR\tOPEN\tHIGH\tLOW\tCLOSE\tVOLUME")
		fmt.Fprintln(w, "----\t----\t----\t----\t---\t-----\t------")
	case "FMV":
		fmt.Fprintln(w, "TIME\tSYMBOL\tFMV")
		fmt.Fprintln(w, "----\t------\t---")
	}
	w.Flush()
}

// printCryptoTrade parses a WebSocket message as an array of crypto trade
// events and prints each one in table format showing time, pair, price,
// size, and exchange.
func printCryptoTrade(message []byte, w *tabwriter.Writer) error {
	var events []cryptoTradeEvent
	if err := json.Unmarshal(message, &events); err != nil {
		return err
	}

	for _, ev := range events {
		if ev.Ev != "XT" {
			continue
		}
		t := time.Unix(0, ev.Timestamp*int64(time.Millisecond))
		fmt.Fprintf(w, "%s\t%s\t%.8f\t%.8f\t%d\n",
			t.Format("15:04:05.000"),
			ev.Pair,
			ev.Price,
			ev.Size,
			ev.Exchange,
		)
		w.Flush()
	}
	return nil
}

// printCryptoQuote parses a WebSocket message as an array of crypto quote
// events and prints each one in table format showing time, pair, bid price,
// bid size, ask price, and ask size.
func printCryptoQuote(message []byte, w *tabwriter.Writer) error {
	var events []cryptoQuoteEvent
	if err := json.Unmarshal(message, &events); err != nil {
		return err
	}

	for _, ev := range events {
		if ev.Ev != "XQ" {
			continue
		}
		t := time.Unix(0, ev.Timestamp*int64(time.Millisecond))
		fmt.Fprintf(w, "%s\t%s\t%.8f\t%.8f\t%.8f\t%.8f\n",
			t.Format("15:04:05.000"),
			ev.Pair,
			ev.BidPrice,
			ev.BidSize,
			ev.AskPrice,
			ev.AskSize,
		)
		w.Flush()
	}
	return nil
}

// printCryptoAggregate parses a WebSocket message as an array of crypto
// aggregate events (per-minute or per-second) and prints each one in table
// format showing time, pair, open, high, low, close, and volume.
func printCryptoAggregate(message []byte, w *tabwriter.Writer) error {
	var events []cryptoAggEvent
	if err := json.Unmarshal(message, &events); err != nil {
		return err
	}

	for _, ev := range events {
		if ev.Ev != "XA" && ev.Ev != "XAS" {
			continue
		}
		t := time.Unix(0, ev.Start*int64(time.Millisecond))
		fmt.Fprintf(w, "%s\t%s\t%.8f\t%.8f\t%.8f\t%.8f\t%.8f\n",
			t.Format("15:04:05.000"),
			ev.Pair,
			ev.Open,
			ev.High,
			ev.Low,
			ev.Close,
			ev.Volume,
		)
		w.Flush()
	}
	return nil
}

// printCryptoFMV parses a WebSocket message as an array of crypto Fair
// Market Value events and prints each one in table format showing time,
// symbol, and the FMV price.
func printCryptoFMV(message []byte, w *tabwriter.Writer) error {
	var events []cryptoFMVEvent
	if err := json.Unmarshal(message, &events); err != nil {
		return err
	}

	for _, ev := range events {
		if ev.Ev != "FMV" {
			continue
		}
		t := time.Unix(0, ev.Timestamp*int64(time.Millisecond))
		fmt.Fprintf(w, "%s\t%s\t%.8f\n",
			t.Format("15:04:05.000"),
			ev.Symbol,
			ev.FMV,
		)
		w.Flush()
	}
	return nil
}

// init registers the crypto WebSocket command and all its subcommands under
// the ws parent command. Each subcommand gets an --all flag to subscribe
// to all available tickers for that channel.
func init() {
	// Add the --all flag to each subcommand for subscribing to all tickers.
	wsCryptoTradesCmd.Flags().Bool("all", false, "Subscribe to all crypto trade events")
	wsCryptoQuotesCmd.Flags().Bool("all", false, "Subscribe to all crypto quote events")
	wsCryptoAggMinuteCmd.Flags().Bool("all", false, "Subscribe to all crypto per-minute aggregates")
	wsCryptoAggSecondCmd.Flags().Bool("all", false, "Subscribe to all crypto per-second aggregates")
	wsCryptoFMVCmd.Flags().Bool("all", false, "Subscribe to all crypto FMV events")

	// Register subcommands under the crypto parent.
	wsCryptoCmd.AddCommand(wsCryptoTradesCmd)
	wsCryptoCmd.AddCommand(wsCryptoQuotesCmd)
	wsCryptoCmd.AddCommand(wsCryptoAggMinuteCmd)
	wsCryptoCmd.AddCommand(wsCryptoAggSecondCmd)
	wsCryptoCmd.AddCommand(wsCryptoFMVCmd)

	// Register the crypto command under the ws parent command.
	wsCmd.AddCommand(wsCryptoCmd)
}
