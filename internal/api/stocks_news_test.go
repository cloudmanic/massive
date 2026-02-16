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

const newsJSON = `{
	"results": [
		{
			"id": "abc123",
			"publisher": {
				"name": "Benzinga",
				"homepage_url": "https://www.benzinga.com/",
				"logo_url": "https://s3.massive.com/public/assets/news/logos/benzinga.svg",
				"favicon_url": "https://s3.massive.com/public/assets/news/favicons/benzinga.ico"
			},
			"title": "Apple Reports Record Earnings",
			"author": "Jane Doe",
			"published_utc": "2026-01-15T10:30:00Z",
			"article_url": "https://www.benzinga.com/apple-earnings",
			"amp_url": "https://www.benzinga.com/amp/apple-earnings",
			"tickers": ["AAPL", "MSFT"],
			"image_url": "https://cdn.benzinga.com/apple.jpg",
			"description": "Apple reported record quarterly earnings driven by strong iPhone sales.",
			"keywords": ["earnings", "tech", "iPhone"],
			"insights": [
				{
					"ticker": "AAPL",
					"sentiment": "positive",
					"sentiment_reasoning": "Record quarterly earnings driven by strong iPhone sales"
				},
				{
					"ticker": "MSFT",
					"sentiment": "neutral",
					"sentiment_reasoning": "Mentioned as competitor in cloud services"
				}
			]
		},
		{
			"id": "def456",
			"publisher": {
				"name": "MarketWatch",
				"homepage_url": "https://www.marketwatch.com/",
				"logo_url": "https://s3.massive.com/public/assets/news/logos/marketwatch.svg",
				"favicon_url": "https://s3.massive.com/public/assets/news/favicons/marketwatch.ico"
			},
			"title": "Tech Stocks Rally After Fed Decision",
			"author": "John Smith",
			"published_utc": "2026-01-14T14:00:00Z",
			"article_url": "https://www.marketwatch.com/tech-rally",
			"tickers": ["AAPL", "GOOG", "AMZN"],
			"description": "Major tech stocks surged following the Federal Reserve decision to hold rates steady.",
			"keywords": ["fed", "rates", "rally"],
			"insights": [
				{
					"ticker": "AAPL",
					"sentiment": "positive",
					"sentiment_reasoning": "Part of broad tech rally after Fed rate decision"
				}
			]
		}
	],
	"status": "OK",
	"request_id": "news123",
	"count": 2,
	"next_url": "https://api.massive.com/v2/reference/news?cursor=YXA9Mg"
}`

// TestGetNews verifies that GetNews correctly parses the API response
// and returns the expected news articles with all fields populated.
func TestGetNews(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/reference/news": newsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := NewsParams{
		Ticker: "AAPL",
		Limit:  "2",
	}

	result, err := client.GetNews(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Count != 2 {
		t.Errorf("expected count 2, got %d", result.Count)
	}

	if result.RequestID != "news123" {
		t.Errorf("expected request_id news123, got %s", result.RequestID)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 articles, got %d", len(result.Results))
	}
}

// TestGetNewsFirstArticle verifies that the first article in the
// response is correctly parsed with all metadata fields.
func TestGetNewsFirstArticle(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/reference/news": newsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetNews(NewsParams{Ticker: "AAPL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	article := result.Results[0]

	if article.ID != "abc123" {
		t.Errorf("expected id abc123, got %s", article.ID)
	}

	if article.Title != "Apple Reports Record Earnings" {
		t.Errorf("expected title 'Apple Reports Record Earnings', got %s", article.Title)
	}

	if article.Author != "Jane Doe" {
		t.Errorf("expected author Jane Doe, got %s", article.Author)
	}

	if article.PublishedUTC != "2026-01-15T10:30:00Z" {
		t.Errorf("expected published_utc 2026-01-15T10:30:00Z, got %s", article.PublishedUTC)
	}

	if article.ArticleURL != "https://www.benzinga.com/apple-earnings" {
		t.Errorf("expected article_url https://www.benzinga.com/apple-earnings, got %s", article.ArticleURL)
	}

	if article.AmpURL != "https://www.benzinga.com/amp/apple-earnings" {
		t.Errorf("expected amp_url https://www.benzinga.com/amp/apple-earnings, got %s", article.AmpURL)
	}

	if article.ImageURL != "https://cdn.benzinga.com/apple.jpg" {
		t.Errorf("expected image_url https://cdn.benzinga.com/apple.jpg, got %s", article.ImageURL)
	}

	if article.Description != "Apple reported record quarterly earnings driven by strong iPhone sales." {
		t.Errorf("unexpected description: %s", article.Description)
	}
}

// TestGetNewsPublisher verifies that the publisher information is
// correctly parsed from the first article in the response.
func TestGetNewsPublisher(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/reference/news": newsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetNews(NewsParams{Ticker: "AAPL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pub := result.Results[0].Publisher

	if pub.Name != "Benzinga" {
		t.Errorf("expected publisher name Benzinga, got %s", pub.Name)
	}

	if pub.HomepageURL != "https://www.benzinga.com/" {
		t.Errorf("expected homepage_url https://www.benzinga.com/, got %s", pub.HomepageURL)
	}

	if pub.LogoURL != "https://s3.massive.com/public/assets/news/logos/benzinga.svg" {
		t.Errorf("expected logo_url, got %s", pub.LogoURL)
	}

	if pub.FaviconURL != "https://s3.massive.com/public/assets/news/favicons/benzinga.ico" {
		t.Errorf("expected favicon_url, got %s", pub.FaviconURL)
	}
}

// TestGetNewsTickers verifies that the tickers array is correctly
// parsed from the first article in the response.
func TestGetNewsTickers(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/reference/news": newsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetNews(NewsParams{Ticker: "AAPL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tickers := result.Results[0].Tickers

	if len(tickers) != 2 {
		t.Fatalf("expected 2 tickers, got %d", len(tickers))
	}

	if tickers[0] != "AAPL" {
		t.Errorf("expected first ticker AAPL, got %s", tickers[0])
	}

	if tickers[1] != "MSFT" {
		t.Errorf("expected second ticker MSFT, got %s", tickers[1])
	}
}

// TestGetNewsKeywords verifies that the keywords array is correctly
// parsed from the first article in the response.
func TestGetNewsKeywords(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/reference/news": newsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetNews(NewsParams{Ticker: "AAPL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	keywords := result.Results[0].Keywords

	if len(keywords) != 3 {
		t.Fatalf("expected 3 keywords, got %d", len(keywords))
	}

	if keywords[0] != "earnings" {
		t.Errorf("expected first keyword 'earnings', got %s", keywords[0])
	}

	if keywords[1] != "tech" {
		t.Errorf("expected second keyword 'tech', got %s", keywords[1])
	}

	if keywords[2] != "iPhone" {
		t.Errorf("expected third keyword 'iPhone', got %s", keywords[2])
	}
}

// TestGetNewsInsights verifies that sentiment insights are correctly
// parsed from the first article in the response.
func TestGetNewsInsights(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/reference/news": newsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetNews(NewsParams{Ticker: "AAPL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	insights := result.Results[0].Insights

	if len(insights) != 2 {
		t.Fatalf("expected 2 insights, got %d", len(insights))
	}

	if insights[0].Ticker != "AAPL" {
		t.Errorf("expected insight ticker AAPL, got %s", insights[0].Ticker)
	}

	if insights[0].Sentiment != "positive" {
		t.Errorf("expected sentiment positive, got %s", insights[0].Sentiment)
	}

	if insights[0].SentimentReasoning != "Record quarterly earnings driven by strong iPhone sales" {
		t.Errorf("unexpected sentiment_reasoning: %s", insights[0].SentimentReasoning)
	}

	if insights[1].Ticker != "MSFT" {
		t.Errorf("expected insight ticker MSFT, got %s", insights[1].Ticker)
	}

	if insights[1].Sentiment != "neutral" {
		t.Errorf("expected sentiment neutral, got %s", insights[1].Sentiment)
	}
}

// TestGetNewsSecondArticle verifies that the second article in the
// response is correctly parsed with its own distinct values.
func TestGetNewsSecondArticle(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/v2/reference/news": newsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetNews(NewsParams{Ticker: "AAPL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	article := result.Results[1]

	if article.ID != "def456" {
		t.Errorf("expected id def456, got %s", article.ID)
	}

	if article.Title != "Tech Stocks Rally After Fed Decision" {
		t.Errorf("expected title 'Tech Stocks Rally After Fed Decision', got %s", article.Title)
	}

	if article.Publisher.Name != "MarketWatch" {
		t.Errorf("expected publisher MarketWatch, got %s", article.Publisher.Name)
	}

	if len(article.Tickers) != 3 {
		t.Fatalf("expected 3 tickers, got %d", len(article.Tickers))
	}

	if article.Tickers[2] != "AMZN" {
		t.Errorf("expected third ticker AMZN, got %s", article.Tickers[2])
	}
}

// TestGetNewsQueryParams verifies that all filter parameters are
// correctly sent to the API endpoint as query parameters.
func TestGetNewsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("ticker") != "AAPL" {
			t.Errorf("expected ticker=AAPL, got %s", q.Get("ticker"))
		}
		if q.Get("order") != "desc" {
			t.Errorf("expected order=desc, got %s", q.Get("order"))
		}
		if q.Get("limit") != "5" {
			t.Errorf("expected limit=5, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "published_utc" {
			t.Errorf("expected sort=published_utc, got %s", q.Get("sort"))
		}
		if q.Get("published_utc.gte") != "2026-01-01" {
			t.Errorf("expected published_utc.gte=2026-01-01, got %s", q.Get("published_utc.gte"))
		}
		if q.Get("published_utc.lte") != "2026-01-31" {
			t.Errorf("expected published_utc.lte=2026-01-31, got %s", q.Get("published_utc.lte"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(newsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetNews(NewsParams{
		Ticker:          "AAPL",
		Order:           "desc",
		Limit:           "5",
		Sort:            "published_utc",
		PublishedUTCGte: "2026-01-01",
		PublishedUTCLte: "2026-01-31",
	})
}

// TestGetNewsRequestPath verifies that GetNews sends requests to
// the correct API path.
func TestGetNewsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(newsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetNews(NewsParams{Ticker: "AAPL"})

	if receivedPath != "/v2/reference/news" {
		t.Errorf("expected path /v2/reference/news, got %s", receivedPath)
	}
}

// TestGetNewsAPIError verifies that GetNews returns an error when
// the API responds with a non-200 status code.
func TestGetNewsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"status":"ERROR","message":"Unauthorized"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetNews(NewsParams{Ticker: "AAPL"})
	if err == nil {
		t.Fatal("expected error for 401 response, got nil")
	}
}

// TestGetNewsEmptyResults verifies that GetNews handles an empty
// results array without error.
func TestGetNewsEmptyResults(t *testing.T) {
	emptyJSON := `{"results":[],"status":"OK","request_id":"abc","count":0}`
	server := mockServer(t, map[string]string{
		"/v2/reference/news": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetNews(NewsParams{Ticker: "ZZZZZZ"})
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

// TestGetNewsPublishedUTCParam verifies that the published_utc
// parameter is correctly sent when provided.
func TestGetNewsPublishedUTCParam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("published_utc") != "2026-01-15" {
			t.Errorf("expected published_utc=2026-01-15, got %s", r.URL.Query().Get("published_utc"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(newsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetNews(NewsParams{PublishedUTC: "2026-01-15"})
}
