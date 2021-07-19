package main

import (
	"context"
	"fmt"
	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
	"github.com/spf13/cobra"
)

type organisationsClient interface {
	List(
		ctx context.Context,
	) ([]*core.Organization, *katapult.Response, error)
}

func organizationsCmd(client organisationsClient) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			orgs, _, err := client.List(cmd.Context())
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			for _, org := range orgs {
				fmt.Fprintf(out, " - %s (%s) [%s]\n", org.Name, org.SubDomain, org.ID)
			}

			return nil
		},
	}
	cmd.AddCommand(list)

	return cmd
}
