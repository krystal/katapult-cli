package main

import (
	"testing"

	"github.com/krystal/go-katapult"

	"github.com/krystal/katapult-cli/config"
)

func TestNewClient_APIKey(t *testing.T) {
	c, err := newClient(&config.Config{
		APIKey: "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	if c.(*katapult.Client).APIKey != "test" {
		t.Fatal("API key is unset")
	}
}

func TestNewClient_BaseURL(t *testing.T) {
	// Test valid base URL.
	c, err := newClient(&config.Config{
		APIKey: "test",
		APIURL: "https://example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if c.(*katapult.Client).APIKey != "test" {
		t.Fatal("API key is set to", c.(*katapult.Client).APIKey)
	}
	s := c.(*katapult.Client).BaseURL.String()
	if s != "https://example.com" {
		t.Fatal("invalid base URL:", s)
	}

	// Test invalid base URL.
	_, err = newClient(&config.Config{
		APIKey: "test",
		APIURL: "this is a test",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}
