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

const optionsSMAJSON = `{
	"results": {
		"underlying": {
			"url": "https://api.polygon.io/v2/aggs/ticker/O:AAPL250117C00150000/range/1/day/1731042000000/1736553600000?limit=50&sort=desc"
		},
		"values": [
			{
				"timestamp": 1736485200000,
				"value": 95.432
			},
			{
				"timestamp": 1736312400000,
				"value": 96.781
			},
			{
				"timestamp": 1736226000000,
				"value": 97.115
			}
		]
	},
	"status": "OK",
	"request_id": "opt-sma-req-001"
}`

const optionsEMAJSON = `{
	"results": {
		"underlying": {
			"url": "https://api.polygon.io/v2/aggs/ticker/O:AAPL250117C00150000/range/1/day/1731042000000/1736553600000?limit=50&sort=desc"
		},
		"values": [
			{
				"timestamp": 1736485200000,
				"value": 94.8765
			},
			{
				"timestamp": 1736312400000,
				"value": 96.2341
			}
		]
	},
	"status": "OK",
	"request_id": "opt-ema-req-001",
	"next_url": "https://api.massive.com/v1/indicators/ema/O:AAPL250117C00150000?cursor=def456"
}`

const optionsRSIJSON = `{
	"results": {
		"underlying": {
			"url": "https://api.polygon.io/v2/aggs/ticker/O:AAPL250117C00150000/range/1/day/1729483200000/1736553600000?limit=68&sort=desc"
		},
		"values": [
			{
				"timestamp": 1736485200000,
				"value": 55.1234
			},
			{
				"timestamp": 1736312400000,
				"value": 52.6789
			},
			{
				"timestamp": 1736226000000,
				"value": 51.3456
			}
		]
	},
	"status": "OK",
	"request_id": "opt-rsi-req-001"
}`

const optionsMACDJSON = `{
	"results": {
		"underlying": {
			"url": "https://api.polygon.io/v2/aggs/ticker/O:AAPL250117C00150000/range/1/day/1724212800000/1736553600000?limit=122&sort=desc"
		},
		"values": [
			{
				"timestamp": 1736485200000,
				"value": 1.2345,
				"signal": 0.9876,
				"histogram": 0.2469
			},
			{
				"timestamp": 1736312400000,
				"value": 1.5678,
				"signal": 1.1234,
				"histogram": 0.4444
			}
		]
	},
	"status": "OK",
	"request_id": "opt-macd-req-001"
}`

// TestGetOptionsSMA verifies that GetOptionsSMA correctly parses the API
// response and returns the expected SMA indicator values for an options contract.
func TestGetOptionsSMA(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/sma/O:AAPL250117C00150000": optionsSMAJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := IndicatorParams{
		TimestampGTE: "2025-01-06",
		TimestampLTE: "2025-01-10",
		Timespan:     "day",
		Window:       "10",
		SeriesType:   "close",
		Limit:        "5",
	}

	result, err := client.GetOptionsSMA("O:AAPL250117C00150000", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "opt-sma-req-001" {
		t.Errorf("expected request_id opt-sma-req-001, got %s", result.RequestID)
	}

	if len(result.Results.Values) != 3 {
		t.Fatalf("expected 3 values, got %d", len(result.Results.Values))
	}

	first := result.Results.Values[0]
	if first.Timestamp != 1736485200000 {
		t.Errorf("expected timestamp 1736485200000, got %d", first.Timestamp)
	}

	if first.Value != 95.432 {
		t.Errorf("expected value 95.432, got %f", first.Value)
	}

	second := result.Results.Values[1]
	if second.Value != 96.781 {
		t.Errorf("expected value 96.781, got %f", second.Value)
	}
}

// TestGetOptionsSMARequestPath verifies that GetOptionsSMA constructs the
// correct API path with the options contract ticker symbol.
func TestGetOptionsSMARequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsSMAJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsSMA("O:TSLA250221P00200000", IndicatorParams{})

	if receivedPath != "/v1/indicators/sma/O:TSLA250221P00200000" {
		t.Errorf("expected path /v1/indicators/sma/O:TSLA250221P00200000, got %s", receivedPath)
	}
}

// TestGetOptionsSMAQueryParams verifies that all SMA query parameters are
// correctly sent to the API endpoint for options indicators.
func TestGetOptionsSMAQueryParams(t *testing.T) {
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
		w.Write([]byte(optionsSMAJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsSMA("O:AAPL250117C00150000", IndicatorParams{
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

// TestGetOptionsSMAUnderlyingURL verifies that the underlying aggregates URL
// is correctly parsed from the options SMA response.
func TestGetOptionsSMAUnderlyingURL(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/sma/O:AAPL250117C00150000": optionsSMAJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsSMA("O:AAPL250117C00150000", IndicatorParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedURL := "https://api.polygon.io/v2/aggs/ticker/O:AAPL250117C00150000/range/1/day/1731042000000/1736553600000?limit=50&sort=desc"
	if result.Results.Underlying.URL != expectedURL {
		t.Errorf("expected underlying URL %s, got %s", expectedURL, result.Results.Underlying.URL)
	}
}

// TestGetOptionsSMAAPIError verifies that GetOptionsSMA returns an error when
// the API responds with a non-200 status code.
func TestGetOptionsSMAAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Data not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetOptionsSMA("O:INVALID", IndicatorParams{})
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// TestGetOptionsSMAEmptyValues verifies that GetOptionsSMA handles an empty
// values array without error.
func TestGetOptionsSMAEmptyValues(t *testing.T) {
	emptyJSON := `{"results":{"underlying":{"url":""},"values":[]},"status":"OK","request_id":"abc"}`
	server := mockServer(t, map[string]string{
		"/v1/indicators/sma/O:AAPL250117C00150000": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsSMA("O:AAPL250117C00150000", IndicatorParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results.Values) != 0 {
		t.Errorf("expected 0 values, got %d", len(result.Results.Values))
	}
}

// TestGetOptionsEMA verifies that GetOptionsEMA correctly parses the API
// response and returns the expected EMA indicator values for an options contract.
func TestGetOptionsEMA(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/ema/O:AAPL250117C00150000": optionsEMAJSON,
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

	result, err := client.GetOptionsEMA("O:AAPL250117C00150000", params)
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
	if first.Timestamp != 1736485200000 {
		t.Errorf("expected timestamp 1736485200000, got %d", first.Timestamp)
	}

	if first.Value != 94.8765 {
		t.Errorf("expected value 94.8765, got %f", first.Value)
	}
}

// TestGetOptionsEMARequestPath verifies that GetOptionsEMA constructs the
// correct API path with the options contract ticker symbol.
func TestGetOptionsEMARequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsEMAJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsEMA("O:TSLA250221P00200000", IndicatorParams{})

	if receivedPath != "/v1/indicators/ema/O:TSLA250221P00200000" {
		t.Errorf("expected path /v1/indicators/ema/O:TSLA250221P00200000, got %s", receivedPath)
	}
}

// TestGetOptionsEMANextURL verifies that the pagination next_url field is
// correctly parsed from the options EMA response.
func TestGetOptionsEMANextURL(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/ema/O:AAPL250117C00150000": optionsEMAJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsEMA("O:AAPL250117C00150000", IndicatorParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.NextURL != "https://api.massive.com/v1/indicators/ema/O:AAPL250117C00150000?cursor=def456" {
		t.Errorf("expected next_url with cursor, got %s", result.NextURL)
	}
}

// TestGetOptionsEMAAPIError verifies that GetOptionsEMA returns an error
// when the API responds with a non-200 status code.
func TestGetOptionsEMAAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetOptionsEMA("O:AAPL250117C00150000", IndicatorParams{})
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

// TestGetOptionsRSI verifies that GetOptionsRSI correctly parses the API
// response and returns the expected RSI indicator values for an options contract.
func TestGetOptionsRSI(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/rsi/O:AAPL250117C00150000": optionsRSIJSON,
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

	result, err := client.GetOptionsRSI("O:AAPL250117C00150000", params)
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
	if first.Timestamp != 1736485200000 {
		t.Errorf("expected timestamp 1736485200000, got %d", first.Timestamp)
	}

	if first.Value != 55.1234 {
		t.Errorf("expected value 55.1234, got %f", first.Value)
	}

	second := result.Results.Values[1]
	if second.Value != 52.6789 {
		t.Errorf("expected value 52.6789, got %f", second.Value)
	}
}

// TestGetOptionsRSIRequestPath verifies that GetOptionsRSI constructs the
// correct API path with the options contract ticker symbol.
func TestGetOptionsRSIRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsRSIJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsRSI("O:SPY250321C00500000", IndicatorParams{})

	if receivedPath != "/v1/indicators/rsi/O:SPY250321C00500000" {
		t.Errorf("expected path /v1/indicators/rsi/O:SPY250321C00500000, got %s", receivedPath)
	}
}

// TestGetOptionsRSIQueryParams verifies that RSI-specific query parameters
// including window and series_type are correctly sent for options indicators.
func TestGetOptionsRSIQueryParams(t *testing.T) {
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
		w.Write([]byte(optionsRSIJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsRSI("O:AAPL250117C00150000", IndicatorParams{
		Window:     "14",
		SeriesType: "close",
		Timespan:   "day",
		Adjusted:   "false",
	})
}

// TestGetOptionsRSIAPIError verifies that GetOptionsRSI returns an error
// when the API responds with a non-200 status code.
func TestGetOptionsRSIAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Ticker not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetOptionsRSI("O:INVALID", IndicatorParams{})
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// TestGetOptionsMACD verifies that GetOptionsMACD correctly parses the API
// response and returns the expected MACD values including signal and histogram
// for an options contract.
func TestGetOptionsMACD(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/macd/O:AAPL250117C00150000": optionsMACDJSON,
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

	result, err := client.GetOptionsMACD("O:AAPL250117C00150000", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "opt-macd-req-001" {
		t.Errorf("expected request_id opt-macd-req-001, got %s", result.RequestID)
	}

	if len(result.Results.Values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(result.Results.Values))
	}

	first := result.Results.Values[0]
	if first.Timestamp != 1736485200000 {
		t.Errorf("expected timestamp 1736485200000, got %d", first.Timestamp)
	}

	if first.Value != 1.2345 {
		t.Errorf("expected value 1.2345, got %f", first.Value)
	}

	if first.Signal != 0.9876 {
		t.Errorf("expected signal 0.9876, got %f", first.Signal)
	}

	if first.Histogram != 0.2469 {
		t.Errorf("expected histogram 0.2469, got %f", first.Histogram)
	}
}

// TestGetOptionsMACDRequestPath verifies that GetOptionsMACD constructs the
// correct API path with the options contract ticker symbol.
func TestGetOptionsMACDRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsMACDJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsMACD("O:NVDA250321C00800000", MACDParams{})

	if receivedPath != "/v1/indicators/macd/O:NVDA250321C00800000" {
		t.Errorf("expected path /v1/indicators/macd/O:NVDA250321C00800000, got %s", receivedPath)
	}
}

// TestGetOptionsMACDQueryParams verifies that MACD-specific query parameters
// including short_window, long_window, and signal_window are correctly sent
// for options indicators.
func TestGetOptionsMACDQueryParams(t *testing.T) {
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
		w.Write([]byte(optionsMACDJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsMACD("O:AAPL250117C00150000", MACDParams{
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

// TestGetOptionsMACDSecondValue verifies that the second MACD value in the
// options response is correctly parsed with its own distinct values.
func TestGetOptionsMACDSecondValue(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/macd/O:AAPL250117C00150000": optionsMACDJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsMACD("O:AAPL250117C00150000", MACDParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := result.Results.Values[1]
	if second.Value != 1.5678 {
		t.Errorf("expected value 1.5678, got %f", second.Value)
	}

	if second.Signal != 1.1234 {
		t.Errorf("expected signal 1.1234, got %f", second.Signal)
	}

	if second.Histogram != 0.4444 {
		t.Errorf("expected histogram 0.4444, got %f", second.Histogram)
	}
}

// TestGetOptionsMACDAPIError verifies that GetOptionsMACD returns an error
// when the API responds with a non-200 status code.
func TestGetOptionsMACDAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"ERROR","message":"Unauthorized"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetOptionsMACD("O:AAPL250117C00150000", MACDParams{})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetOptionsMACDUnderlyingURL verifies that the underlying aggregates URL
// is correctly parsed from the options MACD response.
func TestGetOptionsMACDUnderlyingURL(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/macd/O:AAPL250117C00150000": optionsMACDJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsMACD("O:AAPL250117C00150000", MACDParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedURL := "https://api.polygon.io/v2/aggs/ticker/O:AAPL250117C00150000/range/1/day/1724212800000/1736553600000?limit=122&sort=desc"
	if result.Results.Underlying.URL != expectedURL {
		t.Errorf("expected underlying URL %s, got %s", expectedURL, result.Results.Underlying.URL)
	}
}

// TestGetOptionsMACDEmptyValues verifies that GetOptionsMACD handles an empty
// values array without error.
func TestGetOptionsMACDEmptyValues(t *testing.T) {
	emptyJSON := `{"results":{"underlying":{"url":""},"values":[]},"status":"OK","request_id":"abc"}`
	server := mockServer(t, map[string]string{
		"/v1/indicators/macd/O:AAPL250117C00150000": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsMACD("O:AAPL250117C00150000", MACDParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results.Values) != 0 {
		t.Errorf("expected 0 values, got %d", len(result.Results.Values))
	}
}
