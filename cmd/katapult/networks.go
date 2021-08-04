package main

import (
	"context"
	"fmt"
	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
	"github.com/spf13/cobra"
)

type networksListClient interface {
	List(
		ctx context.Context,
		org core.OrganizationRef,
	) ([]*core.Network, []*core.VirtualNetwork, *katapult.Response, error)
}

func networksCmd(client networksListClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "networks",
		Aliases: []string{"net", "nets"},
		Short:   "Manage networks",
		Long:    "Get information about and manage networks.",
	}

	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Get list of networks available to a Organization",
		Long:    "Get list of networks available to a Organization.",
		RunE: renderOption(func(cmd *cobra.Command, _ []string) (Output, error) {
			id := cmd.Flag("id").Value.String()
			ref := core.OrganizationRef{ID: id}
			if id == "" {
				subdomain := cmd.Flag("subdomain").Value.String()
				if subdomain == "" {
					return nil, fmt.Errorf("both ID and subdomain are unset")
				}
				ref = core.OrganizationRef{SubDomain: subdomain}
			}

			nets, vnets, _, err := client.List(cmd.Context(), ref)
			if err != nil {
				return nil, err
			}

			return genericOutput{
				item: map[string]interface{}{
					"nets":  nets,
					"vnets": vnets,
				},
				tpl:  "",
			}, nil
		}),
	}
	listFlags := list.PersistentFlags()
	listFlags.String("id", "", "The ID of the organization. Preferred over subdomain for lookups.")
	listFlags.String("subdomain", "", "The subdomain of the organization.")
	cmd.AddCommand(list)

	return cmd
}
