package client

import (
	"context"
	"log"
	"net/url"

	"github.com/krystal/go-katapult"
	"github.com/krystal/katapult-cli/pkg/config"
	"golang.org/x/oauth2"
)

func New(ctx context.Context, conf *config.Config) *katapult.Client {
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
