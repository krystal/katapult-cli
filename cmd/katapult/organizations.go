package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(organizationsCmd)
	organizationsCmd.AddCommand(organizationsListCmd)
}

var (
	organizationsCmd = &cobra.Command{
		Use:     "org",
		Aliases: []string{"orgs", "organization", "organizations"},
		Short:   "Manage organizations",
		Long:    "Get information about and manage organizations.",
	}
	organizationsListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Get list of organizations",
		Long:    "Get list of organizations.",
		Run: func(cmd *cobra.Command, args []string) {
			orgs, _, err := cl.Organizations.List(ctx)
			if err != nil {
				er(err)
			}

			for _, org := range orgs {
				fmt.Printf(" - %s (%s) [%s]\n", org.Name, org.SubDomain, org.ID)
			}
		},
	}
)
