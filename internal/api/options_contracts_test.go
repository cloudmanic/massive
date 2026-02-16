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

const optionsContractsJSON = `{
	"results": [
		{
			"cfi": "OCASPS",
			"contract_type": "call",
			"exercise_style": "american",
			"expiration_date": "2026-02-18",
			"primary_exchange": "BATO",
			"shares_per_contract": 100,
			"strike_price": 190,
			"ticker": "O:AAPL260218C00190000",
			"underlying_ticker": "AAPL"
		},
		{
			"cfi": "OCASPS",
			"contract_type": "call",
			"exercise_style": "american",
			"expiration_date": "2026-02-18",
			"primary_exchange": "BATO",
			"shares_per_contract": 100,
			"strike_price": 195,
			"ticker": "O:AAPL260218C00195000",
			"underlying_ticker": "AAPL"
		}
	],
	"status": "OK",
	"request_id": "64574a27abd280ad61a9aaf38d9e1d0e",
	"next_url": "https://api.massive.com/v3/reference/options/contracts?cursor=YXA9Mg"
}`

const optionsContractJSON = `{
	"results": {
		"cfi": "OCASPS",
		"contract_type": "call",
		"exercise_style": "american",
		"expiration_date": "2026-02-18",
		"primary_exchange": "BATO",
		"shares_per_contract": 100,
		"strike_price": 190,
		"ticker": "O:AAPL260218C00190000",
		"underlying_ticker": "AAPL"
	},
	"status": "OK",
	"request_id": "6fef59537c50bcbb32a7190f82106439"
}`

// TestGetOptionsContracts verifies that GetOptionsContracts correctly parses
// the API response and returns the expected list of options contracts including
// ticker, underlying ticker, contract type, strike price, and pagination info.
func TestGetOptionsContracts(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/reference/options/contracts": optionsContractsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := OptionsContractsParams{
		UnderlyingTicker: "AAPL",
		Limit:            "2",
	}

	result, err := client.GetOptionsContracts(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "64574a27abd280ad61a9aaf38d9e1d0e" {
		t.Errorf("expected request_id 64574a27abd280ad61a9aaf38d9e1d0e, got %s", result.RequestID)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 contracts, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.Ticker != "O:AAPL260218C00190000" {
		t.Errorf("expected ticker O:AAPL260218C00190000, got %s", first.Ticker)
	}

	if first.UnderlyingTicker != "AAPL" {
		t.Errorf("expected underlying_ticker AAPL, got %s", first.UnderlyingTicker)
	}

	if first.ContractType != "call" {
		t.Errorf("expected contract_type call, got %s", first.ContractType)
	}

	if first.ExerciseStyle != "american" {
		t.Errorf("expected exercise_style american, got %s", first.ExerciseStyle)
	}

	if first.ExpirationDate != "2026-02-18" {
		t.Errorf("expected expiration_date 2026-02-18, got %s", first.ExpirationDate)
	}

	if first.StrikePrice != 190 {
		t.Errorf("expected strike_price 190, got %f", first.StrikePrice)
	}

	if first.SharesPerContract != 100 {
		t.Errorf("expected shares_per_contract 100, got %d", first.SharesPerContract)
	}

	if first.PrimaryExchange != "BATO" {
		t.Errorf("expected primary_exchange BATO, got %s", first.PrimaryExchange)
	}

	if first.CFI != "OCASPS" {
		t.Errorf("expected cfi OCASPS, got %s", first.CFI)
	}

	second := result.Results[1]
	if second.Ticker != "O:AAPL260218C00195000" {
		t.Errorf("expected ticker O:AAPL260218C00195000, got %s", second.Ticker)
	}

	if second.StrikePrice != 195 {
		t.Errorf("expected strike_price 195, got %f", second.StrikePrice)
	}
}

// TestGetOptionsContractsRequestPath verifies that GetOptionsContracts sends
// requests to the correct /v3/reference/options/contracts API path.
func TestGetOptionsContractsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsContractsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsContracts(OptionsContractsParams{})

	if receivedPath != "/v3/reference/options/contracts" {
		t.Errorf("expected path /v3/reference/options/contracts, got %s", receivedPath)
	}
}

// TestGetOptionsContractsQueryParams verifies that all filter parameters
// are correctly sent to the API endpoint including underlying_ticker,
// contract_type, expiration_date, strike_price, order, limit, and sort.
func TestGetOptionsContractsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("underlying_ticker") != "AAPL" {
			t.Errorf("expected underlying_ticker=AAPL, got %s", q.Get("underlying_ticker"))
		}
		if q.Get("contract_type") != "call" {
			t.Errorf("expected contract_type=call, got %s", q.Get("contract_type"))
		}
		if q.Get("expiration_date") != "2026-02-18" {
			t.Errorf("expected expiration_date=2026-02-18, got %s", q.Get("expiration_date"))
		}
		if q.Get("strike_price") != "190" {
			t.Errorf("expected strike_price=190, got %s", q.Get("strike_price"))
		}
		if q.Get("expired") != "false" {
			t.Errorf("expected expired=false, got %s", q.Get("expired"))
		}
		if q.Get("order") != "asc" {
			t.Errorf("expected order=asc, got %s", q.Get("order"))
		}
		if q.Get("limit") != "100" {
			t.Errorf("expected limit=100, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "ticker" {
			t.Errorf("expected sort=ticker, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsContractsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsContracts(OptionsContractsParams{
		UnderlyingTicker: "AAPL",
		ContractType:     "call",
		ExpirationDate:   "2026-02-18",
		StrikePrice:      "190",
		Expired:          "false",
		Order:            "asc",
		Limit:            "100",
		Sort:             "ticker",
	})
}

// TestGetOptionsContractsRangeParams verifies that the range filter
// parameters (gte, gt, lte, lt) for expiration_date and strike_price
// are correctly sent to the API endpoint.
func TestGetOptionsContractsRangeParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("expiration_date.gte") != "2026-01-01" {
			t.Errorf("expected expiration_date.gte=2026-01-01, got %s", q.Get("expiration_date.gte"))
		}
		if q.Get("expiration_date.lte") != "2026-12-31" {
			t.Errorf("expected expiration_date.lte=2026-12-31, got %s", q.Get("expiration_date.lte"))
		}
		if q.Get("strike_price.gte") != "100" {
			t.Errorf("expected strike_price.gte=100, got %s", q.Get("strike_price.gte"))
		}
		if q.Get("strike_price.lte") != "200" {
			t.Errorf("expected strike_price.lte=200, got %s", q.Get("strike_price.lte"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsContractsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsContracts(OptionsContractsParams{
		ExpirationDateGte: "2026-01-01",
		ExpirationDateLte: "2026-12-31",
		StrikePriceGte:    "100",
		StrikePriceLte:    "200",
	})
}

// TestGetOptionsContractsAsOfParam verifies that the as_of query parameter
// is correctly sent to the API when requesting a historical snapshot.
func TestGetOptionsContractsAsOfParam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("as_of") != "2025-12-31" {
			t.Errorf("expected as_of=2025-12-31, got %s", q.Get("as_of"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsContractsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsContracts(OptionsContractsParams{
		AsOf: "2025-12-31",
	})
}

// TestGetOptionsContractsEmptyResults verifies that GetOptionsContracts
// handles an empty results array without error.
func TestGetOptionsContractsEmptyResults(t *testing.T) {
	emptyJSON := `{"results":[],"status":"OK","request_id":"empty123"}`
	server := mockServer(t, map[string]string{
		"/v3/reference/options/contracts": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsContracts(OptionsContractsParams{
		UnderlyingTicker: "ZZZZZZ",
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

// TestGetOptionsContractsAPIError verifies that GetOptionsContracts returns
// an error when the API responds with a non-200 status code.
func TestGetOptionsContractsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"status":"ERROR","message":"Invalid API key"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetOptionsContracts(OptionsContractsParams{})
	if err == nil {
		t.Fatal("expected error for 401 response, got nil")
	}
}

// TestGetOptionsContract verifies that GetOptionsContract correctly parses
// the API response for a single options contract and returns all expected
// fields including ticker, underlying, type, strike price, and expiration.
func TestGetOptionsContract(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/reference/options/contracts/O:AAPL260218C00190000": optionsContractJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsContract("O:AAPL260218C00190000", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "6fef59537c50bcbb32a7190f82106439" {
		t.Errorf("expected request_id 6fef59537c50bcbb32a7190f82106439, got %s", result.RequestID)
	}

	contract := result.Results
	if contract.Ticker != "O:AAPL260218C00190000" {
		t.Errorf("expected ticker O:AAPL260218C00190000, got %s", contract.Ticker)
	}

	if contract.UnderlyingTicker != "AAPL" {
		t.Errorf("expected underlying_ticker AAPL, got %s", contract.UnderlyingTicker)
	}

	if contract.ContractType != "call" {
		t.Errorf("expected contract_type call, got %s", contract.ContractType)
	}

	if contract.ExerciseStyle != "american" {
		t.Errorf("expected exercise_style american, got %s", contract.ExerciseStyle)
	}

	if contract.ExpirationDate != "2026-02-18" {
		t.Errorf("expected expiration_date 2026-02-18, got %s", contract.ExpirationDate)
	}

	if contract.StrikePrice != 190 {
		t.Errorf("expected strike_price 190, got %f", contract.StrikePrice)
	}

	if contract.SharesPerContract != 100 {
		t.Errorf("expected shares_per_contract 100, got %d", contract.SharesPerContract)
	}

	if contract.PrimaryExchange != "BATO" {
		t.Errorf("expected primary_exchange BATO, got %s", contract.PrimaryExchange)
	}

	if contract.CFI != "OCASPS" {
		t.Errorf("expected cfi OCASPS, got %s", contract.CFI)
	}
}

// TestGetOptionsContractRequestPath verifies that GetOptionsContract sends
// requests to the correct API path with the options ticker in the URL.
func TestGetOptionsContractRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsContractJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsContract("O:TSLA260320P00250000", "")

	expected := "/v3/reference/options/contracts/O:TSLA260320P00250000"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetOptionsContractAsOfParam verifies that the as_of query parameter
// is correctly sent when requesting a historical snapshot of a single contract.
func TestGetOptionsContractAsOfParam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("as_of") != "2025-12-31" {
			t.Errorf("expected as_of=2025-12-31, got %s", q.Get("as_of"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsContractJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsContract("O:AAPL260218C00190000", "2025-12-31")
}

// TestGetOptionsContractAPIError verifies that GetOptionsContract returns
// an error when the API responds with a non-200 status code such as a
// 404 for an invalid options ticker.
func TestGetOptionsContractAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","request_id":"abc","message":"Option Ticker not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetOptionsContract("O:INVALID000000C00000000", "")
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// TestGetOptionsContractAdditionalUnderlyings verifies that GetOptionsContract
// correctly parses the additional_underlyings array when present in the
// API response for contracts with multiple underlying assets.
func TestGetOptionsContractAdditionalUnderlyings(t *testing.T) {
	contractWithUnderlyings := `{
		"results": {
			"cfi": "OCASPS",
			"contract_type": "call",
			"exercise_style": "american",
			"expiration_date": "2026-02-18",
			"primary_exchange": "BATO",
			"shares_per_contract": 100,
			"strike_price": 190,
			"ticker": "O:AAPL260218C00190000",
			"underlying_ticker": "AAPL",
			"correction": 1,
			"additional_underlyings": [
				{
					"underlying": "AAPL",
					"amount": 0.50,
					"type": "equity"
				}
			]
		},
		"status": "OK",
		"request_id": "test123"
	}`

	server := mockServer(t, map[string]string{
		"/v3/reference/options/contracts/O:AAPL260218C00190000": contractWithUnderlyings,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsContract("O:AAPL260218C00190000", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	contract := result.Results
	if contract.Correction != 1 {
		t.Errorf("expected correction 1, got %d", contract.Correction)
	}

	if len(contract.AdditionalUnderlyings) != 1 {
		t.Fatalf("expected 1 additional underlying, got %d", len(contract.AdditionalUnderlyings))
	}

	au := contract.AdditionalUnderlyings[0]
	if au.Underlying != "AAPL" {
		t.Errorf("expected underlying AAPL, got %s", au.Underlying)
	}

	if au.Amount != 0.50 {
		t.Errorf("expected amount 0.50, got %f", au.Amount)
	}

	if au.Type != "equity" {
		t.Errorf("expected type equity, got %s", au.Type)
	}
}
