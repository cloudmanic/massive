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

// ---------------------------------------------------------------------------
// Short Interest test fixtures
// ---------------------------------------------------------------------------

const shortInterestJSON = `{
	"status": "OK",
	"request_id": "330e88220c0f4e6b87f6ce7ad6d8a6d6",
	"count": 2,
	"results": [
		{
			"settlement_date": "2025-01-15",
			"ticker": "AAPL",
			"short_interest": 45746430,
			"avg_daily_volume": 23901107,
			"days_to_cover": 1.91
		},
		{
			"settlement_date": "2025-01-31",
			"ticker": "AAPL",
			"short_interest": 42000000,
			"avg_daily_volume": 25000000,
			"days_to_cover": 1.68
		}
	],
	"next_url": "https://api.massive.com/stocks/v1/short-interest?cursor=abc123"
}`

// ---------------------------------------------------------------------------
// Short Volume test fixtures
// ---------------------------------------------------------------------------

const shortVolumeJSON = `{
	"status": "OK",
	"request_id": "3f52bcf40d8748ea9d088f0ca5bbe434",
	"count": 1,
	"results": [
		{
			"ticker": "AAPL",
			"date": "2025-02-06",
			"total_volume": 16264662,
			"short_volume": 5683713,
			"exempt_volume": 67840,
			"non_exempt_volume": 5615873,
			"short_volume_ratio": 34.95,
			"nyse_short_volume": 356029,
			"nyse_short_volume_exempt": 4308,
			"nasdaq_carteret_short_volume": 5298900,
			"nasdaq_carteret_short_volume_exempt": 63532,
			"nasdaq_chicago_short_volume": 28784,
			"nasdaq_chicago_short_volume_exempt": 0,
			"adf_short_volume": 0,
			"adf_short_volume_exempt": 0
		}
	]
}`

// ---------------------------------------------------------------------------
// Float test fixtures
// ---------------------------------------------------------------------------

const floatJSON = `{
	"status": "OK",
	"request_id": "2035d8c427d64a2fb7a625d2b4d7ea40",
	"results": [
		{
			"ticker": "AAPL",
			"free_float": 14831485766,
			"effective_date": "2025-10-03",
			"free_float_percent": 99.9
		}
	]
}`

// ---------------------------------------------------------------------------
// Balance Sheets test fixtures
// ---------------------------------------------------------------------------

const balanceSheetsJSON = `{
	"status": "OK",
	"request_id": "bs-req-001",
	"results": [
		{
			"cik": "0000320193",
			"tickers": ["AAPL"],
			"period_end": "2024-09-28",
			"filing_date": "2024-11-01",
			"fiscal_year": 2024,
			"fiscal_quarter": 4,
			"timeframe": "annual",
			"total_assets": 364980000000,
			"total_current_assets": 152987000000,
			"total_liabilities": 308030000000,
			"total_current_liabilities": 176392000000,
			"total_equity": 56950000000,
			"total_equity_attributable_to_parent": 56950000000,
			"total_liabilities_and_equity": 364980000000,
			"cash_and_equivalents": 29943000000,
			"short_term_investments": 35228000000,
			"receivables": 66243000000,
			"inventories": 7286000000,
			"property_plant_equipment_net": 44856000000,
			"goodwill": 0,
			"retained_earnings_deficit": -19154000000,
			"accumulated_other_comprehensive_income": -7427000000
		}
	]
}`

// ---------------------------------------------------------------------------
// Income Statements test fixtures
// ---------------------------------------------------------------------------

const incomeStatementsJSON = `{
	"status": "OK",
	"request_id": "is-req-001",
	"results": [
		{
			"cik": "0000320193",
			"tickers": ["AAPL"],
			"period_end": "2024-09-28",
			"filing_date": "2024-11-01",
			"fiscal_year": 2024,
			"fiscal_quarter": 4,
			"timeframe": "annual",
			"revenue": 391035000000,
			"cost_of_revenue": 210352000000,
			"gross_profit": 180683000000,
			"operating_income": 123216000000,
			"income_before_income_taxes": 123485000000,
			"income_taxes": 29749000000,
			"consolidated_net_income_loss": 93736000000,
			"basic_earnings_per_share": 6.11,
			"diluted_earnings_per_share": 6.08,
			"basic_shares_outstanding": 15343783000,
			"diluted_shares_outstanding": 15408095000,
			"ebitda": 134661000000,
			"research_development": 31370000000,
			"selling_general_administrative": 26097000000
		}
	]
}`

// ---------------------------------------------------------------------------
// Cash Flow Statements test fixtures
// ---------------------------------------------------------------------------

const cashFlowStatementsJSON = `{
	"status": "OK",
	"request_id": "cf-req-001",
	"results": [
		{
			"cik": "0000320193",
			"tickers": ["AAPL"],
			"period_end": "2024-09-28",
			"filing_date": "2024-11-01",
			"fiscal_year": 2024,
			"fiscal_quarter": 4,
			"timeframe": "annual",
			"net_cash_from_operating_activities": 118254000000,
			"net_cash_from_investing_activities": -7166000000,
			"net_cash_from_financing_activities": -121983000000,
			"change_in_cash_and_equivalents": -10895000000,
			"net_income": 93736000000,
			"depreciation_depletion_and_amortization": 11445000000,
			"purchase_of_property_plant_and_equipment": -9959000000,
			"dividends": -15234000000,
			"long_term_debt_issuances_repayments": -2750000000
		}
	]
}`

// ---------------------------------------------------------------------------
// Financial Ratios test fixtures
// ---------------------------------------------------------------------------

const ratiosJSON = `{
	"status": "OK",
	"request_id": "ratios-req-001",
	"count": 1,
	"results": [
		{
			"ticker": "AAPL",
			"cik": "0000320193",
			"date": "2025-02-14",
			"price": 244.60,
			"market_cap": 3680000000000,
			"earnings_per_share": 6.08,
			"price_to_earnings": 40.23,
			"price_to_book": 64.62,
			"price_to_sales": 9.41,
			"price_to_cash_flow": 31.12,
			"price_to_free_cash_flow": 33.98,
			"dividend_yield": 0.41,
			"return_on_assets": 25.69,
			"return_on_equity": 164.59,
			"debt_to_equity": 1.87,
			"current": 0.87,
			"quick": 0.83,
			"cash": 0.17,
			"ev_to_sales": 9.93,
			"ev_to_ebitda": 29.51,
			"enterprise_value": 3880000000000,
			"free_cash_flow": 108295000000,
			"average_volume": 55000000
		}
	]
}`

// ---------------------------------------------------------------------------
// Short Interest tests
// ---------------------------------------------------------------------------

// TestGetShortInterest verifies that GetShortInterest correctly parses
// the API response and returns the expected short interest records.
func TestGetShortInterest(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/stocks/v1/short-interest": shortInterestJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetShortInterest(ShortInterestParams{
		Ticker: "AAPL",
		Limit:  "2",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Count != 2 {
		t.Errorf("expected count 2, got %d", result.Count)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.Ticker != "AAPL" {
		t.Errorf("expected ticker AAPL, got %s", first.Ticker)
	}

	if first.SettlementDate != "2025-01-15" {
		t.Errorf("expected settlement_date 2025-01-15, got %s", first.SettlementDate)
	}

	if first.ShortInterest != 45746430 {
		t.Errorf("expected short_interest 45746430, got %d", first.ShortInterest)
	}

	if first.AvgDailyVolume != 23901107 {
		t.Errorf("expected avg_daily_volume 23901107, got %d", first.AvgDailyVolume)
	}

	if first.DaysToCover != 1.91 {
		t.Errorf("expected days_to_cover 1.91, got %f", first.DaysToCover)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}
}

// TestGetShortInterestRequestPath verifies that GetShortInterest
// constructs the correct API path.
func TestGetShortInterestRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(shortInterestJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetShortInterest(ShortInterestParams{Ticker: "MSFT"})

	expected := "/stocks/v1/short-interest"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetShortInterestQueryParams verifies that all query parameters
// are correctly sent to the API endpoint.
func TestGetShortInterestQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("ticker") != "TSLA" {
			t.Errorf("expected ticker=TSLA, got %s", q.Get("ticker"))
		}
		if q.Get("settlement_date") != "2025-01-15" {
			t.Errorf("expected settlement_date=2025-01-15, got %s", q.Get("settlement_date"))
		}
		if q.Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "settlement_date.desc" {
			t.Errorf("expected sort=settlement_date.desc, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(shortInterestJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetShortInterest(ShortInterestParams{
		Ticker:         "TSLA",
		SettlementDate: "2025-01-15",
		Limit:          "50",
		Sort:           "settlement_date.desc",
	})
}

// TestGetShortInterestSecondRecord verifies that the second record in
// the response is correctly parsed with its distinct values.
func TestGetShortInterestSecondRecord(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/stocks/v1/short-interest": shortInterestJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetShortInterest(ShortInterestParams{Ticker: "AAPL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := result.Results[1]
	if second.SettlementDate != "2025-01-31" {
		t.Errorf("expected settlement_date 2025-01-31, got %s", second.SettlementDate)
	}

	if second.ShortInterest != 42000000 {
		t.Errorf("expected short_interest 42000000, got %d", second.ShortInterest)
	}

	if second.DaysToCover != 1.68 {
		t.Errorf("expected days_to_cover 1.68, got %f", second.DaysToCover)
	}
}

// TestGetShortInterestAPIError verifies that GetShortInterest returns
// an error when the API responds with a non-200 status.
func TestGetShortInterestAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"NOT_AUTHORIZED","message":"Unauthorized"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetShortInterest(ShortInterestParams{Ticker: "AAPL"})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// ---------------------------------------------------------------------------
// Short Volume tests
// ---------------------------------------------------------------------------

// TestGetShortVolume verifies that GetShortVolume correctly parses the
// API response and returns the expected short volume data.
func TestGetShortVolume(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/stocks/v1/short-volume": shortVolumeJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetShortVolume(ShortVolumeParams{
		Ticker: "AAPL",
		Limit:  "1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Count != 1 {
		t.Errorf("expected count 1, got %d", result.Count)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	sv := result.Results[0]
	if sv.Ticker != "AAPL" {
		t.Errorf("expected ticker AAPL, got %s", sv.Ticker)
	}

	if sv.Date != "2025-02-06" {
		t.Errorf("expected date 2025-02-06, got %s", sv.Date)
	}

	if sv.TotalVolume != 16264662 {
		t.Errorf("expected total_volume 16264662, got %d", sv.TotalVolume)
	}

	if sv.ShortVolume != 5683713 {
		t.Errorf("expected short_volume 5683713, got %d", sv.ShortVolume)
	}

	if sv.ExemptVolume != 67840 {
		t.Errorf("expected exempt_volume 67840, got %d", sv.ExemptVolume)
	}

	if sv.NonExemptVolume != 5615873 {
		t.Errorf("expected non_exempt_volume 5615873, got %d", sv.NonExemptVolume)
	}

	if sv.ShortVolumeRatio != 34.95 {
		t.Errorf("expected short_volume_ratio 34.95, got %f", sv.ShortVolumeRatio)
	}

	if sv.NYSEShortVolume != 356029 {
		t.Errorf("expected nyse_short_volume 356029, got %d", sv.NYSEShortVolume)
	}

	if sv.NasdaqCarteretShortVolume != 5298900 {
		t.Errorf("expected nasdaq_carteret_short_volume 5298900, got %d", sv.NasdaqCarteretShortVolume)
	}

	if sv.NasdaqChicagoShortVolume != 28784 {
		t.Errorf("expected nasdaq_chicago_short_volume 28784, got %d", sv.NasdaqChicagoShortVolume)
	}
}

// TestGetShortVolumeRequestPath verifies that GetShortVolume constructs
// the correct API path.
func TestGetShortVolumeRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(shortVolumeJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetShortVolume(ShortVolumeParams{Ticker: "TSLA"})

	expected := "/stocks/v1/short-volume"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetShortVolumeQueryParams verifies that all query parameters are
// correctly sent to the API endpoint.
func TestGetShortVolumeQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("ticker") != "MSFT" {
			t.Errorf("expected ticker=MSFT, got %s", q.Get("ticker"))
		}
		if q.Get("date") != "2025-02-06" {
			t.Errorf("expected date=2025-02-06, got %s", q.Get("date"))
		}
		if q.Get("limit") != "100" {
			t.Errorf("expected limit=100, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "date.desc" {
			t.Errorf("expected sort=date.desc, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(shortVolumeJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetShortVolume(ShortVolumeParams{
		Ticker: "MSFT",
		Date:   "2025-02-06",
		Limit:  "100",
		Sort:   "date.desc",
	})
}

// ---------------------------------------------------------------------------
// Float tests
// ---------------------------------------------------------------------------

// TestGetFloat verifies that GetFloat correctly parses the API response
// and returns the expected free float data.
func TestGetFloat(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/stocks/vX/float": floatJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetFloat(FloatParams{
		Ticker: "AAPL",
		Limit:  "1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	f := result.Results[0]
	if f.Ticker != "AAPL" {
		t.Errorf("expected ticker AAPL, got %s", f.Ticker)
	}

	if f.EffectiveDate != "2025-10-03" {
		t.Errorf("expected effective_date 2025-10-03, got %s", f.EffectiveDate)
	}

	if f.FreeFloat != 14831485766 {
		t.Errorf("expected free_float 14831485766, got %d", f.FreeFloat)
	}

	if f.FreeFloatPercent != 99.9 {
		t.Errorf("expected free_float_percent 99.9, got %f", f.FreeFloatPercent)
	}
}

// TestGetFloatRequestPath verifies that GetFloat constructs the correct
// API path.
func TestGetFloatRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(floatJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetFloat(FloatParams{Ticker: "AAPL"})

	expected := "/stocks/vX/float"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetFloatQueryParams verifies that all query parameters are
// correctly sent to the API endpoint.
func TestGetFloatQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("ticker") != "TSLA" {
			t.Errorf("expected ticker=TSLA, got %s", q.Get("ticker"))
		}
		if q.Get("limit") != "5" {
			t.Errorf("expected limit=5, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "ticker.desc" {
			t.Errorf("expected sort=ticker.desc, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(floatJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetFloat(FloatParams{
		Ticker: "TSLA",
		Limit:  "5",
		Sort:   "ticker.desc",
	})
}

// TestGetFloatAPIError verifies that GetFloat returns an error when
// the API responds with a non-200 status.
func TestGetFloatAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Data not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetFloat(FloatParams{Ticker: "INVALID"})
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// ---------------------------------------------------------------------------
// Balance Sheets tests
// ---------------------------------------------------------------------------

// TestGetBalanceSheets verifies that GetBalanceSheets correctly parses
// the API response and returns the expected balance sheet data.
func TestGetBalanceSheets(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/stocks/financials/v1/balance-sheets": balanceSheetsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetBalanceSheets(BalanceSheetsParams{
		Tickers:   "AAPL",
		Timeframe: "annual",
		Limit:     "1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	bs := result.Results[0]
	if bs.CIK != "0000320193" {
		t.Errorf("expected CIK 0000320193, got %s", bs.CIK)
	}

	if len(bs.Tickers) != 1 || bs.Tickers[0] != "AAPL" {
		t.Errorf("expected tickers [AAPL], got %v", bs.Tickers)
	}

	if bs.PeriodEnd != "2024-09-28" {
		t.Errorf("expected period_end 2024-09-28, got %s", bs.PeriodEnd)
	}

	if bs.FiscalYear != 2024 {
		t.Errorf("expected fiscal_year 2024, got %d", bs.FiscalYear)
	}

	if bs.FiscalQuarter != 4 {
		t.Errorf("expected fiscal_quarter 4, got %d", bs.FiscalQuarter)
	}

	if bs.Timeframe != "annual" {
		t.Errorf("expected timeframe annual, got %s", bs.Timeframe)
	}

	if bs.TotalAssets != 364980000000 {
		t.Errorf("expected total_assets 364980000000, got %f", bs.TotalAssets)
	}

	if bs.TotalCurrentAssets != 152987000000 {
		t.Errorf("expected total_current_assets 152987000000, got %f", bs.TotalCurrentAssets)
	}

	if bs.TotalLiabilities != 308030000000 {
		t.Errorf("expected total_liabilities 308030000000, got %f", bs.TotalLiabilities)
	}

	if bs.TotalEquity != 56950000000 {
		t.Errorf("expected total_equity 56950000000, got %f", bs.TotalEquity)
	}

	if bs.CashAndEquivalents != 29943000000 {
		t.Errorf("expected cash_and_equivalents 29943000000, got %f", bs.CashAndEquivalents)
	}
}

// TestGetBalanceSheetsRequestPath verifies that GetBalanceSheets
// constructs the correct API path.
func TestGetBalanceSheetsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(balanceSheetsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetBalanceSheets(BalanceSheetsParams{Tickers: "AAPL"})

	expected := "/stocks/financials/v1/balance-sheets"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetBalanceSheetsQueryParams verifies that all query parameters
// are correctly sent to the API endpoint.
func TestGetBalanceSheetsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("tickers") != "MSFT" {
			t.Errorf("expected tickers=MSFT, got %s", q.Get("tickers"))
		}
		if q.Get("timeframe") != "quarterly" {
			t.Errorf("expected timeframe=quarterly, got %s", q.Get("timeframe"))
		}
		if q.Get("limit") != "10" {
			t.Errorf("expected limit=10, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "period_end.desc" {
			t.Errorf("expected sort=period_end.desc, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(balanceSheetsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetBalanceSheets(BalanceSheetsParams{
		Tickers:   "MSFT",
		Timeframe: "quarterly",
		Limit:     "10",
		Sort:      "period_end.desc",
	})
}

// ---------------------------------------------------------------------------
// Income Statements tests
// ---------------------------------------------------------------------------

// TestGetIncomeStatements verifies that GetIncomeStatements correctly
// parses the API response and returns the expected income statement data.
func TestGetIncomeStatements(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/stocks/financials/v1/income-statements": incomeStatementsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetIncomeStatements(IncomeStatementsParams{
		Tickers:   "AAPL",
		Timeframe: "annual",
		Limit:     "1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	is := result.Results[0]
	if is.CIK != "0000320193" {
		t.Errorf("expected CIK 0000320193, got %s", is.CIK)
	}

	if len(is.Tickers) != 1 || is.Tickers[0] != "AAPL" {
		t.Errorf("expected tickers [AAPL], got %v", is.Tickers)
	}

	if is.Revenue != 391035000000 {
		t.Errorf("expected revenue 391035000000, got %f", is.Revenue)
	}

	if is.GrossProfit != 180683000000 {
		t.Errorf("expected gross_profit 180683000000, got %f", is.GrossProfit)
	}

	if is.OperatingIncome != 123216000000 {
		t.Errorf("expected operating_income 123216000000, got %f", is.OperatingIncome)
	}

	if is.ConsolidatedNetIncomeLoss != 93736000000 {
		t.Errorf("expected consolidated_net_income_loss 93736000000, got %f", is.ConsolidatedNetIncomeLoss)
	}

	if is.BasicEarningsPerShare != 6.11 {
		t.Errorf("expected basic_earnings_per_share 6.11, got %f", is.BasicEarningsPerShare)
	}

	if is.DilutedEarningsPerShare != 6.08 {
		t.Errorf("expected diluted_earnings_per_share 6.08, got %f", is.DilutedEarningsPerShare)
	}

	if is.EBITDA != 134661000000 {
		t.Errorf("expected ebitda 134661000000, got %f", is.EBITDA)
	}

	if is.ResearchDevelopment != 31370000000 {
		t.Errorf("expected research_development 31370000000, got %f", is.ResearchDevelopment)
	}
}

// TestGetIncomeStatementsRequestPath verifies that GetIncomeStatements
// constructs the correct API path.
func TestGetIncomeStatementsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(incomeStatementsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIncomeStatements(IncomeStatementsParams{Tickers: "AAPL"})

	expected := "/stocks/financials/v1/income-statements"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetIncomeStatementsQueryParams verifies that all query parameters
// are correctly sent to the API endpoint.
func TestGetIncomeStatementsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("tickers") != "GOOGL" {
			t.Errorf("expected tickers=GOOGL, got %s", q.Get("tickers"))
		}
		if q.Get("cik") != "0001652044" {
			t.Errorf("expected cik=0001652044, got %s", q.Get("cik"))
		}
		if q.Get("timeframe") != "quarterly" {
			t.Errorf("expected timeframe=quarterly, got %s", q.Get("timeframe"))
		}
		if q.Get("limit") != "4" {
			t.Errorf("expected limit=4, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "period_end.desc" {
			t.Errorf("expected sort=period_end.desc, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(incomeStatementsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetIncomeStatements(IncomeStatementsParams{
		Tickers:   "GOOGL",
		CIK:       "0001652044",
		Timeframe: "quarterly",
		Limit:     "4",
		Sort:      "period_end.desc",
	})
}

// ---------------------------------------------------------------------------
// Cash Flow Statements tests
// ---------------------------------------------------------------------------

// TestGetCashFlowStatements verifies that GetCashFlowStatements correctly
// parses the API response and returns the expected cash flow data.
func TestGetCashFlowStatements(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/stocks/financials/v1/cash-flow-statements": cashFlowStatementsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetCashFlowStatements(CashFlowStatementsParams{
		Tickers:   "AAPL",
		Timeframe: "annual",
		Limit:     "1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	cf := result.Results[0]
	if cf.CIK != "0000320193" {
		t.Errorf("expected CIK 0000320193, got %s", cf.CIK)
	}

	if len(cf.Tickers) != 1 || cf.Tickers[0] != "AAPL" {
		t.Errorf("expected tickers [AAPL], got %v", cf.Tickers)
	}

	if cf.Timeframe != "annual" {
		t.Errorf("expected timeframe annual, got %s", cf.Timeframe)
	}

	if cf.NetCashFromOperatingActivities != 118254000000 {
		t.Errorf("expected net_cash_from_operating_activities 118254000000, got %f", cf.NetCashFromOperatingActivities)
	}

	if cf.NetCashFromInvestingActivities != -7166000000 {
		t.Errorf("expected net_cash_from_investing_activities -7166000000, got %f", cf.NetCashFromInvestingActivities)
	}

	if cf.NetCashFromFinancingActivities != -121983000000 {
		t.Errorf("expected net_cash_from_financing_activities -121983000000, got %f", cf.NetCashFromFinancingActivities)
	}

	if cf.ChangeInCashAndEquivalents != -10895000000 {
		t.Errorf("expected change_in_cash_and_equivalents -10895000000, got %f", cf.ChangeInCashAndEquivalents)
	}

	if cf.NetIncome != 93736000000 {
		t.Errorf("expected net_income 93736000000, got %f", cf.NetIncome)
	}

	if cf.DepreciationDepletionAndAmortization != 11445000000 {
		t.Errorf("expected depreciation_depletion_and_amortization 11445000000, got %f", cf.DepreciationDepletionAndAmortization)
	}

	if cf.Dividends != -15234000000 {
		t.Errorf("expected dividends -15234000000, got %f", cf.Dividends)
	}
}

// TestGetCashFlowStatementsRequestPath verifies that GetCashFlowStatements
// constructs the correct API path.
func TestGetCashFlowStatementsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cashFlowStatementsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCashFlowStatements(CashFlowStatementsParams{Tickers: "AAPL"})

	expected := "/stocks/financials/v1/cash-flow-statements"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetCashFlowStatementsQueryParams verifies that all query parameters
// are correctly sent to the API endpoint.
func TestGetCashFlowStatementsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("tickers") != "AMZN" {
			t.Errorf("expected tickers=AMZN, got %s", q.Get("tickers"))
		}
		if q.Get("timeframe") != "quarterly" {
			t.Errorf("expected timeframe=quarterly, got %s", q.Get("timeframe"))
		}
		if q.Get("limit") != "8" {
			t.Errorf("expected limit=8, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "period_end.desc" {
			t.Errorf("expected sort=period_end.desc, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cashFlowStatementsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetCashFlowStatements(CashFlowStatementsParams{
		Tickers:   "AMZN",
		Timeframe: "quarterly",
		Limit:     "8",
		Sort:      "period_end.desc",
	})
}

// ---------------------------------------------------------------------------
// Financial Ratios tests
// ---------------------------------------------------------------------------

// TestGetRatios verifies that GetRatios correctly parses the API
// response and returns the expected financial ratios data.
func TestGetRatios(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/stocks/financials/v1/ratios": ratiosJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetRatios(RatiosParams{
		Ticker: "AAPL",
		Limit:  "1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Count != 1 {
		t.Errorf("expected count 1, got %d", result.Count)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	r := result.Results[0]
	if r.Ticker != "AAPL" {
		t.Errorf("expected ticker AAPL, got %s", r.Ticker)
	}

	if r.CIK != "0000320193" {
		t.Errorf("expected CIK 0000320193, got %s", r.CIK)
	}

	if r.Date != "2025-02-14" {
		t.Errorf("expected date 2025-02-14, got %s", r.Date)
	}

	if r.Price != 244.60 {
		t.Errorf("expected price 244.60, got %f", r.Price)
	}

	if r.MarketCap != 3680000000000 {
		t.Errorf("expected market_cap 3680000000000, got %f", r.MarketCap)
	}

	if r.EarningsPerShare != 6.08 {
		t.Errorf("expected earnings_per_share 6.08, got %f", r.EarningsPerShare)
	}

	if r.PriceToEarnings != 40.23 {
		t.Errorf("expected price_to_earnings 40.23, got %f", r.PriceToEarnings)
	}

	if r.PriceToBook != 64.62 {
		t.Errorf("expected price_to_book 64.62, got %f", r.PriceToBook)
	}

	if r.DividendYield != 0.41 {
		t.Errorf("expected dividend_yield 0.41, got %f", r.DividendYield)
	}

	if r.ReturnOnAssets != 25.69 {
		t.Errorf("expected return_on_assets 25.69, got %f", r.ReturnOnAssets)
	}

	if r.ReturnOnEquity != 164.59 {
		t.Errorf("expected return_on_equity 164.59, got %f", r.ReturnOnEquity)
	}

	if r.DebtToEquity != 1.87 {
		t.Errorf("expected debt_to_equity 1.87, got %f", r.DebtToEquity)
	}

	if r.Current != 0.87 {
		t.Errorf("expected current 0.87, got %f", r.Current)
	}

	if r.EVToEBITDA != 29.51 {
		t.Errorf("expected ev_to_ebitda 29.51, got %f", r.EVToEBITDA)
	}
}

// TestGetRatiosRequestPath verifies that GetRatios constructs the
// correct API path.
func TestGetRatiosRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(ratiosJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetRatios(RatiosParams{Ticker: "AAPL"})

	expected := "/stocks/financials/v1/ratios"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetRatiosQueryParams verifies that all query parameters are
// correctly sent to the API endpoint.
func TestGetRatiosQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("ticker") != "NVDA" {
			t.Errorf("expected ticker=NVDA, got %s", q.Get("ticker"))
		}
		if q.Get("limit") != "25" {
			t.Errorf("expected limit=25, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "date.desc" {
			t.Errorf("expected sort=date.desc, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(ratiosJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetRatios(RatiosParams{
		Ticker: "NVDA",
		Limit:  "25",
		Sort:   "date.desc",
	})
}

// TestGetRatiosAPIError verifies that GetRatios returns an error when
// the API responds with a non-200 status.
func TestGetRatiosAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"NOT_AUTHORIZED","message":"Upgrade plan"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetRatios(RatiosParams{Ticker: "AAPL"})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}
