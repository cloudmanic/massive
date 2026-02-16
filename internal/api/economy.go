//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

// InflationResult represents a single inflation observation containing
// headline and core measures from both the CPI and PCE indexes.
// Fields may be absent for a given date since CPI and PCE are released
// on different schedules.
type InflationResult struct {
	Date        string  `json:"date"`
	CPI         float64 `json:"cpi"`
	CPICore     float64 `json:"cpi_core"`
	PCE         float64 `json:"pce"`
	PCECore     float64 `json:"pce_core"`
	PCESpending float64 `json:"pce_spending"`
}

// InflationResponse represents the API response from the /fed/v1/inflation
// endpoint. It contains a paginated list of inflation observations with
// CPI and PCE data.
type InflationResponse struct {
	Status    string            `json:"status"`
	RequestID string           `json:"request_id"`
	NextURL   string           `json:"next_url"`
	Results   []InflationResult `json:"results"`
}

// InflationParams holds the query parameters for filtering inflation data
// by date range, sort order, and result count limit.
type InflationParams struct {
	Date    string
	DateGT  string
	DateGTE string
	DateLT  string
	DateLTE string
	Sort    string
	Limit   string
}

// LaborMarketResult represents a single labor market observation including
// unemployment rate, labor force participation, average hourly earnings,
// and job openings data from the Federal Reserve.
type LaborMarketResult struct {
	Date                       string  `json:"date"`
	UnemploymentRate           float64 `json:"unemployment_rate"`
	LaborForceParticipationRate float64 `json:"labor_force_participation_rate"`
	AvgHourlyEarnings          float64 `json:"avg_hourly_earnings"`
	JobOpenings                float64 `json:"job_openings"`
}

// LaborMarketResponse represents the API response from the /fed/v1/labor-market
// endpoint. It contains a paginated list of labor market indicator observations.
type LaborMarketResponse struct {
	Status    string              `json:"status"`
	RequestID string             `json:"request_id"`
	NextURL   string             `json:"next_url"`
	Results   []LaborMarketResult `json:"results"`
}

// LaborMarketParams holds the query parameters for filtering labor market
// data by date range, sort order, and result count limit.
type LaborMarketParams struct {
	Date    string
	DateGT  string
	DateGTE string
	DateLT  string
	DateLTE string
	Sort    string
	Limit   string
}

// TreasuryYieldResult represents a single treasury yield observation with
// yields across multiple maturities from 1-month to 30-year durations.
type TreasuryYieldResult struct {
	Date        string  `json:"date"`
	Yield1Month  float64 `json:"yield_1_month"`
	Yield3Month  float64 `json:"yield_3_month"`
	Yield6Month  float64 `json:"yield_6_month"`
	Yield1Year   float64 `json:"yield_1_year"`
	Yield2Year   float64 `json:"yield_2_year"`
	Yield3Year   float64 `json:"yield_3_year"`
	Yield5Year   float64 `json:"yield_5_year"`
	Yield7Year   float64 `json:"yield_7_year"`
	Yield10Year  float64 `json:"yield_10_year"`
	Yield20Year  float64 `json:"yield_20_year"`
	Yield30Year  float64 `json:"yield_30_year"`
}

// TreasuryYieldResponse represents the API response from the /fed/v1/treasury-yields
// endpoint. It contains a paginated list of daily treasury yield curve observations.
type TreasuryYieldResponse struct {
	Status    string                `json:"status"`
	RequestID string               `json:"request_id"`
	NextURL   string               `json:"next_url"`
	Results   []TreasuryYieldResult `json:"results"`
}

// TreasuryYieldParams holds the query parameters for filtering treasury yield
// data by date range, sort order, and result count limit.
type TreasuryYieldParams struct {
	Date    string
	DateGT  string
	DateGTE string
	DateLT  string
	DateLTE string
	Sort    string
	Limit   string
}

// buildEconomyParams converts the common date filtering, sort, and limit
// parameters used by all economy endpoints into a map suitable for the
// Client.get() method. Empty values are omitted from the map.
func buildEconomyParams(date, dateGT, dateGTE, dateLT, dateLTE, sort, limit string) map[string]string {
	return map[string]string{
		"date":     date,
		"date.gt":  dateGT,
		"date.gte": dateGTE,
		"date.lt":  dateLT,
		"date.lte": dateLTE,
		"sort":     sort,
		"limit":    limit,
	}
}

// GetInflation retrieves inflation indicator data from the Federal Reserve
// including headline and core CPI and PCE measures. Results can be filtered
// by date range and paginated using limit and sort parameters.
func (c *Client) GetInflation(p InflationParams) (*InflationResponse, error) {
	path := "/fed/v1/inflation"
	params := buildEconomyParams(p.Date, p.DateGT, p.DateGTE, p.DateLT, p.DateLTE, p.Sort, p.Limit)

	var result InflationResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetLaborMarket retrieves labor market indicator data from the Federal Reserve
// including unemployment rate, labor force participation, average hourly
// earnings, and job openings. Results can be filtered by date range.
func (c *Client) GetLaborMarket(p LaborMarketParams) (*LaborMarketResponse, error) {
	path := "/fed/v1/labor-market"
	params := buildEconomyParams(p.Date, p.DateGT, p.DateGTE, p.DateLT, p.DateLTE, p.Sort, p.Limit)

	var result LaborMarketResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetTreasuryYields retrieves daily treasury yield curve data from the Federal
// Reserve across multiple maturities from 1-month to 30-year durations.
// Results can be filtered by date range and paginated.
func (c *Client) GetTreasuryYields(p TreasuryYieldParams) (*TreasuryYieldResponse, error) {
	path := "/fed/v1/treasury-yields"
	params := buildEconomyParams(p.Date, p.DateGT, p.DateGTE, p.DateLT, p.DateLTE, p.Sort, p.Limit)

	var result TreasuryYieldResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
