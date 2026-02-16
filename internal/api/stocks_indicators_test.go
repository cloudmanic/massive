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

const smaJSON = `{
	"results": {
		"underlying": {
			"url": "https://api.polygon.io/v2/aggs/ticker/AAPL/range/1/day/1731042000000/1736553600000?limit=52&sort=desc"
		},
		"values": [
			{
				"timestamp": 1736485200000,
				"value": 247.12
			},
			{
				"timestamp": 1736312400000,
				"value": 249.255
			},
			{
				"timestamp": 1736226000000,
				"value": 250.512
			}
		]
	},
	"status": "OK",
	"request_id": "dc8b912001a7547c540f805ea80c1463"
}`

const emaJSON = `{
	"results": {
		"underlying": {
			"url": "https://api.polygon.io/v2/aggs/ticker/AAPL/range/1/day/1731042000000/1736553600000?limit=50&sort=desc"
		},
		"values": [
			{
				"timestamp": 1736485200000,
				"value": 244.9258
			},
			{
				"timestamp": 1736312400000,
				"value": 246.7204
			}
		]
	},
	"status": "OK",
	"request_id": "2ae43aa6603214c666f160063ed699f7",
	"next_url": "https://api.massive.com/v1/indicators/ema/AAPL?cursor=abc123"
}`

const rsiJSON = `{
	"results": {
		"underlying": {
			"url": "https://api.polygon.io/v2/aggs/ticker/AAPL/range/1/day/1729483200000/1736553600000?limit=68&sort=desc"
		},
		"values": [
			{
				"timestamp": 1736485200000,
				"value": 36.6817
			},
			{
				"timestamp": 1736312400000,
				"value": 44.2105
			},
			{
				"timestamp": 1736226000000,
				"value": 43.3054
			}
		]
	},
	"status": "OK",
	"request_id": "e00a03214a982c65f820f020e41844ff"
}`

const macdJSON = `{
	"results": {
		"underlying": {
			"url": "https://api.polygon.io/v2/aggs/ticker/AAPL/range/1/day/1724212800000/1736553600000?limit=122&sort=desc"
		},
		"values": [
			{
				"timestamp": 1736485200000,
				"value": 0.3681,
				"signal": 2.7514,
				"histogram": -2.3833
			},
			{
				"timestamp": 1736312400000,
				"value": 1.2801,
				"signal": 3.3472,
				"histogram": -2.0671
			}
		]
	},
	"status": "OK",
	"request_id": "54f2ddc2d819c98580b3dd90bcd881d2"
}`

// TestGetSMA verifies that GetSMA correctly parses the API response
// and returns the expected SMA indicator values for AAPL.
func TestGetSMA(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/sma/AAPL": smaJSON,
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

	result, err := client.GetSMA("AAPL", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "dc8b912001a7547c540f805ea80c1463" {
		t.Errorf("expected request_id dc8b912001a7547c540f805ea80c1463, got %s", result.RequestID)
	}

	if len(result.Results.Values) != 3 {
		t.Fatalf("expected 3 values, got %d", len(result.Results.Values))
	}

	first := result.Results.Values[0]
	if first.Timestamp != 1736485200000 {
		t.Errorf("expected timestamp 1736485200000, got %d", first.Timestamp)
	}

	if first.Value != 247.12 {
		t.Errorf("expected value 247.12, got %f", first.Value)
	}

	second := result.Results.Values[1]
	if second.Value != 249.255 {
		t.Errorf("expected value 249.255, got %f", second.Value)
	}
}

// TestGetSMARequestPath verifies that GetSMA constructs the correct
// API path with the ticker symbol.
func TestGetSMARequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(smaJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetSMA("MSFT", IndicatorParams{})

	if receivedPath != "/v1/indicators/sma/MSFT" {
		t.Errorf("expected path /v1/indicators/sma/MSFT, got %s", receivedPath)
	}
}

// TestGetSMAQueryParams verifies that all SMA query parameters are
// correctly sent to the API endpoint.
func TestGetSMAQueryParams(t *testing.T) {
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
		w.Write([]byte(smaJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetSMA("AAPL", IndicatorParams{
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

// TestGetSMAUnderlyingURL verifies that the underlying aggregates URL
// is correctly parsed from the SMA response.
func TestGetSMAUnderlyingURL(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/sma/AAPL": smaJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetSMA("AAPL", IndicatorParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedURL := "https://api.polygon.io/v2/aggs/ticker/AAPL/range/1/day/1731042000000/1736553600000?limit=52&sort=desc"
	if result.Results.Underlying.URL != expectedURL {
		t.Errorf("expected underlying URL %s, got %s", expectedURL, result.Results.Underlying.URL)
	}
}

// TestGetSMAAPIError verifies that GetSMA returns an error when the
// API responds with a non-200 status code.
func TestGetSMAAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Data not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetSMA("INVALID", IndicatorParams{})
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// TestGetEMA verifies that GetEMA correctly parses the API response
// and returns the expected EMA indicator values for AAPL.
func TestGetEMA(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/ema/AAPL": emaJSON,
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

	result, err := client.GetEMA("AAPL", params)
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

	if first.Value != 244.9258 {
		t.Errorf("expected value 244.9258, got %f", first.Value)
	}
}

// TestGetEMARequestPath verifies that GetEMA constructs the correct
// API path with the ticker symbol.
func TestGetEMARequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(emaJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetEMA("TSLA", IndicatorParams{})

	if receivedPath != "/v1/indicators/ema/TSLA" {
		t.Errorf("expected path /v1/indicators/ema/TSLA, got %s", receivedPath)
	}
}

// TestGetEMANextURL verifies that the pagination next_url field is
// correctly parsed from the EMA response.
func TestGetEMANextURL(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/ema/AAPL": emaJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetEMA("AAPL", IndicatorParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.NextURL != "https://api.massive.com/v1/indicators/ema/AAPL?cursor=abc123" {
		t.Errorf("expected next_url with cursor, got %s", result.NextURL)
	}
}

// TestGetEMAAPIError verifies that GetEMA returns an error when the
// API responds with a non-200 status code.
func TestGetEMAAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetEMA("AAPL", IndicatorParams{})
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

// TestGetRSI verifies that GetRSI correctly parses the API response
// and returns the expected RSI indicator values for AAPL.
func TestGetRSI(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/rsi/AAPL": rsiJSON,
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

	result, err := client.GetRSI("AAPL", params)
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

	if first.Value != 36.6817 {
		t.Errorf("expected value 36.6817, got %f", first.Value)
	}

	second := result.Results.Values[1]
	if second.Value != 44.2105 {
		t.Errorf("expected value 44.2105, got %f", second.Value)
	}
}

// TestGetRSIRequestPath verifies that GetRSI constructs the correct
// API path with the ticker symbol.
func TestGetRSIRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(rsiJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetRSI("GOOG", IndicatorParams{})

	if receivedPath != "/v1/indicators/rsi/GOOG" {
		t.Errorf("expected path /v1/indicators/rsi/GOOG, got %s", receivedPath)
	}
}

// TestGetRSIQueryParams verifies that RSI-specific query parameters
// including window and series_type are correctly sent.
func TestGetRSIQueryParams(t *testing.T) {
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
		w.Write([]byte(rsiJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetRSI("AAPL", IndicatorParams{
		Window:     "14",
		SeriesType: "close",
		Timespan:   "day",
		Adjusted:   "false",
	})
}

// TestGetRSIAPIError verifies that GetRSI returns an error when the
// API responds with a non-200 status code.
func TestGetRSIAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Ticker not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetRSI("INVALID", IndicatorParams{})
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// TestGetMACD verifies that GetMACD correctly parses the API response
// and returns the expected MACD values including signal and histogram.
func TestGetMACD(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/macd/AAPL": macdJSON,
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

	result, err := client.GetMACD("AAPL", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "54f2ddc2d819c98580b3dd90bcd881d2" {
		t.Errorf("expected request_id 54f2ddc2d819c98580b3dd90bcd881d2, got %s", result.RequestID)
	}

	if len(result.Results.Values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(result.Results.Values))
	}

	first := result.Results.Values[0]
	if first.Timestamp != 1736485200000 {
		t.Errorf("expected timestamp 1736485200000, got %d", first.Timestamp)
	}

	if first.Value != 0.3681 {
		t.Errorf("expected value 0.3681, got %f", first.Value)
	}

	if first.Signal != 2.7514 {
		t.Errorf("expected signal 2.7514, got %f", first.Signal)
	}

	if first.Histogram != -2.3833 {
		t.Errorf("expected histogram -2.3833, got %f", first.Histogram)
	}
}

// TestGetMACDRequestPath verifies that GetMACD constructs the correct
// API path with the ticker symbol.
func TestGetMACDRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(macdJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetMACD("NVDA", MACDParams{})

	if receivedPath != "/v1/indicators/macd/NVDA" {
		t.Errorf("expected path /v1/indicators/macd/NVDA, got %s", receivedPath)
	}
}

// TestGetMACDQueryParams verifies that MACD-specific query parameters
// including short_window, long_window, and signal_window are correctly sent.
func TestGetMACDQueryParams(t *testing.T) {
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
		w.Write([]byte(macdJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetMACD("AAPL", MACDParams{
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

// TestGetMACDSecondValue verifies that the second MACD value in the
// response is correctly parsed with its own distinct values.
func TestGetMACDSecondValue(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/macd/AAPL": macdJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetMACD("AAPL", MACDParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := result.Results.Values[1]
	if second.Value != 1.2801 {
		t.Errorf("expected value 1.2801, got %f", second.Value)
	}

	if second.Signal != 3.3472 {
		t.Errorf("expected signal 3.3472, got %f", second.Signal)
	}

	if second.Histogram != -2.0671 {
		t.Errorf("expected histogram -2.0671, got %f", second.Histogram)
	}
}

// TestGetMACDAPIError verifies that GetMACD returns an error when the
// API responds with a non-200 status code.
func TestGetMACDAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"ERROR","message":"Unauthorized"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetMACD("AAPL", MACDParams{})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetMACDUnderlyingURL verifies that the underlying aggregates URL
// is correctly parsed from the MACD response.
func TestGetMACDUnderlyingURL(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/indicators/macd/AAPL": macdJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetMACD("AAPL", MACDParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedURL := "https://api.polygon.io/v2/aggs/ticker/AAPL/range/1/day/1724212800000/1736553600000?limit=122&sort=desc"
	if result.Results.Underlying.URL != expectedURL {
		t.Errorf("expected underlying URL %s, got %s", expectedURL, result.Results.Underlying.URL)
	}
}

// TestGetSMAEmptyValues verifies that GetSMA handles an empty values
// array without error.
func TestGetSMAEmptyValues(t *testing.T) {
	emptyJSON := `{"results":{"underlying":{"url":""},"values":[]},"status":"OK","request_id":"abc"}`
	server := mockServer(t, map[string]string{
		"/v1/indicators/sma/AAPL": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetSMA("AAPL", IndicatorParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results.Values) != 0 {
		t.Errorf("expected 0 values, got %d", len(result.Results.Values))
	}
}

// TestGetMACDEmptyValues verifies that GetMACD handles an empty values
// array without error.
func TestGetMACDEmptyValues(t *testing.T) {
	emptyJSON := `{"results":{"underlying":{"url":""},"values":[]},"status":"OK","request_id":"abc"}`
	server := mockServer(t, map[string]string{
		"/v1/indicators/macd/AAPL": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetMACD("AAPL", MACDParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results.Values) != 0 {
		t.Errorf("expected 0 values, got %d", len(result.Results.Values))
	}
}
