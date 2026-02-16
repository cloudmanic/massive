//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

// ---------------------------------------------------------------------------
// Short Interest
// ---------------------------------------------------------------------------

// ShortInterestResponse represents the API response for bi-monthly
// aggregated short interest data reported to FINRA by broker-dealers.
type ShortInterestResponse struct {
	Status    string          `json:"status"`
	RequestID string          `json:"request_id"`
	Count     int             `json:"count"`
	NextURL   string          `json:"next_url,omitempty"`
	Results   []ShortInterest `json:"results"`
}

// ShortInterest represents a single short interest record for a ticker
// on a specific settlement date with the estimated days to cover.
type ShortInterest struct {
	Ticker         string  `json:"ticker"`
	SettlementDate string  `json:"settlement_date"`
	ShortInterest  int64   `json:"short_interest"`
	AvgDailyVolume int64   `json:"avg_daily_volume"`
	DaysToCover    float64 `json:"days_to_cover"`
}

// ShortInterestParams holds the query parameters for fetching short
// interest data from the FINRA bi-monthly reports endpoint.
type ShortInterestParams struct {
	Ticker         string
	SettlementDate string
	Limit          string
	Sort           string
}

// GetShortInterest retrieves bi-monthly aggregated short interest data
// reported to FINRA by broker-dealers for the specified filters. The
// data includes short interest counts, average daily volume, and
// estimated days to cover all short positions.
func (c *Client) GetShortInterest(p ShortInterestParams) (*ShortInterestResponse, error) {
	path := "/stocks/v1/short-interest"

	params := map[string]string{
		"ticker":          p.Ticker,
		"settlement_date": p.SettlementDate,
		"limit":           p.Limit,
		"sort":            p.Sort,
	}

	var result ShortInterestResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ---------------------------------------------------------------------------
// Short Volume
// ---------------------------------------------------------------------------

// ShortVolumeResponse represents the API response for daily aggregated
// short sale volume data reported to FINRA from off-exchange venues.
type ShortVolumeResponse struct {
	Status    string        `json:"status"`
	RequestID string        `json:"request_id"`
	Count     int           `json:"count"`
	NextURL   string        `json:"next_url,omitempty"`
	Results   []ShortVolume `json:"results"`
}

// ShortVolume represents a single day's short volume data for a ticker
// broken down by exchange with exempt and non-exempt volumes.
type ShortVolume struct {
	Ticker                        string  `json:"ticker"`
	Date                          string  `json:"date"`
	TotalVolume                   int64   `json:"total_volume"`
	ShortVolume                   int64   `json:"short_volume"`
	ExemptVolume                  int64   `json:"exempt_volume"`
	NonExemptVolume               int64   `json:"non_exempt_volume"`
	ShortVolumeRatio              float64 `json:"short_volume_ratio"`
	NYSEShortVolume               int64   `json:"nyse_short_volume"`
	NYSEShortVolumeExempt         int64   `json:"nyse_short_volume_exempt"`
	NasdaqCarteretShortVolume     int64   `json:"nasdaq_carteret_short_volume"`
	NasdaqCarteretShortVolExempt  int64   `json:"nasdaq_carteret_short_volume_exempt"`
	NasdaqChicagoShortVolume      int64   `json:"nasdaq_chicago_short_volume"`
	NasdaqChicagoShortVolExempt   int64   `json:"nasdaq_chicago_short_volume_exempt"`
	ADFShortVolume                int64   `json:"adf_short_volume"`
	ADFShortVolumeExempt          int64   `json:"adf_short_volume_exempt"`
}

// ShortVolumeParams holds the query parameters for fetching daily
// aggregated short sale volume data from FINRA.
type ShortVolumeParams struct {
	Ticker string
	Date   string
	Limit  string
	Sort   string
}

// GetShortVolume retrieves daily aggregated short sale volume data
// reported to FINRA from off-exchange trading venues and alternative
// trading systems for the specified ticker and date filters.
func (c *Client) GetShortVolume(p ShortVolumeParams) (*ShortVolumeResponse, error) {
	path := "/stocks/v1/short-volume"

	params := map[string]string{
		"ticker": p.Ticker,
		"date":   p.Date,
		"limit":  p.Limit,
		"sort":   p.Sort,
	}

	var result ShortVolumeResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ---------------------------------------------------------------------------
// Float
// ---------------------------------------------------------------------------

// FloatResponse represents the API response for free float data
// showing shares available for public trading.
type FloatResponse struct {
	Status    string      `json:"status"`
	RequestID string      `json:"request_id"`
	NextURL   string      `json:"next_url,omitempty"`
	Results   []FloatData `json:"results"`
}

// FloatData represents the free float data for a single ticker including
// the total free float share count and its percentage of shares outstanding.
type FloatData struct {
	Ticker           string  `json:"ticker"`
	EffectiveDate    string  `json:"effective_date"`
	FreeFloat        int64   `json:"free_float"`
	FreeFloatPercent float64 `json:"free_float_percent"`
}

// FloatParams holds the query parameters for fetching free float data
// for stock tickers.
type FloatParams struct {
	Ticker string
	Limit  string
	Sort   string
}

// GetFloat retrieves the latest free float data for stock tickers. Free
// float represents shares outstanding that are considered available for
// public trading after excluding strategic holdings, insider positions,
// and restricted shares.
func (c *Client) GetFloat(p FloatParams) (*FloatResponse, error) {
	path := "/stocks/vX/float"

	params := map[string]string{
		"ticker": p.Ticker,
		"limit":  p.Limit,
		"sort":   p.Sort,
	}

	var result FloatResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ---------------------------------------------------------------------------
// Balance Sheets
// ---------------------------------------------------------------------------

// BalanceSheetsResponse represents the API response for balance sheet
// data containing quarterly and annual financial positions.
type BalanceSheetsResponse struct {
	Status    string         `json:"status"`
	RequestID string         `json:"request_id"`
	NextURL   string         `json:"next_url,omitempty"`
	Results   []BalanceSheet `json:"results"`
}

// BalanceSheet represents a single balance sheet filing with assets,
// liabilities, and equity data for a specific reporting period.
type BalanceSheet struct {
	CIK                                    string   `json:"cik"`
	Tickers                                []string `json:"tickers"`
	PeriodEnd                              string   `json:"period_end"`
	FilingDate                             string   `json:"filing_date"`
	FiscalYear                             int      `json:"fiscal_year"`
	FiscalQuarter                          int      `json:"fiscal_quarter"`
	Timeframe                              string   `json:"timeframe"`
	TotalAssets                            float64  `json:"total_assets"`
	TotalCurrentAssets                     float64  `json:"total_current_assets"`
	TotalLiabilities                       float64  `json:"total_liabilities"`
	TotalCurrentLiabilities                float64  `json:"total_current_liabilities"`
	TotalEquity                            float64  `json:"total_equity"`
	TotalEquityAttributableToParent        float64  `json:"total_equity_attributable_to_parent"`
	TotalLiabilitiesAndEquity              float64  `json:"total_liabilities_and_equity"`
	CashAndEquivalents                     float64  `json:"cash_and_equivalents"`
	ShortTermInvestments                   float64  `json:"short_term_investments"`
	Receivables                            float64  `json:"receivables"`
	Inventories                            float64  `json:"inventories"`
	OtherCurrentAssets                     float64  `json:"other_current_assets"`
	PropertyPlantEquipmentNet              float64  `json:"property_plant_equipment_net"`
	Goodwill                               float64  `json:"goodwill"`
	IntangibleAssetsNet                    float64  `json:"intangible_assets_net"`
	OtherAssets                            float64  `json:"other_assets"`
	AccountsPayable                        float64  `json:"accounts_payable"`
	AccruedAndOtherCurrentLiabilities      float64  `json:"accrued_and_other_current_liabilities"`
	DeferredRevenueCurrent                 float64  `json:"deferred_revenue_current"`
	DebtCurrent                            float64  `json:"debt_current"`
	LongTermDebtAndCapitalLeaseObligations float64  `json:"long_term_debt_and_capital_lease_obligations"`
	DeferredRevenueNoncurrent              float64  `json:"deferred_revenue_noncurrent"`
	OtherNoncurrentLiabilities             float64  `json:"other_noncurrent_liabilities"`
	CommitmentsAndContingencies            float64  `json:"commitments_and_contingencies"`
	CommonStock                            float64  `json:"common_stock"`
	PreferredStock                         float64  `json:"preferred_stock"`
	AdditionalPaidInCapital                float64  `json:"additional_paid_in_capital"`
	RetainedEarningsDeficit                float64  `json:"retained_earnings_deficit"`
	AccumulatedOtherComprehensiveIncome    float64  `json:"accumulated_other_comprehensive_income"`
	OtherEquity                            float64  `json:"other_equity"`
	TreasuryStock                          float64  `json:"treasury_stock"`
	NoncontrollingInterest                 float64  `json:"noncontrolling_interest"`
}

// BalanceSheetsParams holds the query parameters for fetching balance
// sheet data from the fundamentals endpoint.
type BalanceSheetsParams struct {
	Tickers   string
	CIK       string
	Timeframe string
	Limit     string
	Sort      string
}

// GetBalanceSheets retrieves comprehensive balance sheet data for public
// companies containing quarterly and annual financial positions. The data
// includes assets, liabilities, and equity breakdowns sourced from SEC
// filings.
func (c *Client) GetBalanceSheets(p BalanceSheetsParams) (*BalanceSheetsResponse, error) {
	path := "/stocks/financials/v1/balance-sheets"

	params := map[string]string{
		"tickers":   p.Tickers,
		"cik":       p.CIK,
		"timeframe": p.Timeframe,
		"limit":     p.Limit,
		"sort":      p.Sort,
	}

	var result BalanceSheetsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ---------------------------------------------------------------------------
// Income Statements
// ---------------------------------------------------------------------------

// IncomeStatementsResponse represents the API response for income
// statement data containing revenue, expenses, and net income metrics.
type IncomeStatementsResponse struct {
	Status    string            `json:"status"`
	RequestID string            `json:"request_id"`
	NextURL   string            `json:"next_url,omitempty"`
	Results   []IncomeStatement `json:"results"`
}

// IncomeStatement represents a single income statement filing with
// revenue, expense, and earnings data for a specific reporting period.
type IncomeStatement struct {
	CIK                                     string   `json:"cik"`
	Tickers                                 []string `json:"tickers"`
	PeriodEnd                               string   `json:"period_end"`
	FilingDate                              string   `json:"filing_date"`
	FiscalYear                              int      `json:"fiscal_year"`
	FiscalQuarter                           int      `json:"fiscal_quarter"`
	Timeframe                               string   `json:"timeframe"`
	Revenue                                 float64  `json:"revenue"`
	CostOfRevenue                           float64  `json:"cost_of_revenue"`
	GrossProfit                             float64  `json:"gross_profit"`
	TotalOperatingExpenses                  float64  `json:"total_operating_expenses"`
	OperatingIncome                         float64  `json:"operating_income"`
	InterestIncome                          float64  `json:"interest_income"`
	InterestExpense                         float64  `json:"interest_expense"`
	OtherIncomeExpense                      float64  `json:"other_income_expense"`
	IncomeBeforeIncomeTaxes                 float64  `json:"income_before_income_taxes"`
	IncomeTaxes                             float64  `json:"income_taxes"`
	ConsolidatedNetIncomeLoss               float64  `json:"consolidated_net_income_loss"`
	NetIncomeLossAttributableCommonShareholders float64 `json:"net_income_loss_attributable_common_shareholders"`
	BasicEarningsPerShare                   float64  `json:"basic_earnings_per_share"`
	DilutedEarningsPerShare                 float64  `json:"diluted_earnings_per_share"`
	BasicSharesOutstanding                  float64  `json:"basic_shares_outstanding"`
	DilutedSharesOutstanding                float64  `json:"diluted_shares_outstanding"`
	EBITDA                                  float64  `json:"ebitda"`
	DepreciationDepletionAmortization       float64  `json:"depreciation_depletion_amortization"`
	ResearchDevelopment                     float64  `json:"research_development"`
	SellingGeneralAdministrative            float64  `json:"selling_general_administrative"`
	OtherOperatingExpenses                  float64  `json:"other_operating_expenses"`
	DiscontinuedOperations                  float64  `json:"discontinued_operations"`
	ExtraordinaryItems                      float64  `json:"extraordinary_items"`
	EquityInAffiliates                      float64  `json:"equity_in_affiliates"`
	NoncontrollingInterest                  float64  `json:"noncontrolling_interest"`
	PreferredStockDividendsDeclared         float64  `json:"preferred_stock_dividends_declared"`
	TotalOtherIncomeExpense                 float64  `json:"total_other_income_expense"`
}

// IncomeStatementsParams holds the query parameters for fetching income
// statement data from the fundamentals endpoint.
type IncomeStatementsParams struct {
	Tickers   string
	CIK       string
	Timeframe string
	Limit     string
	Sort      string
}

// GetIncomeStatements retrieves comprehensive income statement data for
// public companies including key metrics such as revenue, expenses, and
// net income. Supports quarterly, annual, and trailing twelve-month
// timeframes sourced from SEC filings.
func (c *Client) GetIncomeStatements(p IncomeStatementsParams) (*IncomeStatementsResponse, error) {
	path := "/stocks/financials/v1/income-statements"

	params := map[string]string{
		"tickers":   p.Tickers,
		"cik":       p.CIK,
		"timeframe": p.Timeframe,
		"limit":     p.Limit,
		"sort":      p.Sort,
	}

	var result IncomeStatementsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ---------------------------------------------------------------------------
// Cash Flow Statements
// ---------------------------------------------------------------------------

// CashFlowStatementsResponse represents the API response for cash flow
// statement data containing operating, investing, and financing activities.
type CashFlowStatementsResponse struct {
	Status    string              `json:"status"`
	RequestID string              `json:"request_id"`
	NextURL   string              `json:"next_url,omitempty"`
	Results   []CashFlowStatement `json:"results"`
}

// CashFlowStatement represents a single cash flow statement filing with
// operating, investing, and financing cash flow breakdowns.
type CashFlowStatement struct {
	CIK                                                string   `json:"cik"`
	Tickers                                            []string `json:"tickers"`
	PeriodEnd                                          string   `json:"period_end"`
	FilingDate                                         string   `json:"filing_date"`
	FiscalYear                                         int      `json:"fiscal_year"`
	FiscalQuarter                                      int      `json:"fiscal_quarter"`
	Timeframe                                          string   `json:"timeframe"`
	NetCashFromOperatingActivities                     float64  `json:"net_cash_from_operating_activities"`
	CashFromOperatingActivitiesContinuingOperations    float64  `json:"cash_from_operating_activities_continuing_operations"`
	NetCashFromOperatingActivitiesDiscontinued          float64  `json:"net_cash_from_operating_activities_discontinued_operations"`
	NetCashFromInvestingActivities                     float64  `json:"net_cash_from_investing_activities"`
	NetCashFromInvestingActivitiesContinuingOperations  float64  `json:"net_cash_from_investing_activities_continuing_operations"`
	NetCashFromInvestingActivitiesDiscontinued          float64  `json:"net_cash_from_investing_activities_discontinued_operations"`
	NetCashFromFinancingActivities                     float64  `json:"net_cash_from_financing_activities"`
	NetCashFromFinancingActivitiesContinuingOperations  float64  `json:"net_cash_from_financing_activities_continuing_operations"`
	NetCashFromFinancingActivitiesDiscontinued          float64  `json:"net_cash_from_financing_activities_discontinued_operations"`
	ChangeInCashAndEquivalents                         float64  `json:"change_in_cash_and_equivalents"`
	NetIncome                                          float64  `json:"net_income"`
	DepreciationDepletionAndAmortization               float64  `json:"depreciation_depletion_and_amortization"`
	ChangeInOtherOperatingAssetsAndLiabilitiesNet      float64  `json:"change_in_other_operating_assets_and_liabilities_net"`
	OtherOperatingActivities                           float64  `json:"other_operating_activities"`
	PurchaseOfPropertyPlantAndEquipment                float64  `json:"purchase_of_property_plant_and_equipment"`
	SaleOfPropertyPlantAndEquipment                    float64  `json:"sale_of_property_plant_and_equipment"`
	OtherInvestingActivities                           float64  `json:"other_investing_activities"`
	ShortTermDebtIssuancesRepayments                   float64  `json:"short_term_debt_issuances_repayments"`
	LongTermDebtIssuancesRepayments                    float64  `json:"long_term_debt_issuances_repayments"`
	Dividends                                          float64  `json:"dividends"`
	OtherFinancingActivities                           float64  `json:"other_financing_activities"`
	EffectOfCurrencyExchangeRate                       float64  `json:"effect_of_currency_exchange_rate"`
	IncomeLossFromDiscontinuedOperations               float64  `json:"income_loss_from_discontinued_operations"`
	NoncontrollingInterests                            float64  `json:"noncontrolling_interests"`
	OtherCashAdjustments                               float64  `json:"other_cash_adjustments"`
}

// CashFlowStatementsParams holds the query parameters for fetching
// cash flow statement data from the fundamentals endpoint.
type CashFlowStatementsParams struct {
	Tickers   string
	CIK       string
	Timeframe string
	Limit     string
	Sort      string
}

// GetCashFlowStatements retrieves comprehensive cash flow statement data
// for public companies including quarterly, annual, and trailing twelve-
// month cash flows. The data covers operating, investing, and financing
// activities sourced from SEC filings.
func (c *Client) GetCashFlowStatements(p CashFlowStatementsParams) (*CashFlowStatementsResponse, error) {
	path := "/stocks/financials/v1/cash-flow-statements"

	params := map[string]string{
		"tickers":   p.Tickers,
		"cik":       p.CIK,
		"timeframe": p.Timeframe,
		"limit":     p.Limit,
		"sort":      p.Sort,
	}

	var result CashFlowStatementsResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ---------------------------------------------------------------------------
// Financial Ratios
// ---------------------------------------------------------------------------

// RatiosResponse represents the API response for financial ratios data
// providing key valuation, profitability, liquidity, and leverage metrics.
type RatiosResponse struct {
	Status    string  `json:"status"`
	RequestID string  `json:"request_id"`
	Count     int     `json:"count"`
	NextURL   string  `json:"next_url,omitempty"`
	Results   []Ratio `json:"results"`
}

// Ratio represents a single financial ratios record for a ticker
// including valuation, profitability, liquidity, and leverage metrics.
type Ratio struct {
	Ticker              string  `json:"ticker"`
	CIK                 string  `json:"cik"`
	Date                string  `json:"date"`
	Price               float64 `json:"price"`
	MarketCap           float64 `json:"market_cap"`
	EarningsPerShare    float64 `json:"earnings_per_share"`
	PriceToEarnings     float64 `json:"price_to_earnings"`
	PriceToBook         float64 `json:"price_to_book"`
	PriceToSales        float64 `json:"price_to_sales"`
	PriceToCashFlow     float64 `json:"price_to_cash_flow"`
	PriceToFreeCashFlow float64 `json:"price_to_free_cash_flow"`
	DividendYield       float64 `json:"dividend_yield"`
	ReturnOnAssets      float64 `json:"return_on_assets"`
	ReturnOnEquity      float64 `json:"return_on_equity"`
	DebtToEquity        float64 `json:"debt_to_equity"`
	Current             float64 `json:"current"`
	Quick               float64 `json:"quick"`
	Cash                float64 `json:"cash"`
	EVToSales           float64 `json:"ev_to_sales"`
	EVToEBITDA          float64 `json:"ev_to_ebitda"`
	EnterpriseValue     float64 `json:"enterprise_value"`
	FreeCashFlow        float64 `json:"free_cash_flow"`
	AverageVolume       float64 `json:"average_volume"`
}

// RatiosParams holds the query parameters for fetching financial
// ratios data from the fundamentals endpoint.
type RatiosParams struct {
	Ticker string
	Limit  string
	Sort   string
}

// GetRatios retrieves comprehensive financial ratios data providing key
// valuation, profitability, liquidity, and leverage metrics for public
// companies. Includes metrics such as P/E, P/B, ROE, ROA, debt-to-equity,
// and current/quick/cash ratios.
func (c *Client) GetRatios(p RatiosParams) (*RatiosResponse, error) {
	path := "/stocks/financials/v1/ratios"

	params := map[string]string{
		"ticker": p.Ticker,
		"limit":  p.Limit,
		"sort":   p.Sort,
	}

	var result RatiosResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
