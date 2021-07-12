package main

import (
	"fmt"
	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
	"github.com/spf13/cobra"
)

func vmFlag(cmd *cobra.Command) (core.VirtualMachineRef, error) {
	id := cmd.Flag("id").Value.String()
	ref := core.VirtualMachineRef{ID: id}
	if id == "" {
		fqdn := cmd.Flag("fqdn").Value.String()
		if fqdn == "" {
			return core.VirtualMachineRef{}, fmt.Errorf("both ID and FQDN are unset")
		}
		return core.VirtualMachineRef{FQDN: fqdn},  nil
	}
	return ref, nil
}

func vmNotFoundHandlingError(err error, resp *katapult.Response) error {
	if resp != nil && resp.StatusCode == 404 {
		return fmt.Errorf("unknown virtual machine")
	}
	return err
}

func virtualMachinesCmd(client *katapult.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vm",
		Aliases: []string{"vms", "virtual-machines", "virtual_machines"},
		Short:   "Get information or do actions with virtual machines",
		Long:    "Get information or do actions with virtual machines.",
	}

	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Get a list of virtual machines from an organisation",
		Long:    "Get a list of virtual machines from an organisation.",
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
				vms, resp, err := core.NewVirtualMachinesClient(client).List(cmd.Context(), ref, &core.ListOptions{
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
	list.Flags().String("id", "", "The ID of the organisation. If set, this takes priority over the sub-domain.")
	list.Flags().String("subdomain", "", "The sub-domain of the organisation.")
	cmd.AddCommand(list)

	poweroff := &cobra.Command{
		Use:   "poweroff",
		Short: "Used to power off a virtual machine.",
		Long:  "Used to power off a virtual machine.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ref, err := vmFlag(cmd)
			if err != nil {
				return err
			}
			_, resp, err := core.NewVirtualMachinesClient(client).Shutdown(cmd.Context(), ref)
			if err != nil {
				return vmNotFoundHandlingError(err, resp)
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Virtual machine successfully powered down.")
			return nil
		},
	}
	cmd.AddCommand(poweroff)

	start := &cobra.Command{
		Use:   "start",
		Short: "Used to start a virtual machine.",
		Long:  "Used to start a virtual machine.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ref, err := vmFlag(cmd)
			if err != nil {
				return err
			}
			_, resp, err := core.NewVirtualMachinesClient(client).Start(cmd.Context(), ref)
			if err != nil {
				return vmNotFoundHandlingError(err, resp)
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Virtual machine successfully started.")
			return nil
		},
	}
	cmd.AddCommand(start)

	stop := &cobra.Command{
		Use:   "stop",
		Short: "Used to stop a virtual machine.",
		Long:  "Used to stop a virtual machine.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ref, err := vmFlag(cmd)
			if err != nil {
				return err
			}
			_, resp, err := core.NewVirtualMachinesClient(client).Stop(cmd.Context(), ref)
			if err != nil {
				return vmNotFoundHandlingError(err, resp)
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Virtual machine successfully stopped.")
			return nil
		},
	}
	cmd.AddCommand(stop)

	reset := &cobra.Command{
		Use:   "reset",
		Short: "Used to reset a virtual machine.",
		Long:  "Used to reset a virtual machine.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ref, err := vmFlag(cmd)
			if err != nil {
				return err
			}
			_, resp, err := core.NewVirtualMachinesClient(client).Reset(cmd.Context(), ref)
			if err != nil {
				return vmNotFoundHandlingError(err, resp)
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Virtual machine successfully reset.")
			return nil
		},
	}
	cmd.AddCommand(reset)

	return cmd
}
