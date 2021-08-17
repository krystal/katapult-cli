package main

import (
	"testing"

	"github.com/krystal/katapult-cli/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name string

		apiKey string
		apiURL string
		output string
	}{
		{
			name:   "empty values",
			apiKey: "",
			apiURL: "",
		},
		{
			name:   "only API URL blank",
			apiKey: "test",
			apiURL: "",
		},
		{
			name:   "both fields present",
			apiKey: "test",
			apiURL: "test",
		},
		{
			name:   "only API key blank",
			apiKey: "",
			apiURL: "test",
		},
		{
			name:   "test JSON output",
			apiKey: "test",
			apiURL: "test",
			output: "json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf, err := config.New()
			assert.NoError(t, err)
			conf.SetDefault("api_key", tt.apiKey)
			conf.SetDefault("api_url", tt.apiURL)
			cmd := configCommand(conf)
			outputFlag = tt.output
			assertCobraCommand(t, cmd, "", "")
			outputFlag = ""
		})
	}
}
