//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

import (
	"fmt"
)

// -------------------------------------------------------------------
// Crypto Open/Close Types
// -------------------------------------------------------------------

// CryptoOpenCloseTrade represents a single trade in the opening or closing
// trades array of a crypto daily open/close response.
type CryptoOpenCloseTrade struct {
	Conditions []int   `json:"c"`
	ID         string  `json:"i"`
	Price      float64 `json:"p"`
	Size       float64 `json:"s"`
	Timestamp  int64   `json:"t"`
	Exchange   int     `json:"x"`
}

// CryptoOpenCloseResponse represents the API response for the daily
// open/close data of a crypto pair on a given date. It includes the
// opening and closing prices along with the individual trades.
type CryptoOpenCloseResponse struct {
	Symbol        string                 `json:"symbol"`
	IsUTC         bool                   `json:"isUTC"`
	Day           string                 `json:"day"`
	Open          float64                `json:"open"`
	Close         float64                `json:"close"`
	OpenTrades    []CryptoOpenCloseTrade `json:"openTrades"`
	ClosingTrades []CryptoOpenCloseTrade `json:"closingTrades"`
}

// -------------------------------------------------------------------
// Crypto Snapshot Types
// -------------------------------------------------------------------

// CryptoSnapshotLastTrade represents the most recent trade in a crypto
// snapshot, including price, size, exchange, conditions, and timestamp.
type CryptoSnapshotLastTrade struct {
	Conditions []int   `json:"conditions"`
	Exchange   int     `json:"exchange"`
	Price      float64 `json:"price"`
	Size       float64 `json:"size"`
	Timestamp  int64   `json:"timestamp"`
}

// CryptoSnapshotTicker represents a single ticker's snapshot data in the
// crypto market. It contains the current day's bar, previous day's bar,
// latest minute bar, the last trade, fair market value, and change values.
type CryptoSnapshotTicker struct {
	Ticker          string                  `json:"ticker"`
	TodaysChange    float64                 `json:"todaysChange"`
	TodaysChangePct float64                 `json:"todaysChangePerc"`
	Updated         int64                   `json:"updated"`
	Day             SnapshotBar             `json:"day"`
	PrevDay         SnapshotBar             `json:"prevDay"`
	Min             SnapshotMinBar          `json:"min"`
	LastTrade       CryptoSnapshotLastTrade `json:"lastTrade"`
	FMV             float64                 `json:"fmv"`
}

// CryptoSnapshotResponse represents the API response for the full crypto
// market snapshot endpoint that returns data for all or filtered tickers.
type CryptoSnapshotResponse struct {
	Status    string                 `json:"status"`
	RequestID string                 `json:"request_id"`
	Tickers   []CryptoSnapshotTicker `json:"tickers"`
}

// CryptoSingleSnapshotResponse represents the API response for a single
// crypto ticker snapshot request.
type CryptoSingleSnapshotResponse struct {
	Status    string               `json:"status"`
	RequestID string               `json:"request_id"`
	Ticker    CryptoSnapshotTicker `json:"ticker"`
}

// CryptoSnapshotParams holds the optional query parameters for fetching
// a crypto market snapshot, allowing filtering by ticker symbols.
type CryptoSnapshotParams struct {
	Tickers string
}

// -------------------------------------------------------------------
// Crypto Unified Snapshot Types
// -------------------------------------------------------------------

// CryptoUnifiedSession represents session OHLC data within a unified
// crypto snapshot response.
type CryptoUnifiedSession struct {
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"change_percent"`
	Close         float64 `json:"close"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	Open          float64 `json:"open"`
	PreviousClose float64 `json:"previous_close"`
}

// CryptoUnifiedSnapshotResult represents a single ticker's data within
// a unified snapshot response from the /v3/snapshot endpoint.
type CryptoUnifiedSnapshotResult struct {
	Ticker       string               `json:"ticker"`
	Name         string               `json:"name"`
	Value        float64              `json:"value"`
	Type         string               `json:"type"`
	Timeframe    string               `json:"timeframe"`
	MarketStatus string               `json:"market_status"`
	LastUpdated  int64                `json:"last_updated"`
	Session      CryptoUnifiedSession `json:"session"`
	FMV          float64              `json:"fmv"`
	Error        string               `json:"error,omitempty"`
	Message      string               `json:"message,omitempty"`
}

// CryptoUnifiedSnapshotResponse represents the API response from the
// /v3/snapshot endpoint for crypto tickers. Contains status, pagination
// info, and an array of snapshot results.
type CryptoUnifiedSnapshotResponse struct {
	Status    string                        `json:"status"`
	RequestID string                        `json:"request_id"`
	NextURL   string                        `json:"next_url,omitempty"`
	Results   []CryptoUnifiedSnapshotResult `json:"results"`
}

// CryptoUnifiedSnapshotParams holds the query parameters for fetching
// a unified snapshot, primarily a comma-separated list of ticker symbols.
type CryptoUnifiedSnapshotParams struct {
	TickerAnyOf string
}

// -------------------------------------------------------------------
// Crypto Trades Types
// -------------------------------------------------------------------

// CryptoTrade represents a single trade record for a crypto pair from
// the /v3/trades endpoint. Fields include exchange-level identifiers,
// price, size, and nanosecond-precision timestamps.
type CryptoTrade struct {
	Conditions           []int   `json:"conditions"`
	Exchange             int     `json:"exchange"`
	ID                   string  `json:"id"`
	ParticipantTimestamp int64   `json:"participant_timestamp"`
	Price                float64 `json:"price"`
	Size                 float64 `json:"size"`
}

// CryptoTradesResponse represents the API response for tick-level crypto
// trade data from the /v3/trades endpoint with pagination support.
type CryptoTradesResponse struct {
	Status    string        `json:"status"`
	RequestID string        `json:"request_id"`
	NextURL   string        `json:"next_url"`
	Results   []CryptoTrade `json:"results"`
}

// CryptoTradesParams holds the query parameters for fetching tick-level
// crypto trade data including timestamp filters and pagination controls.
type CryptoTradesParams struct {
	Timestamp    string
	TimestampGte string
	TimestampGt  string
	TimestampLte string
	TimestampLt  string
	Order        string
	Limit        string
	Sort         string
}

// -------------------------------------------------------------------
// Crypto Last Trade Types
// -------------------------------------------------------------------

// CryptoLastTradeDetail represents the most recent trade for a crypto
// pair from the /v1/last/crypto endpoint with price, size, exchange,
// conditions, and timestamp fields.
type CryptoLastTradeDetail struct {
	Price      float64 `json:"price"`
	Size       float64 `json:"size"`
	Exchange   int     `json:"exchange"`
	Conditions []int   `json:"conditions"`
	Timestamp  int64   `json:"timestamp"`
}

// CryptoLastTradeResponse represents the API response for the most
// recent crypto trade from the /v1/last/crypto/{from}/{to} endpoint.
type CryptoLastTradeResponse struct {
	Status    string                `json:"status"`
	RequestID string               `json:"request_id"`
	Symbol    string                `json:"symbol"`
	Last      CryptoLastTradeDetail `json:"last"`
}

// -------------------------------------------------------------------
// Crypto Condition Codes Types
// -------------------------------------------------------------------

// ConditionCode represents a single condition code with its ID, type,
// name, asset class, and the data types it applies to.
type ConditionCode struct {
	ID            int      `json:"id"`
	Type          string   `json:"type"`
	Name          string   `json:"name"`
	AssetClass    string   `json:"asset_class"`
	DataTypes     []string `json:"data_types"`
	Legacy        bool     `json:"legacy"`
	Abbreviation  string   `json:"abbreviation,omitempty"`
	Description   string   `json:"description,omitempty"`
	ExchangeID    int      `json:"exchange_id,omitempty"`
	SIPMapping    string   `json:"sip_mapping,omitempty"`
}

// ConditionsResponse represents the API response for the reference
// conditions endpoint (/v3/reference/conditions).
type ConditionsResponse struct {
	Status    string          `json:"status"`
	RequestID string          `json:"request_id"`
	Count     int             `json:"count"`
	Results   []ConditionCode `json:"results"`
}

// ConditionsParams holds the optional query parameters for filtering
// condition codes by asset class and data type.
type ConditionsParams struct {
	AssetClass string
	DataType   string
}

// -------------------------------------------------------------------
// Crypto Ticker Overview Types
// -------------------------------------------------------------------

// CryptoTickerOverview represents the detailed reference information for
// a single crypto ticker from the /v3/reference/tickers/{ticker} endpoint.
type CryptoTickerOverview struct {
	Ticker         string `json:"ticker"`
	Name           string `json:"name"`
	Market         string `json:"market"`
	Locale         string `json:"locale"`
	Active         bool   `json:"active"`
	CurrencySymbol string `json:"currency_symbol"`
	CurrencyName   string `json:"currency_name"`
	BaseCurrencySymbol string `json:"base_currency_symbol"`
	BaseCurrencyName   string `json:"base_currency_name"`
	LastUpdatedUTC string `json:"last_updated_utc"`
}

// CryptoTickerOverviewResponse represents the API response for a single
// crypto ticker overview from the /v3/reference/tickers/{ticker} endpoint.
type CryptoTickerOverviewResponse struct {
	Status    string               `json:"status"`
	RequestID string               `json:"request_id"`
	Results   CryptoTickerOverview `json:"results"`
}

// -------------------------------------------------------------------
// Crypto Tickers Params (reuses TickersResponse from stocks.go)
// -------------------------------------------------------------------

// CryptoTickersParams holds the query parameters for searching and
// filtering crypto tickers from the reference endpoint.
type CryptoTickersParams struct {
	Search string
	Active string
	Limit  string
	Sort   string
	Order  string
}

// -------------------------------------------------------------------
// API Methods - Aggregates
// -------------------------------------------------------------------

// GetCryptoBars retrieves custom OHLC aggregate bar data for a specific
// crypto ticker over the time range specified in the BarsParams. This
// reuses the same BarsParams and BarsResponse types used by stocks.
func (c *Client) GetCryptoBars(ticker string, p BarsParams) (*BarsResponse, error) {
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

// GetCryptoDailyMarketSummary retrieves the grouped daily OHLC summary
// for all crypto tickers on the specified date. This reuses the
// MarketSummaryResponse type from stocks.
func (c *Client) GetCryptoDailyMarketSummary(date string, adjusted string) (*MarketSummaryResponse, error) {
	path := fmt.Sprintf("/v2/aggs/grouped/locale/global/market/crypto/%s", date)

	params := map[string]string{
		"adjusted": adjusted,
	}

	var result MarketSummaryResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCryptoDailyTickerSummary retrieves the daily open/close data for
// a specific crypto pair (from/to) on a given date. The response includes
// opening and closing trades along with aggregate OHLC data.
func (c *Client) GetCryptoDailyTickerSummary(from, to, date, adjusted string) (*CryptoOpenCloseResponse, error) {
	path := fmt.Sprintf("/v1/open-close/crypto/%s/%s/%s", from, to, date)

	params := map[string]string{
		"adjusted": adjusted,
	}

	var result CryptoOpenCloseResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCryptoPreviousDayBar retrieves the previous day's OHLC bar data
// for a specific crypto ticker. This reuses the BarsResponse type.
func (c *Client) GetCryptoPreviousDayBar(ticker string, adjusted string) (*BarsResponse, error) {
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

// -------------------------------------------------------------------
// API Methods - Market Operations
// -------------------------------------------------------------------

// GetCryptoConditions retrieves the list of condition codes for the
// crypto asset class from the /v3/reference/conditions endpoint.
func (c *Client) GetCryptoConditions() (*ConditionsResponse, error) {
	path := "/v3/reference/conditions"

	params := map[string]string{
		"asset_class": "crypto",
	}

	var result ConditionsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCryptoExchanges retrieves the list of known exchanges for the
// crypto asset class from the /v3/reference/exchanges endpoint.
func (c *Client) GetCryptoExchanges() (*ExchangesResponse, error) {
	path := "/v3/reference/exchanges"

	params := map[string]string{
		"asset_class": "crypto",
	}

	var result ExchangesResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// -------------------------------------------------------------------
// API Methods - Snapshots
// -------------------------------------------------------------------

// GetCryptoSnapshotFullMarket retrieves snapshot data for all crypto
// tickers or a filtered subset specified by a comma-separated list.
// Each ticker includes day, previous day, minute bars, last trade,
// and fair market value data.
func (c *Client) GetCryptoSnapshotFullMarket(p CryptoSnapshotParams) (*CryptoSnapshotResponse, error) {
	path := "/v2/snapshot/locale/global/markets/crypto/tickers"

	params := map[string]string{
		"tickers": p.Tickers,
	}

	var result CryptoSnapshotResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCryptoSnapshotSingleTicker retrieves the most recent snapshot for
// a single crypto ticker, including the current day's bar, previous day's
// bar, latest minute bar, the last trade, and fair market value.
func (c *Client) GetCryptoSnapshotSingleTicker(ticker string) (*CryptoSingleSnapshotResponse, error) {
	path := fmt.Sprintf("/v2/snapshot/locale/global/markets/crypto/tickers/%s", ticker)

	var result CryptoSingleSnapshotResponse
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCryptoSnapshotTopMovers retrieves the current top crypto gainers or
// losers. The direction parameter must be either "gainers" or "losers".
func (c *Client) GetCryptoSnapshotTopMovers(direction string) (*CryptoSnapshotResponse, error) {
	path := fmt.Sprintf("/v2/snapshot/locale/global/markets/crypto/%s", direction)

	var result CryptoSnapshotResponse
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCryptoUnifiedSnapshot retrieves unified snapshot data for crypto
// tickers from the /v3/snapshot endpoint. Supports filtering by a
// comma-separated list of ticker symbols.
func (c *Client) GetCryptoUnifiedSnapshot(p CryptoUnifiedSnapshotParams) (*CryptoUnifiedSnapshotResponse, error) {
	path := "/v3/snapshot"

	params := map[string]string{
		"ticker.any_of": p.TickerAnyOf,
	}

	var result CryptoUnifiedSnapshotResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// -------------------------------------------------------------------
// API Methods - Technical Indicators
// -------------------------------------------------------------------

// GetCryptoSMA retrieves Simple Moving Average (SMA) data for the
// specified crypto ticker. SMA calculates the arithmetic mean of
// closing prices over a given window period.
func (c *Client) GetCryptoSMA(ticker string, p IndicatorParams) (*IndicatorResponse, error) {
	path := fmt.Sprintf("/v1/indicators/sma/%s", ticker)
	params := indicatorParamsToMap(p)

	var result IndicatorResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCryptoEMA retrieves Exponential Moving Average (EMA) data for the
// specified crypto ticker. EMA places greater weight on recent prices
// compared to SMA for quicker trend detection.
func (c *Client) GetCryptoEMA(ticker string, p IndicatorParams) (*IndicatorResponse, error) {
	path := fmt.Sprintf("/v1/indicators/ema/%s", ticker)
	params := indicatorParamsToMap(p)

	var result IndicatorResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCryptoRSI retrieves Relative Strength Index (RSI) data for the
// specified crypto ticker. RSI measures the speed and magnitude of
// price changes, oscillating between 0 and 100.
func (c *Client) GetCryptoRSI(ticker string, p IndicatorParams) (*IndicatorResponse, error) {
	path := fmt.Sprintf("/v1/indicators/rsi/%s", ticker)
	params := indicatorParamsToMap(p)

	var result IndicatorResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCryptoMACD retrieves Moving Average Convergence/Divergence (MACD)
// data for the specified crypto ticker. MACD is a momentum indicator
// calculated by subtracting the long-period EMA from the short-period EMA.
func (c *Client) GetCryptoMACD(ticker string, p MACDParams) (*MACDResponse, error) {
	path := fmt.Sprintf("/v1/indicators/macd/%s", ticker)
	params := macdParamsToMap(p)

	var result MACDResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// -------------------------------------------------------------------
// API Methods - Tickers
// -------------------------------------------------------------------

// GetCryptoTickers retrieves a list of crypto tickers matching the
// filter criteria specified in the CryptoTickersParams. It hardcodes
// the market=crypto parameter to limit results to crypto assets.
func (c *Client) GetCryptoTickers(p CryptoTickersParams) (*TickersResponse, error) {
	path := "/v3/reference/tickers"

	params := map[string]string{
		"market": "crypto",
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

// GetCryptoTickerOverview retrieves detailed reference information for
// a specific crypto ticker from the /v3/reference/tickers/{ticker} endpoint.
func (c *Client) GetCryptoTickerOverview(ticker string) (*CryptoTickerOverviewResponse, error) {
	path := fmt.Sprintf("/v3/reference/tickers/%s", ticker)

	var result CryptoTickerOverviewResponse
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// -------------------------------------------------------------------
// API Methods - Trades
// -------------------------------------------------------------------

// GetCryptoTrades retrieves tick-level trade data for a specific crypto
// ticker with optional timestamp filtering, sorting, and pagination.
func (c *Client) GetCryptoTrades(ticker string, p CryptoTradesParams) (*CryptoTradesResponse, error) {
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

	var result CryptoTradesResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCryptoLastTrade retrieves the most recent trade for a specific
// crypto pair (from/to currencies) from the /v1/last/crypto endpoint.
func (c *Client) GetCryptoLastTrade(from, to string) (*CryptoLastTradeResponse, error) {
	path := fmt.Sprintf("/v1/last/crypto/%s/%s", from, to)

	var result CryptoLastTradeResponse
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
