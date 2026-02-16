//
// Date: 2026-02-14
// Copyright (c) 2026. All rights reserved.
//

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/cloudmanic/massive-cli/internal/api"
	"github.com/cloudmanic/massive-cli/internal/config"
)

// newClient creates a new Massive API client by loading the API key from
// the environment or config file. Returns an error if no API key is found.
func newClient() (*api.Client, error) {
	apiKey, err := config.GetAPIKey()
	if err != nil {
		return nil, err
	}
	return api.NewClient(apiKey), nil
}

// printJSON formats the given value as indented JSON and prints it to stdout.
// Used when the --output json flag is specified.
func printJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}
