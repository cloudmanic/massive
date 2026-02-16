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

const benzingaNewsJSON = `{
	"results": [
		{
			"benzinga_id": 12345,
			"title": "Apple Unveils New AI Features",
			"body": "<p>Apple today announced a suite of new AI-powered features...</p>",
			"teaser": "Apple announced new AI features for iPhone and Mac.",
			"author": "Jane Smith",
			"published": "2026-02-10T14:30:00Z",
			"last_updated": "2026-02-10T15:00:00Z",
			"url": "https://www.benzinga.com/tech/apple-ai-features",
			"tickers": ["AAPL", "GOOG"],
			"channels": ["Tech", "Large Cap"],
			"tags": ["AI", "Apple", "Technology"],
			"images": ["https://cdn.benzinga.com/apple-ai.jpg"],
			"stocks": ["AAPL"],
			"insights": [
				{
					"ticker": "AAPL",
					"sentiment": "positive",
					"sentiment_reasoning": "New AI features expected to boost iPhone sales"
				}
			]
		},
		{
			"benzinga_id": 12346,
			"title": "Google Responds to Apple AI Push",
			"body": "<p>Google announced competing features...</p>",
			"teaser": "Google unveils its own AI response.",
			"author": "John Doe",
			"published": "2026-02-10T16:00:00Z",
			"last_updated": "2026-02-10T16:30:00Z",
			"url": "https://www.benzinga.com/tech/google-ai-response",
			"tickers": ["GOOG", "AAPL"],
			"channels": ["Tech"],
			"tags": ["AI", "Google"],
			"images": [],
			"stocks": ["GOOG"],
			"insights": []
		}
	],
	"status": "OK",
	"request_id": "bznews123",
	"count": 2,
	"next_url": "https://api.massive.com/benzinga/v2/news?cursor=abc"
}`

const benzingaRatingsJSON = `{
	"results": [
		{
			"benzinga_id": "bz-rating-001",
			"ticker": "AAPL",
			"company_name": "Apple Inc.",
			"date": "2026-02-10",
			"time": "08:30:00",
			"analyst": "John Analyst",
			"firm": "Goldman Sachs",
			"rating": "Buy",
			"rating_action": "upgrades",
			"previous_rating": "Hold",
			"price_target": 250.00,
			"price_target_action": "raises",
			"previous_price_target": 200.00,
			"adjusted_price_target": 248.50,
			"previous_adjusted_price_target": 198.50,
			"price_percent_change": 25.0,
			"currency": "USD",
			"importance": 4,
			"last_updated": "2026-02-10T09:00:00Z",
			"notes": "Upgraded based on strong iPhone sales outlook",
			"benzinga_analyst_id": "analyst-001",
			"benzinga_firm_id": "firm-gs-001",
			"benzinga_calendar_url": "https://www.benzinga.com/calendars/ratings/aapl",
			"benzinga_news_url": "https://www.benzinga.com/news/ratings/aapl"
		},
		{
			"benzinga_id": "bz-rating-002",
			"ticker": "MSFT",
			"company_name": "Microsoft Corporation",
			"date": "2026-02-09",
			"time": "09:15:00",
			"analyst": "Sarah Analyst",
			"firm": "Morgan Stanley",
			"rating": "Overweight",
			"rating_action": "maintains",
			"previous_rating": "Overweight",
			"price_target": 480.00,
			"price_target_action": "maintains",
			"previous_price_target": 480.00,
			"adjusted_price_target": 478.00,
			"previous_adjusted_price_target": 478.00,
			"price_percent_change": 0.0,
			"currency": "USD",
			"importance": 3,
			"last_updated": "2026-02-09T10:00:00Z",
			"notes": "",
			"benzinga_analyst_id": "analyst-002",
			"benzinga_firm_id": "firm-ms-001",
			"benzinga_calendar_url": "",
			"benzinga_news_url": ""
		}
	],
	"status": "OK",
	"request_id": "bzratings123",
	"count": 2,
	"next_url": "https://api.massive.com/benzinga/v1/ratings?cursor=xyz"
}`

const benzingaEarningsJSON = `{
	"results": [
		{
			"benzinga_id": "bz-earn-001",
			"ticker": "AAPL",
			"company_name": "Apple Inc.",
			"date": "2026-01-30",
			"time": "16:30:00",
			"date_status": "confirmed",
			"actual_eps": 2.18,
			"estimated_eps": 2.10,
			"previous_eps": 1.95,
			"eps_surprise": 0.08,
			"eps_surprise_percent": 3.81,
			"actual_revenue": 124500000000,
			"estimated_revenue": 121000000000,
			"previous_revenue": 119500000000,
			"revenue_surprise": 3500000000,
			"revenue_surprise_percent": 2.89,
			"fiscal_period": "Q1",
			"fiscal_year": 2026,
			"importance": 5,
			"currency": "USD",
			"eps_method": "gaap",
			"revenue_method": "gaap",
			"last_updated": "2026-01-30T22:00:00Z",
			"notes": "Record quarterly revenue"
		},
		{
			"benzinga_id": "bz-earn-002",
			"ticker": "MSFT",
			"company_name": "Microsoft Corporation",
			"date": "2026-02-15",
			"time": "00:00:00",
			"date_status": "projected",
			"actual_eps": 0,
			"estimated_eps": 3.25,
			"previous_eps": 3.10,
			"eps_surprise": 0,
			"eps_surprise_percent": 0,
			"actual_revenue": 0,
			"estimated_revenue": 69800000000,
			"previous_revenue": 65600000000,
			"revenue_surprise": 0,
			"revenue_surprise_percent": 0,
			"fiscal_period": "Q2",
			"fiscal_year": 2026,
			"importance": 5,
			"currency": "USD",
			"eps_method": "gaap",
			"revenue_method": "gaap",
			"last_updated": "2026-02-14T12:00:00Z",
			"notes": ""
		}
	],
	"status": "OK",
	"request_id": "bzearn123",
	"count": 2,
	"next_url": "https://api.massive.com/benzinga/v1/earnings?cursor=def"
}`

const benzingaGuidanceJSON = `{
	"results": [
		{
			"benzinga_id": "bz-guide-001",
			"ticker": "AAPL",
			"company_name": "Apple Inc.",
			"date": "2026-01-30",
			"time": "16:30:00",
			"positioning": "primary",
			"eps_method": "gaap",
			"revenue_method": "gaap",
			"estimated_eps_guidance": 2.25,
			"estimated_revenue_guidance": 128000000000,
			"min_eps_guidance": 2.15,
			"max_eps_guidance": 2.35,
			"min_revenue_guidance": 125000000000,
			"max_revenue_guidance": 131000000000,
			"previous_min_eps_guidance": 2.00,
			"previous_max_eps_guidance": 2.20,
			"previous_min_revenue_guidance": 120000000000,
			"previous_max_revenue_guidance": 126000000000,
			"fiscal_period": "Q2",
			"fiscal_year": 2026,
			"importance": 5,
			"currency": "USD",
			"release_type": "official",
			"last_updated": "2026-01-30T22:00:00Z",
			"notes": "Management raised guidance"
		}
	],
	"status": "OK",
	"request_id": "bzguide123",
	"count": 1,
	"next_url": ""
}`

const benzingaAnalystsJSON = `{
	"results": [
		{
			"benzinga_id": "analyst-001",
			"benzinga_firm_id": "firm-gs-001",
			"full_name": "John Analyst",
			"firm_name": "Goldman Sachs",
			"smart_score": 85.5,
			"overall_success_rate": 0.72,
			"overall_avg_return": 12.5,
			"overall_avg_return_percentile": 88.0,
			"total_ratings": 350,
			"total_ratings_percentile": 92.0,
			"last_updated": "2026-02-15T08:00:00Z"
		},
		{
			"benzinga_id": "analyst-002",
			"benzinga_firm_id": "firm-ms-001",
			"full_name": "Sarah Analyst",
			"firm_name": "Morgan Stanley",
			"smart_score": 78.2,
			"overall_success_rate": 0.65,
			"overall_avg_return": 9.8,
			"overall_avg_return_percentile": 75.0,
			"total_ratings": 280,
			"total_ratings_percentile": 85.0,
			"last_updated": "2026-02-14T10:00:00Z"
		}
	],
	"status": "OK",
	"request_id": "bzanalysts123",
	"next_url": ""
}`

// TestGetBenzingaNews verifies that GetBenzingaNews correctly parses the
// API response and returns the expected news articles with all fields populated.
func TestGetBenzingaNews(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/benzinga/v2/news": benzingaNewsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := BenzingaNewsParams{
		Tickers: "AAPL",
		Limit:   "2",
	}

	result, err := client.GetBenzingaNews(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Count != 2 {
		t.Errorf("expected count 2, got %d", result.Count)
	}

	if result.RequestID != "bznews123" {
		t.Errorf("expected request_id bznews123, got %s", result.RequestID)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 articles, got %d", len(result.Results))
	}
}

// TestGetBenzingaNewsFirstArticle verifies that the first article in the
// response is correctly parsed with all metadata fields including title,
// author, published date, tickers, channels, tags, and insights.
func TestGetBenzingaNewsFirstArticle(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/benzinga/v2/news": benzingaNewsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetBenzingaNews(BenzingaNewsParams{Tickers: "AAPL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	article := result.Results[0]

	if article.BenzingaID != 12345 {
		t.Errorf("expected benzinga_id 12345, got %d", article.BenzingaID)
	}

	if article.Title != "Apple Unveils New AI Features" {
		t.Errorf("expected title 'Apple Unveils New AI Features', got %s", article.Title)
	}

	if article.Author != "Jane Smith" {
		t.Errorf("expected author Jane Smith, got %s", article.Author)
	}

	if article.Published != "2026-02-10T14:30:00Z" {
		t.Errorf("expected published 2026-02-10T14:30:00Z, got %s", article.Published)
	}

	if article.URL != "https://www.benzinga.com/tech/apple-ai-features" {
		t.Errorf("expected url https://www.benzinga.com/tech/apple-ai-features, got %s", article.URL)
	}

	if article.Teaser != "Apple announced new AI features for iPhone and Mac." {
		t.Errorf("unexpected teaser: %s", article.Teaser)
	}

	if len(article.Tickers) != 2 {
		t.Fatalf("expected 2 tickers, got %d", len(article.Tickers))
	}

	if article.Tickers[0] != "AAPL" {
		t.Errorf("expected first ticker AAPL, got %s", article.Tickers[0])
	}

	if len(article.Channels) != 2 {
		t.Fatalf("expected 2 channels, got %d", len(article.Channels))
	}

	if article.Channels[0] != "Tech" {
		t.Errorf("expected first channel Tech, got %s", article.Channels[0])
	}

	if len(article.Tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(article.Tags))
	}

	if article.Tags[0] != "AI" {
		t.Errorf("expected first tag AI, got %s", article.Tags[0])
	}
}

// TestGetBenzingaNewsInsights verifies that sentiment insights are correctly
// parsed from the first Benzinga news article in the response.
func TestGetBenzingaNewsInsights(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/benzinga/v2/news": benzingaNewsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetBenzingaNews(BenzingaNewsParams{Tickers: "AAPL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	insights := result.Results[0].Insights

	if len(insights) != 1 {
		t.Fatalf("expected 1 insight, got %d", len(insights))
	}

	if insights[0].Ticker != "AAPL" {
		t.Errorf("expected insight ticker AAPL, got %s", insights[0].Ticker)
	}

	if insights[0].Sentiment != "positive" {
		t.Errorf("expected sentiment positive, got %s", insights[0].Sentiment)
	}

	if insights[0].SentimentReasoning != "New AI features expected to boost iPhone sales" {
		t.Errorf("unexpected sentiment_reasoning: %s", insights[0].SentimentReasoning)
	}
}

// TestGetBenzingaNewsSecondArticle verifies that the second article in the
// response is correctly parsed with its own distinct values.
func TestGetBenzingaNewsSecondArticle(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/benzinga/v2/news": benzingaNewsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetBenzingaNews(BenzingaNewsParams{Tickers: "GOOG"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	article := result.Results[1]

	if article.BenzingaID != 12346 {
		t.Errorf("expected benzinga_id 12346, got %d", article.BenzingaID)
	}

	if article.Title != "Google Responds to Apple AI Push" {
		t.Errorf("expected title 'Google Responds to Apple AI Push', got %s", article.Title)
	}

	if article.Author != "John Doe" {
		t.Errorf("expected author John Doe, got %s", article.Author)
	}

	if len(article.Insights) != 0 {
		t.Errorf("expected 0 insights, got %d", len(article.Insights))
	}
}

// TestGetBenzingaNewsRequestPath verifies that GetBenzingaNews sends
// requests to the correct API path.
func TestGetBenzingaNewsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(benzingaNewsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetBenzingaNews(BenzingaNewsParams{Tickers: "AAPL"})

	if receivedPath != "/benzinga/v2/news" {
		t.Errorf("expected path /benzinga/v2/news, got %s", receivedPath)
	}
}

// TestGetBenzingaNewsQueryParams verifies that all filter parameters are
// correctly sent to the Benzinga news API endpoint as query parameters.
func TestGetBenzingaNewsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("tickers") != "AAPL" {
			t.Errorf("expected tickers=AAPL, got %s", q.Get("tickers"))
		}
		if q.Get("published.gte") != "2026-02-01" {
			t.Errorf("expected published.gte=2026-02-01, got %s", q.Get("published.gte"))
		}
		if q.Get("published.lte") != "2026-02-15" {
			t.Errorf("expected published.lte=2026-02-15, got %s", q.Get("published.lte"))
		}
		if q.Get("channels") != "Tech" {
			t.Errorf("expected channels=Tech, got %s", q.Get("channels"))
		}
		if q.Get("limit") != "10" {
			t.Errorf("expected limit=10, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "published.desc" {
			t.Errorf("expected sort=published.desc, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(benzingaNewsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetBenzingaNews(BenzingaNewsParams{
		Tickers:      "AAPL",
		PublishedGte: "2026-02-01",
		PublishedLte: "2026-02-15",
		Channels:     "Tech",
		Limit:        "10",
		Sort:         "published.desc",
	})
}

// TestGetBenzingaNewsAPIError verifies that GetBenzingaNews returns an error
// when the API responds with a non-200 status code.
func TestGetBenzingaNewsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"NOT_AUTHORIZED","message":"You are not entitled to this data."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetBenzingaNews(BenzingaNewsParams{Tickers: "AAPL"})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetBenzingaNewsEmptyResults verifies that GetBenzingaNews handles
// an empty results array without error.
func TestGetBenzingaNewsEmptyResults(t *testing.T) {
	emptyJSON := `{"results":[],"status":"OK","request_id":"abc","count":0}`
	server := mockServer(t, map[string]string{
		"/benzinga/v2/news": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetBenzingaNews(BenzingaNewsParams{Tickers: "ZZZZZZ"})
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

// TestGetBenzingaRatings verifies that GetBenzingaRatings correctly parses
// the API response and returns the expected analyst ratings data.
func TestGetBenzingaRatings(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/benzinga/v1/ratings": benzingaRatingsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := BenzingaRatingsParams{
		Ticker: "AAPL",
		Limit:  "2",
	}

	result, err := client.GetBenzingaRatings(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Count != 2 {
		t.Errorf("expected count 2, got %d", result.Count)
	}

	if result.RequestID != "bzratings123" {
		t.Errorf("expected request_id bzratings123, got %s", result.RequestID)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 ratings, got %d", len(result.Results))
	}
}

// TestGetBenzingaRatingsFirstRating verifies that the first rating in the
// response is correctly parsed with all fields including analyst, firm,
// rating action, price target, and importance.
func TestGetBenzingaRatingsFirstRating(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/benzinga/v1/ratings": benzingaRatingsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetBenzingaRatings(BenzingaRatingsParams{Ticker: "AAPL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rating := result.Results[0]

	if rating.BenzingaID != "bz-rating-001" {
		t.Errorf("expected benzinga_id bz-rating-001, got %s", rating.BenzingaID)
	}

	if rating.Ticker != "AAPL" {
		t.Errorf("expected ticker AAPL, got %s", rating.Ticker)
	}

	if rating.CompanyName != "Apple Inc." {
		t.Errorf("expected company_name Apple Inc., got %s", rating.CompanyName)
	}

	if rating.Analyst != "John Analyst" {
		t.Errorf("expected analyst John Analyst, got %s", rating.Analyst)
	}

	if rating.Firm != "Goldman Sachs" {
		t.Errorf("expected firm Goldman Sachs, got %s", rating.Firm)
	}

	if rating.Rating != "Buy" {
		t.Errorf("expected rating Buy, got %s", rating.Rating)
	}

	if rating.RatingAction != "upgrades" {
		t.Errorf("expected rating_action upgrades, got %s", rating.RatingAction)
	}

	if rating.PreviousRating != "Hold" {
		t.Errorf("expected previous_rating Hold, got %s", rating.PreviousRating)
	}

	if rating.PriceTarget != 250.00 {
		t.Errorf("expected price_target 250.00, got %f", rating.PriceTarget)
	}

	if rating.PriceTargetAction != "raises" {
		t.Errorf("expected price_target_action raises, got %s", rating.PriceTargetAction)
	}

	if rating.PreviousPriceTarget != 200.00 {
		t.Errorf("expected previous_price_target 200.00, got %f", rating.PreviousPriceTarget)
	}

	if rating.Importance != 4 {
		t.Errorf("expected importance 4, got %d", rating.Importance)
	}

	if rating.Currency != "USD" {
		t.Errorf("expected currency USD, got %s", rating.Currency)
	}
}

// TestGetBenzingaRatingsSecondRating verifies that the second rating in the
// response is correctly parsed with its own distinct values.
func TestGetBenzingaRatingsSecondRating(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/benzinga/v1/ratings": benzingaRatingsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetBenzingaRatings(BenzingaRatingsParams{Ticker: "MSFT"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rating := result.Results[1]

	if rating.Ticker != "MSFT" {
		t.Errorf("expected ticker MSFT, got %s", rating.Ticker)
	}

	if rating.Firm != "Morgan Stanley" {
		t.Errorf("expected firm Morgan Stanley, got %s", rating.Firm)
	}

	if rating.RatingAction != "maintains" {
		t.Errorf("expected rating_action maintains, got %s", rating.RatingAction)
	}

	if rating.PriceTarget != 480.00 {
		t.Errorf("expected price_target 480.00, got %f", rating.PriceTarget)
	}
}

// TestGetBenzingaRatingsRequestPath verifies that GetBenzingaRatings sends
// requests to the correct API path.
func TestGetBenzingaRatingsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(benzingaRatingsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetBenzingaRatings(BenzingaRatingsParams{Ticker: "AAPL"})

	if receivedPath != "/benzinga/v1/ratings" {
		t.Errorf("expected path /benzinga/v1/ratings, got %s", receivedPath)
	}
}

// TestGetBenzingaRatingsQueryParams verifies that all filter parameters are
// correctly sent to the Benzinga ratings API endpoint as query parameters.
func TestGetBenzingaRatingsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("ticker") != "AAPL" {
			t.Errorf("expected ticker=AAPL, got %s", q.Get("ticker"))
		}
		if q.Get("date.gte") != "2026-01-01" {
			t.Errorf("expected date.gte=2026-01-01, got %s", q.Get("date.gte"))
		}
		if q.Get("date.lte") != "2026-02-15" {
			t.Errorf("expected date.lte=2026-02-15, got %s", q.Get("date.lte"))
		}
		if q.Get("rating_action") != "upgrades" {
			t.Errorf("expected rating_action=upgrades, got %s", q.Get("rating_action"))
		}
		if q.Get("importance") != "3" {
			t.Errorf("expected importance=3, got %s", q.Get("importance"))
		}
		if q.Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(benzingaRatingsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetBenzingaRatings(BenzingaRatingsParams{
		Ticker:       "AAPL",
		DateGte:      "2026-01-01",
		DateLte:      "2026-02-15",
		RatingAction: "upgrades",
		Importance:   "3",
		Limit:        "50",
	})
}

// TestGetBenzingaRatingsAPIError verifies that GetBenzingaRatings returns
// an error when the API responds with a non-200 status code.
func TestGetBenzingaRatingsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"NOT_AUTHORIZED","message":"You are not entitled to this data."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetBenzingaRatings(BenzingaRatingsParams{Ticker: "AAPL"})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetBenzingaRatingsEmptyResults verifies that GetBenzingaRatings handles
// an empty results array without error.
func TestGetBenzingaRatingsEmptyResults(t *testing.T) {
	emptyJSON := `{"results":[],"status":"OK","request_id":"abc","count":0}`
	server := mockServer(t, map[string]string{
		"/benzinga/v1/ratings": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetBenzingaRatings(BenzingaRatingsParams{Ticker: "ZZZZZZ"})
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

// TestGetBenzingaEarnings verifies that GetBenzingaEarnings correctly parses
// the API response and returns the expected earnings data.
func TestGetBenzingaEarnings(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/benzinga/v1/earnings": benzingaEarningsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := BenzingaEarningsParams{
		Ticker: "AAPL",
		Limit:  "2",
	}

	result, err := client.GetBenzingaEarnings(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Count != 2 {
		t.Errorf("expected count 2, got %d", result.Count)
	}

	if result.RequestID != "bzearn123" {
		t.Errorf("expected request_id bzearn123, got %s", result.RequestID)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 earnings, got %d", len(result.Results))
	}
}

// TestGetBenzingaEarningsFirstRecord verifies that the first earnings record
// in the response is correctly parsed with all EPS, revenue, and fiscal fields.
func TestGetBenzingaEarningsFirstRecord(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/benzinga/v1/earnings": benzingaEarningsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetBenzingaEarnings(BenzingaEarningsParams{Ticker: "AAPL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	earn := result.Results[0]

	if earn.BenzingaID != "bz-earn-001" {
		t.Errorf("expected benzinga_id bz-earn-001, got %s", earn.BenzingaID)
	}

	if earn.Ticker != "AAPL" {
		t.Errorf("expected ticker AAPL, got %s", earn.Ticker)
	}

	if earn.CompanyName != "Apple Inc." {
		t.Errorf("expected company_name Apple Inc., got %s", earn.CompanyName)
	}

	if earn.DateStatus != "confirmed" {
		t.Errorf("expected date_status confirmed, got %s", earn.DateStatus)
	}

	if earn.ActualEPS != 2.18 {
		t.Errorf("expected actual_eps 2.18, got %f", earn.ActualEPS)
	}

	if earn.EstimatedEPS != 2.10 {
		t.Errorf("expected estimated_eps 2.10, got %f", earn.EstimatedEPS)
	}

	if earn.EPSSurprise != 0.08 {
		t.Errorf("expected eps_surprise 0.08, got %f", earn.EPSSurprise)
	}

	if earn.EPSSurprisePercent != 3.81 {
		t.Errorf("expected eps_surprise_percent 3.81, got %f", earn.EPSSurprisePercent)
	}

	if earn.ActualRevenue != 124500000000 {
		t.Errorf("expected actual_revenue 124500000000, got %f", earn.ActualRevenue)
	}

	if earn.EstimatedRevenue != 121000000000 {
		t.Errorf("expected estimated_revenue 121000000000, got %f", earn.EstimatedRevenue)
	}

	if earn.RevenueSurprise != 3500000000 {
		t.Errorf("expected revenue_surprise 3500000000, got %f", earn.RevenueSurprise)
	}

	if earn.FiscalPeriod != "Q1" {
		t.Errorf("expected fiscal_period Q1, got %s", earn.FiscalPeriod)
	}

	if earn.FiscalYear != 2026 {
		t.Errorf("expected fiscal_year 2026, got %d", earn.FiscalYear)
	}

	if earn.Importance != 5 {
		t.Errorf("expected importance 5, got %d", earn.Importance)
	}

	if earn.EPSMethod != "gaap" {
		t.Errorf("expected eps_method gaap, got %s", earn.EPSMethod)
	}

	if earn.Notes != "Record quarterly revenue" {
		t.Errorf("expected notes 'Record quarterly revenue', got %s", earn.Notes)
	}
}

// TestGetBenzingaEarningsSecondRecord verifies that the second earnings
// record in the response is correctly parsed as a projected earnings entry.
func TestGetBenzingaEarningsSecondRecord(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/benzinga/v1/earnings": benzingaEarningsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetBenzingaEarnings(BenzingaEarningsParams{Ticker: "MSFT"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	earn := result.Results[1]

	if earn.Ticker != "MSFT" {
		t.Errorf("expected ticker MSFT, got %s", earn.Ticker)
	}

	if earn.DateStatus != "projected" {
		t.Errorf("expected date_status projected, got %s", earn.DateStatus)
	}

	if earn.ActualEPS != 0 {
		t.Errorf("expected actual_eps 0, got %f", earn.ActualEPS)
	}

	if earn.EstimatedEPS != 3.25 {
		t.Errorf("expected estimated_eps 3.25, got %f", earn.EstimatedEPS)
	}

	if earn.FiscalPeriod != "Q2" {
		t.Errorf("expected fiscal_period Q2, got %s", earn.FiscalPeriod)
	}
}

// TestGetBenzingaEarningsRequestPath verifies that GetBenzingaEarnings sends
// requests to the correct API path.
func TestGetBenzingaEarningsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(benzingaEarningsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetBenzingaEarnings(BenzingaEarningsParams{Ticker: "AAPL"})

	if receivedPath != "/benzinga/v1/earnings" {
		t.Errorf("expected path /benzinga/v1/earnings, got %s", receivedPath)
	}
}

// TestGetBenzingaEarningsQueryParams verifies that all filter parameters are
// correctly sent to the Benzinga earnings API endpoint as query parameters.
func TestGetBenzingaEarningsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("ticker") != "AAPL" {
			t.Errorf("expected ticker=AAPL, got %s", q.Get("ticker"))
		}
		if q.Get("date.gte") != "2026-01-01" {
			t.Errorf("expected date.gte=2026-01-01, got %s", q.Get("date.gte"))
		}
		if q.Get("date.lte") != "2026-02-15" {
			t.Errorf("expected date.lte=2026-02-15, got %s", q.Get("date.lte"))
		}
		if q.Get("date_status") != "confirmed" {
			t.Errorf("expected date_status=confirmed, got %s", q.Get("date_status"))
		}
		if q.Get("fiscal_period") != "Q1" {
			t.Errorf("expected fiscal_period=Q1, got %s", q.Get("fiscal_period"))
		}
		if q.Get("fiscal_year") != "2026" {
			t.Errorf("expected fiscal_year=2026, got %s", q.Get("fiscal_year"))
		}
		if q.Get("importance") != "5" {
			t.Errorf("expected importance=5, got %s", q.Get("importance"))
		}
		if q.Get("limit") != "25" {
			t.Errorf("expected limit=25, got %s", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(benzingaEarningsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetBenzingaEarnings(BenzingaEarningsParams{
		Ticker:       "AAPL",
		DateGte:      "2026-01-01",
		DateLte:      "2026-02-15",
		DateStatus:   "confirmed",
		FiscalPeriod: "Q1",
		FiscalYear:   "2026",
		Importance:   "5",
		Limit:        "25",
	})
}

// TestGetBenzingaEarningsAPIError verifies that GetBenzingaEarnings returns
// an error when the API responds with a non-200 status code.
func TestGetBenzingaEarningsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"NOT_AUTHORIZED","message":"You are not entitled to this data."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetBenzingaEarnings(BenzingaEarningsParams{Ticker: "AAPL"})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetBenzingaEarningsEmptyResults verifies that GetBenzingaEarnings
// handles an empty results array without error.
func TestGetBenzingaEarningsEmptyResults(t *testing.T) {
	emptyJSON := `{"results":[],"status":"OK","request_id":"abc","count":0}`
	server := mockServer(t, map[string]string{
		"/benzinga/v1/earnings": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetBenzingaEarnings(BenzingaEarningsParams{Ticker: "ZZZZZZ"})
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

// TestGetBenzingaGuidance verifies that GetBenzingaGuidance correctly parses
// the API response and returns the expected corporate guidance data.
func TestGetBenzingaGuidance(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/benzinga/v1/guidance": benzingaGuidanceJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := BenzingaGuidanceParams{
		Ticker: "AAPL",
		Limit:  "1",
	}

	result, err := client.GetBenzingaGuidance(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.Count != 1 {
		t.Errorf("expected count 1, got %d", result.Count)
	}

	if result.RequestID != "bzguide123" {
		t.Errorf("expected request_id bzguide123, got %s", result.RequestID)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 guidance record, got %d", len(result.Results))
	}
}

// TestGetBenzingaGuidanceFirstRecord verifies that the first guidance record
// in the response is correctly parsed with all EPS and revenue guidance fields.
func TestGetBenzingaGuidanceFirstRecord(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/benzinga/v1/guidance": benzingaGuidanceJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetBenzingaGuidance(BenzingaGuidanceParams{Ticker: "AAPL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	guide := result.Results[0]

	if guide.BenzingaID != "bz-guide-001" {
		t.Errorf("expected benzinga_id bz-guide-001, got %s", guide.BenzingaID)
	}

	if guide.Ticker != "AAPL" {
		t.Errorf("expected ticker AAPL, got %s", guide.Ticker)
	}

	if guide.CompanyName != "Apple Inc." {
		t.Errorf("expected company_name Apple Inc., got %s", guide.CompanyName)
	}

	if guide.Positioning != "primary" {
		t.Errorf("expected positioning primary, got %s", guide.Positioning)
	}

	if guide.EstimatedEPSGuidance != 2.25 {
		t.Errorf("expected estimated_eps_guidance 2.25, got %f", guide.EstimatedEPSGuidance)
	}

	if guide.MinEPSGuidance != 2.15 {
		t.Errorf("expected min_eps_guidance 2.15, got %f", guide.MinEPSGuidance)
	}

	if guide.MaxEPSGuidance != 2.35 {
		t.Errorf("expected max_eps_guidance 2.35, got %f", guide.MaxEPSGuidance)
	}

	if guide.EstimatedRevenueGuidance != 128000000000 {
		t.Errorf("expected estimated_revenue_guidance 128000000000, got %f", guide.EstimatedRevenueGuidance)
	}

	if guide.MinRevenueGuidance != 125000000000 {
		t.Errorf("expected min_revenue_guidance 125000000000, got %f", guide.MinRevenueGuidance)
	}

	if guide.MaxRevenueGuidance != 131000000000 {
		t.Errorf("expected max_revenue_guidance 131000000000, got %f", guide.MaxRevenueGuidance)
	}

	if guide.FiscalPeriod != "Q2" {
		t.Errorf("expected fiscal_period Q2, got %s", guide.FiscalPeriod)
	}

	if guide.FiscalYear != 2026 {
		t.Errorf("expected fiscal_year 2026, got %d", guide.FiscalYear)
	}

	if guide.ReleaseType != "official" {
		t.Errorf("expected release_type official, got %s", guide.ReleaseType)
	}

	if guide.Notes != "Management raised guidance" {
		t.Errorf("expected notes 'Management raised guidance', got %s", guide.Notes)
	}
}

// TestGetBenzingaGuidanceRequestPath verifies that GetBenzingaGuidance sends
// requests to the correct API path.
func TestGetBenzingaGuidanceRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(benzingaGuidanceJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetBenzingaGuidance(BenzingaGuidanceParams{Ticker: "AAPL"})

	if receivedPath != "/benzinga/v1/guidance" {
		t.Errorf("expected path /benzinga/v1/guidance, got %s", receivedPath)
	}
}

// TestGetBenzingaGuidanceQueryParams verifies that all filter parameters are
// correctly sent to the Benzinga guidance API endpoint as query parameters.
func TestGetBenzingaGuidanceQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("ticker") != "AAPL" {
			t.Errorf("expected ticker=AAPL, got %s", q.Get("ticker"))
		}
		if q.Get("date.gte") != "2026-01-01" {
			t.Errorf("expected date.gte=2026-01-01, got %s", q.Get("date.gte"))
		}
		if q.Get("positioning") != "primary" {
			t.Errorf("expected positioning=primary, got %s", q.Get("positioning"))
		}
		if q.Get("fiscal_period") != "Q2" {
			t.Errorf("expected fiscal_period=Q2, got %s", q.Get("fiscal_period"))
		}
		if q.Get("limit") != "10" {
			t.Errorf("expected limit=10, got %s", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(benzingaGuidanceJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetBenzingaGuidance(BenzingaGuidanceParams{
		Ticker:       "AAPL",
		DateGte:      "2026-01-01",
		Positioning:  "primary",
		FiscalPeriod: "Q2",
		Limit:        "10",
	})
}

// TestGetBenzingaGuidanceAPIError verifies that GetBenzingaGuidance returns
// an error when the API responds with a non-200 status code.
func TestGetBenzingaGuidanceAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"NOT_AUTHORIZED","message":"You are not entitled to this data."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetBenzingaGuidance(BenzingaGuidanceParams{Ticker: "AAPL"})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// TestGetBenzingaAnalysts verifies that GetBenzingaAnalysts correctly parses
// the API response and returns the expected analyst details data.
func TestGetBenzingaAnalysts(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/benzinga/v1/analysts": benzingaAnalystsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := BenzingaAnalystsParams{
		Limit: "2",
	}

	result, err := client.GetBenzingaAnalysts(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "bzanalysts123" {
		t.Errorf("expected request_id bzanalysts123, got %s", result.RequestID)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 analysts, got %d", len(result.Results))
	}
}

// TestGetBenzingaAnalystsFirstAnalyst verifies that the first analyst in the
// response is correctly parsed with all performance metric fields.
func TestGetBenzingaAnalystsFirstAnalyst(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/benzinga/v1/analysts": benzingaAnalystsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetBenzingaAnalysts(BenzingaAnalystsParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	analyst := result.Results[0]

	if analyst.BenzingaID != "analyst-001" {
		t.Errorf("expected benzinga_id analyst-001, got %s", analyst.BenzingaID)
	}

	if analyst.FullName != "John Analyst" {
		t.Errorf("expected full_name John Analyst, got %s", analyst.FullName)
	}

	if analyst.FirmName != "Goldman Sachs" {
		t.Errorf("expected firm_name Goldman Sachs, got %s", analyst.FirmName)
	}

	if analyst.SmartScore != 85.5 {
		t.Errorf("expected smart_score 85.5, got %f", analyst.SmartScore)
	}

	if analyst.OverallSuccessRate != 0.72 {
		t.Errorf("expected overall_success_rate 0.72, got %f", analyst.OverallSuccessRate)
	}

	if analyst.OverallAvgReturn != 12.5 {
		t.Errorf("expected overall_avg_return 12.5, got %f", analyst.OverallAvgReturn)
	}

	if analyst.TotalRatings != 350 {
		t.Errorf("expected total_ratings 350, got %f", analyst.TotalRatings)
	}
}

// TestGetBenzingaAnalystsSecondAnalyst verifies that the second analyst in
// the response is correctly parsed with its own distinct values.
func TestGetBenzingaAnalystsSecondAnalyst(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/benzinga/v1/analysts": benzingaAnalystsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetBenzingaAnalysts(BenzingaAnalystsParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	analyst := result.Results[1]

	if analyst.FullName != "Sarah Analyst" {
		t.Errorf("expected full_name Sarah Analyst, got %s", analyst.FullName)
	}

	if analyst.FirmName != "Morgan Stanley" {
		t.Errorf("expected firm_name Morgan Stanley, got %s", analyst.FirmName)
	}

	if analyst.SmartScore != 78.2 {
		t.Errorf("expected smart_score 78.2, got %f", analyst.SmartScore)
	}
}

// TestGetBenzingaAnalystsRequestPath verifies that GetBenzingaAnalysts sends
// requests to the correct API path.
func TestGetBenzingaAnalystsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(benzingaAnalystsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetBenzingaAnalysts(BenzingaAnalystsParams{})

	if receivedPath != "/benzinga/v1/analysts" {
		t.Errorf("expected path /benzinga/v1/analysts, got %s", receivedPath)
	}
}

// TestGetBenzingaAnalystsQueryParams verifies that all filter parameters are
// correctly sent to the Benzinga analysts API endpoint as query parameters.
func TestGetBenzingaAnalystsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("benzinga_id") != "analyst-001" {
			t.Errorf("expected benzinga_id=analyst-001, got %s", q.Get("benzinga_id"))
		}
		if q.Get("firm_name") != "Goldman Sachs" {
			t.Errorf("expected firm_name=Goldman Sachs, got %s", q.Get("firm_name"))
		}
		if q.Get("limit") != "10" {
			t.Errorf("expected limit=10, got %s", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(benzingaAnalystsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetBenzingaAnalysts(BenzingaAnalystsParams{
		BenzingaID: "analyst-001",
		FirmName:   "Goldman Sachs",
		Limit:      "10",
	})
}

// TestGetBenzingaAnalystsAPIError verifies that GetBenzingaAnalysts returns
// an error when the API responds with a non-200 status code.
func TestGetBenzingaAnalystsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"NOT_AUTHORIZED","message":"You are not entitled to this data."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetBenzingaAnalysts(BenzingaAnalystsParams{})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}
