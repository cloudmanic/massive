//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

import (
	"fmt"
)

// OptionsContract represents a single options contract from the reference
// data endpoint. It includes the contract's ticker, underlying ticker,
// type (call/put), exercise style, expiration date, strike price, shares
// per contract, primary exchange, CFI code, and any additional underlyings.
type OptionsContract struct {
	Ticker               string                   `json:"ticker"`
	UnderlyingTicker     string                   `json:"underlying_ticker"`
	ContractType         string                   `json:"contract_type"`
	ExerciseStyle        string                   `json:"exercise_style"`
	ExpirationDate       string                   `json:"expiration_date"`
	StrikePrice          float64                  `json:"strike_price"`
	SharesPerContract    int                      `json:"shares_per_contract"`
	PrimaryExchange      string                   `json:"primary_exchange"`
	CFI                  string                   `json:"cfi"`
	Correction           int                      `json:"correction"`
	AdditionalUnderlyings []AdditionalUnderlying  `json:"additional_underlyings"`
}

// AdditionalUnderlying represents an additional underlying asset associated
// with an options contract, including the underlying ticker symbol, the
// amount, and the type of underlying (equity or currency).
type AdditionalUnderlying struct {
	Underlying string  `json:"underlying"`
	Amount     float64 `json:"amount"`
	Type       string  `json:"type"`
}

// OptionsContractsResponse represents the API response for listing options
// contracts from the reference endpoint. It includes pagination support
// via NextURL for retrieving additional pages of results.
type OptionsContractsResponse struct {
	Status    string            `json:"status"`
	RequestID string           `json:"request_id"`
	Results   []OptionsContract `json:"results"`
	NextURL   string            `json:"next_url"`
}

// OptionsContractResponse represents the API response for retrieving a
// single options contract by its ticker. Unlike the list response, the
// Results field is a single OptionsContract object rather than an array.
type OptionsContractResponse struct {
	Status    string          `json:"status"`
	RequestID string         `json:"request_id"`
	Results   OptionsContract `json:"results"`
}

// OptionsContractsParams holds the query parameters for searching and
// filtering options contracts from the reference endpoint. It supports
// filtering by underlying ticker, contract type, expiration date, strike
// price, and various range filters using .gte/.gt/.lte/.lt suffixes.
type OptionsContractsParams struct {
	UnderlyingTicker    string
	ContractType        string
	ExpirationDate      string
	AsOf                string
	StrikePrice         string
	Expired             string
	UnderlyingTickerGte string
	UnderlyingTickerGt  string
	UnderlyingTickerLte string
	UnderlyingTickerLt  string
	ExpirationDateGte   string
	ExpirationDateGt    string
	ExpirationDateLte   string
	ExpirationDateLt    string
	StrikePriceGte      string
	StrikePriceGt       string
	StrikePriceLte      string
	StrikePriceLt       string
	Order               string
	Limit               string
	Sort                string
}

// GetOptionsContracts retrieves a list of options contracts matching the
// filter criteria specified in the OptionsContractsParams. It supports
// filtering by underlying ticker, contract type, expiration date, strike
// price, and various range filters. Results are paginated and the NextURL
// field can be used to fetch additional pages.
func (c *Client) GetOptionsContracts(p OptionsContractsParams) (*OptionsContractsResponse, error) {
	path := "/v3/reference/options/contracts"

	params := map[string]string{
		"underlying_ticker":     p.UnderlyingTicker,
		"contract_type":         p.ContractType,
		"expiration_date":       p.ExpirationDate,
		"as_of":                 p.AsOf,
		"strike_price":          p.StrikePrice,
		"expired":               p.Expired,
		"underlying_ticker.gte": p.UnderlyingTickerGte,
		"underlying_ticker.gt":  p.UnderlyingTickerGt,
		"underlying_ticker.lte": p.UnderlyingTickerLte,
		"underlying_ticker.lt":  p.UnderlyingTickerLt,
		"expiration_date.gte":   p.ExpirationDateGte,
		"expiration_date.gt":    p.ExpirationDateGt,
		"expiration_date.lte":   p.ExpirationDateLte,
		"expiration_date.lt":    p.ExpirationDateLt,
		"strike_price.gte":      p.StrikePriceGte,
		"strike_price.gt":       p.StrikePriceGt,
		"strike_price.lte":      p.StrikePriceLte,
		"strike_price.lt":       p.StrikePriceLt,
		"order":                 p.Order,
		"limit":                 p.Limit,
		"sort":                  p.Sort,
	}

	var result OptionsContractsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOptionsContract retrieves detailed information about a single options
// contract identified by its options ticker (e.g., "O:AAPL260218C00190000").
// The optional asOf parameter allows querying a historical snapshot of the
// contract as of a specific date in YYYY-MM-DD format.
func (c *Client) GetOptionsContract(optionsTicker string, asOf string) (*OptionsContractResponse, error) {
	path := fmt.Sprintf("/v3/reference/options/contracts/%s", optionsTicker)

	params := map[string]string{
		"as_of": asOf,
	}

	var result OptionsContractResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
