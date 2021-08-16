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

//go:embed formatdata/vm/list.txt
var virtualMachineListFormat string

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

			return genericOutput{
				item: allVms,
				tpl:  virtualMachineListFormat,
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
			return genericOutput{
				item: task,
				tpl:  "Virtual machine successfully powered down.\n",
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
			return genericOutput{
				item: task,
				tpl:  "Virtual machine successfully started.\n",
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
			return genericOutput{
				item: task,
				tpl:  "Virtual machine successfully stopped.\n",
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
			return genericOutput{
				item: task,
				tpl:  "Virtual machine successfully reset.\n",
			}, nil
		}),
	}
	reset.Flags().String("id", "", "The ID of the server. If set, this takes priority over the FQDN.")
	reset.Flags().String("fqdn", "", "The FQDN of the server.")
	return reset
}

func virtualMachinesCmd(client virtualMachinesClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vm",
		Aliases: []string{"vms", "virtual-machines", "virtual_machines"},
		Short:   "Get information or do actions with virtual machines",
		Long:    "Get information or do actions with virtual machines.",
	}

	cmd.AddCommand(
		virtualMachinesListCmd(client),
		virtualMachinesPoweroffCmd(client),
		virtualMachinesStartCmd(client),
		virtualMachinesStopCmd(client),
		virtualMachinesResetCmd(client))

	return cmd
}
