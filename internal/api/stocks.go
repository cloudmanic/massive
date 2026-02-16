//
// Date: 2026-02-14
// Copyright (c) 2026. All rights reserved.
//

package api

import (
	"fmt"
)

// OpenCloseResponse represents the API response for daily open/close data
// for a specific stock ticker on a given date.
type OpenCloseResponse struct {
	Status     string  `json:"status"`
	Symbol     string  `json:"symbol"`
	From       string  `json:"from"`
	Open       float64 `json:"open"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Close      float64 `json:"close"`
	Volume     int64   `json:"volume"`
	AfterHours float64 `json:"afterHours"`
	PreMarket  float64 `json:"preMarket"`
}

// BarsResponse represents the API response for OHLC aggregate bar data
// over a custom time range for a specific stock ticker.
type BarsResponse struct {
	Status       string `json:"status"`
	Ticker       string `json:"ticker"`
	Adjusted     bool   `json:"adjusted"`
	QueryCount   int    `json:"queryCount"`
	ResultsCount int    `json:"resultsCount"`
	RequestID    string `json:"request_id"`
	Results      []Bar  `json:"results"`
}

// Bar represents a single OHLC bar with volume and trade data.
// Field names match the abbreviated JSON keys from the API.
type Bar struct {
	Open      float64 `json:"o"`
	High      float64 `json:"h"`
	Low       float64 `json:"l"`
	Close     float64 `json:"c"`
	Volume    float64 `json:"v"`
	VWAP      float64 `json:"vw"`
	Timestamp int64   `json:"t"`
	NumTrades int     `json:"n"`
}

// MarketSummaryResponse represents the API response for a daily grouped
// market summary of all US stocks on a given date.
type MarketSummaryResponse struct {
	Status       string          `json:"status"`
	Adjusted     bool            `json:"adjusted"`
	QueryCount   int             `json:"queryCount"`
	ResultsCount int             `json:"resultsCount"`
	RequestID    string          `json:"request_id"`
	Results      []MarketSummary `json:"results"`
}

// MarketSummary represents a single ticker's daily summary within a
// grouped market response. The Ticker field uses the abbreviated "T" key.
type MarketSummary struct {
	Ticker    string  `json:"T"`
	Open      float64 `json:"o"`
	High      float64 `json:"h"`
	Low       float64 `json:"l"`
	Close     float64 `json:"c"`
	Volume    float64 `json:"v"`
	VWAP      float64 `json:"vw"`
	Timestamp int64   `json:"t"`
	NumTrades int     `json:"n"`
	OTC       bool    `json:"otc"`
}

// TickersResponse represents the API response for listing reference
// ticker data with pagination support via NextURL.
type TickersResponse struct {
	Status    string   `json:"status"`
	Count     int      `json:"count"`
	RequestID string   `json:"request_id"`
	NextURL   string   `json:"next_url"`
	Results   []Ticker `json:"results"`
}

// Ticker represents a single stock ticker's reference data including
// exchange information, identifiers, and active status.
type Ticker struct {
	Ticker          string `json:"ticker"`
	Name            string `json:"name"`
	Market          string `json:"market"`
	Locale          string `json:"locale"`
	PrimaryExchange string `json:"primary_exchange"`
	Type            string `json:"type"`
	Active          bool   `json:"active"`
	CurrencyName    string `json:"currency_name"`
	CIK             string `json:"cik"`
	CompositeFIGI   string `json:"composite_figi"`
	ShareClassFIGI  string `json:"share_class_figi"`
	LastUpdatedUTC  string `json:"last_updated_utc"`
}

// BarsParams holds the query parameters for fetching custom OHLC bar data
// from the aggregates endpoint.
type BarsParams struct {
	Multiplier string
	Timespan   string
	From       string
	To         string
	Adjusted   string
	Sort       string
	Limit      string
}

// TickerParams holds the query parameters for searching and filtering
// stock tickers from the reference endpoint.
type TickerParams struct {
	Ticker   string
	Type     string
	Market   string
	Exchange string
	Search   string
	Active   string
	Sort     string
	Order    string
	Limit    string
}

// MarketSummaryParams holds the query parameters for fetching a daily
// grouped market summary.
type MarketSummaryParams struct {
	Adjusted   string
	IncludeOTC string
}

// GetOpenClose retrieves the daily open, close, high, low, volume, and
// extended hours prices for a specific stock ticker on a given date.
func (c *Client) GetOpenClose(ticker, date string, adjusted string) (*OpenCloseResponse, error) {
	path := fmt.Sprintf("/v1/open-close/%s/%s", ticker, date)
	params := map[string]string{
		"adjusted": adjusted,
	}

	var result OpenCloseResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBars retrieves custom OHLC aggregate bar data for a specific stock
// ticker over the time range specified in the BarsParams.
func (c *Client) GetBars(ticker string, p BarsParams) (*BarsResponse, error) {
	path := fmt.Sprintf("/v2/aggs/ticker/%s/range/%s/%s/%s/%s",
		ticker, p.Multiplier, p.Timespan, p.From, p.To)

	params := map[string]string{
		"adjusted": p.Adjusted,
		"sort":     p.Sort,
		"limit":    p.Limit,
	}

	var result BarsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetMarketSummary retrieves the grouped daily OHLC summary for all US
// stocks on the specified date, with optional OTC inclusion.
func (c *Client) GetMarketSummary(date string, p MarketSummaryParams) (*MarketSummaryResponse, error) {
	path := fmt.Sprintf("/v2/aggs/grouped/locale/us/market/stocks/%s", date)

	params := map[string]string{
		"adjusted":    p.Adjusted,
		"include_otc": p.IncludeOTC,
	}

	var result MarketSummaryResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetTickers retrieves a list of stock tickers matching the filter
// criteria specified in the TickerParams.
func (c *Client) GetTickers(p TickerParams) (*TickersResponse, error) {
	path := "/v3/reference/tickers"

	params := map[string]string{
		"ticker":   p.Ticker,
		"type":     p.Type,
		"market":   p.Market,
		"exchange": p.Exchange,
		"search":   p.Search,
		"active":   p.Active,
		"sort":     p.Sort,
		"order":    p.Order,
		"limit":    p.Limit,
	}

	var result TickersResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
