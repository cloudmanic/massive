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

const indicesSnapshotSingleJSON = `{
	"status": "OK",
	"request_id": "idx123",
	"results": [
		{
			"ticker": "I:SPX",
			"name": "Standard & Poor's 500",
			"value": 6836.17,
			"type": "indices",
			"timeframe": "REAL-TIME",
			"market_status": "closed",
			"last_updated": 1771016540019000000,
			"session": {
				"change": -12.50,
				"change_percent": -0.1828,
				"close": 6836.17,
				"high": 6881.96,
				"low": 6794.55,
				"open": 6834.27,
				"previous_close": 6848.67
			}
		}
	]
}`

const indicesSnapshotMultipleJSON = `{
	"status": "OK",
	"request_id": "idx456",
	"results": [
		{
			"ticker": "I:SPX",
			"name": "Standard & Poor's 500",
			"value": 6836.17,
			"type": "indices",
			"timeframe": "REAL-TIME",
			"market_status": "closed",
			"last_updated": 1771016540019000000,
			"session": {
				"change": -12.50,
				"change_percent": -0.1828,
				"close": 6836.17,
				"high": 6881.96,
				"low": 6794.55,
				"open": 6834.27,
				"previous_close": 6848.67
			}
		},
		{
			"ticker": "I:DJI",
			"name": "Dow Jones Industrial Average",
			"value": 44546.08,
			"type": "indices",
			"timeframe": "REAL-TIME",
			"market_status": "closed",
			"last_updated": 1771016540019000000,
			"session": {
				"change": 342.30,
				"change_percent": 0.7742,
				"close": 44546.08,
				"high": 44600.50,
				"low": 44100.20,
				"open": 44200.00,
				"previous_close": 44203.78
			}
		}
	]
}`

const indicesSnapshotPaginatedJSON = `{
	"status": "OK",
	"request_id": "idx789",
	"next_url": "https://api.massive.com/v3/snapshot/indices?cursor=abc123",
	"results": [
		{
			"ticker": "I:COMP",
			"name": "NASDAQ Composite",
			"value": 19923.45,
			"type": "indices",
			"timeframe": "DELAYED",
			"market_status": "open",
			"last_updated": 1771016540019000000,
			"session": {
				"change": 150.20,
				"change_percent": 0.7598,
				"close": 19923.45,
				"high": 19950.00,
				"low": 19750.00,
				"open": 19800.00,
				"previous_close": 19773.25
			}
		}
	]
}`

// TestGetIndicesSnapshotSingleTicker verifies that GetIndicesSnapshot
// correctly parses the API response for a single index ticker, including
// the status, request ID, and all result fields.
func TestGetIndicesSnapshotSingleTicker(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/indices": indicesSnapshotSingleJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesSnapshot(IndicesSnapshotParams{
		TickerAnyOf: "I:SPX",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "idx123" {
		t.Errorf("expected request_id idx123, got %s", result.RequestID)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	idx := result.Results[0]
	if idx.Ticker != "I:SPX" {
		t.Errorf("expected ticker I:SPX, got %s", idx.Ticker)
	}

	if idx.Name != "Standard & Poor's 500" {
		t.Errorf("expected name Standard & Poor's 500, got %s", idx.Name)
	}

	if idx.Value != 6836.17 {
		t.Errorf("expected value 6836.17, got %f", idx.Value)
	}

	if idx.Type != "indices" {
		t.Errorf("expected type indices, got %s", idx.Type)
	}

	if idx.Timeframe != "REAL-TIME" {
		t.Errorf("expected timeframe REAL-TIME, got %s", idx.Timeframe)
	}

	if idx.MarketStatus != "closed" {
		t.Errorf("expected market_status closed, got %s", idx.MarketStatus)
	}

	if idx.LastUpdated != 1771016540019000000 {
		t.Errorf("expected last_updated 1771016540019000000, got %d", idx.LastUpdated)
	}
}

// TestGetIndicesSnapshotSessionData verifies that the session data within
// a single index snapshot is correctly parsed with change, change_percent,
// open, high, low, close, and previous_close values.
func TestGetIndicesSnapshotSessionData(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/indices": indicesSnapshotSingleJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesSnapshot(IndicesSnapshotParams{
		TickerAnyOf: "I:SPX",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	session := result.Results[0].Session
	if session.Change != -12.50 {
		t.Errorf("expected session change -12.50, got %f", session.Change)
	}

	if session.ChangePercent != -0.1828 {
		t.Errorf("expected session change_percent -0.1828, got %f", session.ChangePercent)
	}

	if session.Close != 6836.17 {
		t.Errorf("expected session close 6836.17, got %f", session.Close)
	}

	if session.High != 6881.96 {
		t.Errorf("expected session high 6881.96, got %f", session.High)
	}

	if session.Low != 6794.55 {
		t.Errorf("expected session low 6794.55, got %f", session.Low)
	}

	if session.Open != 6834.27 {
		t.Errorf("expected session open 6834.27, got %f", session.Open)
	}

	if session.PreviousClose != 6848.67 {
		t.Errorf("expected session previous_close 6848.67, got %f", session.PreviousClose)
	}
}

// TestGetIndicesSnapshotMultipleTickers verifies that GetIndicesSnapshot
// correctly parses a response containing multiple index tickers and that
// each ticker's data is independently and accurately parsed.
func TestGetIndicesSnapshotMultipleTickers(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/indices": indicesSnapshotMultipleJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesSnapshot(IndicesSnapshotParams{
		TickerAnyOf: "I:SPX,I:DJI",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result.Results))
	}

	if result.Results[0].Ticker != "I:SPX" {
		t.Errorf("expected first ticker I:SPX, got %s", result.Results[0].Ticker)
	}

	if result.Results[1].Ticker != "I:DJI" {
		t.Errorf("expected second ticker I:DJI, got %s", result.Results[1].Ticker)
	}
}

// TestGetIndicesSnapshotSecondTicker verifies that the second index in
// a multi-ticker response has its own distinct values parsed correctly,
// including name, value, and session data.
func TestGetIndicesSnapshotSecondTicker(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/indices": indicesSnapshotMultipleJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesSnapshot(IndicesSnapshotParams{
		TickerAnyOf: "I:SPX,I:DJI",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dji := result.Results[1]
	if dji.Name != "Dow Jones Industrial Average" {
		t.Errorf("expected name Dow Jones Industrial Average, got %s", dji.Name)
	}

	if dji.Value != 44546.08 {
		t.Errorf("expected value 44546.08, got %f", dji.Value)
	}

	if dji.Session.Change != 342.30 {
		t.Errorf("expected session change 342.30, got %f", dji.Session.Change)
	}

	if dji.Session.ChangePercent != 0.7742 {
		t.Errorf("expected session change_percent 0.7742, got %f", dji.Session.ChangePercent)
	}

	if dji.Session.PreviousClose != 44203.78 {
		t.Errorf("expected session previous_close 44203.78, got %f", dji.Session.PreviousClose)
	}
}

// TestGetIndicesSnapshotPagination verifies that the next_url field is
// correctly parsed from the API response when pagination is present.
func TestGetIndicesSnapshotPagination(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/indices": indicesSnapshotPaginatedJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesSnapshot(IndicesSnapshotParams{
		Limit: "1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.NextURL != "https://api.massive.com/v3/snapshot/indices?cursor=abc123" {
		t.Errorf("expected next_url with cursor, got %s", result.NextURL)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	idx := result.Results[0]
	if idx.Ticker != "I:COMP" {
		t.Errorf("expected ticker I:COMP, got %s", idx.Ticker)
	}

	if idx.Timeframe != "DELAYED" {
		t.Errorf("expected timeframe DELAYED, got %s", idx.Timeframe)
	}

	if idx.MarketStatus != "open" {
		t.Errorf("expected market_status open, got %s", idx.MarketStatus)
	}
}

// TestGetIndicesSnapshotRequestPath verifies that GetIndicesSnapshot
// constructs the correct API path for the indices snapshot endpoint.
func TestGetIndicesSnapshotRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesSnapshotSingleJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesSnapshot(IndicesSnapshotParams{})

	expected := "/v3/snapshot/indices"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetIndicesSnapshotQueryParams verifies that the ticker.any_of,
// limit, order, and sort query parameters are correctly sent to the API
// when specified in the params struct.
func TestGetIndicesSnapshotQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("ticker.any_of") != "I:SPX,I:DJI" {
			t.Errorf("expected ticker.any_of=I:SPX,I:DJI, got %s", q.Get("ticker.any_of"))
		}
		if q.Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", q.Get("limit"))
		}
		if q.Get("order") != "asc" {
			t.Errorf("expected order=asc, got %s", q.Get("order"))
		}
		if q.Get("sort") != "ticker" {
			t.Errorf("expected sort=ticker, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesSnapshotMultipleJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesSnapshot(IndicesSnapshotParams{
		TickerAnyOf: "I:SPX,I:DJI",
		Limit:       "50",
		Order:       "asc",
		Sort:        "ticker",
	})
}

// TestGetIndicesSnapshotTickerRangeParams verifies that the ticker range
// filter query parameters (ticker.gte, ticker.gt, ticker.lte, ticker.lt)
// are correctly sent to the API when specified.
func TestGetIndicesSnapshotTickerRangeParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("ticker.gte") != "I:A" {
			t.Errorf("expected ticker.gte=I:A, got %s", q.Get("ticker.gte"))
		}
		if q.Get("ticker.lte") != "I:Z" {
			t.Errorf("expected ticker.lte=I:Z, got %s", q.Get("ticker.lte"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesSnapshotSingleJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesSnapshot(IndicesSnapshotParams{
		TickerGte: "I:A",
		TickerLte: "I:Z",
	})
}

// TestGetIndicesSnapshotTickerSearchParam verifies that the ticker search
// query parameter is correctly sent to the API when specified.
func TestGetIndicesSnapshotTickerSearchParam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("ticker") != "I:SPX" {
			t.Errorf("expected ticker=I:SPX, got %s", q.Get("ticker"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesSnapshotSingleJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesSnapshot(IndicesSnapshotParams{
		Ticker: "I:SPX",
	})
}

// TestGetIndicesSnapshotAPIError verifies that GetIndicesSnapshot returns
// an error when the API responds with a non-200 status code.
func TestGetIndicesSnapshotAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"ERROR","message":"Not authorized."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetIndicesSnapshot(IndicesSnapshotParams{})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetIndicesSnapshotEmptyResults verifies that GetIndicesSnapshot
// correctly handles an empty results array from the API.
func TestGetIndicesSnapshotEmptyResults(t *testing.T) {
	emptyJSON := `{"status": "OK", "request_id": "empty123", "results": []}`
	server := mockServer(t, map[string]string{
		"/v3/snapshot/indices": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesSnapshot(IndicesSnapshotParams{
		TickerAnyOf: "I:NOTREAL",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(result.Results))
	}
}
