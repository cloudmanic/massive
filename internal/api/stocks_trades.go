//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

import (
	"fmt"
)

// TradesResponse represents the API response for tick-level trade data
// returned by the /v3/trades/{stockTicker} endpoint. It includes pagination
// via NextURL and a slice of individual Trade records.
type TradesResponse struct {
	Status    string  `json:"status"`
	RequestID string  `json:"request_id"`
	NextURL   string  `json:"next_url"`
	Results   []Trade `json:"results"`
}

// Trade represents a single tick-level trade record for a stock.
// Fields use abbreviated JSON keys matching the Massive API response format.
type Trade struct {
	Conditions           []int   `json:"conditions"`
	Correction           int     `json:"correction"`
	Exchange             int     `json:"exchange"`
	ID                   string  `json:"id"`
	ParticipantTimestamp int64   `json:"participant_timestamp"`
	Price                float64 `json:"price"`
	SequenceNumber       int64   `json:"sequence_number"`
	SipTimestamp         int64   `json:"sip_timestamp"`
	Size                 float64 `json:"size"`
	Tape                 int     `json:"tape"`
	TrfID                int     `json:"trf_id"`
	TrfTimestamp         int64   `json:"trf_timestamp"`
}

// TradesParams holds the query parameters for fetching tick-level trade
// data from the /v3/trades endpoint.
type TradesParams struct {
	Timestamp    string
	TimestampGte string
	TimestampGt  string
	TimestampLte string
	TimestampLt  string
	Order        string
	Limit        string
	Sort         string
}

// LastTradeResponse represents the API response for the most recent trade
// of a stock ticker from the /v2/last/trade/{stocksTicker} endpoint.
type LastTradeResponse struct {
	Status    string    `json:"status"`
	RequestID string    `json:"request_id"`
	Results   LastTrade `json:"results"`
}

// LastTrade represents the most recent trade for a stock ticker.
// Fields use abbreviated single-character JSON keys from the API.
type LastTrade struct {
	Ticker               string  `json:"T"`
	Conditions           []int   `json:"c"`
	Correction           int     `json:"e"`
	TrfTimestamp         int64   `json:"f"`
	ID                   string  `json:"i"`
	Price                float64 `json:"p"`
	SequenceNumber       int64   `json:"q"`
	TrfID                int     `json:"r"`
	Size                 float64 `json:"s"`
	SipTimestamp         int64   `json:"t"`
	Exchange             int     `json:"x"`
	ParticipantTimestamp int64   `json:"y"`
	Tape                 int     `json:"z"`
}

// QuotesResponse represents the API response for tick-level NBBO quote data
// returned by the /v3/quotes/{stockTicker} endpoint. It includes pagination
// via NextURL and a slice of individual Quote records.
type QuotesResponse struct {
	Status    string  `json:"status"`
	RequestID string  `json:"request_id"`
	NextURL   string  `json:"next_url"`
	Results   []Quote `json:"results"`
}

// Quote represents a single NBBO quote record for a stock, containing
// bid/ask prices, sizes, exchange information, and timestamps.
type Quote struct {
	AskExchange          int     `json:"ask_exchange"`
	AskPrice             float64 `json:"ask_price"`
	AskSize              float64 `json:"ask_size"`
	BidExchange          int     `json:"bid_exchange"`
	BidPrice             float64 `json:"bid_price"`
	BidSize              float64 `json:"bid_size"`
	Conditions           []int   `json:"conditions"`
	Indicators           []int   `json:"indicators"`
	ParticipantTimestamp int64   `json:"participant_timestamp"`
	SequenceNumber       int64   `json:"sequence_number"`
	SipTimestamp         int64   `json:"sip_timestamp"`
	Tape                 int     `json:"tape"`
	TrfTimestamp         int64   `json:"trf_timestamp"`
}

// QuotesParams holds the query parameters for fetching tick-level NBBO
// quote data from the /v3/quotes endpoint.
type QuotesParams struct {
	Timestamp    string
	TimestampGte string
	TimestampGt  string
	TimestampLte string
	TimestampLt  string
	Order        string
	Limit        string
	Sort         string
}

// LastQuoteResponse represents the API response for the most recent NBBO
// quote of a stock ticker from the /v2/last/nbbo/{stocksTicker} endpoint.
type LastQuoteResponse struct {
	Status    string    `json:"status"`
	RequestID string    `json:"request_id"`
	Results   LastQuote `json:"results"`
}

// LastQuote represents the most recent NBBO quote for a stock ticker.
// Fields use abbreviated single-character JSON keys from the API where
// uppercase letters represent ask-side data and lowercase represent bid-side.
type LastQuote struct {
	Ticker               string `json:"T"`
	AskPrice             float64 `json:"P"`
	AskSize              int     `json:"S"`
	AskExchange          int     `json:"X"`
	Conditions           []int   `json:"c"`
	TrfTimestamp         int64   `json:"f"`
	Indicators           []int   `json:"i"`
	BidPrice             float64 `json:"p"`
	SequenceNumber       int64   `json:"q"`
	BidSize              int     `json:"s"`
	SipTimestamp         int64   `json:"t"`
	BidExchange          int     `json:"x"`
	ParticipantTimestamp int64   `json:"y"`
	Tape                 int     `json:"z"`
}

// GetTrades retrieves tick-level trade data for a specific stock ticker
// with optional timestamp filtering, sorting, and pagination. Each trade
// record includes price, size, exchange, conditions, and precise timestamps.
func (c *Client) GetTrades(ticker string, p TradesParams) (*TradesResponse, error) {
	path := fmt.Sprintf("/v3/trades/%s", ticker)

	params := map[string]string{
		"timestamp":     p.Timestamp,
		"timestamp.gte": p.TimestampGte,
		"timestamp.gt":  p.TimestampGt,
		"timestamp.lte": p.TimestampLte,
		"timestamp.lt":  p.TimestampLt,
		"order":         p.Order,
		"limit":         p.Limit,
		"sort":          p.Sort,
	}

	var result TradesResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetLastTrade retrieves the most recent trade for a specific stock ticker.
// Returns the last available trade with price, size, exchange, and timestamp
// information useful for monitoring current market activity.
func (c *Client) GetLastTrade(ticker string) (*LastTradeResponse, error) {
	path := fmt.Sprintf("/v2/last/trade/%s", ticker)

	var result LastTradeResponse
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetQuotes retrieves tick-level NBBO quote data for a specific stock ticker
// with optional timestamp filtering, sorting, and pagination. Each quote
// record includes bid/ask prices, sizes, exchange IDs, and precise timestamps.
func (c *Client) GetQuotes(ticker string, p QuotesParams) (*QuotesResponse, error) {
	path := fmt.Sprintf("/v3/quotes/%s", ticker)

	params := map[string]string{
		"timestamp":     p.Timestamp,
		"timestamp.gte": p.TimestampGte,
		"timestamp.gt":  p.TimestampGt,
		"timestamp.lte": p.TimestampLte,
		"timestamp.lt":  p.TimestampLt,
		"order":         p.Order,
		"limit":         p.Limit,
		"sort":          p.Sort,
	}

	var result QuotesResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetLastQuote retrieves the most recent NBBO quote for a specific stock
// ticker. Returns the last available bid/ask prices, sizes, and exchange
// information for real-time market monitoring.
func (c *Client) GetLastQuote(ticker string) (*LastQuoteResponse, error) {
	path := fmt.Sprintf("/v2/last/nbbo/%s", ticker)

	var result LastQuoteResponse
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
