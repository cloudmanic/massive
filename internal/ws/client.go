//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

// defaultBaseURL is the production WebSocket endpoint for real-time streaming
// market data from the Massive API.
const defaultBaseURL = "wss://socket.massive.com"

// StreamConfig holds configuration for a WebSocket stream connection. It specifies
// the API key for authentication, the base URL to connect to, the asset class and
// channel to subscribe to, and the initial set of ticker symbols to stream.
type StreamConfig struct {
	APIKey  string   // API key used for authentication via query parameter
	BaseURL string   // WebSocket base URL; defaults to wss://socket.massive.com
	Asset   string   // Asset class: stocks, options, indices, crypto, forex, futures
	Channel string   // Data channel: T (trades), Q (quotes), AM, A (aggregates), etc.
	Tickers []string // Ticker symbols to subscribe to, or ["*"] for all tickers
}

// Client manages a WebSocket connection for streaming market data from the Massive
// API. It handles connecting, subscribing/unsubscribing to ticker symbols, reading
// incoming messages, and gracefully closing the connection. All write operations
// are protected by a mutex to ensure thread safety.
type Client struct {
	config StreamConfig
	conn   *websocket.Conn
	done   chan struct{}
	mu     sync.Mutex
}

// Message represents a raw WebSocket message event received from the Massive API.
// Each event has an event type identifier and the full raw JSON payload for
// downstream consumers to parse according to their needs.
type Message struct {
	EventType string          `json:"ev"`      // Event type identifier (e.g., "T" for trade, "Q" for quote)
	RawJSON   json.RawMessage `json:"-"`       // The full raw JSON for the event, populated during Listen
}

// subscribeAction is the internal struct used to marshal subscribe and unsubscribe
// JSON messages sent over the WebSocket connection to the Massive API.
type subscribeAction struct {
	Action string `json:"action"` // "subscribe" or "unsubscribe"
	Params string `json:"params"` // Comma-separated ticker symbols (e.g., "T.MSFT,T.AAPL")
}

// NewClient creates a new WebSocket client with the provided StreamConfig. If the
// config's BaseURL is empty, it defaults to the production WebSocket endpoint at
// wss://socket.massive.com. The client is not connected until Connect is called.
func NewClient(config StreamConfig) *Client {
	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURL
	}

	return &Client{
		config: config,
		done:   make(chan struct{}),
	}
}

// buildURL constructs the full WebSocket connection URL from the client's config.
// The URL follows the pattern: {baseURL}/{asset}/{channel}?apiKey={key}. This is
// used internally by Connect to determine where to dial the WebSocket connection.
func (c *Client) buildURL() string {
	return fmt.Sprintf("%s/%s/%s?apiKey=%s", c.config.BaseURL, c.config.Asset, c.config.Channel, c.config.APIKey)
}

// Connect establishes a WebSocket connection to the Massive streaming API using
// the URL built from the client's configuration. The provided context is not
// currently used for dialing but is accepted for future cancellation support.
// Returns an error if the WebSocket handshake fails.
func (c *Client) Connect(ctx context.Context) error {
	url := c.buildURL()

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", url, err)
	}

	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()

	return nil
}

// Subscribe sends a subscribe action to the Massive WebSocket API for the given
// ticker symbols. The tickers are joined into a comma-separated string and sent
// as the "params" field. For example, Subscribe("T.MSFT", "T.AAPL") sends
// {"action":"subscribe","params":"T.MSFT,T.AAPL"}. Returns an error if the
// connection is not established or the message fails to send.
func (c *Client) Subscribe(tickers ...string) error {
	return c.sendAction("subscribe", tickers)
}

// Unsubscribe sends an unsubscribe action to the Massive WebSocket API for the
// given ticker symbols. The tickers are joined into a comma-separated string and
// sent as the "params" field. Returns an error if the connection is not established
// or the message fails to send.
func (c *Client) Unsubscribe(tickers ...string) error {
	return c.sendAction("unsubscribe", tickers)
}

// sendAction marshals and sends a subscribe or unsubscribe action message over the
// WebSocket connection. It joins the provided tickers into a comma-separated string
// for the "params" field. The write operation is protected by a mutex to prevent
// concurrent writes to the WebSocket connection. Returns an error if the connection
// is nil or the write fails.
func (c *Client) sendAction(action string, tickers []string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return fmt.Errorf("websocket connection is not established")
	}

	msg := subscribeAction{
		Action: action,
		Params: strings.Join(tickers, ","),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal %s message: %w", action, err)
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("failed to send %s message: %w", action, err)
	}

	return nil
}

// Listen reads messages from the WebSocket connection in a loop and passes each
// raw message to the provided handler function. The loop terminates when the
// connection is closed (either by the server or by calling Close), when a read
// error occurs, or when the client's done channel is closed. Each message is
// delivered as a raw byte slice so the handler can parse it as needed.
func (c *Client) Listen(handler func([]byte)) error {
	for {
		select {
		case <-c.done:
			return nil
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				// Check if the error is due to a normal close.
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					return nil
				}

				select {
				case <-c.done:
					return nil
				default:
					return fmt.Errorf("read error: %w", err)
				}
			}

			handler(message)
		}
	}
}

// Close gracefully closes the WebSocket connection by sending a close message to
// the server and then closing the underlying connection. It signals the done channel
// to stop the Listen loop. The close operation is protected by a mutex to prevent
// concurrent access. Returns an error if closing the connection fails, or nil if
// the connection was already nil.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Signal the done channel to stop the Listen loop.
	select {
	case <-c.done:
		// Already closed.
	default:
		close(c.done)
	}

	if c.conn == nil {
		return nil
	}

	// Send a close message to the server for a graceful shutdown.
	err := c.conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)
	if err != nil {
		// Even if writing the close message fails, still close the connection.
		c.conn.Close()
		return fmt.Errorf("failed to send close message: %w", err)
	}

	return c.conn.Close()
}
