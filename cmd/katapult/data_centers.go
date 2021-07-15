package main

import (
	"context"
	"fmt"
	"github.com/krystal/go-katapult"

	"github.com/krystal/go-katapult/core"

	"github.com/spf13/cobra"
)

type dataCenterListClient interface {
	List(ctx context.Context) ([]*core.DataCenter, *katapult.Response, error)
}

func listDataCentersCmd(client dataCenterListClient) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List data centers",
		Long:    "List data centers.",
		RunE: func(cmd *cobra.Command, args []string) error {
			dcs, _, err := client.List(cmd.Context())
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

type dataCenterGetClient interface {
	Get(ctx context.Context, ref core.DataCenterRef) (*core.DataCenter, *katapult.Response, error)
}

func getDataCenterCmd(client dataCenterGetClient) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Args:  cobra.ExactArgs(1),
		Short: "Get details for a data center",
		Long:  "Get details for a data center.",
		RunE: func(cmd *cobra.Command, args []string) error {
			dc, resp, err := client.Get(cmd.Context(), core.DataCenterRef{Permalink: args[0]})
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

func dataCentersCmd(client core.RequestMaker) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dc",
		Aliases: []string{"dcs", "data-centers", "data_centers"},
		Short:   "Get information about data centers",
		Long:    "Get information about data centers.",
	}

	cmd.AddCommand(listDataCentersCmd(core.NewDataCentersClient(client)))
	cmd.AddCommand(getDataCenterCmd(core.NewDataCentersClient(client)))

	return cmd
}
