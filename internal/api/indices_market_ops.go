//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

// Note: The indices market operations endpoints use the same API paths and
// response schemas as the stocks market operations endpoints:
//   - Market Status:   GET /v1/marketstatus/now
//   - Market Holidays: GET /v1/marketstatus/upcoming
//
// Because the response structures are identical, the types defined in
// stocks_market_ops.go (MarketStatusResponse, MarketHoliday, etc.) are
// reused here. The methods below provide indices-specific entry points
// that delegate to the same underlying API calls.

// GetIndicesMarketStatus retrieves the current real-time trading status
// of all US exchanges, currency markets, and index groups. This endpoint
// returns the same data as GetMarketStatus but is provided as a separate
// method for indices-focused workflows. The response includes whether
// each index group (S&P, Dow Jones, NASDAQ, MSCI, etc.) is currently
// open, closed, or in extended-hours trading.
func (c *Client) GetIndicesMarketStatus() (*MarketStatusResponse, error) {
	path := "/v1/marketstatus/now"

	var result MarketStatusResponse
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetIndicesMarketHolidays retrieves the list of upcoming market holidays
// and early-close days for NYSE, NASDAQ, and OTC exchanges. This endpoint
// returns the same data as GetMarketHolidays but is provided as a separate
// method for indices-focused workflows. The response is an array of
// MarketHoliday objects sorted by date, with optional open/close times
// for early-close days.
func (c *Client) GetIndicesMarketHolidays() ([]MarketHoliday, error) {
	path := "/v1/marketstatus/upcoming"

	var result []MarketHoliday
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return result, nil
}
