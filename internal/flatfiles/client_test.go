//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package flatfiles

import (
	"compress/gzip"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestNewS3Client verifies that NewS3Client creates a client with the
// correct bucket name and that the underlying S3 client is initialized.
func TestNewS3Client(t *testing.T) {
	client := NewS3Client("access-key", "secret-key", "https://files.massive.com")

	if client.client == nil {
		t.Error("expected S3 client to be initialized")
	}

	if client.bucket != defaultBucket {
		t.Errorf("expected bucket %s, got %s", defaultBucket, client.bucket)
	}
}

// TestNewS3ClientDefaultEndpoint verifies that when an empty endpoint is
// provided, the client falls back to the default S3 endpoint.
func TestNewS3ClientDefaultEndpoint(t *testing.T) {
	client := NewS3Client("access-key", "secret-key", "")

	if client.client == nil {
		t.Error("expected S3 client to be initialized")
	}

	if client.bucket != defaultBucket {
		t.Errorf("expected bucket %s, got %s", defaultBucket, client.bucket)
	}
}

// TestBuildPrefix verifies that BuildPrefix correctly constructs S3 key prefixes
// for various combinations of asset class, data type, year, and month.
func TestBuildPrefix(t *testing.T) {
	tests := []struct {
		name       string
		assetClass string
		dataType   string
		year       string
		month      string
		expected   string
	}{
		{
			name:       "full prefix with all fields",
			assetClass: AssetUSStocks,
			dataType:   DataTypeTrades,
			year:       "2024",
			month:      "03",
			expected:   "us_stocks_sip/trades_v1/2024/03/",
		},
		{
			name:       "prefix with year only",
			assetClass: AssetUSOptions,
			dataType:   DataTypeQuotes,
			year:       "2024",
			month:      "",
			expected:   "us_options_opra/quotes_v1/2024/",
		},
		{
			name:       "prefix without year or month",
			assetClass: AssetCrypto,
			dataType:   DataTypeDayAggs,
			year:       "",
			month:      "",
			expected:   "global_crypto/day_aggs_v1/",
		},
		{
			name:       "month ignored when year is empty",
			assetClass: AssetForex,
			dataType:   DataTypeMinuteAggs,
			year:       "",
			month:      "06",
			expected:   "global_forex/minute_aggs_v1/",
		},
		{
			name:       "indices with day aggregates",
			assetClass: AssetUSIndices,
			dataType:   DataTypeDayAggs,
			year:       "2025",
			month:      "01",
			expected:   "us_indices/day_aggs_v1/2025/01/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildPrefix(tt.assetClass, tt.dataType, tt.year, tt.month)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestBuildKey verifies that BuildKey correctly constructs the full S3 object key
// from an asset class, data type, and date string in YYYY-MM-DD format.
func TestBuildKey(t *testing.T) {
	tests := []struct {
		name       string
		assetClass string
		dataType   string
		date       string
		expected   string
		wantErr    bool
	}{
		{
			name:       "valid stock trades key",
			assetClass: AssetUSStocks,
			dataType:   DataTypeTrades,
			date:       "2024-03-07",
			expected:   "us_stocks_sip/trades_v1/2024/03/2024-03-07.csv.gz",
			wantErr:    false,
		},
		{
			name:       "valid stock day aggregates key",
			assetClass: AssetUSStocks,
			dataType:   DataTypeDayAggs,
			date:       "2024-03-04",
			expected:   "us_stocks_sip/day_aggs_v1/2024/03/2024-03-04.csv.gz",
			wantErr:    false,
		},
		{
			name:       "valid crypto trades key",
			assetClass: AssetCrypto,
			dataType:   DataTypeTrades,
			date:       "2024-02-15",
			expected:   "global_crypto/trades_v1/2024/02/2024-02-15.csv.gz",
			wantErr:    false,
		},
		{
			name:       "valid options quotes key",
			assetClass: AssetUSOptions,
			dataType:   DataTypeQuotes,
			date:       "2025-12-31",
			expected:   "us_options_opra/quotes_v1/2025/12/2025-12-31.csv.gz",
			wantErr:    false,
		},
		{
			name:       "invalid date format - wrong separator",
			assetClass: AssetUSStocks,
			dataType:   DataTypeTrades,
			date:       "2024/03/07",
			wantErr:    true,
		},
		{
			name:       "invalid date format - too short",
			assetClass: AssetUSStocks,
			dataType:   DataTypeTrades,
			date:       "2024-3-7",
			wantErr:    true,
		},
		{
			name:       "invalid date format - empty",
			assetClass: AssetUSStocks,
			dataType:   DataTypeTrades,
			date:       "",
			wantErr:    true,
		},
		{
			name:       "invalid asset class",
			assetClass: "invalid_asset",
			dataType:   DataTypeTrades,
			date:       "2024-03-07",
			wantErr:    true,
		},
		{
			name:       "invalid data type",
			assetClass: AssetUSStocks,
			dataType:   "invalid_type",
			date:       "2024-03-07",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BuildKey(tt.assetClass, tt.dataType, tt.date)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestValidateAssetClass verifies that validateAssetClass accepts all known
// asset class constants and rejects unknown values.
func TestValidateAssetClass(t *testing.T) {
	for _, valid := range ValidAssetClasses {
		if err := validateAssetClass(valid); err != nil {
			t.Errorf("expected %q to be valid, got error: %v", valid, err)
		}
	}

	err := validateAssetClass("not_real")
	if err == nil {
		t.Error("expected error for invalid asset class, got nil")
	}

	if !strings.Contains(err.Error(), "not_real") {
		t.Errorf("error should mention the invalid value, got: %s", err.Error())
	}
}

// TestValidateDataType verifies that validateDataType accepts all known
// data type constants and rejects unknown values.
func TestValidateDataType(t *testing.T) {
	for _, valid := range ValidDataTypes {
		if err := validateDataType(valid); err != nil {
			t.Errorf("expected %q to be valid, got error: %v", valid, err)
		}
	}

	err := validateDataType("invalid_type")
	if err == nil {
		t.Error("expected error for invalid data type, got nil")
	}

	if !strings.Contains(err.Error(), "invalid_type") {
		t.Errorf("error should mention the invalid value, got: %s", err.Error())
	}
}

// ListBucketResult represents the XML response structure returned by the S3
// ListObjectsV2 API. Used to construct mock responses in tests.
type ListBucketResult struct {
	XMLName  xml.Name       `xml:"ListBucketResult"`
	XMLNS    string         `xml:"xmlns,attr"`
	Contents []S3ObjectMock `xml:"Contents"`
}

// S3ObjectMock represents a single object entry in the S3 ListObjectsV2 XML
// response. Used to build mock responses in tests.
type S3ObjectMock struct {
	Key          string `xml:"Key"`
	Size         int64  `xml:"Size"`
	LastModified string `xml:"LastModified"`
}

// TestListFilesWithMockServer verifies that ListFiles correctly parses the S3
// ListObjectsV2 XML response and returns the expected FileInfo structs. It sets
// up a mock HTTP server that returns a valid XML response mimicking the S3 API.
func TestListFilesWithMockServer(t *testing.T) {
	mockResponse := ListBucketResult{
		XMLNS: "http://s3.amazonaws.com/doc/2006-03-01/",
		Contents: []S3ObjectMock{
			{
				Key:          "us_stocks_sip/trades_v1/2024/03/2024-03-04.csv.gz",
				Size:         1048576,
				LastModified: "2024-03-05T00:00:00.000Z",
			},
			{
				Key:          "us_stocks_sip/trades_v1/2024/03/2024-03-05.csv.gz",
				Size:         2097152,
				LastModified: "2024-03-06T00:00:00.000Z",
			},
		},
	}

	xmlBytes, err := xml.MarshalIndent(mockResponse, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mock XML: %v", err)
	}

	xmlStr := xml.Header + string(xmlBytes)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("list-type") != "2" {
			t.Errorf("expected list-type=2, got %s", r.URL.Query().Get("list-type"))
		}

		prefix := r.URL.Query().Get("prefix")
		if prefix != "us_stocks_sip/trades_v1/2024/03/" {
			t.Errorf("expected prefix us_stocks_sip/trades_v1/2024/03/, got %s", prefix)
		}

		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(xmlStr))
	}))
	defer server.Close()

	client := NewS3Client("test-access", "test-secret", server.URL)

	files, err := client.ListFiles(AssetUSStocks, DataTypeTrades, "2024", "03")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}

	if files[0].Key != "us_stocks_sip/trades_v1/2024/03/2024-03-04.csv.gz" {
		t.Errorf("expected key us_stocks_sip/trades_v1/2024/03/2024-03-04.csv.gz, got %s", files[0].Key)
	}

	if files[0].Size != 1048576 {
		t.Errorf("expected size 1048576, got %d", files[0].Size)
	}

	if files[1].Key != "us_stocks_sip/trades_v1/2024/03/2024-03-05.csv.gz" {
		t.Errorf("expected key us_stocks_sip/trades_v1/2024/03/2024-03-05.csv.gz, got %s", files[1].Key)
	}

	if files[1].Size != 2097152 {
		t.Errorf("expected size 2097152, got %d", files[1].Size)
	}
}

// TestListFilesEmptyResult verifies that ListFiles returns an empty slice
// (not an error) when the S3 bucket has no objects matching the prefix.
func TestListFilesEmptyResult(t *testing.T) {
	mockResponse := ListBucketResult{
		XMLNS:    "http://s3.amazonaws.com/doc/2006-03-01/",
		Contents: []S3ObjectMock{},
	}

	xmlBytes, err := xml.MarshalIndent(mockResponse, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mock XML: %v", err)
	}

	xmlStr := xml.Header + string(xmlBytes)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(xmlStr))
	}))
	defer server.Close()

	client := NewS3Client("test-access", "test-secret", server.URL)

	files, err := client.ListFiles(AssetUSIndices, DataTypeDayAggs, "2024", "01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("expected 0 files, got %d", len(files))
	}
}

// TestListFilesInvalidAssetClass verifies that ListFiles returns a validation
// error when an invalid asset class is provided.
func TestListFilesInvalidAssetClass(t *testing.T) {
	client := NewS3Client("test-access", "test-secret", "https://localhost:1")

	_, err := client.ListFiles("bad_asset", DataTypeTrades, "2024", "03")
	if err == nil {
		t.Fatal("expected error for invalid asset class, got nil")
	}

	if !strings.Contains(err.Error(), "invalid asset class") {
		t.Errorf("expected 'invalid asset class' in error, got: %s", err.Error())
	}
}

// TestListFilesInvalidDataType verifies that ListFiles returns a validation
// error when an invalid data type is provided.
func TestListFilesInvalidDataType(t *testing.T) {
	client := NewS3Client("test-access", "test-secret", "https://localhost:1")

	_, err := client.ListFiles(AssetUSStocks, "bad_type", "2024", "03")
	if err == nil {
		t.Fatal("expected error for invalid data type, got nil")
	}

	if !strings.Contains(err.Error(), "invalid data type") {
		t.Errorf("expected 'invalid data type' in error, got: %s", err.Error())
	}
}

// TestDownloadFileWithMockServer verifies that DownloadFile correctly downloads
// content from S3 and writes it to the specified local file path. It uses a mock
// HTTP server that serves gzipped CSV content.
func TestDownloadFileWithMockServer(t *testing.T) {
	csvContent := "ticker,open,high,low,close,volume\nAAPL,150.00,155.00,149.00,154.50,1000000\n"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/flatfiles/us_stocks_sip/day_aggs_v1/2024/03/2024-03-04.csv.gz"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Encoding", "gzip")

		gzWriter := gzip.NewWriter(w)
		gzWriter.Write([]byte(csvContent))
		gzWriter.Close()
	}))
	defer server.Close()

	client := NewS3Client("test-access", "test-secret", server.URL)

	tmpDir := t.TempDir()
	destPath := filepath.Join(tmpDir, "2024-03-04.csv.gz")

	err := client.DownloadFile("us_stocks_sip/day_aggs_v1/2024/03/2024-03-04.csv.gz", destPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	file, err := os.Open(destPath)
	if err != nil {
		t.Fatalf("failed to open downloaded file: %v", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		t.Fatalf("failed to create gzip reader: %v", err)
	}
	defer gzReader.Close()

	data, err := io.ReadAll(gzReader)
	if err != nil {
		t.Fatalf("failed to read decompressed data: %v", err)
	}

	if string(data) != csvContent {
		t.Errorf("expected content %q, got %q", csvContent, string(data))
	}
}

// TestDownloadFileInvalidDestPath verifies that DownloadFile returns an error
// when the destination path is invalid and cannot be created.
func TestDownloadFileInvalidDestPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write([]byte("test data"))
	}))
	defer server.Close()

	client := NewS3Client("test-access", "test-secret", server.URL)

	err := client.DownloadFile("us_stocks_sip/trades_v1/2024/03/2024-03-04.csv.gz", "/nonexistent/dir/file.csv.gz")
	if err == nil {
		t.Fatal("expected error for invalid dest path, got nil")
	}

	if !strings.Contains(err.Error(), "failed to create file") {
		t.Errorf("expected 'failed to create file' in error, got: %s", err.Error())
	}
}

// TestAssetClassConstants verifies that all asset class constants have the
// expected string values matching the S3 directory structure.
func TestAssetClassConstants(t *testing.T) {
	if AssetUSStocks != "us_stocks_sip" {
		t.Errorf("expected us_stocks_sip, got %s", AssetUSStocks)
	}
	if AssetUSOptions != "us_options_opra" {
		t.Errorf("expected us_options_opra, got %s", AssetUSOptions)
	}
	if AssetUSIndices != "us_indices" {
		t.Errorf("expected us_indices, got %s", AssetUSIndices)
	}
	if AssetCrypto != "global_crypto" {
		t.Errorf("expected global_crypto, got %s", AssetCrypto)
	}
	if AssetForex != "global_forex" {
		t.Errorf("expected global_forex, got %s", AssetForex)
	}
}

// TestDataTypeConstants verifies that all data type constants have the
// expected string values matching the S3 directory structure.
func TestDataTypeConstants(t *testing.T) {
	if DataTypeTrades != "trades_v1" {
		t.Errorf("expected trades_v1, got %s", DataTypeTrades)
	}
	if DataTypeQuotes != "quotes_v1" {
		t.Errorf("expected quotes_v1, got %s", DataTypeQuotes)
	}
	if DataTypeDayAggs != "day_aggs_v1" {
		t.Errorf("expected day_aggs_v1, got %s", DataTypeDayAggs)
	}
	if DataTypeMinuteAggs != "minute_aggs_v1" {
		t.Errorf("expected minute_aggs_v1, got %s", DataTypeMinuteAggs)
	}
}

// TestValidAssetClassesSlice verifies that the ValidAssetClasses slice contains
// exactly the five expected asset class prefixes.
func TestValidAssetClassesSlice(t *testing.T) {
	if len(ValidAssetClasses) != 5 {
		t.Fatalf("expected 5 valid asset classes, got %d", len(ValidAssetClasses))
	}

	expected := map[string]bool{
		"us_stocks_sip":   true,
		"us_options_opra": true,
		"us_indices":      true,
		"global_crypto":   true,
		"global_forex":    true,
	}

	for _, ac := range ValidAssetClasses {
		if !expected[ac] {
			t.Errorf("unexpected asset class in ValidAssetClasses: %s", ac)
		}
	}
}

// TestValidDataTypesSlice verifies that the ValidDataTypes slice contains
// exactly the four expected data type subdirectories.
func TestValidDataTypesSlice(t *testing.T) {
	if len(ValidDataTypes) != 4 {
		t.Fatalf("expected 4 valid data types, got %d", len(ValidDataTypes))
	}

	expected := map[string]bool{
		"trades_v1":     true,
		"quotes_v1":     true,
		"day_aggs_v1":   true,
		"minute_aggs_v1": true,
	}

	for _, dt := range ValidDataTypes {
		if !expected[dt] {
			t.Errorf("unexpected data type in ValidDataTypes: %s", dt)
		}
	}
}

// TestListFilesLastModifiedParsing verifies that the LastModified timestamp
// from the S3 response is correctly parsed into a time.Time value.
func TestListFilesLastModifiedParsing(t *testing.T) {
	mockResponse := ListBucketResult{
		XMLNS: "http://s3.amazonaws.com/doc/2006-03-01/",
		Contents: []S3ObjectMock{
			{
				Key:          "us_stocks_sip/day_aggs_v1/2024/03/2024-03-04.csv.gz",
				Size:         500000,
				LastModified: "2024-03-05T12:30:45.000Z",
			},
		},
	}

	xmlBytes, err := xml.MarshalIndent(mockResponse, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mock XML: %v", err)
	}

	xmlStr := xml.Header + string(xmlBytes)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(xmlStr))
	}))
	defer server.Close()

	client := NewS3Client("test-access", "test-secret", server.URL)

	files, err := client.ListFiles(AssetUSStocks, DataTypeDayAggs, "2024", "03")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	expectedTime, _ := time.Parse(time.RFC3339, "2024-03-05T12:30:45Z")
	if !files[0].LastModified.Equal(expectedTime) {
		t.Errorf("expected LastModified %v, got %v", expectedTime, files[0].LastModified)
	}
}
