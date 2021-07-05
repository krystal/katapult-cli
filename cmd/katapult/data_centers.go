package main

import (
	"fmt"
	"github.com/krystal/go-katapult"

	"github.com/spf13/cobra"
)

func dataCentersCmd(client *katapult.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dc",
		Aliases: []string{"dcs", "data-centers", "data_centers"},
		Short:   "Get information about data centers",
		Long:    "Get information about data centers.",
	}

	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List data centers",
		Long:    "List data centers.",
		RunE: func(cmd *cobra.Command, args []string) error {
			dcs, _, err := client.DataCenters.List(cmd.Context())
			if err != nil {
				return err
			}

			for _, dc := range dcs {
				fmt.Printf(
					" - %s (%s) [%s] / %s\n",
					dc.Name, dc.Permalink, dc.ID, dc.Country.Name,
				)
			}

			return nil
		},
	}
	cmd.AddCommand(listCmd)

	getCmd := &cobra.Command{
		Use:   "get",
		Args:  cobra.ExactArgs(1),
		Short: "Get details for a data center",
		Long:  "Get details for a data center.",
		RunE: func(cmd *cobra.Command, args []string) error {
			dc, resp, err := client.DataCenters.Get(cmd.Context(), args[0])
			if err != nil {
				if resp.StatusCode != 404 {
					return err
				} else {
					dc, _, err = client.DataCenters.GetByPermalink(cmd.Context(), args[0])
					if err != nil && resp.StatusCode != 404 {
						return err
					}
				}
			}

			fmt.Printf(
				"%s (%s) [%s] / %s\n",
				dc.Name, dc.Permalink, dc.ID, dc.Country.Name,
			)

			return nil
		},
	}
	cmd.AddCommand(getCmd)

	return cmd
}
