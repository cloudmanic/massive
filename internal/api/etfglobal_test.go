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

const analyticsJSON = `{
	"status": "OK",
	"request_id": "abc123analytics",
	"count": 2,
	"next_url": "https://api.massive.com/etf-global/v1/analytics?cursor=YXA9Mg",
	"results": [
		{
			"composite_ticker": "SPY",
			"effective_date": "2025-12-15",
			"processed_date": "2025-12-16",
			"quant_composite_behavioral": 62.5,
			"quant_composite_fundamental": 55.3,
			"quant_composite_global": 70.1,
			"quant_composite_quality": 80.2,
			"quant_composite_sentiment": 45.8,
			"quant_composite_technical": 68.9,
			"quant_fundamental_div": 50.0,
			"quant_fundamental_pb": 48.5,
			"quant_fundamental_pcf": 60.2,
			"quant_fundamental_pe": 55.1,
			"quant_global_country": 72.0,
			"quant_global_sector": 68.3,
			"quant_grade": "B",
			"quant_quality_diversification": 85.0,
			"quant_quality_firm": 90.5,
			"quant_quality_liquidity": 95.0,
			"quant_sentiment_iv": 40.2,
			"quant_sentiment_pc": 50.5,
			"quant_sentiment_si": 46.7,
			"quant_technical_it": 65.0,
			"quant_technical_lt": 72.3,
			"quant_technical_st": 69.5,
			"quant_total_score": 67.8,
			"reward_score": 72.5,
			"risk_country": 10.0,
			"risk_deviation": 25.3,
			"risk_efficiency": 88.5,
			"risk_liquidity": 95.0,
			"risk_structure": 90.0,
			"risk_total_score": 30.2,
			"risk_volatility": 22.1
		},
		{
			"composite_ticker": "QQQ",
			"effective_date": "2025-12-15",
			"processed_date": "2025-12-16",
			"quant_composite_behavioral": 58.2,
			"quant_composite_fundamental": 50.1,
			"quant_composite_global": 65.4,
			"quant_composite_quality": 78.9,
			"quant_composite_sentiment": 42.3,
			"quant_composite_technical": 71.6,
			"quant_fundamental_div": 35.0,
			"quant_fundamental_pb": 42.0,
			"quant_fundamental_pcf": 55.8,
			"quant_fundamental_pe": 48.3,
			"quant_global_country": 68.0,
			"quant_global_sector": 62.7,
			"quant_grade": "B",
			"quant_quality_diversification": 70.5,
			"quant_quality_firm": 88.0,
			"quant_quality_liquidity": 92.3,
			"quant_sentiment_iv": 38.5,
			"quant_sentiment_pc": 45.0,
			"quant_sentiment_si": 43.1,
			"quant_technical_it": 70.0,
			"quant_technical_lt": 68.5,
			"quant_technical_st": 76.2,
			"quant_total_score": 63.5,
			"reward_score": 68.0,
			"risk_country": 8.5,
			"risk_deviation": 30.1,
			"risk_efficiency": 85.0,
			"risk_liquidity": 93.5,
			"risk_structure": 88.0,
			"risk_total_score": 35.5,
			"risk_volatility": 28.7
		}
	]
}`

const constituentsJSON = `{
	"status": "OK",
	"request_id": "def456constituents",
	"count": 2,
	"next_url": "https://api.massive.com/etf-global/v1/constituents?cursor=YXA9Mg",
	"results": [
		{
			"asset_class": "Equity",
			"composite_ticker": "SPY",
			"constituent_name": "Apple Inc.",
			"constituent_rank": 1,
			"constituent_ticker": "AAPL",
			"country_of_exchange": "United States",
			"currency_traded": "USD",
			"effective_date": "2025-12-15",
			"exchange": "NASDAQ",
			"figi": "BBG000B9XRY4",
			"isin": "US0378331005",
			"market_value": 5250000000.50,
			"processed_date": "2025-12-16",
			"security_type": "Common Stock",
			"sedol": "2046251",
			"shares_held": 45000000.0,
			"us_code": "0378331005",
			"weight": 7.25
		},
		{
			"asset_class": "Equity",
			"composite_ticker": "SPY",
			"constituent_name": "Microsoft Corporation",
			"constituent_rank": 2,
			"constituent_ticker": "MSFT",
			"country_of_exchange": "United States",
			"currency_traded": "USD",
			"effective_date": "2025-12-15",
			"exchange": "NASDAQ",
			"figi": "BBG000BPH459",
			"isin": "US5949181045",
			"market_value": 4800000000.25,
			"processed_date": "2025-12-16",
			"security_type": "Common Stock",
			"sedol": "2588173",
			"shares_held": 38000000.0,
			"us_code": "5949181045",
			"weight": 6.85
		}
	]
}`

// TestGetETFGlobalAnalytics verifies that GetETFGlobalAnalytics correctly
// parses the API response and returns the expected analytics data for
// multiple ETFs including all scoring and risk fields.
func TestGetETFGlobalAnalytics(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/etf-global/v1/analytics": analyticsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := ETFGlobalAnalyticsParams{
		CompositeTicker: "SPY",
		Limit:           "2",
	}

	result, err := client.GetETFGlobalAnalytics(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "abc123analytics" {
		t.Errorf("expected request_id abc123analytics, got %s", result.RequestID)
	}

	if result.Count != 2 {
		t.Errorf("expected count 2, got %d", result.Count)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 analytics results, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.CompositeTicker != "SPY" {
		t.Errorf("expected composite_ticker SPY, got %s", first.CompositeTicker)
	}

	if first.EffectiveDate != "2025-12-15" {
		t.Errorf("expected effective_date 2025-12-15, got %s", first.EffectiveDate)
	}

	if first.ProcessedDate != "2025-12-16" {
		t.Errorf("expected processed_date 2025-12-16, got %s", first.ProcessedDate)
	}

	if first.QuantGrade != "B" {
		t.Errorf("expected quant_grade B, got %s", first.QuantGrade)
	}

	if first.QuantTotalScore != 67.8 {
		t.Errorf("expected quant_total_score 67.8, got %f", first.QuantTotalScore)
	}

	if first.RewardScore != 72.5 {
		t.Errorf("expected reward_score 72.5, got %f", first.RewardScore)
	}

	if first.RiskTotalScore != 30.2 {
		t.Errorf("expected risk_total_score 30.2, got %f", first.RiskTotalScore)
	}

	if first.RiskVolatility != 22.1 {
		t.Errorf("expected risk_volatility 22.1, got %f", first.RiskVolatility)
	}

	if first.RiskLiquidity != 95.0 {
		t.Errorf("expected risk_liquidity 95.0, got %f", first.RiskLiquidity)
	}

	if first.RiskEfficiency != 88.5 {
		t.Errorf("expected risk_efficiency 88.5, got %f", first.RiskEfficiency)
	}

	if first.QuantCompositeTechnical != 68.9 {
		t.Errorf("expected quant_composite_technical 68.9, got %f", first.QuantCompositeTechnical)
	}

	if first.QuantCompositeSentiment != 45.8 {
		t.Errorf("expected quant_composite_sentiment 45.8, got %f", first.QuantCompositeSentiment)
	}

	if first.QuantCompositeBehavioral != 62.5 {
		t.Errorf("expected quant_composite_behavioral 62.5, got %f", first.QuantCompositeBehavioral)
	}

	if first.QuantCompositeFundamental != 55.3 {
		t.Errorf("expected quant_composite_fundamental 55.3, got %f", first.QuantCompositeFundamental)
	}

	if first.QuantCompositeQuality != 80.2 {
		t.Errorf("expected quant_composite_quality 80.2, got %f", first.QuantCompositeQuality)
	}

	if first.QuantCompositeGlobal != 70.1 {
		t.Errorf("expected quant_composite_global 70.1, got %f", first.QuantCompositeGlobal)
	}
}

// TestGetETFGlobalAnalyticsSecondResult verifies that the second analytics
// record in the response is correctly parsed with its own distinct values.
func TestGetETFGlobalAnalyticsSecondResult(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/etf-global/v1/analytics": analyticsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetETFGlobalAnalytics(ETFGlobalAnalyticsParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := result.Results[1]
	if second.CompositeTicker != "QQQ" {
		t.Errorf("expected composite_ticker QQQ, got %s", second.CompositeTicker)
	}

	if second.QuantTotalScore != 63.5 {
		t.Errorf("expected quant_total_score 63.5, got %f", second.QuantTotalScore)
	}

	if second.RewardScore != 68.0 {
		t.Errorf("expected reward_score 68.0, got %f", second.RewardScore)
	}

	if second.RiskTotalScore != 35.5 {
		t.Errorf("expected risk_total_score 35.5, got %f", second.RiskTotalScore)
	}

	if second.RiskVolatility != 28.7 {
		t.Errorf("expected risk_volatility 28.7, got %f", second.RiskVolatility)
	}

	if second.QuantCompositeTechnical != 71.6 {
		t.Errorf("expected quant_composite_technical 71.6, got %f", second.QuantCompositeTechnical)
	}
}

// TestGetETFGlobalAnalyticsRequestPath verifies that GetETFGlobalAnalytics
// constructs the correct API path for the analytics endpoint.
func TestGetETFGlobalAnalyticsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(analyticsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetETFGlobalAnalytics(ETFGlobalAnalyticsParams{CompositeTicker: "SPY"})

	expected := "/etf-global/v1/analytics"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetETFGlobalAnalyticsQueryParams verifies that all filter parameters
// are correctly sent to the analytics API endpoint as query parameters.
func TestGetETFGlobalAnalyticsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("composite_ticker") != "SPY" {
			t.Errorf("expected composite_ticker=SPY, got %s", q.Get("composite_ticker"))
		}
		if q.Get("processed_date") != "2025-12-16" {
			t.Errorf("expected processed_date=2025-12-16, got %s", q.Get("processed_date"))
		}
		if q.Get("effective_date") != "2025-12-15" {
			t.Errorf("expected effective_date=2025-12-15, got %s", q.Get("effective_date"))
		}
		if q.Get("quant_grade") != "A" {
			t.Errorf("expected quant_grade=A, got %s", q.Get("quant_grade"))
		}
		if q.Get("sort") != "quant_total_score.desc" {
			t.Errorf("expected sort=quant_total_score.desc, got %s", q.Get("sort"))
		}
		if q.Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(analyticsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetETFGlobalAnalytics(ETFGlobalAnalyticsParams{
		CompositeTicker: "SPY",
		ProcessedDate:   "2025-12-16",
		EffectiveDate:   "2025-12-15",
		QuantGrade:      "A",
		Sort:            "quant_total_score.desc",
		Limit:           "50",
	})
}

// TestGetETFGlobalAnalyticsAPIError verifies that GetETFGlobalAnalytics
// returns an error when the API responds with a non-200 status code.
func TestGetETFGlobalAnalyticsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"NOT_AUTHORIZED","message":"You are not entitled to this data."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetETFGlobalAnalytics(ETFGlobalAnalyticsParams{CompositeTicker: "SPY"})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetETFGlobalAnalyticsEmptyResults verifies that GetETFGlobalAnalytics
// handles an empty results array without error.
func TestGetETFGlobalAnalyticsEmptyResults(t *testing.T) {
	emptyJSON := `{"status":"OK","request_id":"abc","count":0,"results":[]}`
	server := mockServer(t, map[string]string{
		"/etf-global/v1/analytics": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetETFGlobalAnalytics(ETFGlobalAnalyticsParams{CompositeTicker: "ZZZNOTREAL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Count != 0 {
		t.Errorf("expected count 0, got %d", result.Count)
	}

	if len(result.Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(result.Results))
	}
}

// TestGetETFGlobalAnalyticsFundamentalFields verifies that all individual
// fundamental sub-scores are correctly parsed from the analytics response.
func TestGetETFGlobalAnalyticsFundamentalFields(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/etf-global/v1/analytics": analyticsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetETFGlobalAnalytics(ETFGlobalAnalyticsParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	first := result.Results[0]

	if first.QuantFundamentalDiv != 50.0 {
		t.Errorf("expected quant_fundamental_div 50.0, got %f", first.QuantFundamentalDiv)
	}

	if first.QuantFundamentalPB != 48.5 {
		t.Errorf("expected quant_fundamental_pb 48.5, got %f", first.QuantFundamentalPB)
	}

	if first.QuantFundamentalPCF != 60.2 {
		t.Errorf("expected quant_fundamental_pcf 60.2, got %f", first.QuantFundamentalPCF)
	}

	if first.QuantFundamentalPE != 55.1 {
		t.Errorf("expected quant_fundamental_pe 55.1, got %f", first.QuantFundamentalPE)
	}

	if first.QuantGlobalCountry != 72.0 {
		t.Errorf("expected quant_global_country 72.0, got %f", first.QuantGlobalCountry)
	}

	if first.QuantGlobalSector != 68.3 {
		t.Errorf("expected quant_global_sector 68.3, got %f", first.QuantGlobalSector)
	}

	if first.QuantQualityDiversify != 85.0 {
		t.Errorf("expected quant_quality_diversification 85.0, got %f", first.QuantQualityDiversify)
	}

	if first.QuantQualityFirm != 90.5 {
		t.Errorf("expected quant_quality_firm 90.5, got %f", first.QuantQualityFirm)
	}

	if first.QuantQualityLiquidity != 95.0 {
		t.Errorf("expected quant_quality_liquidity 95.0, got %f", first.QuantQualityLiquidity)
	}

	if first.QuantSentimentIV != 40.2 {
		t.Errorf("expected quant_sentiment_iv 40.2, got %f", first.QuantSentimentIV)
	}

	if first.QuantSentimentPC != 50.5 {
		t.Errorf("expected quant_sentiment_pc 50.5, got %f", first.QuantSentimentPC)
	}

	if first.QuantSentimentSI != 46.7 {
		t.Errorf("expected quant_sentiment_si 46.7, got %f", first.QuantSentimentSI)
	}

	if first.QuantTechnicalIT != 65.0 {
		t.Errorf("expected quant_technical_it 65.0, got %f", first.QuantTechnicalIT)
	}

	if first.QuantTechnicalLT != 72.3 {
		t.Errorf("expected quant_technical_lt 72.3, got %f", first.QuantTechnicalLT)
	}

	if first.QuantTechnicalST != 69.5 {
		t.Errorf("expected quant_technical_st 69.5, got %f", first.QuantTechnicalST)
	}
}

// TestGetETFGlobalAnalyticsRiskFields verifies that all individual risk
// sub-scores are correctly parsed from the analytics response.
func TestGetETFGlobalAnalyticsRiskFields(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/etf-global/v1/analytics": analyticsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetETFGlobalAnalytics(ETFGlobalAnalyticsParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	first := result.Results[0]

	if first.RiskCountry != 10.0 {
		t.Errorf("expected risk_country 10.0, got %f", first.RiskCountry)
	}

	if first.RiskDeviation != 25.3 {
		t.Errorf("expected risk_deviation 25.3, got %f", first.RiskDeviation)
	}

	if first.RiskEfficiency != 88.5 {
		t.Errorf("expected risk_efficiency 88.5, got %f", first.RiskEfficiency)
	}

	if first.RiskStructure != 90.0 {
		t.Errorf("expected risk_structure 90.0, got %f", first.RiskStructure)
	}
}

// TestGetETFGlobalConstituents verifies that GetETFGlobalConstituents correctly
// parses the API response and returns the expected constituent holdings data.
func TestGetETFGlobalConstituents(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/etf-global/v1/constituents": constituentsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := ETFGlobalConstituentsParams{
		CompositeTicker: "SPY",
		Limit:           "2",
	}

	result, err := client.GetETFGlobalConstituents(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "def456constituents" {
		t.Errorf("expected request_id def456constituents, got %s", result.RequestID)
	}

	if result.Count != 2 {
		t.Errorf("expected count 2, got %d", result.Count)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 constituent results, got %d", len(result.Results))
	}

	first := result.Results[0]
	if first.CompositeTicker != "SPY" {
		t.Errorf("expected composite_ticker SPY, got %s", first.CompositeTicker)
	}

	if first.ConstituentTicker != "AAPL" {
		t.Errorf("expected constituent_ticker AAPL, got %s", first.ConstituentTicker)
	}

	if first.ConstituentName != "Apple Inc." {
		t.Errorf("expected constituent_name Apple Inc., got %s", first.ConstituentName)
	}

	if first.ConstituentRank != 1 {
		t.Errorf("expected constituent_rank 1, got %d", first.ConstituentRank)
	}

	if first.AssetClass != "Equity" {
		t.Errorf("expected asset_class Equity, got %s", first.AssetClass)
	}

	if first.SecurityType != "Common Stock" {
		t.Errorf("expected security_type Common Stock, got %s", first.SecurityType)
	}

	if first.Weight != 7.25 {
		t.Errorf("expected weight 7.25, got %f", first.Weight)
	}

	if first.SharesHeld != 45000000.0 {
		t.Errorf("expected shares_held 45000000.0, got %f", first.SharesHeld)
	}

	if first.MarketValue != 5250000000.50 {
		t.Errorf("expected market_value 5250000000.50, got %f", first.MarketValue)
	}

	if first.Exchange != "NASDAQ" {
		t.Errorf("expected exchange NASDAQ, got %s", first.Exchange)
	}

	if first.CountryOfExchange != "United States" {
		t.Errorf("expected country_of_exchange United States, got %s", first.CountryOfExchange)
	}

	if first.CurrencyTraded != "USD" {
		t.Errorf("expected currency_traded USD, got %s", first.CurrencyTraded)
	}

	if first.ISIN != "US0378331005" {
		t.Errorf("expected isin US0378331005, got %s", first.ISIN)
	}

	if first.FIGI != "BBG000B9XRY4" {
		t.Errorf("expected figi BBG000B9XRY4, got %s", first.FIGI)
	}

	if first.SEDOL != "2046251" {
		t.Errorf("expected sedol 2046251, got %s", first.SEDOL)
	}

	if first.USCode != "0378331005" {
		t.Errorf("expected us_code 0378331005, got %s", first.USCode)
	}

	if first.EffectiveDate != "2025-12-15" {
		t.Errorf("expected effective_date 2025-12-15, got %s", first.EffectiveDate)
	}

	if first.ProcessedDate != "2025-12-16" {
		t.Errorf("expected processed_date 2025-12-16, got %s", first.ProcessedDate)
	}
}

// TestGetETFGlobalConstituentsSecondResult verifies that the second
// constituent in the response is correctly parsed with its own values.
func TestGetETFGlobalConstituentsSecondResult(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/etf-global/v1/constituents": constituentsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetETFGlobalConstituents(ETFGlobalConstituentsParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := result.Results[1]
	if second.ConstituentTicker != "MSFT" {
		t.Errorf("expected constituent_ticker MSFT, got %s", second.ConstituentTicker)
	}

	if second.ConstituentName != "Microsoft Corporation" {
		t.Errorf("expected constituent_name Microsoft Corporation, got %s", second.ConstituentName)
	}

	if second.ConstituentRank != 2 {
		t.Errorf("expected constituent_rank 2, got %d", second.ConstituentRank)
	}

	if second.Weight != 6.85 {
		t.Errorf("expected weight 6.85, got %f", second.Weight)
	}

	if second.MarketValue != 4800000000.25 {
		t.Errorf("expected market_value 4800000000.25, got %f", second.MarketValue)
	}

	if second.SharesHeld != 38000000.0 {
		t.Errorf("expected shares_held 38000000.0, got %f", second.SharesHeld)
	}

	if second.ISIN != "US5949181045" {
		t.Errorf("expected isin US5949181045, got %s", second.ISIN)
	}
}

// TestGetETFGlobalConstituentsRequestPath verifies that GetETFGlobalConstituents
// constructs the correct API path for the constituents endpoint.
func TestGetETFGlobalConstituentsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(constituentsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetETFGlobalConstituents(ETFGlobalConstituentsParams{CompositeTicker: "SPY"})

	expected := "/etf-global/v1/constituents"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetETFGlobalConstituentsQueryParams verifies that all filter parameters
// are correctly sent to the constituents API endpoint as query parameters.
func TestGetETFGlobalConstituentsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("composite_ticker") != "SPY" {
			t.Errorf("expected composite_ticker=SPY, got %s", q.Get("composite_ticker"))
		}
		if q.Get("constituent_ticker") != "AAPL" {
			t.Errorf("expected constituent_ticker=AAPL, got %s", q.Get("constituent_ticker"))
		}
		if q.Get("effective_date") != "2025-12-15" {
			t.Errorf("expected effective_date=2025-12-15, got %s", q.Get("effective_date"))
		}
		if q.Get("processed_date") != "2025-12-16" {
			t.Errorf("expected processed_date=2025-12-16, got %s", q.Get("processed_date"))
		}
		if q.Get("isin") != "US0378331005" {
			t.Errorf("expected isin=US0378331005, got %s", q.Get("isin"))
		}
		if q.Get("figi") != "BBG000B9XRY4" {
			t.Errorf("expected figi=BBG000B9XRY4, got %s", q.Get("figi"))
		}
		if q.Get("sedol") != "2046251" {
			t.Errorf("expected sedol=2046251, got %s", q.Get("sedol"))
		}
		if q.Get("us_code") != "0378331005" {
			t.Errorf("expected us_code=0378331005, got %s", q.Get("us_code"))
		}
		if q.Get("sort") != "weight.desc" {
			t.Errorf("expected sort=weight.desc, got %s", q.Get("sort"))
		}
		if q.Get("limit") != "100" {
			t.Errorf("expected limit=100, got %s", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(constituentsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetETFGlobalConstituents(ETFGlobalConstituentsParams{
		CompositeTicker:   "SPY",
		ConstituentTicker: "AAPL",
		EffectiveDate:     "2025-12-15",
		ProcessedDate:     "2025-12-16",
		ISIN:              "US0378331005",
		FIGI:              "BBG000B9XRY4",
		SEDOL:             "2046251",
		USCode:            "0378331005",
		Sort:              "weight.desc",
		Limit:             "100",
	})
}

// TestGetETFGlobalConstituentsAPIError verifies that GetETFGlobalConstituents
// returns an error when the API responds with a non-200 status code.
func TestGetETFGlobalConstituentsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"NOT_AUTHORIZED","message":"You are not entitled to this data."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetETFGlobalConstituents(ETFGlobalConstituentsParams{CompositeTicker: "SPY"})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetETFGlobalConstituentsEmptyResults verifies that GetETFGlobalConstituents
// handles an empty results array without error.
func TestGetETFGlobalConstituentsEmptyResults(t *testing.T) {
	emptyJSON := `{"status":"OK","request_id":"abc","count":0,"results":[]}`
	server := mockServer(t, map[string]string{
		"/etf-global/v1/constituents": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetETFGlobalConstituents(ETFGlobalConstituentsParams{CompositeTicker: "ZZZNOTREAL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Count != 0 {
		t.Errorf("expected count 0, got %d", result.Count)
	}

	if len(result.Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(result.Results))
	}
}
