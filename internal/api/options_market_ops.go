//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

// Note: The options market operations endpoints use the same API paths and
// response schemas as the stocks and indices market operations endpoints:
//   - Market Status:   GET /v1/marketstatus/now
//   - Market Holidays: GET /v1/marketstatus/upcoming
//
// Because the response structures are identical, the types defined in
// stocks_market_ops.go (MarketStatusResponse, MarketHoliday, etc.) are
// reused here. The methods below provide options-specific entry points
// that delegate to the same underlying API calls.

// GetOptionsMarketStatus retrieves the current real-time trading status
// of all US exchanges, currency markets, and index groups. This endpoint
// returns the same data as GetMarketStatus but is provided as a separate
// method for options-focused workflows. The response includes whether
// each exchange (NYSE, NASDAQ, OTC) is currently open, closed, or in
// extended-hours trading, which directly affects options contract trading
// availability.
func (c *Client) GetOptionsMarketStatus() (*MarketStatusResponse, error) {
	path := "/v1/marketstatus/now"

	var result MarketStatusResponse
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOptionsMarketHolidays retrieves the list of upcoming market holidays
// and early-close days for NYSE, NASDAQ, and OTC exchanges. This endpoint
// returns the same data as GetMarketHolidays but is provided as a separate
// method for options-focused workflows. The response is an array of
// MarketHoliday objects sorted by date, with optional open/close times
// for early-close days. Options traders use this to plan around expiration
// dates that fall near market closures.
func (c *Client) GetOptionsMarketHolidays() ([]MarketHoliday, error) {
	path := "/v1/marketstatus/upcoming"

	var result []MarketHoliday
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return result, nil
}
