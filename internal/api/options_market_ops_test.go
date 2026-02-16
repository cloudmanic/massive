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

// optionsMarketStatusJSON is the mock JSON response for the options market
// status endpoint. It mirrors the real API structure with exchanges,
// currencies, and indices groups all reporting their current statuses.
const optionsMarketStatusJSON = `{
	"afterHours": false,
	"currencies": {
		"crypto": "open",
		"fx": "open"
	},
	"earlyHours": true,
	"exchanges": {
		"nasdaq": "open",
		"nyse": "open",
		"otc": "closed"
	},
	"indicesGroups": {
		"s_and_p": "open",
		"societe_generale": "closed",
		"msci": "closed",
		"ftse_russell": "open",
		"mstar": "open",
		"mstarc": "closed",
		"cccy": "open",
		"cgi": "closed",
		"nasdaq": "open",
		"dow_jones": "open"
	},
	"market": "open",
	"serverTime": "2026-02-15T10:30:00-05:00"
}`

// optionsMarketHolidaysJSON is the mock JSON response for the options
// market holidays endpoint. It contains both closed and early-close
// entries to verify proper parsing of optional open/close fields.
const optionsMarketHolidaysJSON = `[
	{
		"date": "2026-04-03",
		"exchange": "NYSE",
		"name": "Good Friday",
		"status": "closed"
	},
	{
		"date": "2026-04-03",
		"exchange": "NASDAQ",
		"name": "Good Friday",
		"status": "closed"
	},
	{
		"date": "2026-11-27",
		"exchange": "NYSE",
		"name": "Thanksgiving",
		"open": "2026-11-27T14:30:00.000Z",
		"close": "2026-11-27T18:00:00.000Z",
		"status": "early-close"
	}
]`

// TestGetOptionsMarketStatus verifies that GetOptionsMarketStatus correctly
// parses the API response including the top-level market fields, early-hours
// flag, and server time.
func TestGetOptionsMarketStatus(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/now": optionsMarketStatusJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsMarketStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Market != "open" {
		t.Errorf("expected market open, got %s", result.Market)
	}

	if result.AfterHours != false {
		t.Error("expected afterHours to be false")
	}

	if result.EarlyHours != true {
		t.Error("expected earlyHours to be true")
	}

	if result.ServerTime != "2026-02-15T10:30:00-05:00" {
		t.Errorf("expected serverTime 2026-02-15T10:30:00-05:00, got %s", result.ServerTime)
	}
}

// TestGetOptionsMarketStatusExchanges verifies that the exchanges nested
// object within the options market status response is correctly parsed
// with the individual NYSE, NASDAQ, and OTC statuses.
func TestGetOptionsMarketStatusExchanges(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/now": optionsMarketStatusJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsMarketStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Exchanges.NYSE != "open" {
		t.Errorf("expected NYSE open, got %s", result.Exchanges.NYSE)
	}

	if result.Exchanges.Nasdaq != "open" {
		t.Errorf("expected NASDAQ open, got %s", result.Exchanges.Nasdaq)
	}

	if result.Exchanges.OTC != "closed" {
		t.Errorf("expected OTC closed, got %s", result.Exchanges.OTC)
	}
}

// TestGetOptionsMarketStatusCurrencies verifies that the currencies nested
// object is parsed correctly with crypto and FX market statuses.
func TestGetOptionsMarketStatusCurrencies(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/now": optionsMarketStatusJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsMarketStatus()
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

// TestGetOptionsMarketStatusIndicesGroups verifies that the indices groups
// nested object is parsed with all index family statuses including S&P,
// Dow Jones, FTSE Russell, and CCCY groups.
func TestGetOptionsMarketStatusIndicesGroups(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/now": optionsMarketStatusJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsMarketStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.IndicesGroups.SAndP != "open" {
		t.Errorf("expected S&P open, got %s", result.IndicesGroups.SAndP)
	}

	if result.IndicesGroups.DowJones != "open" {
		t.Errorf("expected Dow Jones open, got %s", result.IndicesGroups.DowJones)
	}

	if result.IndicesGroups.MStar != "open" {
		t.Errorf("expected MStar open, got %s", result.IndicesGroups.MStar)
	}

	if result.IndicesGroups.CCCY != "open" {
		t.Errorf("expected CCCY open, got %s", result.IndicesGroups.CCCY)
	}

	if result.IndicesGroups.FTSERussell != "open" {
		t.Errorf("expected FTSE Russell open, got %s", result.IndicesGroups.FTSERussell)
	}

	if result.IndicesGroups.MSCI != "closed" {
		t.Errorf("expected MSCI closed, got %s", result.IndicesGroups.MSCI)
	}
}

// TestGetOptionsMarketStatusRequestPath verifies that GetOptionsMarketStatus
// sends the request to the correct API path (/v1/marketstatus/now).
func TestGetOptionsMarketStatusRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsMarketStatusJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsMarketStatus()

	if receivedPath != "/v1/marketstatus/now" {
		t.Errorf("expected path /v1/marketstatus/now, got %s", receivedPath)
	}
}

// TestGetOptionsMarketStatusAPIError verifies that GetOptionsMarketStatus
// returns an error when the API responds with a non-200 status code.
func TestGetOptionsMarketStatusAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"ERROR","message":"Internal server error"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetOptionsMarketStatus()
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

// TestGetOptionsMarketStatusInvalidJSON verifies that GetOptionsMarketStatus
// returns an error when the API responds with malformed JSON.
func TestGetOptionsMarketStatusInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`not valid json`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetOptionsMarketStatus()
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

// TestGetOptionsMarketHolidays verifies that GetOptionsMarketHolidays
// correctly parses the array response containing upcoming market holidays
// with both closed and early-close statuses.
func TestGetOptionsMarketHolidays(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/upcoming": optionsMarketHolidaysJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsMarketHolidays()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 holidays, got %d", len(result))
	}

	first := result[0]
	if first.Date != "2026-04-03" {
		t.Errorf("expected date 2026-04-03, got %s", first.Date)
	}

	if first.Exchange != "NYSE" {
		t.Errorf("expected exchange NYSE, got %s", first.Exchange)
	}

	if first.Name != "Good Friday" {
		t.Errorf("expected name Good Friday, got %s", first.Name)
	}

	if first.Status != "closed" {
		t.Errorf("expected status closed, got %s", first.Status)
	}
}

// TestGetOptionsMarketHolidaysSecondEntry verifies that the second holiday
// in the response is correctly parsed with its distinct exchange value.
func TestGetOptionsMarketHolidaysSecondEntry(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/upcoming": optionsMarketHolidaysJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsMarketHolidays()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := result[1]
	if second.Exchange != "NASDAQ" {
		t.Errorf("expected exchange NASDAQ, got %s", second.Exchange)
	}

	if second.Name != "Good Friday" {
		t.Errorf("expected name Good Friday, got %s", second.Name)
	}

	if second.Status != "closed" {
		t.Errorf("expected status closed, got %s", second.Status)
	}
}

// TestGetOptionsMarketHolidaysEarlyClose verifies that early-close holidays
// include the open and close timestamp fields with their correct values.
func TestGetOptionsMarketHolidaysEarlyClose(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/upcoming": optionsMarketHolidaysJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsMarketHolidays()
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

// TestGetOptionsMarketHolidaysRequestPath verifies that GetOptionsMarketHolidays
// sends the request to the correct API path (/v1/marketstatus/upcoming).
func TestGetOptionsMarketHolidaysRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsMarketHolidaysJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsMarketHolidays()

	if receivedPath != "/v1/marketstatus/upcoming" {
		t.Errorf("expected path /v1/marketstatus/upcoming, got %s", receivedPath)
	}
}

// TestGetOptionsMarketHolidaysAPIError verifies that GetOptionsMarketHolidays
// returns an error when the API responds with a non-200 status code.
func TestGetOptionsMarketHolidaysAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"ERROR","message":"Forbidden"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetOptionsMarketHolidays()
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetOptionsMarketHolidaysInvalidJSON verifies that GetOptionsMarketHolidays
// returns an error when the API responds with malformed JSON.
func TestGetOptionsMarketHolidaysInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{invalid json`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetOptionsMarketHolidays()
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

// TestGetOptionsMarketHolidaysEmptyArray verifies that GetOptionsMarketHolidays
// correctly handles an empty array response with no upcoming holidays.
func TestGetOptionsMarketHolidaysEmptyArray(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v1/marketstatus/upcoming": `[]`,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetOptionsMarketHolidays()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected 0 holidays, got %d", len(result))
	}
}

// TestGetOptionsMarketStatusAPIKeySent verifies that the API key is included
// as a query parameter in the request to the market status endpoint.
func TestGetOptionsMarketStatusAPIKeySent(t *testing.T) {
	var receivedKey string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedKey = r.URL.Query().Get("apiKey")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsMarketStatusJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsMarketStatus()

	if receivedKey != "test-api-key" {
		t.Errorf("expected apiKey=test-api-key, got %s", receivedKey)
	}
}

// TestGetOptionsMarketHolidaysAPIKeySent verifies that the API key is included
// as a query parameter in the request to the market holidays endpoint.
func TestGetOptionsMarketHolidaysAPIKeySent(t *testing.T) {
	var receivedKey string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedKey = r.URL.Query().Get("apiKey")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(optionsMarketHolidaysJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetOptionsMarketHolidays()

	if receivedKey != "test-api-key" {
		t.Errorf("expected apiKey=test-api-key, got %s", receivedKey)
	}
}
