package main

import (
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/krystal/katapult-cli/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newClient(t *testing.T) {
	tests := []struct {
		name string

		apiKey  string
		apiURL  string
		wantErr string
	}{
		{
			name:   "empty API URL",
			apiKey: "test",
		},
		{
			name:    "invalid URL",
			wantErr: "invalid API URL: @@:",
			apiURL:  "@@:",
		},
		{
			name:   "both values",
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

			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)

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
