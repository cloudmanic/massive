//
// Date: 2026-02-14
// Copyright (c) 2026. All rights reserved.
//

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cloudmanic/massive-cli/internal/config"
	"github.com/spf13/cobra"
)

// configCmd is the parent command for all configuration-related subcommands.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Massive CLI configuration",
}

// configInitCmd initializes the CLI configuration by prompting for an API key.
// It first checks the MASSIVE_API_KEY environment variable and offers to use
// that value. The configuration is saved to ~/.config/massive/config.json.
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration with your API key",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			cfg = config.DefaultConfig()
		}

		reader := bufio.NewReader(os.Stdin)

		envKey := os.Getenv("MASSIVE_API_KEY")
		if envKey != "" {
			fmt.Printf("Found API key in environment variable. Use it? [Y/n]: ")
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))
			if answer == "" || answer == "y" || answer == "yes" {
				cfg.APIKey = envKey
			}
		}

		if cfg.APIKey == "" {
			fmt.Print("Enter your Massive API key: ")
			key, _ := reader.ReadString('\n')
			cfg.APIKey = strings.TrimSpace(key)
		}

		if cfg.APIKey == "" {
			return fmt.Errorf("API key cannot be empty")
		}

		fmt.Print("\nConfigure S3 flat file access? [y/N]: ")
		s3Answer, _ := reader.ReadString('\n')
		s3Answer = strings.TrimSpace(strings.ToLower(s3Answer))
		if s3Answer == "y" || s3Answer == "yes" {
			envS3Access := os.Getenv("MASSIVE_S3_ACCESS_KEY")
			if envS3Access != "" {
				fmt.Printf("Found S3 access key in environment variable. Use it? [Y/n]: ")
				answer, _ := reader.ReadString('\n')
				answer = strings.TrimSpace(strings.ToLower(answer))
				if answer == "" || answer == "y" || answer == "yes" {
					cfg.S3AccessKey = envS3Access
				}
			}
			if cfg.S3AccessKey == "" {
				fmt.Print("Enter your S3 Access Key ID: ")
				key, _ := reader.ReadString('\n')
				cfg.S3AccessKey = strings.TrimSpace(key)
			}

			envS3Secret := os.Getenv("MASSIVE_S3_SECRET_KEY")
			if envS3Secret != "" {
				fmt.Printf("Found S3 secret key in environment variable. Use it? [Y/n]: ")
				answer, _ := reader.ReadString('\n')
				answer = strings.TrimSpace(strings.ToLower(answer))
				if answer == "" || answer == "y" || answer == "yes" {
					cfg.S3SecretKey = envS3Secret
				}
			}
			if cfg.S3SecretKey == "" {
				fmt.Print("Enter your S3 Secret Access Key: ")
				key, _ := reader.ReadString('\n')
				cfg.S3SecretKey = strings.TrimSpace(key)
			}
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Println("Configuration saved to ~/.config/massive/config.json")
		return nil
	},
}

// configShowCmd displays the current configuration with the API key partially
// masked for security. Shows the base URL and masked API key.
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		maskedKey := maskString(cfg.APIKey)
		maskedS3Access := maskString(cfg.S3AccessKey)
		maskedS3Secret := maskString(cfg.S3SecretKey)

		fmt.Printf("Base URL:       %s\n", cfg.BaseURL)
		fmt.Printf("API Key:        %s\n", maskedKey)
		fmt.Printf("S3 Endpoint:    %s\n", cfg.S3Endpoint)
		fmt.Printf("S3 Access Key:  %s\n", maskedS3Access)
		fmt.Printf("S3 Secret Key:  %s\n", maskedS3Secret)

		return nil
	},
}

// init registers the config subcommands with the root command.
func init() {
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)
}
