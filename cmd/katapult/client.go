package main

import (
	"fmt"
	"net/url"

	"github.com/krystal/go-katapult"
	"github.com/krystal/katapult-cli/config"
)

// Create a new Katapult client.
func newClient(conf *config.Config) (*katapult.Client, error) {
	a := []katapult.Opt{katapult.WithAPIKey(conf.APIKey)}
	if conf.APIURL != "" {
		apiURL, err := url.Parse(conf.APIURL)
		if err != nil {
			return nil, fmt.Errorf("Invalid API URL: %s\n", conf.APIURL)
		}
		a = append(a, katapult.WithBaseURL(apiURL))
	}
	c, err := katapult.New(a...)
	if err != nil {
		return nil, err
	}
	return c, nil
}
