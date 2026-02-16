//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

// IndicesTicker represents a single index ticker's reference data including
// the ticker symbol, name, market, locale, active status, and the source
// data feed that provides the index data.
type IndicesTicker struct {
	Ticker     string `json:"ticker"`
	Name       string `json:"name"`
	Market     string `json:"market"`
	Locale     string `json:"locale"`
	Active     bool   `json:"active"`
	SourceFeed string `json:"source_feed"`
}

// IndicesTickersResponse represents the API response for listing index
// reference ticker data with pagination support via NextURL.
type IndicesTickersResponse struct {
	Status    string          `json:"status"`
	Count     int             `json:"count"`
	RequestID string          `json:"request_id"`
	NextURL   string          `json:"next_url"`
	Results   []IndicesTicker `json:"results"`
}

// IndicesTickerParams holds the query parameters for searching and filtering
// index tickers from the reference endpoint. The Market field is automatically
// set to "indices" by the GetIndicesTickers method.
type IndicesTickerParams struct {
	Ticker string
	Search string
	Active string
	Sort   string
	Order  string
	Limit  string
}

// GetIndicesTickers retrieves a list of index tickers matching the filter
// criteria specified in the IndicesTickerParams. It uses the same reference
// tickers endpoint as stocks but forces the market parameter to "indices"
// so that only index tickers are returned.
func (c *Client) GetIndicesTickers(p IndicesTickerParams) (*IndicesTickersResponse, error) {
	path := "/v3/reference/tickers"

	params := map[string]string{
		"market": "indices",
		"ticker": p.Ticker,
		"search": p.Search,
		"active": p.Active,
		"sort":   p.Sort,
		"order":  p.Order,
		"limit":  p.Limit,
	}

	var result IndicesTickersResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
