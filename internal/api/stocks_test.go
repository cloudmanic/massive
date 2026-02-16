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

// mockServer creates a test HTTP server that responds to specific paths
// with the provided JSON responses. Returns the server (caller must close).
func mockServer(t *testing.T, routes map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, ok := routes[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"status":"NOT_FOUND","message":"No route matched"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(body))
	}))
}

// newTestClient creates a Client pointed at the given test server URL.
func newTestClient(serverURL string) *Client {
	client := NewClient("test-api-key")
	client.SetBaseURL(serverURL)
	return client
}

const openCloseJSON = `{
	"status": "OK",
	"from": "2025-01-06",
	"symbol": "AAPL",
	"open": 244.31,
	"high": 247.33,
	"low": 243.2,
	"close": 245,
	"volume": 45045571,
	"afterHours": 244.99,
	"preMarket": 243.89
}`

const barsJSON = `{
	"ticker": "AAPL",
	"queryCount": 2,
	"resultsCount": 2,
	"adjusted": true,
	"results": [
		{
			"v": 45045571.0,
			"vw": 245.1726,
			"o": 244.31,
			"c": 245,
			"h": 247.33,
			"l": 243.2,
			"t": 1736139600000,
			"n": 493920
		},
		{
			"v": 40855960.0,
			"vw": 242.9459,
			"o": 242.98,
			"c": 242.21,
			"h": 245.55,
			"l": 241.35,
			"t": 1736226000000,
			"n": 462887
		}
	],
	"status": "OK",
	"request_id": "6f9a8c20663b1bc118723f66e637c046",
	"count": 2
}`

const marketSummaryJSON = `{
	"queryCount": 3,
	"resultsCount": 3,
	"adjusted": true,
	"results": [
		{
			"T": "YXT",
			"v": 23256,
			"vw": 2.127,
			"o": 2.1,
			"c": 2.09,
			"h": 2.2667,
			"l": 2.07,
			"t": 1736197200000,
			"n": 326
		},
		{
			"T": "APT",
			"v": 131308,
			"vw": 5.8533,
			"o": 5.65,
			"c": 5.92,
			"h": 5.99,
			"l": 5.6,
			"t": 1736197200000,
			"n": 784
		},
		{
			"T": "MFC",
			"v": 999319,
			"vw": 30.9212,
			"o": 31.09,
			"c": 30.81,
			"h": 31.25,
			"l": 30.755,
			"t": 1736197200000,
			"n": 9935
		}
	],
	"status": "OK",
	"request_id": "abc123"
}`

const tickersJSON = `{
	"results": [
		{
			"ticker": "AAPI",
			"name": "APPLE ISPORTS GROUP INC",
			"market": "otc",
			"locale": "us",
			"primary_exchange": "OTC Link",
			"type": "CS",
			"active": true,
			"currency_name": "USD",
			"composite_figi": "BBG000CQN9X7",
			"share_class_figi": "BBG001SGHPN2",
			"last_updated_utc": "2024-12-30T07:38:03.949Z"
		},
		{
			"ticker": "AAPL",
			"name": "Apple Inc.",
			"market": "stocks",
			"locale": "us",
			"primary_exchange": "XNAS",
			"type": "CS",
			"active": true,
			"currency_name": "usd",
			"cik": "0000320193",
			"composite_figi": "BBG000B9XRY4",
			"share_class_figi": "BBG001S5N8V8",
			"last_updated_utc": "2026-02-15T07:08:17.692Z"
		}
	],
	"status": "OK",
	"request_id": "39de8d8db5f75e8ffd70bd10dc6cec50",
	"count": 2,
	"next_url": "https://api.massive.com/v3/reference/tickers?cursor=YXA9Mg"
}`

// TestGetOpenClose verifies that GetOpenClose correctly parses the API
// response and returns the expected open/close data for AAPL.
func TestGetOpenClose(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/open-close/AAPL/2025-01-06": openCloseJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOpenClose("AAPL", "2025-01-06", "true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Symbol != "AAPL" {
		t.Errorf("expected symbol AAPL, got %s", result.Symbol)
	}

	if result.From != "2025-01-06" {
		t.Errorf("expected from 2025-01-06, got %s", result.From)
	}

	if result.Open != 244.31 {
		t.Errorf("expected open 244.31, got %f", result.Open)
	}

	if result.High != 247.33 {
		t.Errorf("expected high 247.33, got %f", result.High)
	}

	if result.Low != 243.2 {
		t.Errorf("expected low 243.2, got %f", result.Low)
	}

	if result.Close != 245 {
		t.Errorf("expected close 245, got %f", result.Close)
	}

	if result.Volume != 45045571 {
		t.Errorf("expected volume 45045571, got %d", result.Volume)
	}

	if result.AfterHours != 244.99 {
		t.Errorf("expected afterHours 244.99, got %f", result.AfterHours)
	}

	if result.PreMarket != 243.89 {
		t.Errorf("expected preMarket 243.89, got %f", result.PreMarket)
	}
}

// TestGetOpenCloseRequestPath verifies that GetOpenClose constructs the
// correct API path with the ticker and date.
func TestGetOpenCloseRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(openCloseJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOpenClose("MSFT", "2025-03-15", "false")

	if receivedPath != "/v1/open-close/MSFT/2025-03-15" {
		t.Errorf("expected path /v1/open-close/MSFT/2025-03-15, got %s", receivedPath)
	}
}

// TestGetOpenCloseAdjustedParam verifies that the adjusted query
// parameter is correctly sent to the API.
func TestGetOpenCloseAdjustedParam(t *testing.T) {
	var receivedAdjusted string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAdjusted = r.URL.Query().Get("adjusted")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(openCloseJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOpenClose("AAPL", "2025-01-06", "false")

	if receivedAdjusted != "false" {
		t.Errorf("expected adjusted=false, got %s", receivedAdjusted)
	}
}

// TestGetOpenCloseAPIError verifies that GetOpenClose returns an error
// when the API responds with a non-200 status.
func TestGetOpenCloseAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Data not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetOpenClose("INVALID", "2025-01-06", "true")
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// TestGetBars verifies that GetBars correctly parses the API response
// and returns the expected OHLC bar data.
func TestGetBars(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/aggs/ticker/AAPL/range/1/day/2025-01-06/2025-01-08": barsJSON,
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

	result, err := client.GetBars("AAPL", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Ticker != "AAPL" {
		t.Errorf("expected ticker AAPL, got %s", result.Ticker)
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
	if bar.Open != 244.31 {
		t.Errorf("expected open 244.31, got %f", bar.Open)
	}

	if bar.High != 247.33 {
		t.Errorf("expected high 247.33, got %f", bar.High)
	}

	if bar.Low != 243.2 {
		t.Errorf("expected low 243.2, got %f", bar.Low)
	}

	if bar.Close != 245 {
		t.Errorf("expected close 245, got %f", bar.Close)
	}

	if bar.Volume != 45045571 {
		t.Errorf("expected volume 45045571, got %f", bar.Volume)
	}

	if bar.VWAP != 245.1726 {
		t.Errorf("expected VWAP 245.1726, got %f", bar.VWAP)
	}

	if bar.Timestamp != 1736139600000 {
		t.Errorf("expected timestamp 1736139600000, got %d", bar.Timestamp)
	}

	if bar.NumTrades != 493920 {
		t.Errorf("expected 493920 trades, got %d", bar.NumTrades)
	}
}

// TestGetBarsRequestPath verifies that GetBars constructs the correct
// URL path with ticker, multiplier, timespan, from, and to values.
func TestGetBarsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(barsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	params := BarsParams{
		Multiplier: "5",
		Timespan:   "minute",
		From:       "2025-01-06",
		To:         "2025-01-07",
	}

	client.GetBars("TSLA", params)

	expected := "/v2/aggs/ticker/TSLA/range/5/minute/2025-01-06/2025-01-07"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetBarsQueryParams verifies that GetBars sends the correct query
// parameters including adjusted, sort, and limit.
func TestGetBarsQueryParams(t *testing.T) {
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
		w.Write([]byte(barsJSON))
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

	client.GetBars("AAPL", params)
}

// TestGetBarsSecondBar verifies that the second bar in the response
// is correctly parsed with its own distinct values.
func TestGetBarsSecondBar(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/aggs/ticker/AAPL/range/1/day/2025-01-06/2025-01-08": barsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := BarsParams{
		Multiplier: "1",
		Timespan:   "day",
		From:       "2025-01-06",
		To:         "2025-01-08",
	}

	result, err := client.GetBars("AAPL", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bar := result.Results[1]
	if bar.Open != 242.98 {
		t.Errorf("expected open 242.98, got %f", bar.Open)
	}

	if bar.Close != 242.21 {
		t.Errorf("expected close 242.21, got %f", bar.Close)
	}

	if bar.NumTrades != 462887 {
		t.Errorf("expected 462887 trades, got %d", bar.NumTrades)
	}
}

// TestGetMarketSummary verifies that GetMarketSummary correctly parses
// the grouped daily response with multiple tickers.
func TestGetMarketSummary(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/aggs/grouped/locale/us/market/stocks/2025-01-06": marketSummaryJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := MarketSummaryParams{
		Adjusted:   "true",
		IncludeOTC: "false",
	}

	result, err := client.GetMarketSummary("2025-01-06", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if !result.Adjusted {
		t.Error("expected adjusted to be true")
	}

	if result.ResultsCount != 3 {
		t.Errorf("expected 3 results, got %d", result.ResultsCount)
	}

	if len(result.Results) != 3 {
		t.Fatalf("expected 3 market summaries, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.Ticker != "YXT" {
		t.Errorf("expected ticker YXT, got %s", first.Ticker)
	}

	if first.Open != 2.1 {
		t.Errorf("expected open 2.1, got %f", first.Open)
	}

	if first.Close != 2.09 {
		t.Errorf("expected close 2.09, got %f", first.Close)
	}

	if first.NumTrades != 326 {
		t.Errorf("expected 326 trades, got %d", first.NumTrades)
	}

	third := result.Results[2]
	if third.Ticker != "MFC" {
		t.Errorf("expected ticker MFC, got %s", third.Ticker)
	}

	if third.Volume != 999319 {
		t.Errorf("expected volume 999319, got %f", third.Volume)
	}
}

// TestGetMarketSummaryRequestPath verifies that GetMarketSummary
// constructs the correct API path with the given date.
func TestGetMarketSummaryRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(marketSummaryJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetMarketSummary("2025-06-15", MarketSummaryParams{})

	expected := "/v2/aggs/grouped/locale/us/market/stocks/2025-06-15"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetMarketSummaryQueryParams verifies that the adjusted and
// include_otc parameters are correctly sent.
func TestGetMarketSummaryQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("adjusted") != "false" {
			t.Errorf("expected adjusted=false, got %s", q.Get("adjusted"))
		}
		if q.Get("include_otc") != "true" {
			t.Errorf("expected include_otc=true, got %s", q.Get("include_otc"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(marketSummaryJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetMarketSummary("2025-01-06", MarketSummaryParams{
		Adjusted:   "false",
		IncludeOTC: "true",
	})
}

// TestGetTickers verifies that GetTickers correctly parses the
// reference tickers response including pagination info.
func TestGetTickers(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/reference/tickers": tickersJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := TickerParams{
		Search: "Apple",
		Limit:  "2",
	}

	result, err := client.GetTickers(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Count != 2 {
		t.Errorf("expected count 2, got %d", result.Count)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 tickers, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.Ticker != "AAPI" {
		t.Errorf("expected ticker AAPI, got %s", first.Ticker)
	}

	if first.Name != "APPLE ISPORTS GROUP INC" {
		t.Errorf("expected name APPLE ISPORTS GROUP INC, got %s", first.Name)
	}

	if first.Market != "otc" {
		t.Errorf("expected market otc, got %s", first.Market)
	}

	if !first.Active {
		t.Error("expected active to be true")
	}

	second := result.Results[1]
	if second.Ticker != "AAPL" {
		t.Errorf("expected ticker AAPL, got %s", second.Ticker)
	}

	if second.Name != "Apple Inc." {
		t.Errorf("expected name Apple Inc., got %s", second.Name)
	}

	if second.PrimaryExchange != "XNAS" {
		t.Errorf("expected exchange XNAS, got %s", second.PrimaryExchange)
	}

	if second.CIK != "0000320193" {
		t.Errorf("expected CIK 0000320193, got %s", second.CIK)
	}

	if second.CompositeFIGI != "BBG000B9XRY4" {
		t.Errorf("expected composite FIGI BBG000B9XRY4, got %s", second.CompositeFIGI)
	}
}

// TestGetTickersQueryParams verifies that all filter parameters are
// correctly sent to the API endpoint.
func TestGetTickersQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("search") != "Tesla" {
			t.Errorf("expected search=Tesla, got %s", q.Get("search"))
		}
		if q.Get("type") != "CS" {
			t.Errorf("expected type=CS, got %s", q.Get("type"))
		}
		if q.Get("market") != "stocks" {
			t.Errorf("expected market=stocks, got %s", q.Get("market"))
		}
		if q.Get("exchange") != "XNAS" {
			t.Errorf("expected exchange=XNAS, got %s", q.Get("exchange"))
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
		w.Write([]byte(tickersJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetTickers(TickerParams{
		Search:   "Tesla",
		Type:     "CS",
		Market:   "stocks",
		Exchange: "XNAS",
		Active:   "true",
		Sort:     "name",
		Order:    "desc",
		Limit:    "50",
	})
}

// TestGetTickersEmptyResults verifies that GetTickers handles an
// empty results array without error.
func TestGetTickersEmptyResults(t *testing.T) {
	emptyJSON := `{"results":[],"status":"OK","request_id":"abc","count":0}`
	server := mockServer(t, map[string]string{
		"/v3/reference/tickers": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetTickers(TickerParams{Search: "zzzznotreal"})
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

// TestGetTickersWithTickerFilter verifies that the ticker param is
// sent when filtering by a specific ticker symbol.
func TestGetTickersWithTickerFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("ticker") != "AAPL" {
			t.Errorf("expected ticker=AAPL, got %s", r.URL.Query().Get("ticker"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(tickersJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetTickers(TickerParams{Ticker: "AAPL"})
}
