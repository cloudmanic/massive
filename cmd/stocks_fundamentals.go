//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/cloudmanic/massive-cli/internal/api"
	"github.com/spf13/cobra"
)

// stocksFundamentalsCmd is the parent command for all stock fundamentals
// subcommands including short interest, short volume, float, balance
// sheets, income statements, cash flow statements, and financial ratios.
var stocksFundamentalsCmd = &cobra.Command{
	Use:   "fundamentals",
	Short: "Stock fundamentals data commands",
	Long:  "Access fundamental financial data for stocks including short interest, short volume, float, balance sheets, income statements, cash flow statements, and financial ratios.",
}

// ---------------------------------------------------------------------------
// Short Interest
// ---------------------------------------------------------------------------

// stocksShortInterestCmd retrieves bi-monthly aggregated short interest
// data reported to FINRA by broker-dealers for a specified stock ticker.
// Includes short interest counts, average daily volume, and estimated
// days to cover. Usage: massive stocks fundamentals short-interest --ticker AAPL
var stocksShortInterestCmd = &cobra.Command{
	Use:   "short-interest",
	Short: "Get bi-monthly short interest data from FINRA",
	Long:  "Retrieve bi-monthly aggregated short interest data reported to FINRA by broker-dealers for a specified stock ticker.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker, _ := cmd.Flags().GetString("ticker")
		settlementDate, _ := cmd.Flags().GetString("settlement-date")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.ShortInterestParams{
			Ticker:         strings.ToUpper(ticker),
			SettlementDate: settlementDate,
			Limit:          limit,
			Sort:           sort,
		}

		result, err := client.GetShortInterest(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Short Interest Results: %d\n\n", result.Count)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tSETTLEMENT DATE\tSHORT INTEREST\tAVG DAILY VOL\tDAYS TO COVER")
		fmt.Fprintln(w, "------\t---------------\t--------------\t-------------\t-------------")

		for _, si := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%.2f\n",
				si.Ticker, si.SettlementDate, si.ShortInterest,
				si.AvgDailyVolume, si.DaysToCover)
		}
		w.Flush()

		return nil
	},
}

// ---------------------------------------------------------------------------
// Short Volume
// ---------------------------------------------------------------------------

// stocksShortVolumeCmd retrieves daily aggregated short sale volume data
// reported to FINRA from off-exchange trading venues and alternative
// trading systems. Breaks down volume by exchange.
// Usage: massive stocks fundamentals short-volume --ticker AAPL
var stocksShortVolumeCmd = &cobra.Command{
	Use:   "short-volume",
	Short: "Get daily short sale volume data from FINRA",
	Long:  "Retrieve daily aggregated short sale volume data reported to FINRA from off-exchange trading venues and alternative trading systems.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker, _ := cmd.Flags().GetString("ticker")
		date, _ := cmd.Flags().GetString("date")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.ShortVolumeParams{
			Ticker: strings.ToUpper(ticker),
			Date:   date,
			Limit:  limit,
			Sort:   sort,
		}

		result, err := client.GetShortVolume(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Short Volume Results: %d\n\n", result.Count)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tDATE\tSHORT VOL\tTOTAL VOL\tRATIO\tEXEMPT\tNON-EXEMPT")
		fmt.Fprintln(w, "------\t----\t---------\t---------\t-----\t------\t----------")

		for _, sv := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%.2f%%\t%d\t%d\n",
				sv.Ticker, sv.Date, sv.ShortVolume, sv.TotalVolume,
				sv.ShortVolumeRatio, sv.ExemptVolume, sv.NonExemptVolume)
		}
		w.Flush()

		return nil
	},
}

// ---------------------------------------------------------------------------
// Float
// ---------------------------------------------------------------------------

// stocksFloatCmd retrieves the latest free float data for stock tickers.
// Free float represents shares outstanding available for public trading
// after excluding strategic holdings, insider positions, and restricted shares.
// Usage: massive stocks fundamentals float --ticker AAPL
var stocksFloatCmd = &cobra.Command{
	Use:   "float",
	Short: "Get free float data for stock tickers",
	Long:  "Retrieve the latest free float for a specified stock ticker. Free float represents shares available for public trading after excluding insider and restricted shares.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker, _ := cmd.Flags().GetString("ticker")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.FloatParams{
			Ticker: strings.ToUpper(ticker),
			Limit:  limit,
			Sort:   sort,
		}

		result, err := client.GetFloat(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Float Results: %d\n\n", len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tEFFECTIVE DATE\tFREE FLOAT\tFREE FLOAT %")
		fmt.Fprintln(w, "------\t--------------\t----------\t------------")

		for _, f := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t%d\t%.2f%%\n",
				f.Ticker, f.EffectiveDate, f.FreeFloat, f.FreeFloatPercent)
		}
		w.Flush()

		return nil
	},
}

// ---------------------------------------------------------------------------
// Balance Sheets
// ---------------------------------------------------------------------------

// stocksBalanceSheetsCmd retrieves comprehensive balance sheet data for
// public companies containing quarterly and annual financial positions.
// Includes assets, liabilities, and equity breakdowns sourced from SEC filings.
// Usage: massive stocks fundamentals balance-sheets --tickers AAPL --timeframe annual
var stocksBalanceSheetsCmd = &cobra.Command{
	Use:   "balance-sheets",
	Short: "Get balance sheet data for public companies",
	Long:  "Retrieve comprehensive balance sheet data for public companies, containing quarterly and annual financial positions including assets, liabilities, and equity.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		tickers, _ := cmd.Flags().GetString("tickers")
		cik, _ := cmd.Flags().GetString("cik")
		timeframe, _ := cmd.Flags().GetString("timeframe")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.BalanceSheetsParams{
			Tickers:   strings.ToUpper(tickers),
			CIK:       cik,
			Timeframe: timeframe,
			Limit:     limit,
			Sort:      sort,
		}

		result, err := client.GetBalanceSheets(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Balance Sheet Results: %d\n\n", len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKERS\tPERIOD END\tTIMEFRAME\tTOTAL ASSETS\tTOTAL LIABILITIES\tTOTAL EQUITY\tCASH")
		fmt.Fprintln(w, "-------\t----------\t---------\t------------\t-----------------\t------------\t----")

		for _, bs := range result.Results {
			tickerStr := strings.Join(bs.Tickers, ",")
			fmt.Fprintf(w, "%s\t%s\t%s\t$%.0f\t$%.0f\t$%.0f\t$%.0f\n",
				tickerStr, bs.PeriodEnd, bs.Timeframe,
				bs.TotalAssets, bs.TotalLiabilities, bs.TotalEquity,
				bs.CashAndEquivalents)
		}
		w.Flush()

		return nil
	},
}

// ---------------------------------------------------------------------------
// Income Statements
// ---------------------------------------------------------------------------

// stocksIncomeStatementsCmd retrieves comprehensive income statement data
// for public companies including key metrics such as revenue, expenses,
// and net income. Supports quarterly, annual, and TTM timeframes.
// Usage: massive stocks fundamentals income-statements --tickers AAPL --timeframe annual
var stocksIncomeStatementsCmd = &cobra.Command{
	Use:   "income-statements",
	Short: "Get income statement data for public companies",
	Long:  "Retrieve comprehensive income statement data for public companies, including key metrics such as revenue, expenses, and net income.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		tickers, _ := cmd.Flags().GetString("tickers")
		cik, _ := cmd.Flags().GetString("cik")
		timeframe, _ := cmd.Flags().GetString("timeframe")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.IncomeStatementsParams{
			Tickers:   strings.ToUpper(tickers),
			CIK:       cik,
			Timeframe: timeframe,
			Limit:     limit,
			Sort:      sort,
		}

		result, err := client.GetIncomeStatements(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Income Statement Results: %d\n\n", len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKERS\tPERIOD END\tTIMEFRAME\tREVENUE\tGROSS PROFIT\tOPERATING INCOME\tNET INCOME\tEPS")
		fmt.Fprintln(w, "-------\t----------\t---------\t-------\t------------\t----------------\t----------\t---")

		for _, is := range result.Results {
			tickerStr := strings.Join(is.Tickers, ",")
			fmt.Fprintf(w, "%s\t%s\t%s\t$%.0f\t$%.0f\t$%.0f\t$%.0f\t$%.2f\n",
				tickerStr, is.PeriodEnd, is.Timeframe,
				is.Revenue, is.GrossProfit, is.OperatingIncome,
				is.ConsolidatedNetIncomeLoss, is.DilutedEarningsPerShare)
		}
		w.Flush()

		return nil
	},
}

// ---------------------------------------------------------------------------
// Cash Flow Statements
// ---------------------------------------------------------------------------

// stocksCashFlowStatementsCmd retrieves comprehensive cash flow statement
// data for public companies including operating, investing, and financing
// activities. Supports quarterly, annual, and TTM timeframes.
// Usage: massive stocks fundamentals cash-flow-statements --tickers AAPL --timeframe annual
var stocksCashFlowStatementsCmd = &cobra.Command{
	Use:   "cash-flow-statements",
	Short: "Get cash flow statement data for public companies",
	Long:  "Retrieve comprehensive cash flow statement data for public companies, including quarterly, annual, and trailing twelve-month cash flows.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		tickers, _ := cmd.Flags().GetString("tickers")
		cik, _ := cmd.Flags().GetString("cik")
		timeframe, _ := cmd.Flags().GetString("timeframe")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.CashFlowStatementsParams{
			Tickers:   strings.ToUpper(tickers),
			CIK:       cik,
			Timeframe: timeframe,
			Limit:     limit,
			Sort:      sort,
		}

		result, err := client.GetCashFlowStatements(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Cash Flow Statement Results: %d\n\n", len(result.Results))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKERS\tPERIOD END\tTIMEFRAME\tOPERATING\tINVESTING\tFINANCING\tNET CHANGE")
		fmt.Fprintln(w, "-------\t----------\t---------\t---------\t---------\t---------\t----------")

		for _, cf := range result.Results {
			tickerStr := strings.Join(cf.Tickers, ",")
			fmt.Fprintf(w, "%s\t%s\t%s\t$%.0f\t$%.0f\t$%.0f\t$%.0f\n",
				tickerStr, cf.PeriodEnd, cf.Timeframe,
				cf.NetCashFromOperatingActivities,
				cf.NetCashFromInvestingActivities,
				cf.NetCashFromFinancingActivities,
				cf.ChangeInCashAndEquivalents)
		}
		w.Flush()

		return nil
	},
}

// ---------------------------------------------------------------------------
// Financial Ratios
// ---------------------------------------------------------------------------

// stocksRatiosCmd retrieves comprehensive financial ratios data providing
// key valuation, profitability, liquidity, and leverage metrics for public
// companies. Includes P/E, P/B, ROE, ROA, and other common ratios.
// Usage: massive stocks fundamentals ratios --ticker AAPL
var stocksRatiosCmd = &cobra.Command{
	Use:   "ratios",
	Short: "Get financial ratios for public companies",
	Long:  "Retrieve comprehensive financial ratios data providing key valuation, profitability, liquidity, and leverage metrics for public companies.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		ticker, _ := cmd.Flags().GetString("ticker")
		limit, _ := cmd.Flags().GetString("limit")
		sort, _ := cmd.Flags().GetString("sort")

		params := api.RatiosParams{
			Ticker: strings.ToUpper(ticker),
			Limit:  limit,
			Sort:   sort,
		}

		result, err := client.GetRatios(params)
		if err != nil {
			return err
		}

		if outputFormat == "json" {
			return printJSON(result)
		}

		fmt.Printf("Financial Ratios Results: %d\n\n", result.Count)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TICKER\tDATE\tPRICE\tMKT CAP\tP/E\tP/B\tP/S\tDIV YIELD\tROE\tROA\tD/E\tCURRENT")
		fmt.Fprintln(w, "------\t----\t-----\t-------\t---\t---\t---\t---------\t---\t---\t---\t-------")

		for _, r := range result.Results {
			fmt.Fprintf(w, "%s\t%s\t$%.2f\t$%.0f\t%.2f\t%.2f\t%.2f\t%.2f%%\t%.2f%%\t%.2f%%\t%.2f\t%.2f\n",
				r.Ticker, r.Date, r.Price, r.MarketCap,
				r.PriceToEarnings, r.PriceToBook, r.PriceToSales,
				r.DividendYield, r.ReturnOnEquity, r.ReturnOnAssets,
				r.DebtToEquity, r.Current)
		}
		w.Flush()

		return nil
	},
}

// ---------------------------------------------------------------------------
// init - register all fundamentals subcommands
// ---------------------------------------------------------------------------

// init registers the fundamentals parent command under stocks and adds
// all fundamentals subcommands with their respective flags.
func init() {
	// Register parent fundamentals command under stocks
	stocksCmd.AddCommand(stocksFundamentalsCmd)

	// Short Interest flags
	stocksShortInterestCmd.Flags().String("ticker", "", "Stock ticker symbol")
	stocksShortInterestCmd.Flags().String("settlement-date", "", "Settlement date (YYYY-MM-DD)")
	stocksShortInterestCmd.Flags().String("limit", "10", "Number of results to return (max 50000)")
	stocksShortInterestCmd.Flags().String("sort", "", "Sort order (e.g., settlement_date.desc)")
	stocksFundamentalsCmd.AddCommand(stocksShortInterestCmd)

	// Short Volume flags
	stocksShortVolumeCmd.Flags().String("ticker", "", "Stock ticker symbol")
	stocksShortVolumeCmd.Flags().String("date", "", "Date (YYYY-MM-DD)")
	stocksShortVolumeCmd.Flags().String("limit", "10", "Number of results to return (max 50000)")
	stocksShortVolumeCmd.Flags().String("sort", "", "Sort order (e.g., date.desc)")
	stocksFundamentalsCmd.AddCommand(stocksShortVolumeCmd)

	// Float flags
	stocksFloatCmd.Flags().String("ticker", "", "Stock ticker symbol")
	stocksFloatCmd.Flags().String("limit", "100", "Number of results to return (max 5000)")
	stocksFloatCmd.Flags().String("sort", "ticker.asc", "Sort order (e.g., ticker.desc)")
	stocksFundamentalsCmd.AddCommand(stocksFloatCmd)

	// Balance Sheets flags
	stocksBalanceSheetsCmd.Flags().String("tickers", "", "Stock ticker symbol(s)")
	stocksBalanceSheetsCmd.Flags().String("cik", "", "SEC CIK identifier")
	stocksBalanceSheetsCmd.Flags().String("timeframe", "", "Timeframe (quarterly, annual)")
	stocksBalanceSheetsCmd.Flags().String("limit", "100", "Number of results to return (max 50000)")
	stocksBalanceSheetsCmd.Flags().String("sort", "period_end.asc", "Sort order (e.g., period_end.desc)")
	stocksFundamentalsCmd.AddCommand(stocksBalanceSheetsCmd)

	// Income Statements flags
	stocksIncomeStatementsCmd.Flags().String("tickers", "", "Stock ticker symbol(s)")
	stocksIncomeStatementsCmd.Flags().String("cik", "", "SEC CIK identifier")
	stocksIncomeStatementsCmd.Flags().String("timeframe", "", "Timeframe (quarterly, annual, trailing_twelve_months)")
	stocksIncomeStatementsCmd.Flags().String("limit", "100", "Number of results to return (max 50000)")
	stocksIncomeStatementsCmd.Flags().String("sort", "period_end.asc", "Sort order (e.g., period_end.desc)")
	stocksFundamentalsCmd.AddCommand(stocksIncomeStatementsCmd)

	// Cash Flow Statements flags
	stocksCashFlowStatementsCmd.Flags().String("tickers", "", "Stock ticker symbol(s)")
	stocksCashFlowStatementsCmd.Flags().String("cik", "", "SEC CIK identifier")
	stocksCashFlowStatementsCmd.Flags().String("timeframe", "", "Timeframe (quarterly, annual, trailing_twelve_months)")
	stocksCashFlowStatementsCmd.Flags().String("limit", "100", "Number of results to return (max 50000)")
	stocksCashFlowStatementsCmd.Flags().String("sort", "period_end.asc", "Sort order (e.g., period_end.desc)")
	stocksFundamentalsCmd.AddCommand(stocksCashFlowStatementsCmd)

	// Financial Ratios flags
	stocksRatiosCmd.Flags().String("ticker", "", "Stock ticker symbol")
	stocksRatiosCmd.Flags().String("limit", "100", "Number of results to return (max 50000)")
	stocksRatiosCmd.Flags().String("sort", "", "Sort order (e.g., date.desc)")
	stocksFundamentalsCmd.AddCommand(stocksRatiosCmd)
}
