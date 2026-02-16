//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

import (
	"fmt"
)

// OptionsTradesResponse represents the API response for tick-level trade data
// returned by the /v3/trades/{optionsTicker} endpoint. It includes pagination
// via NextURL and a slice of individual OptionsTrade records.
type OptionsTradesResponse struct {
	Status    string         `json:"status"`
	RequestID string         `json:"request_id"`
	NextURL   string         `json:"next_url"`
	Results   []OptionsTrade `json:"results"`
}

// OptionsTrade represents a single tick-level trade record for an options contract.
// Fields include price, size, exchange, conditions, correction indicator, and
// nanosecond-precision timestamps from both the participant and the SIP.
type OptionsTrade struct {
	Conditions           []int   `json:"conditions"`
	Correction           int     `json:"correction"`
	Exchange             int     `json:"exchange"`
	ParticipantTimestamp int64   `json:"participant_timestamp"`
	Price                float64 `json:"price"`
	SequenceNumber       int64   `json:"sequence_number"`
	SipTimestamp         int64   `json:"sip_timestamp"`
	Size                 float64 `json:"size"`
}

// OptionsTradesParams holds the query parameters for fetching tick-level trade
// data from the /v3/trades/{optionsTicker} endpoint. Supports timestamp range
// filtering, sorting, and pagination controls.
type OptionsTradesParams struct {
	Timestamp    string
	TimestampGte string
	TimestampGt  string
	TimestampLte string
	TimestampLt  string
	Order        string
	Limit        string
	Sort         string
}

// OptionsLastTradeResponse represents the API response for the most recent trade
// of an options contract from the /v2/last/trade/{optionsTicker} endpoint.
type OptionsLastTradeResponse struct {
	Status    string           `json:"status"`
	RequestID string           `json:"request_id"`
	Results   OptionsLastTrade `json:"results"`
}

// OptionsLastTrade represents the most recent trade for an options contract.
// Fields use abbreviated single-character JSON keys from the API, including
// ticker symbol, price, size, exchange, trade ID, and multiple timestamps.
type OptionsLastTrade struct {
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

// OptionsQuotesResponse represents the API response for tick-level NBBO quote
// data returned by the /v3/quotes/{optionsTicker} endpoint. It includes
// pagination via NextURL and a slice of individual OptionsQuote records.
type OptionsQuotesResponse struct {
	Status    string         `json:"status"`
	RequestID string         `json:"request_id"`
	NextURL   string         `json:"next_url"`
	Results   []OptionsQuote `json:"results"`
}

// OptionsQuote represents a single NBBO quote record for an options contract,
// containing bid/ask prices, sizes, exchange information, sequence number,
// and nanosecond-precision SIP timestamp.
type OptionsQuote struct {
	AskExchange    int     `json:"ask_exchange"`
	AskPrice       float64 `json:"ask_price"`
	AskSize        float64 `json:"ask_size"`
	BidExchange    int     `json:"bid_exchange"`
	BidPrice       float64 `json:"bid_price"`
	BidSize        float64 `json:"bid_size"`
	SequenceNumber int64   `json:"sequence_number"`
	SipTimestamp   int64   `json:"sip_timestamp"`
}

// OptionsQuotesParams holds the query parameters for fetching tick-level NBBO
// quote data from the /v3/quotes/{optionsTicker} endpoint. Supports timestamp
// range filtering, sorting, and pagination controls.
type OptionsQuotesParams struct {
	Timestamp    string
	TimestampGte string
	TimestampGt  string
	TimestampLte string
	TimestampLt  string
	Order        string
	Limit        string
	Sort         string
}

// OptionsLastQuoteResponse represents the API response for the most recent
// NBBO quote of an options contract from the /v2/last/nbbo/{optionsTicker}
// endpoint.
type OptionsLastQuoteResponse struct {
	Status    string           `json:"status"`
	RequestID string           `json:"request_id"`
	Results   OptionsLastQuote `json:"results"`
}

// OptionsLastQuote represents the most recent NBBO quote for an options contract.
// Fields use abbreviated single-character JSON keys from the API where uppercase
// letters represent ask-side data and lowercase represent bid-side data.
type OptionsLastQuote struct {
	Ticker               string  `json:"T"`
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

// GetOptionsTrades retrieves tick-level trade data for a specific options contract
// ticker with optional timestamp filtering, sorting, and pagination. Each trade
// record includes price, size, exchange, conditions, and precise timestamps.
func (c *Client) GetOptionsTrades(ticker string, p OptionsTradesParams) (*OptionsTradesResponse, error) {
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

	var result OptionsTradesResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOptionsLastTrade retrieves the most recent trade for a specific options
// contract ticker. Returns the last available trade with price, size, exchange,
// and timestamp information useful for monitoring current options market activity.
func (c *Client) GetOptionsLastTrade(ticker string) (*OptionsLastTradeResponse, error) {
	path := fmt.Sprintf("/v2/last/trade/%s", ticker)

	var result OptionsLastTradeResponse
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOptionsQuotes retrieves tick-level NBBO quote data for a specific options
// contract ticker with optional timestamp filtering, sorting, and pagination.
// Each quote record includes bid/ask prices, sizes, exchange IDs, and precise
// timestamps.
func (c *Client) GetOptionsQuotes(ticker string, p OptionsQuotesParams) (*OptionsQuotesResponse, error) {
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

	var result OptionsQuotesResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOptionsLastQuote retrieves the most recent NBBO quote for a specific
// options contract ticker. Returns the last available bid/ask prices, sizes,
// and exchange information for real-time options market monitoring.
func (c *Client) GetOptionsLastQuote(ticker string) (*OptionsLastQuoteResponse, error) {
	path := fmt.Sprintf("/v2/last/nbbo/%s", ticker)

	var result OptionsLastQuoteResponse
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
