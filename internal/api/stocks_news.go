//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

// NewsResponse represents the API response for stock news articles.
// It includes pagination support via NextURL and a list of news results.
type NewsResponse struct {
	Status    string        `json:"status"`
	Count     int           `json:"count"`
	RequestID string        `json:"request_id"`
	NextURL   string        `json:"next_url"`
	Results   []NewsArticle `json:"results"`
}

// NewsArticle represents a single news article returned by the API,
// including metadata such as title, author, publisher, associated
// tickers, keywords, and sentiment insights.
type NewsArticle struct {
	ID           string         `json:"id"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	ArticleURL   string         `json:"article_url"`
	AmpURL       string         `json:"amp_url"`
	Author       string         `json:"author"`
	PublishedUTC string         `json:"published_utc"`
	ImageURL     string         `json:"image_url"`
	Keywords     []string       `json:"keywords"`
	Tickers      []string       `json:"tickers"`
	Insights     []NewsInsight  `json:"insights"`
	Publisher    NewsPublisher  `json:"publisher"`
}

// NewsInsight represents a sentiment analysis insight for a specific
// ticker mentioned in a news article. It includes the sentiment
// classification and reasoning.
type NewsInsight struct {
	Ticker             string `json:"ticker"`
	Sentiment          string `json:"sentiment"`
	SentimentReasoning string `json:"sentiment_reasoning"`
}

// NewsPublisher represents the publishing source of a news article
// including its name, homepage, logo, and favicon URLs.
type NewsPublisher struct {
	Name        string `json:"name"`
	HomepageURL string `json:"homepage_url"`
	LogoURL     string `json:"logo_url"`
	FaviconURL  string `json:"favicon_url"`
}

// NewsParams holds the query parameters for fetching stock news
// from the reference news endpoint. All fields are optional.
type NewsParams struct {
	Ticker          string
	PublishedUTC    string
	PublishedUTCGte string
	PublishedUTCLte string
	Order           string
	Limit           string
	Sort            string
}

// GetNews retrieves stock news articles from the Massive API with
// optional filtering by ticker symbol, publication date range, and
// sorting. Returns paginated results matching the specified criteria.
func (c *Client) GetNews(p NewsParams) (*NewsResponse, error) {
	path := "/v2/reference/news"

	params := map[string]string{
		"ticker":            p.Ticker,
		"published_utc":     p.PublishedUTC,
		"published_utc.gte": p.PublishedUTCGte,
		"published_utc.lte": p.PublishedUTCLte,
		"order":             p.Order,
		"limit":             p.Limit,
		"sort":              p.Sort,
	}

	var result NewsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
