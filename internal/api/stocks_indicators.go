//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

import (
	"fmt"
)

// IndicatorValue represents a single data point returned by the SMA, EMA,
// or RSI technical indicator endpoints. Each value is paired with a
// millisecond timestamp indicating the period it was calculated for.
type IndicatorValue struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// MACDValue represents a single data point returned by the MACD technical
// indicator endpoint. It includes the MACD line value, the signal line
// value, and the histogram (difference between MACD and signal).
type MACDValue struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
	Signal    float64 `json:"signal"`
	Histogram float64 `json:"histogram"`
}

// IndicatorUnderlying holds a reference URL to the underlying aggregate
// data used to compute the technical indicator values.
type IndicatorUnderlying struct {
	URL string `json:"url"`
}

// IndicatorResults contains the computed indicator values along with a
// reference to the underlying aggregate data source.
type IndicatorResults struct {
	Underlying IndicatorUnderlying `json:"underlying"`
	Values     []IndicatorValue    `json:"values"`
}

// MACDResults contains the computed MACD indicator values along with a
// reference to the underlying aggregate data source.
type MACDResults struct {
	Underlying IndicatorUnderlying `json:"underlying"`
	Values     []MACDValue         `json:"values"`
}

// IndicatorResponse represents the API response for the SMA, EMA, and RSI
// technical indicator endpoints. It includes pagination support via NextURL.
type IndicatorResponse struct {
	Status    string           `json:"status"`
	RequestID string           `json:"request_id"`
	NextURL   string           `json:"next_url,omitempty"`
	Results   IndicatorResults `json:"results"`
}

// MACDResponse represents the API response for the MACD technical indicator
// endpoint. MACD values include the MACD line, signal line, and histogram.
type MACDResponse struct {
	Status    string      `json:"status"`
	RequestID string      `json:"request_id"`
	NextURL   string      `json:"next_url,omitempty"`
	Results   MACDResults `json:"results"`
}

// IndicatorParams holds the common query parameters shared by the SMA, EMA,
// and RSI technical indicator endpoints. These control the time range,
// calculation window, and result pagination.
type IndicatorParams struct {
	TimestampGTE    string
	TimestampGT     string
	TimestampLTE    string
	TimestampLT     string
	Timespan        string
	Adjusted        string
	Window          string
	SeriesType      string
	ExpandUnderlying string
	Order           string
	Limit           string
}

// MACDParams holds the query parameters for the MACD technical indicator
// endpoint. MACD uses three window parameters (short, long, signal) instead
// of a single window.
type MACDParams struct {
	TimestampGTE    string
	TimestampGT     string
	TimestampLTE    string
	TimestampLT     string
	Timespan        string
	Adjusted        string
	ShortWindow     string
	LongWindow      string
	SignalWindow    string
	SeriesType      string
	ExpandUnderlying string
	Order           string
	Limit           string
}

// indicatorParamsToMap converts an IndicatorParams struct into a map of
// query parameter key-value pairs suitable for passing to the client's
// get method. Empty values are excluded automatically by the client.
func indicatorParamsToMap(p IndicatorParams) map[string]string {
	return map[string]string{
		"timestamp.gte":    p.TimestampGTE,
		"timestamp.gt":     p.TimestampGT,
		"timestamp.lte":    p.TimestampLTE,
		"timestamp.lt":     p.TimestampLT,
		"timespan":         p.Timespan,
		"adjusted":         p.Adjusted,
		"window":           p.Window,
		"series_type":      p.SeriesType,
		"expand_underlying": p.ExpandUnderlying,
		"order":            p.Order,
		"limit":            p.Limit,
	}
}

// macdParamsToMap converts a MACDParams struct into a map of query parameter
// key-value pairs suitable for passing to the client's get method. Empty
// values are excluded automatically by the client.
func macdParamsToMap(p MACDParams) map[string]string {
	return map[string]string{
		"timestamp.gte":    p.TimestampGTE,
		"timestamp.gt":     p.TimestampGT,
		"timestamp.lte":    p.TimestampLTE,
		"timestamp.lt":     p.TimestampLT,
		"timespan":         p.Timespan,
		"adjusted":         p.Adjusted,
		"short_window":     p.ShortWindow,
		"long_window":      p.LongWindow,
		"signal_window":    p.SignalWindow,
		"series_type":      p.SeriesType,
		"expand_underlying": p.ExpandUnderlying,
		"order":            p.Order,
		"limit":            p.Limit,
	}
}

// GetSMA retrieves Simple Moving Average (SMA) data for the specified stock
// ticker. SMA calculates the arithmetic mean of closing prices over a given
// window period, providing a smoothed trend line for the stock.
func (c *Client) GetSMA(ticker string, p IndicatorParams) (*IndicatorResponse, error) {
	path := fmt.Sprintf("/v1/indicators/sma/%s", ticker)
	params := indicatorParamsToMap(p)

	var result IndicatorResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetEMA retrieves Exponential Moving Average (EMA) data for the specified
// stock ticker. EMA places greater weight on recent prices compared to SMA,
// enabling quicker trend detection and more responsive signals.
func (c *Client) GetEMA(ticker string, p IndicatorParams) (*IndicatorResponse, error) {
	path := fmt.Sprintf("/v1/indicators/ema/%s", ticker)
	params := indicatorParamsToMap(p)

	var result IndicatorResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetRSI retrieves Relative Strength Index (RSI) data for the specified stock
// ticker. RSI measures the speed and magnitude of price changes, oscillating
// between 0 and 100 to help identify overbought or oversold conditions.
func (c *Client) GetRSI(ticker string, p IndicatorParams) (*IndicatorResponse, error) {
	path := fmt.Sprintf("/v1/indicators/rsi/%s", ticker)
	params := indicatorParamsToMap(p)

	var result IndicatorResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetMACD retrieves Moving Average Convergence/Divergence (MACD) data for the
// specified stock ticker. MACD is a momentum indicator calculated by subtracting
// the long-period EMA from the short-period EMA. The response includes the MACD
// line, signal line, and histogram values.
func (c *Client) GetMACD(ticker string, p MACDParams) (*MACDResponse, error) {
	path := fmt.Sprintf("/v1/indicators/macd/%s", ticker)
	params := macdParamsToMap(p)

	var result MACDResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
