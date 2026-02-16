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

const optionsBarsJSON = `{
	"ticker": "O:AAPL250221C00230000",
	"queryCount": 3,
	"resultsCount": 3,
	"adjusted": true,
	"count": 3,
	"results": [
		{
			"v": 9162,
			"vw": 3.1012,
			"o": 3.75,
			"c": 2.53,
			"h": 4.18,
			"l": 2.39,
			"t": 1739163600000,
			"n": 1716
		},
		{
			"v": 11082,
			"vw": 4.9512,
			"o": 2.80,
			"c": 5.25,
			"h": 7.16,
			"l": 2.80,
			"t": 1739250000000,
			"n": 2252
		},
		{
			"v": 4423,
			"vw": 6.2639,
			"o": 4.55,
			"c": 8.22,
			"h": 8.22,
			"l": 4.15,
			"t": 1739336400000,
			"n": 1064
		}
	],
	"status": "OK",
	"request_id": "f3991f6b445aaeaae7b2e3c0bc6056e8"
}`

const optionsDailyTickerSummaryJSON = `{
	"status": "OK",
	"from": "2025-02-10",
	"symbol": "O:AAPL250221C00230000",
	"open": 3.75,
	"high": 4.18,
	"low": 2.39,
	"close": 2.53,
	"volume": 9162,
	"afterHours": 2.53,
	"preMarket": 3.75
}`

const optionsPreviousDayBarJSON = `{
	"ticker": "O:AAPL250221C00230000",
	"queryCount": 1,
	"resultsCount": 1,
	"adjusted": true,
	"count": 1,
	"results": [
		{
			"T": "O:AAPL250221C00230000",
			"v": 3791,
			"vw": 16.3893,
			"o": 16.50,
			"c": 15.68,
			"h": 18.61,
			"l": 15.25,
			"t": 1740171600000,
			"n": 351
		}
	],
	"status": "OK",
	"request_id": "2cc2c860c9f4d613f312725bdeaafb98"
}`

// TestGetOptionsBars verifies that GetOptionsBars correctly parses the API
// response and returns the expected OHLC bar data for an options contract.
func TestGetOptionsBars(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/aggs/ticker/O:AAPL250221C00230000/range/1/day/2025-02-10/2025-02-14": optionsBarsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := OptionsBarsParams{
		Multiplier: "1",
		Timespan:   "day",
		From:       "2025-02-10",
		To:         "2025-02-14",
		Adjusted:   "true",
		Sort:       "asc",
		Limit:      "5000",
	}

	result, err := client.GetOptionsBars("O:AAPL250221C00230000", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Ticker != "O:AAPL250221C00230000" {
		t.Errorf("expected ticker O:AAPL250221C00230000, got %s", result.Ticker)
	}

	if result.Adjusted != true {
		t.Errorf("expected adjusted true, got %v", result.Adjusted)
	}

	if result.QueryCount != 3 {
		t.Errorf("expected queryCount 3, got %d", result.QueryCount)
	}

	if result.ResultsCount != 3 {
		t.Errorf("expected resultsCount 3, got %d", result.ResultsCount)
	}

	if len(result.Results) != 3 {
		t.Fatalf("expected 3 bars, got %d", len(result.Results))
	}

	bar := result.Results[0]
	if bar.Open != 3.75 {
		t.Errorf("expected open 3.75, got %f", bar.Open)
	}

	if bar.High != 4.18 {
		t.Errorf("expected high 4.18, got %f", bar.High)
	}

	if bar.Low != 2.39 {
		t.Errorf("expected low 2.39, got %f", bar.Low)
	}

	if bar.Close != 2.53 {
		t.Errorf("expected close 2.53, got %f", bar.Close)
	}

	if bar.Volume != 9162 {
		t.Errorf("expected volume 9162, got %f", bar.Volume)
	}

	if bar.VWAP != 3.1012 {
		t.Errorf("expected vwap 3.1012, got %f", bar.VWAP)
	}

	if bar.Timestamp != 1739163600000 {
		t.Errorf("expected timestamp 1739163600000, got %d", bar.Timestamp)
	}

	if bar.NumTrades != 1716 {
		t.Errorf("expected numTrades 1716, got %d", bar.NumTrades)
	}
}

// TestGetOptionsBarsRequestPath verifies that GetOptionsBars constructs the
// correct URL path with the options ticker, multiplier, timespan, from, and
// to values embedded in the path.
func TestGetOptionsBarsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsBarsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	params := OptionsBarsParams{
		Multiplier: "5",
		Timespan:   "minute",
		From:       "2025-02-10",
		To:         "2025-02-11",
	}

	client.GetOptionsBars("O:SPY250221P00500000", params)

	expected := "/v2/aggs/ticker/O:SPY250221P00500000/range/5/minute/2025-02-10/2025-02-11"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetOptionsBarsQueryParams verifies that GetOptionsBars sends the
// correct query parameters including adjusted, sort, and limit.
func TestGetOptionsBarsQueryParams(t *testing.T) {
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
		w.Write([]byte(optionsBarsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	params := OptionsBarsParams{
		Multiplier: "1",
		Timespan:   "day",
		From:       "2025-02-10",
		To:         "2025-02-14",
		Adjusted:   "false",
		Sort:       "desc",
		Limit:      "100",
	}

	client.GetOptionsBars("O:AAPL250221C00230000", params)
}

// TestGetOptionsBarsSecondBar verifies that the second bar in the response
// is correctly parsed with its own distinct values for all fields.
func TestGetOptionsBarsSecondBar(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/aggs/ticker/O:AAPL250221C00230000/range/1/day/2025-02-10/2025-02-14": optionsBarsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := OptionsBarsParams{
		Multiplier: "1",
		Timespan:   "day",
		From:       "2025-02-10",
		To:         "2025-02-14",
	}

	result, err := client.GetOptionsBars("O:AAPL250221C00230000", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bar := result.Results[1]
	if bar.Open != 2.80 {
		t.Errorf("expected open 2.80, got %f", bar.Open)
	}

	if bar.Close != 5.25 {
		t.Errorf("expected close 5.25, got %f", bar.Close)
	}

	if bar.Volume != 11082 {
		t.Errorf("expected volume 11082, got %f", bar.Volume)
	}

	if bar.VWAP != 4.9512 {
		t.Errorf("expected vwap 4.9512, got %f", bar.VWAP)
	}

	if bar.Timestamp != 1739250000000 {
		t.Errorf("expected timestamp 1739250000000, got %d", bar.Timestamp)
	}

	if bar.NumTrades != 2252 {
		t.Errorf("expected numTrades 2252, got %d", bar.NumTrades)
	}
}

// TestGetOptionsBarsThirdBar verifies that the third bar in the response
// is correctly parsed with its own distinct values for all fields.
func TestGetOptionsBarsThirdBar(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/aggs/ticker/O:AAPL250221C00230000/range/1/day/2025-02-10/2025-02-14": optionsBarsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := OptionsBarsParams{
		Multiplier: "1",
		Timespan:   "day",
		From:       "2025-02-10",
		To:         "2025-02-14",
	}

	result, err := client.GetOptionsBars("O:AAPL250221C00230000", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bar := result.Results[2]
	if bar.Open != 4.55 {
		t.Errorf("expected open 4.55, got %f", bar.Open)
	}

	if bar.Close != 8.22 {
		t.Errorf("expected close 8.22, got %f", bar.Close)
	}

	if bar.High != 8.22 {
		t.Errorf("expected high 8.22, got %f", bar.High)
	}

	if bar.Low != 4.15 {
		t.Errorf("expected low 4.15, got %f", bar.Low)
	}

	if bar.NumTrades != 1064 {
		t.Errorf("expected numTrades 1064, got %d", bar.NumTrades)
	}
}

// TestGetOptionsBarsAPIError verifies that GetOptionsBars returns an error
// when the API responds with a non-200 status code.
func TestGetOptionsBarsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Data not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	params := OptionsBarsParams{
		Multiplier: "1",
		Timespan:   "day",
		From:       "2025-02-10",
		To:         "2025-02-14",
	}

	_, err := client.GetOptionsBars("O:INVALID", params)
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// TestGetOptionsDailyTickerSummary verifies that GetOptionsDailyTickerSummary
// correctly parses the API response and returns the expected daily OHLC data
// for an options contract.
func TestGetOptionsDailyTickerSummary(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/open-close/O:AAPL250221C00230000/2025-02-10": optionsDailyTickerSummaryJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsDailyTickerSummary("O:AAPL250221C00230000", "2025-02-10", "true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Symbol != "O:AAPL250221C00230000" {
		t.Errorf("expected symbol O:AAPL250221C00230000, got %s", result.Symbol)
	}

	if result.From != "2025-02-10" {
		t.Errorf("expected from 2025-02-10, got %s", result.From)
	}

	if result.Open != 3.75 {
		t.Errorf("expected open 3.75, got %f", result.Open)
	}

	if result.High != 4.18 {
		t.Errorf("expected high 4.18, got %f", result.High)
	}

	if result.Low != 2.39 {
		t.Errorf("expected low 2.39, got %f", result.Low)
	}

	if result.Close != 2.53 {
		t.Errorf("expected close 2.53, got %f", result.Close)
	}

	if result.Volume != 9162 {
		t.Errorf("expected volume 9162, got %f", result.Volume)
	}

	if result.AfterHours != 2.53 {
		t.Errorf("expected afterHours 2.53, got %f", result.AfterHours)
	}

	if result.PreMarket != 3.75 {
		t.Errorf("expected preMarket 3.75, got %f", result.PreMarket)
	}
}

// TestGetOptionsDailyTickerSummaryRequestPath verifies that
// GetOptionsDailyTickerSummary constructs the correct API path with the
// options ticker and date.
func TestGetOptionsDailyTickerSummaryRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsDailyTickerSummaryJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsDailyTickerSummary("O:SPY250221P00500000", "2025-03-15", "true")

	expected := "/v1/open-close/O:SPY250221P00500000/2025-03-15"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetOptionsDailyTickerSummaryQueryParams verifies that
// GetOptionsDailyTickerSummary sends the adjusted query parameter correctly.
func TestGetOptionsDailyTickerSummaryQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("adjusted") != "false" {
			t.Errorf("expected adjusted=false, got %s", q.Get("adjusted"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsDailyTickerSummaryJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsDailyTickerSummary("O:AAPL250221C00230000", "2025-02-10", "false")
}

// TestGetOptionsDailyTickerSummaryAPIError verifies that
// GetOptionsDailyTickerSummary returns an error when the API responds
// with a non-200 status code.
func TestGetOptionsDailyTickerSummaryAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Data not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetOptionsDailyTickerSummary("O:INVALID", "2025-02-10", "true")
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// TestGetOptionsPreviousDayBar verifies that GetOptionsPreviousDayBar
// correctly parses the API response and returns the expected previous
// day OHLC data for an options contract.
func TestGetOptionsPreviousDayBar(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/aggs/ticker/O:AAPL250221C00230000/prev": optionsPreviousDayBarJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsPreviousDayBar("O:AAPL250221C00230000", "true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Ticker != "O:AAPL250221C00230000" {
		t.Errorf("expected ticker O:AAPL250221C00230000, got %s", result.Ticker)
	}

	if result.Adjusted != true {
		t.Errorf("expected adjusted true, got %v", result.Adjusted)
	}

	if result.QueryCount != 1 {
		t.Errorf("expected queryCount 1, got %d", result.QueryCount)
	}

	if result.ResultsCount != 1 {
		t.Errorf("expected resultsCount 1, got %d", result.ResultsCount)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 bar, got %d", len(result.Results))
	}

	bar := result.Results[0]
	if bar.Ticker != "O:AAPL250221C00230000" {
		t.Errorf("expected bar ticker O:AAPL250221C00230000, got %s", bar.Ticker)
	}

	if bar.Open != 16.50 {
		t.Errorf("expected open 16.50, got %f", bar.Open)
	}

	if bar.High != 18.61 {
		t.Errorf("expected high 18.61, got %f", bar.High)
	}

	if bar.Low != 15.25 {
		t.Errorf("expected low 15.25, got %f", bar.Low)
	}

	if bar.Close != 15.68 {
		t.Errorf("expected close 15.68, got %f", bar.Close)
	}

	if bar.Volume != 3791 {
		t.Errorf("expected volume 3791, got %f", bar.Volume)
	}

	if bar.VWAP != 16.3893 {
		t.Errorf("expected vwap 16.3893, got %f", bar.VWAP)
	}

	if bar.Timestamp != 1740171600000 {
		t.Errorf("expected timestamp 1740171600000, got %d", bar.Timestamp)
	}

	if bar.NumTrades != 351 {
		t.Errorf("expected numTrades 351, got %d", bar.NumTrades)
	}
}

// TestGetOptionsPreviousDayBarRequestPath verifies that
// GetOptionsPreviousDayBar constructs the correct API path with the
// options contract ticker symbol.
func TestGetOptionsPreviousDayBarRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsPreviousDayBarJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsPreviousDayBar("O:SPY250221P00500000", "true")

	expected := "/v2/aggs/ticker/O:SPY250221P00500000/prev"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetOptionsPreviousDayBarQueryParams verifies that
// GetOptionsPreviousDayBar sends the adjusted query parameter correctly.
func TestGetOptionsPreviousDayBarQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("adjusted") != "false" {
			t.Errorf("expected adjusted=false, got %s", q.Get("adjusted"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsPreviousDayBarJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsPreviousDayBar("O:AAPL250221C00230000", "false")
}

// TestGetOptionsPreviousDayBarAPIError verifies that
// GetOptionsPreviousDayBar returns an error when the API responds
// with a non-200 status code.
func TestGetOptionsPreviousDayBarAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"ERROR","message":"Forbidden."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetOptionsPreviousDayBar("O:INVALID", "true")
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}
