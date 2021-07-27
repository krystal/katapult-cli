package main

import (
	"context"
	"fmt"
	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
	"github.com/krystal/katapult-cli/cmd/katapult/console"
	"github.com/spf13/cobra"
)

type virtualMachinePackagesClient interface {
	List(
		ctx context.Context,
		opts *core.ListOptions,
	) ([]*core.VirtualMachinePackage, *katapult.Response, error)
}

func listAllVMPackages(ctx context.Context, vmPackagesClient virtualMachinePackagesClient) ([]*core.VirtualMachinePackage, error) {
	totalPages := 1
	allPackages := make([]*core.VirtualMachinePackage, 0)
	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		packages, resp, err := vmPackagesClient.List(ctx, &core.ListOptions{Page: pageNum})
		if err != nil {
			return nil, err
		}
		if resp.Pagination != nil {
			totalPages = resp.Pagination.TotalPages
		}
		allPackages = append(allPackages, packages...)
	}
	return allPackages, nil
}

// TODO: Move this!
func createCmd(orgsClient organisationsClient, dcsClient dataCentersClient, vmPackagesClient virtualMachinePackagesClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "!! TEMPORARY UNTIL VM COMMANDS EXIST !!",
		Long:    "!! TEMPORARY UNTIL VM COMMANDS EXIST !!",
		RunE:    func(cmd *cobra.Command, args []string) error {
			// TODO: Accept argument from env var.

			// List the organisations.
			orgs, _, err := orgsClient.List(cmd.Context())
			if err != nil {
				return err
			}

			// Create a fuzzy searcher for organisations.
			orgStrs := make([]string, len(orgs))
			for i, org := range orgs {
				orgStrs[i] = fmt.Sprintf("%s (%s) [%s]", org.Name, org.SubDomain, org.ID)
			}
			orgStr := console.FuzzySelector("Which organisation would you like to deploy the VM in?", orgStrs, cmd.InOrStdin())
			index := 0
			for i, v := range orgStrs {
				if v == orgStr {
					index = i
					break
				}
			}
			org := orgs[index]

			// List the datacenters.
			dcs, _, err := dcsClient.List(cmd.Context())
			if err != nil {
				return err
			}

			// Create a fuzzy searcher for data centres.
			dcStrs := make([]string, len(dcs))
			for i, dc := range dcs {
				dcStrs[i] = fmt.Sprintf("%s (%s) [%s] / %s", dc.Name, dc.Permalink, dc.ID, dc.Country.Name)
			}
			dcStr := console.FuzzySelector("Which DC would you like to deploy the VM in?", dcStrs, cmd.InOrStdin())
			for i, v := range dcStrs {
				if v == dcStr {
					index = i
					break
				}
			}
			dc := dcs[index]

			// List the packages.
			packages, err := listAllVMPackages(cmd.Context(), vmPackagesClient)
			if err != nil {
				return err
			}
			packageStrs := make([]string, len(packages))
			for i, package_ := range packages {
				packageStrs[i] = fmt.Sprintf(
					"%s (%d cores, %d GB RAM) [%s]", package_.Name, package_.CPUCores,
					package_.MemoryInGB, package_.ID)
			}
			packageStr := console.FuzzySelector("Which VM package would you like?", packageStrs, cmd.InOrStdin())
			for i, v := range packageStrs {
				if v == packageStr {
					index = i
					break
				}
			}
			package_ := packages[index]

			// TODO
			fmt.Println(org, dc, packages, package_)

			// Return no errors.
			return nil
		},
	}

	return cmd
}
