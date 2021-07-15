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
			var apiUrl string
			if c.(*katapult.Client).BaseURL != nil {
				apiUrl = c.(*katapult.Client).BaseURL.String()
			}
			expectedApiUrl := tt.apiURL
			if expectedApiUrl == "" {
				expectedApiUrl = "https://api.katapult.io"
			}
			assert.Equal(t, expectedApiUrl, apiUrl)
		})
	}
}
