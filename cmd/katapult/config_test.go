package main

import (
	"testing"

	"github.com/krystal/katapult-cli/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name string

		apiToken string
		apiURL   string
		output   string
	}{
		{
			name:     "empty values",
			apiToken: "",
			apiURL:   "",
		},
		{
			name:     "only API URL blank",
			apiToken: "test",
			apiURL:   "",
		},
		{
			name:     "both fields present",
			apiToken: "test",
			apiURL:   "test",
		},
		{
			name:     "only API key blank",
			apiToken: "",
			apiURL:   "test",
		},
		{
			name:     "json output",
			apiToken: "testKey",
			apiURL:   "testURL",
			output:   "json",
		},
		{
			name:     "yaml output",
			apiToken: "testKey",
			apiURL:   "testURL",
			output:   "yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf, err := config.New()
			assert.NoError(t, err)
			conf.SetDefault("api_token", tt.apiToken)
			conf.SetDefault("api_url", tt.apiURL)
			cmd := configCommand(conf)
			outputFlag = tt.output
			assertCobraCommand(t, cmd, "", "")
			outputFlag = ""
		})
	}
}
