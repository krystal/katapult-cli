package main

import (
	"fmt"
	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
	"github.com/krystal/katapult-cli/config"
	"net/url"
)

// Create a new Katapult client.
func newClient(conf *config.Config) (core.RequestMaker, error) {
	a := []katapult.Opt{katapult.WithAPIKey(conf.APIKey)}
	if conf.APIURL != "" {
		apiURL, err := url.Parse(conf.APIURL)
		if err != nil {
			return nil, fmt.Errorf("invalid API URL: %s", conf.APIURL)
		}
		a = append(a, katapult.WithBaseURL(apiURL))
	}
	c, err := katapult.New(a...)
	if err != nil {
		return nil, err
	}
	return c, nil
}
