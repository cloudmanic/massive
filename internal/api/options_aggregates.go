//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

import (
	"fmt"
)

// OptionsBarsResponse represents the API response for OHLC aggregate bar
// data over a custom time range for a specific options contract ticker.
// Options bars include volume, VWAP, and trade count fields.
type OptionsBarsResponse struct {
	Status       string       `json:"status"`
	Ticker       string       `json:"ticker"`
	Adjusted     bool         `json:"adjusted"`
	QueryCount   int          `json:"queryCount"`
	ResultsCount int          `json:"resultsCount"`
	RequestID    string       `json:"request_id"`
	Count        int          `json:"count"`
	Results      []OptionsBar `json:"results"`
}

// OptionsBar represents a single OHLC bar for an options contract. Each bar
// contains open, high, low, close, volume, volume-weighted average price,
// a millisecond Unix timestamp, and the number of trades in the window.
type OptionsBar struct {
	Open      float64 `json:"o"`
	High      float64 `json:"h"`
	Low       float64 `json:"l"`
	Close     float64 `json:"c"`
	Volume    float64 `json:"v"`
	VWAP      float64 `json:"vw"`
	Timestamp int64   `json:"t"`
	NumTrades int     `json:"n"`
}

// OptionsDailyTickerSummaryResponse represents the API response for daily
// open/close data for a specific options contract on a given date, including
// pre-market and after-hours trade prices when available.
type OptionsDailyTickerSummaryResponse struct {
	Status     string  `json:"status"`
	Symbol     string  `json:"symbol"`
	From       string  `json:"from"`
	Open       float64 `json:"open"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Close      float64 `json:"close"`
	Volume     float64 `json:"volume"`
	AfterHours float64 `json:"afterHours"`
	PreMarket  float64 `json:"preMarket"`
}

// OptionsPreviousDayBarResponse represents the API response for the
// previous trading day's OHLC data for a specific options contract ticker.
type OptionsPreviousDayBarResponse struct {
	Status       string                  `json:"status"`
	Ticker       string                  `json:"ticker"`
	Adjusted     bool                    `json:"adjusted"`
	QueryCount   int                     `json:"queryCount"`
	ResultsCount int                     `json:"resultsCount"`
	RequestID    string                  `json:"request_id"`
	Count        int                     `json:"count"`
	Results      []OptionsPreviousDayBar `json:"results"`
}

// OptionsPreviousDayBar represents a single previous-day OHLC bar for an
// options contract. It includes the ticker symbol ("T" key) along with
// open, high, low, close, volume, VWAP, timestamp, and number of trades.
type OptionsPreviousDayBar struct {
	Ticker    string  `json:"T"`
	Open      float64 `json:"o"`
	High      float64 `json:"h"`
	Low       float64 `json:"l"`
	Close     float64 `json:"c"`
	Volume    float64 `json:"v"`
	VWAP      float64 `json:"vw"`
	Timestamp int64   `json:"t"`
	NumTrades int     `json:"n"`
}

// OptionsBarsParams holds the query parameters for fetching custom OHLC
// bar data from the options aggregates endpoint. The Multiplier and Timespan
// fields are used to build the URL path, while Adjusted, Sort, and Limit
// are sent as query parameters.
type OptionsBarsParams struct {
	Multiplier string
	Timespan   string
	From       string
	To         string
	Adjusted   string
	Sort       string
	Limit      string
}

// GetOptionsBars retrieves custom OHLC aggregate bar data for a specific
// options contract ticker over the time range specified in the OptionsBarsParams.
// The endpoint path includes the ticker, multiplier, timespan, from, and to
// values. Adjusted, sort, and limit are passed as query parameters.
func (c *Client) GetOptionsBars(ticker string, p OptionsBarsParams) (*OptionsBarsResponse, error) {
	path := fmt.Sprintf("/v2/aggs/ticker/%s/range/%s/%s/%s/%s",
		ticker, p.Multiplier, p.Timespan, p.From, p.To)

	params := map[string]string{
		"adjusted": p.Adjusted,
		"sort":     p.Sort,
		"limit":    p.Limit,
	}

	var result OptionsBarsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOptionsDailyTickerSummary retrieves the daily open, close, high, low,
// volume, and extended hours prices for a specific options contract on a
// given date. The adjusted parameter controls whether results are adjusted
// for splits.
func (c *Client) GetOptionsDailyTickerSummary(ticker, date string, adjusted string) (*OptionsDailyTickerSummaryResponse, error) {
	path := fmt.Sprintf("/v1/open-close/%s/%s", ticker, date)

	params := map[string]string{
		"adjusted": adjusted,
	}

	var result OptionsDailyTickerSummaryResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOptionsPreviousDayBar retrieves the previous trading day's open, high,
// low, close, volume, VWAP, and trade count for a specified options contract
// ticker. The adjusted parameter controls whether results are adjusted for
// splits.
func (c *Client) GetOptionsPreviousDayBar(ticker string, adjusted string) (*OptionsPreviousDayBarResponse, error) {
	path := fmt.Sprintf("/v2/aggs/ticker/%s/prev", ticker)

	params := map[string]string{
		"adjusted": adjusted,
	}

	var result OptionsPreviousDayBarResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
