//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

// MarketStatusExchanges holds the open/closed status for each major
// stock exchange (NYSE, NASDAQ, OTC) as reported by the market status API.
type MarketStatusExchanges struct {
	Nasdaq string `json:"nasdaq"`
	NYSE   string `json:"nyse"`
	OTC    string `json:"otc"`
}

// MarketStatusCurrencies holds the open/closed status for the crypto
// and foreign exchange currency markets.
type MarketStatusCurrencies struct {
	Crypto string `json:"crypto"`
	FX     string `json:"fx"`
}

// MarketStatusIndicesGroups holds the open/closed status for various
// index groups such as S&P, Dow Jones, NASDAQ, MSCI, and others.
type MarketStatusIndicesGroups struct {
	SAndP           string `json:"s_and_p"`
	SocieteGenerale string `json:"societe_generale"`
	MSCI            string `json:"msci"`
	FTSERussell     string `json:"ftse_russell"`
	MStar           string `json:"mstar"`
	MStarC          string `json:"mstarc"`
	CCCY            string `json:"cccy"`
	CGI             string `json:"cgi"`
	Nasdaq          string `json:"nasdaq"`
	DowJones        string `json:"dow_jones"`
}

// MarketStatusResponse represents the API response from the market status
// endpoint (/v1/marketstatus/now). It provides a real-time snapshot of
// whether exchanges, currencies, and indices are currently open or closed.
type MarketStatusResponse struct {
	AfterHours    bool                      `json:"afterHours"`
	Currencies    MarketStatusCurrencies    `json:"currencies"`
	EarlyHours    bool                      `json:"earlyHours"`
	Exchanges     MarketStatusExchanges     `json:"exchanges"`
	IndicesGroups MarketStatusIndicesGroups  `json:"indicesGroups"`
	Market        string                    `json:"market"`
	ServerTime    string                    `json:"serverTime"`
}

// MarketHoliday represents a single upcoming market holiday or early-close
// day for a specific exchange. The Open and Close fields are only populated
// when the Status is "early-close".
type MarketHoliday struct {
	Date     string `json:"date"`
	Exchange string `json:"exchange"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Open     string `json:"open,omitempty"`
	Close    string `json:"close,omitempty"`
}

// ExchangesResponse represents the API response from the reference
// exchanges endpoint (/v3/reference/exchanges). It contains a list of
// known exchanges along with pagination metadata.
type ExchangesResponse struct {
	Status    string     `json:"status"`
	RequestID string     `json:"request_id"`
	Count     int        `json:"count"`
	Results   []Exchange `json:"results"`
}

// Exchange represents a single exchange with its identifiers, name,
// asset class, locale, and other reference attributes.
type Exchange struct {
	ID             int    `json:"id"`
	Type           string `json:"type"`
	AssetClass     string `json:"asset_class"`
	Locale         string `json:"locale"`
	Name           string `json:"name"`
	Acronym        string `json:"acronym,omitempty"`
	MIC            string `json:"mic,omitempty"`
	OperatingMIC   string `json:"operating_mic,omitempty"`
	ParticipantID  string `json:"participant_id,omitempty"`
	URL            string `json:"url,omitempty"`
}

// ExchangesParams holds the optional query parameters for filtering
// exchanges by asset class and locale.
type ExchangesParams struct {
	AssetClass string
	Locale     string
}

// GetMarketStatus retrieves the current real-time status of all US stock
// exchanges, currency markets, and index groups. This includes whether
// the market is in regular hours, after hours, or early hours trading.
func (c *Client) GetMarketStatus() (*MarketStatusResponse, error) {
	path := "/v1/marketstatus/now"

	var result MarketStatusResponse
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetMarketHolidays retrieves the list of upcoming market holidays and
// early-close days for NYSE, NASDAQ, and OTC exchanges. The response
// is an array of MarketHoliday objects sorted by date.
func (c *Client) GetMarketHolidays() ([]MarketHoliday, error) {
	path := "/v1/marketstatus/upcoming"

	var result []MarketHoliday
	if err := c.get(path, nil, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetExchanges retrieves a list of known exchanges filtered by the
// optional asset class and locale parameters. Each exchange includes
// identifiers like MIC codes, participant IDs, and URLs.
func (c *Client) GetExchanges(p ExchangesParams) (*ExchangesResponse, error) {
	path := "/v3/reference/exchanges"

	params := map[string]string{
		"asset_class": p.AssetClass,
		"locale":      p.Locale,
	}

	var result ExchangesResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
