package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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

			if strings.ToLower(cmd.Flag("output").Value.String()) == "json" {
				if nets == nil {
					nets = []*core.Network{}
				}
				if vnets == nil {
					vnets = []*core.VirtualNetwork{}
				}
				j, err := json.Marshal(map[string]interface{}{
					"networks":         nets,
					"virtual_networks": vnets,
				})
				if err != nil {
					return err
				}
				_, _ = cmd.OutOrStdout().Write(append(j, '\n'))
			} else {
				out := cmd.OutOrStdout()
				_, _ = fmt.Fprintln(out, "Networks:")
				for _, net := range nets {
					_, _ = fmt.Fprintf(out, " - %s [%s]\n", net.Name, net.ID)
				}

				if len(vnets) > 0 {
					_, _ = fmt.Fprintln(out, "Virtual Networks:")
					for _, net := range vnets {
						_, _ = fmt.Fprintf(out, " - %s [%s]\n", net.Name, net.ID)
					}
				}
			}

			return nil
		},
	}
	listFlags := list.PersistentFlags()
	listFlags.String("id", "", "The ID of the organization. Preferred over subdomain for lookups.")
	listFlags.String("subdomain", "", "The subdomain of the organization.")
	listFlags.StringP("output", "o", "text", "Defines the output type of the data centers. Can be text or json.")
	cmd.AddCommand(list)

	return cmd
}
