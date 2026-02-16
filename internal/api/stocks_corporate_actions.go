//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

// DividendsResponse represents the API response for listing historical
// cash dividend distributions. It includes pagination support via NextURL.
type DividendsResponse struct {
	Status    string     `json:"status"`
	RequestID string     `json:"request_id"`
	NextURL   string     `json:"next_url,omitempty"`
	Results   []Dividend `json:"results"`
}

// Dividend represents a single cash dividend distribution event for a
// stock ticker. It includes key dates (declaration, ex-dividend, record,
// pay), the cash payout amount, frequency classification, distribution
// type, and adjustment factors for normalizing historical data after splits.
type Dividend struct {
	ID                        string  `json:"id"`
	Ticker                    string  `json:"ticker"`
	DeclarationDate           string  `json:"declaration_date,omitempty"`
	ExDividendDate            string  `json:"ex_dividend_date"`
	RecordDate                string  `json:"record_date"`
	PayDate                   string  `json:"pay_date"`
	Frequency                 int     `json:"frequency"`
	CashAmount                float64 `json:"cash_amount"`
	Currency                  string  `json:"currency"`
	DistributionType          string  `json:"distribution_type"`
	HistoricalAdjustmentFactor float64 `json:"historical_adjustment_factor"`
	SplitAdjustedCashAmount   float64 `json:"split_adjusted_cash_amount"`
}

// SplitsResponse represents the API response for listing historical
// stock split events. It includes pagination support via NextURL.
type SplitsResponse struct {
	Status    string  `json:"status"`
	RequestID string  `json:"request_id"`
	NextURL   string  `json:"next_url,omitempty"`
	Results   []Split `json:"results"`
}

// Split represents a single stock split event including the execution
// date, the split ratio (split_from and split_to), the type of adjustment
// (forward_split, reverse_split, or stock_dividend), and a cumulative
// historical adjustment factor for normalizing historical price data.
type Split struct {
	ID                        string  `json:"id"`
	Ticker                    string  `json:"ticker"`
	ExecutionDate             string  `json:"execution_date"`
	SplitFrom                 float64 `json:"split_from"`
	SplitTo                   float64 `json:"split_to"`
	AdjustmentType            string  `json:"adjustment_type"`
	HistoricalAdjustmentFactor float64 `json:"historical_adjustment_factor"`
}

// DividendsParams holds the query parameters for fetching historical
// dividend data from the dividends endpoint. Supports filtering by
// ticker, ex-dividend date range, frequency, distribution type, and
// result ordering/limiting.
type DividendsParams struct {
	Ticker           string
	ExDividendDate   string
	ExDividendDateGT string
	ExDividendDateGTE string
	ExDividendDateLT string
	ExDividendDateLTE string
	Frequency        string
	DistributionType string
	Sort             string
	Limit            string
}

// SplitsParams holds the query parameters for fetching historical
// stock split data from the splits endpoint. Supports filtering by
// ticker, execution date range, adjustment type, and result ordering/limiting.
type SplitsParams struct {
	Ticker           string
	ExecutionDate    string
	ExecutionDateGT  string
	ExecutionDateGTE string
	ExecutionDateLT  string
	ExecutionDateLTE string
	AdjustmentType   string
	Sort             string
	Limit            string
}

// GetDividends retrieves a list of historical cash dividend distributions
// matching the filter criteria specified in the DividendsParams. Results
// include payout amounts, key dates, frequency, and split-adjusted values.
func (c *Client) GetDividends(p DividendsParams) (*DividendsResponse, error) {
	path := "/stocks/v1/dividends"

	params := map[string]string{
		"ticker":                p.Ticker,
		"ex_dividend_date":      p.ExDividendDate,
		"ex_dividend_date.gt":   p.ExDividendDateGT,
		"ex_dividend_date.gte":  p.ExDividendDateGTE,
		"ex_dividend_date.lt":   p.ExDividendDateLT,
		"ex_dividend_date.lte":  p.ExDividendDateLTE,
		"frequency":             p.Frequency,
		"distribution_type":     p.DistributionType,
		"sort":                  p.Sort,
		"limit":                 p.Limit,
	}

	var result DividendsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSplits retrieves a list of historical stock split events matching
// the filter criteria specified in the SplitsParams. Results include
// execution dates, split ratios, adjustment types, and historical
// adjustment factors for price normalization.
func (c *Client) GetSplits(p SplitsParams) (*SplitsResponse, error) {
	path := "/stocks/v1/splits"

	params := map[string]string{
		"ticker":              p.Ticker,
		"execution_date":      p.ExecutionDate,
		"execution_date.gt":   p.ExecutionDateGT,
		"execution_date.gte":  p.ExecutionDateGTE,
		"execution_date.lt":   p.ExecutionDateLT,
		"execution_date.lte":  p.ExecutionDateLTE,
		"adjustment_type":     p.AdjustmentType,
		"sort":                p.Sort,
		"limit":               p.Limit,
	}

	var result SplitsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
