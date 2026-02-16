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

const singleTickerSnapshotJSON = `{
	"status": "OK",
	"request_id": "abc123",
	"ticker": {
		"ticker": "AAPL",
		"todaysChange": -6.43,
		"todaysChangePerc": -2.46,
		"updated": 1771030800000000000,
		"day": {
			"o": 262.01,
			"h": 262.23,
			"l": 255.45,
			"c": 255.78,
			"v": 56291457.0,
			"vw": 257.98
		},
		"min": {
			"av": 56291457.0,
			"t": 1771030740000,
			"n": 53,
			"o": 255.31,
			"h": 255.32,
			"l": 255.30,
			"c": 255.30,
			"v": 2689,
			"vw": 255.31
		},
		"prevDay": {
			"o": 275.59,
			"h": 275.72,
			"l": 260.18,
			"c": 261.73,
			"v": 81063749.0,
			"vw": 264.27
		}
	}
}`

const allTickersSnapshotJSON = `{
	"status": "OK",
	"request_id": "def456",
	"count": 2,
	"tickers": [
		{
			"ticker": "AAPL",
			"todaysChange": -6.43,
			"todaysChangePerc": -2.46,
			"updated": 1771030800000000000,
			"day": {
				"o": 262.01,
				"h": 262.23,
				"l": 255.45,
				"c": 255.78,
				"v": 56291457.0,
				"vw": 257.98
			},
			"min": {
				"av": 56291457.0,
				"t": 1771030740000,
				"n": 53,
				"o": 255.31,
				"h": 255.32,
				"l": 255.30,
				"c": 255.30,
				"v": 2689,
				"vw": 255.31
			},
			"prevDay": {
				"o": 275.59,
				"h": 275.72,
				"l": 260.18,
				"c": 261.73,
				"v": 81063749.0,
				"vw": 264.27
			}
		},
		{
			"ticker": "MSFT",
			"todaysChange": -1.69,
			"todaysChangePerc": -0.42,
			"updated": 1771030800000000000,
			"day": {
				"o": 404.45,
				"h": 405.54,
				"l": 398.05,
				"c": 401.32,
				"v": 34138860.0,
				"vw": 401.93
			},
			"min": {
				"av": 34138860.0,
				"t": 1771030740000,
				"n": 70,
				"o": 400.18,
				"h": 400.20,
				"l": 400.15,
				"c": 400.15,
				"v": 1935,
				"vw": 400.17
			},
			"prevDay": {
				"o": 405.00,
				"h": 406.20,
				"l": 398.01,
				"c": 401.84,
				"v": 40722419.0,
				"vw": 402.36
			}
		}
	]
}`

const gainersSnapshotJSON = `{
	"status": "OK",
	"request_id": "ghi789",
	"tickers": [
		{
			"ticker": "RIME",
			"todaysChange": 3.51,
			"todaysChangePerc": 325.00,
			"updated": 1771030800000000000,
			"day": {
				"o": 1.30,
				"h": 3.65,
				"l": 1.16,
				"c": 3.48,
				"v": 167399557.0,
				"vw": 2.62
			},
			"min": {
				"av": 167399557.0,
				"t": 1771030740000,
				"n": 757,
				"o": 4.51,
				"h": 4.60,
				"l": 4.43,
				"c": 4.59,
				"v": 210468,
				"vw": 4.50
			},
			"prevDay": {
				"o": 0.84,
				"h": 1.51,
				"l": 0.83,
				"c": 1.08,
				"v": 53988053.0,
				"vw": 1.15
			}
		},
		{
			"ticker": "ATOM",
			"todaysChange": 1.47,
			"todaysChangePerc": 61.51,
			"updated": 1771030800000000000,
			"day": {
				"o": 2.78,
				"h": 4.02,
				"l": 2.67,
				"c": 3.92,
				"v": 28258702.0,
				"vw": 3.58
			},
			"min": {
				"av": 28258702.0,
				"t": 1771030740000,
				"n": 100,
				"o": 3.90,
				"h": 3.95,
				"l": 3.88,
				"c": 3.92,
				"v": 50000,
				"vw": 3.91
			},
			"prevDay": {
				"o": 2.50,
				"h": 2.60,
				"l": 2.30,
				"c": 2.39,
				"v": 5000000.0,
				"vw": 2.45
			}
		}
	]
}`

// TestGetSnapshotTicker verifies that GetSnapshotTicker correctly parses
// the API response and returns the expected snapshot data for AAPL,
// including the day bar, minute bar, previous day bar, and change values.
func TestGetSnapshotTicker(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/snapshot/locale/us/markets/stocks/tickers/AAPL": singleTickerSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetSnapshotTicker("AAPL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "abc123" {
		t.Errorf("expected request_id abc123, got %s", result.RequestID)
	}

	if result.Ticker.Ticker != "AAPL" {
		t.Errorf("expected ticker AAPL, got %s", result.Ticker.Ticker)
	}

	if result.Ticker.TodaysChange != -6.43 {
		t.Errorf("expected todaysChange -6.43, got %f", result.Ticker.TodaysChange)
	}

	if result.Ticker.TodaysChangePct != -2.46 {
		t.Errorf("expected todaysChangePerc -2.46, got %f", result.Ticker.TodaysChangePct)
	}

	if result.Ticker.Updated != 1771030800000000000 {
		t.Errorf("expected updated 1771030800000000000, got %d", result.Ticker.Updated)
	}
}

// TestGetSnapshotTickerDayBar verifies that the day bar within a single
// ticker snapshot is correctly parsed with open, high, low, close, volume,
// and volume-weighted average price.
func TestGetSnapshotTickerDayBar(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/snapshot/locale/us/markets/stocks/tickers/AAPL": singleTickerSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetSnapshotTicker("AAPL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	day := result.Ticker.Day
	if day.Open != 262.01 {
		t.Errorf("expected day open 262.01, got %f", day.Open)
	}

	if day.High != 262.23 {
		t.Errorf("expected day high 262.23, got %f", day.High)
	}

	if day.Low != 255.45 {
		t.Errorf("expected day low 255.45, got %f", day.Low)
	}

	if day.Close != 255.78 {
		t.Errorf("expected day close 255.78, got %f", day.Close)
	}

	if day.Volume != 56291457 {
		t.Errorf("expected day volume 56291457, got %f", day.Volume)
	}

	if day.VWAP != 257.98 {
		t.Errorf("expected day vwap 257.98, got %f", day.VWAP)
	}
}

// TestGetSnapshotTickerPrevDay verifies that the previous day's bar
// data is correctly parsed from the single ticker snapshot response.
func TestGetSnapshotTickerPrevDay(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/snapshot/locale/us/markets/stocks/tickers/AAPL": singleTickerSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetSnapshotTicker("AAPL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	prev := result.Ticker.PrevDay
	if prev.Open != 275.59 {
		t.Errorf("expected prevDay open 275.59, got %f", prev.Open)
	}

	if prev.Close != 261.73 {
		t.Errorf("expected prevDay close 261.73, got %f", prev.Close)
	}

	if prev.Volume != 81063749 {
		t.Errorf("expected prevDay volume 81063749, got %f", prev.Volume)
	}
}

// TestGetSnapshotTickerMinBar verifies that the most recent minute bar
// is correctly parsed, including accumulated volume and transaction count.
func TestGetSnapshotTickerMinBar(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/snapshot/locale/us/markets/stocks/tickers/AAPL": singleTickerSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetSnapshotTicker("AAPL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	min := result.Ticker.Min
	if min.Open != 255.31 {
		t.Errorf("expected min open 255.31, got %f", min.Open)
	}

	if min.AccumulatedVolume != 56291457 {
		t.Errorf("expected min accumulated volume 56291457, got %f", min.AccumulatedVolume)
	}

	if min.Timestamp != 1771030740000 {
		t.Errorf("expected min timestamp 1771030740000, got %d", min.Timestamp)
	}

	if min.NumTransactions != 53 {
		t.Errorf("expected min transactions 53, got %d", min.NumTransactions)
	}

	if min.Volume != 2689 {
		t.Errorf("expected min volume 2689, got %f", min.Volume)
	}
}

// TestGetSnapshotTickerRequestPath verifies that GetSnapshotTicker
// constructs the correct API path with the ticker symbol.
func TestGetSnapshotTickerRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(singleTickerSnapshotJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetSnapshotTicker("MSFT")

	expected := "/v2/snapshot/locale/us/markets/stocks/tickers/MSFT"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetSnapshotTickerAPIError verifies that GetSnapshotTicker returns
// an error when the API responds with a non-200 status code.
func TestGetSnapshotTickerAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Ticker not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetSnapshotTicker("INVALID")
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// TestGetSnapshotAllTickers verifies that GetSnapshotAllTickers correctly
// parses the multi-ticker snapshot response with count and ticker array.
func TestGetSnapshotAllTickers(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/snapshot/locale/us/markets/stocks/tickers": allTickersSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := AllTickersSnapshotParams{
		Tickers: "AAPL,MSFT",
	}

	result, err := client.GetSnapshotAllTickers(params)
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

	if result.Tickers[0].Ticker != "AAPL" {
		t.Errorf("expected first ticker AAPL, got %s", result.Tickers[0].Ticker)
	}

	if result.Tickers[1].Ticker != "MSFT" {
		t.Errorf("expected second ticker MSFT, got %s", result.Tickers[1].Ticker)
	}
}

// TestGetSnapshotAllTickersSecondTicker verifies that the second ticker
// in the multi-ticker response has its own distinct values parsed correctly.
func TestGetSnapshotAllTickersSecondTicker(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/snapshot/locale/us/markets/stocks/tickers": allTickersSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetSnapshotAllTickers(AllTickersSnapshotParams{Tickers: "AAPL,MSFT"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msft := result.Tickers[1]
	if msft.TodaysChange != -1.69 {
		t.Errorf("expected MSFT todaysChange -1.69, got %f", msft.TodaysChange)
	}

	if msft.Day.Open != 404.45 {
		t.Errorf("expected MSFT day open 404.45, got %f", msft.Day.Open)
	}

	if msft.Day.Close != 401.32 {
		t.Errorf("expected MSFT day close 401.32, got %f", msft.Day.Close)
	}

	if msft.PrevDay.Close != 401.84 {
		t.Errorf("expected MSFT prevDay close 401.84, got %f", msft.PrevDay.Close)
	}
}

// TestGetSnapshotAllTickersQueryParams verifies that the tickers and
// include_otc query parameters are correctly sent to the API.
func TestGetSnapshotAllTickersQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("tickers") != "AAPL,TSLA" {
			t.Errorf("expected tickers=AAPL,TSLA, got %s", q.Get("tickers"))
		}
		if q.Get("include_otc") != "true" {
			t.Errorf("expected include_otc=true, got %s", q.Get("include_otc"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(allTickersSnapshotJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetSnapshotAllTickers(AllTickersSnapshotParams{
		Tickers:    "AAPL,TSLA",
		IncludeOTC: "true",
	})
}

// TestGetSnapshotAllTickersRequestPath verifies that GetSnapshotAllTickers
// constructs the correct API path for the full market snapshot endpoint.
func TestGetSnapshotAllTickersRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(allTickersSnapshotJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetSnapshotAllTickers(AllTickersSnapshotParams{})

	expected := "/v2/snapshot/locale/us/markets/stocks/tickers"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetSnapshotGainersLosersGainers verifies that GetSnapshotGainersLosers
// correctly parses the gainers response with ticker data sorted by gain percentage.
func TestGetSnapshotGainersLosersGainers(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/snapshot/locale/us/markets/stocks/gainers": gainersSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetSnapshotGainersLosers("gainers", GainersLosersParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Tickers) != 2 {
		t.Fatalf("expected 2 tickers, got %d", len(result.Tickers))
	}

	first := result.Tickers[0]
	if first.Ticker != "RIME" {
		t.Errorf("expected first ticker RIME, got %s", first.Ticker)
	}

	if first.TodaysChangePct != 325.00 {
		t.Errorf("expected todaysChangePerc 325.00, got %f", first.TodaysChangePct)
	}

	if first.TodaysChange != 3.51 {
		t.Errorf("expected todaysChange 3.51, got %f", first.TodaysChange)
	}

	if first.Day.Open != 1.30 {
		t.Errorf("expected day open 1.30, got %f", first.Day.Open)
	}

	if first.Day.Close != 3.48 {
		t.Errorf("expected day close 3.48, got %f", first.Day.Close)
	}
}

// TestGetSnapshotGainersLosersLosers verifies that GetSnapshotGainersLosers
// constructs the correct path when the direction is set to "losers".
func TestGetSnapshotGainersLosersLosers(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(gainersSnapshotJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetSnapshotGainersLosers("losers", GainersLosersParams{})

	expected := "/v2/snapshot/locale/us/markets/stocks/losers"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetSnapshotGainersLosersGainersPath verifies that GetSnapshotGainersLosers
// constructs the correct path when the direction is set to "gainers".
func TestGetSnapshotGainersLosersGainersPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(gainersSnapshotJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetSnapshotGainersLosers("gainers", GainersLosersParams{})

	expected := "/v2/snapshot/locale/us/markets/stocks/gainers"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetSnapshotGainersLosersIncludeOTC verifies that the include_otc
// query parameter is correctly sent when requesting gainers or losers.
func TestGetSnapshotGainersLosersIncludeOTC(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("include_otc") != "true" {
			t.Errorf("expected include_otc=true, got %s", q.Get("include_otc"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(gainersSnapshotJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetSnapshotGainersLosers("gainers", GainersLosersParams{
		IncludeOTC: "true",
	})
}

// TestGetSnapshotGainersLosersSecondTicker verifies that the second
// ticker in the gainers/losers response is correctly parsed.
func TestGetSnapshotGainersLosersSecondTicker(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/snapshot/locale/us/markets/stocks/gainers": gainersSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetSnapshotGainersLosers("gainers", GainersLosersParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := result.Tickers[1]
	if second.Ticker != "ATOM" {
		t.Errorf("expected second ticker ATOM, got %s", second.Ticker)
	}

	if second.TodaysChangePct != 61.51 {
		t.Errorf("expected todaysChangePerc 61.51, got %f", second.TodaysChangePct)
	}

	if second.Day.High != 4.02 {
		t.Errorf("expected day high 4.02, got %f", second.Day.High)
	}

	if second.PrevDay.Close != 2.39 {
		t.Errorf("expected prevDay close 2.39, got %f", second.PrevDay.Close)
	}
}

// TestGetSnapshotGainersLosersAPIError verifies that GetSnapshotGainersLosers
// returns an error when the API responds with a non-200 status code.
func TestGetSnapshotGainersLosersAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"ERROR","message":"Not authorized."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetSnapshotGainersLosers("gainers", GainersLosersParams{})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}
