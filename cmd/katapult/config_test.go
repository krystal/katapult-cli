package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/krystal/katapult-cli/config"
)

const mockConfigFormat = `---
api_key: %s
api_url: %s

`

func TestConfig(t *testing.T) {
	conf, err := config.New()
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string

		apiKey string
		apiURL string
	}{
		{
			name:   "Check that the config command displays blank values properly",
			apiKey: "",
			apiURL: "",
		},
		{
			name:   "Check that the config command displays only API URL being blank properly",
			apiKey: "test",
			apiURL: "",
		},
		{
			name:   "Check that the config command handles both fields being present properly",
			apiKey: "test",
			apiURL: "test",
		},
		{
			name:   "Check that the config command displays only API key being blank properly",
			apiKey: "",
			apiURL: "test",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("key: %s, url: %s", tt.apiKey, tt.apiURL), func(t *testing.T) {
			conf.APIKey = tt.apiKey
			conf.APIURL = tt.apiURL
			conf.SetDefault("api_key", tt.apiKey)
			conf.SetDefault("api_url", tt.apiURL)
			cmd := configCommand(conf)
			stdout := &bytes.Buffer{}
			cmd.SetOut(stdout)
			if err := cmd.RunE(cmd, []string{}); err != nil {
				t.Fatal(err)
			}
			mockAPIKey := tt.apiKey
			if mockAPIKey == "" {
				mockAPIKey = "\"\""
			}
			mockAPIURL := tt.apiURL
			if mockAPIURL == "" {
				mockAPIURL = "\"\""
			}
			mockOutput := fmt.Sprintf(mockConfigFormat, mockAPIKey, mockAPIURL)
			if stdout.String() != mockOutput {
				t.Fatal("invalid result:\n\n", stdout.String())
			}
		})
	}
}
