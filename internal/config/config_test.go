//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package config

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestDir creates a temp directory and sets the config override
// so tests don't touch the real config. Returns a cleanup function.
func setupTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	SetConfigDir(dir)
	t.Cleanup(func() { SetConfigDir("") })
	return dir
}

// TestDefaultConfig verifies that DefaultConfig returns the expected
// default base URL and an empty API key.
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.BaseURL != "https://api.massive.com" {
		t.Errorf("expected base URL https://api.massive.com, got %s", cfg.BaseURL)
	}

	if cfg.APIKey != "" {
		t.Errorf("expected empty API key, got %s", cfg.APIKey)
	}
}

// TestLoadNoConfigFile verifies that Load returns a default config
// when no config file exists on disk.
func TestLoadNoConfigFile(t *testing.T) {
	setupTestDir(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.BaseURL != "https://api.massive.com" {
		t.Errorf("expected default base URL, got %s", cfg.BaseURL)
	}

	if cfg.APIKey != "" {
		t.Errorf("expected empty API key, got %s", cfg.APIKey)
	}
}

// TestSaveAndLoad verifies that saving a config and loading it back
// produces identical values.
func TestSaveAndLoad(t *testing.T) {
	setupTestDir(t)

	original := &Config{
		APIKey:  "test-api-key-12345",
		BaseURL: "https://api.massive.com",
	}

	if err := Save(original); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loaded.APIKey != original.APIKey {
		t.Errorf("expected API key %s, got %s", original.APIKey, loaded.APIKey)
	}

	if loaded.BaseURL != original.BaseURL {
		t.Errorf("expected base URL %s, got %s", original.BaseURL, loaded.BaseURL)
	}
}

// TestSaveCreatesDirectory verifies that Save creates the config
// directory if it does not already exist.
func TestSaveCreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	nestedDir := filepath.Join(dir, "nested", "config")
	SetConfigDir(nestedDir)
	t.Cleanup(func() { SetConfigDir("") })

	cfg := &Config{
		APIKey:  "test-key",
		BaseURL: "https://api.massive.com",
	}

	if err := Save(cfg); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	if _, err := os.Stat(filepath.Join(nestedDir, configFile)); os.IsNotExist(err) {
		t.Error("expected config file to be created")
	}
}

// TestSaveFilePermissions verifies that the config file is written
// with 0600 permissions to protect the API key.
func TestSaveFilePermissions(t *testing.T) {
	setupTestDir(t)

	cfg := &Config{
		APIKey:  "secret-key",
		BaseURL: "https://api.massive.com",
	}

	if err := Save(cfg); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	dir, _ := configDirPath()
	info, err := os.Stat(filepath.Join(dir, configFile))
	if err != nil {
		t.Fatalf("failed to stat config file: %v", err)
	}

	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("expected file permissions 0600, got %04o", perm)
	}
}

// TestLoadInvalidJSON verifies that Load returns an error when the
// config file contains invalid JSON.
func TestLoadInvalidJSON(t *testing.T) {
	dir := setupTestDir(t)

	if err := os.WriteFile(filepath.Join(dir, configFile), []byte("not json"), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := Load()
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

// TestGetAPIKeyFromEnv verifies that GetAPIKey returns the value from
// the MASSIVE_API_KEY environment variable when it is set.
func TestGetAPIKeyFromEnv(t *testing.T) {
	setupTestDir(t)

	t.Setenv("MASSIVE_API_KEY", "env-test-key")

	key, err := GetAPIKey()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if key != "env-test-key" {
		t.Errorf("expected env-test-key, got %s", key)
	}
}

// TestGetAPIKeyFromConfig verifies that GetAPIKey falls back to the
// config file when the environment variable is not set.
func TestGetAPIKeyFromConfig(t *testing.T) {
	setupTestDir(t)

	t.Setenv("MASSIVE_API_KEY", "")

	cfg := &Config{
		APIKey:  "config-test-key",
		BaseURL: "https://api.massive.com",
	}
	if err := Save(cfg); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	key, err := GetAPIKey()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if key != "config-test-key" {
		t.Errorf("expected config-test-key, got %s", key)
	}
}

// TestGetAPIKeyEnvTakesPrecedence verifies that the environment variable
// takes priority over a config file API key.
func TestGetAPIKeyEnvTakesPrecedence(t *testing.T) {
	setupTestDir(t)

	cfg := &Config{
		APIKey:  "config-key",
		BaseURL: "https://api.massive.com",
	}
	if err := Save(cfg); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	t.Setenv("MASSIVE_API_KEY", "env-key")

	key, err := GetAPIKey()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if key != "env-key" {
		t.Errorf("expected env-key, got %s", key)
	}
}

// TestGetAPIKeyNotConfigured verifies that GetAPIKey returns an error
// when no API key is set in either the environment or config file.
func TestGetAPIKeyNotConfigured(t *testing.T) {
	setupTestDir(t)

	t.Setenv("MASSIVE_API_KEY", "")

	_, err := GetAPIKey()
	if err == nil {
		t.Error("expected error when no API key is configured, got nil")
	}
}

// TestSaveOverwritesExisting verifies that saving a config overwrites
// any previously saved configuration.
func TestSaveOverwritesExisting(t *testing.T) {
	setupTestDir(t)

	first := &Config{APIKey: "first-key", BaseURL: "https://api.massive.com"}
	if err := Save(first); err != nil {
		t.Fatalf("failed to save first config: %v", err)
	}

	second := &Config{APIKey: "second-key", BaseURL: "https://custom.example.com"}
	if err := Save(second); err != nil {
		t.Fatalf("failed to save second config: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loaded.APIKey != "second-key" {
		t.Errorf("expected second-key, got %s", loaded.APIKey)
	}

	if loaded.BaseURL != "https://custom.example.com" {
		t.Errorf("expected custom base URL, got %s", loaded.BaseURL)
	}
}
