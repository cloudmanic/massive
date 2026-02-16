//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- Futures Aggregate Bars Test Data ---

const futuresAggsJSON = `{
	"request_id": "fut-agg-123",
	"status": "OK",
	"results": [
		{
			"close": 4150.25,
			"dollar_volume": 125000000.50,
			"high": 4175.00,
			"low": 4130.50,
			"open": 4140.00,
			"session_end_date": "2025-03-15",
			"settlement_price": 4148.75,
			"ticker": "ESH5",
			"transactions": 85432,
			"volume": 30125.0,
			"window_start": 1710460800000000000
		},
		{
			"close": 4155.50,
			"dollar_volume": 130000000.00,
			"high": 4180.25,
			"low": 4135.00,
			"open": 4150.25,
			"session_end_date": "2025-03-16",
			"settlement_price": 4153.00,
			"ticker": "ESH5",
			"transactions": 92100,
			"volume": 31500.0,
			"window_start": 1710547200000000000
		}
	]
}`

// --- Futures Contracts Test Data ---

const futuresContractsJSON = `{
	"next_url": "https://api.massive.com/futures/vX/contracts?cursor=abc123",
	"request_id": "fut-con-456",
	"status": "OK",
	"results": [
		{
			"active": true,
			"date": "2025-03-15",
			"days_to_maturity": 45,
			"first_trade_date": "2024-09-15",
			"group_code": "ES",
			"last_trade_date": "2025-06-20",
			"max_order_quantity": 10000,
			"min_order_quantity": 1,
			"name": "E-mini S&P 500 Futures",
			"product_code": "ES",
			"settlement_date": "2025-06-20",
			"settlement_tick_size": 0.25,
			"spread_tick_size": 0.05,
			"ticker": "ESM5",
			"trade_tick_size": 0.25,
			"trading_venue": "CME",
			"type": "futures"
		},
		{
			"active": true,
			"date": "2025-03-15",
			"days_to_maturity": 135,
			"first_trade_date": "2024-12-15",
			"group_code": "ES",
			"last_trade_date": "2025-09-19",
			"max_order_quantity": 10000,
			"min_order_quantity": 1,
			"name": "E-mini S&P 500 Futures",
			"product_code": "ES",
			"settlement_date": "2025-09-19",
			"settlement_tick_size": 0.25,
			"spread_tick_size": 0.05,
			"ticker": "ESU5",
			"trade_tick_size": 0.25,
			"trading_venue": "CME",
			"type": "futures"
		}
	]
}`

// --- Futures Products Test Data ---

const futuresProductsJSON = `{
	"request_id": "fut-prod-789",
	"status": "OK",
	"results": [
		{
			"asset_class": "equity_index",
			"asset_sub_class": "large_cap",
			"date": "2025-03-15",
			"last_updated": "2025-03-15T10:00:00Z",
			"name": "E-mini S&P 500",
			"price_quotation": "USD per index point",
			"product_code": "ES",
			"sector": "index",
			"settlement_currency_code": "USD",
			"settlement_method": "cash",
			"settlement_type": "financial",
			"sub_sector": "equity",
			"trade_currency_code": "USD",
			"trading_venue": "CME",
			"type": "futures",
			"unit_of_measure": "index_point",
			"unit_of_measure_qty": 50.0
		}
	]
}`

// --- Futures Schedules Test Data ---

const futuresSchedulesJSON = `{
	"request_id": "fut-sched-101",
	"status": "OK",
	"results": [
		{
			"event": "settlement",
			"product_code": "ES",
			"product_name": "E-mini S&P 500",
			"session_end_date": "2025-03-15",
			"timestamp": "2025-03-15T16:00:00Z",
			"trading_venue": "CME"
		},
		{
			"event": "last_trade",
			"product_code": "ES",
			"product_name": "E-mini S&P 500",
			"session_end_date": "2025-06-20",
			"timestamp": "2025-06-20T09:30:00Z",
			"trading_venue": "CME"
		}
	]
}`

// --- Futures Exchanges Test Data ---

const futuresExchangesJSON = `{
	"count": 2,
	"results": [
		{
			"id": 1,
			"name": "Chicago Mercantile Exchange",
			"acronym": "CME",
			"mic": "XCME",
			"operating_mic": "XCME",
			"locale": "us",
			"type": "exchange",
			"url": "https://www.cmegroup.com"
		},
		{
			"id": 2,
			"name": "Intercontinental Exchange",
			"acronym": "ICE",
			"mic": "XICE",
			"operating_mic": "XICE",
			"locale": "us",
			"type": "exchange",
			"url": "https://www.theice.com"
		}
	]
}`

// --- Futures Snapshot Test Data ---

const futuresSnapshotJSON = `{
	"count": 1,
	"results": [
		{
			"details": {
				"open_interest": 245000,
				"settlement_date": "2025-06-20"
			},
			"last_minute": {
				"open": 4150.00,
				"high": 4152.50,
				"low": 4149.75,
				"close": 4151.25,
				"volume": 1250.0
			},
			"last_quote": {
				"ask_price": 4151.50,
				"ask_size": 125.0,
				"bid_price": 4151.00,
				"bid_size": 200.0,
				"ask_timestamp": 1710460800000000000,
				"bid_timestamp": 1710460800000000000
			},
			"last_trade": {
				"price": 4151.25,
				"size": 5.0,
				"timestamp": 1710460800000000000
			},
			"session": {
				"change": 11.25,
				"close": 4151.25,
				"high": 4175.00,
				"low": 4130.50,
				"open": 4140.00,
				"settlement_price": 4148.75,
				"volume": 30125.0
			},
			"product_code": "ES",
			"ticker": "ESM5"
		}
	]
}`

// --- Futures Trades Test Data ---

const futuresTradesJSON = `{
	"request_id": "fut-trades-201",
	"status": "OK",
	"results": [
		{
			"price": 4151.25,
			"report_sequence": 100001,
			"sequence_number": 500001,
			"session_end_date": "2025-03-15",
			"size": 5.0,
			"ticker": "ESM5",
			"timestamp": 1710460800000000000
		},
		{
			"price": 4151.50,
			"report_sequence": 100002,
			"sequence_number": 500002,
			"session_end_date": "2025-03-15",
			"size": 10.0,
			"ticker": "ESM5",
			"timestamp": 1710460801000000000
		}
	]
}`

// --- Futures Quotes Test Data ---

const futuresQuotesJSON = `{
	"request_id": "fut-quotes-301",
	"status": "OK",
	"results": [
		{
			"ask_price": 4151.50,
			"ask_size": 125.0,
			"ask_timestamp": 1710460800000000000,
			"bid_price": 4151.00,
			"bid_size": 200.0,
			"bid_timestamp": 1710460800000000000,
			"report_sequence": 200001,
			"sequence_number": 600001,
			"session_end_date": "2025-03-15",
			"ticker": "ESM5",
			"timestamp": 1710460800000000000
		},
		{
			"ask_price": 4151.75,
			"ask_size": 100.0,
			"ask_timestamp": 1710460801000000000,
			"bid_price": 4151.25,
			"bid_size": 150.0,
			"bid_timestamp": 1710460801000000000,
			"report_sequence": 200002,
			"sequence_number": 600002,
			"session_end_date": "2025-03-15",
			"ticker": "ESM5",
			"timestamp": 1710460801000000000
		}
	]
}`

// =====================
// Aggregate Bars Tests
// =====================

// TestGetFuturesAggs verifies that GetFuturesAggs correctly parses the API
// response and returns the expected aggregate bar data for a futures ticker.
func TestGetFuturesAggs(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/futures/vX/aggs/ESH5": futuresAggsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := FuturesAggParams{
		Resolution: "1day",
		Limit:      "10",
		Sort:       "asc",
	}

	result, err := client.GetFuturesAggs("ESH5", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "fut-agg-123" {
		t.Errorf("expected request_id fut-agg-123, got %s", result.RequestID)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 bars, got %d", len(result.Results))
	}

	bar := result.Results[0]
	if bar.Ticker != "ESH5" {
		t.Errorf("expected ticker ESH5, got %s", bar.Ticker)
	}

	if bar.Open != 4140.00 {
		t.Errorf("expected open 4140.00, got %f", bar.Open)
	}

	if bar.High != 4175.00 {
		t.Errorf("expected high 4175.00, got %f", bar.High)
	}

	if bar.Low != 4130.50 {
		t.Errorf("expected low 4130.50, got %f", bar.Low)
	}

	if bar.Close != 4150.25 {
		t.Errorf("expected close 4150.25, got %f", bar.Close)
	}

	if bar.Volume != 30125.0 {
		t.Errorf("expected volume 30125.0, got %f", bar.Volume)
	}

	if bar.DollarVolume != 125000000.50 {
		t.Errorf("expected dollar_volume 125000000.50, got %f", bar.DollarVolume)
	}

	if bar.SettlementPrice != 4148.75 {
		t.Errorf("expected settlement_price 4148.75, got %f", bar.SettlementPrice)
	}

	if bar.Transactions != 85432 {
		t.Errorf("expected transactions 85432, got %d", bar.Transactions)
	}

	if bar.SessionEndDate != "2025-03-15" {
		t.Errorf("expected session_end_date 2025-03-15, got %s", bar.SessionEndDate)
	}

	if bar.WindowStart != 1710460800000000000 {
		t.Errorf("expected window_start 1710460800000000000, got %d", bar.WindowStart)
	}
}

// TestGetFuturesAggsSecondBar verifies that the second bar in the aggregate
// response is correctly parsed with its own distinct values.
func TestGetFuturesAggsSecondBar(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/futures/vX/aggs/ESH5": futuresAggsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetFuturesAggs("ESH5", FuturesAggParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bar := result.Results[1]
	if bar.Close != 4155.50 {
		t.Errorf("expected close 4155.50, got %f", bar.Close)
	}

	if bar.Transactions != 92100 {
		t.Errorf("expected transactions 92100, got %d", bar.Transactions)
	}

	if bar.SessionEndDate != "2025-03-16" {
		t.Errorf("expected session_end_date 2025-03-16, got %s", bar.SessionEndDate)
	}
}

// TestGetFuturesAggsRequestPath verifies that GetFuturesAggs constructs the
// correct API path with the ticker in the URL.
func TestGetFuturesAggsRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(futuresAggsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetFuturesAggs("NQM5", FuturesAggParams{})

	expected := "/futures/vX/aggs/NQM5"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetFuturesAggsQueryParams verifies that all query parameters are
// correctly sent to the API endpoint for futures aggregates.
func TestGetFuturesAggsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("resolution") != "15mins" {
			t.Errorf("expected resolution=15mins, got %s", q.Get("resolution"))
		}
		if q.Get("window_start") != "2025-03-15" {
			t.Errorf("expected window_start=2025-03-15, got %s", q.Get("window_start"))
		}
		if q.Get("window_start.gte") != "2025-03-10" {
			t.Errorf("expected window_start.gte=2025-03-10, got %s", q.Get("window_start.gte"))
		}
		if q.Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "desc" {
			t.Errorf("expected sort=desc, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(futuresAggsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetFuturesAggs("ESH5", FuturesAggParams{
		Resolution:     "15mins",
		WindowStart:    "2025-03-15",
		WindowStartGte: "2025-03-10",
		Limit:          "50",
		Sort:           "desc",
	})
}

// TestGetFuturesAggsAPIError verifies that GetFuturesAggs returns an error
// when the API responds with a non-200 status.
func TestGetFuturesAggsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Ticker not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetFuturesAggs("INVALID", FuturesAggParams{})
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// =====================
// Contracts Tests
// =====================

// TestGetFuturesContracts verifies that GetFuturesContracts correctly parses
// the API response and returns the expected contract data with pagination.
func TestGetFuturesContracts(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/futures/vX/contracts": futuresContractsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := FuturesContractsParams{
		ProductCode: "ES",
		Active:      "true",
		Limit:       "10",
	}

	result, err := client.GetFuturesContracts(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "fut-con-456" {
		t.Errorf("expected request_id fut-con-456, got %s", result.RequestID)
	}

	if result.NextURL == "" {
		t.Error("expected next_url to be populated")
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 contracts, got %d", len(result.Results))
	}

	c := result.Results[0]
	if c.Ticker != "ESM5" {
		t.Errorf("expected ticker ESM5, got %s", c.Ticker)
	}

	if c.Name != "E-mini S&P 500 Futures" {
		t.Errorf("expected name E-mini S&P 500 Futures, got %s", c.Name)
	}

	if !c.Active {
		t.Error("expected active to be true")
	}

	if c.ProductCode != "ES" {
		t.Errorf("expected product_code ES, got %s", c.ProductCode)
	}

	if c.DaysToMaturity != 45 {
		t.Errorf("expected days_to_maturity 45, got %d", c.DaysToMaturity)
	}

	if c.TradingVenue != "CME" {
		t.Errorf("expected trading_venue CME, got %s", c.TradingVenue)
	}

	if c.TradeTickSize != 0.25 {
		t.Errorf("expected trade_tick_size 0.25, got %f", c.TradeTickSize)
	}

	if c.Type != "futures" {
		t.Errorf("expected type futures, got %s", c.Type)
	}

	if c.FirstTradeDate != "2024-09-15" {
		t.Errorf("expected first_trade_date 2024-09-15, got %s", c.FirstTradeDate)
	}

	if c.LastTradeDate != "2025-06-20" {
		t.Errorf("expected last_trade_date 2025-06-20, got %s", c.LastTradeDate)
	}

	if c.SettlementDate != "2025-06-20" {
		t.Errorf("expected settlement_date 2025-06-20, got %s", c.SettlementDate)
	}
}

// TestGetFuturesContractsSecondContract verifies that the second contract
// in the response is correctly parsed with its own distinct values.
func TestGetFuturesContractsSecondContract(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/futures/vX/contracts": futuresContractsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetFuturesContracts(FuturesContractsParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	c := result.Results[1]
	if c.Ticker != "ESU5" {
		t.Errorf("expected ticker ESU5, got %s", c.Ticker)
	}

	if c.DaysToMaturity != 135 {
		t.Errorf("expected days_to_maturity 135, got %d", c.DaysToMaturity)
	}
}

// TestGetFuturesContractsQueryParams verifies that all filter parameters are
// correctly sent to the API endpoint for futures contracts.
func TestGetFuturesContractsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("product_code") != "ES" {
			t.Errorf("expected product_code=ES, got %s", q.Get("product_code"))
		}
		if q.Get("ticker") != "ESM5" {
			t.Errorf("expected ticker=ESM5, got %s", q.Get("ticker"))
		}
		if q.Get("active") != "true" {
			t.Errorf("expected active=true, got %s", q.Get("active"))
		}
		if q.Get("type") != "futures" {
			t.Errorf("expected type=futures, got %s", q.Get("type"))
		}
		if q.Get("limit") != "25" {
			t.Errorf("expected limit=25, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "ticker" {
			t.Errorf("expected sort=ticker, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(futuresContractsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetFuturesContracts(FuturesContractsParams{
		ProductCode: "ES",
		Ticker:      "ESM5",
		Active:      "true",
		Type:        "futures",
		Limit:       "25",
		Sort:        "ticker",
	})
}

// TestGetFuturesContractsAPIError verifies that GetFuturesContracts returns
// an error when the API responds with a non-200 status.
func TestGetFuturesContractsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"ERROR","message":"Internal server error"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetFuturesContracts(FuturesContractsParams{})
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

// =====================
// Products Tests
// =====================

// TestGetFuturesProducts verifies that GetFuturesProducts correctly parses
// the API response and returns the expected product data.
func TestGetFuturesProducts(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/futures/vX/products": futuresProductsJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := FuturesProductsParams{
		ProductCode: "ES",
		Limit:       "10",
	}

	result, err := client.GetFuturesProducts(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "fut-prod-789" {
		t.Errorf("expected request_id fut-prod-789, got %s", result.RequestID)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 product, got %d", len(result.Results))
	}

	p := result.Results[0]
	if p.Name != "E-mini S&P 500" {
		t.Errorf("expected name E-mini S&P 500, got %s", p.Name)
	}

	if p.ProductCode != "ES" {
		t.Errorf("expected product_code ES, got %s", p.ProductCode)
	}

	if p.AssetClass != "equity_index" {
		t.Errorf("expected asset_class equity_index, got %s", p.AssetClass)
	}

	if p.AssetSubClass != "large_cap" {
		t.Errorf("expected asset_sub_class large_cap, got %s", p.AssetSubClass)
	}

	if p.Sector != "index" {
		t.Errorf("expected sector index, got %s", p.Sector)
	}

	if p.SubSector != "equity" {
		t.Errorf("expected sub_sector equity, got %s", p.SubSector)
	}

	if p.TradingVenue != "CME" {
		t.Errorf("expected trading_venue CME, got %s", p.TradingVenue)
	}

	if p.SettlementMethod != "cash" {
		t.Errorf("expected settlement_method cash, got %s", p.SettlementMethod)
	}

	if p.SettlementCurrencyCode != "USD" {
		t.Errorf("expected settlement_currency_code USD, got %s", p.SettlementCurrencyCode)
	}

	if p.UnitOfMeasureQty != 50.0 {
		t.Errorf("expected unit_of_measure_qty 50.0, got %f", p.UnitOfMeasureQty)
	}

	if p.Type != "futures" {
		t.Errorf("expected type futures, got %s", p.Type)
	}
}

// TestGetFuturesProductsQueryParams verifies that all filter parameters are
// correctly sent to the API endpoint for futures products.
func TestGetFuturesProductsQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("name") != "E-mini" {
			t.Errorf("expected name=E-mini, got %s", q.Get("name"))
		}
		if q.Get("product_code") != "ES" {
			t.Errorf("expected product_code=ES, got %s", q.Get("product_code"))
		}
		if q.Get("sector") != "index" {
			t.Errorf("expected sector=index, got %s", q.Get("sector"))
		}
		if q.Get("asset_class") != "equity_index" {
			t.Errorf("expected asset_class=equity_index, got %s", q.Get("asset_class"))
		}
		if q.Get("trading_venue") != "CME" {
			t.Errorf("expected trading_venue=CME, got %s", q.Get("trading_venue"))
		}
		if q.Get("type") != "futures" {
			t.Errorf("expected type=futures, got %s", q.Get("type"))
		}
		if q.Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "name" {
			t.Errorf("expected sort=name, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(futuresProductsJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetFuturesProducts(FuturesProductsParams{
		Name:         "E-mini",
		ProductCode:  "ES",
		Sector:       "index",
		AssetClass:   "equity_index",
		TradingVenue: "CME",
		Type:         "futures",
		Limit:        "50",
		Sort:         "name",
	})
}

// TestGetFuturesProductsEmptyResults verifies that GetFuturesProducts
// handles an empty results array without error.
func TestGetFuturesProductsEmptyResults(t *testing.T) {
	emptyJSON := `{"request_id":"abc","status":"OK","results":[]}`
	server := mockServer(t, map[string]string{
		"/futures/vX/products": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetFuturesProducts(FuturesProductsParams{Name: "nonexistent"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(result.Results))
	}
}

// =====================
// Schedules Tests
// =====================

// TestGetFuturesSchedules verifies that GetFuturesSchedules correctly
// parses the API response and returns the expected schedule events.
func TestGetFuturesSchedules(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/futures/vX/schedules": futuresSchedulesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := FuturesSchedulesParams{
		ProductCode: "ES",
		Limit:       "10",
	}

	result, err := client.GetFuturesSchedules(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "fut-sched-101" {
		t.Errorf("expected request_id fut-sched-101, got %s", result.RequestID)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 schedules, got %d", len(result.Results))
	}

	s := result.Results[0]
	if s.Event != "settlement" {
		t.Errorf("expected event settlement, got %s", s.Event)
	}

	if s.ProductCode != "ES" {
		t.Errorf("expected product_code ES, got %s", s.ProductCode)
	}

	if s.ProductName != "E-mini S&P 500" {
		t.Errorf("expected product_name E-mini S&P 500, got %s", s.ProductName)
	}

	if s.SessionEndDate != "2025-03-15" {
		t.Errorf("expected session_end_date 2025-03-15, got %s", s.SessionEndDate)
	}

	if s.TradingVenue != "CME" {
		t.Errorf("expected trading_venue CME, got %s", s.TradingVenue)
	}

	if s.Timestamp != "2025-03-15T16:00:00Z" {
		t.Errorf("expected timestamp 2025-03-15T16:00:00Z, got %s", s.Timestamp)
	}
}

// TestGetFuturesSchedulesSecondEvent verifies that the second schedule
// event in the response is correctly parsed.
func TestGetFuturesSchedulesSecondEvent(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/futures/vX/schedules": futuresSchedulesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetFuturesSchedules(FuturesSchedulesParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := result.Results[1]
	if s.Event != "last_trade" {
		t.Errorf("expected event last_trade, got %s", s.Event)
	}

	if s.SessionEndDate != "2025-06-20" {
		t.Errorf("expected session_end_date 2025-06-20, got %s", s.SessionEndDate)
	}
}

// TestGetFuturesSchedulesQueryParams verifies that all query parameters are
// correctly sent to the API endpoint for futures schedules.
func TestGetFuturesSchedulesQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("product_code") != "ES" {
			t.Errorf("expected product_code=ES, got %s", q.Get("product_code"))
		}
		if q.Get("session_end_date") != "2025-03-15" {
			t.Errorf("expected session_end_date=2025-03-15, got %s", q.Get("session_end_date"))
		}
		if q.Get("trading_venue") != "CME" {
			t.Errorf("expected trading_venue=CME, got %s", q.Get("trading_venue"))
		}
		if q.Get("limit") != "20" {
			t.Errorf("expected limit=20, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "timestamp" {
			t.Errorf("expected sort=timestamp, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(futuresSchedulesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetFuturesSchedules(FuturesSchedulesParams{
		ProductCode:    "ES",
		SessionEndDate: "2025-03-15",
		TradingVenue:   "CME",
		Limit:          "20",
		Sort:           "timestamp",
	})
}

// =====================
// Exchanges Tests
// =====================

// TestGetFuturesExchanges verifies that GetFuturesExchanges correctly
// parses the API response and returns the expected exchange data.
func TestGetFuturesExchanges(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/futures/vX/exchanges": futuresExchangesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := FuturesExchangesParams{
		Limit: "10",
	}

	result, err := client.GetFuturesExchanges(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Count != 2 {
		t.Errorf("expected count 2, got %d", result.Count)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 exchanges, got %d", len(result.Results))
	}

	e := result.Results[0]
	if e.ID != 1 {
		t.Errorf("expected id 1, got %d", e.ID)
	}

	if e.Name != "Chicago Mercantile Exchange" {
		t.Errorf("expected name Chicago Mercantile Exchange, got %s", e.Name)
	}

	if e.Acronym != "CME" {
		t.Errorf("expected acronym CME, got %s", e.Acronym)
	}

	if e.MIC != "XCME" {
		t.Errorf("expected mic XCME, got %s", e.MIC)
	}

	if e.OperatingMIC != "XCME" {
		t.Errorf("expected operating_mic XCME, got %s", e.OperatingMIC)
	}

	if e.Locale != "us" {
		t.Errorf("expected locale us, got %s", e.Locale)
	}

	if e.Type != "exchange" {
		t.Errorf("expected type exchange, got %s", e.Type)
	}

	if e.URL != "https://www.cmegroup.com" {
		t.Errorf("expected url https://www.cmegroup.com, got %s", e.URL)
	}
}

// TestGetFuturesExchangesSecondExchange verifies that the second exchange
// in the response is correctly parsed.
func TestGetFuturesExchangesSecondExchange(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/futures/vX/exchanges": futuresExchangesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetFuturesExchanges(FuturesExchangesParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e := result.Results[1]
	if e.Name != "Intercontinental Exchange" {
		t.Errorf("expected name Intercontinental Exchange, got %s", e.Name)
	}

	if e.Acronym != "ICE" {
		t.Errorf("expected acronym ICE, got %s", e.Acronym)
	}
}

// TestGetFuturesExchangesQueryParams verifies that the limit parameter is
// correctly sent to the API endpoint for futures exchanges.
func TestGetFuturesExchangesQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("limit") != "5" {
			t.Errorf("expected limit=5, got %s", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(futuresExchangesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetFuturesExchanges(FuturesExchangesParams{Limit: "5"})
}

// TestGetFuturesExchangesAPIError verifies that GetFuturesExchanges returns
// an error when the API responds with a non-200 status.
func TestGetFuturesExchangesAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"status":"FORBIDDEN","message":"Unauthorized"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetFuturesExchanges(FuturesExchangesParams{})
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
}

// =====================
// Snapshot Tests
// =====================

// TestGetFuturesSnapshot verifies that GetFuturesSnapshot correctly parses
// the API response including nested details, minute, quote, trade, and session.
func TestGetFuturesSnapshot(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/futures/vX/snapshot": futuresSnapshotJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := FuturesSnapshotParams{
		ProductCode: "ES",
		Limit:       "10",
	}

	result, err := client.GetFuturesSnapshot(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Count != 1 {
		t.Errorf("expected count 1, got %d", result.Count)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(result.Results))
	}

	snap := result.Results[0]
	if snap.Ticker != "ESM5" {
		t.Errorf("expected ticker ESM5, got %s", snap.Ticker)
	}

	if snap.ProductCode != "ES" {
		t.Errorf("expected product_code ES, got %s", snap.ProductCode)
	}

	// Details
	if snap.Details.OpenInterest != 245000 {
		t.Errorf("expected open_interest 245000, got %d", snap.Details.OpenInterest)
	}

	if snap.Details.SettlementDate != "2025-06-20" {
		t.Errorf("expected settlement_date 2025-06-20, got %s", snap.Details.SettlementDate)
	}

	// Last Minute
	if snap.LastMinute.Open != 4150.00 {
		t.Errorf("expected last_minute open 4150.00, got %f", snap.LastMinute.Open)
	}

	if snap.LastMinute.Close != 4151.25 {
		t.Errorf("expected last_minute close 4151.25, got %f", snap.LastMinute.Close)
	}

	if snap.LastMinute.Volume != 1250.0 {
		t.Errorf("expected last_minute volume 1250.0, got %f", snap.LastMinute.Volume)
	}

	// Last Quote
	if snap.LastQuote.AskPrice != 4151.50 {
		t.Errorf("expected ask_price 4151.50, got %f", snap.LastQuote.AskPrice)
	}

	if snap.LastQuote.BidPrice != 4151.00 {
		t.Errorf("expected bid_price 4151.00, got %f", snap.LastQuote.BidPrice)
	}

	if snap.LastQuote.AskSize != 125.0 {
		t.Errorf("expected ask_size 125.0, got %f", snap.LastQuote.AskSize)
	}

	if snap.LastQuote.BidSize != 200.0 {
		t.Errorf("expected bid_size 200.0, got %f", snap.LastQuote.BidSize)
	}

	// Last Trade
	if snap.LastTrade.Price != 4151.25 {
		t.Errorf("expected last_trade price 4151.25, got %f", snap.LastTrade.Price)
	}

	if snap.LastTrade.Size != 5.0 {
		t.Errorf("expected last_trade size 5.0, got %f", snap.LastTrade.Size)
	}

	// Session
	if snap.Session.Change != 11.25 {
		t.Errorf("expected session change 11.25, got %f", snap.Session.Change)
	}

	if snap.Session.Open != 4140.00 {
		t.Errorf("expected session open 4140.00, got %f", snap.Session.Open)
	}

	if snap.Session.High != 4175.00 {
		t.Errorf("expected session high 4175.00, got %f", snap.Session.High)
	}

	if snap.Session.Low != 4130.50 {
		t.Errorf("expected session low 4130.50, got %f", snap.Session.Low)
	}

	if snap.Session.Close != 4151.25 {
		t.Errorf("expected session close 4151.25, got %f", snap.Session.Close)
	}

	if snap.Session.SettlementPrice != 4148.75 {
		t.Errorf("expected session settlement_price 4148.75, got %f", snap.Session.SettlementPrice)
	}

	if snap.Session.Volume != 30125.0 {
		t.Errorf("expected session volume 30125.0, got %f", snap.Session.Volume)
	}
}

// TestGetFuturesSnapshotQueryParams verifies that all query parameters are
// correctly sent to the API endpoint for futures snapshots.
func TestGetFuturesSnapshotQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("product_code") != "ES" {
			t.Errorf("expected product_code=ES, got %s", q.Get("product_code"))
		}
		if q.Get("ticker") != "ESM5" {
			t.Errorf("expected ticker=ESM5, got %s", q.Get("ticker"))
		}
		if q.Get("limit") != "5" {
			t.Errorf("expected limit=5, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "ticker" {
			t.Errorf("expected sort=ticker, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(futuresSnapshotJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetFuturesSnapshot(FuturesSnapshotParams{
		ProductCode: "ES",
		Ticker:      "ESM5",
		Limit:       "5",
		Sort:        "ticker",
	})
}

// TestGetFuturesSnapshotAPIError verifies that GetFuturesSnapshot returns
// an error when the API responds with a non-200 status.
func TestGetFuturesSnapshotAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"BAD_REQUEST","message":"Invalid parameters"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetFuturesSnapshot(FuturesSnapshotParams{})
	if err == nil {
		t.Fatal("expected error for 400 response, got nil")
	}
}

// =====================
// Trades Tests
// =====================

// TestGetFuturesTrades verifies that GetFuturesTrades correctly parses
// the API response and returns the expected trade data.
func TestGetFuturesTrades(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/futures/vX/trades/ESM5": futuresTradesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := FuturesTradesParams{
		SessionEndDate: "2025-03-15",
		Limit:          "10",
		Sort:           "timestamp",
	}

	result, err := client.GetFuturesTrades("ESM5", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "fut-trades-201" {
		t.Errorf("expected request_id fut-trades-201, got %s", result.RequestID)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 trades, got %d", len(result.Results))
	}

	trade := result.Results[0]
	if trade.Ticker != "ESM5" {
		t.Errorf("expected ticker ESM5, got %s", trade.Ticker)
	}

	if trade.Price != 4151.25 {
		t.Errorf("expected price 4151.25, got %f", trade.Price)
	}

	if trade.Size != 5.0 {
		t.Errorf("expected size 5.0, got %f", trade.Size)
	}

	if trade.ReportSequence != 100001 {
		t.Errorf("expected report_sequence 100001, got %d", trade.ReportSequence)
	}

	if trade.SequenceNumber != 500001 {
		t.Errorf("expected sequence_number 500001, got %d", trade.SequenceNumber)
	}

	if trade.SessionEndDate != "2025-03-15" {
		t.Errorf("expected session_end_date 2025-03-15, got %s", trade.SessionEndDate)
	}

	if trade.Timestamp != 1710460800000000000 {
		t.Errorf("expected timestamp 1710460800000000000, got %d", trade.Timestamp)
	}
}

// TestGetFuturesTradesSecondTrade verifies that the second trade in the
// response is correctly parsed with its own distinct values.
func TestGetFuturesTradesSecondTrade(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/futures/vX/trades/ESM5": futuresTradesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetFuturesTrades("ESM5", FuturesTradesParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	trade := result.Results[1]
	if trade.Price != 4151.50 {
		t.Errorf("expected price 4151.50, got %f", trade.Price)
	}

	if trade.Size != 10.0 {
		t.Errorf("expected size 10.0, got %f", trade.Size)
	}

	if trade.SequenceNumber != 500002 {
		t.Errorf("expected sequence_number 500002, got %d", trade.SequenceNumber)
	}
}

// TestGetFuturesTradesRequestPath verifies that GetFuturesTrades constructs
// the correct API path with the ticker in the URL.
func TestGetFuturesTradesRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(futuresTradesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetFuturesTrades("NQM5", FuturesTradesParams{})

	expected := "/futures/vX/trades/NQM5"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetFuturesTradesQueryParams verifies that all query parameters are
// correctly sent to the API endpoint for futures trades.
func TestGetFuturesTradesQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("session_end_date") != "2025-03-15" {
			t.Errorf("expected session_end_date=2025-03-15, got %s", q.Get("session_end_date"))
		}
		if q.Get("timestamp") != "2025-03-15" {
			t.Errorf("expected timestamp=2025-03-15, got %s", q.Get("timestamp"))
		}
		if q.Get("timestamp.gte") != "2025-03-10" {
			t.Errorf("expected timestamp.gte=2025-03-10, got %s", q.Get("timestamp.gte"))
		}
		if q.Get("limit") != "100" {
			t.Errorf("expected limit=100, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "timestamp" {
			t.Errorf("expected sort=timestamp, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(futuresTradesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetFuturesTrades("ESM5", FuturesTradesParams{
		Timestamp:      "2025-03-15",
		TimestampGte:   "2025-03-10",
		SessionEndDate: "2025-03-15",
		Limit:          "100",
		Sort:           "timestamp",
	})
}

// TestGetFuturesTradesAPIError verifies that GetFuturesTrades returns
// an error when the API responds with a non-200 status.
func TestGetFuturesTradesAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Ticker not found."}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetFuturesTrades("INVALID", FuturesTradesParams{})
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

// =====================
// Quotes Tests
// =====================

// TestGetFuturesQuotes verifies that GetFuturesQuotes correctly parses
// the API response and returns the expected quote data.
func TestGetFuturesQuotes(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/futures/vX/quotes/ESM5": futuresQuotesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	params := FuturesQuotesParams{
		SessionEndDate: "2025-03-15",
		Limit:          "10",
		Sort:           "timestamp",
	}

	result, err := client.GetFuturesQuotes("ESM5", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "OK" {
		t.Errorf("expected status OK, got %s", result.Status)
	}

	if result.RequestID != "fut-quotes-301" {
		t.Errorf("expected request_id fut-quotes-301, got %s", result.RequestID)
	}

	if len(result.Results) != 2 {
		t.Fatalf("expected 2 quotes, got %d", len(result.Results))
	}

	quote := result.Results[0]
	if quote.Ticker != "ESM5" {
		t.Errorf("expected ticker ESM5, got %s", quote.Ticker)
	}

	if quote.AskPrice != 4151.50 {
		t.Errorf("expected ask_price 4151.50, got %f", quote.AskPrice)
	}

	if quote.AskSize != 125.0 {
		t.Errorf("expected ask_size 125.0, got %f", quote.AskSize)
	}

	if quote.BidPrice != 4151.00 {
		t.Errorf("expected bid_price 4151.00, got %f", quote.BidPrice)
	}

	if quote.BidSize != 200.0 {
		t.Errorf("expected bid_size 200.0, got %f", quote.BidSize)
	}

	if quote.ReportSequence != 200001 {
		t.Errorf("expected report_sequence 200001, got %d", quote.ReportSequence)
	}

	if quote.SequenceNumber != 600001 {
		t.Errorf("expected sequence_number 600001, got %d", quote.SequenceNumber)
	}

	if quote.SessionEndDate != "2025-03-15" {
		t.Errorf("expected session_end_date 2025-03-15, got %s", quote.SessionEndDate)
	}

	if quote.AskTimestamp != 1710460800000000000 {
		t.Errorf("expected ask_timestamp 1710460800000000000, got %d", quote.AskTimestamp)
	}

	if quote.BidTimestamp != 1710460800000000000 {
		t.Errorf("expected bid_timestamp 1710460800000000000, got %d", quote.BidTimestamp)
	}

	if quote.Timestamp != 1710460800000000000 {
		t.Errorf("expected timestamp 1710460800000000000, got %d", quote.Timestamp)
	}
}

// TestGetFuturesQuotesSecondQuote verifies that the second quote in the
// response is correctly parsed with its own distinct values.
func TestGetFuturesQuotesSecondQuote(t *testing.T) {
	server := mockServer(t, map[string]string{
		"/futures/vX/quotes/ESM5": futuresQuotesJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetFuturesQuotes("ESM5", FuturesQuotesParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	quote := result.Results[1]
	if quote.AskPrice != 4151.75 {
		t.Errorf("expected ask_price 4151.75, got %f", quote.AskPrice)
	}

	if quote.BidPrice != 4151.25 {
		t.Errorf("expected bid_price 4151.25, got %f", quote.BidPrice)
	}

	if quote.AskSize != 100.0 {
		t.Errorf("expected ask_size 100.0, got %f", quote.AskSize)
	}

	if quote.BidSize != 150.0 {
		t.Errorf("expected bid_size 150.0, got %f", quote.BidSize)
	}

	if quote.SequenceNumber != 600002 {
		t.Errorf("expected sequence_number 600002, got %d", quote.SequenceNumber)
	}
}

// TestGetFuturesQuotesRequestPath verifies that GetFuturesQuotes constructs
// the correct API path with the ticker in the URL.
func TestGetFuturesQuotesRequestPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(futuresQuotesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetFuturesQuotes("CLM5", FuturesQuotesParams{})

	expected := "/futures/vX/quotes/CLM5"
	if receivedPath != expected {
		t.Errorf("expected path %s, got %s", expected, receivedPath)
	}
}

// TestGetFuturesQuotesQueryParams verifies that all query parameters are
// correctly sent to the API endpoint for futures quotes.
func TestGetFuturesQuotesQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("session_end_date") != "2025-03-15" {
			t.Errorf("expected session_end_date=2025-03-15, got %s", q.Get("session_end_date"))
		}
		if q.Get("timestamp") != "2025-03-15" {
			t.Errorf("expected timestamp=2025-03-15, got %s", q.Get("timestamp"))
		}
		if q.Get("timestamp.lte") != "2025-03-20" {
			t.Errorf("expected timestamp.lte=2025-03-20, got %s", q.Get("timestamp.lte"))
		}
		if q.Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", q.Get("limit"))
		}
		if q.Get("sort") != "timestamp" {
			t.Errorf("expected sort=timestamp, got %s", q.Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(futuresQuotesJSON))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.GetFuturesQuotes("ESM5", FuturesQuotesParams{
		Timestamp:      "2025-03-15",
		TimestampLte:   "2025-03-20",
		SessionEndDate: "2025-03-15",
		Limit:          "50",
		Sort:           "timestamp",
	})
}

// TestGetFuturesQuotesAPIError verifies that GetFuturesQuotes returns
// an error when the API responds with a non-200 status.
func TestGetFuturesQuotesAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"status":"UNAUTHORIZED","message":"Invalid API key"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	_, err := client.GetFuturesQuotes("ESM5", FuturesQuotesParams{})
	if err == nil {
		t.Fatal("expected error for 401 response, got nil")
	}
}

// TestGetFuturesQuotesEmptyResults verifies that GetFuturesQuotes handles
// an empty results array without error.
func TestGetFuturesQuotesEmptyResults(t *testing.T) {
	emptyJSON := `{"request_id":"abc","status":"OK","results":[]}`
	server := mockServer(t, map[string]string{
		"/futures/vX/quotes/ESM5": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetFuturesQuotes("ESM5", FuturesQuotesParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(result.Results))
	}
}

// TestGetFuturesTradesEmptyResults verifies that GetFuturesTrades handles
// an empty results array without error.
func TestGetFuturesTradesEmptyResults(t *testing.T) {
	emptyJSON := `{"request_id":"abc","status":"OK","results":[]}`
	server := mockServer(t, map[string]string{
		"/futures/vX/trades/ESM5": emptyJSON,
	})
	defer server.Close()

	client := newTestClient(server.URL)
	result, err := client.GetFuturesTrades("ESM5", FuturesTradesParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(result.Results))
	}
}
