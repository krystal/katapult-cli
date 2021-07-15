package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/krystal/go-katapult"

	"github.com/krystal/katapult-cli/config"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name string

		apiKey string
		apiURL string
	}{
		{
			name:   "Handle blank API URL",
			apiKey: "test",
		},
		{
			name:   "Handle both values",
			apiKey: "test",
			apiURL: "https://example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := newClient(&config.Config{
				APIKey: tt.apiKey,
				APIURL: tt.apiURL,
			})
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.apiKey, c.(*katapult.Client).APIKey)
			var apiURL string
			if c.(*katapult.Client).BaseURL != nil {
				apiURL = c.(*katapult.Client).BaseURL.String()
			}
			expectedAPIURL := tt.apiURL
			if expectedAPIURL == "" {
				expectedAPIURL = "https://api.katapult.io"
			}
			assert.Equal(t, expectedAPIURL, apiURL)
		})
	}
}
