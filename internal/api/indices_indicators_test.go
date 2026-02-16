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

const indicesSMAJSON = `{
	"results": {
		"underlying": {
			"url": "https://api.polygon.io/v2/aggs/ticker/I:SPX/range/1/day/1731024000000/1736553600000?limit=50&sort=desc"
		},
		"values": [
			{
				"timestamp": 1736488800000,
				"value": 5923.772
			},
			{
				"timestamp": 1736316000000,
				"value": 5945.072
			},
			{
				"timestamp": 1736229600000,
				"value": 5950.654
			}
		]
	},
	"status": "OK",
	"request_id": "14b7b5fc629e4e75ec33538c4ec00c91"
}`

const indicesEMAJSON = `{
	"results": {
		"underlying": {
			"url": "https://api.polygon.io/v2/aggs/ticker/I:SPX/range/1/day/1731024000000/1736553600000?limit=50&sort=desc"
		},
		"values": [
			{
				"timestamp": 1736488800000,
				"value": 5935.1234
			},
			{
				"timestamp": 1736316000000,
				"value": 5948.5678
			}
		]
	},
	"status": "OK",
	"request_id": "ema-indices-req-001",
	"next_url": "https://api.massive.com/v1/indicators/ema/I:SPX?cursor=xyz789"
}`

const indicesRSIJSON = `{
	"results": {
		"underlying": {
			"url": "https://api.polygon.io/v2/aggs/ticker/I:SPX/range/1/day/1729483200000/1736553600000?limit=68&sort=desc"
		},
		"values": [
			{
				"timestamp": 1736488800000,
				"value": 42.3567
			},
			{
				"timestamp": 1736316000000,
				"value": 48.9012
			},
			{
				"timestamp": 1736229600000,
				"value": 47.2345
			}
		]
	},
	"status": "OK",
	"request_id": "rsi-indices-req-001"
}`

const indicesMACDJSON = `{
	"results": {
		"underlying": {
			"url": "https://api.polygon.io/v2/aggs/ticker/I:SPX/range/1/day/1724212800000/1736553600000?limit=122&sort=desc"
		},
		"values": [
			{
				"timestamp": 1736488800000,
				"value": 15.4321,
				"signal": 22.8765,
				"histogram": -7.4444
			},
			{
				"timestamp": 1736316000000,
				"value": 18.6543,
				"signal": 24.1234,
				"histogram": -5.4691
			}
		]
	},
	"status": "OK",
	"request_id": "macd-indices-req-001"
}`

// TestGetIndicesSMA verifies that GetIndicesSMA correctly parses the API
// response and returns the expected SMA indicator values for I:SPX.
func TestGetIndicesSMA(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/sma/I:SPX": indicesSMAJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := IndicatorParams{
		TimestampGTE: "2025-01-06",
		TimestampLTE: "2025-01-10",
		Timespan:     "day",
		Window:       "10",
		SeriesType:   "close",
		Limit:        "3",
	}

	result, err := client.GetIndicesSMA("I:SPX", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "14b7b5fc629e4e75ec33538c4ec00c91" {
		t.Errorf("expected request_id 14b7b5fc629e4e75ec33538c4ec00c91, got %s", result.RequestID)
	}

	if len(result.Results.Values) != 3 {
		t.Fatalf("expected 3 values, got %d", len(result.Results.Values))
	}

	first := result.Results.Values[0]
	if first.Timestamp != 1736488800000 {
		t.Errorf("expected timestamp 1736488800000, got %d", first.Timestamp)
	}

	if first.Value != 5923.772 {
		t.Errorf("expected value 5923.772, got %f", first.Value)
	}

	second := result.Results.Values[1]
	if second.Value != 5945.072 {
		t.Errorf("expected value 5945.072, got %f", second.Value)
	}
}

// TestGetIndicesSMARequestPath verifies that GetIndicesSMA constructs the
// correct API path with the index ticker symbol including the I: prefix.
func TestGetIndicesSMARequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesSMAJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesSMA("I:DJI", IndicatorParams{})

	if receivedPath != "/v1/indicators/sma/I:DJI" {
		t.Errorf("expected path /v1/indicators/sma/I:DJI, got %s", receivedPath)
	}
}

// TestGetIndicesSMAQueryParams verifies that all SMA query parameters are
// correctly sent to the API endpoint for index indicators.
func TestGetIndicesSMAQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("timestamp.gte") != "2025-01-06" {
			t.Errorf("expected timestamp.gte=2025-01-06, got %s", q.Get("timestamp.gte"))
		}
		if q.Get("timestamp.lte") != "2025-01-10" {
			t.Errorf("expected timestamp.lte=2025-01-10, got %s", q.Get("timestamp.lte"))
		}
		if q.Get("timespan") != "day" {
			t.Errorf("expected timespan=day, got %s", q.Get("timespan"))
		}
		if q.Get("window") != "10" {
			t.Errorf("expected window=10, got %s", q.Get("window"))
		}
		if q.Get("series_type") != "close" {
			t.Errorf("expected series_type=close, got %s", q.Get("series_type"))
		}
		if q.Get("adjusted") != "true" {
			t.Errorf("expected adjusted=true, got %s", q.Get("adjusted"))
		}
		if q.Get("order") != "desc" {
			t.Errorf("expected order=desc, got %s", q.Get("order"))
		}
		if q.Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesSMAJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesSMA("I:SPX", IndicatorParams{
		TimestampGTE: "2025-01-06",
		TimestampLTE: "2025-01-10",
		Timespan:     "day",
		Window:       "10",
		SeriesType:   "close",
		Adjusted:     "true",
		Order:        "desc",
		Limit:        "50",
	})
}

// TestGetIndicesSMAUnderlyingURL verifies that the underlying aggregates URL
// is correctly parsed from the indices SMA response.
func TestGetIndicesSMAUnderlyingURL(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/sma/I:SPX": indicesSMAJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesSMA("I:SPX", IndicatorParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedURL := "https://api.polygon.io/v2/aggs/ticker/I:SPX/range/1/day/1731024000000/1736553600000?limit=50&sort=desc"
	if result.Results.Underlying.URL != expectedURL {
		t.Errorf("expected underlying URL %s, got %s", expectedURL, result.Results.Underlying.URL)
	}
}

// TestGetIndicesSMAAPIError verifies that GetIndicesSMA returns an error
// when the API responds with a non-200 status code.
func TestGetIndicesSMAAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Data not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetIndicesSMA("I:INVALID", IndicatorParams{})
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// TestGetIndicesSMAEmptyValues verifies that GetIndicesSMA handles an empty
// values array without error.
func TestGetIndicesSMAEmptyValues(t *testing.T) {
	emptyJSON := `{"results":{"underlying":{"url":""},"values":[]},"status":"OK","request_id":"abc"}`
	server := mockServer(t, map[string]string{
		"/v1/indicators/sma/I:SPX": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesSMA("I:SPX", IndicatorParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results.Values) != 0 {
		t.Errorf("expected 0 values, got %d", len(result.Results.Values))
	}
}

// TestGetIndicesEMA verifies that GetIndicesEMA correctly parses the API
// response and returns the expected EMA indicator values for I:SPX.
func TestGetIndicesEMA(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/ema/I:SPX": indicesEMAJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := IndicatorParams{
		TimestampGTE: "2025-01-06",
		TimestampLTE: "2025-01-10",
		Timespan:     "day",
		Window:       "10",
		SeriesType:   "close",
		Limit:        "3",
	}

	result, err := client.GetIndicesEMA("I:SPX", params)
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
	if first.Timestamp != 1736488800000 {
		t.Errorf("expected timestamp 1736488800000, got %d", first.Timestamp)
	}

	if first.Value != 5935.1234 {
		t.Errorf("expected value 5935.1234, got %f", first.Value)
	}
}

// TestGetIndicesEMARequestPath verifies that GetIndicesEMA constructs the
// correct API path with the index ticker symbol.
func TestGetIndicesEMARequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesEMAJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesEMA("I:NDX", IndicatorParams{})

	if receivedPath != "/v1/indicators/ema/I:NDX" {
		t.Errorf("expected path /v1/indicators/ema/I:NDX, got %s", receivedPath)
	}
}

// TestGetIndicesEMANextURL verifies that the pagination next_url field is
// correctly parsed from the indices EMA response.
func TestGetIndicesEMANextURL(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/ema/I:SPX": indicesEMAJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesEMA("I:SPX", IndicatorParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.NextURL != "https://api.massive.com/v1/indicators/ema/I:SPX?cursor=xyz789" {
		t.Errorf("expected next_url with cursor, got %s", result.NextURL)
	}
}

// TestGetIndicesEMAAPIError verifies that GetIndicesEMA returns an error
// when the API responds with a non-200 status code.
func TestGetIndicesEMAAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetIndicesEMA("I:SPX", IndicatorParams{})
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

// TestGetIndicesRSI verifies that GetIndicesRSI correctly parses the API
// response and returns the expected RSI indicator values for I:SPX.
func TestGetIndicesRSI(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/rsi/I:SPX": indicesRSIJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := IndicatorParams{
		TimestampGTE: "2025-01-06",
		TimestampLTE: "2025-01-10",
		Timespan:     "day",
		Window:       "14",
		SeriesType:   "close",
		Limit:        "5",
	}

	result, err := client.GetIndicesRSI("I:SPX", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Results.Values) != 3 {
		t.Fatalf("expected 3 values, got %d", len(result.Results.Values))
	}

	first := result.Results.Values[0]
	if first.Timestamp != 1736488800000 {
		t.Errorf("expected timestamp 1736488800000, got %d", first.Timestamp)
	}

	if first.Value != 42.3567 {
		t.Errorf("expected value 42.3567, got %f", first.Value)
	}

	second := result.Results.Values[1]
	if second.Value != 48.9012 {
		t.Errorf("expected value 48.9012, got %f", second.Value)
	}
}

// TestGetIndicesRSIRequestPath verifies that GetIndicesRSI constructs the
// correct API path with the index ticker symbol.
func TestGetIndicesRSIRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesRSIJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesRSI("I:VIX", IndicatorParams{})

	if receivedPath != "/v1/indicators/rsi/I:VIX" {
		t.Errorf("expected path /v1/indicators/rsi/I:VIX, got %s", receivedPath)
	}
}

// TestGetIndicesRSIQueryParams verifies that RSI-specific query parameters
// including window and series_type are correctly sent for indices.
func TestGetIndicesRSIQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("window") != "14" {
			t.Errorf("expected window=14, got %s", q.Get("window"))
		}
		if q.Get("series_type") != "close" {
			t.Errorf("expected series_type=close, got %s", q.Get("series_type"))
		}
		if q.Get("timespan") != "day" {
			t.Errorf("expected timespan=day, got %s", q.Get("timespan"))
		}
		if q.Get("adjusted") != "false" {
			t.Errorf("expected adjusted=false, got %s", q.Get("adjusted"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesRSIJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesRSI("I:SPX", IndicatorParams{
		Window:     "14",
		SeriesType: "close",
		Timespan:   "day",
		Adjusted:   "false",
	})
}

// TestGetIndicesRSIAPIError verifies that GetIndicesRSI returns an error
// when the API responds with a non-200 status code.
func TestGetIndicesRSIAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Ticker not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetIndicesRSI("I:INVALID", IndicatorParams{})
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// TestGetIndicesMACD verifies that GetIndicesMACD correctly parses the API
// response and returns the expected MACD values including signal and histogram.
func TestGetIndicesMACD(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/macd/I:SPX": indicesMACDJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := MACDParams{
		TimestampGTE: "2025-01-06",
		TimestampLTE: "2025-01-10",
		Timespan:     "day",
		ShortWindow:  "12",
		LongWindow:   "26",
		SignalWindow: "9",
		SeriesType:   "close",
		Limit:        "5",
	}

	result, err := client.GetIndicesMACD("I:SPX", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "macd-indices-req-001" {
		t.Errorf("expected request_id macd-indices-req-001, got %s", result.RequestID)
	}

	if len(result.Results.Values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(result.Results.Values))
	}

	first := result.Results.Values[0]
	if first.Timestamp != 1736488800000 {
		t.Errorf("expected timestamp 1736488800000, got %d", first.Timestamp)
	}

	if first.Value != 15.4321 {
		t.Errorf("expected value 15.4321, got %f", first.Value)
	}

	if first.Signal != 22.8765 {
		t.Errorf("expected signal 22.8765, got %f", first.Signal)
	}

	if first.Histogram != -7.4444 {
		t.Errorf("expected histogram -7.4444, got %f", first.Histogram)
	}
}

// TestGetIndicesMACDRequestPath verifies that GetIndicesMACD constructs the
// correct API path with the index ticker symbol.
func TestGetIndicesMACDRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesMACDJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesMACD("I:DJI", MACDParams{})

	if receivedPath != "/v1/indicators/macd/I:DJI" {
		t.Errorf("expected path /v1/indicators/macd/I:DJI, got %s", receivedPath)
	}
}

// TestGetIndicesMACDQueryParams verifies that MACD-specific query parameters
// including short_window, long_window, and signal_window are correctly sent
// for index indicators.
func TestGetIndicesMACDQueryParams(t *testing.T) {
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
		if q.Get("series_type") != "close" {
			t.Errorf("expected series_type=close, got %s", q.Get("series_type"))
		}
		if q.Get("timespan") != "day" {
			t.Errorf("expected timespan=day, got %s", q.Get("timespan"))
		}
		if q.Get("timestamp.gte") != "2025-01-01" {
			t.Errorf("expected timestamp.gte=2025-01-01, got %s", q.Get("timestamp.gte"))
		}
		if q.Get("timestamp.lte") != "2025-01-31" {
			t.Errorf("expected timestamp.lte=2025-01-31, got %s", q.Get("timestamp.lte"))
		}
		if q.Get("limit") != "100" {
			t.Errorf("expected limit=100, got %s", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesMACDJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesMACD("I:SPX", MACDParams{
		TimestampGTE: "2025-01-01",
		TimestampLTE: "2025-01-31",
		Timespan:     "day",
		ShortWindow:  "12",
		LongWindow:   "26",
		SignalWindow: "9",
		SeriesType:   "close",
		Limit:        "100",
	})
}

// TestGetIndicesMACDSecondValue verifies that the second MACD value in the
// indices response is correctly parsed with its own distinct values.
func TestGetIndicesMACDSecondValue(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/macd/I:SPX": indicesMACDJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesMACD("I:SPX", MACDParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := result.Results.Values[1]
	if second.Value != 18.6543 {
		t.Errorf("expected value 18.6543, got %f", second.Value)
	}

	if second.Signal != 24.1234 {
		t.Errorf("expected signal 24.1234, got %f", second.Signal)
	}

	if second.Histogram != -5.4691 {
		t.Errorf("expected histogram -5.4691, got %f", second.Histogram)
	}
}

// TestGetIndicesMACDAPIError verifies that GetIndicesMACD returns an error
// when the API responds with a non-200 status code.
func TestGetIndicesMACDAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"ERROR","message":"Unauthorized"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetIndicesMACD("I:SPX", MACDParams{})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetIndicesMACDUnderlyingURL verifies that the underlying aggregates URL
// is correctly parsed from the indices MACD response.
func TestGetIndicesMACDUnderlyingURL(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/macd/I:SPX": indicesMACDJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesMACD("I:SPX", MACDParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedURL := "https://api.polygon.io/v2/aggs/ticker/I:SPX/range/1/day/1724212800000/1736553600000?limit=122&sort=desc"
	if result.Results.Underlying.URL != expectedURL {
		t.Errorf("expected underlying URL %s, got %s", expectedURL, result.Results.Underlying.URL)
	}
}

// TestGetIndicesMACDEmptyValues verifies that GetIndicesMACD handles an empty
// values array without error.
func TestGetIndicesMACDEmptyValues(t *testing.T) {
	emptyJSON := `{"results":{"underlying":{"url":""},"values":[]},"status":"OK","request_id":"abc"}`
	server := mockServer(t, map[string]string{
		"/v1/indicators/macd/I:SPX": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesMACD("I:SPX", MACDParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results.Values) != 0 {
		t.Errorf("expected 0 values, got %d", len(result.Results.Values))
	}
}
