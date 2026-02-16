//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

import (
	"fmt"
)

// IndicesBarsResponse represents the API response for OHLC aggregate bar
// data over a custom time range for a specific index ticker. Unlike stock
// bars, index bars do not include volume, VWAP, or trade count fields.
type IndicesBarsResponse struct {
	Status       string       `json:"status"`
	Ticker       string       `json:"ticker"`
	QueryCount   int          `json:"queryCount"`
	ResultsCount int          `json:"resultsCount"`
	RequestID    string       `json:"request_id"`
	Count        int          `json:"count"`
	Results      []IndicesBar `json:"results"`
}

// IndicesBar represents a single OHLC bar for an index. Index bars contain
// open, high, low, close, and a millisecond Unix timestamp but do not
// include volume, VWAP, or number of trades like stock bars do.
type IndicesBar struct {
	Open      float64 `json:"o"`
	High      float64 `json:"h"`
	Low       float64 `json:"l"`
	Close     float64 `json:"c"`
	Timestamp int64   `json:"t"`
}

// IndicesDailyTickerSummaryResponse represents the API response for daily
// open/close data for a specific index ticker on a given date, including
// pre-market and after-hours prices when available.
type IndicesDailyTickerSummaryResponse struct {
	Status     string  `json:"status"`
	Symbol     string  `json:"symbol"`
	From       string  `json:"from"`
	Open       float64 `json:"open"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Close      float64 `json:"close"`
	AfterHours float64 `json:"afterHours"`
	PreMarket  float64 `json:"preMarket"`
}

// IndicesPreviousDayBarResponse represents the API response for the
// previous trading day's OHLC data for a specific index ticker.
type IndicesPreviousDayBarResponse struct {
	Status       string                `json:"status"`
	Ticker       string                `json:"ticker"`
	QueryCount   int                   `json:"queryCount"`
	ResultsCount int                   `json:"resultsCount"`
	RequestID    string                `json:"request_id"`
	Count        int                   `json:"count"`
	Results      []IndicesPreviousDayBar `json:"results"`
}

// IndicesPreviousDayBar represents a single previous-day OHLC bar for an
// index. It includes the ticker symbol ("T" key) along with open, high,
// low, close, and a millisecond Unix timestamp.
type IndicesPreviousDayBar struct {
	Ticker    string  `json:"T"`
	Open      float64 `json:"o"`
	High      float64 `json:"h"`
	Low       float64 `json:"l"`
	Close     float64 `json:"c"`
	Timestamp int64   `json:"t"`
}

// IndicesBarsParams holds the query parameters for fetching custom OHLC
// bar data from the indices aggregates endpoint. The Multiplier and Timespan
// fields are used to build the URL path, while Sort and Limit are sent as
// query parameters.
type IndicesBarsParams struct {
	Multiplier string
	Timespan   string
	From       string
	To         string
	Sort       string
	Limit      string
}

// GetIndicesBars retrieves custom OHLC aggregate bar data for a specific
// index ticker over the time range specified in the IndicesBarsParams.
// The endpoint path includes the ticker, multiplier, timespan, from, and
// to values. Sort and limit are passed as query parameters.
func (c *Client) GetIndicesBars(ticker string, p IndicesBarsParams) (*IndicesBarsResponse, error) {
	path := fmt.Sprintf("/v2/aggs/ticker/%s/range/%s/%s/%s/%s",
		ticker, p.Multiplier, p.Timespan, p.From, p.To)

	params := map[string]string{
		"sort":  p.Sort,
		"limit": p.Limit,
	}

	var result IndicesBarsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetIndicesDailyTickerSummary retrieves the daily open, close, high, low,
// and extended hours prices for a specific index ticker on a given date.
// This mirrors the stocks open-close endpoint but for index tickers.
func (c *Client) GetIndicesDailyTickerSummary(ticker, date string) (*IndicesDailyTickerSummaryResponse, error) {
	path := fmt.Sprintf("/v1/open-close/%s/%s", ticker, date)

	var result IndicesDailyTickerSummaryResponse
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetIndicesPreviousDayBar retrieves the previous trading day's open, high,
// low, and close data for a specified index ticker. This is useful for
// quickly checking the most recent completed session's price data.
func (c *Client) GetIndicesPreviousDayBar(ticker string) (*IndicesPreviousDayBarResponse, error) {
	path := fmt.Sprintf("/v2/aggs/ticker/%s/prev", ticker)

	var result IndicesPreviousDayBarResponse
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
