//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

// IndicesSnapshotSession represents the trading session data for an index
// snapshot, including open, high, low, close values and the calculated
// change and change percent from the previous close.
type IndicesSnapshotSession struct {
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"change_percent"`
	Close         float64 `json:"close"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	Open          float64 `json:"open"`
	PreviousClose float64 `json:"previous_close"`
}

// IndicesSnapshotResult represents a single index snapshot entry returned
// by the Massive API. It contains the index ticker, name, current value,
// market status, last update timestamp, asset type, data timeframe, and
// the session data with OHLC and change metrics.
type IndicesSnapshotResult struct {
	Ticker       string                 `json:"ticker"`
	Name         string                 `json:"name"`
	Value        float64                `json:"value"`
	Type         string                 `json:"type"`
	Timeframe    string                 `json:"timeframe"`
	MarketStatus string                 `json:"market_status"`
	LastUpdated  int64                  `json:"last_updated"`
	Session      IndicesSnapshotSession `json:"session"`
	Error        string                 `json:"error,omitempty"`
	Message      string                 `json:"message,omitempty"`
}

// IndicesSnapshotResponse represents the API response for an indices
// snapshot request from the /v3/snapshot/indices endpoint. It contains
// the response status, request ID, an optional pagination URL, and an
// array of index snapshot results.
type IndicesSnapshotResponse struct {
	Status    string                  `json:"status"`
	RequestID string                  `json:"request_id"`
	NextURL   string                  `json:"next_url,omitempty"`
	Results   []IndicesSnapshotResult `json:"results"`
}

// IndicesSnapshotParams holds the optional query parameters for fetching
// index snapshots. TickerAnyOf accepts a comma-separated list of tickers
// (up to 250). The ticker range filters allow lexicographic filtering.
// Order and Sort control result ordering, and Limit sets the maximum
// number of results returned (default 10, max 250).
type IndicesSnapshotParams struct {
	TickerAnyOf string
	Ticker      string
	TickerGte   string
	TickerGt    string
	TickerLte   string
	TickerLt    string
	Order       string
	Limit       string
	Sort        string
}

// GetIndicesSnapshot retrieves snapshot data for one or more indices from
// the /v3/snapshot/indices endpoint. The response includes each index's
// current value, trading session metrics (open, high, low, close, change),
// market status, and last update timestamp. Results can be filtered by
// ticker symbols and paginated using limit and sort parameters.
func (c *Client) GetIndicesSnapshot(p IndicesSnapshotParams) (*IndicesSnapshotResponse, error) {
	path := "/v3/snapshot/indices"

	params := map[string]string{
		"ticker.any_of": p.TickerAnyOf,
		"ticker":        p.Ticker,
		"ticker.gte":    p.TickerGte,
		"ticker.gt":     p.TickerGt,
		"ticker.lte":    p.TickerLte,
		"ticker.lt":     p.TickerLt,
		"order":         p.Order,
		"limit":         p.Limit,
		"sort":          p.Sort,
	}

	var result IndicesSnapshotResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
