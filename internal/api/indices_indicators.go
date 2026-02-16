//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

import (
	"fmt"
)

// GetIndicesSMA retrieves Simple Moving Average (SMA) data for the specified
// index ticker (e.g., I:SPX). SMA calculates the arithmetic mean of values
// over a given window period, providing a smoothed trend line for the index.
// The response format is identical to the stocks SMA endpoint, so we reuse
// the IndicatorResponse type.
func (c *Client) GetIndicesSMA(ticker string, p IndicatorParams) (*IndicatorResponse, error) {
	path := fmt.Sprintf("/v1/indicators/sma/%s", ticker)
	params := indicatorParamsToMap(p)

	var result IndicatorResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetIndicesEMA retrieves Exponential Moving Average (EMA) data for the
// specified index ticker (e.g., I:SPX). EMA places greater weight on recent
// values compared to SMA, enabling quicker trend detection and more responsive
// signals for the index.
func (c *Client) GetIndicesEMA(ticker string, p IndicatorParams) (*IndicatorResponse, error) {
	path := fmt.Sprintf("/v1/indicators/ema/%s", ticker)
	params := indicatorParamsToMap(p)

	var result IndicatorResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetIndicesRSI retrieves Relative Strength Index (RSI) data for the specified
// index ticker (e.g., I:SPX). RSI measures the speed and magnitude of price
// changes, oscillating between 0 and 100 to help identify overbought or
// oversold conditions for the index.
func (c *Client) GetIndicesRSI(ticker string, p IndicatorParams) (*IndicatorResponse, error) {
	path := fmt.Sprintf("/v1/indicators/rsi/%s", ticker)
	params := indicatorParamsToMap(p)

	var result IndicatorResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetIndicesMACD retrieves Moving Average Convergence/Divergence (MACD) data
// for the specified index ticker (e.g., I:SPX). MACD is a momentum indicator
// calculated by subtracting the long-period EMA from the short-period EMA.
// The response includes the MACD line, signal line, and histogram values.
func (c *Client) GetIndicesMACD(ticker string, p MACDParams) (*MACDResponse, error) {
	path := fmt.Sprintf("/v1/indicators/macd/%s", ticker)
	params := macdParamsToMap(p)

	var result MACDResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
