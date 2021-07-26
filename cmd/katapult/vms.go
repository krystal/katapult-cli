package main

import (
	"context"
	"errors"
	"fmt"

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

func virtualMachinesListCmd(client virtualMachinesClient) *cobra.Command {
	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Get a list of virtual machines from an organization",
		Long:    "Get a list of virtual machines from an organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ref := core.OrganizationRef{ID: cmd.Flags().Lookup("id").Value.String()}
			if ref.ID == "" {
				subdomain := cmd.Flags().Lookup("subdomain").Value.String()
				if subdomain == "" {
					return fmt.Errorf("both ID and subdomain are unset")
				}
				ref.SubDomain = subdomain
			}

			stdout := cmd.OutOrStdout()

			page := 1
			perPage := 30
			for {
				vms, resp, err := client.List(cmd.Context(), ref, &core.ListOptions{
					Page:    page,
					PerPage: perPage,
				})
				if err != nil {
					return err
				}

				// Ensure all the data for the pagination is correct.
				pagination := resp.Pagination
				page = pagination.CurrentPage + 1
				perPage = pagination.PerPage
				end := pagination.TotalPages == 0 || pagination.CurrentPage == pagination.TotalPages

				for _, vm := range vms {
					fqdn := vm.FQDN
					if fqdn == "" {
						fqdn = "<no fqdn specified>"
					}
					_, _ = fmt.Fprintf(stdout, " - %s (%s) [%s]: %s\n", vm.Name, fqdn, vm.ID, vm.Package.Name)
				}

				if end {
					// We are done with the pages.
					break
				}
			}

			return nil
		},
	}
	list.Flags().String("id", "", "The ID of the organization. If set, this takes priority over the sub-domain.")
	list.Flags().String("subdomain", "", "The sub-domain of the organization.")
	return list
}

func virtualMachinesPoweroffCmd(client virtualMachinesClient) *cobra.Command {
	poweroff := &cobra.Command{
		Use:   "poweroff",
		Short: "Used to power off a virtual machine.",
		Long:  "Used to power off a virtual machine.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ref, err := getVMRef(cmd)
			if err != nil {
				return err
			}
			if _, _, err = client.Shutdown(cmd.Context(), ref); err != nil {
				return vmNotFoundHandlingError(err)
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Virtual machine successfully powered down.")
			return nil
		},
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
		RunE: func(cmd *cobra.Command, args []string) error {
			ref, err := getVMRef(cmd)
			if err != nil {
				return err
			}
			if _, _, err = client.Start(cmd.Context(), ref); err != nil {
				return vmNotFoundHandlingError(err)
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Virtual machine successfully started.")
			return nil
		},
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
		RunE: func(cmd *cobra.Command, args []string) error {
			ref, err := getVMRef(cmd)
			if err != nil {
				return err
			}
			if _, _, err = client.Stop(cmd.Context(), ref); err != nil {
				return vmNotFoundHandlingError(err)
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Virtual machine successfully stopped.")
			return nil
		},
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
		RunE: func(cmd *cobra.Command, args []string) error {
			ref, err := getVMRef(cmd)
			if err != nil {
				return err
			}
			if _, _, err = client.Reset(cmd.Context(), ref); err != nil {
				return vmNotFoundHandlingError(err)
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Virtual machine successfully reset.")
			return nil
		},
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
