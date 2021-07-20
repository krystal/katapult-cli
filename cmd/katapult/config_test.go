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
	tests := []struct {
		name string

		apiKey string
		apiURL string
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
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("key: %s, url: %s", tt.apiKey, tt.apiURL), func(t *testing.T) {
			conf, err := config.New()
			if err != nil {
				t.Fatal(err)
			}
			conf.SetDefault("api_key", tt.apiKey)
			conf.SetDefault("api_url", tt.apiURL)
			cmd := configCommand(conf)
			out := &bytes.Buffer{}
			cmd.SetOut(out)
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
			if out.String() != mockOutput {
				t.Fatal("invalid result:\n\n", out.String())
			}
		})
	}
}
