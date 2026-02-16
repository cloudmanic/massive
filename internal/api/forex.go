//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

import (
	"fmt"
)

// --- Aggregates (reuse BarsResponse, Bar, MarketSummaryResponse types from stocks.go) ---

// ForexBarsParams holds the query parameters for fetching custom OHLC bar data
// from the forex aggregates endpoint.
type ForexBarsParams struct {
	Multiplier string
	Timespan   string
	From       string
	To         string
	Adjusted   string
	Sort       string
	Limit      string
}

// ForexMarketSummaryParams holds the query parameters for fetching a daily
// grouped forex market summary.
type ForexMarketSummaryParams struct {
	Adjusted string
}

// --- Currency Conversion ---

// ForexConversionLast holds the last quote data within a currency conversion
// response, including ask, bid, exchange identifier, and timestamp.
type ForexConversionLast struct {
	Ask       float64 `json:"ask"`
	Bid       float64 `json:"bid"`
	Exchange  int     `json:"exchange"`
	Timestamp int64   `json:"timestamp"`
}

// ForexConversionResponse represents the API response for converting one
// currency to another. It includes the converted amount, initial amount,
// and the last quote data used for the conversion.
type ForexConversionResponse struct {
	Converted     float64             `json:"converted"`
	From          string              `json:"from"`
	InitialAmount float64             `json:"initialAmount"`
	Last          ForexConversionLast `json:"last"`
	RequestID     string              `json:"request_id"`
	Status        string              `json:"status"`
	Symbol        string              `json:"symbol"`
	To            string              `json:"to"`
}

// ForexConversionParams holds the optional query parameters for the currency
// conversion endpoint, including the amount to convert and decimal precision.
type ForexConversionParams struct {
	Amount    string
	Precision string
}

// --- Quotes ---

// ForexQuote represents a single forex quote record containing ask and bid
// prices, exchange identifiers, and a participant timestamp.
type ForexQuote struct {
	AskExchange          int     `json:"ask_exchange"`
	AskPrice             float64 `json:"ask_price"`
	BidExchange          int     `json:"bid_exchange"`
	BidPrice             float64 `json:"bid_price"`
	ParticipantTimestamp int64   `json:"participant_timestamp"`
}

// ForexQuotesResponse represents the API response for forex tick-level quote
// data returned by the /v3/quotes/{fxTicker} endpoint. It includes pagination
// via NextURL and a slice of individual ForexQuote records.
type ForexQuotesResponse struct {
	Status    string       `json:"status"`
	RequestID string       `json:"request_id"`
	NextURL   string       `json:"next_url"`
	Results   []ForexQuote `json:"results"`
}

// ForexQuotesParams holds the query parameters for fetching forex tick-level
// quote data with optional timestamp filtering, sorting, and pagination.
type ForexQuotesParams struct {
	Timestamp    string
	TimestampGte string
	TimestampGt  string
	TimestampLte string
	TimestampLt  string
	Order        string
	Limit        string
	Sort         string
}

// ForexLastQuoteLast holds the last quote data within a forex last quote
// response, containing the most recent ask, bid, exchange, and timestamp.
type ForexLastQuoteLast struct {
	Ask       float64 `json:"ask"`
	Bid       float64 `json:"bid"`
	Exchange  int     `json:"exchange"`
	Timestamp int64   `json:"timestamp"`
}

// ForexLastQuoteResponse represents the API response for the most recent
// forex quote for a currency pair from the /v1/last_quote/currencies endpoint.
type ForexLastQuoteResponse struct {
	Last      ForexLastQuoteLast `json:"last"`
	RequestID string             `json:"request_id"`
	Status    string             `json:"status"`
	Symbol    string             `json:"symbol"`
}

// --- Snapshots ---

// ForexSnapshotDay holds the day-level OHLC data for a forex snapshot ticker.
type ForexSnapshotDay struct {
	Open  float64 `json:"o"`
	High  float64 `json:"h"`
	Low   float64 `json:"l"`
	Close float64 `json:"c"`
}

// ForexSnapshotLastQuote holds the last quote data within a forex snapshot,
// including ask, bid, exchange, and timestamp values.
type ForexSnapshotLastQuote struct {
	Ask       float64 `json:"a"`
	Bid       float64 `json:"b"`
	Exchange  int     `json:"x"`
	Timestamp int64   `json:"t"`
}

// ForexSnapshotTicker represents a single forex ticker's snapshot data
// containing the ticker symbol, day bar, last quote, previous day bar,
// and the calculated change values.
type ForexSnapshotTicker struct {
	Ticker          string                 `json:"ticker"`
	TodaysChange    float64                `json:"todaysChange"`
	TodaysChangePct float64                `json:"todaysChangePerc"`
	Updated         int64                  `json:"updated"`
	Day             ForexSnapshotDay       `json:"day"`
	LastQuote       ForexSnapshotLastQuote `json:"lastQuote"`
	PrevDay         ForexSnapshotDay       `json:"prevDay"`
}

// ForexSnapshotAllResponse represents the API response for a full forex
// market or filtered multi-ticker snapshot request.
type ForexSnapshotAllResponse struct {
	Status    string                `json:"status"`
	RequestID string                `json:"request_id"`
	Count     int                   `json:"count"`
	Tickers   []ForexSnapshotTicker `json:"tickers"`
}

// ForexSnapshotSingleResponse represents the API response for a single
// forex ticker snapshot request.
type ForexSnapshotSingleResponse struct {
	Status    string              `json:"status"`
	RequestID string              `json:"request_id"`
	Ticker    ForexSnapshotTicker `json:"ticker"`
}

// ForexSnapshotGainersLosersResponse represents the API response for the
// top forex market movers (gainers or losers) snapshot request.
type ForexSnapshotGainersLosersResponse struct {
	Status    string                `json:"status"`
	RequestID string                `json:"request_id"`
	Tickers   []ForexSnapshotTicker `json:"tickers"`
}

// ForexSnapshotAllParams holds the optional query parameters for fetching
// a full forex market or filtered multi-ticker snapshot.
type ForexSnapshotAllParams struct {
	Tickers string
}

// UnifiedSnapshotResult represents a single ticker result from the unified
// snapshot endpoint (/v3/snapshot).
type UnifiedSnapshotResult struct {
	Ticker          string           `json:"ticker"`
	TodaysChange    float64          `json:"todaysChange"`
	TodaysChangePct float64          `json:"todaysChangePerc"`
	Updated         int64            `json:"updated"`
	Day             ForexSnapshotDay `json:"day"`
	PrevDay         ForexSnapshotDay `json:"prevDay"`
}

// UnifiedSnapshotResponse represents the API response from the unified
// snapshot endpoint (/v3/snapshot) which supports any market type.
type UnifiedSnapshotResponse struct {
	Status    string                  `json:"status"`
	RequestID string                  `json:"request_id"`
	Results   []UnifiedSnapshotResult `json:"results"`
}

// --- Tickers (reuse TickersResponse/Ticker types from stocks.go) ---

// ForexTickerParams holds the query parameters for searching and filtering
// forex tickers from the reference endpoint.
type ForexTickerParams struct {
	Search string
	Active string
	Limit  string
	Sort   string
	Order  string
}

// ForexTickerOverviewResponse represents the API response for detailed
// information about a specific forex ticker from the reference endpoint.
type ForexTickerOverviewResponse struct {
	Status    string              `json:"status"`
	RequestID string              `json:"request_id"`
	Results   ForexTickerOverview `json:"results"`
}

// ForexTickerOverview represents the detailed reference data for a specific
// forex ticker, including market, locale, currency info, and active status.
type ForexTickerOverview struct {
	Ticker         string `json:"ticker"`
	Name           string `json:"name"`
	Market         string `json:"market"`
	Locale         string `json:"locale"`
	Active         bool   `json:"active"`
	CurrencySymbol string `json:"currency_symbol"`
	CurrencyName   string `json:"currency_name"`
	BaseCurrencySymbol string `json:"base_currency_symbol"`
	BaseCurrencyName   string `json:"base_currency_name"`
}

// --- API Methods ---

// GetForexBars retrieves custom OHLC aggregate bar data for a specific forex
// ticker over the time range specified in the ForexBarsParams. The response
// uses the shared BarsResponse type since the data format is identical.
func (c *Client) GetForexBars(ticker string, p ForexBarsParams) (*BarsResponse, error) {
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

// GetForexDailyMarketSummary retrieves the grouped daily OHLC summary for
// all forex tickers on the specified date. The response uses the shared
// MarketSummaryResponse type since the format is identical to stocks.
func (c *Client) GetForexDailyMarketSummary(date string, p ForexMarketSummaryParams) (*MarketSummaryResponse, error) {
	path := fmt.Sprintf("/v2/aggs/grouped/locale/global/market/fx/%s", date)

	params := map[string]string{
		"adjusted": p.Adjusted,
	}

	var result MarketSummaryResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetForexPreviousDayBar retrieves the previous day's OHLC bar data for a
// specific forex ticker. The response uses the shared BarsResponse type
// since the data format matches the aggregates endpoint.
func (c *Client) GetForexPreviousDayBar(ticker string, adjusted string) (*BarsResponse, error) {
	path := fmt.Sprintf("/v2/aggs/ticker/%s/prev", ticker)

	params := map[string]string{
		"adjusted": adjusted,
	}

	var result BarsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetForexConversion converts a specified amount from one currency to another
// using the latest exchange rate. The amount and precision parameters control
// the quantity converted and the number of decimal places in the result.
func (c *Client) GetForexConversion(from, to string, p ForexConversionParams) (*ForexConversionResponse, error) {
	path := fmt.Sprintf("/v1/conversion/%s/%s", from, to)

	params := map[string]string{
		"amount":    p.Amount,
		"precision": p.Precision,
	}

	var result ForexConversionResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetForexExchanges retrieves a list of known forex exchanges by calling the
// shared exchanges endpoint with the asset class filter set to "fx".
func (c *Client) GetForexExchanges() (*ExchangesResponse, error) {
	return c.GetExchanges(ExchangesParams{AssetClass: "fx"})
}

// GetForexMarketHolidays retrieves the list of upcoming market holidays.
// This is a convenience wrapper around the shared GetMarketHolidays method
// since forex and stock holidays come from the same endpoint.
func (c *Client) GetForexMarketHolidays() ([]MarketHoliday, error) {
	return c.GetMarketHolidays()
}

// GetForexMarketStatus retrieves the current real-time market status
// including forex exchange status. This is a convenience wrapper around
// the shared GetMarketStatus method.
func (c *Client) GetForexMarketStatus() (*MarketStatusResponse, error) {
	return c.GetMarketStatus()
}

// GetForexQuotes retrieves tick-level quote data for a specific forex ticker
// with optional timestamp filtering, sorting, and pagination. Each quote
// record includes ask/bid prices, exchange IDs, and participant timestamps.
func (c *Client) GetForexQuotes(ticker string, p ForexQuotesParams) (*ForexQuotesResponse, error) {
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

	var result ForexQuotesResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetForexLastQuote retrieves the most recent forex quote for a currency pair
// specified by the from and to currency codes. Returns the last available
// ask/bid prices with exchange and timestamp information.
func (c *Client) GetForexLastQuote(from, to string) (*ForexLastQuoteResponse, error) {
	path := fmt.Sprintf("/v1/last_quote/currencies/%s/%s", from, to)

	var result ForexLastQuoteResponse
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetForexSnapshotAll retrieves snapshot data for all forex tickers or a
// filtered subset specified by a comma-separated list of tickers in the
// params. Each snapshot includes day, previous day, and last quote data.
func (c *Client) GetForexSnapshotAll(p ForexSnapshotAllParams) (*ForexSnapshotAllResponse, error) {
	path := "/v2/snapshot/locale/global/markets/forex/tickers"

	params := map[string]string{
		"tickers": p.Tickers,
	}

	var result ForexSnapshotAllResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetForexSnapshotTicker retrieves the most recent snapshot for a single
// forex ticker, including the current day's bar, previous day's bar, and
// the last available quote data.
func (c *Client) GetForexSnapshotTicker(ticker string) (*ForexSnapshotSingleResponse, error) {
	path := fmt.Sprintf("/v2/snapshot/locale/global/markets/forex/tickers/%s", ticker)

	var result ForexSnapshotSingleResponse
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetForexGainersLosers retrieves the top forex market movers by direction.
// The direction parameter must be either "gainers" or "losers" and determines
// which set of movers is returned from the forex snapshot endpoint.
func (c *Client) GetForexGainersLosers(direction string) (*ForexSnapshotGainersLosersResponse, error) {
	path := fmt.Sprintf("/v2/snapshot/locale/global/markets/forex/%s", direction)

	var result ForexSnapshotGainersLosersResponse
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetForexUnifiedSnapshot retrieves snapshot data for the specified forex
// tickers using the unified snapshot endpoint (/v3/snapshot). The tickers
// parameter is a comma-separated list of forex ticker symbols.
func (c *Client) GetForexUnifiedSnapshot(tickers string) (*UnifiedSnapshotResponse, error) {
	path := "/v3/snapshot"

	params := map[string]string{
		"ticker.any_of": tickers,
	}

	var result UnifiedSnapshotResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetForexSMA retrieves Simple Moving Average (SMA) data for the specified
// forex ticker. SMA calculates the arithmetic mean of prices over a given
// window period, providing a smoothed trend line. Reuses the shared
// IndicatorResponse and IndicatorParams types.
func (c *Client) GetForexSMA(ticker string, p IndicatorParams) (*IndicatorResponse, error) {
	path := fmt.Sprintf("/v1/indicators/sma/%s", ticker)
	params := indicatorParamsToMap(p)

	var result IndicatorResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetForexEMA retrieves Exponential Moving Average (EMA) data for the
// specified forex ticker. EMA places greater weight on recent prices
// compared to SMA for quicker trend detection. Reuses the shared
// IndicatorResponse and IndicatorParams types.
func (c *Client) GetForexEMA(ticker string, p IndicatorParams) (*IndicatorResponse, error) {
	path := fmt.Sprintf("/v1/indicators/ema/%s", ticker)
	params := indicatorParamsToMap(p)

	var result IndicatorResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetForexRSI retrieves Relative Strength Index (RSI) data for the specified
// forex ticker. RSI measures the speed and magnitude of price changes,
// oscillating between 0 and 100 to identify overbought or oversold
// conditions. Reuses the shared IndicatorResponse and IndicatorParams types.
func (c *Client) GetForexRSI(ticker string, p IndicatorParams) (*IndicatorResponse, error) {
	path := fmt.Sprintf("/v1/indicators/rsi/%s", ticker)
	params := indicatorParamsToMap(p)

	var result IndicatorResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetForexMACD retrieves Moving Average Convergence/Divergence (MACD) data
// for the specified forex ticker. MACD is a momentum indicator calculated by
// subtracting the long-period EMA from the short-period EMA. The response
// includes the MACD line, signal line, and histogram values.
func (c *Client) GetForexMACD(ticker string, p MACDParams) (*MACDResponse, error) {
	path := fmt.Sprintf("/v1/indicators/macd/%s", ticker)
	params := macdParamsToMap(p)

	var result MACDResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetForexTickers retrieves a list of forex tickers matching the filter
// criteria specified in the ForexTickerParams. This calls the shared
// reference tickers endpoint with market=fx.
func (c *Client) GetForexTickers(p ForexTickerParams) (*TickersResponse, error) {
	path := "/v3/reference/tickers"

	params := map[string]string{
		"market": "fx",
		"search": p.Search,
		"active": p.Active,
		"limit":  p.Limit,
		"sort":   p.Sort,
		"order":  p.Order,
	}

	var result TickersResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetForexTickerOverview retrieves detailed reference data for a specific
// forex ticker including market, locale, currency names, and active status.
func (c *Client) GetForexTickerOverview(ticker string) (*ForexTickerOverviewResponse, error) {
	path := fmt.Sprintf("/v3/reference/tickers/%s", ticker)

	var result ForexTickerOverviewResponse
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
