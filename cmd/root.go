//
// Date: 2026-02-14
// Copyright (c) 2026. All rights reserved.
//

package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var outputFormat string

// rootCmd is the base command for the Massive CLI. All subcommands
// are registered as children of this command.
var rootCmd = &cobra.Command{
	Use:   "massive",
	Short: "CLI for the Massive financial data API",
	Long:  "A command-line interface for interacting with the Massive API to access stocks, crypto, forex, and other financial data.",
}

// Execute runs the root command and exits with a non-zero status code
// if any error occurs during command execution.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// init registers persistent flags and loads environment variables from
// the .env file if present. The output flag controls whether results
// are displayed as a table or raw JSON.
func init() {
	cobra.OnInitialize(loadEnv)
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, json)")
}

// loadEnv attempts to load environment variables from a .env file in
// the current working directory. Errors are silently ignored since the
// .env file is optional.
func loadEnv() {
	_ = godotenv.Load()
}
