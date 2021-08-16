package main

import (
	"context"

	_ "embed"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
	"github.com/spf13/cobra"
)

type organizationsClient interface {
	List(
		ctx context.Context,
	) ([]*core.Organization, *katapult.Response, error)
}

//go:embed formatdata/orgs/list.txt
var organizationsListFormat string

func organizationsCmd(client organizationsClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "org",
		Aliases: []string{"orgs", "organization", "organizations"},
		Short:   "Manage organizations",
		Long:    "Get information about and manage organizations.",
	}

	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Get list of organizations",
		Long:    "Get list of organizations.",
		RunE: outputWrapper(func(cmd *cobra.Command, args []string) (Output, error) {
			orgs, _, err := client.List(cmd.Context())
			if err != nil {
				return nil, err
			}

			return &genericOutput{
				item: orgs,
				tpl:  organizationsListFormat,
			}, nil
		}),
	}
	cmd.AddCommand(list)

	return cmd
}
