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

type virtualMachineDiskTemplatesClient interface {
	List(
		ctx context.Context,
		org core.OrganizationRef,
		opts *core.DiskTemplateListOptions,
	) ([]*core.DiskTemplate, *katapult.Response, error)
}

func listAllDiskTemplates(ctx context.Context, org core.OrganizationRef, diskTemplatesClient virtualMachineDiskTemplatesClient) ([]*core.DiskTemplate, error) {
	totalPages := 1
	allImages := make([]*core.DiskTemplate, 0)
	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		images, resp, err := diskTemplatesClient.List(ctx, org, &core.DiskTemplateListOptions{Page: pageNum, IncludeUniversal: true})
		if err != nil {
			return nil, err
		}
		if resp.Pagination != nil {
			totalPages = resp.Pagination.TotalPages
		}
		allImages = append(allImages, images...)
	}
	return allImages, nil
}

func getStringIndex(needle string, haystack []string) int {
	for i, v := range haystack {
		if v == needle {
			return i
		}
	}
	return -1
}

// TODO: Move this!
func createCmd(
	orgsClient organisationsClient, dcsClient dataCentersClient,
	vmPackagesClient virtualMachinePackagesClient,
	diskTemplatesClient virtualMachineDiskTemplatesClient,
) *cobra.Command {
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
			index := getStringIndex(orgStr, orgStrs)
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
			index = getStringIndex(dcStr, dcStrs)
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
			index = getStringIndex(packageStr, packageStrs)
			package_ := packages[index]

			// Ask about the distribution.
			distributions, err := listAllDiskTemplates(cmd.Context(), core.OrganizationRef{ID: org.ID}, diskTemplatesClient)
			if err != nil {
				return err
			}
			distributionStrs := make([]string, len(distributions))
			for i, distribution := range distributions {
				distributionStrs[i] = fmt.Sprintf("%s [%s]", distribution.Name, distribution.ID)
			}
			distributionStr := console.FuzzySelector("Which distribution would you like?", distributionStrs, cmd.InOrStdin())
			index = getStringIndex(distributionStr, distributionStrs)
			distribution := distributions[index]

			// TODO
			fmt.Println(org, dc, package_, distribution)

			// Return no errors.
			return nil
		},
	}

	return cmd
}
