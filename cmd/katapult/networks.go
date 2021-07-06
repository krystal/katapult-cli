package main

import (
	"fmt"
	"github.com/krystal/go-katapult/core"

	"github.com/krystal/go-katapult"

	"github.com/spf13/cobra"
)

func networksCmd(client *katapult.Client) *cobra.Command {
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
					fmt.Println("Both ID and subdomain are unset.")
					return nil
				}
				ref = core.OrganizationRef{SubDomain: subdomain}
			}

			nets, vnets, _, err := core.NewNetworksClient(client).List(cmd.Context(), ref)
			if err != nil {
				return err
			}

			fmt.Println("Networks:")
			for _, net := range nets {
				fmt.Printf(" - %s [%s]\n", net.Name, net.ID)
			}

			if len(vnets) > 0 {
				fmt.Println("Virtual Networks:")
				for _, net := range vnets {
					fmt.Printf(" - %s [%s]\n", net.Name, net.ID)
				}
			}

			return nil
		},
	}
	listFlags := list.PersistentFlags()
	listFlags.String("id", "", "The ID of the organisation. Preferred over subdomain for lookups.")
	listFlags.String("subdomain", "", "The subdomain of the organisation.")
	cmd.AddCommand(list)

	return cmd
}
