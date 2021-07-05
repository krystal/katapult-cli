package main

import (
	"context"
	"log"
	"net/url"

	"github.com/krystal/go-katapult"
	"github.com/krystal/katapult-cli/config"
	"golang.org/x/oauth2"
)

// Create a new Katapult client.
func newClient(ctx context.Context, conf *config.Config) *katapult.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: conf.APIKey},
	)
	tc := oauth2.NewClient(ctx, ts)

	c := katapult.NewClient(tc)

	if conf.APIURL != "" {
		apiURL, err := url.Parse(conf.APIURL)
		if err != nil {
			log.Fatalf("Invalid API URL: %s\n", conf.APIURL)
		}
		c.SetBaseURL(apiURL)
	}

	return c
}
