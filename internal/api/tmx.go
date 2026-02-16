//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

// TMXCorporateEventsResponse represents the API response for listing
// corporate events from the TMX/Wall Street Horizon data feed. It
// includes pagination support via NextURL for iterating through large
// result sets.
type TMXCorporateEventsResponse struct {
	Status    string              `json:"status"`
	Count     int                 `json:"count"`
	RequestID string              `json:"request_id"`
	NextURL   string              `json:"next_url,omitempty"`
	Results   []TMXCorporateEvent `json:"results"`
}

// TMXCorporateEvent represents a single corporate event record from
// the TMX/Wall Street Horizon calendar. Events include earnings
// announcements, dividend dates, investor conferences, stock splits,
// and other material corporate actions. Each event includes the
// company name, scheduled date, event type, current status, and
// exchange venue information.
type TMXCorporateEvent struct {
	CompanyName  string `json:"company_name"`
	Date         string `json:"date"`
	ISIN         string `json:"isin"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	Ticker       string `json:"ticker"`
	TMXCompanyID int    `json:"tmx_company_id"`
	TMXRecordID  string `json:"tmx_record_id"`
	TradingVenue string `json:"trading_venue"`
	Type         string `json:"type"`
	URL          string `json:"url"`
}

// TMXCorporateEventsParams holds the query parameters for fetching
// corporate events from the TMX endpoint. Supports filtering by date,
// event type, status, ticker, ISIN, trading venue, and TMX-specific
// identifiers. Date, type, status, ticker, ISIN, trading venue, and
// TMX record ID all support range operators (gt, gte, lt, lte) and
// the any_of multi-value filter.
type TMXCorporateEventsParams struct {
	Date            string
	DateAnyOf       string
	DateGT          string
	DateGTE         string
	DateLT          string
	DateLTE         string
	Type            string
	TypeAnyOf       string
	TypeGT          string
	TypeGTE         string
	TypeLT          string
	TypeLTE         string
	Status          string
	StatusAnyOf     string
	StatusGT        string
	StatusGTE       string
	StatusLT        string
	StatusLTE       string
	Ticker          string
	TickerAnyOf     string
	TickerGT        string
	TickerGTE       string
	TickerLT        string
	TickerLTE       string
	ISIN            string
	ISINAnyOf       string
	ISINGT          string
	ISINGTE         string
	ISINLT          string
	ISINLTE         string
	TradingVenue    string
	TradingVenueAnyOf string
	TradingVenueGT  string
	TradingVenueGTE string
	TradingVenueLT  string
	TradingVenueLTE string
	TMXCompanyID    string
	TMXCompanyIDGT  string
	TMXCompanyIDGTE string
	TMXCompanyIDLT  string
	TMXCompanyIDLTE string
	TMXRecordID     string
	TMXRecordIDAnyOf string
	TMXRecordIDGT   string
	TMXRecordIDGTE  string
	TMXRecordIDLT   string
	TMXRecordIDLTE  string
	Sort            string
	Limit           string
}

// GetTMXCorporateEvents retrieves a list of corporate events from the
// TMX/Wall Street Horizon data feed matching the filter criteria specified
// in the TMXCorporateEventsParams. Results include earnings announcements,
// dividend dates, investor conferences, stock splits, and other material
// corporate actions. Supports pagination via the NextURL field in the response.
func (c *Client) GetTMXCorporateEvents(p TMXCorporateEventsParams) (*TMXCorporateEventsResponse, error) {
	path := "/tmx/v1/corporate-events"

	params := map[string]string{
		"date":                   p.Date,
		"date.any_of":            p.DateAnyOf,
		"date.gt":                p.DateGT,
		"date.gte":               p.DateGTE,
		"date.lt":                p.DateLT,
		"date.lte":               p.DateLTE,
		"type":                   p.Type,
		"type.any_of":            p.TypeAnyOf,
		"type.gt":                p.TypeGT,
		"type.gte":               p.TypeGTE,
		"type.lt":                p.TypeLT,
		"type.lte":               p.TypeLTE,
		"status":                 p.Status,
		"status.any_of":          p.StatusAnyOf,
		"status.gt":              p.StatusGT,
		"status.gte":             p.StatusGTE,
		"status.lt":              p.StatusLT,
		"status.lte":             p.StatusLTE,
		"ticker":                 p.Ticker,
		"ticker.any_of":          p.TickerAnyOf,
		"ticker.gt":              p.TickerGT,
		"ticker.gte":             p.TickerGTE,
		"ticker.lt":              p.TickerLT,
		"ticker.lte":             p.TickerLTE,
		"isin":                   p.ISIN,
		"isin.any_of":            p.ISINAnyOf,
		"isin.gt":                p.ISINGT,
		"isin.gte":               p.ISINGTE,
		"isin.lt":                p.ISINLT,
		"isin.lte":               p.ISINLTE,
		"trading_venue":          p.TradingVenue,
		"trading_venue.any_of":   p.TradingVenueAnyOf,
		"trading_venue.gt":       p.TradingVenueGT,
		"trading_venue.gte":      p.TradingVenueGTE,
		"trading_venue.lt":       p.TradingVenueLT,
		"trading_venue.lte":      p.TradingVenueLTE,
		"tmx_company_id":         p.TMXCompanyID,
		"tmx_company_id.gt":      p.TMXCompanyIDGT,
		"tmx_company_id.gte":     p.TMXCompanyIDGTE,
		"tmx_company_id.lt":      p.TMXCompanyIDLT,
		"tmx_company_id.lte":     p.TMXCompanyIDLTE,
		"tmx_record_id":          p.TMXRecordID,
		"tmx_record_id.any_of":   p.TMXRecordIDAnyOf,
		"tmx_record_id.gt":       p.TMXRecordIDGT,
		"tmx_record_id.gte":      p.TMXRecordIDGTE,
		"tmx_record_id.lt":       p.TMXRecordIDLT,
		"tmx_record_id.lte":      p.TMXRecordIDLTE,
		"sort":                   p.Sort,
		"limit":                  p.Limit,
	}

	var result TMXCorporateEventsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
