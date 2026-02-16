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

const secFilingSectionsJSON = `{
	"status": "OK",
	"request_id": "72d89ba2c12644e586279dbad8d53922",
	"results": [
		{
			"cik": "0000320193",
			"ticker": "AAPL",
			"section": "risk_factors",
			"filing_date": "2025-10-31",
			"period_end": "2025-09-27",
			"text": "Item 1A. Risk Factors\nThe following summarizes factors that could have a material adverse effect.",
			"filing_url": "https://www.sec.gov/Archives/edgar/data/320193/0000320193-25-000079.txt"
		},
		{
			"cik": "0000320193",
			"ticker": "AAPL",
			"section": "business",
			"filing_date": "2025-10-31",
			"period_end": "2025-09-27",
			"text": "Item 1. Business\nThe Company designs, manufactures and markets smartphones and other products.",
			"filing_url": "https://www.sec.gov/Archives/edgar/data/320193/0000320193-25-000079.txt"
		}
	],
	"next_url": "https://api.massive.com/stocks/filings/10-K/vX/sections?cursor=abc123"
}`

const riskFactorsJSON = `{
	"status": "OK",
	"request_id": "803c79037bf0402abb120314e7c3d9ea",
	"results": [
		{
			"cik": "0000320193",
			"ticker": "AAPL",
			"primary_category": "financial_and_market",
			"secondary_category": "capital_structure_and_performance",
			"tertiary_category": "dividend_policy_and_capital_allocation",
			"filing_date": "2024-11-01",
			"supporting_text": "The Company believes the price of its stock should reflect expectations that its cash dividend will continue at current levels or grow."
		},
		{
			"cik": "0000320193",
			"ticker": "AAPL",
			"primary_category": "operational_and_infrastructure",
			"secondary_category": "supply_chain_and_logistics",
			"tertiary_category": "supplier_concentration_and_dependence",
			"filing_date": "2024-11-01",
			"supporting_text": "The Company depends on component and product manufacturing and logistical services provided by outsourcing partners."
		}
	],
	"next_url": "https://api.massive.com/stocks/filings/vX/risk-factors?cursor=def456"
}`

const riskCategoriesJSON = `{
	"status": "OK",
	"request_id": "73eb31445acb49088c0c9cd64ca80a31",
	"results": [
		{
			"primary_category": "governance_and_stakeholder",
			"secondary_category": "organizational_and_management",
			"tertiary_category": "performance_management_and_accountability",
			"description": "Risk from inadequate performance management systems, unclear accountability structures.",
			"taxonomy": 1.0
		},
		{
			"primary_category": "governance_and_stakeholder",
			"secondary_category": "organizational_and_management",
			"tertiary_category": "communication_and_coordination",
			"description": "Risk from poor internal communication, lack of coordination between departments.",
			"taxonomy": 1.0
		}
	],
	"next_url": "https://api.massive.com/stocks/taxonomies/vX/risk-factors?cursor=ghi789"
}`

// TestGetSECFilingSections verifies that GetSECFilingSections correctly
// parses the API response and returns the expected section data for AAPL
// 10-K filings.
func TestGetSECFilingSections(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/stocks/filings/10-K/vX/sections": secFilingSectionsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := SECFilingSectionsParams{
		Ticker: "AAPL",
		Limit:  "2",
	}

	result, err := client.GetSECFilingSections(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "72d89ba2c12644e586279dbad8d53922" {
		t.Errorf("expected request_id 72d89ba2c12644e586279dbad8d53922, got %s", result.RequestID)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.CIK != "0000320193" {
		t.Errorf("expected CIK 0000320193, got %s", first.CIK)
	}

	if first.Ticker != "AAPL" {
		t.Errorf("expected ticker AAPL, got %s", first.Ticker)
	}

	if first.Section != "risk_factors" {
		t.Errorf("expected section risk_factors, got %s", first.Section)
	}

	if first.FilingDate != "2025-10-31" {
		t.Errorf("expected filing_date 2025-10-31, got %s", first.FilingDate)
	}

	if first.PeriodEnd != "2025-09-27" {
		t.Errorf("expected period_end 2025-09-27, got %s", first.PeriodEnd)
	}

	if first.Text == "" {
		t.Error("expected text to be populated")
	}

	if first.FilingURL != "https://www.sec.gov/Archives/edgar/data/320193/0000320193-25-000079.txt" {
		t.Errorf("expected filing_url to match SEC URL, got %s", first.FilingURL)
	}

	second := result.Results[1]
	if second.Section != "business" {
		t.Errorf("expected section business, got %s", second.Section)
	}
}

// TestGetSECFilingSectionsRequestPath verifies that GetSECFilingSections
// constructs the correct API path for the 10-K sections endpoint.
func TestGetSECFilingSectionsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(secFilingSectionsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetSECFilingSections(SECFilingSectionsParams{Ticker: "AAPL"})

	expected := "/stocks/filings/10-K/vX/sections"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetSECFilingSectionsQueryParams verifies that all filter parameters
// are correctly sent to the API endpoint as query string values.
func TestGetSECFilingSectionsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("ticker") != "MSFT" {
			t.Errorf("expected ticker=MSFT, got %s", q.Get("ticker"))
		}
		if q.Get("section") != "business" {
			t.Errorf("expected section=business, got %s", q.Get("section"))
		}
		if q.Get("filing_date.gt") != "2024-01-01" {
			t.Errorf("expected filing_date.gt=2024-01-01, got %s", q.Get("filing_date.gt"))
		}
		if q.Get("filing_date.lt") != "2025-12-31" {
			t.Errorf("expected filing_date.lt=2025-12-31, got %s", q.Get("filing_date.lt"))
		}
		if q.Get("period_end.gt") != "2024-06-01" {
			t.Errorf("expected period_end.gt=2024-06-01, got %s", q.Get("period_end.gt"))
		}
		if q.Get("limit") != "5" {
			t.Errorf("expected limit=5, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "filing_date.asc" {
			t.Errorf("expected sort=filing_date.asc, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(secFilingSectionsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetSECFilingSections(SECFilingSectionsParams{
		Ticker:       "MSFT",
		Section:      "business",
		FilingDateGt: "2024-01-01",
		FilingDateLt: "2025-12-31",
		PeriodEndGt:  "2024-06-01",
		Limit:        "5",
		Sort:         "filing_date.asc",
	})
}

// TestGetSECFilingSectionsAPIError verifies that GetSECFilingSections
// returns an error when the API responds with a non-200 status.
func TestGetSECFilingSectionsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"ERROR","message":"Insufficient permissions"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetSECFilingSections(SECFilingSectionsParams{Ticker: "AAPL"})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetSECFilingSectionsWithCIK verifies that filtering by CIK instead
// of ticker correctly sends the cik query parameter.
func TestGetSECFilingSectionsWithCIK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("cik") != "0000320193" {
			t.Errorf("expected cik=0000320193, got %s", r.URL.Query().Get("cik"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(secFilingSectionsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetSECFilingSections(SECFilingSectionsParams{CIK: "0000320193"})
}

// TestGetRiskFactors verifies that GetRiskFactors correctly parses the
// API response and returns the expected categorized risk factor data.
func TestGetRiskFactors(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/stocks/filings/vX/risk-factors": riskFactorsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := RiskFactorsParams{
		Ticker: "AAPL",
		Limit:  "2",
	}

	result, err := client.GetRiskFactors(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.CIK != "0000320193" {
		t.Errorf("expected CIK 0000320193, got %s", first.CIK)
	}

	if first.Ticker != "AAPL" {
		t.Errorf("expected ticker AAPL, got %s", first.Ticker)
	}

	if first.PrimaryCategory != "financial_and_market" {
		t.Errorf("expected primary_category financial_and_market, got %s", first.PrimaryCategory)
	}

	if first.SecondaryCategory != "capital_structure_and_performance" {
		t.Errorf("expected secondary_category capital_structure_and_performance, got %s", first.SecondaryCategory)
	}

	if first.TertiaryCategory != "dividend_policy_and_capital_allocation" {
		t.Errorf("expected tertiary_category dividend_policy_and_capital_allocation, got %s", first.TertiaryCategory)
	}

	if first.FilingDate != "2024-11-01" {
		t.Errorf("expected filing_date 2024-11-01, got %s", first.FilingDate)
	}

	if first.SupportingText == "" {
		t.Error("expected supporting_text to be populated")
	}

	second := result.Results[1]
	if second.PrimaryCategory != "operational_and_infrastructure" {
		t.Errorf("expected primary_category operational_and_infrastructure, got %s", second.PrimaryCategory)
	}

	if second.SecondaryCategory != "supply_chain_and_logistics" {
		t.Errorf("expected secondary_category supply_chain_and_logistics, got %s", second.SecondaryCategory)
	}
}

// TestGetRiskFactorsRequestPath verifies that GetRiskFactors constructs
// the correct API path for the risk factors endpoint.
func TestGetRiskFactorsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(riskFactorsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetRiskFactors(RiskFactorsParams{Ticker: "AAPL"})

	expected := "/stocks/filings/vX/risk-factors"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetRiskFactorsQueryParams verifies that all filter parameters are
// correctly sent to the risk factors API endpoint.
func TestGetRiskFactorsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("ticker") != "TSLA" {
			t.Errorf("expected ticker=TSLA, got %s", q.Get("ticker"))
		}
		if q.Get("cik") != "0001318605" {
			t.Errorf("expected cik=0001318605, got %s", q.Get("cik"))
		}
		if q.Get("filing_date.gt") != "2023-01-01" {
			t.Errorf("expected filing_date.gt=2023-01-01, got %s", q.Get("filing_date.gt"))
		}
		if q.Get("filing_date.lt") != "2025-01-01" {
			t.Errorf("expected filing_date.lt=2025-01-01, got %s", q.Get("filing_date.lt"))
		}
		if q.Get("limit") != "10" {
			t.Errorf("expected limit=10, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "filing_date.asc" {
			t.Errorf("expected sort=filing_date.asc, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(riskFactorsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetRiskFactors(RiskFactorsParams{
		Ticker:       "TSLA",
		CIK:          "0001318605",
		FilingDateGt: "2023-01-01",
		FilingDateLt: "2025-01-01",
		Limit:        "10",
		Sort:         "filing_date.asc",
	})
}

// TestGetRiskFactorsAPIError verifies that GetRiskFactors returns an
// error when the API responds with a non-200 status code.
func TestGetRiskFactorsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"ERROR","message":"Internal server error"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetRiskFactors(RiskFactorsParams{Ticker: "AAPL"})
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

// TestGetRiskCategories verifies that GetRiskCategories correctly parses
// the API response and returns the expected taxonomy entries.
func TestGetRiskCategories(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/stocks/taxonomies/vX/risk-factors": riskCategoriesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := RiskCategoriesParams{
		Limit: "2",
	}

	result, err := client.GetRiskCategories(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.PrimaryCategory != "governance_and_stakeholder" {
		t.Errorf("expected primary_category governance_and_stakeholder, got %s", first.PrimaryCategory)
	}

	if first.SecondaryCategory != "organizational_and_management" {
		t.Errorf("expected secondary_category organizational_and_management, got %s", first.SecondaryCategory)
	}

	if first.TertiaryCategory != "performance_management_and_accountability" {
		t.Errorf("expected tertiary_category performance_management_and_accountability, got %s", first.TertiaryCategory)
	}

	if first.Description == "" {
		t.Error("expected description to be populated")
	}

	if first.Taxonomy != 1.0 {
		t.Errorf("expected taxonomy 1.0, got %f", first.Taxonomy)
	}

	second := result.Results[1]
	if second.TertiaryCategory != "communication_and_coordination" {
		t.Errorf("expected tertiary_category communication_and_coordination, got %s", second.TertiaryCategory)
	}
}

// TestGetRiskCategoriesRequestPath verifies that GetRiskCategories
// constructs the correct API path for the taxonomies endpoint.
func TestGetRiskCategoriesRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(riskCategoriesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetRiskCategories(RiskCategoriesParams{})

	expected := "/stocks/taxonomies/vX/risk-factors"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetRiskCategoriesQueryParams verifies that all filter parameters
// are correctly sent to the risk categories API endpoint.
func TestGetRiskCategoriesQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("primary_category") != "financial_and_market" {
			t.Errorf("expected primary_category=financial_and_market, got %s", q.Get("primary_category"))
		}
		if q.Get("secondary_category") != "capital_structure_and_performance" {
			t.Errorf("expected secondary_category=capital_structure_and_performance, got %s", q.Get("secondary_category"))
		}
		if q.Get("tertiary_category") != "dividend_policy_and_capital_allocation" {
			t.Errorf("expected tertiary_category=dividend_policy_and_capital_allocation, got %s", q.Get("tertiary_category"))
		}
		if q.Get("taxonomy") != "1.0" {
			t.Errorf("expected taxonomy=1.0, got %s", q.Get("taxonomy"))
		}
		if q.Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "primary_category.asc" {
			t.Errorf("expected sort=primary_category.asc, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(riskCategoriesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetRiskCategories(RiskCategoriesParams{
		PrimaryCategory:   "financial_and_market",
		SecondaryCategory: "capital_structure_and_performance",
		TertiaryCategory:  "dividend_policy_and_capital_allocation",
		Taxonomy:          "1.0",
		Limit:             "50",
		Sort:              "primary_category.asc",
	})
}

// TestGetRiskCategoriesAPIError verifies that GetRiskCategories returns
// an error when the API responds with a non-200 status code.
func TestGetRiskCategoriesAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Resource not found"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetRiskCategories(RiskCategoriesParams{})
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// TestGetRiskCategoriesEmptyResults verifies that GetRiskCategories
// handles an empty results array without error.
func TestGetRiskCategoriesEmptyResults(t *testing.T) {
	emptyJSON := `{"status":"OK","request_id":"abc123","results":[]}`
	server := mockServer(t, map[string]string{
		"/stocks/taxonomies/vX/risk-factors": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetRiskCategories(RiskCategoriesParams{
		PrimaryCategory: "nonexistent_category",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(result.Results))
	}
}
