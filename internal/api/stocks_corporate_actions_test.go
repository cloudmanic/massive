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

const dividendsJSON = `{
	"status": "OK",
	"request_id": "5a67a998c03f40dabf20fdc9c1fe6079",
	"results": [
		{
			"id": "E2ab18998fec423fdb40a4f3bdb6df573c30c51c00039eb26bcda869f951ea88f",
			"ticker": "AAPL",
			"record_date": "2012-08-13",
			"pay_date": "2012-08-16",
			"ex_dividend_date": "2012-08-09",
			"frequency": 0,
			"cash_amount": 2.65,
			"currency": "USD",
			"distribution_type": "unknown",
			"historical_adjustment_factor": 0.838964,
			"split_adjusted_cash_amount": 0.094643
		},
		{
			"id": "E11cfbce7b91c73bba1c5601ebf5e9eee8029555e840c84ff7e95545acdf45d71",
			"ticker": "AAPL",
			"record_date": "2012-11-12",
			"pay_date": "2012-11-15",
			"declaration_date": "2012-10-25",
			"ex_dividend_date": "2012-11-07",
			"frequency": 0,
			"cash_amount": 2.65,
			"currency": "USD",
			"distribution_type": "unknown",
			"historical_adjustment_factor": 0.842566,
			"split_adjusted_cash_amount": 0.094643
		}
	],
	"next_url": "https://api.massive.com/stocks/v1/dividends?cursor=AQwPBEFBUEwCAQABBQABAQIAAgEPBEFBUEwBCftGeg=="
}`

const splitsJSON = `{
	"status": "OK",
	"request_id": "36538f14c2ee4f98b68bb4d968a85be4",
	"results": [
		{
			"id": "E36416cce743c3964c5da63e1ef1626c0aece30fb47302eea5a49c0055c04e8d0",
			"execution_date": "2020-08-31",
			"split_from": 1.0,
			"split_to": 4.0,
			"ticker": "AAPL",
			"adjustment_type": "forward_split",
			"historical_adjustment_factor": 0.25
		},
		{
			"id": "E91a6b74ca1a9dcbce26a1f34e24ae26ba2c6359822ccf901ecd827f419137654",
			"execution_date": "2014-06-09",
			"split_from": 1.0,
			"split_to": 7.0,
			"ticker": "AAPL",
			"adjustment_type": "forward_split",
			"historical_adjustment_factor": 0.035714
		}
	],
	"next_url": "https://api.massive.com/stocks/v1/splits?cursor=AQcPBEFBUEwCAQAABAAAAQIAAgEJ-8x-AQ8EQUFQTA=="
}`

// TestGetDividends verifies that GetDividends correctly parses the API
// response and returns the expected dividend data for AAPL including
// all fields such as cash amount, dates, and adjustment factors.
func TestGetDividends(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/stocks/v1/dividends": dividendsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := DividendsParams{
		Ticker: "AAPL",
		Limit:  "2",
	}

	result, err := client.GetDividends(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "5a67a998c03f40dabf20fdc9c1fe6079" {
		t.Errorf("expected request_id 5a67a998c03f40dabf20fdc9c1fe6079, got %s", result.RequestID)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 dividends, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.Ticker != "AAPL" {
		t.Errorf("expected ticker AAPL, got %s", first.Ticker)
	}

	if first.ExDividendDate != "2012-08-09" {
		t.Errorf("expected ex_dividend_date 2012-08-09, got %s", first.ExDividendDate)
	}

	if first.RecordDate != "2012-08-13" {
		t.Errorf("expected record_date 2012-08-13, got %s", first.RecordDate)
	}

	if first.PayDate != "2012-08-16" {
		t.Errorf("expected pay_date 2012-08-16, got %s", first.PayDate)
	}

	if first.CashAmount != 2.65 {
		t.Errorf("expected cash_amount 2.65, got %f", first.CashAmount)
	}

	if first.Currency != "USD" {
		t.Errorf("expected currency USD, got %s", first.Currency)
	}

	if first.DistributionType != "unknown" {
		t.Errorf("expected distribution_type unknown, got %s", first.DistributionType)
	}

	if first.Frequency != 0 {
		t.Errorf("expected frequency 0, got %d", first.Frequency)
	}

	if first.HistoricalAdjustmentFactor != 0.838964 {
		t.Errorf("expected historical_adjustment_factor 0.838964, got %f", first.HistoricalAdjustmentFactor)
	}

	if first.SplitAdjustedCashAmount != 0.094643 {
		t.Errorf("expected split_adjusted_cash_amount 0.094643, got %f", first.SplitAdjustedCashAmount)
	}

	if first.DeclarationDate != "" {
		t.Errorf("expected empty declaration_date for first dividend, got %s", first.DeclarationDate)
	}
}

// TestGetDividendsSecondResult verifies that the second dividend in the
// response is correctly parsed with its own distinct values, including
// the declaration_date field which is present on this record.
func TestGetDividendsSecondResult(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/stocks/v1/dividends": dividendsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetDividends(DividendsParams{Ticker: "AAPL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := result.Results[1]
	if second.DeclarationDate != "2012-10-25" {
		t.Errorf("expected declaration_date 2012-10-25, got %s", second.DeclarationDate)
	}

	if second.ExDividendDate != "2012-11-07" {
		t.Errorf("expected ex_dividend_date 2012-11-07, got %s", second.ExDividendDate)
	}

	if second.RecordDate != "2012-11-12" {
		t.Errorf("expected record_date 2012-11-12, got %s", second.RecordDate)
	}

	if second.PayDate != "2012-11-15" {
		t.Errorf("expected pay_date 2012-11-15, got %s", second.PayDate)
	}

	if second.HistoricalAdjustmentFactor != 0.842566 {
		t.Errorf("expected historical_adjustment_factor 0.842566, got %f", second.HistoricalAdjustmentFactor)
	}
}

// TestGetDividendsRequestPath verifies that GetDividends constructs
// the correct API path for the dividends endpoint.
func TestGetDividendsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(dividendsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetDividends(DividendsParams{Ticker: "MSFT"})

	if receivedPath != "/stocks/v1/dividends" {
		t.Errorf("expected path /stocks/v1/dividends, got %s", receivedPath)
	}
}

// TestGetDividendsQueryParams verifies that all filter parameters are
// correctly sent to the API endpoint as query parameters, including
// ticker, date range filters, frequency, and distribution type.
func TestGetDividendsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("ticker") != "AAPL" {
			t.Errorf("expected ticker=AAPL, got %s", q.Get("ticker"))
		}
		if q.Get("ex_dividend_date.gte") != "2024-01-01" {
			t.Errorf("expected ex_dividend_date.gte=2024-01-01, got %s", q.Get("ex_dividend_date.gte"))
		}
		if q.Get("ex_dividend_date.lte") != "2024-12-31" {
			t.Errorf("expected ex_dividend_date.lte=2024-12-31, got %s", q.Get("ex_dividend_date.lte"))
		}
		if q.Get("frequency") != "4" {
			t.Errorf("expected frequency=4, got %s", q.Get("frequency"))
		}
		if q.Get("distribution_type") != "recurring" {
			t.Errorf("expected distribution_type=recurring, got %s", q.Get("distribution_type"))
		}
		if q.Get("sort") != "ex_dividend_date.desc" {
			t.Errorf("expected sort=ex_dividend_date.desc, got %s", q.Get("sort"))
		}
		if q.Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(dividendsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetDividends(DividendsParams{
		Ticker:            "AAPL",
		ExDividendDateGTE: "2024-01-01",
		ExDividendDateLTE: "2024-12-31",
		Frequency:         "4",
		DistributionType:  "recurring",
		Sort:              "ex_dividend_date.desc",
		Limit:             "50",
	})
}

// TestGetDividendsAPIError verifies that GetDividends returns an error
// when the API responds with a non-200 status code.
func TestGetDividendsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"ERROR","message":"Access denied."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetDividends(DividendsParams{Ticker: "AAPL"})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetDividendsEmptyResults verifies that GetDividends handles an
// empty results array without error when no dividends match the filter.
func TestGetDividendsEmptyResults(t *testing.T) {
	emptyJSON := `{"status":"OK","request_id":"abc123","results":[]}`
	server := mockServer(t, map[string]string{
		"/stocks/v1/dividends": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetDividends(DividendsParams{Ticker: "ZZZZNOTREAL"})
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

// TestGetSplits verifies that GetSplits correctly parses the API
// response and returns the expected stock split data for AAPL including
// execution dates, split ratios, and adjustment factors.
func TestGetSplits(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/stocks/v1/splits": splitsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := SplitsParams{
		Ticker: "AAPL",
		Limit:  "2",
	}

	result, err := client.GetSplits(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "36538f14c2ee4f98b68bb4d968a85be4" {
		t.Errorf("expected request_id 36538f14c2ee4f98b68bb4d968a85be4, got %s", result.RequestID)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 splits, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.Ticker != "AAPL" {
		t.Errorf("expected ticker AAPL, got %s", first.Ticker)
	}

	if first.ExecutionDate != "2020-08-31" {
		t.Errorf("expected execution_date 2020-08-31, got %s", first.ExecutionDate)
	}

	if first.SplitFrom != 1.0 {
		t.Errorf("expected split_from 1.0, got %f", first.SplitFrom)
	}

	if first.SplitTo != 4.0 {
		t.Errorf("expected split_to 4.0, got %f", first.SplitTo)
	}

	if first.AdjustmentType != "forward_split" {
		t.Errorf("expected adjustment_type forward_split, got %s", first.AdjustmentType)
	}

	if first.HistoricalAdjustmentFactor != 0.25 {
		t.Errorf("expected historical_adjustment_factor 0.25, got %f", first.HistoricalAdjustmentFactor)
	}
}

// TestGetSplitsSecondResult verifies that the second split in the
// response is correctly parsed with its own distinct values including
// the 7-for-1 split ratio.
func TestGetSplitsSecondResult(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/stocks/v1/splits": splitsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetSplits(SplitsParams{Ticker: "AAPL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := result.Results[1]
	if second.ExecutionDate != "2014-06-09" {
		t.Errorf("expected execution_date 2014-06-09, got %s", second.ExecutionDate)
	}

	if second.SplitFrom != 1.0 {
		t.Errorf("expected split_from 1.0, got %f", second.SplitFrom)
	}

	if second.SplitTo != 7.0 {
		t.Errorf("expected split_to 7.0, got %f", second.SplitTo)
	}

	if second.AdjustmentType != "forward_split" {
		t.Errorf("expected adjustment_type forward_split, got %s", second.AdjustmentType)
	}

	if second.HistoricalAdjustmentFactor != 0.035714 {
		t.Errorf("expected historical_adjustment_factor 0.035714, got %f", second.HistoricalAdjustmentFactor)
	}
}

// TestGetSplitsRequestPath verifies that GetSplits constructs the
// correct API path for the splits endpoint.
func TestGetSplitsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(splitsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetSplits(SplitsParams{Ticker: "TSLA"})

	if receivedPath != "/stocks/v1/splits" {
		t.Errorf("expected path /stocks/v1/splits, got %s", receivedPath)
	}
}

// TestGetSplitsQueryParams verifies that all filter parameters are
// correctly sent to the API endpoint as query parameters, including
// ticker, date range filters, and adjustment type.
func TestGetSplitsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("ticker") != "AAPL" {
			t.Errorf("expected ticker=AAPL, got %s", q.Get("ticker"))
		}
		if q.Get("execution_date.gte") != "2020-01-01" {
			t.Errorf("expected execution_date.gte=2020-01-01, got %s", q.Get("execution_date.gte"))
		}
		if q.Get("execution_date.lte") != "2025-12-31" {
			t.Errorf("expected execution_date.lte=2025-12-31, got %s", q.Get("execution_date.lte"))
		}
		if q.Get("adjustment_type") != "forward_split" {
			t.Errorf("expected adjustment_type=forward_split, got %s", q.Get("adjustment_type"))
		}
		if q.Get("sort") != "execution_date.desc" {
			t.Errorf("expected sort=execution_date.desc, got %s", q.Get("sort"))
		}
		if q.Get("limit") != "100" {
			t.Errorf("expected limit=100, got %s", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(splitsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetSplits(SplitsParams{
		Ticker:           "AAPL",
		ExecutionDateGTE: "2020-01-01",
		ExecutionDateLTE: "2025-12-31",
		AdjustmentType:   "forward_split",
		Sort:             "execution_date.desc",
		Limit:            "100",
	})
}

// TestGetSplitsAPIError verifies that GetSplits returns an error
// when the API responds with a non-200 status code.
func TestGetSplitsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"ERROR","message":"Access denied."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetSplits(SplitsParams{Ticker: "AAPL"})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetSplitsEmptyResults verifies that GetSplits handles an empty
// results array without error when no splits match the filter criteria.
func TestGetSplitsEmptyResults(t *testing.T) {
	emptyJSON := `{"status":"OK","request_id":"abc123","results":[]}`
	server := mockServer(t, map[string]string{
		"/stocks/v1/splits": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetSplits(SplitsParams{Ticker: "ZZZZNOTREAL"})
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
