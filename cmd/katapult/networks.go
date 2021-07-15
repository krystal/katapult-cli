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
		RunE: func(cmd *cobra.Command, _ []string) error {
			id := cmd.Flag("id").Value.String()
			ref := core.OrganizationRef{ID: id}
			if id == "" {
				subdomain := cmd.Flag("subdomain").Value.String()
				if subdomain == "" {
					return fmt.Errorf("both ID and subdomain are unset")
				}
				ref = core.OrganizationRef{SubDomain: subdomain}
			}

			nets, vnets, _, err := client.List(cmd.Context(), ref)
			if err != nil {
				return err
			}

			stdout := cmd.OutOrStdout()
			_, _ = fmt.Fprintln(stdout, "Networks:")
			for _, net := range nets {
				_, _ = fmt.Fprintf(stdout, " - %s [%s]\n", net.Name, net.ID)
			}

			if len(vnets) > 0 {
				_, _ = fmt.Fprintln(stdout, "Virtual Networks:")
				for _, net := range vnets {
					_, _ = fmt.Fprintf(stdout, " - %s [%s]\n", net.Name, net.ID)
				}
			}

			return nil
		},
	}
	listFlags := list.PersistentFlags()
	listFlags.String("id", "", "The ID of the organization. Preferred over subdomain for lookups.")
	listFlags.String("subdomain", "", "The subdomain of the organization.")
	cmd.AddCommand(list)

	return cmd
}
