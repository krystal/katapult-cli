package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(networksCmd)
	networksCmd.AddCommand(networksListCmd)
}

var (
	networksCmd = &cobra.Command{
		Use:     "networks",
		Aliases: []string{"net", "nets"},
		Short:   "Manage networks",
		Long:    "Get information about and manage networks.",
	}
	networksListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(1),
		Short:   "Get list of networks available to a Organization",
		Long:    "Get list of networks available to a Organization.",
		Run: func(cmd *cobra.Command, args []string) {
			nets, vnets, _, err := cl.Networks.List(ctx, args[0])
			if err != nil {
				er(err)
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
		},
	}
)
