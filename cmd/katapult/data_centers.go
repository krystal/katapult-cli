package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
	"github.com/spf13/cobra"
)

type dataCentersClient interface {
	List(ctx context.Context) ([]*core.DataCenter, *katapult.Response, error)
	Get(ctx context.Context, ref core.DataCenterRef) (*core.DataCenter, *katapult.Response, error)
}

func listDataCentersCmd(client dataCentersClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List data centers",
		Long:    "List data centers.",
		RunE: func(cmd *cobra.Command, args []string) error {
			dcs, _, err := client.List(cmd.Context())
			if err != nil {
				return err
			}

			if strings.ToLower(cmd.Flag("output").Value.String()) == "json" {
				j, err := json.Marshal(dcs)
				if err != nil {
					return err
				}
				_, _ = cmd.OutOrStdout().Write(append(j, '\n'))
			} else {
				for _, dc := range dcs {
					_, _ = fmt.Fprintf(
						cmd.OutOrStdout(),
						" - %s (%s) [%s] / %s\n",
						dc.Name, dc.Permalink, dc.ID, dc.Country.Name,
					)
				}
			}

			return nil
		},
	}

	flags := cmd.PersistentFlags()
	flags.StringP("output", "o", "text", "Defines the output type of the data centers. Can be text or json.")

	return cmd
}

func getDataCenterCmd(client dataCentersClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Args:  cobra.ExactArgs(1),
		Short: "Get details for a data center",
		Long:  "Get details for a data center.",
		RunE: func(cmd *cobra.Command, args []string) error {
			dc, _, err := client.Get(cmd.Context(), core.DataCenterRef{Permalink: args[0]})
			if err != nil {
				if errors.Is(err, katapult.ErrNotFound) {
					return fmt.Errorf("unknown datacentre")
				}
				return err
			}

			if strings.ToLower(cmd.Flag("output").Value.String()) == "json" {
				j, err := json.Marshal(dc)
				if err != nil {
					return err
				}
				_, _ = cmd.OutOrStdout().Write(append(j, '\n'))
			} else {
				_, _ = fmt.Fprintf(
					cmd.OutOrStdout(),
					"%s (%s) [%s] / %s\n",
					dc.Name, dc.Permalink, dc.ID, dc.Country.Name,
				)
			}

			return nil
		},
	}

	flags := cmd.PersistentFlags()
	flags.StringP("output", "o", "text", "Defines the output type of the data centers. Can be text or json.")

	return cmd
}

func dataCentersCmd(client dataCentersClient) *cobra.Command {
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
