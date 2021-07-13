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
		apiKey string
		apiURL string
	}{
		{
			apiKey: "",
			apiURL: "",
		},
		{
			apiKey: "test",
			apiURL: "",
		},
		{
			apiKey: "test",
			apiURL: "test",
		},
		{
			apiKey: "",
			apiURL: "test",
		},
	}
	for _, tt := range tests {
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
	}
}
