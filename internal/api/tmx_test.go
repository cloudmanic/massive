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

const corporateEventsJSON = `{
	"status": "OK",
	"count": 3,
	"request_id": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
	"next_url": "https://api.massive.com/tmx/v1/corporate-events?cursor=YWN0aXZlPXRydWUmZGF0ZT0yMDI1",
	"results": [
		{
			"company_name": "Apple Inc.",
			"date": "2025-04-30",
			"isin": "US0378331005",
			"name": "Q2 2025 Earnings",
			"status": "confirmed",
			"ticker": "AAPL",
			"tmx_company_id": 12345,
			"tmx_record_id": "REC-001-AAPL",
			"trading_venue": "XNAS",
			"type": "earnings_announcement_date",
			"url": "https://example.com/aapl-earnings"
		},
		{
			"company_name": "Microsoft Corporation",
			"date": "2025-05-15",
			"isin": "US5949181045",
			"name": "Quarterly Dividend",
			"status": "approved",
			"ticker": "MSFT",
			"tmx_company_id": 67890,
			"tmx_record_id": "REC-002-MSFT",
			"trading_venue": "XNAS",
			"type": "dividend",
			"url": "https://example.com/msft-dividend"
		},
		{
			"company_name": "Tesla Inc.",
			"date": "2025-06-01",
			"isin": "US88160R1014",
			"name": "Annual Shareholder Meeting",
			"status": "unconfirmed",
			"ticker": "TSLA",
			"tmx_company_id": 11223,
			"tmx_record_id": "REC-003-TSLA",
			"trading_venue": "XNAS",
			"type": "investor_conference",
			"url": "https://example.com/tsla-meeting"
		}
	]
}`

// TestGetTMXCorporateEvents verifies that GetTMXCorporateEvents correctly
// parses the API response and returns the expected corporate event data
// including all fields such as company name, date, event type, and status.
func TestGetTMXCorporateEvents(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/tmx/v1/corporate-events": corporateEventsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := TMXCorporateEventsParams{
		Ticker: "AAPL",
		Limit:  "3",
	}

	result, err := client.GetTMXCorporateEvents(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Count != 3 {
		t.Errorf("expected count 3, got %d", result.Count)
	}

	if result.RequestID != "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6" {
		t.Errorf("expected request_id a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6, got %s", result.RequestID)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 3 {
		t.Fatalf("expected 3 events, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.CompanyName != "Apple Inc." {
		t.Errorf("expected company_name Apple Inc., got %s", first.CompanyName)
	}

	if first.Date != "2025-04-30" {
		t.Errorf("expected date 2025-04-30, got %s", first.Date)
	}

	if first.ISIN != "US0378331005" {
		t.Errorf("expected isin US0378331005, got %s", first.ISIN)
	}

	if first.Name != "Q2 2025 Earnings" {
		t.Errorf("expected name Q2 2025 Earnings, got %s", first.Name)
	}

	if first.Status != "confirmed" {
		t.Errorf("expected status confirmed, got %s", first.Status)
	}

	if first.Ticker != "AAPL" {
		t.Errorf("expected ticker AAPL, got %s", first.Ticker)
	}

	if first.TMXCompanyID != 12345 {
		t.Errorf("expected tmx_company_id 12345, got %d", first.TMXCompanyID)
	}

	if first.TMXRecordID != "REC-001-AAPL" {
		t.Errorf("expected tmx_record_id REC-001-AAPL, got %s", first.TMXRecordID)
	}

	if first.TradingVenue != "XNAS" {
		t.Errorf("expected trading_venue XNAS, got %s", first.TradingVenue)
	}

	if first.Type != "earnings_announcement_date" {
		t.Errorf("expected type earnings_announcement_date, got %s", first.Type)
	}

	if first.URL != "https://example.com/aapl-earnings" {
		t.Errorf("expected url https://example.com/aapl-earnings, got %s", first.URL)
	}
}

// TestGetTMXCorporateEventsSecondResult verifies that the second event
// in the response is correctly parsed with its own distinct values
// including the dividend type and approved status.
func TestGetTMXCorporateEventsSecondResult(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/tmx/v1/corporate-events": corporateEventsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetTMXCorporateEvents(TMXCorporateEventsParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := result.Results[1]
	if second.CompanyName != "Microsoft Corporation" {
		t.Errorf("expected company_name Microsoft Corporation, got %s", second.CompanyName)
	}

	if second.Date != "2025-05-15" {
		t.Errorf("expected date 2025-05-15, got %s", second.Date)
	}

	if second.Ticker != "MSFT" {
		t.Errorf("expected ticker MSFT, got %s", second.Ticker)
	}

	if second.Type != "dividend" {
		t.Errorf("expected type dividend, got %s", second.Type)
	}

	if second.Status != "approved" {
		t.Errorf("expected status approved, got %s", second.Status)
	}

	if second.TMXCompanyID != 67890 {
		t.Errorf("expected tmx_company_id 67890, got %d", second.TMXCompanyID)
	}

	if second.ISIN != "US5949181045" {
		t.Errorf("expected isin US5949181045, got %s", second.ISIN)
	}
}

// TestGetTMXCorporateEventsThirdResult verifies that the third event
// in the response is correctly parsed with its own distinct values
// including the investor_conference type and unconfirmed status.
func TestGetTMXCorporateEventsThirdResult(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/tmx/v1/corporate-events": corporateEventsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetTMXCorporateEvents(TMXCorporateEventsParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	third := result.Results[2]
	if third.CompanyName != "Tesla Inc." {
		t.Errorf("expected company_name Tesla Inc., got %s", third.CompanyName)
	}

	if third.Ticker != "TSLA" {
		t.Errorf("expected ticker TSLA, got %s", third.Ticker)
	}

	if third.Type != "investor_conference" {
		t.Errorf("expected type investor_conference, got %s", third.Type)
	}

	if third.Status != "unconfirmed" {
		t.Errorf("expected status unconfirmed, got %s", third.Status)
	}

	if third.TMXCompanyID != 11223 {
		t.Errorf("expected tmx_company_id 11223, got %d", third.TMXCompanyID)
	}
}

// TestGetTMXCorporateEventsRequestPath verifies that GetTMXCorporateEvents
// constructs the correct API path for the corporate events endpoint.
func TestGetTMXCorporateEventsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(corporateEventsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetTMXCorporateEvents(TMXCorporateEventsParams{Ticker: "AAPL"})

	if receivedPath != "/tmx/v1/corporate-events" {
		t.Errorf("expected path /tmx/v1/corporate-events, got %s", receivedPath)
	}
}

// TestGetTMXCorporateEventsQueryParams verifies that all filter parameters
// are correctly sent to the API endpoint as query parameters, including
// ticker, date range filters, type, status, and limit.
func TestGetTMXCorporateEventsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("ticker") != "AAPL" {
			t.Errorf("expected ticker=AAPL, got %s", q.Get("ticker"))
		}
		if q.Get("date.gte") != "2025-01-01" {
			t.Errorf("expected date.gte=2025-01-01, got %s", q.Get("date.gte"))
		}
		if q.Get("date.lte") != "2025-12-31" {
			t.Errorf("expected date.lte=2025-12-31, got %s", q.Get("date.lte"))
		}
		if q.Get("type") != "earnings_announcement_date" {
			t.Errorf("expected type=earnings_announcement_date, got %s", q.Get("type"))
		}
		if q.Get("status") != "confirmed" {
			t.Errorf("expected status=confirmed, got %s", q.Get("status"))
		}
		if q.Get("sort") != "date.desc" {
			t.Errorf("expected sort=date.desc, got %s", q.Get("sort"))
		}
		if q.Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(corporateEventsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetTMXCorporateEvents(TMXCorporateEventsParams{
		Ticker:  "AAPL",
		DateGTE: "2025-01-01",
		DateLTE: "2025-12-31",
		Type:    "earnings_announcement_date",
		Status:  "confirmed",
		Sort:    "date.desc",
		Limit:   "50",
	})
}

// TestGetTMXCorporateEventsAdvancedParams verifies that advanced filter
// parameters like ISIN, trading venue, any_of, and TMX-specific IDs
// are correctly sent to the API endpoint.
func TestGetTMXCorporateEventsAdvancedParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("isin") != "US0378331005" {
			t.Errorf("expected isin=US0378331005, got %s", q.Get("isin"))
		}
		if q.Get("trading_venue") != "XNAS" {
			t.Errorf("expected trading_venue=XNAS, got %s", q.Get("trading_venue"))
		}
		if q.Get("type.any_of") != "earnings_announcement_date,dividend" {
			t.Errorf("expected type.any_of=earnings_announcement_date,dividend, got %s", q.Get("type.any_of"))
		}
		if q.Get("tmx_company_id") != "12345" {
			t.Errorf("expected tmx_company_id=12345, got %s", q.Get("tmx_company_id"))
		}
		if q.Get("tmx_record_id") != "REC-001-AAPL" {
			t.Errorf("expected tmx_record_id=REC-001-AAPL, got %s", q.Get("tmx_record_id"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(corporateEventsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetTMXCorporateEvents(TMXCorporateEventsParams{
		ISIN:         "US0378331005",
		TradingVenue: "XNAS",
		TypeAnyOf:    "earnings_announcement_date,dividend",
		TMXCompanyID: "12345",
		TMXRecordID:  "REC-001-AAPL",
	})
}

// TestGetTMXCorporateEventsAPIError verifies that GetTMXCorporateEvents
// returns an error when the API responds with a non-200 status code.
func TestGetTMXCorporateEventsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"NOT_AUTHORIZED","message":"You are not entitled to this data."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetTMXCorporateEvents(TMXCorporateEventsParams{Ticker: "AAPL"})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetTMXCorporateEventsEmptyResults verifies that GetTMXCorporateEvents
// handles an empty results array without error when no events match the
// filter criteria.
func TestGetTMXCorporateEventsEmptyResults(t *testing.T) {
	emptyJSON := `{"status":"OK","count":0,"request_id":"abc123","results":[]}`
	server := mockServer(t, map[string]string{
		"/tmx/v1/corporate-events": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetTMXCorporateEvents(TMXCorporateEventsParams{
		Ticker: "ZZZZNOTREAL",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Count != 0 {
		t.Errorf("expected count 0, got %d", result.Count)
	}

	if len(result.Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(result.Results))
	}
}

// TestGetTMXCorporateEventsDateRangeParams verifies that date range
// comparison operators (gt, gte, lt, lte) are correctly sent as
// query parameters to the API.
func TestGetTMXCorporateEventsDateRangeParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("date.gt") != "2025-01-01" {
			t.Errorf("expected date.gt=2025-01-01, got %s", q.Get("date.gt"))
		}
		if q.Get("date.lt") != "2025-06-30" {
			t.Errorf("expected date.lt=2025-06-30, got %s", q.Get("date.lt"))
		}
		if q.Get("date") != "2025-03-15" {
			t.Errorf("expected date=2025-03-15, got %s", q.Get("date"))
		}
		if q.Get("date.any_of") != "2025-03-15,2025-04-15" {
			t.Errorf("expected date.any_of=2025-03-15,2025-04-15, got %s", q.Get("date.any_of"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(corporateEventsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetTMXCorporateEvents(TMXCorporateEventsParams{
		Date:      "2025-03-15",
		DateAnyOf: "2025-03-15,2025-04-15",
		DateGT:    "2025-01-01",
		DateLT:    "2025-06-30",
	})
}

// TestGetTMXCorporateEventsStatusParams verifies that status filtering
// parameters including exact match, any_of, and range operators are
// correctly sent as query parameters.
func TestGetTMXCorporateEventsStatusParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("status.any_of") != "confirmed,approved" {
			t.Errorf("expected status.any_of=confirmed,approved, got %s", q.Get("status.any_of"))
		}
		if q.Get("ticker.any_of") != "AAPL,MSFT,TSLA" {
			t.Errorf("expected ticker.any_of=AAPL,MSFT,TSLA, got %s", q.Get("ticker.any_of"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(corporateEventsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetTMXCorporateEvents(TMXCorporateEventsParams{
		StatusAnyOf: "confirmed,approved",
		TickerAnyOf: "AAPL,MSFT,TSLA",
	})
}
