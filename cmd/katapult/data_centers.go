package main

import (
	"context"
	"errors"
	"fmt"

	_ "embed"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
	"github.com/spf13/cobra"
)

type dataCentersClient interface {
	List(ctx context.Context) ([]*core.DataCenter, *katapult.Response, error)
	Get(ctx context.Context, ref core.DataCenterRef) (*core.DataCenter, *katapult.Response, error)
}

//go:embed formatdata/dcs/list.txt
var dataCentersFormat string

//go:embed formatdata/dcs/get.txt
var getDataCenterFormat string

func listDataCentersCmd(client dataCentersClient) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List data centers",
		Long:    "List data centers.",
		RunE: outputWrapper(func(cmd *cobra.Command, args []string) (Output, error) {
			dcs, _, err := client.List(cmd.Context())
			if err != nil {
				return nil, err
			}

			return &genericOutput{
				item:                dcs,
				defaultTextTemplate: dataCentersFormat,
			}, nil
		}),
	}
}

func getDataCenterCmd(client dataCentersClient) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Args:  cobra.ExactArgs(1),
		Short: "Get details for a data center",
		Long:  "Get details for a data center.",
		RunE: outputWrapper(func(cmd *cobra.Command, args []string) (Output, error) {
			dc, _, err := client.Get(cmd.Context(), core.DataCenterRef{Permalink: args[0]})
			if err != nil {
				if errors.Is(err, katapult.ErrNotFound) {
					return nil, fmt.Errorf("unknown datacentre")
				}
				return nil, err
			}

			return &genericOutput{item: dc, defaultTextTemplate: getDataCenterFormat}, nil
		}),
	}
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
