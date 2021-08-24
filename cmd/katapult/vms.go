package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/buger/goterm"
	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/buildspec"
	"github.com/krystal/go-katapult/core"
	"github.com/krystal/katapult-cli/cmd/katapult/console"
	"github.com/spf13/cobra"
)

type virtualMachinesClient interface {
	List(
		ctx context.Context,
		org core.OrganizationRef,
		opts *core.ListOptions,
	) ([]*core.VirtualMachine, *katapult.Response, error)

	Delete(
		ctx context.Context,
		ref core.VirtualMachineRef,
	) (*core.TrashObject, *katapult.Response, error)

	Shutdown(
		ctx context.Context,
		ref core.VirtualMachineRef,
	) (*core.Task, *katapult.Response, error)

	Start(
		ctx context.Context,
		ref core.VirtualMachineRef,
	) (*core.Task, *katapult.Response, error)

	Stop(
		ctx context.Context,
		ref core.VirtualMachineRef,
	) (*core.Task, *katapult.Response, error)

	Reset(
		ctx context.Context,
		ref core.VirtualMachineRef,
	) (*core.Task, *katapult.Response, error)
}

func getVMRef(cmd *cobra.Command) (core.VirtualMachineRef, error) {
	id := cmd.Flag("id").Value.String()
	ref := core.VirtualMachineRef{ID: id}
	if id == "" {
		fqdn := cmd.Flag("fqdn").Value.String()
		if fqdn == "" {
			return core.VirtualMachineRef{}, fmt.Errorf("both ID and FQDN are unset")
		}
		return core.VirtualMachineRef{FQDN: fqdn}, nil
	}
	return ref, nil
}

func vmNotFoundHandlingError(err error) error {
	if errors.Is(err, core.ErrVirtualMachineNotFound) {
		return fmt.Errorf("unknown virtual machine")
	}
	return err
}

const virtualMachineListFormat = `{{ Table (StringSlice "Name" "FQDN") (MultipleRows . "Name" "FQDN") }}`

func virtualMachinesListCmd(client virtualMachinesClient) *cobra.Command {
	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Get a list of virtual machines from an organization",
		Long: "Get a list of virtual machines from an organization. By default, " +
			"the argument is used as the sub-domain and is used if the ID is not specified.",
		RunE: outputWrapper(func(cmd *cobra.Command, args []string) (Output, error) {
			ref := core.OrganizationRef{ID: cmd.Flags().Lookup("org-id").Value.String()}
			if ref.ID == "" {
				if len(args) == 0 || args[0] == "" {
					return nil, fmt.Errorf("both ID and subdomain are unset")
				}
				ref.SubDomain = args[0]
			}

			totalPages := 1
			allVms := make([]*core.VirtualMachine, 0)
			for pageNum := 1; pageNum <= totalPages; pageNum++ {
				vms, resp, err := client.List(
					cmd.Context(), ref, &core.ListOptions{Page: pageNum},
				)
				if err != nil {
					return nil, err
				}
				if resp.Pagination != nil {
					totalPages = resp.Pagination.TotalPages
				}
				allVms = append(allVms, vms...)
			}

			return &genericOutput{
				item:                allVms,
				defaultTextTemplate: virtualMachineListFormat,
			}, nil
		}),
	}
	list.Flags().String("org-id", "", "The ID of the organization. If set, this takes priority over the sub-domain.")
	return list
}

func virtualMachinesPoweroffCmd(client virtualMachinesClient) *cobra.Command {
	poweroff := &cobra.Command{
		Use:   "poweroff",
		Short: "Used to power off a virtual machine.",
		Long:  "Used to power off a virtual machine.",
		RunE: outputWrapper(func(cmd *cobra.Command, args []string) (Output, error) {
			ref, err := getVMRef(cmd)
			if err != nil {
				return nil, err
			}
			task, _, err := client.Shutdown(cmd.Context(), ref)
			if err != nil {
				return nil, vmNotFoundHandlingError(err)
			}
			return &genericOutput{
				item:                task,
				defaultTextTemplate: "Virtual machine successfully powered down.\n",
			}, nil
		}),
	}
	poweroff.Flags().String("id", "", "The ID of the server. If set, this takes priority over the FQDN.")
	poweroff.Flags().String("fqdn", "", "The FQDN of the server.")
	return poweroff
}

func virtualMachinesStartCmd(client virtualMachinesClient) *cobra.Command {
	start := &cobra.Command{
		Use:   "start",
		Short: "Used to start a virtual machine.",
		Long:  "Used to start a virtual machine.",
		RunE: outputWrapper(func(cmd *cobra.Command, args []string) (Output, error) {
			ref, err := getVMRef(cmd)
			if err != nil {
				return nil, err
			}
			task, _, err := client.Start(cmd.Context(), ref)
			if err != nil {
				return nil, vmNotFoundHandlingError(err)
			}
			return &genericOutput{
				item:                task,
				defaultTextTemplate: "Virtual machine successfully started.\n",
			}, nil
		}),
	}
	start.Flags().String("id", "", "The ID of the server. If set, this takes priority over the FQDN.")
	start.Flags().String("fqdn", "", "The FQDN of the server.")
	return start
}

func virtualMachinesStopCmd(client virtualMachinesClient) *cobra.Command {
	stop := &cobra.Command{
		Use:   "stop",
		Short: "Used to stop a virtual machine.",
		Long:  "Used to stop a virtual machine.",
		RunE: outputWrapper(func(cmd *cobra.Command, args []string) (Output, error) {
			ref, err := getVMRef(cmd)
			if err != nil {
				return nil, err
			}
			task, _, err := client.Stop(cmd.Context(), ref)
			if err != nil {
				return nil, vmNotFoundHandlingError(err)
			}
			return &genericOutput{
				item:                task,
				defaultTextTemplate: "Virtual machine successfully stopped.\n",
			}, nil
		}),
	}
	stop.Flags().String("id", "", "The ID of the server. If set, this takes priority over the FQDN.")
	stop.Flags().String("fqdn", "", "The FQDN of the server.")
	return stop
}

func virtualMachinesResetCmd(client virtualMachinesClient) *cobra.Command {
	reset := &cobra.Command{
		Use:   "reset",
		Short: "Used to reset a virtual machine.",
		Long:  "Used to reset a virtual machine.",
		RunE: outputWrapper(func(cmd *cobra.Command, args []string) (Output, error) {
			ref, err := getVMRef(cmd)
			if err != nil {
				return nil, err
			}
			task, _, err := client.Reset(cmd.Context(), ref)
			if err != nil {
				return nil, vmNotFoundHandlingError(err)
			}
			return &genericOutput{
				item:                task,
				defaultTextTemplate: "Virtual machine successfully reset.\n",
			}, nil
		}),
	}
	reset.Flags().String("id", "", "The ID of the server. If set, this takes priority over the FQDN.")
	reset.Flags().String("fqdn", "", "The FQDN of the server.")
	return reset
}

type virtualMachinePackagesClient interface {
	List(
		ctx context.Context,
		opts *core.ListOptions,
	) ([]*core.VirtualMachinePackage, *katapult.Response, error)
}

func listAllVMPackages(ctx context.Context,
	vmPackagesClient virtualMachinePackagesClient) ([]*core.VirtualMachinePackage, error) {
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

func listAllIPAddresses(ctx context.Context, org core.OrganizationRef,
	ipAddressesClient virtualMachineIPAddressesClient) ([]*core.IPAddress, error) {
	totalPages := 1
	allAddresses := make([]*core.IPAddress, 0)
	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		addresses, resp, err := ipAddressesClient.List(ctx, org, &core.ListOptions{Page: pageNum})
		if err != nil {
			return nil, err
		}
		if resp.Pagination != nil {
			totalPages = resp.Pagination.TotalPages
		}
		allAddresses = append(allAddresses, addresses...)
	}
	return allAddresses, nil
}

type virtualMachineDiskTemplatesClient interface {
	List(
		ctx context.Context,
		org core.OrganizationRef,
		opts *core.DiskTemplateListOptions,
	) ([]*core.DiskTemplate, *katapult.Response, error)
}

func listAllDiskTemplates(ctx context.Context, org core.OrganizationRef,
	diskTemplatesClient virtualMachineDiskTemplatesClient) ([]*core.DiskTemplate, error) {
	totalPages := 1
	allImages := make([]*core.DiskTemplate, 0)
	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		images, resp, err := diskTemplatesClient.List(
			ctx, org, &core.DiskTemplateListOptions{Page: pageNum, IncludeUniversal: true})
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

func listAllTags(ctx context.Context, org core.OrganizationRef, tagsClient tagsClient) ([]*core.Tag, error) {
	totalPages := 1
	allTags := make([]*core.Tag, 0)
	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		tags, resp, err := tagsClient.List(ctx, org, &core.ListOptions{Page: pageNum})
		if err != nil {
			return nil, err
		}
		if resp.Pagination != nil {
			totalPages = resp.Pagination.TotalPages
		}
		allTags = append(allTags, tags...)
	}
	return allTags, nil
}

func listAllSSHKeys(ctx context.Context, org core.OrganizationRef,
	sshKeysClient sshKeysListClient) ([]*core.AuthSSHKey, error) {
	totalPages := 1
	allKeys := make([]*core.AuthSSHKey, 0)
	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		keys, resp, err := sshKeysClient.List(ctx, org, &core.ListOptions{Page: pageNum})
		if err != nil {
			return nil, err
		}
		if resp.Pagination != nil {
			totalPages = resp.Pagination.TotalPages
		}
		allKeys = append(allKeys, keys...)
	}
	return allKeys, nil
}

func getStringIndex(needle string, haystack []string) int {
	for i, v := range haystack {
		if v == needle {
			return i
		}
	}
	return -1
}

func getArrayIndex(needle []string, haystack [][]string) int {
	for i, v := range haystack {
		if &v[0] == &needle[0] {
			return i
		}
	}
	return -1
}

type virtualMachineIPAddressesClient interface {
	List(
		ctx context.Context,
		org core.OrganizationRef,
		opts *core.ListOptions,
	) ([]*core.IPAddress, *katapult.Response, error)
}

type sshKeysListClient interface {
	List(
		ctx context.Context,
		org core.OrganizationRef,
		opts *core.ListOptions,
	) ([]*core.AuthSSHKey, *katapult.Response, error)
}

type tagsClient interface {
	List(
		ctx context.Context,
		org core.OrganizationRef,
		opts *core.ListOptions,
	) ([]*core.Tag, *katapult.Response, error)
}

type virtualMachinesBuilderClient interface {
	CreateFromSpec(
		ctx context.Context,
		org core.OrganizationRef,
		spec *buildspec.VirtualMachineSpec,
	) (*core.VirtualMachineBuild, *katapult.Response, error)
}

//nolint:funlen,gocyclo
func virtualMachinesCreateCmd(
	orgsClient organisationsListClient,
	dcsClient dataCentersClient,
	vmPackagesClient virtualMachinePackagesClient,
	diskTemplatesClient virtualMachineDiskTemplatesClient,
	ipAddressesClient virtualMachineIPAddressesClient,
	sshKeysClient sshKeysListClient,
	tagsClient tagsClient,
	vmBuilderClient virtualMachinesBuilderClient,
	terminal console.TerminalInterface,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Allows you to create a VM.",
		Long:  "Allows you to create a VM.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Accept argument from env var.

			// List the organizations.
			orgs, _, err := orgsClient.List(cmd.Context())
			if err != nil {
				return err
			}

			// Create a fuzzy searcher for organizations.
			orgRows := make([][]string, len(orgs))
			for i, org := range orgs {
				orgRows[i] = []string{org.Name, org.SubDomain}
			}
			orgArr := console.FuzzyTableSelector(
				"Which organization would you like to deploy the VM in?",
				[]string{"Name", "Subdomain"}, orgRows, cmd.InOrStdin(), terminal)
			index := getArrayIndex(orgArr, orgRows)
			org := orgs[index]

			// List the datacenters.
			dcs, _, err := dcsClient.List(cmd.Context())
			if err != nil {
				return err
			}

			// Create a fuzzy searcher for data centers.
			dcRows := make([][]string, len(dcs))
			for i, dc := range dcs {
				dcRows[i] = []string{dc.Name, dc.Country.Name}
			}
			dcArr := console.FuzzyTableSelector(
				"Which DC would you like to deploy the VM in?", []string{"Name", "Country"}, dcRows,
				cmd.InOrStdin(), terminal)
			index = getArrayIndex(dcArr, dcRows)
			dc := dcs[index]

			// List the packages.
			packages, err := listAllVMPackages(cmd.Context(), vmPackagesClient)
			if err != nil {
				return err
			}
			packageRows := make([][]string, len(packages))
			for i, packageItem := range packages {
				packageRows[i] = []string{
					packageItem.Name, strconv.Itoa(packageItem.CPUCores),
					strconv.Itoa(packageItem.MemoryInGB) + "GB",
				}
			}
			packageArr := console.FuzzyTableSelector(
				"Which VM package would you like?", []string{"Name", "CPU Cores", "RAM"}, packageRows,
				cmd.InOrStdin(), terminal)
			index = getArrayIndex(packageArr, packageRows)
			packageResult := packages[index]

			// Ask about the distribution.
			distributions, err := listAllDiskTemplates(
				cmd.Context(), core.OrganizationRef{ID: org.ID}, diskTemplatesClient)
			if err != nil {
				return err
			}
			distributionStrs := make([]string, len(distributions))
			for i, distribution := range distributions {
				distributionStrs[i] = distribution.Name
			}
			distributionStr := console.FuzzySelector(
				"Which distribution would you like?", distributionStrs,
				cmd.InOrStdin(), terminal)
			index = getStringIndex(distributionStr, distributionStrs)
			distribution := distributions[index]

			// Handle networking if there's IP addresses.
			ips, err := listAllIPAddresses(cmd.Context(), core.OrganizationRef{ID: org.ID}, ipAddressesClient)
			if err != nil {
				return err
			}
			allIps := ips
			ips = make([]*core.IPAddress, 0)
			for _, v := range allIps {
				if v.AllocationID == "" {
					ips = append(ips, v)
				}
			}
			selectedIps := []*core.IPAddress{}
			if len(ips) != 0 {
				ipRows := make([][]string, len(ips))
				for i, ip := range ips {
					ipRows[i] = []string{ip.Address, ip.ReverseDNS}
				}
				selectedIPRows := console.FuzzyTableMultiSelector(
					"Please select any IP addresses you wish to add.",
					[]string{"Address", "Reverse DNS"}, ipRows, cmd.InOrStdin(), terminal)
				selectedIps = make([]*core.IPAddress, len(selectedIPRows))
				for i, arr := range selectedIPRows {
					selectedIps[i] = ips[getArrayIndex(arr, ipRows)]
				}
			}

			// List the SSH keys.
			keys, err := listAllSSHKeys(cmd.Context(), core.OrganizationRef{ID: org.ID}, sshKeysClient)
			if err != nil {
				return err
			}
			keyIds := []string{}
			if len(keys) != 0 {
				keyRows := make([][]string, len(keys))
				for i, key := range keys {
					keyRows[i] = []string{key.Name, key.Fingerprint}
				}
				selectedKeys := console.FuzzyTableMultiSelector(
					"Which organization SSH keys do you wish to add?", []string{"Name", "Fingerprint"},
					keyRows, cmd.InOrStdin(), terminal)
				keyIds = make([]string, len(selectedKeys))
				for i, arr := range selectedKeys {
					keyIds[i] = keys[getArrayIndex(arr, keyRows)].ID
				}
			}

			// Ask for the tags.
			tags, err := listAllTags(cmd.Context(), core.OrganizationRef{ID: org.ID}, tagsClient)
			if err != nil {
				return err
			}
			tagStrs := make([]string, len(tags))
			for i, v := range tags {
				tagStrs[i] = v.Name
			}
			tagIds := []string{}
			if len(tags) != 0 {
				selectedTags := console.FuzzyMultiSelector(
					"Do you wish to add any tags?", tagStrs, cmd.InOrStdin(), terminal)
				tagIds = make([]string, len(selectedTags))
				for i, tagName := range selectedTags {
					tagIds[i] = tags[getStringIndex(tagName, selectedTags)].ID
				}
			}

			// Clear the terminal.
			goterm.Clear()
			goterm.Flush()

			// Get the buffered stdin.
			bufferedStdin := bufio.NewReader(cmd.InOrStdin())

			// Ask for the name.
			name := console.Question(
				"What would you like the virtual machine to be called?", false, bufferedStdin, cmd.OutOrStdout())

			// Ask for the hostname.
			hostname := console.Question(
				"If you want a hostname, what do you want it to be?", true, bufferedStdin, cmd.OutOrStdout())

			// Ask for the description.
			description := console.Question(
				"If you want a description, what do you want it to be?", true, bufferedStdin, cmd.OutOrStdout())

			// Build the virtual machine spec.
			ifaces := make([]*buildspec.NetworkInterface, len(selectedIps))
			for i, ip := range selectedIps {
				if ip.Network == nil {
					return errors.New("ip address not assigned to network")
				}
				ifaces[i] = &buildspec.NetworkInterface{
					IPAddressAllocations: []*buildspec.IPAddressAllocation{
						{
							IPAddress: &buildspec.IPAddress{ID: ip.ID},
							Type:      buildspec.ExistingIPAddressAllocation,
						},
					},
					Network: &buildspec.Network{ID: ip.Network.ID},
				}
			}
			spec := &buildspec.VirtualMachineSpec{
				DataCenter: &buildspec.DataCenter{ID: dc.ID},
				Resources:  &buildspec.Resources{Package: &buildspec.Package{ID: packageResult.ID}},
				DiskTemplate: &buildspec.DiskTemplate{ID: distribution.ID, Options: []*buildspec.DiskTemplateOption{
					{
						Key:   "install_agent",
						Value: "true",
					},
				}},
				NetworkInterfaces: ifaces,
				Hostname:          hostname,
				Name:              name,
				Description:       description,
				AuthorizedKeys:    &buildspec.AuthorizedKeys{SSHKeys: keyIds},
				Tags:              tagIds,
			}

			// ✨ Build the virtual machine.
			_, _, err = vmBuilderClient.CreateFromSpec(cmd.Context(), core.OrganizationRef{ID: org.ID}, spec)
			if err != nil {
				return err
			}

			// Return no errors.
			return nil
		},
	}

	// Return the command.
	return cmd
}

func virtualMachinesCmd(vmClient virtualMachinesClient,
	orgsClient organisationsListClient, dcsClient dataCentersClient,
	vmPackagesClient virtualMachinePackagesClient,
	diskTemplatesClient virtualMachineDiskTemplatesClient,
	ipAddressesClient virtualMachineIPAddressesClient,
	sshKeysClient sshKeysListClient,
	tagsClient tagsClient,
	vmBuilderClient virtualMachinesBuilderClient,
	terminal console.TerminalInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vm",
		Aliases: []string{"vms", "virtual-machines", "virtual_machines"},
		Short:   "Get information or do actions with virtual machines",
		Long:    "Get information or do actions with virtual machines.",
	}

	cmd.AddCommand(
		virtualMachinesListCmd(vmClient),
		virtualMachinesPoweroffCmd(vmClient),
		virtualMachinesStartCmd(vmClient),
		virtualMachinesStopCmd(vmClient),
		virtualMachinesResetCmd(vmClient),
		virtualMachinesCreateCmd(orgsClient, dcsClient, vmPackagesClient,
			diskTemplatesClient, ipAddressesClient, sshKeysClient,
			tagsClient, vmBuilderClient, terminal))

	return cmd
}
