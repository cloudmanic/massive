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

// indicesMarketStatusJSON is the mock JSON response for the indices market
// status endpoint. It mirrors the real API structure with exchanges,
// currencies, and indices groups all reporting their current statuses.
const indicesMarketStatusJSON = `{
	"afterHours": true,
	"currencies": {
		"crypto": "open",
		"fx": "closed"
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
	"market": "extended-hours",
	"serverTime": "2026-02-15T18:30:00-05:00"
}`

// indicesMarketHolidaysJSON is the mock JSON response for the indices
// market holidays endpoint. It contains both closed and early-close
// entries to verify proper parsing of optional open/close fields.
const indicesMarketHolidaysJSON = `[
	{
		"date": "2026-05-25",
		"exchange": "NYSE",
		"name": "Memorial Day",
		"status": "closed"
	},
	{
		"date": "2026-05-25",
		"exchange": "NASDAQ",
		"name": "Memorial Day",
		"status": "closed"
	},
	{
		"date": "2026-07-03",
		"exchange": "NYSE",
		"name": "Independence Day",
		"open": "2026-07-03T14:30:00.000Z",
		"close": "2026-07-03T18:00:00.000Z",
		"status": "early-close"
	}
]`

// TestGetIndicesMarketStatus verifies that GetIndicesMarketStatus correctly
// parses the API response including the top-level market fields, after-hours
// flag, and server time.
func TestGetIndicesMarketStatus(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/now": indicesMarketStatusJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesMarketStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Market != "extended-hours" {
		t.Errorf("expected market extended-hours, got %s", result.Market)
	}

	if result.AfterHours != true {
		t.Error("expected afterHours to be true")
	}

	if result.EarlyHours != false {
		t.Error("expected earlyHours to be false")
	}

	if result.ServerTime != "2026-02-15T18:30:00-05:00" {
		t.Errorf("expected serverTime 2026-02-15T18:30:00-05:00, got %s", result.ServerTime)
	}
}

// TestGetIndicesMarketStatusExchanges verifies that the exchanges nested
// object within the indices market status response is correctly parsed
// with the individual NYSE, NASDAQ, and OTC statuses.
func TestGetIndicesMarketStatusExchanges(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/now": indicesMarketStatusJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesMarketStatus()
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

// TestGetIndicesMarketStatusCurrencies verifies that the currencies nested
// object is parsed correctly with crypto and FX market statuses.
func TestGetIndicesMarketStatusCurrencies(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/now": indicesMarketStatusJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesMarketStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Currencies.Crypto != "open" {
		t.Errorf("expected crypto open, got %s", result.Currencies.Crypto)
	}

	if result.Currencies.FX != "closed" {
		t.Errorf("expected fx closed, got %s", result.Currencies.FX)
	}
}

// TestGetIndicesMarketStatusIndicesGroups verifies that the indices groups
// nested object is parsed with all index family statuses including S&P,
// Dow Jones, MStar, and CCCY groups.
func TestGetIndicesMarketStatusIndicesGroups(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/now": indicesMarketStatusJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesMarketStatus()
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

	if result.IndicesGroups.FTSERussell != "closed" {
		t.Errorf("expected FTSE Russell closed, got %s", result.IndicesGroups.FTSERussell)
	}

	if result.IndicesGroups.MSCI != "closed" {
		t.Errorf("expected MSCI closed, got %s", result.IndicesGroups.MSCI)
	}
}

// TestGetIndicesMarketStatusRequestPath verifies that GetIndicesMarketStatus
// sends the request to the correct API path (/v1/marketstatus/now).
func TestGetIndicesMarketStatusRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesMarketStatusJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesMarketStatus()

	if receivedPath != "/v1/marketstatus/now" {
		t.Errorf("expected path /v1/marketstatus/now, got %s", receivedPath)
	}
}

// TestGetIndicesMarketStatusAPIError verifies that GetIndicesMarketStatus
// returns an error when the API responds with a non-200 status code.
func TestGetIndicesMarketStatusAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"ERROR","message":"Internal server error"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetIndicesMarketStatus()
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

// TestGetIndicesMarketStatusInvalidJSON verifies that GetIndicesMarketStatus
// returns an error when the API responds with malformed JSON.
func TestGetIndicesMarketStatusInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`not valid json`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetIndicesMarketStatus()
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

// TestGetIndicesMarketHolidays verifies that GetIndicesMarketHolidays
// correctly parses the array response containing upcoming market holidays
// with both closed and early-close statuses.
func TestGetIndicesMarketHolidays(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/upcoming": indicesMarketHolidaysJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesMarketHolidays()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 holidays, got %d", len(result))
	}

	first := result[0]
	if first.Date != "2026-05-25" {
		t.Errorf("expected date 2026-05-25, got %s", first.Date)
	}

	if first.Exchange != "NYSE" {
		t.Errorf("expected exchange NYSE, got %s", first.Exchange)
	}

	if first.Name != "Memorial Day" {
		t.Errorf("expected name Memorial Day, got %s", first.Name)
	}

	if first.Status != "closed" {
		t.Errorf("expected status closed, got %s", first.Status)
	}
}

// TestGetIndicesMarketHolidaysSecondEntry verifies that the second holiday
// in the response is correctly parsed with its distinct exchange value.
func TestGetIndicesMarketHolidaysSecondEntry(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/upcoming": indicesMarketHolidaysJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesMarketHolidays()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := result[1]
	if second.Exchange != "NASDAQ" {
		t.Errorf("expected exchange NASDAQ, got %s", second.Exchange)
	}

	if second.Name != "Memorial Day" {
		t.Errorf("expected name Memorial Day, got %s", second.Name)
	}

	if second.Status != "closed" {
		t.Errorf("expected status closed, got %s", second.Status)
	}
}

// TestGetIndicesMarketHolidaysEarlyClose verifies that early-close holidays
// include the open and close timestamp fields with their correct values.
func TestGetIndicesMarketHolidaysEarlyClose(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/upcoming": indicesMarketHolidaysJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesMarketHolidays()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	earlyClose := result[2]
	if earlyClose.Status != "early-close" {
		t.Errorf("expected status early-close, got %s", earlyClose.Status)
	}

	if earlyClose.Open != "2026-07-03T14:30:00.000Z" {
		t.Errorf("expected open 2026-07-03T14:30:00.000Z, got %s", earlyClose.Open)
	}

	if earlyClose.Close != "2026-07-03T18:00:00.000Z" {
		t.Errorf("expected close 2026-07-03T18:00:00.000Z, got %s", earlyClose.Close)
	}

	if earlyClose.Exchange != "NYSE" {
		t.Errorf("expected exchange NYSE, got %s", earlyClose.Exchange)
	}

	if earlyClose.Name != "Independence Day" {
		t.Errorf("expected name Independence Day, got %s", earlyClose.Name)
	}
}

// TestGetIndicesMarketHolidaysRequestPath verifies that GetIndicesMarketHolidays
// sends the request to the correct API path (/v1/marketstatus/upcoming).
func TestGetIndicesMarketHolidaysRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesMarketHolidaysJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesMarketHolidays()

	if receivedPath != "/v1/marketstatus/upcoming" {
		t.Errorf("expected path /v1/marketstatus/upcoming, got %s", receivedPath)
	}
}

// TestGetIndicesMarketHolidaysAPIError verifies that GetIndicesMarketHolidays
// returns an error when the API responds with a non-200 status code.
func TestGetIndicesMarketHolidaysAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"ERROR","message":"Forbidden"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetIndicesMarketHolidays()
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetIndicesMarketHolidaysInvalidJSON verifies that GetIndicesMarketHolidays
// returns an error when the API responds with malformed JSON.
func TestGetIndicesMarketHolidaysInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{invalid json`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetIndicesMarketHolidays()
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

// TestGetIndicesMarketHolidaysEmptyArray verifies that GetIndicesMarketHolidays
// correctly handles an empty array response with no upcoming holidays.
func TestGetIndicesMarketHolidaysEmptyArray(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/upcoming": `[]`,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIndicesMarketHolidays()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected 0 holidays, got %d", len(result))
	}
}

// TestGetIndicesMarketStatusAPIKeySent verifies that the API key is included
// as a query parameter in the request to the market status endpoint.
func TestGetIndicesMarketStatusAPIKeySent(t *testing.T) {
	var receivedKey string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedKey = r.URL.Query().Get("apiKey")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesMarketStatusJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesMarketStatus()

	if receivedKey != "test-api-key" {
		t.Errorf("expected apiKey=test-api-key, got %s", receivedKey)
	}
}

// TestGetIndicesMarketHolidaysAPIKeySent verifies that the API key is included
// as a query parameter in the request to the market holidays endpoint.
func TestGetIndicesMarketHolidaysAPIKeySent(t *testing.T) {
	var receivedKey string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedKey = r.URL.Query().Get("apiKey")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(indicesMarketHolidaysJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIndicesMarketHolidays()

	if receivedKey != "test-api-key" {
		t.Errorf("expected apiKey=test-api-key, got %s", receivedKey)
	}
}
