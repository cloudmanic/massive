//
// Date: 2026-02-14
// Copyright (c) 2026. All rights reserved.
//

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	configDir  = ".config/massive"
	configFile = "config.json"
)

// Config holds the application configuration including API credentials
// and the base URL for the Massive API.
type Config struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
}

// DefaultConfig returns a Config with default values. The base URL defaults
// to the Massive API production endpoint.
func DefaultConfig() *Config {
	return &Config{
		BaseURL: "https://api.massive.com",
	}
}

// configPath returns the full filesystem path to the config file
// located at ~/.config/massive/config.json.
func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, configDir, configFile), nil
}

// configDirPath returns the full filesystem path to the config directory
// located at ~/.config/massive/.
func configDirPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, configDir), nil
}

// Load reads the configuration from disk. If the config file does not exist,
// it returns a default configuration. Returns an error if the file exists
// but cannot be read or parsed.
func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

// Save writes the configuration to disk at ~/.config/massive/config.json.
// It creates the config directory if it does not exist. The file is written
// with 0600 permissions to protect the API key.
func Save(cfg *Config) error {
	dir, err := configDirPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetAPIKey returns the API key by checking the MASSIVE_API_KEY environment
// variable first, then falling back to the config file. Returns an error
// if no API key is found in either location.
func GetAPIKey() (string, error) {
	if key := os.Getenv("MASSIVE_API_KEY"); key != "" {
		return key, nil
	}

	cfg, err := Load()
	if err != nil {
		return "", err
	}

	if cfg.APIKey == "" {
		return "", fmt.Errorf("API key not configured. Run 'massive config init' or set MASSIVE_API_KEY environment variable")
	}

	return cfg.APIKey, nil
}
