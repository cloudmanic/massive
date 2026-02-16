//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

// SECFilingSectionsResponse represents the API response for retrieving
// plain-text content of specific sections from SEC 10-K filings. Includes
// pagination support via NextURL.
type SECFilingSectionsResponse struct {
	Status    string             `json:"status"`
	RequestID string             `json:"request_id"`
	NextURL   string             `json:"next_url"`
	Results   []SECFilingSection `json:"results"`
}

// SECFilingSection represents a single section extracted from a SEC 10-K
// filing, including the section identifier, raw text content, and filing
// metadata such as dates and the source URL on SEC.gov.
type SECFilingSection struct {
	CIK        string `json:"cik"`
	Ticker     string `json:"ticker"`
	Section    string `json:"section"`
	FilingDate string `json:"filing_date"`
	PeriodEnd  string `json:"period_end"`
	Text       string `json:"text"`
	FilingURL  string `json:"filing_url"`
}

// SECFilingSectionsParams holds the query parameters for fetching 10-K
// section content from the filings endpoint. Supports filtering by ticker,
// CIK, section type, filing date ranges, and period end date ranges.
type SECFilingSectionsParams struct {
	Ticker       string
	CIK          string
	Section      string
	FilingDate   string
	FilingDateGt string
	FilingDateLt string
	PeriodEnd    string
	PeriodEndGt  string
	PeriodEndLt  string
	Limit        string
	Sort         string
}

// RiskFactorsResponse represents the API response for retrieving
// standardized, machine-readable risk factor disclosures from SEC filings.
// Includes pagination support via NextURL.
type RiskFactorsResponse struct {
	Status    string       `json:"status"`
	RequestID string       `json:"request_id"`
	NextURL   string       `json:"next_url"`
	Results   []RiskFactor `json:"results"`
}

// RiskFactor represents a single categorized risk factor extracted from
// a company's SEC filing. Each risk is classified into a three-level
// taxonomy (primary, secondary, tertiary) with supporting text from the
// original filing.
type RiskFactor struct {
	CIK               string `json:"cik"`
	Ticker             string `json:"ticker"`
	PrimaryCategory    string `json:"primary_category"`
	SecondaryCategory  string `json:"secondary_category"`
	TertiaryCategory   string `json:"tertiary_category"`
	FilingDate         string `json:"filing_date"`
	SupportingText     string `json:"supporting_text"`
}

// RiskFactorsParams holds the query parameters for fetching risk factor
// disclosures. Supports filtering by ticker, CIK, and filing date ranges.
type RiskFactorsParams struct {
	Ticker       string
	CIK          string
	FilingDate   string
	FilingDateGt string
	FilingDateLt string
	Limit        string
	Sort         string
}

// RiskCategoriesResponse represents the API response for retrieving the
// taxonomy used to classify risk factors in SEC filings. Includes
// pagination support via NextURL.
type RiskCategoriesResponse struct {
	Status    string         `json:"status"`
	RequestID string         `json:"request_id"`
	NextURL   string         `json:"next_url"`
	Results   []RiskCategory `json:"results"`
}

// RiskCategory represents a single entry in the risk factor taxonomy,
// describing a three-level classification with a human-readable
// description and a taxonomy version identifier.
type RiskCategory struct {
	PrimaryCategory   string  `json:"primary_category"`
	SecondaryCategory string  `json:"secondary_category"`
	TertiaryCategory  string  `json:"tertiary_category"`
	Description       string  `json:"description"`
	Taxonomy          float64 `json:"taxonomy"`
}

// RiskCategoriesParams holds the query parameters for fetching the risk
// factor taxonomy. Supports filtering by category levels and taxonomy version.
type RiskCategoriesParams struct {
	PrimaryCategory   string
	SecondaryCategory string
	TertiaryCategory  string
	Taxonomy          string
	Limit             string
	Sort              string
}

// GetSECFilingSections retrieves plain-text content of specific sections
// from SEC 10-K filings for a given ticker or CIK. Supports filtering by
// section type (e.g., business, risk_factors), filing date, and period end
// date. Results are paginated and sorted by period_end descending by default.
func (c *Client) GetSECFilingSections(p SECFilingSectionsParams) (*SECFilingSectionsResponse, error) {
	path := "/stocks/filings/10-K/vX/sections"

	params := map[string]string{
		"ticker":          p.Ticker,
		"cik":             p.CIK,
		"section":         p.Section,
		"filing_date":     p.FilingDate,
		"filing_date.gt":  p.FilingDateGt,
		"filing_date.lt":  p.FilingDateLt,
		"period_end":      p.PeriodEnd,
		"period_end.gt":   p.PeriodEndGt,
		"period_end.lt":   p.PeriodEndLt,
		"limit":           p.Limit,
		"sort":            p.Sort,
	}

	var result SECFilingSectionsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetRiskFactors retrieves standardized, machine-readable risk factor
// disclosures from SEC filings. Each risk factor is classified using a
// three-level taxonomy and includes supporting text from the original
// filing. Supports filtering by ticker, CIK, and filing date ranges.
func (c *Client) GetRiskFactors(p RiskFactorsParams) (*RiskFactorsResponse, error) {
	path := "/stocks/filings/vX/risk-factors"

	params := map[string]string{
		"ticker":         p.Ticker,
		"cik":            p.CIK,
		"filing_date":    p.FilingDate,
		"filing_date.gt": p.FilingDateGt,
		"filing_date.lt": p.FilingDateLt,
		"limit":          p.Limit,
		"sort":           p.Sort,
	}

	var result RiskFactorsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetRiskCategories retrieves the taxonomy used to classify risk factors
// in SEC filing disclosures. Each category entry includes a three-level
// classification hierarchy with a description and taxonomy version number.
// Supports filtering by primary, secondary, and tertiary category values.
func (c *Client) GetRiskCategories(p RiskCategoriesParams) (*RiskCategoriesResponse, error) {
	path := "/stocks/taxonomies/vX/risk-factors"

	params := map[string]string{
		"primary_category":   p.PrimaryCategory,
		"secondary_category": p.SecondaryCategory,
		"tertiary_category":  p.TertiaryCategory,
		"taxonomy":           p.Taxonomy,
		"limit":              p.Limit,
		"sort":               p.Sort,
	}

	var result RiskCategoriesResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
