package main

import (
	"fmt"
	"testing"

	"github.com/krystal/katapult-cli/config"
	"github.com/stretchr/testify/assert"
)

const mockConfigFormat = `---
api_key: %s
api_url: %s

`

func TestConfig(t *testing.T) {
	formatYaml := func(apiKey, apiURL string) string {
		expectedAPIKey := apiKey
		if expectedAPIKey == "" {
			expectedAPIKey = "\"\""
		}
		expectedAPIURL := apiURL
		if expectedAPIURL == "" {
			expectedAPIURL = "\"\""
		}
		return fmt.Sprintf(mockConfigFormat, expectedAPIKey, expectedAPIURL)
	}

	tests := []struct {
		name string

		apiKey string
		apiURL string
		args   []string
		wants  string
	}{
		{
			name:   "empty values",
			apiKey: "",
			apiURL: "",
			wants: formatYaml("", ""),
		},
		{
			name:   "only API URL blank",
			apiKey: "test",
			apiURL: "",
			wants: formatYaml("test", ""),
		},
		{
			name:   "both fields present",
			apiKey: "test",
			apiURL: "test",
			wants:  formatYaml("test", "test"),
		},
		{
			name:   "only API key blank",
			apiKey: "",
			apiURL: "test",
			wants:  formatYaml("", "test"),
		},
		{
			name:   "test JSON output",
			apiKey: "test",
			apiURL: "test",
			args:   []string{"-o", "json"},
			wants:  `{"api_key":"test","api_url":"test"}
`,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("key: %s, url: %s", tt.apiKey, tt.apiURL), func(t *testing.T) {
			conf, err := config.New()
			assert.NoError(t, err)
			conf.SetDefault("api_key", tt.apiKey)
			conf.SetDefault("api_url", tt.apiURL)
			cmd := configCommand(conf)
			if tt.args == nil {
				tt.args = []string{}
			}
			cmd.SetArgs(tt.args)
			assertCobraCommand(t, cmd, "", tt.wants, "")
		})
	}
}

func TestConfig_JSON(t *testing.T) {

}
