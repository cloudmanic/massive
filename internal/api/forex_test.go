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

// --- Test JSON fixtures ---

const forexBarsJSON = `{
	"ticker": "C:EURUSD",
	"queryCount": 2,
	"resultsCount": 2,
	"adjusted": true,
	"results": [
		{
			"v": 50000,
			"vw": 1.0855,
			"o": 1.0850,
			"c": 1.0860,
			"h": 1.0870,
			"l": 1.0840,
			"t": 1736139600000,
			"n": 1200
		},
		{
			"v": 45000,
			"vw": 1.0862,
			"o": 1.0860,
			"c": 1.0855,
			"h": 1.0880,
			"l": 1.0845,
			"t": 1736226000000,
			"n": 1100
		}
	],
	"status": "OK",
	"request_id": "forex_bars_req_001"
}`

const forexMarketSummaryJSON = `{
	"queryCount": 2,
	"resultsCount": 2,
	"adjusted": true,
	"results": [
		{
			"T": "C:EURUSD",
			"v": 120000,
			"vw": 1.0855,
			"o": 1.0850,
			"c": 1.0860,
			"h": 1.0870,
			"l": 1.0840,
			"t": 1736197200000,
			"n": 5000
		},
		{
			"T": "C:GBPUSD",
			"v": 95000,
			"vw": 1.2710,
			"o": 1.2700,
			"c": 1.2720,
			"h": 1.2740,
			"l": 1.2690,
			"t": 1736197200000,
			"n": 4200
		}
	],
	"status": "OK",
	"request_id": "forex_summary_req_001"
}`

const forexPrevDayJSON = `{
	"ticker": "C:EURUSD",
	"queryCount": 1,
	"resultsCount": 1,
	"adjusted": true,
	"results": [
		{
			"v": 60000,
			"vw": 1.0848,
			"o": 1.0845,
			"c": 1.0850,
			"h": 1.0865,
			"l": 1.0835,
			"t": 1736139600000,
			"n": 1500
		}
	],
	"status": "OK",
	"request_id": "forex_prev_req_001"
}`

const forexConversionJSON = `{
	"converted": 108.50,
	"from": "USD",
	"initialAmount": 100,
	"last": {
		"ask": 1.0855,
		"bid": 1.0850,
		"exchange": 48,
		"timestamp": 1736139600000
	},
	"request_id": "forex_conv_req_001",
	"status": "OK",
	"symbol": "USD/EUR",
	"to": "EUR"
}`

const forexQuotesJSON = `{
	"status": "OK",
	"request_id": "forex_quotes_req_001",
	"next_url": "https://api.massive.com/v3/quotes/C:EURUSD?cursor=abc",
	"results": [
		{
			"ask_exchange": 48,
			"ask_price": 1.0855,
			"bid_exchange": 48,
			"bid_price": 1.0850,
			"participant_timestamp": 1736139600000
		},
		{
			"ask_exchange": 48,
			"ask_price": 1.0860,
			"bid_exchange": 48,
			"bid_price": 1.0852,
			"participant_timestamp": 1736139601000
		}
	]
}`

const forexLastQuoteJSON = `{
	"last": {
		"ask": 1.0855,
		"bid": 1.0850,
		"exchange": 48,
		"timestamp": 1736139600000
	},
	"request_id": "forex_lq_req_001",
	"status": "OK",
	"symbol": "EUR/USD"
}`

const forexSnapshotAllJSON = `{
	"status": "OK",
	"request_id": "forex_snap_all_001",
	"count": 2,
	"tickers": [
		{
			"ticker": "C:EURUSD",
			"todaysChange": 0.0010,
			"todaysChangePerc": 0.092,
			"updated": 1736139600000,
			"day": {"o": 1.0850, "h": 1.0870, "l": 1.0840, "c": 1.0860},
			"lastQuote": {"a": 1.0855, "b": 1.0850, "x": 48, "t": 1736139600000},
			"prevDay": {"o": 1.0840, "h": 1.0860, "l": 1.0830, "c": 1.0850}
		},
		{
			"ticker": "C:GBPUSD",
			"todaysChange": 0.0020,
			"todaysChangePerc": 0.157,
			"updated": 1736139600000,
			"day": {"o": 1.2700, "h": 1.2740, "l": 1.2690, "c": 1.2720},
			"lastQuote": {"a": 1.2715, "b": 1.2710, "x": 48, "t": 1736139600000},
			"prevDay": {"o": 1.2690, "h": 1.2720, "l": 1.2680, "c": 1.2700}
		}
	]
}`

const forexSnapshotSingleJSON = `{
	"status": "OK",
	"request_id": "forex_snap_single_001",
	"ticker": {
		"ticker": "C:EURUSD",
		"todaysChange": 0.0010,
		"todaysChangePerc": 0.092,
		"updated": 1736139600000,
		"day": {"o": 1.0850, "h": 1.0870, "l": 1.0840, "c": 1.0860},
		"lastQuote": {"a": 1.0855, "b": 1.0850, "x": 48, "t": 1736139600000},
		"prevDay": {"o": 1.0840, "h": 1.0860, "l": 1.0830, "c": 1.0850}
	}
}`

const forexGainersJSON = `{
	"status": "OK",
	"request_id": "forex_gainers_001",
	"tickers": [
		{
			"ticker": "C:EURUSD",
			"todaysChange": 0.0050,
			"todaysChangePerc": 0.461,
			"updated": 1736139600000,
			"day": {"o": 1.0850, "h": 1.0920, "l": 1.0840, "c": 1.0900},
			"lastQuote": {"a": 1.0905, "b": 1.0900, "x": 48, "t": 1736139600000},
			"prevDay": {"o": 1.0840, "h": 1.0860, "l": 1.0830, "c": 1.0850}
		}
	]
}`

const forexUnifiedSnapshotJSON = `{
	"status": "OK",
	"request_id": "forex_unified_001",
	"results": [
		{
			"ticker": "C:EURUSD",
			"todaysChange": 0.0010,
			"todaysChangePerc": 0.092,
			"updated": 1736139600000,
			"day": {"o": 1.0850, "h": 1.0870, "l": 1.0840, "c": 1.0860},
			"prevDay": {"o": 1.0840, "h": 1.0860, "l": 1.0830, "c": 1.0850}
		}
	]
}`

const forexIndicatorJSON = `{
	"status": "OK",
	"request_id": "forex_ind_001",
	"results": {
		"underlying": {"url": "https://api.massive.com/v2/aggs/ticker/C:EURUSD/range/1/day/2025-01-01/2025-01-31"},
		"values": [
			{"timestamp": 1736139600000, "value": 1.0852},
			{"timestamp": 1736226000000, "value": 1.0855}
		]
	}
}`

const forexMACDJSON = `{
	"status": "OK",
	"request_id": "forex_macd_001",
	"results": {
		"underlying": {"url": "https://api.massive.com/v2/aggs/ticker/C:EURUSD/range/1/day/2025-01-01/2025-01-31"},
		"values": [
			{"timestamp": 1736139600000, "value": 0.0012, "signal": 0.0008, "histogram": 0.0004},
			{"timestamp": 1736226000000, "value": 0.0015, "signal": 0.0010, "histogram": 0.0005}
		]
	}
}`

const forexTickersJSON = `{
	"results": [
		{
			"ticker": "C:EURUSD",
			"name": "Euro - United States Dollar",
			"market": "fx",
			"locale": "global",
			"active": true,
			"currency_name": "United States Dollar",
			"last_updated_utc": "2026-02-15T07:08:17.692Z"
		},
		{
			"ticker": "C:GBPUSD",
			"name": "British Pound - United States Dollar",
			"market": "fx",
			"locale": "global",
			"active": true,
			"currency_name": "United States Dollar",
			"last_updated_utc": "2026-02-15T07:08:17.692Z"
		}
	],
	"status": "OK",
	"request_id": "forex_tickers_001",
	"count": 2,
	"next_url": "https://api.massive.com/v3/reference/tickers?cursor=abc"
}`

const forexTickerOverviewJSON = `{
	"status": "OK",
	"request_id": "forex_overview_001",
	"results": {
		"ticker": "C:EURUSD",
		"name": "Euro - United States Dollar",
		"market": "fx",
		"locale": "global",
		"active": true,
		"currency_symbol": "USD",
		"currency_name": "United States Dollar",
		"base_currency_symbol": "EUR",
		"base_currency_name": "Euro"
	}
}`

const forexExchangesJSON = `{
	"status": "OK",
	"request_id": "forex_exch_001",
	"count": 1,
	"results": [
		{
			"id": 48,
			"type": "exchange",
			"asset_class": "fx",
			"locale": "global",
			"name": "Forex Exchange",
			"acronym": "FX",
			"mic": "FXEX"
		}
	]
}`

// --- Aggregates Tests ---

// TestGetForexBars verifies that GetForexBars correctly parses the API
// response and returns the expected OHLC bar data for a forex ticker.
func TestGetForexBars(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/aggs/ticker/C:EURUSD/range/1/day/2025-01-06/2025-01-08": forexBarsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := ForexBarsParams{
		Multiplier: "1",
		Timespan:   "day",
		From:       "2025-01-06",
		To:         "2025-01-08",
		Adjusted:   "true",
		Sort:       "asc",
		Limit:      "2",
	}

	result, err := client.GetForexBars("C:EURUSD", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Ticker != "C:EURUSD" {
		t.Errorf("expected ticker C:EURUSD, got %s", result.Ticker)
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
	if bar.Open != 1.0850 {
		t.Errorf("expected open 1.0850, got %f", bar.Open)
	}

	if bar.High != 1.0870 {
		t.Errorf("expected high 1.0870, got %f", bar.High)
	}

	if bar.Low != 1.0840 {
		t.Errorf("expected low 1.0840, got %f", bar.Low)
	}

	if bar.Close != 1.0860 {
		t.Errorf("expected close 1.0860, got %f", bar.Close)
	}

	if bar.Volume != 50000 {
		t.Errorf("expected volume 50000, got %f", bar.Volume)
	}

	if bar.VWAP != 1.0855 {
		t.Errorf("expected VWAP 1.0855, got %f", bar.VWAP)
	}

	if bar.NumTrades != 1200 {
		t.Errorf("expected 1200 trades, got %d", bar.NumTrades)
	}
}

// TestGetForexBarsRequestPath verifies that GetForexBars constructs the
// correct URL path with ticker, multiplier, timespan, from, and to values.
func TestGetForexBarsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexBarsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	params := ForexBarsParams{
		Multiplier: "5",
		Timespan:   "minute",
		From:       "2025-01-06",
		To:         "2025-01-07",
	}

	client.GetForexBars("C:GBPUSD", params)

	expected := "/v2/aggs/ticker/C:GBPUSD/range/5/minute/2025-01-06/2025-01-07"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetForexBarsQueryParams verifies that GetForexBars sends the correct
// query parameters including adjusted, sort, and limit.
func TestGetForexBarsQueryParams(t *testing.T) {
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
		w.Write([]byte(forexBarsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	params := ForexBarsParams{
		Multiplier: "1",
		Timespan:   "day",
		From:       "2025-01-06",
		To:         "2025-01-08",
		Adjusted:   "false",
		Sort:       "desc",
		Limit:      "100",
	}

	client.GetForexBars("C:EURUSD", params)
}

// TestGetForexBarsSecondBar verifies that the second bar in the response
// is correctly parsed with its own distinct values.
func TestGetForexBarsSecondBar(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/aggs/ticker/C:EURUSD/range/1/day/2025-01-06/2025-01-08": forexBarsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := ForexBarsParams{
		Multiplier: "1",
		Timespan:   "day",
		From:       "2025-01-06",
		To:         "2025-01-08",
	}

	result, err := client.GetForexBars("C:EURUSD", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bar := result.Results[1]
	if bar.Open != 1.0860 {
		t.Errorf("expected open 1.0860, got %f", bar.Open)
	}

	if bar.Close != 1.0855 {
		t.Errorf("expected close 1.0855, got %f", bar.Close)
	}

	if bar.NumTrades != 1100 {
		t.Errorf("expected 1100 trades, got %d", bar.NumTrades)
	}
}

// TestGetForexBarsAPIError verifies that GetForexBars returns an error
// when the API responds with a non-200 status.
func TestGetForexBarsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Data not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	params := ForexBarsParams{
		Multiplier: "1",
		Timespan:   "day",
		From:       "2025-01-06",
		To:         "2025-01-08",
	}

	_, err := client.GetForexBars("C:INVALID", params)
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// --- Daily Market Summary Tests ---

// TestGetForexDailyMarketSummary verifies that GetForexDailyMarketSummary
// correctly parses the grouped daily response with multiple forex tickers.
func TestGetForexDailyMarketSummary(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/aggs/grouped/locale/global/market/fx/2025-01-06": forexMarketSummaryJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := ForexMarketSummaryParams{
		Adjusted: "true",
	}

	result, err := client.GetForexDailyMarketSummary("2025-01-06", params)
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
	if first.Ticker != "C:EURUSD" {
		t.Errorf("expected ticker C:EURUSD, got %s", first.Ticker)
	}

	if first.Open != 1.0850 {
		t.Errorf("expected open 1.0850, got %f", first.Open)
	}

	second := result.Results[1]
	if second.Ticker != "C:GBPUSD" {
		t.Errorf("expected ticker C:GBPUSD, got %s", second.Ticker)
	}
}

// TestGetForexDailyMarketSummaryRequestPath verifies that the correct
// API path is constructed for forex market summary requests.
func TestGetForexDailyMarketSummaryRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexMarketSummaryJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetForexDailyMarketSummary("2025-06-15", ForexMarketSummaryParams{})

	expected := "/v2/aggs/grouped/locale/global/market/fx/2025-06-15"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetForexDailyMarketSummaryAdjustedParam verifies that the adjusted
// query parameter is correctly sent to the API.
func TestGetForexDailyMarketSummaryAdjustedParam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("adjusted") != "false" {
			t.Errorf("expected adjusted=false, got %s", q.Get("adjusted"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexMarketSummaryJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetForexDailyMarketSummary("2025-01-06", ForexMarketSummaryParams{
		Adjusted: "false",
	})
}

// --- Previous Day Bar Tests ---

// TestGetForexPreviousDayBar verifies that GetForexPreviousDayBar correctly
// parses the API response for a forex ticker's previous day bar data.
func TestGetForexPreviousDayBar(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/aggs/ticker/C:EURUSD/prev": forexPrevDayJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetForexPreviousDayBar("C:EURUSD", "true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Ticker != "C:EURUSD" {
		t.Errorf("expected ticker C:EURUSD, got %s", result.Ticker)
	}

	if result.ResultsCount != 1 {
		t.Errorf("expected 1 result, got %d", result.ResultsCount)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 bar, got %d", len(result.Results))
	}

	bar := result.Results[0]
	if bar.Open != 1.0845 {
		t.Errorf("expected open 1.0845, got %f", bar.Open)
	}

	if bar.Close != 1.0850 {
		t.Errorf("expected close 1.0850, got %f", bar.Close)
	}

	if bar.Volume != 60000 {
		t.Errorf("expected volume 60000, got %f", bar.Volume)
	}
}

// TestGetForexPreviousDayBarRequestPath verifies that GetForexPreviousDayBar
// constructs the correct API path with the forex ticker.
func TestGetForexPreviousDayBarRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexPrevDayJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetForexPreviousDayBar("C:GBPUSD", "true")

	expected := "/v2/aggs/ticker/C:GBPUSD/prev"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// --- Currency Conversion Tests ---

// TestGetForexConversion verifies that GetForexConversion correctly parses
// the API response for currency conversion including the converted amount,
// initial amount, and last quote data.
func TestGetForexConversion(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/conversion/USD/EUR": forexConversionJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := ForexConversionParams{
		Amount:    "100",
		Precision: "2",
	}

	result, err := client.GetForexConversion("USD", "EUR", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.From != "USD" {
		t.Errorf("expected from USD, got %s", result.From)
	}

	if result.To != "EUR" {
		t.Errorf("expected to EUR, got %s", result.To)
	}

	if result.InitialAmount != 100 {
		t.Errorf("expected initialAmount 100, got %f", result.InitialAmount)
	}

	if result.Converted != 108.50 {
		t.Errorf("expected converted 108.50, got %f", result.Converted)
	}

	if result.Symbol != "USD/EUR" {
		t.Errorf("expected symbol USD/EUR, got %s", result.Symbol)
	}

	if result.Last.Ask != 1.0855 {
		t.Errorf("expected last ask 1.0855, got %f", result.Last.Ask)
	}

	if result.Last.Bid != 1.0850 {
		t.Errorf("expected last bid 1.0850, got %f", result.Last.Bid)
	}

	if result.Last.Exchange != 48 {
		t.Errorf("expected last exchange 48, got %d", result.Last.Exchange)
	}

	if result.Last.Timestamp != 1736139600000 {
		t.Errorf("expected last timestamp 1736139600000, got %d", result.Last.Timestamp)
	}
}

// TestGetForexConversionRequestPath verifies that GetForexConversion constructs
// the correct API path with the from and to currency codes.
func TestGetForexConversionRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexConversionJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetForexConversion("GBP", "JPY", ForexConversionParams{})

	expected := "/v1/conversion/GBP/JPY"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetForexConversionQueryParams verifies that the amount and precision
// query parameters are correctly sent to the API.
func TestGetForexConversionQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("amount") != "500" {
			t.Errorf("expected amount=500, got %s", q.Get("amount"))
		}
		if q.Get("precision") != "4" {
			t.Errorf("expected precision=4, got %s", q.Get("precision"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexConversionJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetForexConversion("USD", "EUR", ForexConversionParams{
		Amount:    "500",
		Precision: "4",
	})
}

// TestGetForexConversionAPIError verifies that GetForexConversion returns an
// error when the API responds with a non-200 status.
func TestGetForexConversionAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"ERROR","message":"Invalid currency pair."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetForexConversion("INVALID", "XXX", ForexConversionParams{})
	if err == nil {
		t.Fatal("expected error for 400 response, got nil")
	}
}

// --- Exchanges Tests ---

// TestGetForexExchanges verifies that GetForexExchanges correctly calls the
// shared exchanges endpoint with asset_class=fx and parses the response.
func TestGetForexExchanges(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/reference/exchanges" {
			t.Errorf("expected path /v3/reference/exchanges, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("asset_class") != "fx" {
			t.Errorf("expected asset_class=fx, got %s", r.URL.Query().Get("asset_class"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexExchangesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetForexExchanges()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Count != 1 {
		t.Errorf("expected count 1, got %d", result.Count)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 exchange, got %d", len(result.Results))
	}

	e := result.Results[0]
	if e.Name != "Forex Exchange" {
		t.Errorf("expected name Forex Exchange, got %s", e.Name)
	}

	if e.AssetClass != "fx" {
		t.Errorf("expected asset_class fx, got %s", e.AssetClass)
	}
}

// --- Market Status Tests ---

// TestGetForexMarketStatus verifies that GetForexMarketStatus correctly
// delegates to the shared GetMarketStatus method.
func TestGetForexMarketStatus(t *testing.T) {
	statusJSON := `{
		"afterHours": false,
		"currencies": {"crypto": "open", "fx": "open"},
		"earlyHours": false,
		"exchanges": {"nasdaq": "open", "nyse": "open", "otc": "open"},
		"indicesGroups": {"s_and_p": "open", "societe_generale": "open", "msci": "open", "ftse_russell": "open", "mstar": "open", "mstarc": "open", "cccy": "open", "cgi": "open", "nasdaq": "open", "dow_jones": "open"},
		"market": "open",
		"serverTime": "2025-01-06T12:00:00-05:00"
	}`

	server := mockServer(t, map[string]string{
		"/v1/marketstatus/now": statusJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetForexMarketStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Market != "open" {
		t.Errorf("expected market open, got %s", result.Market)
	}

	if result.Currencies.FX != "open" {
		t.Errorf("expected fx open, got %s", result.Currencies.FX)
	}
}

// --- Market Holidays Tests ---

// TestGetForexMarketHolidays verifies that GetForexMarketHolidays correctly
// delegates to the shared GetMarketHolidays method and parses the response.
func TestGetForexMarketHolidays(t *testing.T) {
	holidaysJSON := `[
		{
			"date": "2025-01-20",
			"exchange": "FOREX",
			"name": "Martin Luther King Jr. Day",
			"status": "closed"
		}
	]`

	server := mockServer(t, map[string]string{
		"/v1/marketstatus/upcoming": holidaysJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetForexMarketHolidays()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 holiday, got %d", len(result))
	}

	if result[0].Name != "Martin Luther King Jr. Day" {
		t.Errorf("expected MLK Day, got %s", result[0].Name)
	}

	if result[0].Status != "closed" {
		t.Errorf("expected status closed, got %s", result[0].Status)
	}
}

// --- Quotes Tests ---

// TestGetForexQuotes verifies that GetForexQuotes correctly parses the API
// response and returns the expected forex quote data.
func TestGetForexQuotes(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/quotes/C:EURUSD": forexQuotesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := ForexQuotesParams{
		Limit: "2",
		Sort:  "timestamp",
		Order: "asc",
	}

	result, err := client.GetForexQuotes("C:EURUSD", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 quotes, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.AskPrice != 1.0855 {
		t.Errorf("expected ask_price 1.0855, got %f", first.AskPrice)
	}

	if first.BidPrice != 1.0850 {
		t.Errorf("expected bid_price 1.0850, got %f", first.BidPrice)
	}

	if first.AskExchange != 48 {
		t.Errorf("expected ask_exchange 48, got %d", first.AskExchange)
	}

	if first.ParticipantTimestamp != 1736139600000 {
		t.Errorf("expected participant_timestamp 1736139600000, got %d", first.ParticipantTimestamp)
	}

	second := result.Results[1]
	if second.AskPrice != 1.0860 {
		t.Errorf("expected ask_price 1.0860, got %f", second.AskPrice)
	}

	if second.BidPrice != 1.0852 {
		t.Errorf("expected bid_price 1.0852, got %f", second.BidPrice)
	}
}

// TestGetForexQuotesRequestPath verifies that GetForexQuotes constructs
// the correct API path with the forex ticker.
func TestGetForexQuotesRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexQuotesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetForexQuotes("C:GBPUSD", ForexQuotesParams{})

	expected := "/v3/quotes/C:GBPUSD"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetForexQuotesQueryParams verifies that all query parameters are
// correctly sent to the forex quotes endpoint.
func TestGetForexQuotesQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("timestamp.gte") != "2025-01-06" {
			t.Errorf("expected timestamp.gte=2025-01-06, got %s", q.Get("timestamp.gte"))
		}
		if q.Get("timestamp.lte") != "2025-01-08" {
			t.Errorf("expected timestamp.lte=2025-01-08, got %s", q.Get("timestamp.lte"))
		}
		if q.Get("order") != "desc" {
			t.Errorf("expected order=desc, got %s", q.Get("order"))
		}
		if q.Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "timestamp" {
			t.Errorf("expected sort=timestamp, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexQuotesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetForexQuotes("C:EURUSD", ForexQuotesParams{
		TimestampGte: "2025-01-06",
		TimestampLte: "2025-01-08",
		Order:        "desc",
		Limit:        "50",
		Sort:         "timestamp",
	})
}

// --- Last Quote Tests ---

// TestGetForexLastQuote verifies that GetForexLastQuote correctly parses
// the API response for the most recent forex quote of a currency pair.
func TestGetForexLastQuote(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/last_quote/currencies/EUR/USD": forexLastQuoteJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetForexLastQuote("EUR", "USD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Symbol != "EUR/USD" {
		t.Errorf("expected symbol EUR/USD, got %s", result.Symbol)
	}

	if result.Last.Ask != 1.0855 {
		t.Errorf("expected ask 1.0855, got %f", result.Last.Ask)
	}

	if result.Last.Bid != 1.0850 {
		t.Errorf("expected bid 1.0850, got %f", result.Last.Bid)
	}

	if result.Last.Exchange != 48 {
		t.Errorf("expected exchange 48, got %d", result.Last.Exchange)
	}

	if result.Last.Timestamp != 1736139600000 {
		t.Errorf("expected timestamp 1736139600000, got %d", result.Last.Timestamp)
	}
}

// TestGetForexLastQuoteRequestPath verifies that GetForexLastQuote constructs
// the correct API path with the from and to currency codes.
func TestGetForexLastQuoteRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexLastQuoteJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetForexLastQuote("GBP", "JPY")

	expected := "/v1/last_quote/currencies/GBP/JPY"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetForexLastQuoteAPIError verifies that GetForexLastQuote returns an
// error when the API responds with a non-200 status.
func TestGetForexLastQuoteAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Currency pair not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetForexLastQuote("INVALID", "XXX")
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// --- Snapshot Tests ---

// TestGetForexSnapshotAll verifies that GetForexSnapshotAll correctly parses
// the API response for a full market forex snapshot.
func TestGetForexSnapshotAll(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/snapshot/locale/global/markets/forex/tickers": forexSnapshotAllJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := ForexSnapshotAllParams{
		Tickers: "C:EURUSD,C:GBPUSD",
	}

	result, err := client.GetForexSnapshotAll(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Count != 2 {
		t.Errorf("expected count 2, got %d", result.Count)
	}

	if len(result.Tickers) != 2 {
		t.Fatalf("expected 2 tickers, got %d", len(result.Tickers))
	}

	first := result.Tickers[0]
	if first.Ticker != "C:EURUSD" {
		t.Errorf("expected ticker C:EURUSD, got %s", first.Ticker)
	}

	if first.Day.Open != 1.0850 {
		t.Errorf("expected day open 1.0850, got %f", first.Day.Open)
	}

	if first.LastQuote.Ask != 1.0855 {
		t.Errorf("expected last quote ask 1.0855, got %f", first.LastQuote.Ask)
	}

	if first.PrevDay.Close != 1.0850 {
		t.Errorf("expected prev day close 1.0850, got %f", first.PrevDay.Close)
	}
}

// TestGetForexSnapshotAllTickersParam verifies that the tickers query parameter
// is correctly sent to the forex snapshot endpoint.
func TestGetForexSnapshotAllTickersParam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("tickers") != "C:EURUSD,C:GBPUSD" {
			t.Errorf("expected tickers=C:EURUSD,C:GBPUSD, got %s", r.URL.Query().Get("tickers"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexSnapshotAllJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetForexSnapshotAll(ForexSnapshotAllParams{
		Tickers: "C:EURUSD,C:GBPUSD",
	})
}

// TestGetForexSnapshotTicker verifies that GetForexSnapshotTicker correctly
// parses the API response for a single forex ticker snapshot.
func TestGetForexSnapshotTicker(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/snapshot/locale/global/markets/forex/tickers/C:EURUSD": forexSnapshotSingleJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetForexSnapshotTicker("C:EURUSD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	t2 := result.Ticker
	if t2.Ticker != "C:EURUSD" {
		t.Errorf("expected ticker C:EURUSD, got %s", t2.Ticker)
	}

	if t2.TodaysChange != 0.0010 {
		t.Errorf("expected todaysChange 0.0010, got %f", t2.TodaysChange)
	}

	if t2.Day.Open != 1.0850 {
		t.Errorf("expected day open 1.0850, got %f", t2.Day.Open)
	}

	if t2.LastQuote.Ask != 1.0855 {
		t.Errorf("expected last quote ask 1.0855, got %f", t2.LastQuote.Ask)
	}
}

// TestGetForexSnapshotTickerRequestPath verifies that GetForexSnapshotTicker
// constructs the correct API path with the forex ticker.
func TestGetForexSnapshotTickerRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexSnapshotSingleJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetForexSnapshotTicker("C:GBPUSD")

	expected := "/v2/snapshot/locale/global/markets/forex/tickers/C:GBPUSD"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// --- Gainers/Losers Tests ---

// TestGetForexGainersLosersGainers verifies that GetForexGainersLosers
// correctly parses the API response for the top forex gainers.
func TestGetForexGainersLosersGainers(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/snapshot/locale/global/markets/forex/gainers": forexGainersJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetForexGainersLosers("gainers")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Tickers) != 1 {
		t.Fatalf("expected 1 ticker, got %d", len(result.Tickers))
	}

	ticker := result.Tickers[0]
	if ticker.Ticker != "C:EURUSD" {
		t.Errorf("expected ticker C:EURUSD, got %s", ticker.Ticker)
	}

	if ticker.TodaysChangePct != 0.461 {
		t.Errorf("expected todaysChangePerc 0.461, got %f", ticker.TodaysChangePct)
	}
}

// TestGetForexGainersLosersRequestPath verifies that GetForexGainersLosers
// constructs the correct API path for both gainers and losers directions.
func TestGetForexGainersLosersRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexGainersJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	client.GetForexGainersLosers("losers")
	expected := "/v2/snapshot/locale/global/markets/forex/losers"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}

	client.GetForexGainersLosers("gainers")
	expected = "/v2/snapshot/locale/global/markets/forex/gainers"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// --- Unified Snapshot Tests ---

// TestGetForexUnifiedSnapshot verifies that GetForexUnifiedSnapshot correctly
// parses the API response from the unified snapshot endpoint.
func TestGetForexUnifiedSnapshot(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot": forexUnifiedSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetForexUnifiedSnapshot("C:EURUSD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	r := result.Results[0]
	if r.Ticker != "C:EURUSD" {
		t.Errorf("expected ticker C:EURUSD, got %s", r.Ticker)
	}

	if r.Day.Open != 1.0850 {
		t.Errorf("expected day open 1.0850, got %f", r.Day.Open)
	}

	if r.PrevDay.Close != 1.0850 {
		t.Errorf("expected prev day close 1.0850, got %f", r.PrevDay.Close)
	}
}

// TestGetForexUnifiedSnapshotQueryParams verifies that the ticker.any_of
// query parameter is correctly sent to the unified snapshot endpoint.
func TestGetForexUnifiedSnapshotQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("ticker.any_of") != "C:EURUSD,C:GBPUSD" {
			t.Errorf("expected ticker.any_of=C:EURUSD,C:GBPUSD, got %s", q.Get("ticker.any_of"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexUnifiedSnapshotJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetForexUnifiedSnapshot("C:EURUSD,C:GBPUSD")
}

// --- Technical Indicator Tests ---

// TestGetForexSMA verifies that GetForexSMA correctly parses the API response
// and returns the expected SMA indicator values for a forex ticker.
func TestGetForexSMA(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/sma/C:EURUSD": forexIndicatorJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := IndicatorParams{
		TimestampGTE: "2025-01-06",
		TimestampLTE: "2025-01-31",
		Timespan:     "day",
		Window:       "10",
		SeriesType:   "close",
		Limit:        "2",
	}

	result, err := client.GetForexSMA("C:EURUSD", params)
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
	if first.Value != 1.0852 {
		t.Errorf("expected value 1.0852, got %f", first.Value)
	}

	if first.Timestamp != 1736139600000 {
		t.Errorf("expected timestamp 1736139600000, got %d", first.Timestamp)
	}
}

// TestGetForexSMARequestPath verifies that GetForexSMA constructs the
// correct API path with the forex ticker.
func TestGetForexSMARequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexIndicatorJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetForexSMA("C:GBPUSD", IndicatorParams{})

	expected := "/v1/indicators/sma/C:GBPUSD"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetForexSMAQueryParams verifies that all indicator query parameters
// are correctly sent to the SMA endpoint.
func TestGetForexSMAQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("timestamp.gte") != "2025-01-06" {
			t.Errorf("expected timestamp.gte=2025-01-06, got %s", q.Get("timestamp.gte"))
		}
		if q.Get("timestamp.lte") != "2025-01-31" {
			t.Errorf("expected timestamp.lte=2025-01-31, got %s", q.Get("timestamp.lte"))
		}
		if q.Get("timespan") != "day" {
			t.Errorf("expected timespan=day, got %s", q.Get("timespan"))
		}
		if q.Get("adjusted") != "true" {
			t.Errorf("expected adjusted=true, got %s", q.Get("adjusted"))
		}
		if q.Get("window") != "20" {
			t.Errorf("expected window=20, got %s", q.Get("window"))
		}
		if q.Get("series_type") != "close" {
			t.Errorf("expected series_type=close, got %s", q.Get("series_type"))
		}
		if q.Get("order") != "desc" {
			t.Errorf("expected order=desc, got %s", q.Get("order"))
		}
		if q.Get("limit") != "10" {
			t.Errorf("expected limit=10, got %s", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexIndicatorJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetForexSMA("C:EURUSD", IndicatorParams{
		TimestampGTE: "2025-01-06",
		TimestampLTE: "2025-01-31",
		Timespan:     "day",
		Adjusted:     "true",
		Window:       "20",
		SeriesType:   "close",
		Order:        "desc",
		Limit:        "10",
	})
}

// TestGetForexEMA verifies that GetForexEMA correctly parses the API response
// and returns the expected EMA indicator values for a forex ticker.
func TestGetForexEMA(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/ema/C:EURUSD": forexIndicatorJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := IndicatorParams{
		TimestampGTE: "2025-01-06",
		TimestampLTE: "2025-01-31",
		Timespan:     "day",
		Window:       "10",
	}

	result, err := client.GetForexEMA("C:EURUSD", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Results.Values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(result.Results.Values))
	}

	if result.Results.Values[0].Value != 1.0852 {
		t.Errorf("expected value 1.0852, got %f", result.Results.Values[0].Value)
	}
}

// TestGetForexEMARequestPath verifies that GetForexEMA constructs the
// correct API path with the forex ticker.
func TestGetForexEMARequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexIndicatorJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetForexEMA("C:EURUSD", IndicatorParams{})

	expected := "/v1/indicators/ema/C:EURUSD"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetForexRSI verifies that GetForexRSI correctly parses the API response
// and returns the expected RSI indicator values for a forex ticker.
func TestGetForexRSI(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/rsi/C:EURUSD": forexIndicatorJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := IndicatorParams{
		TimestampGTE: "2025-01-06",
		TimestampLTE: "2025-01-31",
		Timespan:     "day",
		Window:       "14",
	}

	result, err := client.GetForexRSI("C:EURUSD", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Results.Values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(result.Results.Values))
	}
}

// TestGetForexRSIRequestPath verifies that GetForexRSI constructs the
// correct API path with the forex ticker.
func TestGetForexRSIRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexIndicatorJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetForexRSI("C:EURUSD", IndicatorParams{})

	expected := "/v1/indicators/rsi/C:EURUSD"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetForexMACD verifies that GetForexMACD correctly parses the API
// response and returns the expected MACD indicator values for a forex ticker.
func TestGetForexMACD(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/macd/C:EURUSD": forexMACDJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := MACDParams{
		TimestampGTE: "2025-01-06",
		TimestampLTE: "2025-01-31",
		Timespan:     "day",
		ShortWindow:  "12",
		LongWindow:   "26",
		SignalWindow: "9",
	}

	result, err := client.GetForexMACD("C:EURUSD", params)
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
	if first.Value != 0.0012 {
		t.Errorf("expected MACD value 0.0012, got %f", first.Value)
	}

	if first.Signal != 0.0008 {
		t.Errorf("expected signal 0.0008, got %f", first.Signal)
	}

	if first.Histogram != 0.0004 {
		t.Errorf("expected histogram 0.0004, got %f", first.Histogram)
	}
}

// TestGetForexMACDRequestPath verifies that GetForexMACD constructs the
// correct API path with the forex ticker.
func TestGetForexMACDRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexMACDJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetForexMACD("C:GBPUSD", MACDParams{})

	expected := "/v1/indicators/macd/C:GBPUSD"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetForexMACDQueryParams verifies that all MACD-specific query parameters
// are correctly sent to the MACD endpoint.
func TestGetForexMACDQueryParams(t *testing.T) {
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
		w.Write([]byte(forexMACDJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetForexMACD("C:EURUSD", MACDParams{
		ShortWindow:  "12",
		LongWindow:   "26",
		SignalWindow: "9",
	})
}

// --- Tickers Tests ---

// TestGetForexTickers verifies that GetForexTickers correctly parses the
// reference tickers response for forex market tickers.
func TestGetForexTickers(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/reference/tickers": forexTickersJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := ForexTickerParams{
		Search: "EUR",
		Limit:  "2",
	}

	result, err := client.GetForexTickers(params)
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
	if first.Ticker != "C:EURUSD" {
		t.Errorf("expected ticker C:EURUSD, got %s", first.Ticker)
	}

	if first.Name != "Euro - United States Dollar" {
		t.Errorf("expected name Euro - United States Dollar, got %s", first.Name)
	}

	if first.Market != "fx" {
		t.Errorf("expected market fx, got %s", first.Market)
	}

	if !first.Active {
		t.Error("expected active to be true")
	}
}

// TestGetForexTickersMarketParam verifies that the market=fx parameter
// is always sent when querying forex tickers.
func TestGetForexTickersMarketParam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("market") != "fx" {
			t.Errorf("expected market=fx, got %s", q.Get("market"))
		}
		if q.Get("search") != "GBP" {
			t.Errorf("expected search=GBP, got %s", q.Get("search"))
		}
		if q.Get("active") != "true" {
			t.Errorf("expected active=true, got %s", q.Get("active"))
		}
		if q.Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexTickersJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetForexTickers(ForexTickerParams{
		Search: "GBP",
		Active: "true",
		Limit:  "50",
	})
}

// TestGetForexTickersEmptyResults verifies that GetForexTickers handles an
// empty results array without error.
func TestGetForexTickersEmptyResults(t *testing.T) {
	emptyJSON := `{"results":[],"status":"OK","request_id":"abc","count":0}`
	server := mockServer(t, map[string]string{
		"/v3/reference/tickers": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetForexTickers(ForexTickerParams{Search: "zzzznotreal"})
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

// --- Ticker Overview Tests ---

// TestGetForexTickerOverview verifies that GetForexTickerOverview correctly
// parses the API response for detailed forex ticker reference data.
func TestGetForexTickerOverview(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/reference/tickers/C:EURUSD": forexTickerOverviewJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetForexTickerOverview("C:EURUSD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	r := result.Results
	if r.Ticker != "C:EURUSD" {
		t.Errorf("expected ticker C:EURUSD, got %s", r.Ticker)
	}

	if r.Name != "Euro - United States Dollar" {
		t.Errorf("expected name Euro - United States Dollar, got %s", r.Name)
	}

	if r.Market != "fx" {
		t.Errorf("expected market fx, got %s", r.Market)
	}

	if r.Locale != "global" {
		t.Errorf("expected locale global, got %s", r.Locale)
	}

	if !r.Active {
		t.Error("expected active to be true")
	}

	if r.CurrencyName != "United States Dollar" {
		t.Errorf("expected currency_name United States Dollar, got %s", r.CurrencyName)
	}

	if r.BaseCurrencyName != "Euro" {
		t.Errorf("expected base_currency_name Euro, got %s", r.BaseCurrencyName)
	}

	if r.CurrencySymbol != "USD" {
		t.Errorf("expected currency_symbol USD, got %s", r.CurrencySymbol)
	}

	if r.BaseCurrencySymbol != "EUR" {
		t.Errorf("expected base_currency_symbol EUR, got %s", r.BaseCurrencySymbol)
	}
}

// TestGetForexTickerOverviewRequestPath verifies that GetForexTickerOverview
// constructs the correct API path with the forex ticker.
func TestGetForexTickerOverviewRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(forexTickerOverviewJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetForexTickerOverview("C:GBPJPY")

	expected := "/v3/reference/tickers/C:GBPJPY"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetForexTickerOverviewAPIError verifies that GetForexTickerOverview
// returns an error when the API responds with a non-200 status.
func TestGetForexTickerOverviewAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Ticker not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetForexTickerOverview("C:INVALID")
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}
