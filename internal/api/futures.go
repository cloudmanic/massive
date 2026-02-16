//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

import (
	"fmt"
)

// --- Aggregate Bars ---

// FuturesAggResponse represents the API response for futures aggregate
// bar data including request metadata and an array of bar results.
type FuturesAggResponse struct {
	RequestID string       `json:"request_id"`
	Status    string       `json:"status"`
	Results   []FuturesBar `json:"results"`
}

// FuturesBar represents a single futures OHLC aggregate bar with
// settlement price, volume, dollar volume, and nanosecond window start.
type FuturesBar struct {
	Close           float64 `json:"close"`
	DollarVolume    float64 `json:"dollar_volume"`
	High            float64 `json:"high"`
	Low             float64 `json:"low"`
	Open            float64 `json:"open"`
	SessionEndDate  string  `json:"session_end_date"`
	SettlementPrice float64 `json:"settlement_price"`
	Ticker          string  `json:"ticker"`
	Transactions    int64   `json:"transactions"`
	Volume          float64 `json:"volume"`
	WindowStart     int64   `json:"window_start"`
}

// FuturesAggParams holds the query parameters for fetching futures
// aggregate bar data from the aggregates endpoint.
type FuturesAggParams struct {
	Resolution     string
	WindowStart    string
	WindowStartGte string
	WindowStartGt  string
	WindowStartLte string
	WindowStartLt  string
	Limit          string
	Sort           string
}

// GetFuturesAggs retrieves aggregate bar data for a specific futures ticker
// with configurable resolution, time window, sorting, and result limits.
func (c *Client) GetFuturesAggs(ticker string, p FuturesAggParams) (*FuturesAggResponse, error) {
	path := fmt.Sprintf("/futures/vX/aggs/%s", ticker)

	params := map[string]string{
		"resolution":       p.Resolution,
		"window_start":     p.WindowStart,
		"window_start.gte": p.WindowStartGte,
		"window_start.gt":  p.WindowStartGt,
		"window_start.lte": p.WindowStartLte,
		"window_start.lt":  p.WindowStartLt,
		"limit":            p.Limit,
		"sort":             p.Sort,
	}

	var result FuturesAggResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// --- Contracts ---

// FuturesContractsResponse represents the API response for listing futures
// contracts with pagination support via NextURL.
type FuturesContractsResponse struct {
	NextURL   string            `json:"next_url"`
	RequestID string            `json:"request_id"`
	Status    string            `json:"status"`
	Results   []FuturesContract `json:"results"`
}

// FuturesContract represents a single futures contract with details
// including trade dates, settlement info, tick sizes, and trading venue.
type FuturesContract struct {
	Active             bool    `json:"active"`
	Date               string  `json:"date"`
	DaysToMaturity     int     `json:"days_to_maturity"`
	FirstTradeDate     string  `json:"first_trade_date"`
	GroupCode          string  `json:"group_code"`
	LastTradeDate      string  `json:"last_trade_date"`
	MaxOrderQuantity   float64 `json:"max_order_quantity"`
	MinOrderQuantity   float64 `json:"min_order_quantity"`
	Name               string  `json:"name"`
	ProductCode        string  `json:"product_code"`
	SettlementDate     string  `json:"settlement_date"`
	SettlementTickSize float64 `json:"settlement_tick_size"`
	SpreadTickSize     float64 `json:"spread_tick_size"`
	Ticker             string  `json:"ticker"`
	TradeTickSize      float64 `json:"trade_tick_size"`
	TradingVenue       string  `json:"trading_venue"`
	Type               string  `json:"type"`
}

// FuturesContractsParams holds the query parameters for filtering and
// paginating the list of futures contracts.
type FuturesContractsParams struct {
	Date           string
	ProductCode    string
	Ticker         string
	Active         string
	Type           string
	FirstTradeDate string
	LastTradeDate  string
	Limit          string
	Sort           string
}

// GetFuturesContracts retrieves a list of futures contracts matching the
// provided filter criteria including product code, ticker, type, and dates.
func (c *Client) GetFuturesContracts(p FuturesContractsParams) (*FuturesContractsResponse, error) {
	path := "/futures/vX/contracts"

	params := map[string]string{
		"date":             p.Date,
		"product_code":     p.ProductCode,
		"ticker":           p.Ticker,
		"active":           p.Active,
		"type":             p.Type,
		"first_trade_date": p.FirstTradeDate,
		"last_trade_date":  p.LastTradeDate,
		"limit":            p.Limit,
		"sort":             p.Sort,
	}

	var result FuturesContractsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// --- Products ---

// FuturesProductsResponse represents the API response for listing futures
// products with request metadata and an array of product results.
type FuturesProductsResponse struct {
	RequestID string           `json:"request_id"`
	Status    string           `json:"status"`
	Results   []FuturesProduct `json:"results"`
}

// FuturesProduct represents a single futures product definition with
// asset classification, settlement details, and unit of measure information.
type FuturesProduct struct {
	AssetClass             string  `json:"asset_class"`
	AssetSubClass          string  `json:"asset_sub_class"`
	Date                   string  `json:"date"`
	LastUpdated            string  `json:"last_updated"`
	Name                   string  `json:"name"`
	PriceQuotation         string  `json:"price_quotation"`
	ProductCode            string  `json:"product_code"`
	Sector                 string  `json:"sector"`
	SettlementCurrencyCode string  `json:"settlement_currency_code"`
	SettlementMethod       string  `json:"settlement_method"`
	SettlementType         string  `json:"settlement_type"`
	SubSector              string  `json:"sub_sector"`
	TradeCurrencyCode      string  `json:"trade_currency_code"`
	TradingVenue           string  `json:"trading_venue"`
	Type                   string  `json:"type"`
	UnitOfMeasure          string  `json:"unit_of_measure"`
	UnitOfMeasureQty       float64 `json:"unit_of_measure_qty"`
}

// FuturesProductsParams holds the query parameters for filtering futures
// products by name, code, sector, asset class, venue, and other attributes.
type FuturesProductsParams struct {
	Name          string
	ProductCode   string
	Date          string
	TradingVenue  string
	Sector        string
	SubSector     string
	AssetClass    string
	AssetSubClass string
	Type          string
	Limit         string
	Sort          string
}

// GetFuturesProducts retrieves a list of futures products matching the
// provided filter criteria such as sector, asset class, and trading venue.
func (c *Client) GetFuturesProducts(p FuturesProductsParams) (*FuturesProductsResponse, error) {
	path := "/futures/vX/products"

	params := map[string]string{
		"name":            p.Name,
		"product_code":    p.ProductCode,
		"date":            p.Date,
		"trading_venue":   p.TradingVenue,
		"sector":          p.Sector,
		"sub_sector":      p.SubSector,
		"asset_class":     p.AssetClass,
		"asset_sub_class": p.AssetSubClass,
		"type":            p.Type,
		"limit":           p.Limit,
		"sort":            p.Sort,
	}

	var result FuturesProductsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// --- Schedules ---

// FuturesSchedulesResponse represents the API response for listing futures
// schedule events with request metadata and an array of schedule results.
type FuturesSchedulesResponse struct {
	RequestID string            `json:"request_id"`
	Status    string            `json:"status"`
	Results   []FuturesSchedule `json:"results"`
}

// FuturesSchedule represents a single futures schedule event such as a
// trading halt, settlement, or session boundary for a given product.
type FuturesSchedule struct {
	Event          string `json:"event"`
	ProductCode    string `json:"product_code"`
	ProductName    string `json:"product_name"`
	SessionEndDate string `json:"session_end_date"`
	Timestamp      string `json:"timestamp"`
	TradingVenue   string `json:"trading_venue"`
}

// FuturesSchedulesParams holds the query parameters for filtering futures
// schedules by product code, session end date, and trading venue.
type FuturesSchedulesParams struct {
	ProductCode    string
	SessionEndDate string
	TradingVenue   string
	Limit          string
	Sort           string
}

// GetFuturesSchedules retrieves a list of futures schedule events matching
// the provided filters for product code, session date, and venue.
func (c *Client) GetFuturesSchedules(p FuturesSchedulesParams) (*FuturesSchedulesResponse, error) {
	path := "/futures/vX/schedules"

	params := map[string]string{
		"product_code":     p.ProductCode,
		"session_end_date": p.SessionEndDate,
		"trading_venue":    p.TradingVenue,
		"limit":            p.Limit,
		"sort":             p.Sort,
	}

	var result FuturesSchedulesResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// --- Exchanges ---

// FuturesExchangesResponse represents the API response for listing futures
// exchanges with a count and array of exchange results.
type FuturesExchangesResponse struct {
	Count   int               `json:"count"`
	Results []FuturesExchange `json:"results"`
}

// FuturesExchange represents a single futures exchange with identifying
// information including MIC codes, acronym, locale, type, and URL.
type FuturesExchange struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Acronym      string `json:"acronym"`
	MIC          string `json:"mic"`
	OperatingMIC string `json:"operating_mic"`
	Locale       string `json:"locale"`
	Type         string `json:"type"`
	URL          string `json:"url"`
}

// FuturesExchangesParams holds the query parameters for limiting the
// number of futures exchanges returned.
type FuturesExchangesParams struct {
	Limit string
}

// GetFuturesExchanges retrieves a list of known futures exchanges with
// their identifiers, MIC codes, and metadata.
func (c *Client) GetFuturesExchanges(p FuturesExchangesParams) (*FuturesExchangesResponse, error) {
	path := "/futures/vX/exchanges"

	params := map[string]string{
		"limit": p.Limit,
	}

	var result FuturesExchangesResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// --- Snapshots ---

// FuturesSnapshotResponse represents the API response for futures contract
// snapshots including a count and array of snapshot results.
type FuturesSnapshotResponse struct {
	Count   int                      `json:"count"`
	Results []FuturesSnapshotContract `json:"results"`
}

// FuturesSnapshotContract represents a single futures contract snapshot
// with nested details, last minute bar, last quote, last trade, and session.
type FuturesSnapshotContract struct {
	Details     FuturesSnapshotDetails   `json:"details"`
	LastMinute  FuturesSnapshotMinute    `json:"last_minute"`
	LastQuote   FuturesSnapshotLastQuote `json:"last_quote"`
	LastTrade   FuturesSnapshotLastTrade `json:"last_trade"`
	Session     FuturesSnapshotSession   `json:"session"`
	ProductCode string                   `json:"product_code"`
	Ticker      string                   `json:"ticker"`
}

// FuturesSnapshotDetails holds static contract details within a snapshot
// including open interest and settlement date.
type FuturesSnapshotDetails struct {
	OpenInterest   int64  `json:"open_interest"`
	SettlementDate string `json:"settlement_date"`
}

// FuturesSnapshotMinute holds the most recent minute bar data within a
// futures snapshot including OHLCV values.
type FuturesSnapshotMinute struct {
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

// FuturesSnapshotLastQuote holds the most recent quote data within a
// futures snapshot including bid/ask prices, sizes, and timestamps.
type FuturesSnapshotLastQuote struct {
	AskPrice     float64 `json:"ask_price"`
	AskSize      float64 `json:"ask_size"`
	BidPrice     float64 `json:"bid_price"`
	BidSize      float64 `json:"bid_size"`
	AskTimestamp int64   `json:"ask_timestamp"`
	BidTimestamp int64   `json:"bid_timestamp"`
}

// FuturesSnapshotLastTrade holds the most recent trade data within a
// futures snapshot including price, size, and nanosecond timestamp.
type FuturesSnapshotLastTrade struct {
	Price     float64 `json:"price"`
	Size      float64 `json:"size"`
	Timestamp int64   `json:"timestamp"`
}

// FuturesSnapshotSession holds the current trading session data within
// a futures snapshot including OHLC, settlement price, change, and volume.
type FuturesSnapshotSession struct {
	Change          float64 `json:"change"`
	Close           float64 `json:"close"`
	High            float64 `json:"high"`
	Low             float64 `json:"low"`
	Open            float64 `json:"open"`
	SettlementPrice float64 `json:"settlement_price"`
	Volume          float64 `json:"volume"`
}

// FuturesSnapshotParams holds the query parameters for filtering futures
// contract snapshots by product code, ticker, limit, and sort order.
type FuturesSnapshotParams struct {
	ProductCode string
	Ticker      string
	Limit       string
	Sort        string
}

// GetFuturesSnapshot retrieves snapshot data for futures contracts matching
// the provided product code and ticker filters with pagination support.
func (c *Client) GetFuturesSnapshot(p FuturesSnapshotParams) (*FuturesSnapshotResponse, error) {
	path := "/futures/vX/snapshot"

	params := map[string]string{
		"product_code": p.ProductCode,
		"ticker":       p.Ticker,
		"limit":        p.Limit,
		"sort":         p.Sort,
	}

	var result FuturesSnapshotResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// --- Trades ---

// FuturesTradesResponse represents the API response for futures trade data
// with request metadata and an array of trade results.
type FuturesTradesResponse struct {
	RequestID string         `json:"request_id"`
	Status    string         `json:"status"`
	Results   []FuturesTrade `json:"results"`
}

// FuturesTrade represents a single futures trade with price, size, sequence
// numbers, session end date, ticker, and nanosecond timestamp.
type FuturesTrade struct {
	Price          float64 `json:"price"`
	ReportSequence int64   `json:"report_sequence"`
	SequenceNumber int64   `json:"sequence_number"`
	SessionEndDate string  `json:"session_end_date"`
	Size           float64 `json:"size"`
	Ticker         string  `json:"ticker"`
	Timestamp      int64   `json:"timestamp"`
}

// FuturesTradesParams holds the query parameters for filtering futures
// trades by timestamp, session end date, limit, and sort order.
type FuturesTradesParams struct {
	Timestamp      string
	TimestampGte   string
	TimestampGt    string
	TimestampLte   string
	TimestampLt    string
	SessionEndDate string
	Limit          string
	Sort           string
}

// GetFuturesTrades retrieves tick-level trade data for a specific futures
// ticker with optional timestamp and session date filtering.
func (c *Client) GetFuturesTrades(ticker string, p FuturesTradesParams) (*FuturesTradesResponse, error) {
	path := fmt.Sprintf("/futures/vX/trades/%s", ticker)

	params := map[string]string{
		"timestamp":        p.Timestamp,
		"timestamp.gte":    p.TimestampGte,
		"timestamp.gt":     p.TimestampGt,
		"timestamp.lte":    p.TimestampLte,
		"timestamp.lt":     p.TimestampLt,
		"session_end_date": p.SessionEndDate,
		"limit":            p.Limit,
		"sort":             p.Sort,
	}

	var result FuturesTradesResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// --- Quotes ---

// FuturesQuotesResponse represents the API response for futures quote data
// with request metadata and an array of quote results.
type FuturesQuotesResponse struct {
	RequestID string         `json:"request_id"`
	Status    string         `json:"status"`
	Results   []FuturesQuote `json:"results"`
}

// FuturesQuote represents a single futures quote with bid/ask prices and
// sizes, sequence numbers, session end date, and nanosecond timestamps.
type FuturesQuote struct {
	AskPrice       float64 `json:"ask_price"`
	AskSize        float64 `json:"ask_size"`
	AskTimestamp   int64   `json:"ask_timestamp"`
	BidPrice       float64 `json:"bid_price"`
	BidSize        float64 `json:"bid_size"`
	BidTimestamp   int64   `json:"bid_timestamp"`
	ReportSequence int64   `json:"report_sequence"`
	SequenceNumber int64   `json:"sequence_number"`
	SessionEndDate string  `json:"session_end_date"`
	Ticker         string  `json:"ticker"`
	Timestamp      int64   `json:"timestamp"`
}

// FuturesQuotesParams holds the query parameters for filtering futures
// quotes by timestamp, session end date, limit, and sort order.
type FuturesQuotesParams struct {
	Timestamp      string
	TimestampGte   string
	TimestampGt    string
	TimestampLte   string
	TimestampLt    string
	SessionEndDate string
	Limit          string
	Sort           string
}

// GetFuturesQuotes retrieves tick-level quote data for a specific futures
// ticker with optional timestamp and session date filtering.
func (c *Client) GetFuturesQuotes(ticker string, p FuturesQuotesParams) (*FuturesQuotesResponse, error) {
	path := fmt.Sprintf("/futures/vX/quotes/%s", ticker)

	params := map[string]string{
		"timestamp":        p.Timestamp,
		"timestamp.gte":    p.TimestampGte,
		"timestamp.gt":     p.TimestampGt,
		"timestamp.lte":    p.TimestampLte,
		"timestamp.lt":     p.TimestampLt,
		"session_end_date": p.SessionEndDate,
		"limit":            p.Limit,
		"sort":             p.Sort,
	}

	var result FuturesQuotesResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
