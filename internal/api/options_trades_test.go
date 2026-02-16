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

const optionsTradesJSON = `{
	"status": "OK",
	"request_id": "abc123opttrades",
	"next_url": "https://api.massive.com/v3/trades/O:AAPL250221C00230000?cursor=YXA9MQ",
	"results": [
		{
			"conditions": [209],
			"correction": 0,
			"exchange": 46,
			"participant_timestamp": 1401715883806000000,
			"price": 6.91,
			"sequence_number": 987654,
			"sip_timestamp": 1401715883806000000,
			"size": 1
		},
		{
			"conditions": [227],
			"correction": 0,
			"exchange": 312,
			"participant_timestamp": 1401715884000000000,
			"price": 7.05,
			"sequence_number": 987655,
			"sip_timestamp": 1401715884100000000,
			"size": 5
		}
	]
}`

const optionsLastTradeJSON = `{
	"status": "OK",
	"request_id": "abc123optlasttrade",
	"results": {
		"T": "O:TSLA210903C00700000",
		"c": [227],
		"e": 0,
		"f": 1617901342969796400,
		"i": "52983525029461",
		"p": 115.55,
		"q": 1325541950,
		"r": 202,
		"s": 25,
		"t": 1617901342969834000,
		"x": 312,
		"y": 1617901342969834000,
		"z": 3
	}
}`

const optionsQuotesJSON = `{
	"status": "OK",
	"request_id": "abc123optquotes",
	"next_url": "https://api.massive.com/v3/quotes/O:AAPL250221C00230000?cursor=YXA9MQ",
	"results": [
		{
			"ask_exchange": 302,
			"ask_price": 7.10,
			"ask_size": 20,
			"bid_exchange": 316,
			"bid_price": 6.80,
			"bid_size": 15,
			"sequence_number": 5554321,
			"sip_timestamp": 1401715883806000000
		},
		{
			"ask_exchange": 312,
			"ask_price": 7.15,
			"ask_size": 10,
			"bid_exchange": 302,
			"bid_price": 6.85,
			"bid_size": 25,
			"sequence_number": 5554322,
			"sip_timestamp": 1401715884100000000
		}
	]
}`

const optionsLastQuoteJSON = `{
	"status": "OK",
	"request_id": "abc123optlastquote",
	"results": {
		"T": "O:TSLA210903C00700000",
		"P": 7.10,
		"S": 20,
		"X": 302,
		"c": [1],
		"f": 1617901342969796400,
		"i": [0],
		"p": 6.80,
		"q": 5554321,
		"s": 15,
		"t": 1617901342969834000,
		"x": 316,
		"y": 1617901342969834000,
		"z": 3
	}
}`

// TestGetOptionsTrades verifies that GetOptionsTrades correctly parses the API
// response and returns the expected tick-level trade data for an options contract.
func TestGetOptionsTrades(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/trades/O:AAPL250221C00230000": optionsTradesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := OptionsTradesParams{
		Timestamp: "2025-01-06",
		Limit:     "2",
	}

	result, err := client.GetOptionsTrades("O:AAPL250221C00230000", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "abc123opttrades" {
		t.Errorf("expected request_id abc123opttrades, got %s", result.RequestID)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 trades, got %d", len(result.Results))
	}

	trade := result.Results[0]
	if trade.Price != 6.91 {
		t.Errorf("expected price 6.91, got %f", trade.Price)
	}

	if trade.Size != 1 {
		t.Errorf("expected size 1, got %f", trade.Size)
	}

	if trade.Exchange != 46 {
		t.Errorf("expected exchange 46, got %d", trade.Exchange)
	}

	if trade.SequenceNumber != 987654 {
		t.Errorf("expected sequence_number 987654, got %d", trade.SequenceNumber)
	}

	if len(trade.Conditions) != 1 {
		t.Errorf("expected 1 condition, got %d", len(trade.Conditions))
	}

	if trade.Conditions[0] != 209 {
		t.Errorf("expected first condition 209, got %d", trade.Conditions[0])
	}
}

// TestGetOptionsTradesSecondResult verifies that the second trade record in
// the response is correctly parsed with its own distinct values.
func TestGetOptionsTradesSecondResult(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/trades/O:AAPL250221C00230000": optionsTradesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsTrades("O:AAPL250221C00230000", OptionsTradesParams{Timestamp: "2025-01-06"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	trade := result.Results[1]
	if trade.Price != 7.05 {
		t.Errorf("expected price 7.05, got %f", trade.Price)
	}

	if trade.Size != 5 {
		t.Errorf("expected size 5, got %f", trade.Size)
	}

	if trade.Exchange != 312 {
		t.Errorf("expected exchange 312, got %d", trade.Exchange)
	}
}

// TestGetOptionsTradesRequestPath verifies that GetOptionsTrades constructs the
// correct API path with the options ticker symbol.
func TestGetOptionsTradesRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsTradesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsTrades("O:SPY250221P00500000", OptionsTradesParams{})

	if receivedPath != "/v3/trades/O:SPY250221P00500000" {
		t.Errorf("expected path /v3/trades/O:SPY250221P00500000, got %s", receivedPath)
	}
}

// TestGetOptionsTradesQueryParams verifies that all query parameters including
// timestamp range filters are correctly sent to the API.
func TestGetOptionsTradesQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("timestamp") != "2025-01-06" {
			t.Errorf("expected timestamp=2025-01-06, got %s", q.Get("timestamp"))
		}
		if q.Get("order") != "desc" {
			t.Errorf("expected order=desc, got %s", q.Get("order"))
		}
		if q.Get("limit") != "100" {
			t.Errorf("expected limit=100, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "timestamp" {
			t.Errorf("expected sort=timestamp, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsTradesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsTrades("O:AAPL250221C00230000", OptionsTradesParams{
		Timestamp: "2025-01-06",
		Order:     "desc",
		Limit:     "100",
		Sort:      "timestamp",
	})
}

// TestGetOptionsTradesTimestampRangeParams verifies that the timestamp range
// query parameters (gte, gt, lte, lt) are correctly sent to the API.
func TestGetOptionsTradesTimestampRangeParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("timestamp.gte") != "2025-01-06" {
			t.Errorf("expected timestamp.gte=2025-01-06, got %s", q.Get("timestamp.gte"))
		}
		if q.Get("timestamp.lte") != "2025-01-08" {
			t.Errorf("expected timestamp.lte=2025-01-08, got %s", q.Get("timestamp.lte"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsTradesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsTrades("O:AAPL250221C00230000", OptionsTradesParams{
		TimestampGte: "2025-01-06",
		TimestampLte: "2025-01-08",
	})
}

// TestGetOptionsTradesAPIError verifies that GetOptionsTrades returns an error
// when the API responds with a non-200 status code.
func TestGetOptionsTradesAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Data not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetOptionsTrades("O:INVALID", OptionsTradesParams{})
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// TestGetOptionsLastTrade verifies that GetOptionsLastTrade correctly parses
// the API response and returns the expected last trade data for an options contract.
func TestGetOptionsLastTrade(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/last/trade/O:TSLA210903C00700000": optionsLastTradeJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsLastTrade("O:TSLA210903C00700000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "abc123optlasttrade" {
		t.Errorf("expected request_id abc123optlasttrade, got %s", result.RequestID)
	}

	trade := result.Results
	if trade.Ticker != "O:TSLA210903C00700000" {
		t.Errorf("expected ticker O:TSLA210903C00700000, got %s", trade.Ticker)
	}

	if trade.Price != 115.55 {
		t.Errorf("expected price 115.55, got %f", trade.Price)
	}

	if trade.Size != 25 {
		t.Errorf("expected size 25, got %f", trade.Size)
	}

	if trade.Exchange != 312 {
		t.Errorf("expected exchange 312, got %d", trade.Exchange)
	}

	if trade.ID != "52983525029461" {
		t.Errorf("expected id 52983525029461, got %s", trade.ID)
	}

	if trade.Tape != 3 {
		t.Errorf("expected tape 3, got %d", trade.Tape)
	}

	if trade.SequenceNumber != 1325541950 {
		t.Errorf("expected sequence_number 1325541950, got %d", trade.SequenceNumber)
	}

	if trade.Correction != 0 {
		t.Errorf("expected correction 0, got %d", trade.Correction)
	}

	if len(trade.Conditions) != 1 {
		t.Errorf("expected 1 condition, got %d", len(trade.Conditions))
	}

	if trade.Conditions[0] != 227 {
		t.Errorf("expected first condition 227, got %d", trade.Conditions[0])
	}
}

// TestGetOptionsLastTradeRequestPath verifies that GetOptionsLastTrade constructs
// the correct API path with the options ticker symbol.
func TestGetOptionsLastTradeRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsLastTradeJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsLastTrade("O:SPY250221P00500000")

	if receivedPath != "/v2/last/trade/O:SPY250221P00500000" {
		t.Errorf("expected path /v2/last/trade/O:SPY250221P00500000, got %s", receivedPath)
	}
}

// TestGetOptionsLastTradeAPIError verifies that GetOptionsLastTrade returns an
// error when the API responds with a non-200 status code.
func TestGetOptionsLastTradeAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"status":"NOT_AUTHORIZED","message":"Not authorized."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetOptionsLastTrade("O:AAPL250221C00230000")
	if err == nil {
		t.Fatal("expected error for 401 response, got nil")
	}
}

// TestGetOptionsQuotes verifies that GetOptionsQuotes correctly parses the API
// response and returns the expected tick-level NBBO quote data for an options contract.
func TestGetOptionsQuotes(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/quotes/O:AAPL250221C00230000": optionsQuotesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := OptionsQuotesParams{
		Timestamp: "2025-01-06",
		Limit:     "2",
	}

	result, err := client.GetOptionsQuotes("O:AAPL250221C00230000", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "abc123optquotes" {
		t.Errorf("expected request_id abc123optquotes, got %s", result.RequestID)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 quotes, got %d", len(result.Results))
	}

	quote := result.Results[0]
	if quote.AskPrice != 7.10 {
		t.Errorf("expected ask_price 7.10, got %f", quote.AskPrice)
	}

	if quote.AskSize != 20 {
		t.Errorf("expected ask_size 20, got %f", quote.AskSize)
	}

	if quote.AskExchange != 302 {
		t.Errorf("expected ask_exchange 302, got %d", quote.AskExchange)
	}

	if quote.BidPrice != 6.80 {
		t.Errorf("expected bid_price 6.80, got %f", quote.BidPrice)
	}

	if quote.BidSize != 15 {
		t.Errorf("expected bid_size 15, got %f", quote.BidSize)
	}

	if quote.BidExchange != 316 {
		t.Errorf("expected bid_exchange 316, got %d", quote.BidExchange)
	}

	if quote.SequenceNumber != 5554321 {
		t.Errorf("expected sequence_number 5554321, got %d", quote.SequenceNumber)
	}
}

// TestGetOptionsQuotesSecondResult verifies that the second quote record in
// the response is correctly parsed with its own distinct values.
func TestGetOptionsQuotesSecondResult(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/quotes/O:AAPL250221C00230000": optionsQuotesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsQuotes("O:AAPL250221C00230000", OptionsQuotesParams{Timestamp: "2025-01-06"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	quote := result.Results[1]
	if quote.AskPrice != 7.15 {
		t.Errorf("expected ask_price 7.15, got %f", quote.AskPrice)
	}

	if quote.AskSize != 10 {
		t.Errorf("expected ask_size 10, got %f", quote.AskSize)
	}

	if quote.BidPrice != 6.85 {
		t.Errorf("expected bid_price 6.85, got %f", quote.BidPrice)
	}

	if quote.BidSize != 25 {
		t.Errorf("expected bid_size 25, got %f", quote.BidSize)
	}
}

// TestGetOptionsQuotesRequestPath verifies that GetOptionsQuotes constructs the
// correct API path with the options ticker symbol.
func TestGetOptionsQuotesRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsQuotesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsQuotes("O:SPY250221P00500000", OptionsQuotesParams{})

	if receivedPath != "/v3/quotes/O:SPY250221P00500000" {
		t.Errorf("expected path /v3/quotes/O:SPY250221P00500000, got %s", receivedPath)
	}
}

// TestGetOptionsQuotesQueryParams verifies that all query parameters including
// timestamp range filters are correctly sent to the API.
func TestGetOptionsQuotesQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("timestamp") != "2025-01-06" {
			t.Errorf("expected timestamp=2025-01-06, got %s", q.Get("timestamp"))
		}
		if q.Get("order") != "asc" {
			t.Errorf("expected order=asc, got %s", q.Get("order"))
		}
		if q.Get("limit") != "500" {
			t.Errorf("expected limit=500, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "timestamp" {
			t.Errorf("expected sort=timestamp, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsQuotesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsQuotes("O:AAPL250221C00230000", OptionsQuotesParams{
		Timestamp: "2025-01-06",
		Order:     "asc",
		Limit:     "500",
		Sort:      "timestamp",
	})
}

// TestGetOptionsQuotesTimestampRangeParams verifies that the timestamp range
// query parameters (gte, gt, lte, lt) are correctly sent to the API.
func TestGetOptionsQuotesTimestampRangeParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("timestamp.gt") != "2025-01-06" {
			t.Errorf("expected timestamp.gt=2025-01-06, got %s", q.Get("timestamp.gt"))
		}
		if q.Get("timestamp.lt") != "2025-01-08" {
			t.Errorf("expected timestamp.lt=2025-01-08, got %s", q.Get("timestamp.lt"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsQuotesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsQuotes("O:AAPL250221C00230000", OptionsQuotesParams{
		TimestampGt: "2025-01-06",
		TimestampLt: "2025-01-08",
	})
}

// TestGetOptionsQuotesAPIError verifies that GetOptionsQuotes returns an error
// when the API responds with a non-200 status code.
func TestGetOptionsQuotesAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Data not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetOptionsQuotes("O:INVALID", OptionsQuotesParams{})
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// TestGetOptionsLastQuote verifies that GetOptionsLastQuote correctly parses the
// API response and returns the expected last NBBO quote data for an options contract.
func TestGetOptionsLastQuote(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/last/nbbo/O:TSLA210903C00700000": optionsLastQuoteJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsLastQuote("O:TSLA210903C00700000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "abc123optlastquote" {
		t.Errorf("expected request_id abc123optlastquote, got %s", result.RequestID)
	}

	quote := result.Results
	if quote.Ticker != "O:TSLA210903C00700000" {
		t.Errorf("expected ticker O:TSLA210903C00700000, got %s", quote.Ticker)
	}

	if quote.AskPrice != 7.10 {
		t.Errorf("expected ask_price 7.10, got %f", quote.AskPrice)
	}

	if quote.AskSize != 20 {
		t.Errorf("expected ask_size 20, got %d", quote.AskSize)
	}

	if quote.AskExchange != 302 {
		t.Errorf("expected ask_exchange 302, got %d", quote.AskExchange)
	}

	if quote.BidPrice != 6.80 {
		t.Errorf("expected bid_price 6.80, got %f", quote.BidPrice)
	}

	if quote.BidSize != 15 {
		t.Errorf("expected bid_size 15, got %d", quote.BidSize)
	}

	if quote.BidExchange != 316 {
		t.Errorf("expected bid_exchange 316, got %d", quote.BidExchange)
	}

	if quote.Tape != 3 {
		t.Errorf("expected tape 3, got %d", quote.Tape)
	}

	if quote.SequenceNumber != 5554321 {
		t.Errorf("expected sequence_number 5554321, got %d", quote.SequenceNumber)
	}
}

// TestGetOptionsLastQuoteRequestPath verifies that GetOptionsLastQuote constructs
// the correct API path with the options ticker symbol.
func TestGetOptionsLastQuoteRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsLastQuoteJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsLastQuote("O:SPY250221P00500000")

	if receivedPath != "/v2/last/nbbo/O:SPY250221P00500000" {
		t.Errorf("expected path /v2/last/nbbo/O:SPY250221P00500000, got %s", receivedPath)
	}
}

// TestGetOptionsLastQuoteAPIError verifies that GetOptionsLastQuote returns an
// error when the API responds with a non-200 status code.
func TestGetOptionsLastQuoteAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"status":"NOT_AUTHORIZED","message":"Not authorized."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetOptionsLastQuote("O:AAPL250221C00230000")
	if err == nil {
		t.Fatal("expected error for 401 response, got nil")
	}
}
