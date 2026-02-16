//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/cloudmanic/massive-cli/internal/config"
	"github.com/cloudmanic/massive-cli/internal/flatfiles"
	"github.com/spf13/cobra"
)

// assetClassMap maps user-friendly asset class names to their S3 prefix constants.
// Users type short names like "stocks" and this map resolves them to the full
// S3 directory names like "us_stocks_sip".
var assetClassMap = map[string]string{
	"stocks":  flatfiles.AssetUSStocks,
	"options": flatfiles.AssetUSOptions,
	"indices": flatfiles.AssetUSIndices,
	"crypto":  flatfiles.AssetCrypto,
	"forex":   flatfiles.AssetForex,
}

// dataTypeMap maps user-friendly data type names to their S3 prefix constants.
// Users type short names like "trades" and this map resolves them to the full
// S3 directory names like "trades_v1".
var dataTypeMap = map[string]string{
	"trades":      flatfiles.DataTypeTrades,
	"quotes":      flatfiles.DataTypeQuotes,
	"day-aggs":    flatfiles.DataTypeDayAggs,
	"minute-aggs": flatfiles.DataTypeMinuteAggs,
}

// filesCmd is the parent command for all flat file (S3) data subcommands
// including listing available files, downloading data, and viewing asset
// classes and data types.
var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "Flat file (S3) data commands",
	Long:  "Commands for listing, browsing, and downloading flat file data from the Massive S3-compatible storage endpoint.",
}

// filesListCmd lists available flat files for a given asset class and data type.
// It requires a year flag and optionally accepts a month flag to narrow results.
// Output can be formatted as a table (default) or JSON via the --output flag.
// Usage: massive files list stocks trades --year 2024 --month 01
var filesListCmd = &cobra.Command{
	Use:   "list [asset] [datatype]",
	Short: "List available flat files for an asset class and data type",
	Long:  "Lists all available flat files (gzipped CSVs) for a given asset class and data type. Requires --year and optionally --month to filter results.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		s3Client, err := newS3Client()
		if err != nil {
			return err
		}

		assetClass, err := resolveAssetClass(args[0])
		if err != nil {
			return err
		}

		dataType, err := resolveDataType(args[1])
		if err != nil {
			return err
		}

		year, _ := cmd.Flags().GetString("year")
		month, _ := cmd.Flags().GetString("month")

		files, err := s3Client.ListFiles(assetClass, dataType, year, month)
		if err != nil {
			return fmt.Errorf("failed to list files: %w", err)
		}

		if len(files) == 0 {
			fmt.Println("No files found for the given criteria.")
			return nil
		}

		if outputFormat == "json" {
			return printJSON(files)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "FILE\tSIZE\tLAST MODIFIED")
		fmt.Fprintln(w, "----\t----\t-------------")

		for _, f := range files {
			fmt.Fprintf(w, "%s\t%s\t%s\n",
				f.Key,
				formatFileSize(f.Size),
				f.LastModified.Format("2006-01-02 15:04:05"),
			)
		}
		w.Flush()

		return nil
	},
}

// filesDownloadCmd downloads a specific flat file for a given asset class,
// data type, and date. The date must be in YYYY-MM-DD format. Files are
// saved to the current directory by default, or to the path specified by
// the --output-dir flag.
// Usage: massive files download stocks trades 2024-01-15 --output-dir ./data
var filesDownloadCmd = &cobra.Command{
	Use:   "download [asset] [datatype] [date]",
	Short: "Download a flat file for a specific date",
	Long:  "Downloads a specific flat file (gzipped CSV) for a given asset class, data type, and date (YYYY-MM-DD format).",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		s3Client, err := newS3Client()
		if err != nil {
			return err
		}

		assetClass, err := resolveAssetClass(args[0])
		if err != nil {
			return err
		}

		dataType, err := resolveDataType(args[1])
		if err != nil {
			return err
		}

		date := args[2]
		outputDir, _ := cmd.Flags().GetString("output-dir")

		key, err := flatfiles.BuildKey(assetClass, dataType, date)
		if err != nil {
			return fmt.Errorf("failed to build file key: %w", err)
		}

		filename := filepath.Base(key)
		destPath := filepath.Join(outputDir, filename)

		fmt.Printf("Downloading %s ...\n", key)

		if err := s3Client.DownloadFile(key, destPath); err != nil {
			return fmt.Errorf("failed to download file: %w", err)
		}

		fmt.Printf("Successfully downloaded to %s\n", destPath)
		return nil
	},
}

// filesAssetsCmd lists all available asset classes with their user-friendly
// names and corresponding S3 prefix values. Supports table and JSON output.
// Usage: massive files assets
var filesAssetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "List all available asset classes",
	Long:  "Displays all available asset classes with their friendly names and S3 prefix values.",
	RunE: func(cmd *cobra.Command, args []string) error {
		type assetEntry struct {
			Name   string `json:"name"`
			Prefix string `json:"prefix"`
		}

		entries := []assetEntry{
			{Name: "stocks", Prefix: flatfiles.AssetUSStocks},
			{Name: "options", Prefix: flatfiles.AssetUSOptions},
			{Name: "indices", Prefix: flatfiles.AssetUSIndices},
			{Name: "crypto", Prefix: flatfiles.AssetCrypto},
			{Name: "forex", Prefix: flatfiles.AssetForex},
		}

		if outputFormat == "json" {
			return printJSON(entries)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tS3 PREFIX")
		fmt.Fprintln(w, "----\t---------")

		for _, e := range entries {
			fmt.Fprintf(w, "%s\t%s\n", e.Name, e.Prefix)
		}
		w.Flush()

		return nil
	},
}

// filesTypesCmd lists all available data types with their user-friendly
// names and corresponding S3 prefix values. Supports table and JSON output.
// Usage: massive files types
var filesTypesCmd = &cobra.Command{
	Use:   "types",
	Short: "List all available data types",
	Long:  "Displays all available data types with their friendly names and S3 prefix values.",
	RunE: func(cmd *cobra.Command, args []string) error {
		type dataTypeEntry struct {
			Name   string `json:"name"`
			Prefix string `json:"prefix"`
		}

		entries := []dataTypeEntry{
			{Name: "trades", Prefix: flatfiles.DataTypeTrades},
			{Name: "quotes", Prefix: flatfiles.DataTypeQuotes},
			{Name: "day-aggs", Prefix: flatfiles.DataTypeDayAggs},
			{Name: "minute-aggs", Prefix: flatfiles.DataTypeMinuteAggs},
		}

		if outputFormat == "json" {
			return printJSON(entries)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tS3 PREFIX")
		fmt.Fprintln(w, "----\t---------")

		for _, e := range entries {
			fmt.Fprintf(w, "%s\t%s\n", e.Name, e.Prefix)
		}
		w.Flush()

		return nil
	},
}

// init registers the files parent command and all its subcommands with the
// root command. It also sets up command-specific flags for list and download.
func init() {
	// Register flags for the list subcommand
	filesListCmd.Flags().String("year", "", "Year to list files for (YYYY) [required]")
	filesListCmd.Flags().String("month", "", "Month to list files for (MM, optional)")
	filesListCmd.MarkFlagRequired("year")

	// Register flags for the download subcommand
	filesDownloadCmd.Flags().String("output-dir", ".", "Directory to save downloaded file")

	// Wire up subcommands under the files parent
	filesCmd.AddCommand(filesListCmd)
	filesCmd.AddCommand(filesDownloadCmd)
	filesCmd.AddCommand(filesAssetsCmd)
	filesCmd.AddCommand(filesTypesCmd)

	// Register files command under root
	rootCmd.AddCommand(filesCmd)
}

// newS3Client creates a new S3 client for accessing Massive flat files.
// It checks for MASSIVE_S3_ACCESS_KEY and MASSIVE_S3_SECRET_KEY environment
// variables first, then falls back to values stored in the config file.
// Returns an error if no S3 credentials are found in either location.
func newS3Client() (*flatfiles.S3Client, error) {
	accessKey := os.Getenv("MASSIVE_S3_ACCESS_KEY")
	secretKey := os.Getenv("MASSIVE_S3_SECRET_KEY")

	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Fall back to config file values if env vars are not set
	if accessKey == "" {
		accessKey = cfg.S3AccessKey
	}
	if secretKey == "" {
		secretKey = cfg.S3SecretKey
	}

	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("S3 credentials not configured. Set MASSIVE_S3_ACCESS_KEY and MASSIVE_S3_SECRET_KEY environment variables, or run 'massive config init' to save them")
	}

	endpoint := cfg.S3Endpoint
	if endpoint == "" {
		endpoint = "https://files.massive.com"
	}

	return flatfiles.NewS3Client(accessKey, secretKey, endpoint), nil
}

// resolveAssetClass converts a user-friendly asset class name (e.g., "stocks")
// to its corresponding S3 prefix constant (e.g., "us_stocks_sip"). Returns an
// error with a list of valid names if the input is not recognized.
func resolveAssetClass(name string) (string, error) {
	if v, ok := assetClassMap[name]; ok {
		return v, nil
	}
	return "", fmt.Errorf("unknown asset class %q. Valid values: stocks, options, indices, crypto, forex", name)
}

// resolveDataType converts a user-friendly data type name (e.g., "trades")
// to its corresponding S3 prefix constant (e.g., "trades_v1"). Returns an
// error with a list of valid names if the input is not recognized.
func resolveDataType(name string) (string, error) {
	if v, ok := dataTypeMap[name]; ok {
		return v, nil
	}
	return "", fmt.Errorf("unknown data type %q. Valid values: trades, quotes, day-aggs, minute-aggs", name)
}

// formatFileSize converts a file size in bytes to a human-readable string
// with appropriate units (B, KB, MB, GB). Uses 1024-based units for
// accurate binary size representation.
func formatFileSize(bytes int64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
	)

	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
