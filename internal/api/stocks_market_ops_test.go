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

const marketStatusJSON = `{
	"afterHours": false,
	"currencies": {
		"crypto": "open",
		"fx": "open"
	},
	"earlyHours": false,
	"exchanges": {
		"nasdaq": "closed",
		"nyse": "closed",
		"otc": "closed"
	},
	"indicesGroups": {
		"s_and_p": "closed",
		"societe_generale": "closed",
		"msci": "closed",
		"ftse_russell": "closed",
		"mstar": "open",
		"mstarc": "open",
		"cccy": "open",
		"cgi": "closed",
		"nasdaq": "closed",
		"dow_jones": "closed"
	},
	"market": "closed",
	"serverTime": "2026-02-15T20:19:16-05:00"
}`

const marketHolidaysJSON = `[
	{
		"date": "2026-02-16",
		"exchange": "NYSE",
		"name": "Washington's Birthday",
		"status": "closed"
	},
	{
		"date": "2026-02-16",
		"exchange": "NASDAQ",
		"name": "Washington's Birthday",
		"status": "closed"
	},
	{
		"close": "2026-11-27T18:00:00.000Z",
		"date": "2026-11-27",
		"exchange": "NYSE",
		"name": "Thanksgiving",
		"open": "2026-11-27T14:30:00.000Z",
		"status": "early-close"
	}
]`

const exchangesJSON = `{
	"results": [
		{
			"id": 1,
			"type": "exchange",
			"asset_class": "stocks",
			"locale": "us",
			"name": "NYSE American, LLC",
			"acronym": "AMEX",
			"mic": "XASE",
			"operating_mic": "XNYS",
			"participant_id": "A",
			"url": "https://www.nyse.com/markets/nyse-american"
		},
		{
			"id": 10,
			"type": "exchange",
			"asset_class": "stocks",
			"locale": "us",
			"name": "New York Stock Exchange",
			"mic": "XNYS",
			"operating_mic": "XNYS",
			"participant_id": "N",
			"url": "https://www.nyse.com"
		},
		{
			"id": 12,
			"type": "exchange",
			"asset_class": "stocks",
			"locale": "us",
			"name": "Nasdaq",
			"mic": "XNAS",
			"operating_mic": "XNAS",
			"participant_id": "T",
			"url": "https://www.nasdaq.com"
		}
	],
	"status": "OK",
	"request_id": "abc123",
	"count": 3
}`

// TestGetMarketStatus verifies that GetMarketStatus correctly parses
// the API response including exchange statuses, currency statuses,
// indices groups, and top-level market fields.
func TestGetMarketStatus(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/now": marketStatusJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetMarketStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Market != "closed" {
		t.Errorf("expected market closed, got %s", result.Market)
	}

	if result.AfterHours != false {
		t.Error("expected afterHours to be false")
	}

	if result.EarlyHours != false {
		t.Error("expected earlyHours to be false")
	}

	if result.ServerTime != "2026-02-15T20:19:16-05:00" {
		t.Errorf("expected serverTime 2026-02-15T20:19:16-05:00, got %s", result.ServerTime)
	}
}

// TestGetMarketStatusExchanges verifies that the exchanges nested
// object within the market status response is correctly parsed with
// the individual NYSE, NASDAQ, and OTC statuses.
func TestGetMarketStatusExchanges(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/now": marketStatusJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetMarketStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Exchanges.NYSE != "closed" {
		t.Errorf("expected NYSE closed, got %s", result.Exchanges.NYSE)
	}

	if result.Exchanges.Nasdaq != "closed" {
		t.Errorf("expected NASDAQ closed, got %s", result.Exchanges.Nasdaq)
	}

	if result.Exchanges.OTC != "closed" {
		t.Errorf("expected OTC closed, got %s", result.Exchanges.OTC)
	}
}

// TestGetMarketStatusCurrencies verifies that the currencies nested
// object is parsed correctly with crypto and FX market statuses.
func TestGetMarketStatusCurrencies(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/now": marketStatusJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetMarketStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Currencies.Crypto != "open" {
		t.Errorf("expected crypto open, got %s", result.Currencies.Crypto)
	}

	if result.Currencies.FX != "open" {
		t.Errorf("expected fx open, got %s", result.Currencies.FX)
	}
}

// TestGetMarketStatusIndicesGroups verifies that the indices groups
// nested object is parsed with all index family statuses.
func TestGetMarketStatusIndicesGroups(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/now": marketStatusJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetMarketStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.IndicesGroups.SAndP != "closed" {
		t.Errorf("expected S&P closed, got %s", result.IndicesGroups.SAndP)
	}

	if result.IndicesGroups.DowJones != "closed" {
		t.Errorf("expected Dow Jones closed, got %s", result.IndicesGroups.DowJones)
	}

	if result.IndicesGroups.MStar != "open" {
		t.Errorf("expected MStar open, got %s", result.IndicesGroups.MStar)
	}

	if result.IndicesGroups.CCCY != "open" {
		t.Errorf("expected CCCY open, got %s", result.IndicesGroups.CCCY)
	}
}

// TestGetMarketStatusRequestPath verifies that GetMarketStatus sends
// the request to the correct API path.
func TestGetMarketStatusRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(marketStatusJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetMarketStatus()

	if receivedPath != "/v1/marketstatus/now" {
		t.Errorf("expected path /v1/marketstatus/now, got %s", receivedPath)
	}
}

// TestGetMarketStatusAPIError verifies that GetMarketStatus returns
// an error when the API responds with a non-200 status code.
func TestGetMarketStatusAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"ERROR","message":"Internal server error"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetMarketStatus()
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

// TestGetMarketHolidays verifies that GetMarketHolidays correctly
// parses the array response containing upcoming market holidays
// with both closed and early-close statuses.
func TestGetMarketHolidays(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/upcoming": marketHolidaysJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetMarketHolidays()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 holidays, got %d", len(result))
	}

	first := result[0]
	if first.Date != "2026-02-16" {
		t.Errorf("expected date 2026-02-16, got %s", first.Date)
	}

	if first.Exchange != "NYSE" {
		t.Errorf("expected exchange NYSE, got %s", first.Exchange)
	}

	if first.Name != "Washington's Birthday" {
		t.Errorf("expected name Washington's Birthday, got %s", first.Name)
	}

	if first.Status != "closed" {
		t.Errorf("expected status closed, got %s", first.Status)
	}
}

// TestGetMarketHolidaysEarlyClose verifies that early-close holidays
// include the open and close timestamp fields.
func TestGetMarketHolidaysEarlyClose(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/upcoming": marketHolidaysJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetMarketHolidays()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	earlyClose := result[2]
	if earlyClose.Status != "early-close" {
		t.Errorf("expected status early-close, got %s", earlyClose.Status)
	}

	if earlyClose.Open != "2026-11-27T14:30:00.000Z" {
		t.Errorf("expected open 2026-11-27T14:30:00.000Z, got %s", earlyClose.Open)
	}

	if earlyClose.Close != "2026-11-27T18:00:00.000Z" {
		t.Errorf("expected close 2026-11-27T18:00:00.000Z, got %s", earlyClose.Close)
	}

	if earlyClose.Exchange != "NYSE" {
		t.Errorf("expected exchange NYSE, got %s", earlyClose.Exchange)
	}

	if earlyClose.Name != "Thanksgiving" {
		t.Errorf("expected name Thanksgiving, got %s", earlyClose.Name)
	}
}

// TestGetMarketHolidaysRequestPath verifies that GetMarketHolidays
// sends the request to the correct API path.
func TestGetMarketHolidaysRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(marketHolidaysJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetMarketHolidays()

	if receivedPath != "/v1/marketstatus/upcoming" {
		t.Errorf("expected path /v1/marketstatus/upcoming, got %s", receivedPath)
	}
}

// TestGetMarketHolidaysAPIError verifies that GetMarketHolidays returns
// an error when the API responds with a non-200 status code.
func TestGetMarketHolidaysAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"ERROR","message":"Forbidden"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetMarketHolidays()
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetExchanges verifies that GetExchanges correctly parses the
// API response including the list of exchanges with their metadata.
func TestGetExchanges(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/reference/exchanges": exchangesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetExchanges(ExchangesParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Count != 3 {
		t.Errorf("expected count 3, got %d", result.Count)
	}

	if len(result.Results) != 3 {
		t.Fatalf("expected 3 exchanges, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.ID != 1 {
		t.Errorf("expected id 1, got %d", first.ID)
	}

	if first.Name != "NYSE American, LLC" {
		t.Errorf("expected name NYSE American, LLC, got %s", first.Name)
	}

	if first.Acronym != "AMEX" {
		t.Errorf("expected acronym AMEX, got %s", first.Acronym)
	}

	if first.MIC != "XASE" {
		t.Errorf("expected MIC XASE, got %s", first.MIC)
	}

	if first.AssetClass != "stocks" {
		t.Errorf("expected asset_class stocks, got %s", first.AssetClass)
	}

	if first.Locale != "us" {
		t.Errorf("expected locale us, got %s", first.Locale)
	}

	if first.Type != "exchange" {
		t.Errorf("expected type exchange, got %s", first.Type)
	}
}

// TestGetExchangesSecondAndThird verifies that multiple exchanges in
// the response are correctly parsed with their distinct values.
func TestGetExchangesSecondAndThird(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v3/reference/exchanges": exchangesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetExchanges(ExchangesParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	nyse := result.Results[1]
	if nyse.ID != 10 {
		t.Errorf("expected id 10, got %d", nyse.ID)
	}

	if nyse.Name != "New York Stock Exchange" {
		t.Errorf("expected name New York Stock Exchange, got %s", nyse.Name)
	}

	if nyse.MIC != "XNYS" {
		t.Errorf("expected MIC XNYS, got %s", nyse.MIC)
	}

	if nyse.ParticipantID != "N" {
		t.Errorf("expected participant_id N, got %s", nyse.ParticipantID)
	}

	nasdaq := result.Results[2]
	if nasdaq.ID != 12 {
		t.Errorf("expected id 12, got %d", nasdaq.ID)
	}

	if nasdaq.Name != "Nasdaq" {
		t.Errorf("expected name Nasdaq, got %s", nasdaq.Name)
	}

	if nasdaq.MIC != "XNAS" {
		t.Errorf("expected MIC XNAS, got %s", nasdaq.MIC)
	}

	if nasdaq.URL != "https://www.nasdaq.com" {
		t.Errorf("expected url https://www.nasdaq.com, got %s", nasdaq.URL)
	}
}

// TestGetExchangesQueryParams verifies that the asset_class and locale
// filter parameters are correctly sent to the API endpoint.
func TestGetExchangesQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("asset_class") != "stocks" {
			t.Errorf("expected asset_class=stocks, got %s", q.Get("asset_class"))
		}
		if q.Get("locale") != "us" {
			t.Errorf("expected locale=us, got %s", q.Get("locale"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(exchangesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetExchanges(ExchangesParams{
		AssetClass: "stocks",
		Locale:     "us",
	})
}

// TestGetExchangesRequestPath verifies that GetExchanges sends the
// request to the correct API path.
func TestGetExchangesRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(exchangesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetExchanges(ExchangesParams{})

	if receivedPath != "/v3/reference/exchanges" {
		t.Errorf("expected path /v3/reference/exchanges, got %s", receivedPath)
	}
}

// TestGetExchangesAPIError verifies that GetExchanges returns an error
// when the API responds with a non-200 status code.
func TestGetExchangesAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"status":"ERROR","message":"Unauthorized"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetExchanges(ExchangesParams{})
	if err == nil {
		t.Fatal("expected error for 401 response, got nil")
	}
}

// TestGetExchangesEmptyParams verifies that GetExchanges works correctly
// when no filter parameters are provided, sending no extra query params.
func TestGetExchangesEmptyParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("asset_class") != "" {
			t.Errorf("expected no asset_class param, got %s", q.Get("asset_class"))
		}
		if q.Get("locale") != "" {
			t.Errorf("expected no locale param, got %s", q.Get("locale"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(exchangesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetExchanges(ExchangesParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}
}
