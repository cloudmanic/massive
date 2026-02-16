//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package flatfiles

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Default S3 endpoint for Massive flat files.
const defaultS3Endpoint = "https://files.massive.com"

// Default bucket name where all flat files are stored.
const defaultBucket = "flatfiles"

// Asset prefix constants represent the top-level directory names in S3
// that organize data by asset class.
const (
	AssetUSStocks  = "us_stocks_sip"
	AssetUSOptions = "us_options_opra"
	AssetUSIndices = "us_indices"
	AssetCrypto    = "global_crypto"
	AssetForex     = "global_forex"
)

// Data type constants represent the subdirectories under each asset prefix
// that organize data by the kind of market data contained in the files.
const (
	DataTypeTrades     = "trades_v1"
	DataTypeQuotes     = "quotes_v1"
	DataTypeDayAggs    = "day_aggs_v1"
	DataTypeMinuteAggs = "minute_aggs_v1"
)

// ValidAssetClasses contains all supported asset class prefixes. This is used
// for input validation when building S3 key paths.
var ValidAssetClasses = []string{
	AssetUSStocks,
	AssetUSOptions,
	AssetUSIndices,
	AssetCrypto,
	AssetForex,
}

// ValidDataTypes contains all supported data type subdirectories. This is used
// for input validation when building S3 key paths.
var ValidDataTypes = []string{
	DataTypeTrades,
	DataTypeQuotes,
	DataTypeDayAggs,
	DataTypeMinuteAggs,
}

// FileInfo represents metadata about a single file stored in S3. It contains
// the full S3 object key, the file size in bytes, and the last modification timestamp.
type FileInfo struct {
	Key          string
	Size         int64
	LastModified time.Time
}

// S3Client wraps the AWS S3 service client to provide convenient methods for
// listing and downloading flat files from the Massive S3-compatible endpoint.
type S3Client struct {
	client *s3.Client
	bucket string
}

// NewS3Client creates a new S3Client configured to communicate with the Massive
// flat files endpoint. It sets up static credentials using the provided access key
// and secret key, and configures the client to use path-style addressing which is
// required for S3-compatible endpoints.
func NewS3Client(accessKey, secretKey, endpoint string) *S3Client {
	if endpoint == "" {
		endpoint = defaultS3Endpoint
	}

	client := s3.New(s3.Options{
		Region:       "us-east-1",
		BaseEndpoint: aws.String(endpoint),
		Credentials:  credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		UsePathStyle: true,
	})

	return &S3Client{
		client: client,
		bucket: defaultBucket,
	}
}

// ListFiles lists all available files in S3 that match the given asset class,
// data type, year, and month. It constructs the appropriate S3 key prefix and
// performs a ListObjectsV2 call to retrieve matching objects. Returns a slice
// of FileInfo structs containing metadata about each file found, or an error
// if the S3 request fails.
func (s *S3Client) ListFiles(assetClass, dataType, year, month string) ([]FileInfo, error) {
	if err := validateAssetClass(assetClass); err != nil {
		return nil, err
	}

	if err := validateDataType(dataType); err != nil {
		return nil, err
	}

	prefix := BuildPrefix(assetClass, dataType, year, month)

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	}

	result, err := s.client.ListObjectsV2(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects with prefix %s: %w", prefix, err)
	}

	var files []FileInfo
	for _, obj := range result.Contents {
		files = append(files, FileInfo{
			Key:          aws.ToString(obj.Key),
			Size:         aws.ToInt64(obj.Size),
			LastModified: aws.ToTime(obj.LastModified),
		})
	}

	return files, nil
}

// DownloadFile downloads a single file from S3 by its full object key and writes
// it to the specified destination path on the local filesystem. It creates the
// destination file with standard permissions (0644). Returns an error if the S3
// request fails, the file cannot be created, or the data cannot be written.
func (s *S3Client) DownloadFile(key, destPath string) error {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	result, err := s.client.GetObject(context.Background(), input)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", key, err)
	}
	defer result.Body.Close()

	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", destPath, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, result.Body); err != nil {
		return fmt.Errorf("failed to write file %s: %w", destPath, err)
	}

	return nil
}

// BuildPrefix constructs the S3 key prefix used to list files for a given asset
// class, data type, year, and month. The resulting prefix follows the pattern
// "{assetClass}/{dataType}/{year}/{month}/". If year or month are empty, the
// prefix is truncated at that level, allowing broader listing queries.
func BuildPrefix(assetClass, dataType, year, month string) string {
	prefix := assetClass + "/" + dataType + "/"

	if year != "" {
		prefix += year + "/"
	}

	if year != "" && month != "" {
		prefix += month + "/"
	}

	return prefix
}

// BuildKey constructs the full S3 object key for a specific date's data file.
// It parses the date string (expected format YYYY-MM-DD) to extract the year
// and month components, then builds the full path following the pattern
// "{assetClass}/{dataType}/{year}/{month}/{date}.csv.gz". Returns an error
// if the date format is invalid.
func BuildKey(assetClass, dataType, date string) (string, error) {
	if err := validateAssetClass(assetClass); err != nil {
		return "", err
	}

	if err := validateDataType(dataType); err != nil {
		return "", err
	}

	parts := strings.Split(date, "-")
	if len(parts) != 3 || len(parts[0]) != 4 || len(parts[1]) != 2 || len(parts[2]) != 2 {
		return "", fmt.Errorf("invalid date format %q, expected YYYY-MM-DD", date)
	}

	year := parts[0]
	month := parts[1]

	key := fmt.Sprintf("%s/%s/%s/%s/%s.csv.gz", assetClass, dataType, year, month, date)

	return key, nil
}

// validateAssetClass checks whether the provided asset class string matches one
// of the known valid asset class prefixes. Returns an error with the list of
// valid options if the value is not recognized.
func validateAssetClass(assetClass string) error {
	for _, valid := range ValidAssetClasses {
		if assetClass == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid asset class %q, must be one of: %s", assetClass, strings.Join(ValidAssetClasses, ", "))
}

// validateDataType checks whether the provided data type string matches one of
// the known valid data type subdirectories. Returns an error with the list of
// valid options if the value is not recognized.
func validateDataType(dataType string) error {
	for _, valid := range ValidDataTypes {
		if dataType == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid data type %q, must be one of: %s", dataType, strings.Join(ValidDataTypes, ", "))
}
