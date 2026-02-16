//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

// ETFGlobalAnalyticsResponse represents the API response for the ETF Global
// analytics endpoint, which returns quantitative scoring, risk, and reward
// data for exchange-traded funds.
type ETFGlobalAnalyticsResponse struct {
	Status    string               `json:"status"`
	RequestID string               `json:"request_id"`
	Count     int                  `json:"count"`
	NextURL   string               `json:"next_url"`
	Results   []ETFGlobalAnalytics `json:"results"`
}

// ETFGlobalAnalytics represents a single ETF's analytics record containing
// comprehensive quantitative scores, risk metrics, and reward assessments
// provided by ETF Global.
type ETFGlobalAnalytics struct {
	CompositeTicker           string  `json:"composite_ticker"`
	EffectiveDate             string  `json:"effective_date"`
	ProcessedDate             string  `json:"processed_date"`
	QuantCompositeBehavioral  float64 `json:"quant_composite_behavioral"`
	QuantCompositeFundamental float64 `json:"quant_composite_fundamental"`
	QuantCompositeGlobal      float64 `json:"quant_composite_global"`
	QuantCompositeQuality     float64 `json:"quant_composite_quality"`
	QuantCompositeSentiment   float64 `json:"quant_composite_sentiment"`
	QuantCompositeTechnical   float64 `json:"quant_composite_technical"`
	QuantFundamentalDiv       float64 `json:"quant_fundamental_div"`
	QuantFundamentalPB        float64 `json:"quant_fundamental_pb"`
	QuantFundamentalPCF       float64 `json:"quant_fundamental_pcf"`
	QuantFundamentalPE        float64 `json:"quant_fundamental_pe"`
	QuantGlobalCountry        float64 `json:"quant_global_country"`
	QuantGlobalSector         float64 `json:"quant_global_sector"`
	QuantGrade                string  `json:"quant_grade"`
	QuantQualityDiversify     float64 `json:"quant_quality_diversification"`
	QuantQualityFirm          float64 `json:"quant_quality_firm"`
	QuantQualityLiquidity     float64 `json:"quant_quality_liquidity"`
	QuantSentimentIV          float64 `json:"quant_sentiment_iv"`
	QuantSentimentPC          float64 `json:"quant_sentiment_pc"`
	QuantSentimentSI          float64 `json:"quant_sentiment_si"`
	QuantTechnicalIT          float64 `json:"quant_technical_it"`
	QuantTechnicalLT          float64 `json:"quant_technical_lt"`
	QuantTechnicalST          float64 `json:"quant_technical_st"`
	QuantTotalScore           float64 `json:"quant_total_score"`
	RewardScore               float64 `json:"reward_score"`
	RiskCountry               float64 `json:"risk_country"`
	RiskDeviation             float64 `json:"risk_deviation"`
	RiskEfficiency            float64 `json:"risk_efficiency"`
	RiskLiquidity             float64 `json:"risk_liquidity"`
	RiskStructure             float64 `json:"risk_structure"`
	RiskTotalScore            float64 `json:"risk_total_score"`
	RiskVolatility            float64 `json:"risk_volatility"`
}

// ETFGlobalAnalyticsParams holds the query parameters for fetching ETF Global
// analytics data. All fields are optional and support comparison operators
// (e.g., composite_ticker.any_of, risk_total_score.gte).
type ETFGlobalAnalyticsParams struct {
	CompositeTicker string
	ProcessedDate   string
	EffectiveDate   string
	RiskTotalScore  string
	RewardScore     string
	QuantTotalScore string
	QuantGrade      string
	Sort            string
	Limit           string
}

// ETFGlobalConstituentsResponse represents the API response for the ETF Global
// constituents endpoint, which returns the underlying holdings of an ETF
// including weight, shares held, and security identifiers.
type ETFGlobalConstituentsResponse struct {
	Status    string                  `json:"status"`
	RequestID string                  `json:"request_id"`
	Count     int                     `json:"count"`
	NextURL   string                  `json:"next_url"`
	Results   []ETFGlobalConstituent  `json:"results"`
}

// ETFGlobalConstituent represents a single constituent holding within an ETF,
// including its weight, market value, and various security identifiers such
// as ISIN, FIGI, SEDOL, and US code.
type ETFGlobalConstituent struct {
	AssetClass        string  `json:"asset_class"`
	CompositeTicker   string  `json:"composite_ticker"`
	ConstituentName   string  `json:"constituent_name"`
	ConstituentRank   int     `json:"constituent_rank"`
	ConstituentTicker string  `json:"constituent_ticker"`
	CountryOfExchange string  `json:"country_of_exchange"`
	CurrencyTraded    string  `json:"currency_traded"`
	EffectiveDate     string  `json:"effective_date"`
	Exchange          string  `json:"exchange"`
	FIGI              string  `json:"figi"`
	ISIN              string  `json:"isin"`
	MarketValue       float64 `json:"market_value"`
	ProcessedDate     string  `json:"processed_date"`
	SecurityType      string  `json:"security_type"`
	SEDOL             string  `json:"sedol"`
	SharesHeld        float64 `json:"shares_held"`
	USCode            string  `json:"us_code"`
	Weight            float64 `json:"weight"`
}

// ETFGlobalConstituentsParams holds the query parameters for fetching ETF Global
// constituent holdings. Supports filtering by composite ticker, constituent
// ticker, effective date, and various security identifiers.
type ETFGlobalConstituentsParams struct {
	CompositeTicker   string
	ConstituentTicker string
	EffectiveDate     string
	ProcessedDate     string
	USCode            string
	ISIN              string
	FIGI              string
	SEDOL             string
	Sort              string
	Limit             string
}

// GetETFGlobalAnalytics retrieves ETF Global analytics data including
// quantitative scores, risk assessments, and reward metrics for ETFs
// matching the specified filter criteria.
func (c *Client) GetETFGlobalAnalytics(p ETFGlobalAnalyticsParams) (*ETFGlobalAnalyticsResponse, error) {
	path := "/etf-global/v1/analytics"

	params := map[string]string{
		"composite_ticker": p.CompositeTicker,
		"processed_date":   p.ProcessedDate,
		"effective_date":   p.EffectiveDate,
		"risk_total_score": p.RiskTotalScore,
		"reward_score":     p.RewardScore,
		"quant_total_score": p.QuantTotalScore,
		"quant_grade":      p.QuantGrade,
		"sort":             p.Sort,
		"limit":            p.Limit,
	}

	var result ETFGlobalAnalyticsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetETFGlobalConstituents retrieves the underlying constituent holdings
// of ETFs from ETF Global, including weights, market values, share counts,
// and security identifiers for each position.
func (c *Client) GetETFGlobalConstituents(p ETFGlobalConstituentsParams) (*ETFGlobalConstituentsResponse, error) {
	path := "/etf-global/v1/constituents"

	params := map[string]string{
		"composite_ticker":   p.CompositeTicker,
		"constituent_ticker": p.ConstituentTicker,
		"effective_date":     p.EffectiveDate,
		"processed_date":     p.ProcessedDate,
		"us_code":            p.USCode,
		"isin":               p.ISIN,
		"figi":               p.FIGI,
		"sedol":              p.SEDOL,
		"sort":               p.Sort,
		"limit":              p.Limit,
	}

	var result ETFGlobalConstituentsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
