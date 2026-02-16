//
// Date: 2026-02-15
// Copyright (c) 2026. All rights reserved.
//

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestNewClient verifies that NewClient creates a client with the
// correct default base URL and the provided API key.
func TestNewClient(t *testing.T) {
	client := NewClient("test-key")

	if client.baseURL != defaultBaseURL {
		t.Errorf("expected base URL %s, got %s", defaultBaseURL, client.baseURL)
	}

	if client.apiKey != "test-key" {
		t.Errorf("expected API key test-key, got %s", client.apiKey)
	}

	if client.httpClient == nil {
		t.Error("expected httpClient to be initialized")
	}
}

// TestSetBaseURL verifies that SetBaseURL correctly overrides the
// client's base URL for pointing at mock servers.
func TestSetBaseURL(t *testing.T) {
	client := NewClient("test-key")
	client.SetBaseURL("http://localhost:9999")

	if client.baseURL != "http://localhost:9999" {
		t.Errorf("expected http://localhost:9999, got %s", client.baseURL)
	}
}

// TestGetAddsAPIKey verifies that the client appends the apiKey query
// parameter to every outgoing request.
func TestGetAddsAPIKey(t *testing.T) {
	var receivedKey string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedKey = r.URL.Query().Get("apiKey")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"OK"}`))
	}))
	defer server.Close()

	client := NewClient("my-secret-key")
	client.SetBaseURL(server.URL)

	var result map[string]interface{}
	err := client.get("/test", nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedKey != "my-secret-key" {
		t.Errorf("expected apiKey=my-secret-key, got %s", receivedKey)
	}
}

// TestGetAddsQueryParams verifies that additional query parameters are
// correctly appended to the request URL alongside the API key.
func TestGetAddsQueryParams(t *testing.T) {
	var receivedParams map[string]string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedParams = map[string]string{
			"apiKey": r.URL.Query().Get("apiKey"),
			"search": r.URL.Query().Get("search"),
			"limit":  r.URL.Query().Get("limit"),
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"OK"}`))
	}))
	defer server.Close()

	client := NewClient("key")
	client.SetBaseURL(server.URL)

	params := map[string]string{
		"search": "Apple",
		"limit":  "10",
	}

	var result map[string]interface{}
	err := client.get("/test", params, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedParams["search"] != "Apple" {
		t.Errorf("expected search=Apple, got %s", receivedParams["search"])
	}

	if receivedParams["limit"] != "10" {
		t.Errorf("expected limit=10, got %s", receivedParams["limit"])
	}
}

// TestGetSkipsEmptyParams verifies that empty string parameters are not
// included in the request URL.
func TestGetSkipsEmptyParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("empty") != "" {
			t.Error("empty param should not be sent")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"OK"}`))
	}))
	defer server.Close()

	client := NewClient("key")
	client.SetBaseURL(server.URL)

	params := map[string]string{
		"empty":  "",
		"filled": "value",
	}

	var result map[string]interface{}
	client.get("/test", params, &result)
}

// TestGetHandlesNon200Status verifies that the client returns an error
// containing the status code and response body for non-200 responses.
func TestGetHandlesNon200Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"NOT_FOUND","message":"Data not found."}`))
	}))
	defer server.Close()

	client := NewClient("key")
	client.SetBaseURL(server.URL)

	var result map[string]interface{}
	err := client.get("/test", nil, &result)
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}

	expected := "API error (status 404)"
	if len(err.Error()) < len(expected) || err.Error()[:len(expected)] != expected {
		t.Errorf("expected error to start with %q, got %q", expected, err.Error())
	}
}

// TestGetHandles500Error verifies that server errors are properly
// reported with the status code and body.
func TestGetHandles500Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
	}))
	defer server.Close()

	client := NewClient("key")
	client.SetBaseURL(server.URL)

	var result map[string]interface{}
	err := client.get("/test", nil, &result)
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}

	expected := "API error (status 500)"
	if len(err.Error()) < len(expected) || err.Error()[:len(expected)] != expected {
		t.Errorf("expected error to start with %q, got %q", expected, err.Error())
	}
}

// TestGetHandlesInvalidJSON verifies that the client returns an error
// when the response body contains invalid JSON.
func TestGetHandlesInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`not valid json`))
	}))
	defer server.Close()

	client := NewClient("key")
	client.SetBaseURL(server.URL)

	var result map[string]interface{}
	err := client.get("/test", nil, &result)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

// TestGetHandlesConnectionError verifies that the client returns an
// error when it cannot connect to the server.
func TestGetHandlesConnectionError(t *testing.T) {
	client := NewClient("key")
	client.SetBaseURL("http://localhost:1")

	var result map[string]interface{}
	err := client.get("/test", nil, &result)
	if err == nil {
		t.Fatal("expected connection error, got nil")
	}
}

// TestGetSendsCorrectPath verifies that the request path is correctly
// constructed from the base URL and the provided path.
func TestGetSendsCorrectPath(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"OK"}`))
	}))
	defer server.Close()

	client := NewClient("key")
	client.SetBaseURL(server.URL)

	var result map[string]interface{}
	client.get("/v1/open-close/AAPL/2025-01-06", nil, &result)

	if receivedPath != "/v1/open-close/AAPL/2025-01-06" {
		t.Errorf("expected path /v1/open-close/AAPL/2025-01-06, got %s", receivedPath)
	}
}
