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

const tradesJSON = `{
	"status": "OK",
	"request_id": "abc123trades",
	"next_url": "https://api.massive.com/v3/trades/AAPL?cursor=YXA9MQ",
	"results": [
		{
			"conditions": [12, 37],
			"correction": 0,
			"exchange": 11,
			"id": "52983525029461",
			"participant_timestamp": 1736182800000000000,
			"price": 244.50,
			"sequence_number": 1234567,
			"sip_timestamp": 1736182800100000000,
			"size": 100,
			"tape": 3,
			"trf_id": 0,
			"trf_timestamp": 1736182800200000000
		},
		{
			"conditions": [12],
			"correction": 0,
			"exchange": 4,
			"id": "52983525029462",
			"participant_timestamp": 1736182801000000000,
			"price": 244.55,
			"sequence_number": 1234568,
			"sip_timestamp": 1736182801100000000,
			"size": 50,
			"tape": 3,
			"trf_id": 0,
			"trf_timestamp": 1736182801200000000
		}
	]
}`

const lastTradeJSON = `{
	"status": "OK",
	"request_id": "abc123lasttrade",
	"results": {
		"T": "AAPL",
		"c": [12, 37],
		"e": 0,
		"f": 1736182800200000000,
		"i": "52983525029461",
		"p": 244.50,
		"q": 1234567,
		"r": 0,
		"s": 100,
		"t": 1736182800100000000,
		"x": 11,
		"y": 1736182800000000000,
		"z": 3
	}
}`

const quotesJSON = `{
	"status": "OK",
	"request_id": "abc123quotes",
	"next_url": "https://api.massive.com/v3/quotes/AAPL?cursor=YXA9MQ",
	"results": [
		{
			"ask_exchange": 11,
			"ask_price": 244.55,
			"ask_size": 200,
			"bid_exchange": 19,
			"bid_price": 244.50,
			"bid_size": 300,
			"conditions": [1],
			"indicators": [0],
			"participant_timestamp": 1736182800000000000,
			"sequence_number": 9876543,
			"sip_timestamp": 1736182800100000000,
			"tape": 3,
			"trf_timestamp": 1736182800200000000
		},
		{
			"ask_exchange": 4,
			"ask_price": 244.60,
			"ask_size": 150,
			"bid_exchange": 11,
			"bid_price": 244.52,
			"bid_size": 250,
			"conditions": [1],
			"indicators": [0],
			"participant_timestamp": 1736182801000000000,
			"sequence_number": 9876544,
			"sip_timestamp": 1736182801100000000,
			"tape": 3,
			"trf_timestamp": 1736182801200000000
		}
	]
}`

const lastQuoteJSON = `{
	"status": "OK",
	"request_id": "abc123lastquote",
	"results": {
		"T": "AAPL",
		"P": 244.55,
		"S": 200,
		"X": 11,
		"c": [1],
		"f": 1736182800200000000,
		"i": [0],
		"p": 244.50,
		"q": 9876543,
		"s": 300,
		"t": 1736182800100000000,
		"x": 19,
		"y": 1736182800000000000,
		"z": 3
	}
}`

// TestGetTrades verifies that GetTrades correctly parses the API response
// and returns the expected tick-level trade data for AAPL.
func TestGetTrades(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/trades/AAPL": tradesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := TradesParams{
		Timestamp: "2025-01-06",
		Limit:     "2",
	}

	result, err := client.GetTrades("AAPL", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "abc123trades" {
		t.Errorf("expected request_id abc123trades, got %s", result.RequestID)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 trades, got %d", len(result.Results))
	}

	trade := result.Results[0]
	if trade.Price != 244.50 {
		t.Errorf("expected price 244.50, got %f", trade.Price)
	}

	if trade.Size != 100 {
		t.Errorf("expected size 100, got %f", trade.Size)
	}

	if trade.Exchange != 11 {
		t.Errorf("expected exchange 11, got %d", trade.Exchange)
	}

	if trade.ID != "52983525029461" {
		t.Errorf("expected id 52983525029461, got %s", trade.ID)
	}

	if trade.Tape != 3 {
		t.Errorf("expected tape 3, got %d", trade.Tape)
	}

	if trade.SequenceNumber != 1234567 {
		t.Errorf("expected sequence_number 1234567, got %d", trade.SequenceNumber)
	}

	if len(trade.Conditions) != 2 {
		t.Errorf("expected 2 conditions, got %d", len(trade.Conditions))
	}

	if trade.Conditions[0] != 12 {
		t.Errorf("expected first condition 12, got %d", trade.Conditions[0])
	}
}

// TestGetTradesSecondResult verifies that the second trade record in the
// response is correctly parsed with its own distinct values.
func TestGetTradesSecondResult(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/trades/AAPL": tradesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetTrades("AAPL", TradesParams{Timestamp: "2025-01-06"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	trade := result.Results[1]
	if trade.Price != 244.55 {
		t.Errorf("expected price 244.55, got %f", trade.Price)
	}

	if trade.Size != 50 {
		t.Errorf("expected size 50, got %f", trade.Size)
	}

	if trade.Exchange != 4 {
		t.Errorf("expected exchange 4, got %d", trade.Exchange)
	}

	if trade.ID != "52983525029462" {
		t.Errorf("expected id 52983525029462, got %s", trade.ID)
	}
}

// TestGetTradesRequestPath verifies that GetTrades constructs the correct
// API path with the ticker symbol.
func TestGetTradesRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(tradesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetTrades("MSFT", TradesParams{})

	if receivedPath != "/v3/trades/MSFT" {
		t.Errorf("expected path /v3/trades/MSFT, got %s", receivedPath)
	}
}

// TestGetTradesQueryParams verifies that all query parameters including
// timestamp range filters are correctly sent to the API.
func TestGetTradesQueryParams(t *testing.T) {
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
		w.Write([]byte(tradesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetTrades("AAPL", TradesParams{
		Timestamp: "2025-01-06",
		Order:     "desc",
		Limit:     "100",
		Sort:      "timestamp",
	})
}

// TestGetTradesTimestampRangeParams verifies that the timestamp range
// query parameters (gte, gt, lte, lt) are correctly sent to the API.
func TestGetTradesTimestampRangeParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("timestamp.gte") != "2025-01-06" {
			t.Errorf("expected timestamp.gte=2025-01-06, got %s", q.Get("timestamp.gte"))
		}
		if q.Get("timestamp.lte") != "2025-01-08" {
			t.Errorf("expected timestamp.lte=2025-01-08, got %s", q.Get("timestamp.lte"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(tradesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetTrades("AAPL", TradesParams{
		TimestampGte: "2025-01-06",
		TimestampLte: "2025-01-08",
	})
}

// TestGetTradesAPIError verifies that GetTrades returns an error when
// the API responds with a non-200 status code.
func TestGetTradesAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Data not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetTrades("INVALID", TradesParams{})
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// TestGetLastTrade verifies that GetLastTrade correctly parses the API
// response and returns the expected last trade data for AAPL.
func TestGetLastTrade(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/last/trade/AAPL": lastTradeJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetLastTrade("AAPL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "abc123lasttrade" {
		t.Errorf("expected request_id abc123lasttrade, got %s", result.RequestID)
	}

	trade := result.Results
	if trade.Ticker != "AAPL" {
		t.Errorf("expected ticker AAPL, got %s", trade.Ticker)
	}

	if trade.Price != 244.50 {
		t.Errorf("expected price 244.50, got %f", trade.Price)
	}

	if trade.Size != 100 {
		t.Errorf("expected size 100, got %f", trade.Size)
	}

	if trade.Exchange != 11 {
		t.Errorf("expected exchange 11, got %d", trade.Exchange)
	}

	if trade.ID != "52983525029461" {
		t.Errorf("expected id 52983525029461, got %s", trade.ID)
	}

	if trade.Tape != 3 {
		t.Errorf("expected tape 3, got %d", trade.Tape)
	}

	if trade.SequenceNumber != 1234567 {
		t.Errorf("expected sequence_number 1234567, got %d", trade.SequenceNumber)
	}

	if trade.Correction != 0 {
		t.Errorf("expected correction 0, got %d", trade.Correction)
	}

	if len(trade.Conditions) != 2 {
		t.Errorf("expected 2 conditions, got %d", len(trade.Conditions))
	}
}

// TestGetLastTradeRequestPath verifies that GetLastTrade constructs
// the correct API path with the ticker symbol.
func TestGetLastTradeRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(lastTradeJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetLastTrade("TSLA")

	if receivedPath != "/v2/last/trade/TSLA" {
		t.Errorf("expected path /v2/last/trade/TSLA, got %s", receivedPath)
	}
}

// TestGetLastTradeAPIError verifies that GetLastTrade returns an error
// when the API responds with a non-200 status code.
func TestGetLastTradeAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"status":"NOT_AUTHORIZED","message":"Not authorized."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetLastTrade("AAPL")
	if err == nil {
		t.Fatal("expected error for 401 response, got nil")
	}
}

// TestGetQuotes verifies that GetQuotes correctly parses the API response
// and returns the expected tick-level NBBO quote data for AAPL.
func TestGetQuotes(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/quotes/AAPL": quotesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := QuotesParams{
		Timestamp: "2025-01-06",
		Limit:     "2",
	}

	result, err := client.GetQuotes("AAPL", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "abc123quotes" {
		t.Errorf("expected request_id abc123quotes, got %s", result.RequestID)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 quotes, got %d", len(result.Results))
	}

	quote := result.Results[0]
	if quote.AskPrice != 244.55 {
		t.Errorf("expected ask_price 244.55, got %f", quote.AskPrice)
	}

	if quote.AskSize != 200 {
		t.Errorf("expected ask_size 200, got %f", quote.AskSize)
	}

	if quote.AskExchange != 11 {
		t.Errorf("expected ask_exchange 11, got %d", quote.AskExchange)
	}

	if quote.BidPrice != 244.50 {
		t.Errorf("expected bid_price 244.50, got %f", quote.BidPrice)
	}

	if quote.BidSize != 300 {
		t.Errorf("expected bid_size 300, got %f", quote.BidSize)
	}

	if quote.BidExchange != 19 {
		t.Errorf("expected bid_exchange 19, got %d", quote.BidExchange)
	}

	if quote.Tape != 3 {
		t.Errorf("expected tape 3, got %d", quote.Tape)
	}

	if quote.SequenceNumber != 9876543 {
		t.Errorf("expected sequence_number 9876543, got %d", quote.SequenceNumber)
	}
}

// TestGetQuotesSecondResult verifies that the second quote record in the
// response is correctly parsed with its own distinct values.
func TestGetQuotesSecondResult(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/quotes/AAPL": quotesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetQuotes("AAPL", QuotesParams{Timestamp: "2025-01-06"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	quote := result.Results[1]
	if quote.AskPrice != 244.60 {
		t.Errorf("expected ask_price 244.60, got %f", quote.AskPrice)
	}

	if quote.AskSize != 150 {
		t.Errorf("expected ask_size 150, got %f", quote.AskSize)
	}

	if quote.BidPrice != 244.52 {
		t.Errorf("expected bid_price 244.52, got %f", quote.BidPrice)
	}

	if quote.BidSize != 250 {
		t.Errorf("expected bid_size 250, got %f", quote.BidSize)
	}
}

// TestGetQuotesRequestPath verifies that GetQuotes constructs the correct
// API path with the ticker symbol.
func TestGetQuotesRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(quotesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetQuotes("MSFT", QuotesParams{})

	if receivedPath != "/v3/quotes/MSFT" {
		t.Errorf("expected path /v3/quotes/MSFT, got %s", receivedPath)
	}
}

// TestGetQuotesQueryParams verifies that all query parameters including
// timestamp range filters are correctly sent to the API.
func TestGetQuotesQueryParams(t *testing.T) {
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
		w.Write([]byte(quotesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetQuotes("AAPL", QuotesParams{
		Timestamp: "2025-01-06",
		Order:     "asc",
		Limit:     "500",
		Sort:      "timestamp",
	})
}

// TestGetQuotesTimestampRangeParams verifies that the timestamp range
// query parameters (gte, gt, lte, lt) are correctly sent to the API.
func TestGetQuotesTimestampRangeParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("timestamp.gt") != "2025-01-06" {
			t.Errorf("expected timestamp.gt=2025-01-06, got %s", q.Get("timestamp.gt"))
		}
		if q.Get("timestamp.lt") != "2025-01-08" {
			t.Errorf("expected timestamp.lt=2025-01-08, got %s", q.Get("timestamp.lt"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(quotesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetQuotes("AAPL", QuotesParams{
		TimestampGt: "2025-01-06",
		TimestampLt: "2025-01-08",
	})
}

// TestGetQuotesAPIError verifies that GetQuotes returns an error when
// the API responds with a non-200 status code.
func TestGetQuotesAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Data not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetQuotes("INVALID", QuotesParams{})
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// TestGetLastQuote verifies that GetLastQuote correctly parses the API
// response and returns the expected last NBBO quote data for AAPL.
func TestGetLastQuote(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/last/nbbo/AAPL": lastQuoteJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetLastQuote("AAPL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "abc123lastquote" {
		t.Errorf("expected request_id abc123lastquote, got %s", result.RequestID)
	}

	quote := result.Results
	if quote.Ticker != "AAPL" {
		t.Errorf("expected ticker AAPL, got %s", quote.Ticker)
	}

	if quote.AskPrice != 244.55 {
		t.Errorf("expected ask_price 244.55, got %f", quote.AskPrice)
	}

	if quote.AskSize != 200 {
		t.Errorf("expected ask_size 200, got %d", quote.AskSize)
	}

	if quote.AskExchange != 11 {
		t.Errorf("expected ask_exchange 11, got %d", quote.AskExchange)
	}

	if quote.BidPrice != 244.50 {
		t.Errorf("expected bid_price 244.50, got %f", quote.BidPrice)
	}

	if quote.BidSize != 300 {
		t.Errorf("expected bid_size 300, got %d", quote.BidSize)
	}

	if quote.BidExchange != 19 {
		t.Errorf("expected bid_exchange 19, got %d", quote.BidExchange)
	}

	if quote.Tape != 3 {
		t.Errorf("expected tape 3, got %d", quote.Tape)
	}

	if quote.SequenceNumber != 9876543 {
		t.Errorf("expected sequence_number 9876543, got %d", quote.SequenceNumber)
	}
}

// TestGetLastQuoteRequestPath verifies that GetLastQuote constructs
// the correct API path with the ticker symbol.
func TestGetLastQuoteRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(lastQuoteJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetLastQuote("TSLA")

	if receivedPath != "/v2/last/nbbo/TSLA" {
		t.Errorf("expected path /v2/last/nbbo/TSLA, got %s", receivedPath)
	}
}

// TestGetLastQuoteAPIError verifies that GetLastQuote returns an error
// when the API responds with a non-200 status code.
func TestGetLastQuoteAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"status":"NOT_AUTHORIZED","message":"Not authorized."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetLastQuote("AAPL")
	if err == nil {
		t.Fatal("expected error for 401 response, got nil")
	}
}
