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

const indicesTickersJSON = `{
	"results": [
		{
			"ticker": "I:A1BSC",
			"name": "Dow Jones Americas Basic Materials Index",
			"market": "indices",
			"locale": "us",
			"active": true,
			"source_feed": "CMEMarketDataPlatformDowJones"
		},
		{
			"ticker": "I:SPX",
			"name": "S&P 500",
			"market": "indices",
			"locale": "us",
			"active": true,
			"source_feed": "CboeGlobalIndicesMain"
		}
	],
	"status": "OK",
	"request_id": "abc123indices",
	"count": 2,
	"next_url": "https://api.massive.com/v3/reference/tickers?cursor=YXA9Mg&market=indices"
}`

// TestGetIndicesTickers verifies that GetIndicesTickers correctly parses
// the API response and returns the expected index ticker data including
// ticker symbol, name, market, source feed, and pagination info.
func TestGetIndicesTickers(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/reference/tickers": indicesTickersJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := IndicesTickerParams{
		Search: "Dow Jones",
		Limit:  "2",
	}

	result, err := client.GetIndicesTickers(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Count != 2 {
		t.Errorf("expected count 2, got %d", result.Count)
	}

	if result.RequestID != "abc123indices" {
		t.Errorf("expected request_id abc123indices, got %s", result.RequestID)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 tickers, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.Ticker != "I:A1BSC" {
		t.Errorf("expected ticker I:A1BSC, got %s", first.Ticker)
	}

	if first.Name != "Dow Jones Americas Basic Materials Index" {
		t.Errorf("expected name Dow Jones Americas Basic Materials Index, got %s", first.Name)
	}

	if first.Market != "indices" {
		t.Errorf("expected market indices, got %s", first.Market)
	}

	if first.Locale != "us" {
		t.Errorf("expected locale us, got %s", first.Locale)
	}

	if !first.Active {
		t.Error("expected active to be true")
	}

	if first.SourceFeed != "CMEMarketDataPlatformDowJones" {
		t.Errorf("expected source_feed CMEMarketDataPlatformDowJones, got %s", first.SourceFeed)
	}

	second := result.Results[1]
	if second.Ticker != "I:SPX" {
		t.Errorf("expected ticker I:SPX, got %s", second.Ticker)
	}

	if second.Name != "S&P 500" {
		t.Errorf("expected name S&P 500, got %s", second.Name)
	}

	if second.SourceFeed != "CboeGlobalIndicesMain" {
		t.Errorf("expected source_feed CboeGlobalIndicesMain, got %s", second.SourceFeed)
	}
}

// TestGetIndicesTickersMarketParam verifies that GetIndicesTickers always
// sends market=indices to the API regardless of other parameters.
func TestGetIndicesTickersMarketParam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("market") != "indices" {
			t.Errorf("expected market=indices, got %s", q.Get("market"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesTickersJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesTickers(IndicesTickerParams{})
}

// TestGetIndicesTickersQueryParams verifies that all filter parameters
// are correctly sent to the API endpoint including search, active,
// sort, order, and limit.
func TestGetIndicesTickersQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("market") != "indices" {
			t.Errorf("expected market=indices, got %s", q.Get("market"))
		}
		if q.Get("search") != "S&P" {
			t.Errorf("expected search=S&P, got %s", q.Get("search"))
		}
		if q.Get("active") != "true" {
			t.Errorf("expected active=true, got %s", q.Get("active"))
		}
		if q.Get("sort") != "name" {
			t.Errorf("expected sort=name, got %s", q.Get("sort"))
		}
		if q.Get("order") != "desc" {
			t.Errorf("expected order=desc, got %s", q.Get("order"))
		}
		if q.Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesTickersJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesTickers(IndicesTickerParams{
		Search: "S&P",
		Active: "true",
		Sort:   "name",
		Order:  "desc",
		Limit:  "50",
	})
}

// TestGetIndicesTickersWithTickerFilter verifies that the ticker param
// is sent when filtering by a specific index ticker symbol.
func TestGetIndicesTickersWithTickerFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("ticker") != "I:SPX" {
			t.Errorf("expected ticker=I:SPX, got %s", r.URL.Query().Get("ticker"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesTickersJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesTickers(IndicesTickerParams{Ticker: "I:SPX"})
}

// TestGetIndicesTickersEmptyResults verifies that GetIndicesTickers handles
// an empty results array without error.
func TestGetIndicesTickersEmptyResults(t *testing.T) {
	emptyJSON := `{"results":[],"status":"OK","request_id":"empty123","count":0}`
	server := mockServer(t, map[string]string{
		"/v3/reference/tickers": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesTickers(IndicesTickerParams{Search: "zzzznotreal"})
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

// TestGetIndicesTickersAPIError verifies that GetIndicesTickers returns
// an error when the API responds with a non-200 status code.
func TestGetIndicesTickersAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"status":"ERROR","message":"Invalid API key"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetIndicesTickers(IndicesTickerParams{})
	if err == nil {
		t.Fatal("expected error for 401 response, got nil")
	}
}

// TestGetIndicesTickersRequestPath verifies that GetIndicesTickers sends
// requests to the correct /v3/reference/tickers API path.
func TestGetIndicesTickersRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesTickersJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesTickers(IndicesTickerParams{})

	if receivedPath != "/v3/reference/tickers" {
		t.Errorf("expected path /v3/reference/tickers, got %s", receivedPath)
	}
}
