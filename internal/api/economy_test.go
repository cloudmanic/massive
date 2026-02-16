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

const inflationJSON = `{
	"status": "OK",
	"request_id": "60fa9df936064cb79eee45e5c5ed4e2a",
	"results": [
		{
			"date": "2025-11-01",
			"cpi": 325.063,
			"cpi_core": 331.043,
			"pce": 128.093,
			"pce_core": 127.422,
			"pce_spending": 21409.7
		},
		{
			"date": "2025-09-01",
			"cpi": 324.245,
			"cpi_core": 330.418,
			"pce": 127.625,
			"pce_core": 126.954,
			"pce_spending": 21202.4
		}
	],
	"next_url": "https://api.massive.com/fed/v1/inflation?cursor=AAEAAAABAgABAQ8KMjAyNS0wOS0wMQ=="
}`

const laborMarketJSON = `{
	"status": "OK",
	"request_id": "7001492ee53347d48ee3e43992e9a5dd",
	"results": [
		{
			"date": "2026-01-01",
			"unemployment_rate": 4.3,
			"labor_force_participation_rate": 62.5,
			"avg_hourly_earnings": 37.17
		},
		{
			"date": "2025-12-01",
			"unemployment_rate": 4.4,
			"labor_force_participation_rate": 62.4,
			"avg_hourly_earnings": 37.02,
			"job_openings": 6542.0
		}
	],
	"next_url": "https://api.massive.com/fed/v1/labor-market?cursor=AAEAAAABAgABAQ8KMjAyNS0xMi0wMQ=="
}`

const treasuryYieldsJSON = `{
	"status": "OK",
	"request_id": "c33732ef5b614ef7990f3899dbbc6da2",
	"results": [
		{
			"date": "2026-02-12",
			"yield_1_month": 3.72,
			"yield_3_month": 3.70,
			"yield_1_year": 3.45,
			"yield_2_year": 3.47,
			"yield_5_year": 3.67,
			"yield_10_year": 4.09,
			"yield_30_year": 4.72
		},
		{
			"date": "2026-02-11",
			"yield_1_month": 3.71,
			"yield_3_month": 3.70,
			"yield_1_year": 3.47,
			"yield_2_year": 3.52,
			"yield_5_year": 3.75,
			"yield_10_year": 4.18,
			"yield_30_year": 4.82
		}
	],
	"next_url": "https://api.massive.com/fed/v1/treasury-yields?cursor=AAEAAAABAgABAQ8KMjAyNi0wMi0xMQ=="
}`

// TestGetInflation verifies that GetInflation correctly parses the API
// response and returns the expected inflation data with CPI and PCE fields.
func TestGetInflation(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/fed/v1/inflation": inflationJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetInflation(InflationParams{
		Sort:  "date.desc",
		Limit: "2",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "60fa9df936064cb79eee45e5c5ed4e2a" {
		t.Errorf("expected request_id 60fa9df936064cb79eee45e5c5ed4e2a, got %s", result.RequestID)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.Date != "2025-11-01" {
		t.Errorf("expected date 2025-11-01, got %s", first.Date)
	}

	if first.CPI != 325.063 {
		t.Errorf("expected CPI 325.063, got %f", first.CPI)
	}

	if first.CPICore != 331.043 {
		t.Errorf("expected CPI core 331.043, got %f", first.CPICore)
	}

	if first.PCE != 128.093 {
		t.Errorf("expected PCE 128.093, got %f", first.PCE)
	}

	if first.PCECore != 127.422 {
		t.Errorf("expected PCE core 127.422, got %f", first.PCECore)
	}

	if first.PCESpending != 21409.7 {
		t.Errorf("expected PCE spending 21409.7, got %f", first.PCESpending)
	}
}

// TestGetInflationRequestPath verifies that GetInflation constructs the
// correct API path for the inflation endpoint.
func TestGetInflationRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(inflationJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetInflation(InflationParams{})

	if receivedPath != "/fed/v1/inflation" {
		t.Errorf("expected path /fed/v1/inflation, got %s", receivedPath)
	}
}

// TestGetInflationQueryParams verifies that all filter parameters including
// date range operators, sort, and limit are correctly sent to the API.
func TestGetInflationQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("date.gte") != "2025-01-01" {
			t.Errorf("expected date.gte=2025-01-01, got %s", q.Get("date.gte"))
		}
		if q.Get("date.lte") != "2025-12-31" {
			t.Errorf("expected date.lte=2025-12-31, got %s", q.Get("date.lte"))
		}
		if q.Get("sort") != "date.asc" {
			t.Errorf("expected sort=date.asc, got %s", q.Get("sort"))
		}
		if q.Get("limit") != "10" {
			t.Errorf("expected limit=10, got %s", q.Get("limit"))
		}
		if q.Get("apiKey") != "test-api-key" {
			t.Errorf("expected apiKey=test-api-key, got %s", q.Get("apiKey"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(inflationJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetInflation(InflationParams{
		DateGTE: "2025-01-01",
		DateLTE: "2025-12-31",
		Sort:    "date.asc",
		Limit:   "10",
	})
}

// TestGetInflationDateFilter verifies that the exact date parameter is
// correctly sent when filtering for a specific observation date.
func TestGetInflationDateFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("date") != "2025-06-01" {
			t.Errorf("expected date=2025-06-01, got %s", q.Get("date"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(inflationJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetInflation(InflationParams{Date: "2025-06-01"})
}

// TestGetInflationAPIError verifies that GetInflation returns an error
// when the API responds with a non-200 status code.
func TestGetInflationAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"ERROR","message":"Forbidden"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetInflation(InflationParams{})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetInflationSecondResult verifies that the second result in the
// inflation response is correctly parsed with its own distinct values.
func TestGetInflationSecondResult(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/fed/v1/inflation": inflationJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetInflation(InflationParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := result.Results[1]
	if second.Date != "2025-09-01" {
		t.Errorf("expected date 2025-09-01, got %s", second.Date)
	}

	if second.CPI != 324.245 {
		t.Errorf("expected CPI 324.245, got %f", second.CPI)
	}

	if second.PCESpending != 21202.4 {
		t.Errorf("expected PCE spending 21202.4, got %f", second.PCESpending)
	}
}

// TestGetLaborMarket verifies that GetLaborMarket correctly parses the API
// response and returns the expected labor market indicator data.
func TestGetLaborMarket(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/fed/v1/labor-market": laborMarketJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetLaborMarket(LaborMarketParams{
		Sort:  "date.desc",
		Limit: "2",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "7001492ee53347d48ee3e43992e9a5dd" {
		t.Errorf("expected request_id 7001492ee53347d48ee3e43992e9a5dd, got %s", result.RequestID)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.Date != "2026-01-01" {
		t.Errorf("expected date 2026-01-01, got %s", first.Date)
	}

	if first.UnemploymentRate != 4.3 {
		t.Errorf("expected unemployment rate 4.3, got %f", first.UnemploymentRate)
	}

	if first.LaborForceParticipationRate != 62.5 {
		t.Errorf("expected labor force participation rate 62.5, got %f", first.LaborForceParticipationRate)
	}

	if first.AvgHourlyEarnings != 37.17 {
		t.Errorf("expected avg hourly earnings 37.17, got %f", first.AvgHourlyEarnings)
	}
}

// TestGetLaborMarketRequestPath verifies that GetLaborMarket constructs
// the correct API path for the labor market endpoint.
func TestGetLaborMarketRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(laborMarketJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetLaborMarket(LaborMarketParams{})

	if receivedPath != "/fed/v1/labor-market" {
		t.Errorf("expected path /fed/v1/labor-market, got %s", receivedPath)
	}
}

// TestGetLaborMarketQueryParams verifies that all filter parameters including
// date range operators, sort, and limit are correctly sent to the API.
func TestGetLaborMarketQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("date.gt") != "2025-01-01" {
			t.Errorf("expected date.gt=2025-01-01, got %s", q.Get("date.gt"))
		}
		if q.Get("date.lt") != "2025-12-31" {
			t.Errorf("expected date.lt=2025-12-31, got %s", q.Get("date.lt"))
		}
		if q.Get("sort") != "date.desc" {
			t.Errorf("expected sort=date.desc, got %s", q.Get("sort"))
		}
		if q.Get("limit") != "5" {
			t.Errorf("expected limit=5, got %s", q.Get("limit"))
		}
		if q.Get("apiKey") != "test-api-key" {
			t.Errorf("expected apiKey=test-api-key, got %s", q.Get("apiKey"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(laborMarketJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetLaborMarket(LaborMarketParams{
		DateGT: "2025-01-01",
		DateLT: "2025-12-31",
		Sort:   "date.desc",
		Limit:  "5",
	})
}

// TestGetLaborMarketSecondResult verifies that the second result in the
// labor market response is correctly parsed including job openings data.
func TestGetLaborMarketSecondResult(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/fed/v1/labor-market": laborMarketJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetLaborMarket(LaborMarketParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := result.Results[1]
	if second.Date != "2025-12-01" {
		t.Errorf("expected date 2025-12-01, got %s", second.Date)
	}

	if second.UnemploymentRate != 4.4 {
		t.Errorf("expected unemployment rate 4.4, got %f", second.UnemploymentRate)
	}

	if second.AvgHourlyEarnings != 37.02 {
		t.Errorf("expected avg hourly earnings 37.02, got %f", second.AvgHourlyEarnings)
	}

	if second.JobOpenings != 6542.0 {
		t.Errorf("expected job openings 6542, got %f", second.JobOpenings)
	}
}

// TestGetLaborMarketAPIError verifies that GetLaborMarket returns an error
// when the API responds with a non-200 status code.
func TestGetLaborMarketAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"ERROR","message":"Internal Server Error"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetLaborMarket(LaborMarketParams{})
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

// TestGetTreasuryYields verifies that GetTreasuryYields correctly parses
// the API response and returns the expected yield curve data across
// multiple maturities.
func TestGetTreasuryYields(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/fed/v1/treasury-yields": treasuryYieldsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetTreasuryYields(TreasuryYieldParams{
		Sort:  "date.desc",
		Limit: "2",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "c33732ef5b614ef7990f3899dbbc6da2" {
		t.Errorf("expected request_id c33732ef5b614ef7990f3899dbbc6da2, got %s", result.RequestID)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.Date != "2026-02-12" {
		t.Errorf("expected date 2026-02-12, got %s", first.Date)
	}

	if first.Yield1Month != 3.72 {
		t.Errorf("expected yield_1_month 3.72, got %f", first.Yield1Month)
	}

	if first.Yield3Month != 3.70 {
		t.Errorf("expected yield_3_month 3.70, got %f", first.Yield3Month)
	}

	if first.Yield1Year != 3.45 {
		t.Errorf("expected yield_1_year 3.45, got %f", first.Yield1Year)
	}

	if first.Yield2Year != 3.47 {
		t.Errorf("expected yield_2_year 3.47, got %f", first.Yield2Year)
	}

	if first.Yield5Year != 3.67 {
		t.Errorf("expected yield_5_year 3.67, got %f", first.Yield5Year)
	}

	if first.Yield10Year != 4.09 {
		t.Errorf("expected yield_10_year 4.09, got %f", first.Yield10Year)
	}

	if first.Yield30Year != 4.72 {
		t.Errorf("expected yield_30_year 4.72, got %f", first.Yield30Year)
	}
}

// TestGetTreasuryYieldsRequestPath verifies that GetTreasuryYields
// constructs the correct API path for the treasury yields endpoint.
func TestGetTreasuryYieldsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(treasuryYieldsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetTreasuryYields(TreasuryYieldParams{})

	if receivedPath != "/fed/v1/treasury-yields" {
		t.Errorf("expected path /fed/v1/treasury-yields, got %s", receivedPath)
	}
}

// TestGetTreasuryYieldsQueryParams verifies that all filter parameters
// including date range operators, sort, and limit are correctly sent.
func TestGetTreasuryYieldsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("date.gte") != "2026-01-01" {
			t.Errorf("expected date.gte=2026-01-01, got %s", q.Get("date.gte"))
		}
		if q.Get("date.lte") != "2026-02-15" {
			t.Errorf("expected date.lte=2026-02-15, got %s", q.Get("date.lte"))
		}
		if q.Get("sort") != "date.asc" {
			t.Errorf("expected sort=date.asc, got %s", q.Get("sort"))
		}
		if q.Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", q.Get("limit"))
		}
		if q.Get("apiKey") != "test-api-key" {
			t.Errorf("expected apiKey=test-api-key, got %s", q.Get("apiKey"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(treasuryYieldsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetTreasuryYields(TreasuryYieldParams{
		DateGTE: "2026-01-01",
		DateLTE: "2026-02-15",
		Sort:    "date.asc",
		Limit:   "50",
	})
}

// TestGetTreasuryYieldsSecondResult verifies that the second result in the
// treasury yields response is correctly parsed with its own distinct values.
func TestGetTreasuryYieldsSecondResult(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/fed/v1/treasury-yields": treasuryYieldsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetTreasuryYields(TreasuryYieldParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := result.Results[1]
	if second.Date != "2026-02-11" {
		t.Errorf("expected date 2026-02-11, got %s", second.Date)
	}

	if second.Yield1Month != 3.71 {
		t.Errorf("expected yield_1_month 3.71, got %f", second.Yield1Month)
	}

	if second.Yield2Year != 3.52 {
		t.Errorf("expected yield_2_year 3.52, got %f", second.Yield2Year)
	}

	if second.Yield10Year != 4.18 {
		t.Errorf("expected yield_10_year 4.18, got %f", second.Yield10Year)
	}

	if second.Yield30Year != 4.82 {
		t.Errorf("expected yield_30_year 4.82, got %f", second.Yield30Year)
	}
}

// TestGetTreasuryYieldsAPIError verifies that GetTreasuryYields returns
// an error when the API responds with a non-200 status code.
func TestGetTreasuryYieldsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"ERROR","message":"Forbidden"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetTreasuryYields(TreasuryYieldParams{})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetTreasuryYieldsEmptyResults verifies that GetTreasuryYields handles
// an empty results array without error.
func TestGetTreasuryYieldsEmptyResults(t *testing.T) {
	emptyJSON := `{"status":"OK","request_id":"abc123","results":[]}`
	server := mockServer(t, map[string]string{
		"/fed/v1/treasury-yields": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetTreasuryYields(TreasuryYieldParams{Date: "2020-01-01"})
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

// TestGetInflationEmptyResults verifies that GetInflation handles an
// empty results array without error.
func TestGetInflationEmptyResults(t *testing.T) {
	emptyJSON := `{"status":"OK","request_id":"abc123","results":[]}`
	server := mockServer(t, map[string]string{
		"/fed/v1/inflation": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetInflation(InflationParams{Date: "2020-01-01"})
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

// TestGetLaborMarketEmptyResults verifies that GetLaborMarket handles an
// empty results array without error.
func TestGetLaborMarketEmptyResults(t *testing.T) {
	emptyJSON := `{"status":"OK","request_id":"abc123","results":[]}`
	server := mockServer(t, map[string]string{
		"/fed/v1/labor-market": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetLaborMarket(LaborMarketParams{Date: "2020-01-01"})
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
