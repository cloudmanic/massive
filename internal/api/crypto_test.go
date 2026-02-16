//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// -------------------------------------------------------------------
// Mock JSON Responses
// -------------------------------------------------------------------

const cryptoBarsJSON = `{
	"ticker": "X:BTCUSD",
	"queryCount": 2,
	"resultsCount": 2,
	"adjusted": true,
	"results": [
		{
			"v": 123456.78,
			"vw": 43250.5,
			"o": 43000.00,
			"c": 43500.00,
			"h": 43800.00,
			"l": 42900.00,
			"t": 1736139600000,
			"n": 15000
		},
		{
			"v": 98765.43,
			"vw": 43600.25,
			"o": 43500.00,
			"c": 43700.00,
			"h": 44000.00,
			"l": 43400.00,
			"t": 1736226000000,
			"n": 12500
		}
	],
	"status": "OK",
	"request_id": "crypto-bars-123"
}`

const cryptoDailyMarketSummaryJSON = `{
	"queryCount": 2,
	"resultsCount": 2,
	"adjusted": true,
	"results": [
		{
			"T": "X:BTCUSD",
			"v": 500000.0,
			"vw": 43100.50,
			"o": 43000.00,
			"c": 43200.00,
			"h": 43500.00,
			"l": 42800.00,
			"t": 1736139600000,
			"n": 25000
		},
		{
			"T": "X:ETHUSD",
			"v": 1200000.0,
			"vw": 2250.75,
			"o": 2200.00,
			"c": 2300.00,
			"h": 2350.00,
			"l": 2180.00,
			"t": 1736139600000,
			"n": 18000
		}
	],
	"status": "OK",
	"request_id": "crypto-market-123"
}`

const cryptoDailyTickerSummaryJSON = `{
	"symbol": "X:BTCUSD",
	"isUTC": true,
	"day": "2025-01-06",
	"open": 43000.00,
	"close": 43500.00,
	"openTrades": [
		{
			"c": [1, 2],
			"i": "trade-open-1",
			"p": 43000.00,
			"s": 0.5,
			"t": 1736139600000,
			"x": 1
		}
	],
	"closingTrades": [
		{
			"c": [3],
			"i": "trade-close-1",
			"p": 43500.00,
			"s": 1.2,
			"t": 1736225999000,
			"x": 2
		}
	]
}`

const cryptoPreviousDayBarJSON = `{
	"ticker": "X:BTCUSD",
	"queryCount": 1,
	"resultsCount": 1,
	"adjusted": true,
	"results": [
		{
			"v": 250000.0,
			"vw": 42500.00,
			"o": 42000.00,
			"c": 43000.00,
			"h": 43100.00,
			"l": 41900.00,
			"t": 1736053200000,
			"n": 20000
		}
	],
	"status": "OK",
	"request_id": "crypto-prev-123"
}`

const cryptoConditionsJSON = `{
	"results": [
		{
			"id": 1,
			"type": "sale_condition",
			"name": "Regular Sale",
			"asset_class": "crypto",
			"data_types": ["trade"],
			"legacy": false
		},
		{
			"id": 2,
			"type": "sale_condition",
			"name": "Block Trade",
			"asset_class": "crypto",
			"data_types": ["trade"],
			"legacy": false
		}
	],
	"status": "OK",
	"request_id": "conditions-123",
	"count": 2
}`

const cryptoExchangesJSON = `{
	"results": [
		{
			"id": 1,
			"type": "exchange",
			"asset_class": "crypto",
			"locale": "global",
			"name": "Coinbase",
			"acronym": "COINBASE",
			"url": "https://www.coinbase.com"
		},
		{
			"id": 2,
			"type": "exchange",
			"asset_class": "crypto",
			"locale": "global",
			"name": "Binance",
			"acronym": "BINANCE",
			"url": "https://www.binance.com"
		}
	],
	"status": "OK",
	"request_id": "crypto-exchanges-123",
	"count": 2
}`

const cryptoSnapshotFullMarketJSON = `{
	"status": "OK",
	"request_id": "crypto-snap-market-123",
	"tickers": [
		{
			"ticker": "X:BTCUSD",
			"todaysChange": 500.00,
			"todaysChangePerc": 1.16,
			"updated": 1736225999000,
			"day": {
				"o": 43000.00,
				"h": 43800.00,
				"l": 42900.00,
				"c": 43500.00,
				"v": 123456.78,
				"vw": 43250.50
			},
			"prevDay": {
				"o": 42000.00,
				"h": 43100.00,
				"l": 41900.00,
				"c": 43000.00,
				"v": 250000.00,
				"vw": 42500.00
			},
			"min": {
				"o": 43480.00,
				"h": 43510.00,
				"l": 43470.00,
				"c": 43500.00,
				"v": 150.00,
				"vw": 43490.00,
				"t": 1736225940000,
				"n": 25,
				"av": 123456.78
			},
			"lastTrade": {
				"conditions": [1],
				"exchange": 1,
				"price": 43500.00,
				"size": 0.5,
				"timestamp": 1736225999000
			},
			"fmv": 43495.00
		},
		{
			"ticker": "X:ETHUSD",
			"todaysChange": 100.00,
			"todaysChangePerc": 4.55,
			"updated": 1736225999000,
			"day": {
				"o": 2200.00,
				"h": 2350.00,
				"l": 2180.00,
				"c": 2300.00,
				"v": 1200000.00,
				"vw": 2250.75
			},
			"prevDay": {
				"o": 2100.00,
				"h": 2250.00,
				"l": 2080.00,
				"c": 2200.00,
				"v": 900000.00,
				"vw": 2150.00
			},
			"min": {
				"o": 2295.00,
				"h": 2305.00,
				"l": 2290.00,
				"c": 2300.00,
				"v": 500.00,
				"vw": 2298.00,
				"t": 1736225940000,
				"n": 12,
				"av": 1200000.00
			},
			"lastTrade": {
				"conditions": [1],
				"exchange": 2,
				"price": 2300.00,
				"size": 10.0,
				"timestamp": 1736225999000
			},
			"fmv": 2299.50
		}
	]
}`

const cryptoSnapshotSingleTickerJSON = `{
	"status": "OK",
	"request_id": "crypto-snap-single-123",
	"ticker": {
		"ticker": "X:BTCUSD",
		"todaysChange": 500.00,
		"todaysChangePerc": 1.16,
		"updated": 1736225999000,
		"day": {
			"o": 43000.00,
			"h": 43800.00,
			"l": 42900.00,
			"c": 43500.00,
			"v": 123456.78,
			"vw": 43250.50
		},
		"prevDay": {
			"o": 42000.00,
			"h": 43100.00,
			"l": 41900.00,
			"c": 43000.00,
			"v": 250000.00,
			"vw": 42500.00
		},
		"min": {
			"o": 43480.00,
			"h": 43510.00,
			"l": 43470.00,
			"c": 43500.00,
			"v": 150.00,
			"vw": 43490.00,
			"t": 1736225940000,
			"n": 25,
			"av": 123456.78
		},
		"lastTrade": {
			"conditions": [1],
			"exchange": 1,
			"price": 43500.00,
			"size": 0.5,
			"timestamp": 1736225999000
		},
		"fmv": 43495.00
	}
}`

const cryptoSnapshotGainersJSON = `{
	"status": "OK",
	"request_id": "crypto-gainers-123",
	"tickers": [
		{
			"ticker": "X:SOLUSD",
			"todaysChange": 15.00,
			"todaysChangePerc": 12.50,
			"updated": 1736225999000,
			"day": {
				"o": 120.00,
				"h": 140.00,
				"l": 118.00,
				"c": 135.00,
				"v": 5000000.00,
				"vw": 130.00
			},
			"prevDay": {
				"o": 110.00,
				"h": 125.00,
				"l": 108.00,
				"c": 120.00,
				"v": 3000000.00,
				"vw": 115.00
			},
			"min": {
				"o": 134.00,
				"h": 135.50,
				"l": 133.50,
				"c": 135.00,
				"v": 10000.00,
				"vw": 134.50,
				"t": 1736225940000,
				"n": 50,
				"av": 5000000.00
			},
			"lastTrade": {
				"conditions": [1],
				"exchange": 1,
				"price": 135.00,
				"size": 100.0,
				"timestamp": 1736225999000
			},
			"fmv": 134.95
		}
	]
}`

const cryptoUnifiedSnapshotJSON = `{
	"status": "OK",
	"request_id": "crypto-unified-123",
	"results": [
		{
			"ticker": "X:BTCUSD",
			"name": "Bitcoin - United States Dollar",
			"value": 43500.00,
			"type": "crypto",
			"timeframe": "REAL-TIME",
			"market_status": "open",
			"last_updated": 1736225999000,
			"session": {
				"change": 500.00,
				"change_percent": 1.16,
				"close": 43500.00,
				"high": 43800.00,
				"low": 42900.00,
				"open": 43000.00,
				"previous_close": 43000.00
			},
			"fmv": 43495.00
		}
	]
}`

const cryptoSMAJSON = `{
	"status": "OK",
	"request_id": "crypto-sma-123",
	"results": {
		"underlying": {
			"url": "https://api.massive.com/v2/aggs/ticker/X:BTCUSD/range/1/day/2025-01-01/2025-01-10"
		},
		"values": [
			{
				"timestamp": 1736139600000,
				"value": 43250.50
			},
			{
				"timestamp": 1736226000000,
				"value": 43375.25
			}
		]
	}
}`

const cryptoEMAJSON = `{
	"status": "OK",
	"request_id": "crypto-ema-123",
	"results": {
		"underlying": {
			"url": "https://api.massive.com/v2/aggs/ticker/X:BTCUSD/range/1/day/2025-01-01/2025-01-10"
		},
		"values": [
			{
				"timestamp": 1736139600000,
				"value": 43300.75
			},
			{
				"timestamp": 1736226000000,
				"value": 43425.50
			}
		]
	}
}`

const cryptoRSIJSON = `{
	"status": "OK",
	"request_id": "crypto-rsi-123",
	"results": {
		"underlying": {
			"url": "https://api.massive.com/v2/aggs/ticker/X:BTCUSD/range/1/day/2025-01-01/2025-01-10"
		},
		"values": [
			{
				"timestamp": 1736139600000,
				"value": 62.50
			},
			{
				"timestamp": 1736226000000,
				"value": 65.30
			}
		]
	}
}`

const cryptoMACDJSON = `{
	"status": "OK",
	"request_id": "crypto-macd-123",
	"results": {
		"underlying": {
			"url": "https://api.massive.com/v2/aggs/ticker/X:BTCUSD/range/1/day/2025-01-01/2025-01-10"
		},
		"values": [
			{
				"timestamp": 1736139600000,
				"value": 150.25,
				"signal": 120.50,
				"histogram": 29.75
			},
			{
				"timestamp": 1736226000000,
				"value": 175.00,
				"signal": 135.75,
				"histogram": 39.25
			}
		]
	}
}`

const cryptoTickersJSON = `{
	"results": [
		{
			"ticker": "X:BTCUSD",
			"name": "Bitcoin - United States Dollar",
			"market": "crypto",
			"locale": "global",
			"active": true,
			"currency_name": "United States Dollar",
			"last_updated_utc": "2026-02-15T07:08:17.692Z"
		},
		{
			"ticker": "X:ETHUSD",
			"name": "Ethereum - United States Dollar",
			"market": "crypto",
			"locale": "global",
			"active": true,
			"currency_name": "United States Dollar",
			"last_updated_utc": "2026-02-15T07:08:17.692Z"
		}
	],
	"status": "OK",
	"request_id": "crypto-tickers-123",
	"count": 2,
	"next_url": "https://api.massive.com/v3/reference/tickers?cursor=YXA9Mg"
}`

const cryptoTickerOverviewJSON = `{
	"status": "OK",
	"request_id": "crypto-overview-123",
	"results": {
		"ticker": "X:BTCUSD",
		"name": "Bitcoin - United States Dollar",
		"market": "crypto",
		"locale": "global",
		"active": true,
		"currency_symbol": "USD",
		"currency_name": "United States Dollar",
		"base_currency_symbol": "BTC",
		"base_currency_name": "Bitcoin",
		"last_updated_utc": "2026-02-15T07:08:17.692Z"
	}
}`

const cryptoTradesJSON = `{
	"status": "OK",
	"request_id": "crypto-trades-123",
	"next_url": "https://api.massive.com/v3/trades/X:BTCUSD?cursor=abc",
	"results": [
		{
			"conditions": [1],
			"exchange": 1,
			"id": "trade-1",
			"participant_timestamp": 1736225999000000000,
			"price": 43500.00,
			"size": 0.5
		},
		{
			"conditions": [1, 2],
			"exchange": 2,
			"id": "trade-2",
			"participant_timestamp": 1736225998000000000,
			"price": 43499.50,
			"size": 1.25
		}
	]
}`

const cryptoLastTradeJSON = `{
	"status": "OK",
	"request_id": "crypto-last-trade-123",
	"symbol": "X:BTCUSD",
	"last": {
		"price": 43500.00,
		"size": 0.5,
		"exchange": 1,
		"conditions": [1],
		"timestamp": 1736225999000
	}
}`

// -------------------------------------------------------------------
// Aggregates Tests
// -------------------------------------------------------------------

// TestGetCryptoBars verifies that GetCryptoBars correctly parses the
// API response and returns the expected OHLC bar data for a crypto ticker.
func TestGetCryptoBars(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/aggs/ticker/X:BTCUSD/range/1/day/2025-01-06/2025-01-08": cryptoBarsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := BarsParams{
		Multiplier: "1",
		Timespan:   "day",
		From:       "2025-01-06",
		To:         "2025-01-08",
		Adjusted:   "true",
		Sort:       "asc",
		Limit:      "2",
	}

	result, err := client.GetCryptoBars("X:BTCUSD", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Ticker != "X:BTCUSD" {
		t.Errorf("expected ticker X:BTCUSD, got %s", result.Ticker)
	}

	if !result.Adjusted {
		t.Error("expected adjusted to be true")
	}

	if result.ResultsCount != 2 {
		t.Errorf("expected 2 results, got %d", result.ResultsCount)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 bars, got %d", len(result.Results))
	}

	bar := result.Results[0]
	if bar.Open != 43000.00 {
		t.Errorf("expected open 43000.00, got %f", bar.Open)
	}

	if bar.High != 43800.00 {
		t.Errorf("expected high 43800.00, got %f", bar.High)
	}

	if bar.Low != 42900.00 {
		t.Errorf("expected low 42900.00, got %f", bar.Low)
	}

	if bar.Close != 43500.00 {
		t.Errorf("expected close 43500.00, got %f", bar.Close)
	}

	if bar.Volume != 123456.78 {
		t.Errorf("expected volume 123456.78, got %f", bar.Volume)
	}

	if bar.VWAP != 43250.5 {
		t.Errorf("expected VWAP 43250.5, got %f", bar.VWAP)
	}

	if bar.Timestamp != 1736139600000 {
		t.Errorf("expected timestamp 1736139600000, got %d", bar.Timestamp)
	}

	if bar.NumTrades != 15000 {
		t.Errorf("expected 15000 trades, got %d", bar.NumTrades)
	}
}

// TestGetCryptoBarsRequestPath verifies that GetCryptoBars constructs
// the correct URL path with the crypto ticker, multiplier, timespan,
// from, and to values.
func TestGetCryptoBarsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoBarsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	params := BarsParams{
		Multiplier: "5",
		Timespan:   "minute",
		From:       "2025-01-06",
		To:         "2025-01-07",
	}

	client.GetCryptoBars("X:ETHUSD", params)

	expected := "/v2/aggs/ticker/X:ETHUSD/range/5/minute/2025-01-06/2025-01-07"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetCryptoBarsQueryParams verifies that GetCryptoBars sends the
// correct query parameters including adjusted, sort, and limit.
func TestGetCryptoBarsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("adjusted") != "false" {
			t.Errorf("expected adjusted=false, got %s", q.Get("adjusted"))
		}
		if q.Get("sort") != "desc" {
			t.Errorf("expected sort=desc, got %s", q.Get("sort"))
		}
		if q.Get("limit") != "100" {
			t.Errorf("expected limit=100, got %s", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoBarsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	params := BarsParams{
		Multiplier: "1",
		Timespan:   "day",
		From:       "2025-01-06",
		To:         "2025-01-08",
		Adjusted:   "false",
		Sort:       "desc",
		Limit:      "100",
	}

	client.GetCryptoBars("X:BTCUSD", params)
}

// TestGetCryptoBarsSecondBar verifies that the second bar in the response
// is correctly parsed with its own distinct values.
func TestGetCryptoBarsSecondBar(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/aggs/ticker/X:BTCUSD/range/1/day/2025-01-06/2025-01-08": cryptoBarsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := BarsParams{
		Multiplier: "1",
		Timespan:   "day",
		From:       "2025-01-06",
		To:         "2025-01-08",
	}

	result, err := client.GetCryptoBars("X:BTCUSD", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bar := result.Results[1]
	if bar.Open != 43500.00 {
		t.Errorf("expected open 43500.00, got %f", bar.Open)
	}

	if bar.Close != 43700.00 {
		t.Errorf("expected close 43700.00, got %f", bar.Close)
	}

	if bar.NumTrades != 12500 {
		t.Errorf("expected 12500 trades, got %d", bar.NumTrades)
	}
}

// TestGetCryptoBarsAPIError verifies that GetCryptoBars returns an error
// when the API responds with a non-200 status code.
func TestGetCryptoBarsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Data not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	params := BarsParams{
		Multiplier: "1",
		Timespan:   "day",
		From:       "2025-01-06",
		To:         "2025-01-08",
	}

	_, err := client.GetCryptoBars("X:INVALID", params)
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// TestGetCryptoDailyMarketSummary verifies that GetCryptoDailyMarketSummary
// correctly parses the grouped daily response with multiple crypto tickers.
func TestGetCryptoDailyMarketSummary(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/aggs/grouped/locale/global/market/crypto/2025-01-06": cryptoDailyMarketSummaryJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetCryptoDailyMarketSummary("2025-01-06", "true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if !result.Adjusted {
		t.Error("expected adjusted to be true")
	}

	if result.ResultsCount != 2 {
		t.Errorf("expected 2 results, got %d", result.ResultsCount)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 market summaries, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.Ticker != "X:BTCUSD" {
		t.Errorf("expected ticker X:BTCUSD, got %s", first.Ticker)
	}

	if first.Open != 43000.00 {
		t.Errorf("expected open 43000.00, got %f", first.Open)
	}

	if first.Close != 43200.00 {
		t.Errorf("expected close 43200.00, got %f", first.Close)
	}

	second := result.Results[1]
	if second.Ticker != "X:ETHUSD" {
		t.Errorf("expected ticker X:ETHUSD, got %s", second.Ticker)
	}

	if second.Volume != 1200000.0 {
		t.Errorf("expected volume 1200000, got %f", second.Volume)
	}
}

// TestGetCryptoDailyMarketSummaryRequestPath verifies that
// GetCryptoDailyMarketSummary constructs the correct API path.
func TestGetCryptoDailyMarketSummaryRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoDailyMarketSummaryJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoDailyMarketSummary("2025-06-15", "true")

	expected := "/v2/aggs/grouped/locale/global/market/crypto/2025-06-15"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetCryptoDailyMarketSummaryQueryParams verifies that the adjusted
// parameter is correctly sent to the API.
func TestGetCryptoDailyMarketSummaryQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("adjusted") != "false" {
			t.Errorf("expected adjusted=false, got %s", q.Get("adjusted"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoDailyMarketSummaryJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoDailyMarketSummary("2025-01-06", "false")
}

// TestGetCryptoDailyTickerSummary verifies that GetCryptoDailyTickerSummary
// correctly parses the open/close response for a crypto pair.
func TestGetCryptoDailyTickerSummary(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/open-close/crypto/BTC/USD/2025-01-06": cryptoDailyTickerSummaryJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetCryptoDailyTickerSummary("BTC", "USD", "2025-01-06", "true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Symbol != "X:BTCUSD" {
		t.Errorf("expected symbol X:BTCUSD, got %s", result.Symbol)
	}

	if !result.IsUTC {
		t.Error("expected isUTC to be true")
	}

	if result.Day != "2025-01-06" {
		t.Errorf("expected day 2025-01-06, got %s", result.Day)
	}

	if result.Open != 43000.00 {
		t.Errorf("expected open 43000.00, got %f", result.Open)
	}

	if result.Close != 43500.00 {
		t.Errorf("expected close 43500.00, got %f", result.Close)
	}
}

// TestGetCryptoDailyTickerSummaryOpenTrades verifies that the open trades
// array in the daily ticker summary is correctly parsed.
func TestGetCryptoDailyTickerSummaryOpenTrades(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/open-close/crypto/BTC/USD/2025-01-06": cryptoDailyTickerSummaryJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetCryptoDailyTickerSummary("BTC", "USD", "2025-01-06", "true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.OpenTrades) != 1 {
		t.Fatalf("expected 1 open trade, got %d", len(result.OpenTrades))
	}

	trade := result.OpenTrades[0]
	if trade.ID != "trade-open-1" {
		t.Errorf("expected trade ID trade-open-1, got %s", trade.ID)
	}

	if trade.Price != 43000.00 {
		t.Errorf("expected price 43000.00, got %f", trade.Price)
	}

	if trade.Size != 0.5 {
		t.Errorf("expected size 0.5, got %f", trade.Size)
	}

	if trade.Exchange != 1 {
		t.Errorf("expected exchange 1, got %d", trade.Exchange)
	}

	if len(trade.Conditions) != 2 {
		t.Errorf("expected 2 conditions, got %d", len(trade.Conditions))
	}
}

// TestGetCryptoDailyTickerSummaryClosingTrades verifies that the closing
// trades array in the daily ticker summary is correctly parsed.
func TestGetCryptoDailyTickerSummaryClosingTrades(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/open-close/crypto/BTC/USD/2025-01-06": cryptoDailyTickerSummaryJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetCryptoDailyTickerSummary("BTC", "USD", "2025-01-06", "true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.ClosingTrades) != 1 {
		t.Fatalf("expected 1 closing trade, got %d", len(result.ClosingTrades))
	}

	trade := result.ClosingTrades[0]
	if trade.ID != "trade-close-1" {
		t.Errorf("expected trade ID trade-close-1, got %s", trade.ID)
	}

	if trade.Price != 43500.00 {
		t.Errorf("expected price 43500.00, got %f", trade.Price)
	}

	if trade.Size != 1.2 {
		t.Errorf("expected size 1.2, got %f", trade.Size)
	}
}

// TestGetCryptoDailyTickerSummaryRequestPath verifies that the API path
// is correctly constructed with the from/to currency pair and date.
func TestGetCryptoDailyTickerSummaryRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoDailyTickerSummaryJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoDailyTickerSummary("ETH", "USD", "2025-03-15", "true")

	expected := "/v1/open-close/crypto/ETH/USD/2025-03-15"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetCryptoPreviousDayBar verifies that GetCryptoPreviousDayBar
// correctly parses the API response for the previous day's bar data.
func TestGetCryptoPreviousDayBar(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/aggs/ticker/X:BTCUSD/prev": cryptoPreviousDayBarJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetCryptoPreviousDayBar("X:BTCUSD", "true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Ticker != "X:BTCUSD" {
		t.Errorf("expected ticker X:BTCUSD, got %s", result.Ticker)
	}

	if result.ResultsCount != 1 {
		t.Errorf("expected 1 result, got %d", result.ResultsCount)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 bar, got %d", len(result.Results))
	}

	bar := result.Results[0]
	if bar.Open != 42000.00 {
		t.Errorf("expected open 42000.00, got %f", bar.Open)
	}

	if bar.Close != 43000.00 {
		t.Errorf("expected close 43000.00, got %f", bar.Close)
	}

	if bar.High != 43100.00 {
		t.Errorf("expected high 43100.00, got %f", bar.High)
	}

	if bar.Low != 41900.00 {
		t.Errorf("expected low 41900.00, got %f", bar.Low)
	}

	if bar.Volume != 250000.0 {
		t.Errorf("expected volume 250000, got %f", bar.Volume)
	}
}

// TestGetCryptoPreviousDayBarRequestPath verifies that
// GetCryptoPreviousDayBar constructs the correct API path.
func TestGetCryptoPreviousDayBarRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoPreviousDayBarJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoPreviousDayBar("X:ETHUSD", "true")

	expected := "/v2/aggs/ticker/X:ETHUSD/prev"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetCryptoPreviousDayBarAdjustedParam verifies that the adjusted
// query parameter is correctly sent to the API.
func TestGetCryptoPreviousDayBarAdjustedParam(t *testing.T) {
	var receivedAdjusted string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAdjusted = r.URL.Query().Get("adjusted")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoPreviousDayBarJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoPreviousDayBar("X:BTCUSD", "false")

	if receivedAdjusted != "false" {
		t.Errorf("expected adjusted=false, got %s", receivedAdjusted)
	}
}

// -------------------------------------------------------------------
// Market Operations Tests
// -------------------------------------------------------------------

// TestGetCryptoConditions verifies that GetCryptoConditions correctly
// parses the API response containing crypto condition codes.
func TestGetCryptoConditions(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/reference/conditions": cryptoConditionsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetCryptoConditions()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Count != 2 {
		t.Errorf("expected count 2, got %d", result.Count)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 conditions, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.ID != 1 {
		t.Errorf("expected id 1, got %d", first.ID)
	}

	if first.Name != "Regular Sale" {
		t.Errorf("expected name Regular Sale, got %s", first.Name)
	}

	if first.AssetClass != "crypto" {
		t.Errorf("expected asset_class crypto, got %s", first.AssetClass)
	}

	if first.Type != "sale_condition" {
		t.Errorf("expected type sale_condition, got %s", first.Type)
	}
}

// TestGetCryptoConditionsRequestPath verifies that GetCryptoConditions
// sends the request to the correct API path with asset_class=crypto.
func TestGetCryptoConditionsRequestPath(t *testing.T) {
	var receivedPath string
	var receivedAssetClass string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedAssetClass = r.URL.Query().Get("asset_class")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoConditionsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoConditions()

	if receivedPath != "/v3/reference/conditions" {
		t.Errorf("expected path /v3/reference/conditions, got %s", receivedPath)
	}

	if receivedAssetClass != "crypto" {
		t.Errorf("expected asset_class=crypto, got %s", receivedAssetClass)
	}
}

// TestGetCryptoConditionsAPIError verifies that GetCryptoConditions
// returns an error when the API responds with a non-200 status.
func TestGetCryptoConditionsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"status":"ERROR","message":"Unauthorized"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetCryptoConditions()
	if err == nil {
		t.Fatal("expected error for 401 response, got nil")
	}
}

// TestGetCryptoExchanges verifies that GetCryptoExchanges correctly
// parses the API response containing crypto exchange information.
func TestGetCryptoExchanges(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/reference/exchanges": cryptoExchangesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetCryptoExchanges()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Count != 2 {
		t.Errorf("expected count 2, got %d", result.Count)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 exchanges, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.ID != 1 {
		t.Errorf("expected id 1, got %d", first.ID)
	}

	if first.Name != "Coinbase" {
		t.Errorf("expected name Coinbase, got %s", first.Name)
	}

	if first.AssetClass != "crypto" {
		t.Errorf("expected asset_class crypto, got %s", first.AssetClass)
	}

	second := result.Results[1]
	if second.Name != "Binance" {
		t.Errorf("expected name Binance, got %s", second.Name)
	}

	if second.Acronym != "BINANCE" {
		t.Errorf("expected acronym BINANCE, got %s", second.Acronym)
	}
}

// TestGetCryptoExchangesRequestPath verifies that GetCryptoExchanges
// sends the request with asset_class=crypto.
func TestGetCryptoExchangesRequestPath(t *testing.T) {
	var receivedAssetClass string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAssetClass = r.URL.Query().Get("asset_class")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoExchangesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoExchanges()

	if receivedAssetClass != "crypto" {
		t.Errorf("expected asset_class=crypto, got %s", receivedAssetClass)
	}
}

// -------------------------------------------------------------------
// Snapshot Tests
// -------------------------------------------------------------------

// TestGetCryptoSnapshotFullMarket verifies that GetCryptoSnapshotFullMarket
// correctly parses the full market snapshot response with multiple tickers.
func TestGetCryptoSnapshotFullMarket(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/snapshot/locale/global/markets/crypto/tickers": cryptoSnapshotFullMarketJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetCryptoSnapshotFullMarket(CryptoSnapshotParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Tickers) != 2 {
		t.Fatalf("expected 2 tickers, got %d", len(result.Tickers))
	}

	btc := result.Tickers[0]
	if btc.Ticker != "X:BTCUSD" {
		t.Errorf("expected ticker X:BTCUSD, got %s", btc.Ticker)
	}

	if btc.TodaysChange != 500.00 {
		t.Errorf("expected todaysChange 500.00, got %f", btc.TodaysChange)
	}

	if btc.TodaysChangePct != 1.16 {
		t.Errorf("expected todaysChangePerc 1.16, got %f", btc.TodaysChangePct)
	}

	if btc.Day.Open != 43000.00 {
		t.Errorf("expected day open 43000.00, got %f", btc.Day.Open)
	}

	if btc.Day.Close != 43500.00 {
		t.Errorf("expected day close 43500.00, got %f", btc.Day.Close)
	}

	if btc.PrevDay.Close != 43000.00 {
		t.Errorf("expected prevDay close 43000.00, got %f", btc.PrevDay.Close)
	}

	if btc.LastTrade.Price != 43500.00 {
		t.Errorf("expected lastTrade price 43500.00, got %f", btc.LastTrade.Price)
	}

	if btc.FMV != 43495.00 {
		t.Errorf("expected fmv 43495.00, got %f", btc.FMV)
	}

	eth := result.Tickers[1]
	if eth.Ticker != "X:ETHUSD" {
		t.Errorf("expected ticker X:ETHUSD, got %s", eth.Ticker)
	}

	if eth.TodaysChangePct != 4.55 {
		t.Errorf("expected todaysChangePerc 4.55, got %f", eth.TodaysChangePct)
	}
}

// TestGetCryptoSnapshotFullMarketWithTickers verifies that the tickers
// query parameter is correctly sent when filtering the snapshot.
func TestGetCryptoSnapshotFullMarketWithTickers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("tickers") != "X:BTCUSD,X:ETHUSD" {
			t.Errorf("expected tickers=X:BTCUSD,X:ETHUSD, got %s", q.Get("tickers"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoSnapshotFullMarketJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoSnapshotFullMarket(CryptoSnapshotParams{
		Tickers: "X:BTCUSD,X:ETHUSD",
	})
}

// TestGetCryptoSnapshotSingleTicker verifies that GetCryptoSnapshotSingleTicker
// correctly parses the single ticker snapshot response.
func TestGetCryptoSnapshotSingleTicker(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/snapshot/locale/global/markets/crypto/tickers/X:BTCUSD": cryptoSnapshotSingleTickerJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetCryptoSnapshotSingleTicker("X:BTCUSD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Ticker.Ticker != "X:BTCUSD" {
		t.Errorf("expected ticker X:BTCUSD, got %s", result.Ticker.Ticker)
	}

	if result.Ticker.TodaysChange != 500.00 {
		t.Errorf("expected todaysChange 500.00, got %f", result.Ticker.TodaysChange)
	}

	if result.Ticker.Day.Open != 43000.00 {
		t.Errorf("expected day open 43000.00, got %f", result.Ticker.Day.Open)
	}

	if result.Ticker.LastTrade.Price != 43500.00 {
		t.Errorf("expected lastTrade price 43500.00, got %f", result.Ticker.LastTrade.Price)
	}

	if result.Ticker.FMV != 43495.00 {
		t.Errorf("expected fmv 43495.00, got %f", result.Ticker.FMV)
	}

	if result.Ticker.Min.NumTransactions != 25 {
		t.Errorf("expected min numTransactions 25, got %d", result.Ticker.Min.NumTransactions)
	}
}

// TestGetCryptoSnapshotSingleTickerRequestPath verifies the correct
// API path is constructed with the crypto ticker.
func TestGetCryptoSnapshotSingleTickerRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoSnapshotSingleTickerJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoSnapshotSingleTicker("X:ETHUSD")

	expected := "/v2/snapshot/locale/global/markets/crypto/tickers/X:ETHUSD"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetCryptoSnapshotTopMoversGainers verifies that GetCryptoSnapshotTopMovers
// correctly retrieves and parses the top gainers snapshot.
func TestGetCryptoSnapshotTopMoversGainers(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/snapshot/locale/global/markets/crypto/gainers": cryptoSnapshotGainersJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetCryptoSnapshotTopMovers("gainers")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Tickers) != 1 {
		t.Fatalf("expected 1 ticker, got %d", len(result.Tickers))
	}

	sol := result.Tickers[0]
	if sol.Ticker != "X:SOLUSD" {
		t.Errorf("expected ticker X:SOLUSD, got %s", sol.Ticker)
	}

	if sol.TodaysChangePct != 12.50 {
		t.Errorf("expected todaysChangePerc 12.50, got %f", sol.TodaysChangePct)
	}

	if sol.Day.Close != 135.00 {
		t.Errorf("expected day close 135.00, got %f", sol.Day.Close)
	}
}

// TestGetCryptoSnapshotTopMoversLosersPath verifies that the losers
// direction correctly constructs the API path.
func TestGetCryptoSnapshotTopMoversLosersPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoSnapshotGainersJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoSnapshotTopMovers("losers")

	expected := "/v2/snapshot/locale/global/markets/crypto/losers"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetCryptoUnifiedSnapshot verifies that GetCryptoUnifiedSnapshot
// correctly parses the unified snapshot response.
func TestGetCryptoUnifiedSnapshot(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot": cryptoUnifiedSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetCryptoUnifiedSnapshot(CryptoUnifiedSnapshotParams{
		TickerAnyOf: "X:BTCUSD",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	snap := result.Results[0]
	if snap.Ticker != "X:BTCUSD" {
		t.Errorf("expected ticker X:BTCUSD, got %s", snap.Ticker)
	}

	if snap.Name != "Bitcoin - United States Dollar" {
		t.Errorf("expected name Bitcoin - United States Dollar, got %s", snap.Name)
	}

	if snap.Value != 43500.00 {
		t.Errorf("expected value 43500.00, got %f", snap.Value)
	}

	if snap.MarketStatus != "open" {
		t.Errorf("expected market_status open, got %s", snap.MarketStatus)
	}

	if snap.Session.Change != 500.00 {
		t.Errorf("expected session change 500.00, got %f", snap.Session.Change)
	}

	if snap.Session.ChangePercent != 1.16 {
		t.Errorf("expected session change_percent 1.16, got %f", snap.Session.ChangePercent)
	}

	if snap.Session.Close != 43500.00 {
		t.Errorf("expected session close 43500.00, got %f", snap.Session.Close)
	}

	if snap.Session.PreviousClose != 43000.00 {
		t.Errorf("expected session previous_close 43000.00, got %f", snap.Session.PreviousClose)
	}

	if snap.FMV != 43495.00 {
		t.Errorf("expected fmv 43495.00, got %f", snap.FMV)
	}
}

// TestGetCryptoUnifiedSnapshotQueryParams verifies that the ticker.any_of
// query parameter is correctly sent to the API.
func TestGetCryptoUnifiedSnapshotQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("ticker.any_of") != "X:BTCUSD,X:ETHUSD" {
			t.Errorf("expected ticker.any_of=X:BTCUSD,X:ETHUSD, got %s", q.Get("ticker.any_of"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoUnifiedSnapshotJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoUnifiedSnapshot(CryptoUnifiedSnapshotParams{
		TickerAnyOf: "X:BTCUSD,X:ETHUSD",
	})
}

// -------------------------------------------------------------------
// Technical Indicator Tests
// -------------------------------------------------------------------

// TestGetCryptoSMA verifies that GetCryptoSMA correctly parses the
// SMA indicator response for a crypto ticker.
func TestGetCryptoSMA(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/sma/X:BTCUSD": cryptoSMAJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := IndicatorParams{
		TimestampGTE: "2025-01-01",
		TimestampLTE: "2025-01-10",
		Timespan:     "day",
		Adjusted:     "true",
		Window:       "10",
		SeriesType:   "close",
		Order:        "desc",
		Limit:        "10",
	}

	result, err := client.GetCryptoSMA("X:BTCUSD", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Results.Values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(result.Results.Values))
	}

	first := result.Results.Values[0]
	if first.Timestamp != 1736139600000 {
		t.Errorf("expected timestamp 1736139600000, got %d", first.Timestamp)
	}

	if first.Value != 43250.50 {
		t.Errorf("expected value 43250.50, got %f", first.Value)
	}
}

// TestGetCryptoSMARequestPath verifies that GetCryptoSMA constructs
// the correct API path for the crypto ticker.
func TestGetCryptoSMARequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoSMAJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoSMA("X:ETHUSD", IndicatorParams{})

	if receivedPath != "/v1/indicators/sma/X:ETHUSD" {
		t.Errorf("expected path /v1/indicators/sma/X:ETHUSD, got %s", receivedPath)
	}
}

// TestGetCryptoSMAQueryParams verifies that the indicator query parameters
// are correctly sent to the API endpoint.
func TestGetCryptoSMAQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("timestamp.gte") != "2025-01-01" {
			t.Errorf("expected timestamp.gte=2025-01-01, got %s", q.Get("timestamp.gte"))
		}
		if q.Get("timestamp.lte") != "2025-01-10" {
			t.Errorf("expected timestamp.lte=2025-01-10, got %s", q.Get("timestamp.lte"))
		}
		if q.Get("timespan") != "day" {
			t.Errorf("expected timespan=day, got %s", q.Get("timespan"))
		}
		if q.Get("window") != "20" {
			t.Errorf("expected window=20, got %s", q.Get("window"))
		}
		if q.Get("series_type") != "close" {
			t.Errorf("expected series_type=close, got %s", q.Get("series_type"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoSMAJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoSMA("X:BTCUSD", IndicatorParams{
		TimestampGTE: "2025-01-01",
		TimestampLTE: "2025-01-10",
		Timespan:     "day",
		Window:       "20",
		SeriesType:   "close",
	})
}

// TestGetCryptoEMA verifies that GetCryptoEMA correctly parses the
// EMA indicator response for a crypto ticker.
func TestGetCryptoEMA(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/ema/X:BTCUSD": cryptoEMAJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetCryptoEMA("X:BTCUSD", IndicatorParams{
		TimestampGTE: "2025-01-01",
		TimestampLTE: "2025-01-10",
		Timespan:     "day",
		Window:       "10",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Results.Values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(result.Results.Values))
	}

	first := result.Results.Values[0]
	if first.Value != 43300.75 {
		t.Errorf("expected value 43300.75, got %f", first.Value)
	}
}

// TestGetCryptoEMARequestPath verifies that GetCryptoEMA constructs
// the correct API path for the crypto ticker.
func TestGetCryptoEMARequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoEMAJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoEMA("X:BTCUSD", IndicatorParams{})

	if receivedPath != "/v1/indicators/ema/X:BTCUSD" {
		t.Errorf("expected path /v1/indicators/ema/X:BTCUSD, got %s", receivedPath)
	}
}

// TestGetCryptoRSI verifies that GetCryptoRSI correctly parses the
// RSI indicator response for a crypto ticker.
func TestGetCryptoRSI(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/rsi/X:BTCUSD": cryptoRSIJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetCryptoRSI("X:BTCUSD", IndicatorParams{
		TimestampGTE: "2025-01-01",
		TimestampLTE: "2025-01-10",
		Timespan:     "day",
		Window:       "14",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Results.Values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(result.Results.Values))
	}

	first := result.Results.Values[0]
	if first.Value != 62.50 {
		t.Errorf("expected value 62.50, got %f", first.Value)
	}

	second := result.Results.Values[1]
	if second.Value != 65.30 {
		t.Errorf("expected value 65.30, got %f", second.Value)
	}
}

// TestGetCryptoRSIRequestPath verifies that GetCryptoRSI constructs
// the correct API path for the crypto ticker.
func TestGetCryptoRSIRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoRSIJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoRSI("X:BTCUSD", IndicatorParams{})

	if receivedPath != "/v1/indicators/rsi/X:BTCUSD" {
		t.Errorf("expected path /v1/indicators/rsi/X:BTCUSD, got %s", receivedPath)
	}
}

// TestGetCryptoMACD verifies that GetCryptoMACD correctly parses the
// MACD indicator response for a crypto ticker including MACD line,
// signal line, and histogram values.
func TestGetCryptoMACD(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/macd/X:BTCUSD": cryptoMACDJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := MACDParams{
		TimestampGTE: "2025-01-01",
		TimestampLTE: "2025-01-10",
		Timespan:     "day",
		Adjusted:     "true",
		ShortWindow:  "12",
		LongWindow:   "26",
		SignalWindow: "9",
		SeriesType:   "close",
		Order:        "desc",
		Limit:        "10",
	}

	result, err := client.GetCryptoMACD("X:BTCUSD", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Results.Values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(result.Results.Values))
	}

	first := result.Results.Values[0]
	if first.Value != 150.25 {
		t.Errorf("expected value 150.25, got %f", first.Value)
	}

	if first.Signal != 120.50 {
		t.Errorf("expected signal 120.50, got %f", first.Signal)
	}

	if first.Histogram != 29.75 {
		t.Errorf("expected histogram 29.75, got %f", first.Histogram)
	}

	second := result.Results.Values[1]
	if second.Value != 175.00 {
		t.Errorf("expected value 175.00, got %f", second.Value)
	}

	if second.Signal != 135.75 {
		t.Errorf("expected signal 135.75, got %f", second.Signal)
	}

	if second.Histogram != 39.25 {
		t.Errorf("expected histogram 39.25, got %f", second.Histogram)
	}
}

// TestGetCryptoMACDRequestPath verifies that GetCryptoMACD constructs
// the correct API path for the crypto ticker.
func TestGetCryptoMACDRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoMACDJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoMACD("X:BTCUSD", MACDParams{})

	if receivedPath != "/v1/indicators/macd/X:BTCUSD" {
		t.Errorf("expected path /v1/indicators/macd/X:BTCUSD, got %s", receivedPath)
	}
}

// TestGetCryptoMACDQueryParams verifies that the MACD-specific query
// parameters (short_window, long_window, signal_window) are correctly
// sent to the API endpoint.
func TestGetCryptoMACDQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("short_window") != "12" {
			t.Errorf("expected short_window=12, got %s", q.Get("short_window"))
		}
		if q.Get("long_window") != "26" {
			t.Errorf("expected long_window=26, got %s", q.Get("long_window"))
		}
		if q.Get("signal_window") != "9" {
			t.Errorf("expected signal_window=9, got %s", q.Get("signal_window"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoMACDJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoMACD("X:BTCUSD", MACDParams{
		ShortWindow:  "12",
		LongWindow:   "26",
		SignalWindow: "9",
	})
}

// -------------------------------------------------------------------
// Tickers Tests
// -------------------------------------------------------------------

// TestGetCryptoTickers verifies that GetCryptoTickers correctly parses
// the reference tickers response filtered by market=crypto.
func TestGetCryptoTickers(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/reference/tickers": cryptoTickersJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetCryptoTickers(CryptoTickersParams{
		Search: "Bitcoin",
		Limit:  "2",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Count != 2 {
		t.Errorf("expected count 2, got %d", result.Count)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 tickers, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.Ticker != "X:BTCUSD" {
		t.Errorf("expected ticker X:BTCUSD, got %s", first.Ticker)
	}

	if first.Name != "Bitcoin - United States Dollar" {
		t.Errorf("expected name Bitcoin - United States Dollar, got %s", first.Name)
	}

	if first.Market != "crypto" {
		t.Errorf("expected market crypto, got %s", first.Market)
	}

	if !first.Active {
		t.Error("expected active to be true")
	}

	second := result.Results[1]
	if second.Ticker != "X:ETHUSD" {
		t.Errorf("expected ticker X:ETHUSD, got %s", second.Ticker)
	}
}

// TestGetCryptoTickersQueryParams verifies that the market=crypto
// parameter is always sent along with the user-provided filters.
func TestGetCryptoTickersQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("market") != "crypto" {
			t.Errorf("expected market=crypto, got %s", q.Get("market"))
		}
		if q.Get("search") != "Ethereum" {
			t.Errorf("expected search=Ethereum, got %s", q.Get("search"))
		}
		if q.Get("active") != "true" {
			t.Errorf("expected active=true, got %s", q.Get("active"))
		}
		if q.Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "ticker" {
			t.Errorf("expected sort=ticker, got %s", q.Get("sort"))
		}
		if q.Get("order") != "asc" {
			t.Errorf("expected order=asc, got %s", q.Get("order"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoTickersJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoTickers(CryptoTickersParams{
		Search: "Ethereum",
		Active: "true",
		Limit:  "50",
		Sort:   "ticker",
		Order:  "asc",
	})
}

// TestGetCryptoTickersEmptyResults verifies that GetCryptoTickers handles
// an empty results array without error.
func TestGetCryptoTickersEmptyResults(t *testing.T) {
	emptyJSON := `{"results":[],"status":"OK","request_id":"abc","count":0}`
	server := mockServer(t, map[string]string{
		"/v3/reference/tickers": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetCryptoTickers(CryptoTickersParams{Search: "zzzznotreal"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Count != 0 {
		t.Errorf("expected count 0, got %d", result.Count)
	}

	if len(result.Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(result.Results))
	}
}

// TestGetCryptoTickerOverview verifies that GetCryptoTickerOverview
// correctly parses the detailed reference information for a crypto ticker.
func TestGetCryptoTickerOverview(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/reference/tickers/X:BTCUSD": cryptoTickerOverviewJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetCryptoTickerOverview("X:BTCUSD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	r := result.Results
	if r.Ticker != "X:BTCUSD" {
		t.Errorf("expected ticker X:BTCUSD, got %s", r.Ticker)
	}

	if r.Name != "Bitcoin - United States Dollar" {
		t.Errorf("expected name Bitcoin - United States Dollar, got %s", r.Name)
	}

	if r.Market != "crypto" {
		t.Errorf("expected market crypto, got %s", r.Market)
	}

	if !r.Active {
		t.Error("expected active to be true")
	}

	if r.CurrencyName != "United States Dollar" {
		t.Errorf("expected currency_name United States Dollar, got %s", r.CurrencyName)
	}

	if r.BaseCurrencySymbol != "BTC" {
		t.Errorf("expected base_currency_symbol BTC, got %s", r.BaseCurrencySymbol)
	}

	if r.BaseCurrencyName != "Bitcoin" {
		t.Errorf("expected base_currency_name Bitcoin, got %s", r.BaseCurrencyName)
	}
}

// TestGetCryptoTickerOverviewRequestPath verifies that
// GetCryptoTickerOverview constructs the correct API path.
func TestGetCryptoTickerOverviewRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoTickerOverviewJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoTickerOverview("X:ETHUSD")

	expected := "/v3/reference/tickers/X:ETHUSD"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetCryptoTickerOverviewAPIError verifies that GetCryptoTickerOverview
// returns an error when the API responds with a non-200 status.
func TestGetCryptoTickerOverviewAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Ticker not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetCryptoTickerOverview("X:INVALID")
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// -------------------------------------------------------------------
// Trades Tests
// -------------------------------------------------------------------

// TestGetCryptoTrades verifies that GetCryptoTrades correctly parses
// the tick-level trade data response for a crypto ticker.
func TestGetCryptoTrades(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/trades/X:BTCUSD": cryptoTradesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := CryptoTradesParams{
		Limit: "10",
		Order: "desc",
	}

	result, err := client.GetCryptoTrades("X:BTCUSD", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 trades, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.ID != "trade-1" {
		t.Errorf("expected trade ID trade-1, got %s", first.ID)
	}

	if first.Price != 43500.00 {
		t.Errorf("expected price 43500.00, got %f", first.Price)
	}

	if first.Size != 0.5 {
		t.Errorf("expected size 0.5, got %f", first.Size)
	}

	if first.Exchange != 1 {
		t.Errorf("expected exchange 1, got %d", first.Exchange)
	}

	second := result.Results[1]
	if second.ID != "trade-2" {
		t.Errorf("expected trade ID trade-2, got %s", second.ID)
	}

	if second.Price != 43499.50 {
		t.Errorf("expected price 43499.50, got %f", second.Price)
	}

	if len(second.Conditions) != 2 {
		t.Errorf("expected 2 conditions, got %d", len(second.Conditions))
	}
}

// TestGetCryptoTradesRequestPath verifies that GetCryptoTrades constructs
// the correct API path with the crypto ticker.
func TestGetCryptoTradesRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoTradesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoTrades("X:ETHUSD", CryptoTradesParams{})

	expected := "/v3/trades/X:ETHUSD"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetCryptoTradesQueryParams verifies that all timestamp filter,
// ordering, and pagination parameters are correctly sent to the API.
func TestGetCryptoTradesQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("timestamp") != "2025-01-06" {
			t.Errorf("expected timestamp=2025-01-06, got %s", q.Get("timestamp"))
		}
		if q.Get("order") != "asc" {
			t.Errorf("expected order=asc, got %s", q.Get("order"))
		}
		if q.Get("limit") != "500" {
			t.Errorf("expected limit=500, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "timestamp" {
			t.Errorf("expected sort=timestamp, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoTradesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoTrades("X:BTCUSD", CryptoTradesParams{
		Timestamp: "2025-01-06",
		Order:     "asc",
		Limit:     "500",
		Sort:      "timestamp",
	})
}

// TestGetCryptoTradesTimestampFilters verifies that the timestamp range
// filter parameters (gte, gt, lte, lt) are correctly sent.
func TestGetCryptoTradesTimestampFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("timestamp.gte") != "2025-01-01" {
			t.Errorf("expected timestamp.gte=2025-01-01, got %s", q.Get("timestamp.gte"))
		}
		if q.Get("timestamp.lte") != "2025-01-31" {
			t.Errorf("expected timestamp.lte=2025-01-31, got %s", q.Get("timestamp.lte"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoTradesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoTrades("X:BTCUSD", CryptoTradesParams{
		TimestampGte: "2025-01-01",
		TimestampLte: "2025-01-31",
	})
}

// TestGetCryptoTradesAPIError verifies that GetCryptoTrades returns an
// error when the API responds with a non-200 status code.
func TestGetCryptoTradesAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"ERROR","message":"Forbidden"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetCryptoTrades("X:BTCUSD", CryptoTradesParams{})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetCryptoLastTrade verifies that GetCryptoLastTrade correctly
// parses the most recent trade response for a crypto pair.
func TestGetCryptoLastTrade(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/last/crypto/BTC/USD": cryptoLastTradeJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetCryptoLastTrade("BTC", "USD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Symbol != "X:BTCUSD" {
		t.Errorf("expected symbol X:BTCUSD, got %s", result.Symbol)
	}

	if result.Last.Price != 43500.00 {
		t.Errorf("expected price 43500.00, got %f", result.Last.Price)
	}

	if result.Last.Size != 0.5 {
		t.Errorf("expected size 0.5, got %f", result.Last.Size)
	}

	if result.Last.Exchange != 1 {
		t.Errorf("expected exchange 1, got %d", result.Last.Exchange)
	}

	if result.Last.Timestamp != 1736225999000 {
		t.Errorf("expected timestamp 1736225999000, got %d", result.Last.Timestamp)
	}

	if len(result.Last.Conditions) != 1 {
		t.Errorf("expected 1 condition, got %d", len(result.Last.Conditions))
	}
}

// TestGetCryptoLastTradeRequestPath verifies that GetCryptoLastTrade
// constructs the correct API path with the from/to currency pair.
func TestGetCryptoLastTradeRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cryptoLastTradeJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCryptoLastTrade("ETH", "USD")

	expected := "/v1/last/crypto/ETH/USD"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetCryptoLastTradeAPIError verifies that GetCryptoLastTrade returns
// an error when the API responds with a non-200 status code.
func TestGetCryptoLastTradeAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"No data found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetCryptoLastTrade("INVALID", "USD")
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}
