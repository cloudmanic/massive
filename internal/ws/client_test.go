//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// upgrader is the WebSocket upgrader used by mock servers in tests. It accepts
// all origins to simplify test setup.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// TestNewClientDefaults verifies that NewClient sets the default base URL when
// none is provided in the StreamConfig, and that all other config fields are
// preserved as given.
func TestNewClientDefaults(t *testing.T) {
	client := NewClient(StreamConfig{
		APIKey:  "test-key",
		Asset:   "stocks",
		Channel: "T",
		Tickers: []string{"MSFT"},
	})

	if client.config.BaseURL != defaultBaseURL {
		t.Errorf("expected default base URL %s, got %s", defaultBaseURL, client.config.BaseURL)
	}

	if client.config.APIKey != "test-key" {
		t.Errorf("expected API key test-key, got %s", client.config.APIKey)
	}

	if client.config.Asset != "stocks" {
		t.Errorf("expected asset stocks, got %s", client.config.Asset)
	}

	if client.config.Channel != "T" {
		t.Errorf("expected channel T, got %s", client.config.Channel)
	}

	if client.done == nil {
		t.Error("expected done channel to be initialized")
	}
}

// TestNewClientCustomBaseURL verifies that NewClient preserves a custom base URL
// when one is explicitly provided in the StreamConfig instead of using the default.
func TestNewClientCustomBaseURL(t *testing.T) {
	client := NewClient(StreamConfig{
		APIKey:  "test-key",
		BaseURL: "wss://delayed.massive.com",
		Asset:   "stocks",
		Channel: "T",
	})

	if client.config.BaseURL != "wss://delayed.massive.com" {
		t.Errorf("expected custom base URL wss://delayed.massive.com, got %s", client.config.BaseURL)
	}
}

// TestBuildURL verifies that the client constructs the correct WebSocket URL
// from the config fields, following the pattern {baseURL}/{asset}/{channel}?apiKey={key}.
func TestBuildURL(t *testing.T) {
	client := NewClient(StreamConfig{
		APIKey:  "my-api-key",
		BaseURL: "wss://socket.massive.com",
		Asset:   "stocks",
		Channel: "T",
	})

	url := client.buildURL()
	expected := "wss://socket.massive.com/stocks/T?apiKey=my-api-key"

	if url != expected {
		t.Errorf("expected URL %s, got %s", expected, url)
	}
}

// TestBuildURLWithDifferentAssets verifies that the URL builder correctly handles
// different asset class and channel combinations in the URL path.
func TestBuildURLWithDifferentAssets(t *testing.T) {
	tests := []struct {
		name     string
		asset    string
		channel  string
		expected string
	}{
		{
			name:     "crypto trades",
			asset:    "crypto",
			channel:  "T",
			expected: "wss://socket.massive.com/crypto/T?apiKey=key",
		},
		{
			name:     "options quotes",
			asset:    "options",
			channel:  "Q",
			expected: "wss://socket.massive.com/options/Q?apiKey=key",
		},
		{
			name:     "stocks aggregates",
			asset:    "stocks",
			channel:  "AM",
			expected: "wss://socket.massive.com/stocks/AM?apiKey=key",
		},
		{
			name:     "forex trades",
			asset:    "forex",
			channel:  "T",
			expected: "wss://socket.massive.com/forex/T?apiKey=key",
		},
		{
			name:     "indices aggregates",
			asset:    "indices",
			channel:  "A",
			expected: "wss://socket.massive.com/indices/A?apiKey=key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(StreamConfig{
				APIKey:  "key",
				Asset:   tt.asset,
				Channel: tt.channel,
			})

			url := client.buildURL()
			if url != tt.expected {
				t.Errorf("expected URL %s, got %s", tt.expected, url)
			}
		})
	}
}

// TestConnectToMockServer verifies that Connect successfully establishes a
// WebSocket connection to a mock server and that the connection is non-nil
// after connecting.
func TestConnectToMockServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("failed to upgrade: %v", err)
			return
		}
		defer conn.Close()

		// Keep the connection open until the test completes.
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	// Convert http:// to ws:// for the WebSocket client.
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	client := NewClient(StreamConfig{
		APIKey:  "test-key",
		BaseURL: wsURL,
		Asset:   "stocks",
		Channel: "T",
	})

	err := client.Connect(context.Background())
	if err != nil {
		t.Fatalf("unexpected error connecting: %v", err)
	}
	defer client.Close()

	if client.conn == nil {
		t.Error("expected connection to be established, got nil")
	}
}

// TestConnectURLIncludesAPIKey verifies that the WebSocket connection URL includes
// the apiKey query parameter so the server can authenticate the client.
func TestConnectURLIncludesAPIKey(t *testing.T) {
	var receivedAPIKey string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAPIKey = r.URL.Query().Get("apiKey")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("failed to upgrade: %v", err)
			return
		}
		defer conn.Close()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	client := NewClient(StreamConfig{
		APIKey:  "secret-api-key-123",
		BaseURL: wsURL,
		Asset:   "stocks",
		Channel: "T",
	})

	err := client.Connect(context.Background())
	if err != nil {
		t.Fatalf("unexpected error connecting: %v", err)
	}
	defer client.Close()

	if receivedAPIKey != "secret-api-key-123" {
		t.Errorf("expected apiKey=secret-api-key-123, got %s", receivedAPIKey)
	}
}

// TestConnectURLIncludesAssetAndChannel verifies that the WebSocket connection URL
// includes the correct asset class and channel in the path.
func TestConnectURLIncludesAssetAndChannel(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("failed to upgrade: %v", err)
			return
		}
		defer conn.Close()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	client := NewClient(StreamConfig{
		APIKey:  "key",
		BaseURL: wsURL,
		Asset:   "options",
		Channel: "Q",
	})

	err := client.Connect(context.Background())
	if err != nil {
		t.Fatalf("unexpected error connecting: %v", err)
	}
	defer client.Close()

	if receivedPath != "/options/Q" {
		t.Errorf("expected path /options/Q, got %s", receivedPath)
	}
}

// TestConnectFailsWithBadURL verifies that Connect returns an error when given
// an invalid WebSocket URL that cannot be dialed.
func TestConnectFailsWithBadURL(t *testing.T) {
	client := NewClient(StreamConfig{
		APIKey:  "key",
		BaseURL: "ws://localhost:1",
		Asset:   "stocks",
		Channel: "T",
	})

	err := client.Connect(context.Background())
	if err == nil {
		t.Fatal("expected connection error, got nil")
	}
}

// TestSubscribeSendsCorrectJSON verifies that Subscribe sends a properly formatted
// JSON message with the "subscribe" action and comma-separated ticker params.
func TestSubscribeSendsCorrectJSON(t *testing.T) {
	receivedCh := make(chan []byte, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("failed to upgrade: %v", err)
			return
		}
		defer conn.Close()

		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		receivedCh <- msg

		// Keep connection open.
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	client := NewClient(StreamConfig{
		APIKey:  "key",
		BaseURL: wsURL,
		Asset:   "stocks",
		Channel: "T",
	})

	if err := client.Connect(context.Background()); err != nil {
		t.Fatalf("unexpected error connecting: %v", err)
	}
	defer client.Close()

	if err := client.Subscribe("T.MSFT", "T.AAPL"); err != nil {
		t.Fatalf("unexpected error subscribing: %v", err)
	}

	select {
	case msg := <-receivedCh:
		var action subscribeAction
		if err := json.Unmarshal(msg, &action); err != nil {
			t.Fatalf("failed to parse subscribe message: %v", err)
		}

		if action.Action != "subscribe" {
			t.Errorf("expected action subscribe, got %s", action.Action)
		}

		if action.Params != "T.MSFT,T.AAPL" {
			t.Errorf("expected params T.MSFT,T.AAPL, got %s", action.Params)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for subscribe message")
	}
}

// TestSubscribeSingleTicker verifies that Subscribe correctly sends a subscribe
// message for a single ticker without trailing commas.
func TestSubscribeSingleTicker(t *testing.T) {
	receivedCh := make(chan []byte, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("failed to upgrade: %v", err)
			return
		}
		defer conn.Close()

		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		receivedCh <- msg

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	client := NewClient(StreamConfig{
		APIKey:  "key",
		BaseURL: wsURL,
		Asset:   "stocks",
		Channel: "T",
	})

	if err := client.Connect(context.Background()); err != nil {
		t.Fatalf("unexpected error connecting: %v", err)
	}
	defer client.Close()

	if err := client.Subscribe("T.MSFT"); err != nil {
		t.Fatalf("unexpected error subscribing: %v", err)
	}

	select {
	case msg := <-receivedCh:
		var action subscribeAction
		if err := json.Unmarshal(msg, &action); err != nil {
			t.Fatalf("failed to parse subscribe message: %v", err)
		}

		if action.Params != "T.MSFT" {
			t.Errorf("expected params T.MSFT, got %s", action.Params)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for subscribe message")
	}
}

// TestSubscribeWithoutConnection verifies that Subscribe returns an error when
// called before Connect has established a WebSocket connection.
func TestSubscribeWithoutConnection(t *testing.T) {
	client := NewClient(StreamConfig{
		APIKey:  "key",
		Asset:   "stocks",
		Channel: "T",
	})

	err := client.Subscribe("T.MSFT")
	if err == nil {
		t.Fatal("expected error when subscribing without connection, got nil")
	}

	if !strings.Contains(err.Error(), "not established") {
		t.Errorf("expected error about connection not established, got: %s", err.Error())
	}
}

// TestUnsubscribeSendsCorrectJSON verifies that Unsubscribe sends a properly
// formatted JSON message with the "unsubscribe" action and comma-separated
// ticker params.
func TestUnsubscribeSendsCorrectJSON(t *testing.T) {
	receivedCh := make(chan []byte, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("failed to upgrade: %v", err)
			return
		}
		defer conn.Close()

		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		receivedCh <- msg

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	client := NewClient(StreamConfig{
		APIKey:  "key",
		BaseURL: wsURL,
		Asset:   "stocks",
		Channel: "T",
	})

	if err := client.Connect(context.Background()); err != nil {
		t.Fatalf("unexpected error connecting: %v", err)
	}
	defer client.Close()

	if err := client.Unsubscribe("T.MSFT", "T.GOOG"); err != nil {
		t.Fatalf("unexpected error unsubscribing: %v", err)
	}

	select {
	case msg := <-receivedCh:
		var action subscribeAction
		if err := json.Unmarshal(msg, &action); err != nil {
			t.Fatalf("failed to parse unsubscribe message: %v", err)
		}

		if action.Action != "unsubscribe" {
			t.Errorf("expected action unsubscribe, got %s", action.Action)
		}

		if action.Params != "T.MSFT,T.GOOG" {
			t.Errorf("expected params T.MSFT,T.GOOG, got %s", action.Params)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for unsubscribe message")
	}
}

// TestUnsubscribeWithoutConnection verifies that Unsubscribe returns an error
// when called before Connect has established a WebSocket connection.
func TestUnsubscribeWithoutConnection(t *testing.T) {
	client := NewClient(StreamConfig{
		APIKey:  "key",
		Asset:   "stocks",
		Channel: "T",
	})

	err := client.Unsubscribe("T.MSFT")
	if err == nil {
		t.Fatal("expected error when unsubscribing without connection, got nil")
	}
}

// TestListenReceivesMessages verifies that Listen correctly reads messages from
// the WebSocket connection and delivers each one to the provided handler function.
func TestListenReceivesMessages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("failed to upgrade: %v", err)
			return
		}
		defer conn.Close()

		// Send test messages to the client.
		messages := []string{
			`[{"ev":"T","sym":"MSFT","p":420.50}]`,
			`[{"ev":"T","sym":"AAPL","p":185.25}]`,
			`[{"ev":"T","sym":"GOOG","p":140.00}]`,
		}

		for _, msg := range messages {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				return
			}
		}

		// Close the connection after sending all messages.
		conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	client := NewClient(StreamConfig{
		APIKey:  "key",
		BaseURL: wsURL,
		Asset:   "stocks",
		Channel: "T",
	})

	if err := client.Connect(context.Background()); err != nil {
		t.Fatalf("unexpected error connecting: %v", err)
	}
	defer client.Close()

	var mu sync.Mutex
	var received []string

	err := client.Listen(func(msg []byte) {
		mu.Lock()
		received = append(received, string(msg))
		mu.Unlock()
	})

	if err != nil {
		t.Fatalf("unexpected error from Listen: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if len(received) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(received))
	}

	if !strings.Contains(received[0], "MSFT") {
		t.Errorf("expected first message to contain MSFT, got %s", received[0])
	}

	if !strings.Contains(received[1], "AAPL") {
		t.Errorf("expected second message to contain AAPL, got %s", received[1])
	}

	if !strings.Contains(received[2], "GOOG") {
		t.Errorf("expected third message to contain GOOG, got %s", received[2])
	}
}

// TestListenStopsOnClose verifies that the Listen loop terminates cleanly when
// Close is called, allowing the client to shut down without hanging.
func TestListenStopsOnClose(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("failed to upgrade: %v", err)
			return
		}
		defer conn.Close()

		// Keep connection open and read messages until the client disconnects.
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	client := NewClient(StreamConfig{
		APIKey:  "key",
		BaseURL: wsURL,
		Asset:   "stocks",
		Channel: "T",
	})

	if err := client.Connect(context.Background()); err != nil {
		t.Fatalf("unexpected error connecting: %v", err)
	}

	listenDone := make(chan error, 1)
	go func() {
		listenDone <- client.Listen(func(msg []byte) {})
	}()

	// Give Listen a moment to start reading.
	time.Sleep(100 * time.Millisecond)

	if err := client.Close(); err != nil {
		t.Fatalf("unexpected error closing: %v", err)
	}

	select {
	case err := <-listenDone:
		if err != nil {
			t.Fatalf("expected Listen to return nil after close, got: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for Listen to return after Close")
	}
}

// TestCloseWithoutConnect verifies that Close returns nil without error when
// called on a client that was never connected.
func TestCloseWithoutConnect(t *testing.T) {
	client := NewClient(StreamConfig{
		APIKey:  "key",
		Asset:   "stocks",
		Channel: "T",
	})

	err := client.Close()
	if err != nil {
		t.Fatalf("expected nil error when closing without connection, got: %v", err)
	}
}

// TestCloseCalledTwice verifies that calling Close multiple times does not panic
// or return an unexpected error.
func TestCloseCalledTwice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	client := NewClient(StreamConfig{
		APIKey:  "key",
		BaseURL: wsURL,
		Asset:   "stocks",
		Channel: "T",
	})

	if err := client.Connect(context.Background()); err != nil {
		t.Fatalf("unexpected error connecting: %v", err)
	}

	// First close should succeed.
	client.Close()

	// Second close should not panic.
	client.Close()
}

// TestSubscribeWildcard verifies that subscribing with the wildcard "*" symbol
// sends the correct params format for receiving all ticker events.
func TestSubscribeWildcard(t *testing.T) {
	receivedCh := make(chan []byte, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("failed to upgrade: %v", err)
			return
		}
		defer conn.Close()

		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		receivedCh <- msg

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	client := NewClient(StreamConfig{
		APIKey:  "key",
		BaseURL: wsURL,
		Asset:   "stocks",
		Channel: "T",
	})

	if err := client.Connect(context.Background()); err != nil {
		t.Fatalf("unexpected error connecting: %v", err)
	}
	defer client.Close()

	if err := client.Subscribe("*"); err != nil {
		t.Fatalf("unexpected error subscribing: %v", err)
	}

	select {
	case msg := <-receivedCh:
		var action subscribeAction
		if err := json.Unmarshal(msg, &action); err != nil {
			t.Fatalf("failed to parse subscribe message: %v", err)
		}

		if action.Params != "*" {
			t.Errorf("expected params *, got %s", action.Params)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for subscribe message")
	}
}
