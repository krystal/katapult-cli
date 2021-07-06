package main

import (
	"bytes"
	"fmt"
	"github.com/krystal/katapult-cli/config"
	"testing"
)

const mockConfigFormat = `---
api_key: %s
api_url: %s

`

var mockConfigSuite = []struct{
	apiKey string
	apiUrl string
}{
	{
		apiKey: "",
		apiUrl: "",
	},
	{
		apiKey: "test",
		apiUrl: "",
	},
	{
		apiKey: "test",
		apiUrl: "test",
	},
	{
		apiKey: "",
		apiUrl: "test",
	},
}

func TestConfig(t *testing.T) {
	conf, err := config.New()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range mockConfigSuite {
		conf.APIKey = v.apiKey
		conf.APIURL = v.apiUrl
		conf.SetDefault("api_key", v.apiKey)
		conf.SetDefault("api_url", v.apiUrl)
		cmd := configCommand(conf)
		stdout := &bytes.Buffer{}
		cmd.SetOut(stdout)
		if err := cmd.RunE(cmd, []string{}); err != nil {
			t.Fatal(err)
		}
		mockApiKey := v.apiKey
		if mockApiKey == "" {
			mockApiKey = "\"\""
		}
		mockApiUrl := v.apiUrl
		if mockApiUrl == "" {
			mockApiUrl = "\"\""
		}
		mockOutput := fmt.Sprintf(mockConfigFormat, mockApiKey, mockApiUrl)
		if stdout.String() != mockOutput {
			t.Fatal("invalid result:\n\n", stdout.String())
		}
	}
}
