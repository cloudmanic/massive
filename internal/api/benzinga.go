//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

// BenzingaNewsResponse represents the API response for Benzinga news articles.
// It includes pagination support via NextURL and a list of news article results.
type BenzingaNewsResponse struct {
	Status    string                `json:"status"`
	Count     int                   `json:"count"`
	RequestID string                `json:"request_id"`
	NextURL   string                `json:"next_url"`
	Results   []BenzingaNewsArticle `json:"results"`
}

// BenzingaNewsArticle represents a single Benzinga news article returned by the API,
// including metadata such as title, author, body content, associated tickers,
// channels, tags, and sentiment insights.
type BenzingaNewsArticle struct {
	BenzingaID  int                  `json:"benzinga_id"`
	Title       string               `json:"title"`
	Body        string               `json:"body"`
	Teaser      string               `json:"teaser"`
	Author      string               `json:"author"`
	Published   string               `json:"published"`
	LastUpdated string               `json:"last_updated"`
	URL         string               `json:"url"`
	Tickers     []string             `json:"tickers"`
	Channels    []string             `json:"channels"`
	Tags        []string             `json:"tags"`
	Images      []string             `json:"images"`
	Stocks      []string             `json:"stocks"`
	Insights    []BenzingaNewsInsight `json:"insights"`
}

// BenzingaNewsInsight represents a sentiment analysis insight for a specific
// ticker mentioned in a Benzinga news article, including the sentiment
// classification and reasoning behind it.
type BenzingaNewsInsight struct {
	Ticker             string `json:"ticker"`
	Sentiment          string `json:"sentiment"`
	SentimentReasoning string `json:"sentiment_reasoning"`
}

// BenzingaNewsParams holds the query parameters for fetching Benzinga news
// articles from the API. All fields are optional and support various
// filtering options including date ranges, tickers, channels, and tags.
type BenzingaNewsParams struct {
	Tickers      string
	TickersAnyOf string
	Published    string
	PublishedGt  string
	PublishedGte string
	PublishedLt  string
	PublishedLte string
	Channels     string
	Tags         string
	Author       string
	Limit        string
	Sort         string
}

// BenzingaRatingsResponse represents the API response for Benzinga analyst
// ratings data. It includes pagination support via NextURL and a list of
// rating results with analyst details.
type BenzingaRatingsResponse struct {
	Status    string            `json:"status"`
	Count     int               `json:"count"`
	RequestID string            `json:"request_id"`
	NextURL   string            `json:"next_url"`
	Results   []BenzingaRating  `json:"results"`
}

// BenzingaRating represents a single analyst rating entry from the Benzinga
// API, including the analyst name, firm, rating action, price targets,
// and associated metadata.
type BenzingaRating struct {
	BenzingaID                  string  `json:"benzinga_id"`
	Ticker                      string  `json:"ticker"`
	CompanyName                 string  `json:"company_name"`
	Date                        string  `json:"date"`
	Time                        string  `json:"time"`
	Analyst                     string  `json:"analyst"`
	Firm                        string  `json:"firm"`
	Rating                      string  `json:"rating"`
	RatingAction                string  `json:"rating_action"`
	PreviousRating              string  `json:"previous_rating"`
	PriceTarget                 float64 `json:"price_target"`
	PriceTargetAction           string  `json:"price_target_action"`
	PreviousPriceTarget         float64 `json:"previous_price_target"`
	AdjustedPriceTarget         float64 `json:"adjusted_price_target"`
	PreviousAdjustedPriceTarget float64 `json:"previous_adjusted_price_target"`
	PricePercentChange          float64 `json:"price_percent_change"`
	Currency                    string  `json:"currency"`
	Importance                  int     `json:"importance"`
	LastUpdated                 string  `json:"last_updated"`
	Notes                       string  `json:"notes"`
	BenzingaAnalystID           string  `json:"benzinga_analyst_id"`
	BenzingaFirmID              string  `json:"benzinga_firm_id"`
	BenzingaCalendarURL         string  `json:"benzinga_calendar_url"`
	BenzingaNewsURL             string  `json:"benzinga_news_url"`
}

// BenzingaRatingsParams holds the query parameters for fetching Benzinga
// analyst ratings from the API. All fields are optional and support
// filtering by ticker, date range, rating action, and price target action.
type BenzingaRatingsParams struct {
	Ticker            string
	TickerAnyOf       string
	Date              string
	DateGt            string
	DateGte           string
	DateLt            string
	DateLte           string
	Importance        string
	RatingAction      string
	PriceTargetAction string
	Limit             string
	Sort              string
}

// BenzingaEarningsResponse represents the API response for Benzinga earnings
// data. It includes pagination support via NextURL and a list of earnings
// results with EPS and revenue details.
type BenzingaEarningsResponse struct {
	Status    string             `json:"status"`
	Count     int                `json:"count"`
	RequestID string             `json:"request_id"`
	NextURL   string             `json:"next_url"`
	Results   []BenzingaEarnings `json:"results"`
}

// BenzingaEarnings represents a single earnings record from the Benzinga
// API, including actual and estimated EPS/revenue figures, surprise metrics,
// fiscal period information, and reporting metadata.
type BenzingaEarnings struct {
	BenzingaID             string  `json:"benzinga_id"`
	Ticker                 string  `json:"ticker"`
	CompanyName            string  `json:"company_name"`
	Date                   string  `json:"date"`
	Time                   string  `json:"time"`
	DateStatus             string  `json:"date_status"`
	ActualEPS              float64 `json:"actual_eps"`
	EstimatedEPS           float64 `json:"estimated_eps"`
	PreviousEPS            float64 `json:"previous_eps"`
	EPSSurprise            float64 `json:"eps_surprise"`
	EPSSurprisePercent     float64 `json:"eps_surprise_percent"`
	ActualRevenue          float64 `json:"actual_revenue"`
	EstimatedRevenue       float64 `json:"estimated_revenue"`
	PreviousRevenue        float64 `json:"previous_revenue"`
	RevenueSurprise        float64 `json:"revenue_surprise"`
	RevenueSurprisePercent float64 `json:"revenue_surprise_percent"`
	FiscalPeriod           string  `json:"fiscal_period"`
	FiscalYear             int     `json:"fiscal_year"`
	Importance             int     `json:"importance"`
	Currency               string  `json:"currency"`
	EPSMethod              string  `json:"eps_method"`
	RevenueMethod          string  `json:"revenue_method"`
	LastUpdated            string  `json:"last_updated"`
	Notes                  string  `json:"notes"`
}

// BenzingaEarningsParams holds the query parameters for fetching Benzinga
// earnings data from the API. All fields are optional and support filtering
// by ticker, date range, fiscal period, and importance level.
type BenzingaEarningsParams struct {
	Ticker       string
	TickerAnyOf  string
	Date         string
	DateGt       string
	DateGte      string
	DateLt       string
	DateLte      string
	DateStatus   string
	FiscalYear   string
	FiscalPeriod string
	Importance   string
	Limit        string
	Sort         string
}

// BenzingaGuidanceResponse represents the API response for Benzinga corporate
// guidance data. It includes pagination support via NextURL and a list of
// guidance results with EPS and revenue projections.
type BenzingaGuidanceResponse struct {
	Status    string              `json:"status"`
	Count     int                 `json:"count"`
	RequestID string              `json:"request_id"`
	NextURL   string              `json:"next_url"`
	Results   []BenzingaGuidance  `json:"results"`
}

// BenzingaGuidance represents a single corporate guidance record from the
// Benzinga API, including projected EPS and revenue ranges, fiscal period
// information, and company metadata.
type BenzingaGuidance struct {
	BenzingaID                  string  `json:"benzinga_id"`
	Ticker                      string  `json:"ticker"`
	CompanyName                 string  `json:"company_name"`
	Date                        string  `json:"date"`
	Time                        string  `json:"time"`
	Positioning                 string  `json:"positioning"`
	EPSMethod                   string  `json:"eps_method"`
	RevenueMethod               string  `json:"revenue_method"`
	EstimatedEPSGuidance        float64 `json:"estimated_eps_guidance"`
	EstimatedRevenueGuidance    float64 `json:"estimated_revenue_guidance"`
	MinEPSGuidance              float64 `json:"min_eps_guidance"`
	MaxEPSGuidance              float64 `json:"max_eps_guidance"`
	MinRevenueGuidance          float64 `json:"min_revenue_guidance"`
	MaxRevenueGuidance          float64 `json:"max_revenue_guidance"`
	PreviousMinEPSGuidance      float64 `json:"previous_min_eps_guidance"`
	PreviousMaxEPSGuidance      float64 `json:"previous_max_eps_guidance"`
	PreviousMinRevenueGuidance  float64 `json:"previous_min_revenue_guidance"`
	PreviousMaxRevenueGuidance  float64 `json:"previous_max_revenue_guidance"`
	FiscalPeriod                string  `json:"fiscal_period"`
	FiscalYear                  int     `json:"fiscal_year"`
	Importance                  int     `json:"importance"`
	Currency                    string  `json:"currency"`
	ReleaseType                 string  `json:"release_type"`
	LastUpdated                 string  `json:"last_updated"`
	Notes                       string  `json:"notes"`
}

// BenzingaGuidanceParams holds the query parameters for fetching Benzinga
// corporate guidance data from the API. All fields are optional and support
// filtering by ticker, date range, fiscal period, positioning, and importance.
type BenzingaGuidanceParams struct {
	Ticker       string
	TickerAnyOf  string
	Date         string
	DateGt       string
	DateGte      string
	DateLt       string
	DateLte      string
	Positioning  string
	FiscalYear   string
	FiscalPeriod string
	Importance   string
	Limit        string
	Sort         string
}

// BenzingaAnalystsResponse represents the API response for Benzinga analyst
// details data. It includes pagination support via NextURL and a list of
// analyst profiles with performance metrics.
type BenzingaAnalystsResponse struct {
	Status    string             `json:"status"`
	RequestID string             `json:"request_id"`
	NextURL   string             `json:"next_url"`
	Results   []BenzingaAnalyst  `json:"results"`
}

// BenzingaAnalyst represents a single analyst profile from the Benzinga
// API, including the analyst name, firm affiliation, performance metrics
// such as success rate and average return, and a smart score.
type BenzingaAnalyst struct {
	BenzingaID                   string  `json:"benzinga_id"`
	BenzingaFirmID               string  `json:"benzinga_firm_id"`
	FullName                     string  `json:"full_name"`
	FirmName                     string  `json:"firm_name"`
	SmartScore                   float64 `json:"smart_score"`
	OverallSuccessRate           float64 `json:"overall_success_rate"`
	OverallAvgReturn             float64 `json:"overall_avg_return"`
	OverallAvgReturnPercentile   float64 `json:"overall_avg_return_percentile"`
	TotalRatings                 float64 `json:"total_ratings"`
	TotalRatingsPercentile       float64 `json:"total_ratings_percentile"`
	LastUpdated                  string  `json:"last_updated"`
}

// BenzingaAnalystsParams holds the query parameters for fetching Benzinga
// analyst details from the API. All fields are optional and support
// filtering by analyst ID, firm ID, name, and firm name.
type BenzingaAnalystsParams struct {
	BenzingaID     string
	BenzingaFirmID string
	FullName       string
	FirmName       string
	Limit          string
	Sort           string
}

// GetBenzingaNews retrieves Benzinga news articles from the Massive API
// with optional filtering by tickers, publication date range, channels,
// tags, and author. Returns paginated results matching the specified criteria.
func (c *Client) GetBenzingaNews(p BenzingaNewsParams) (*BenzingaNewsResponse, error) {
	path := "/benzinga/v2/news"

	params := map[string]string{
		"tickers":        p.Tickers,
		"tickers.any_of": p.TickersAnyOf,
		"published":      p.Published,
		"published.gt":   p.PublishedGt,
		"published.gte":  p.PublishedGte,
		"published.lt":   p.PublishedLt,
		"published.lte":  p.PublishedLte,
		"channels":       p.Channels,
		"tags":           p.Tags,
		"author":         p.Author,
		"limit":          p.Limit,
		"sort":           p.Sort,
	}

	var result BenzingaNewsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBenzingaRatings retrieves Benzinga analyst ratings from the Massive API
// with optional filtering by ticker, date range, rating action, price target
// action, and importance level. Returns paginated results.
func (c *Client) GetBenzingaRatings(p BenzingaRatingsParams) (*BenzingaRatingsResponse, error) {
	path := "/benzinga/v1/ratings"

	params := map[string]string{
		"ticker":              p.Ticker,
		"ticker.any_of":      p.TickerAnyOf,
		"date":               p.Date,
		"date.gt":            p.DateGt,
		"date.gte":           p.DateGte,
		"date.lt":            p.DateLt,
		"date.lte":           p.DateLte,
		"importance":         p.Importance,
		"rating_action":      p.RatingAction,
		"price_target_action": p.PriceTargetAction,
		"limit":              p.Limit,
		"sort":               p.Sort,
	}

	var result BenzingaRatingsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBenzingaEarnings retrieves Benzinga earnings data from the Massive API
// with optional filtering by ticker, date range, fiscal period, date status,
// and importance level. Returns paginated results with EPS and revenue details.
func (c *Client) GetBenzingaEarnings(p BenzingaEarningsParams) (*BenzingaEarningsResponse, error) {
	path := "/benzinga/v1/earnings"

	params := map[string]string{
		"ticker":        p.Ticker,
		"ticker.any_of": p.TickerAnyOf,
		"date":          p.Date,
		"date.gt":       p.DateGt,
		"date.gte":      p.DateGte,
		"date.lt":       p.DateLt,
		"date.lte":      p.DateLte,
		"date_status":   p.DateStatus,
		"fiscal_year":   p.FiscalYear,
		"fiscal_period": p.FiscalPeriod,
		"importance":    p.Importance,
		"limit":         p.Limit,
		"sort":          p.Sort,
	}

	var result BenzingaEarningsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBenzingaGuidance retrieves Benzinga corporate guidance data from the
// Massive API with optional filtering by ticker, date range, fiscal period,
// positioning, and importance level. Returns paginated results with
// projected EPS and revenue ranges.
func (c *Client) GetBenzingaGuidance(p BenzingaGuidanceParams) (*BenzingaGuidanceResponse, error) {
	path := "/benzinga/v1/guidance"

	params := map[string]string{
		"ticker":        p.Ticker,
		"ticker.any_of": p.TickerAnyOf,
		"date":          p.Date,
		"date.gt":       p.DateGt,
		"date.gte":      p.DateGte,
		"date.lt":       p.DateLt,
		"date.lte":      p.DateLte,
		"positioning":   p.Positioning,
		"fiscal_year":   p.FiscalYear,
		"fiscal_period": p.FiscalPeriod,
		"importance":    p.Importance,
		"limit":         p.Limit,
		"sort":          p.Sort,
	}

	var result BenzingaGuidanceResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBenzingaAnalysts retrieves Benzinga analyst details from the Massive API
// with optional filtering by analyst ID, firm ID, name, and firm name.
// Returns paginated results with analyst performance metrics.
func (c *Client) GetBenzingaAnalysts(p BenzingaAnalystsParams) (*BenzingaAnalystsResponse, error) {
	path := "/benzinga/v1/analysts"

	params := map[string]string{
		"benzinga_id":      p.BenzingaID,
		"benzinga_firm_id": p.BenzingaFirmID,
		"full_name":        p.FullName,
		"firm_name":        p.FirmName,
		"limit":            p.Limit,
		"sort":             p.Sort,
	}

	var result BenzingaAnalystsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
