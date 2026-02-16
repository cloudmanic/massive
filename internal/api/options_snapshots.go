//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

import (
	"fmt"
)

// OptionSnapshotDay represents the daily OHLC bar data for an options
// contract snapshot, including change values, previous close, volume,
// VWAP, and the last update timestamp in nanoseconds.
type OptionSnapshotDay struct {
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"change_percent"`
	Close         float64 `json:"close"`
	High          float64 `json:"high"`
	LastUpdated   int64   `json:"last_updated"`
	Low           float64 `json:"low"`
	Open          float64 `json:"open"`
	PreviousClose float64 `json:"previous_close"`
	Volume        float64 `json:"volume"`
	VWAP          float64 `json:"vwap"`
}

// OptionSnapshotDetails holds the contract specification details for
// an options contract including the contract type (call/put), exercise
// style, expiration date, shares per contract, strike price, and ticker.
type OptionSnapshotDetails struct {
	ContractType      string  `json:"contract_type"`
	ExerciseStyle     string  `json:"exercise_style"`
	ExpirationDate    string  `json:"expiration_date"`
	SharesPerContract float64 `json:"shares_per_contract"`
	StrikePrice       float64 `json:"strike_price"`
	Ticker            string  `json:"ticker"`
}

// OptionSnapshotGreeks contains the option Greeks values (delta, gamma,
// theta, vega) that measure the sensitivity of the option's price to
// various factors.
type OptionSnapshotGreeks struct {
	Delta float64 `json:"delta"`
	Gamma float64 `json:"gamma"`
	Theta float64 `json:"theta"`
	Vega  float64 `json:"vega"`
}

// OptionSnapshotLastQuote represents the most recent quote data for an
// options contract, including bid/ask prices and sizes, midpoint, the
// last update timestamp, and the data timeframe.
type OptionSnapshotLastQuote struct {
	Ask         float64 `json:"ask"`
	AskSize     float64 `json:"ask_size"`
	Bid         float64 `json:"bid"`
	BidSize     float64 `json:"bid_size"`
	LastUpdated int64   `json:"last_updated"`
	Midpoint    float64 `json:"midpoint"`
	Timeframe   string  `json:"timeframe"`
}

// OptionSnapshotLastTrade represents the most recent trade for an options
// contract, including the trade price, size, exchange, SIP timestamp,
// trade conditions, and the data timeframe.
type OptionSnapshotLastTrade struct {
	Conditions   []int   `json:"conditions"`
	Exchange     int     `json:"exchange"`
	Price        float64 `json:"price"`
	SipTimestamp int64   `json:"sip_timestamp"`
	Size         float64 `json:"size"`
	Timeframe    string  `json:"timeframe"`
}

// OptionSnapshotUnderlyingAsset contains information about the underlying
// stock for an options contract, including the current price, the change
// needed to reach the break-even price, ticker symbol, last update
// timestamp, and data timeframe.
type OptionSnapshotUnderlyingAsset struct {
	ChangeToBreakEven float64 `json:"change_to_break_even"`
	LastUpdated       int64   `json:"last_updated"`
	Price             float64 `json:"price"`
	Ticker            string  `json:"ticker"`
	Timeframe         string  `json:"timeframe"`
}

// OptionSnapshotResult represents a single options contract snapshot
// containing the break-even price, day bar, contract details, Greeks,
// implied volatility, last quote, last trade, open interest, and
// underlying asset information.
type OptionSnapshotResult struct {
	BreakEvenPrice    float64                       `json:"break_even_price"`
	Day               OptionSnapshotDay             `json:"day"`
	Details           OptionSnapshotDetails         `json:"details"`
	FMV               float64                       `json:"fmv"`
	FMVLastUpdated    int64                         `json:"fmv_last_updated"`
	Greeks            OptionSnapshotGreeks          `json:"greeks"`
	ImpliedVolatility float64                       `json:"implied_volatility"`
	LastQuote         OptionSnapshotLastQuote       `json:"last_quote"`
	LastTrade         OptionSnapshotLastTrade       `json:"last_trade"`
	OpenInterest      float64                       `json:"open_interest"`
	UnderlyingAsset   OptionSnapshotUnderlyingAsset `json:"underlying_asset"`
}

// OptionsChainSnapshotResponse represents the API response for the
// options chain snapshot endpoint (/v3/snapshot/options/{underlyingAsset}).
// It contains an array of option contract snapshots for a given underlying.
type OptionsChainSnapshotResponse struct {
	Status    string                 `json:"status"`
	RequestID string                 `json:"request_id"`
	NextURL   string                 `json:"next_url"`
	Results   []OptionSnapshotResult `json:"results"`
}

// OptionContractSnapshotResponse represents the API response for a single
// option contract snapshot endpoint
// (/v3/snapshot/options/{underlyingAsset}/{optionContract}).
// It contains the snapshot data for one specific option contract.
type OptionContractSnapshotResponse struct {
	Status    string               `json:"status"`
	RequestID string               `json:"request_id"`
	Results   OptionSnapshotResult `json:"results"`
}

// OptionsChainSnapshotParams holds the optional query parameters for
// fetching the options chain snapshot for an underlying asset. Supports
// filtering by strike price, expiration date, contract type, and pagination.
type OptionsChainSnapshotParams struct {
	StrikePrice        string
	ExpirationDate     string
	ContractType       string
	StrikePriceGTE     string
	StrikePriceGT      string
	StrikePriceLTE     string
	StrikePriceLT      string
	ExpirationDateGTE  string
	ExpirationDateGT   string
	ExpirationDateLTE  string
	ExpirationDateLT   string
	Order              string
	Limit              string
	Sort               string
}

// GetOptionsChainSnapshot retrieves snapshot data for all options contracts
// associated with a given underlying asset ticker. The response includes
// day bar, contract details, Greeks, implied volatility, last quote,
// last trade, open interest, and underlying asset data for each contract.
// Supports filtering by strike price, expiration date, and contract type.
func (c *Client) GetOptionsChainSnapshot(underlyingAsset string, p OptionsChainSnapshotParams) (*OptionsChainSnapshotResponse, error) {
	path := fmt.Sprintf("/v3/snapshot/options/%s", underlyingAsset)

	params := map[string]string{
		"strike_price":        p.StrikePrice,
		"expiration_date":     p.ExpirationDate,
		"contract_type":       p.ContractType,
		"strike_price.gte":    p.StrikePriceGTE,
		"strike_price.gt":     p.StrikePriceGT,
		"strike_price.lte":    p.StrikePriceLTE,
		"strike_price.lt":     p.StrikePriceLT,
		"expiration_date.gte": p.ExpirationDateGTE,
		"expiration_date.gt":  p.ExpirationDateGT,
		"expiration_date.lte": p.ExpirationDateLTE,
		"expiration_date.lt":  p.ExpirationDateLT,
		"order":               p.Order,
		"limit":               p.Limit,
		"sort":                p.Sort,
	}

	var result OptionsChainSnapshotResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOptionContractSnapshot retrieves the most recent snapshot for a
// single option contract identified by the underlying asset ticker and
// the option contract ticker. The snapshot includes the day bar, contract
// details, Greeks, implied volatility, last quote, last trade, open
// interest, and underlying asset information.
func (c *Client) GetOptionContractSnapshot(underlyingAsset, optionContract string) (*OptionContractSnapshotResponse, error) {
	path := fmt.Sprintf("/v3/snapshot/options/%s/%s", underlyingAsset, optionContract)

	var result OptionContractSnapshotResponse
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
