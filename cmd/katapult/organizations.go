package main

import (
	"fmt"

	"github.com/krystal/go-katapult/core"

	"github.com/spf13/cobra"
)

func organizationsCmd(client core.RequestMaker) *cobra.Command {
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
			orgs, _, err := core.NewOrganizationsClient(client).List(cmd.Context())
			if err != nil {
				return err
			}

			stdout := cmd.OutOrStdout()
			for _, org := range orgs {
				fmt.Fprintf(stdout, " - %s (%s) [%s]\n", org.Name, org.SubDomain, org.ID)
			}

			return nil
		},
	}
	cmd.AddCommand(list)

	return cmd
}
