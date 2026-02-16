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

const optionsChainSnapshotJSON = `{
	"status": "OK",
	"request_id": "opt-chain-123",
	"next_url": "https://api.massive.com/v3/snapshot/options/AAPL?cursor=abc123",
	"results": [
		{
			"break_even_price": 262.77,
			"day": {
				"change": -3.78,
				"change_percent": -22.83,
				"close": 12.77,
				"high": 16.55,
				"last_updated": 1771016400000000000,
				"low": 12.70,
				"open": 16.55,
				"previous_close": 16.55,
				"volume": 612,
				"vwap": 14.4929
			},
			"details": {
				"contract_type": "call",
				"exercise_style": "american",
				"expiration_date": "2026-03-20",
				"shares_per_contract": 100,
				"strike_price": 250,
				"ticker": "O:AAPL260320C00250000"
			},
			"greeks": {
				"delta": 0.6184,
				"gamma": 0.0161,
				"theta": -0.1524,
				"vega": 0.2970
			},
			"implied_volatility": 0.3117,
			"last_quote": {
				"ask": 12.90,
				"ask_size": 10,
				"bid": 12.65,
				"bid_size": 15,
				"last_updated": 1771016400000000000,
				"midpoint": 12.775,
				"timeframe": "REAL-TIME"
			},
			"last_trade": {
				"conditions": [209],
				"exchange": 316,
				"price": 12.77,
				"sip_timestamp": 1771016399000000000,
				"size": 5,
				"timeframe": "REAL-TIME"
			},
			"open_interest": 13765,
			"underlying_asset": {
				"change_to_break_even": 12.77,
				"last_updated": 1771016400000000000,
				"price": 250.00,
				"ticker": "AAPL",
				"timeframe": "REAL-TIME"
			}
		},
		{
			"break_even_price": 243.85,
			"day": {
				"change": 1.55,
				"change_percent": 33.70,
				"close": 6.15,
				"high": 6.27,
				"last_updated": 1771016400000000000,
				"low": 4.20,
				"open": 4.60,
				"previous_close": 4.60,
				"volume": 7049,
				"vwap": 5.1963
			},
			"details": {
				"contract_type": "put",
				"exercise_style": "american",
				"expiration_date": "2026-03-20",
				"shares_per_contract": 100,
				"strike_price": 250,
				"ticker": "O:AAPL260320P00250000"
			},
			"greeks": {
				"delta": -0.3777,
				"gamma": 0.0174,
				"theta": -0.1213,
				"vega": 0.2968
			},
			"implied_volatility": 0.2895,
			"open_interest": 21410,
			"underlying_asset": {
				"change_to_break_even": -6.15,
				"last_updated": 1771016400000000000,
				"price": 250.00,
				"ticker": "AAPL",
				"timeframe": "REAL-TIME"
			}
		}
	]
}`

const optionContractSnapshotJSON = `{
	"status": "OK",
	"request_id": "opt-contract-456",
	"results": {
		"break_even_price": 262.77,
		"day": {
			"change": -3.78,
			"change_percent": -22.83,
			"close": 12.77,
			"high": 16.55,
			"last_updated": 1771016400000000000,
			"low": 12.70,
			"open": 16.55,
			"previous_close": 16.55,
			"volume": 612,
			"vwap": 14.4929
		},
		"details": {
			"contract_type": "call",
			"exercise_style": "american",
			"expiration_date": "2026-03-20",
			"shares_per_contract": 100,
			"strike_price": 250,
			"ticker": "O:AAPL260320C00250000"
		},
		"greeks": {
			"delta": 0.6184,
			"gamma": 0.0161,
			"theta": -0.1524,
			"vega": 0.2970
		},
		"implied_volatility": 0.3117,
		"last_quote": {
			"ask": 12.90,
			"ask_size": 10,
			"bid": 12.65,
			"bid_size": 15,
			"last_updated": 1771016400000000000,
			"midpoint": 12.775,
			"timeframe": "REAL-TIME"
		},
		"last_trade": {
			"conditions": [209],
			"exchange": 316,
			"price": 12.77,
			"sip_timestamp": 1771016399000000000,
			"size": 5,
			"timeframe": "REAL-TIME"
		},
		"open_interest": 13765,
		"underlying_asset": {
			"change_to_break_even": 12.77,
			"last_updated": 1771016400000000000,
			"price": 250.00,
			"ticker": "AAPL",
			"timeframe": "REAL-TIME"
		}
	}
}`

// TestGetOptionsChainSnapshot verifies that GetOptionsChainSnapshot correctly
// parses the API response and returns the expected chain data for AAPL,
// including the results array, status, and request ID.
func TestGetOptionsChainSnapshot(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/options/AAPL": optionsChainSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsChainSnapshot("AAPL", OptionsChainSnapshotParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "opt-chain-123" {
		t.Errorf("expected request_id opt-chain-123, got %s", result.RequestID)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result.Results))
	}

	if result.NextURL != "https://api.massive.com/v3/snapshot/options/AAPL?cursor=abc123" {
		t.Errorf("expected next_url with cursor, got %s", result.NextURL)
	}
}

// TestGetOptionsChainSnapshotFirstResult verifies that the first option
// contract in the chain snapshot response has its details, day bar,
// Greeks, and implied volatility correctly parsed.
func TestGetOptionsChainSnapshotFirstResult(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/options/AAPL": optionsChainSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsChainSnapshot("AAPL", OptionsChainSnapshotParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	first := result.Results[0]

	if first.BreakEvenPrice != 262.77 {
		t.Errorf("expected break_even_price 262.77, got %f", first.BreakEvenPrice)
	}

	if first.Details.ContractType != "call" {
		t.Errorf("expected contract_type call, got %s", first.Details.ContractType)
	}

	if first.Details.StrikePrice != 250 {
		t.Errorf("expected strike_price 250, got %f", first.Details.StrikePrice)
	}

	if first.Details.ExpirationDate != "2026-03-20" {
		t.Errorf("expected expiration_date 2026-03-20, got %s", first.Details.ExpirationDate)
	}

	if first.Details.Ticker != "O:AAPL260320C00250000" {
		t.Errorf("expected ticker O:AAPL260320C00250000, got %s", first.Details.Ticker)
	}

	if first.Details.ExerciseStyle != "american" {
		t.Errorf("expected exercise_style american, got %s", first.Details.ExerciseStyle)
	}

	if first.Details.SharesPerContract != 100 {
		t.Errorf("expected shares_per_contract 100, got %f", first.Details.SharesPerContract)
	}
}

// TestGetOptionsChainSnapshotDayBar verifies that the day bar within the
// first option contract snapshot is correctly parsed with change, OHLC,
// volume, and VWAP values.
func TestGetOptionsChainSnapshotDayBar(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/options/AAPL": optionsChainSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsChainSnapshot("AAPL", OptionsChainSnapshotParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	day := result.Results[0].Day

	if day.Change != -3.78 {
		t.Errorf("expected day change -3.78, got %f", day.Change)
	}

	if day.ChangePercent != -22.83 {
		t.Errorf("expected day change_percent -22.83, got %f", day.ChangePercent)
	}

	if day.Close != 12.77 {
		t.Errorf("expected day close 12.77, got %f", day.Close)
	}

	if day.High != 16.55 {
		t.Errorf("expected day high 16.55, got %f", day.High)
	}

	if day.Low != 12.70 {
		t.Errorf("expected day low 12.70, got %f", day.Low)
	}

	if day.Open != 16.55 {
		t.Errorf("expected day open 16.55, got %f", day.Open)
	}

	if day.PreviousClose != 16.55 {
		t.Errorf("expected day previous_close 16.55, got %f", day.PreviousClose)
	}

	if day.Volume != 612 {
		t.Errorf("expected day volume 612, got %f", day.Volume)
	}

	if day.VWAP != 14.4929 {
		t.Errorf("expected day vwap 14.4929, got %f", day.VWAP)
	}

	if day.LastUpdated != 1771016400000000000 {
		t.Errorf("expected day last_updated 1771016400000000000, got %d", day.LastUpdated)
	}
}

// TestGetOptionsChainSnapshotGreeks verifies that the Greeks values
// (delta, gamma, theta, vega) are correctly parsed from the first
// option contract in the chain snapshot response.
func TestGetOptionsChainSnapshotGreeks(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/options/AAPL": optionsChainSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsChainSnapshot("AAPL", OptionsChainSnapshotParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	greeks := result.Results[0].Greeks

	if greeks.Delta != 0.6184 {
		t.Errorf("expected delta 0.6184, got %f", greeks.Delta)
	}

	if greeks.Gamma != 0.0161 {
		t.Errorf("expected gamma 0.0161, got %f", greeks.Gamma)
	}

	if greeks.Theta != -0.1524 {
		t.Errorf("expected theta -0.1524, got %f", greeks.Theta)
	}

	if greeks.Vega != 0.2970 {
		t.Errorf("expected vega 0.2970, got %f", greeks.Vega)
	}
}

// TestGetOptionsChainSnapshotImpliedVolatility verifies that the implied
// volatility value is correctly parsed from the chain snapshot response.
func TestGetOptionsChainSnapshotImpliedVolatility(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/options/AAPL": optionsChainSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsChainSnapshot("AAPL", OptionsChainSnapshotParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Results[0].ImpliedVolatility != 0.3117 {
		t.Errorf("expected implied_volatility 0.3117, got %f", result.Results[0].ImpliedVolatility)
	}
}

// TestGetOptionsChainSnapshotOpenInterest verifies that the open interest
// value is correctly parsed from the chain snapshot response.
func TestGetOptionsChainSnapshotOpenInterest(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/options/AAPL": optionsChainSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsChainSnapshot("AAPL", OptionsChainSnapshotParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Results[0].OpenInterest != 13765 {
		t.Errorf("expected open_interest 13765, got %f", result.Results[0].OpenInterest)
	}
}

// TestGetOptionsChainSnapshotLastQuote verifies that the last quote data
// (bid, ask, sizes, midpoint) is correctly parsed from the snapshot response.
func TestGetOptionsChainSnapshotLastQuote(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/options/AAPL": optionsChainSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsChainSnapshot("AAPL", OptionsChainSnapshotParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	quote := result.Results[0].LastQuote

	if quote.Ask != 12.90 {
		t.Errorf("expected ask 12.90, got %f", quote.Ask)
	}

	if quote.AskSize != 10 {
		t.Errorf("expected ask_size 10, got %f", quote.AskSize)
	}

	if quote.Bid != 12.65 {
		t.Errorf("expected bid 12.65, got %f", quote.Bid)
	}

	if quote.BidSize != 15 {
		t.Errorf("expected bid_size 15, got %f", quote.BidSize)
	}

	if quote.Midpoint != 12.775 {
		t.Errorf("expected midpoint 12.775, got %f", quote.Midpoint)
	}

	if quote.Timeframe != "REAL-TIME" {
		t.Errorf("expected timeframe REAL-TIME, got %s", quote.Timeframe)
	}
}

// TestGetOptionsChainSnapshotLastTrade verifies that the last trade data
// (price, size, exchange, conditions) is correctly parsed from the response.
func TestGetOptionsChainSnapshotLastTrade(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/options/AAPL": optionsChainSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsChainSnapshot("AAPL", OptionsChainSnapshotParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	trade := result.Results[0].LastTrade

	if trade.Price != 12.77 {
		t.Errorf("expected price 12.77, got %f", trade.Price)
	}

	if trade.Size != 5 {
		t.Errorf("expected size 5, got %f", trade.Size)
	}

	if trade.Exchange != 316 {
		t.Errorf("expected exchange 316, got %d", trade.Exchange)
	}

	if len(trade.Conditions) != 1 || trade.Conditions[0] != 209 {
		t.Errorf("expected conditions [209], got %v", trade.Conditions)
	}

	if trade.SipTimestamp != 1771016399000000000 {
		t.Errorf("expected sip_timestamp 1771016399000000000, got %d", trade.SipTimestamp)
	}

	if trade.Timeframe != "REAL-TIME" {
		t.Errorf("expected timeframe REAL-TIME, got %s", trade.Timeframe)
	}
}

// TestGetOptionsChainSnapshotUnderlyingAsset verifies that the underlying
// asset data (ticker, price, change_to_break_even) is correctly parsed.
func TestGetOptionsChainSnapshotUnderlyingAsset(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/options/AAPL": optionsChainSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsChainSnapshot("AAPL", OptionsChainSnapshotParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ua := result.Results[0].UnderlyingAsset

	if ua.Ticker != "AAPL" {
		t.Errorf("expected underlying ticker AAPL, got %s", ua.Ticker)
	}

	if ua.Price != 250.00 {
		t.Errorf("expected underlying price 250.00, got %f", ua.Price)
	}

	if ua.ChangeToBreakEven != 12.77 {
		t.Errorf("expected change_to_break_even 12.77, got %f", ua.ChangeToBreakEven)
	}

	if ua.Timeframe != "REAL-TIME" {
		t.Errorf("expected timeframe REAL-TIME, got %s", ua.Timeframe)
	}
}

// TestGetOptionsChainSnapshotSecondResult verifies that the second option
// contract (put) in the chain snapshot response is correctly parsed with
// its own distinct values.
func TestGetOptionsChainSnapshotSecondResult(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/options/AAPL": optionsChainSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsChainSnapshot("AAPL", OptionsChainSnapshotParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := result.Results[1]

	if second.Details.ContractType != "put" {
		t.Errorf("expected contract_type put, got %s", second.Details.ContractType)
	}

	if second.Details.Ticker != "O:AAPL260320P00250000" {
		t.Errorf("expected ticker O:AAPL260320P00250000, got %s", second.Details.Ticker)
	}

	if second.BreakEvenPrice != 243.85 {
		t.Errorf("expected break_even_price 243.85, got %f", second.BreakEvenPrice)
	}

	if second.Greeks.Delta != -0.3777 {
		t.Errorf("expected delta -0.3777, got %f", second.Greeks.Delta)
	}

	if second.ImpliedVolatility != 0.2895 {
		t.Errorf("expected implied_volatility 0.2895, got %f", second.ImpliedVolatility)
	}

	if second.OpenInterest != 21410 {
		t.Errorf("expected open_interest 21410, got %f", second.OpenInterest)
	}

	if second.Day.Volume != 7049 {
		t.Errorf("expected day volume 7049, got %f", second.Day.Volume)
	}
}

// TestGetOptionsChainSnapshotRequestPath verifies that GetOptionsChainSnapshot
// constructs the correct API path with the underlying asset ticker.
func TestGetOptionsChainSnapshotRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsChainSnapshotJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsChainSnapshot("TSLA", OptionsChainSnapshotParams{})

	expected := "/v3/snapshot/options/TSLA"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetOptionsChainSnapshotQueryParams verifies that the filter query
// parameters (strike_price, expiration_date, contract_type, limit, order,
// sort) are correctly sent to the API.
func TestGetOptionsChainSnapshotQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("strike_price") != "250" {
			t.Errorf("expected strike_price=250, got %s", q.Get("strike_price"))
		}
		if q.Get("expiration_date") != "2026-03-20" {
			t.Errorf("expected expiration_date=2026-03-20, got %s", q.Get("expiration_date"))
		}
		if q.Get("contract_type") != "call" {
			t.Errorf("expected contract_type=call, got %s", q.Get("contract_type"))
		}
		if q.Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", q.Get("limit"))
		}
		if q.Get("order") != "asc" {
			t.Errorf("expected order=asc, got %s", q.Get("order"))
		}
		if q.Get("sort") != "strike_price" {
			t.Errorf("expected sort=strike_price, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsChainSnapshotJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsChainSnapshot("AAPL", OptionsChainSnapshotParams{
		StrikePrice:    "250",
		ExpirationDate: "2026-03-20",
		ContractType:   "call",
		Limit:          "50",
		Order:          "asc",
		Sort:           "strike_price",
	})
}

// TestGetOptionsChainSnapshotRangeParams verifies that the range filter
// query parameters (strike_price.gte, expiration_date.lte, etc.) are
// correctly sent to the API.
func TestGetOptionsChainSnapshotRangeParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("strike_price.gte") != "200" {
			t.Errorf("expected strike_price.gte=200, got %s", q.Get("strike_price.gte"))
		}
		if q.Get("strike_price.lte") != "300" {
			t.Errorf("expected strike_price.lte=300, got %s", q.Get("strike_price.lte"))
		}
		if q.Get("expiration_date.gte") != "2026-03-01" {
			t.Errorf("expected expiration_date.gte=2026-03-01, got %s", q.Get("expiration_date.gte"))
		}
		if q.Get("expiration_date.lte") != "2026-06-30" {
			t.Errorf("expected expiration_date.lte=2026-06-30, got %s", q.Get("expiration_date.lte"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsChainSnapshotJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsChainSnapshot("AAPL", OptionsChainSnapshotParams{
		StrikePriceGTE:    "200",
		StrikePriceLTE:    "300",
		ExpirationDateGTE: "2026-03-01",
		ExpirationDateLTE: "2026-06-30",
	})
}

// TestGetOptionsChainSnapshotAPIError verifies that GetOptionsChainSnapshot
// returns an error when the API responds with a non-200 status code.
func TestGetOptionsChainSnapshotAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"ERROR","message":"Not authorized."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetOptionsChainSnapshot("AAPL", OptionsChainSnapshotParams{})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetOptionContractSnapshot verifies that GetOptionContractSnapshot
// correctly parses the API response for a single option contract,
// including the status, request ID, and contract details.
func TestGetOptionContractSnapshot(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/options/AAPL/O:AAPL260320C00250000": optionContractSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionContractSnapshot("AAPL", "O:AAPL260320C00250000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "opt-contract-456" {
		t.Errorf("expected request_id opt-contract-456, got %s", result.RequestID)
	}
}

// TestGetOptionContractSnapshotDetails verifies that the contract details
// (type, style, expiration, shares, strike, ticker) are correctly parsed
// from the single contract snapshot response.
func TestGetOptionContractSnapshotDetails(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/options/AAPL/O:AAPL260320C00250000": optionContractSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionContractSnapshot("AAPL", "O:AAPL260320C00250000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	details := result.Results.Details

	if details.ContractType != "call" {
		t.Errorf("expected contract_type call, got %s", details.ContractType)
	}

	if details.ExerciseStyle != "american" {
		t.Errorf("expected exercise_style american, got %s", details.ExerciseStyle)
	}

	if details.ExpirationDate != "2026-03-20" {
		t.Errorf("expected expiration_date 2026-03-20, got %s", details.ExpirationDate)
	}

	if details.SharesPerContract != 100 {
		t.Errorf("expected shares_per_contract 100, got %f", details.SharesPerContract)
	}

	if details.StrikePrice != 250 {
		t.Errorf("expected strike_price 250, got %f", details.StrikePrice)
	}

	if details.Ticker != "O:AAPL260320C00250000" {
		t.Errorf("expected ticker O:AAPL260320C00250000, got %s", details.Ticker)
	}
}

// TestGetOptionContractSnapshotGreeks verifies that the Greeks values
// are correctly parsed from the single contract snapshot response.
func TestGetOptionContractSnapshotGreeks(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/options/AAPL/O:AAPL260320C00250000": optionContractSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionContractSnapshot("AAPL", "O:AAPL260320C00250000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	greeks := result.Results.Greeks

	if greeks.Delta != 0.6184 {
		t.Errorf("expected delta 0.6184, got %f", greeks.Delta)
	}

	if greeks.Gamma != 0.0161 {
		t.Errorf("expected gamma 0.0161, got %f", greeks.Gamma)
	}

	if greeks.Theta != -0.1524 {
		t.Errorf("expected theta -0.1524, got %f", greeks.Theta)
	}

	if greeks.Vega != 0.2970 {
		t.Errorf("expected vega 0.2970, got %f", greeks.Vega)
	}
}

// TestGetOptionContractSnapshotDayBar verifies that the day bar within
// the single contract snapshot is correctly parsed.
func TestGetOptionContractSnapshotDayBar(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/options/AAPL/O:AAPL260320C00250000": optionContractSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionContractSnapshot("AAPL", "O:AAPL260320C00250000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	day := result.Results.Day

	if day.Close != 12.77 {
		t.Errorf("expected day close 12.77, got %f", day.Close)
	}

	if day.High != 16.55 {
		t.Errorf("expected day high 16.55, got %f", day.High)
	}

	if day.Volume != 612 {
		t.Errorf("expected day volume 612, got %f", day.Volume)
	}

	if day.VWAP != 14.4929 {
		t.Errorf("expected day vwap 14.4929, got %f", day.VWAP)
	}
}

// TestGetOptionContractSnapshotBreakEven verifies that the break-even
// price and open interest are correctly parsed from the response.
func TestGetOptionContractSnapshotBreakEven(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/options/AAPL/O:AAPL260320C00250000": optionContractSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionContractSnapshot("AAPL", "O:AAPL260320C00250000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Results.BreakEvenPrice != 262.77 {
		t.Errorf("expected break_even_price 262.77, got %f", result.Results.BreakEvenPrice)
	}

	if result.Results.OpenInterest != 13765 {
		t.Errorf("expected open_interest 13765, got %f", result.Results.OpenInterest)
	}

	if result.Results.ImpliedVolatility != 0.3117 {
		t.Errorf("expected implied_volatility 0.3117, got %f", result.Results.ImpliedVolatility)
	}
}

// TestGetOptionContractSnapshotRequestPath verifies that
// GetOptionContractSnapshot constructs the correct API path with both
// the underlying asset and option contract tickers.
func TestGetOptionContractSnapshotRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionContractSnapshotJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionContractSnapshot("TSLA", "O:TSLA260320C00200000")

	expected := "/v3/snapshot/options/TSLA/O:TSLA260320C00200000"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetOptionContractSnapshotAPIError verifies that
// GetOptionContractSnapshot returns an error when the API responds with
// a non-200 status code.
func TestGetOptionContractSnapshotAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"ERROR","error":"Options contract not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetOptionContractSnapshot("AAPL", "O:INVALID")
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// TestGetOptionContractSnapshotUnderlyingAsset verifies that the
// underlying asset data is correctly parsed from the single contract
// snapshot response.
func TestGetOptionContractSnapshotUnderlyingAsset(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/snapshot/options/AAPL/O:AAPL260320C00250000": optionContractSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionContractSnapshot("AAPL", "O:AAPL260320C00250000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ua := result.Results.UnderlyingAsset

	if ua.Ticker != "AAPL" {
		t.Errorf("expected underlying ticker AAPL, got %s", ua.Ticker)
	}

	if ua.Price != 250.00 {
		t.Errorf("expected underlying price 250.00, got %f", ua.Price)
	}

	if ua.ChangeToBreakEven != 12.77 {
		t.Errorf("expected change_to_break_even 12.77, got %f", ua.ChangeToBreakEven)
	}
}
