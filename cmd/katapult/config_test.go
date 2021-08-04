package main

import (
	"fmt"
	"testing"

	"github.com/krystal/katapult-cli/config"
	"github.com/stretchr/testify/assert"
)

const mockConfigFormat = `api_key: %s
api_url: %s
`

func TestConfig(t *testing.T) {
	tests := []struct {
		name string

		apiKey string
		apiURL string
		output string
		wants  string
	}{
		{
			name:   "empty values",
			apiKey: "",
			apiURL: "",
			wants:   fmt.Sprintf(mockConfigFormat, "", ""),
		},
		{
			name:   "only API URL blank",
			apiKey: "test",
			apiURL: "",
			wants:  fmt.Sprintf(mockConfigFormat, "test", ""),
		},
		{
			name:   "both fields present",
			apiKey: "test",
			apiURL: "test",
			wants:  fmt.Sprintf(mockConfigFormat, "test", "test"),
		},
		{
			name:   "only API key blank",
			apiKey: "",
			apiURL: "test",
			wants:  fmt.Sprintf(mockConfigFormat, "", "test"),
		},
		{
			name:   "test JSON output",
			apiKey: "test",
			apiURL: "test",
			output: "json",
			wants:  getTestData(t, "test_JSON_output.json"),
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("key: %s, url: %s", tt.apiKey, tt.apiURL), func(t *testing.T) {
			conf, err := config.New()
			assert.NoError(t, err)
			conf.SetDefault("api_key", tt.apiKey)
			conf.SetDefault("api_url", tt.apiURL)
			cmd := configCommand(conf)
			outputFlag = tt.output
			assertCobraCommand(t, cmd, "", tt.wants, "")
			outputFlag = ""
		})
	}
}
