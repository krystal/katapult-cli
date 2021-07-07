package main

import (
	"fmt"

	"github.com/krystal/go-katapult/core"

	"github.com/krystal/go-katapult"

	"github.com/spf13/cobra"
)

func listDataCentersCmd(client *katapult.Client) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List data centers",
		Long:    "List data centers.",
		RunE: func(cmd *cobra.Command, args []string) error {
			dcs, _, err := core.NewDataCentersClient(client).List(cmd.Context())
			if err != nil {
				return err
			}

			for _, dc := range dcs {
				_, _ = fmt.Fprintf(
					cmd.OutOrStdout(),
					" - %s (%s) [%s] / %s\n",
					dc.Name, dc.Permalink, dc.ID, dc.Country.Name,
				)
			}

			return nil
		},
	}
}

func getDataCenterCmd(client *katapult.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Args:  cobra.ExactArgs(1),
		Short: "Get details for a data center",
		Long:  "Get details for a data center.",
		RunE: func(cmd *cobra.Command, args []string) error {
			dcClient := core.NewDataCentersClient(client)
			dc, resp, err := dcClient.Get(cmd.Context(), core.DataCenterRef{Permalink: args[0]})
			if err != nil {
				if resp != nil && resp.StatusCode == 404 {
					return fmt.Errorf("unknown datacentre")
				}
				return err
			}

			_, _ = fmt.Fprintf(
				cmd.OutOrStdout(),
				"%s (%s) [%s] / %s\n",
				dc.Name, dc.Permalink, dc.ID, dc.Country.Name,
			)

			return nil
		},
	}
}

func dataCentersCmd(client *katapult.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dc",
		Aliases: []string{"dcs", "data-centers", "data_centers"},
		Short:   "Get information about data centers",
		Long:    "Get information about data centers.",
	}

	cmd.AddCommand(listDataCentersCmd(client))
	cmd.AddCommand(getDataCenterCmd(client))

	return cmd
}
