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

var mockConfigSuite = []struct {
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

func TestConfig(t *testing.T) {
	conf, err := config.New()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range mockConfigSuite {
		conf.APIKey = v.apiKey
		conf.APIURL = v.apiURL
		conf.SetDefault("api_key", v.apiKey)
		conf.SetDefault("api_url", v.apiURL)
		cmd := configCommand(conf)
		stdout := &bytes.Buffer{}
		cmd.SetOut(stdout)
		if err := cmd.RunE(cmd, []string{}); err != nil {
			t.Fatal(err)
		}
		mockAPIKey := v.apiKey
		if mockAPIKey == "" {
			mockAPIKey = "\"\""
		}
		mockAPIURL := v.apiURL
		if mockAPIURL == "" {
			mockAPIURL = "\"\""
		}
		mockOutput := fmt.Sprintf(mockConfigFormat, mockAPIKey, mockAPIURL)
		if stdout.String() != mockOutput {
			t.Fatal("invalid result:\n\n", stdout.String())
		}
	}
}
