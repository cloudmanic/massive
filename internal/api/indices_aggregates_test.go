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

const indicesBarsJSON = `{
	"ticker": "I:SPX",
	"queryCount": 3,
	"resultsCount": 3,
	"count": 3,
	"results": [
		{
			"o": 5982.81,
			"c": 5975.38,
			"h": 6021.04,
			"l": 5960.01,
			"t": 1736143200000
		},
		{
			"o": 5993.26,
			"c": 5909.03,
			"h": 6000.68,
			"l": 5890.68,
			"t": 1736229600000
		},
		{
			"o": 5910.66,
			"c": 5918.25,
			"h": 5927.89,
			"l": 5874.78,
			"t": 1736316000000
		}
	],
	"status": "OK",
	"request_id": "68f5c5b03a697c32b7d30112ba933379"
}`

const indicesDailyTickerSummaryJSON = `{
	"status": "OK",
	"from": "2025-01-06",
	"symbol": "I:SPX",
	"open": 5982.81,
	"high": 6021.04,
	"low": 5960.01,
	"close": 5975.38,
	"afterHours": 5975.38,
	"preMarket": 5982.81
}`

const indicesPreviousDayBarJSON = `{
	"ticker": "I:SPX",
	"queryCount": 1,
	"resultsCount": 1,
	"count": 1,
	"results": [
		{
			"T": "I:SPX",
			"o": 6834.27,
			"c": 6836.17,
			"h": 6881.96,
			"l": 6794.55,
			"t": 1771016400000
		}
	],
	"status": "OK",
	"request_id": "695f18c4181906dc6e5f8487b4a783ec"
}`

// TestGetIndicesBars verifies that GetIndicesBars correctly parses the API
// response and returns the expected OHLC bar data for I:SPX.
func TestGetIndicesBars(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/aggs/ticker/I:SPX/range/1/day/2025-01-06/2025-01-08": indicesBarsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := IndicesBarsParams{
		Multiplier: "1",
		Timespan:   "day",
		From:       "2025-01-06",
		To:         "2025-01-08",
		Sort:       "asc",
		Limit:      "5000",
	}

	result, err := client.GetIndicesBars("I:SPX", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Ticker != "I:SPX" {
		t.Errorf("expected ticker I:SPX, got %s", result.Ticker)
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
	if bar.Open != 5982.81 {
		t.Errorf("expected open 5982.81, got %f", bar.Open)
	}

	if bar.High != 6021.04 {
		t.Errorf("expected high 6021.04, got %f", bar.High)
	}

	if bar.Low != 5960.01 {
		t.Errorf("expected low 5960.01, got %f", bar.Low)
	}

	if bar.Close != 5975.38 {
		t.Errorf("expected close 5975.38, got %f", bar.Close)
	}

	if bar.Timestamp != 1736143200000 {
		t.Errorf("expected timestamp 1736143200000, got %d", bar.Timestamp)
	}
}

// TestGetIndicesBarsRequestPath verifies that GetIndicesBars constructs the
// correct URL path with ticker, multiplier, timespan, from, and to values.
func TestGetIndicesBarsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesBarsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	params := IndicesBarsParams{
		Multiplier: "5",
		Timespan:   "minute",
		From:       "2025-01-06",
		To:         "2025-01-07",
	}

	client.GetIndicesBars("I:NDX", params)

	expected := "/v2/aggs/ticker/I:NDX/range/5/minute/2025-01-06/2025-01-07"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetIndicesBarsQueryParams verifies that GetIndicesBars sends the
// correct query parameters including sort and limit.
func TestGetIndicesBarsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("sort") != "desc" {
			t.Errorf("expected sort=desc, got %s", q.Get("sort"))
		}
		if q.Get("limit") != "100" {
			t.Errorf("expected limit=100, got %s", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesBarsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	params := IndicesBarsParams{
		Multiplier: "1",
		Timespan:   "day",
		From:       "2025-01-06",
		To:         "2025-01-08",
		Sort:       "desc",
		Limit:      "100",
	}

	client.GetIndicesBars("I:SPX", params)
}

// TestGetIndicesBarsSecondBar verifies that the second bar in the response
// is correctly parsed with its own distinct values.
func TestGetIndicesBarsSecondBar(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/aggs/ticker/I:SPX/range/1/day/2025-01-06/2025-01-08": indicesBarsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := IndicesBarsParams{
		Multiplier: "1",
		Timespan:   "day",
		From:       "2025-01-06",
		To:         "2025-01-08",
	}

	result, err := client.GetIndicesBars("I:SPX", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bar := result.Results[1]
	if bar.Open != 5993.26 {
		t.Errorf("expected open 5993.26, got %f", bar.Open)
	}

	if bar.Close != 5909.03 {
		t.Errorf("expected close 5909.03, got %f", bar.Close)
	}

	if bar.Timestamp != 1736229600000 {
		t.Errorf("expected timestamp 1736229600000, got %d", bar.Timestamp)
	}
}

// TestGetIndicesBarsThirdBar verifies that the third bar in the response
// is correctly parsed with its own distinct values.
func TestGetIndicesBarsThirdBar(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/aggs/ticker/I:SPX/range/1/day/2025-01-06/2025-01-08": indicesBarsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := IndicesBarsParams{
		Multiplier: "1",
		Timespan:   "day",
		From:       "2025-01-06",
		To:         "2025-01-08",
	}

	result, err := client.GetIndicesBars("I:SPX", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bar := result.Results[2]
	if bar.Open != 5910.66 {
		t.Errorf("expected open 5910.66, got %f", bar.Open)
	}

	if bar.Close != 5918.25 {
		t.Errorf("expected close 5918.25, got %f", bar.Close)
	}

	if bar.High != 5927.89 {
		t.Errorf("expected high 5927.89, got %f", bar.High)
	}

	if bar.Low != 5874.78 {
		t.Errorf("expected low 5874.78, got %f", bar.Low)
	}
}

// TestGetIndicesBarsAPIError verifies that GetIndicesBars returns an error
// when the API responds with a non-200 status code.
func TestGetIndicesBarsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Data not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	params := IndicesBarsParams{
		Multiplier: "1",
		Timespan:   "day",
		From:       "2025-01-06",
		To:         "2025-01-08",
	}

	_, err := client.GetIndicesBars("I:INVALID", params)
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// TestGetIndicesDailyTickerSummary verifies that GetIndicesDailyTickerSummary
// correctly parses the API response and returns the expected daily OHLC data
// for an index ticker.
func TestGetIndicesDailyTickerSummary(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/open-close/I:SPX/2025-01-06": indicesDailyTickerSummaryJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesDailyTickerSummary("I:SPX", "2025-01-06")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Symbol != "I:SPX" {
		t.Errorf("expected symbol I:SPX, got %s", result.Symbol)
	}

	if result.From != "2025-01-06" {
		t.Errorf("expected from 2025-01-06, got %s", result.From)
	}

	if result.Open != 5982.81 {
		t.Errorf("expected open 5982.81, got %f", result.Open)
	}

	if result.High != 6021.04 {
		t.Errorf("expected high 6021.04, got %f", result.High)
	}

	if result.Low != 5960.01 {
		t.Errorf("expected low 5960.01, got %f", result.Low)
	}

	if result.Close != 5975.38 {
		t.Errorf("expected close 5975.38, got %f", result.Close)
	}

	if result.AfterHours != 5975.38 {
		t.Errorf("expected afterHours 5975.38, got %f", result.AfterHours)
	}

	if result.PreMarket != 5982.81 {
		t.Errorf("expected preMarket 5982.81, got %f", result.PreMarket)
	}
}

// TestGetIndicesDailyTickerSummaryRequestPath verifies that
// GetIndicesDailyTickerSummary constructs the correct API path with
// the ticker and date.
func TestGetIndicesDailyTickerSummaryRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesDailyTickerSummaryJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesDailyTickerSummary("I:NDX", "2025-03-15")

	expected := "/v1/open-close/I:NDX/2025-03-15"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetIndicesDailyTickerSummaryAPIError verifies that
// GetIndicesDailyTickerSummary returns an error when the API responds
// with a non-200 status code.
func TestGetIndicesDailyTickerSummaryAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Data not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetIndicesDailyTickerSummary("I:INVALID", "2025-01-06")
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// TestGetIndicesPreviousDayBar verifies that GetIndicesPreviousDayBar
// correctly parses the API response and returns the expected previous
// day OHLC data for an index ticker.
func TestGetIndicesPreviousDayBar(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/aggs/ticker/I:SPX/prev": indicesPreviousDayBarJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesPreviousDayBar("I:SPX")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Ticker != "I:SPX" {
		t.Errorf("expected ticker I:SPX, got %s", result.Ticker)
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
	if bar.Ticker != "I:SPX" {
		t.Errorf("expected bar ticker I:SPX, got %s", bar.Ticker)
	}

	if bar.Open != 6834.27 {
		t.Errorf("expected open 6834.27, got %f", bar.Open)
	}

	if bar.High != 6881.96 {
		t.Errorf("expected high 6881.96, got %f", bar.High)
	}

	if bar.Low != 6794.55 {
		t.Errorf("expected low 6794.55, got %f", bar.Low)
	}

	if bar.Close != 6836.17 {
		t.Errorf("expected close 6836.17, got %f", bar.Close)
	}

	if bar.Timestamp != 1771016400000 {
		t.Errorf("expected timestamp 1771016400000, got %d", bar.Timestamp)
	}
}

// TestGetIndicesPreviousDayBarRequestPath verifies that
// GetIndicesPreviousDayBar constructs the correct API path with the
// ticker symbol.
func TestGetIndicesPreviousDayBarRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesPreviousDayBarJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesPreviousDayBar("I:NDX")

	expected := "/v2/aggs/ticker/I:NDX/prev"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetIndicesPreviousDayBarAPIError verifies that
// GetIndicesPreviousDayBar returns an error when the API responds
// with a non-200 status code.
func TestGetIndicesPreviousDayBarAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"ERROR","message":"Forbidden."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetIndicesPreviousDayBar("I:INVALID")
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}
