//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

import (
	"fmt"
)

// SnapshotBar represents OHLC bar data with volume and volume-weighted
// average price. Used in day, prevDay, and min snapshot sections.
type SnapshotBar struct {
	Open   float64 `json:"o"`
	High   float64 `json:"h"`
	Low    float64 `json:"l"`
	Close  float64 `json:"c"`
	Volume float64 `json:"v"`
	VWAP   float64 `json:"vw"`
}

// SnapshotMinBar represents the most recent minute bar with additional
// fields for accumulated volume, timestamp, and number of transactions.
type SnapshotMinBar struct {
	Open              float64 `json:"o"`
	High              float64 `json:"h"`
	Low               float64 `json:"l"`
	Close             float64 `json:"c"`
	Volume            float64 `json:"v"`
	VWAP              float64 `json:"vw"`
	Timestamp         int64   `json:"t"`
	NumTransactions   int     `json:"n"`
	AccumulatedVolume float64 `json:"av"`
}

// SnapshotTicker represents a single ticker's snapshot data containing
// the current day's bar, previous day's bar, latest minute bar, the
// calculated change values, and the last update timestamp.
type SnapshotTicker struct {
	Ticker          string      `json:"ticker"`
	TodaysChange    float64     `json:"todaysChange"`
	TodaysChangePct float64     `json:"todaysChangePerc"`
	Updated         int64       `json:"updated"`
	Day             SnapshotBar `json:"day"`
	PrevDay         SnapshotBar `json:"prevDay"`
	Min             SnapshotMinBar `json:"min"`
}

// SingleTickerSnapshotResponse represents the API response for a single
// ticker snapshot request from the /v2/snapshot endpoint.
type SingleTickerSnapshotResponse struct {
	Status    string         `json:"status"`
	RequestID string         `json:"request_id"`
	Ticker    SnapshotTicker `json:"ticker"`
}

// AllTickersSnapshotResponse represents the API response for a full
// market or filtered multi-ticker snapshot request.
type AllTickersSnapshotResponse struct {
	Status    string           `json:"status"`
	RequestID string           `json:"request_id"`
	Count     int              `json:"count"`
	Tickers   []SnapshotTicker `json:"tickers"`
}

// GainersLosersSnapshotResponse represents the API response for the
// top market movers (gainers or losers) snapshot request.
type GainersLosersSnapshotResponse struct {
	Status    string           `json:"status"`
	RequestID string           `json:"request_id"`
	Tickers   []SnapshotTicker `json:"tickers"`
}

// AllTickersSnapshotParams holds the optional query parameters for
// fetching a full market or filtered multi-ticker snapshot.
type AllTickersSnapshotParams struct {
	Tickers    string
	IncludeOTC string
}

// GainersLosersParams holds the optional query parameters for fetching
// the top market movers (gainers or losers) snapshot.
type GainersLosersParams struct {
	IncludeOTC string
}

// GetSnapshotTicker retrieves the most recent snapshot for a single
// stock ticker, including the current day's bar, previous day's bar,
// latest minute bar, and the day's price change values.
func (c *Client) GetSnapshotTicker(ticker string) (*SingleTickerSnapshotResponse, error) {
	path := fmt.Sprintf("/v2/snapshot/locale/us/markets/stocks/tickers/%s", ticker)

	var result SingleTickerSnapshotResponse
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSnapshotAllTickers retrieves snapshot data for all US stock tickers
// or a filtered subset specified by a comma-separated list in the params.
// The response includes day, previous day, and minute bars for each ticker.
func (c *Client) GetSnapshotAllTickers(p AllTickersSnapshotParams) (*AllTickersSnapshotResponse, error) {
	path := "/v2/snapshot/locale/us/markets/stocks/tickers"

	params := map[string]string{
		"tickers":     p.Tickers,
		"include_otc": p.IncludeOTC,
	}

	var result AllTickersSnapshotResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSnapshotGainersLosers retrieves the current top 20 gainers or losers
// in the US stock market. The direction parameter must be either "gainers"
// or "losers" and determines which set of movers is returned.
func (c *Client) GetSnapshotGainersLosers(direction string, p GainersLosersParams) (*GainersLosersSnapshotResponse, error) {
	path := fmt.Sprintf("/v2/snapshot/locale/us/markets/stocks/%s", direction)

	params := map[string]string{
		"include_otc": p.IncludeOTC,
	}

	var result GainersLosersSnapshotResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
