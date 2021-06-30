package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dataCentersCmd)
	dataCentersCmd.AddCommand(dataCentersListCmd)
	dataCentersCmd.AddCommand(dataCentersGetCmd)
}

var (
	dataCentersCmd = &cobra.Command{
		Use:     "dc",
		Aliases: []string{"dcs", "data-centers", "data_centers"},
		Short:   "Get information about data centers",
		Long:    "Get information about data centers.",
	}

	dataCentersListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List data centers",
		Long:    "List data centers.",
		Run: func(cmd *cobra.Command, args []string) {
			dcs, _, err := cl.DataCenters.List(ctx)
			if err != nil {
				er(err)
			}

			for _, dc := range dcs {
				fmt.Printf(
					" - %s (%s) [%s] / %s\n",
					dc.Name, dc.Permalink, dc.ID, dc.Country.Name,
				)
			}
		},
	}

	dataCentersGetCmd = &cobra.Command{
		Use:   "get",
		Args:  cobra.ExactArgs(1),
		Short: "Get details for a data center",
		Long:  "Get details for a data center.",
		Run: func(cmd *cobra.Command, args []string) {
			dc, resp, err := cl.DataCenters.Get(ctx, args[0])
			if err != nil {
				if resp.StatusCode != 404 {
					er(err)
				} else {
					dc, _, err = cl.DataCenters.GetByPermalink(ctx, args[0])
					if err != nil && resp.StatusCode != 404 {
						er(err)
					}
				}
			}

			fmt.Printf(
				"%s (%s) [%s] / %s\n",
				dc.Name, dc.Permalink, dc.ID, dc.Country.Name,
			)
		},
	}
)
