
package main

import (
	"bytes"
	"context"
	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

const stdoutVMsList = ` - test (test.example.com) [1]: test
`

type vmPages [][]*core.VirtualMachine

type vmsClient struct {
	// Defines the power state of the VM's.
	powerStates map[string]bool

	// Defines the organisation ID -> vmPages.
	organisationIDPages map[string]vmPages

	// Defines the organisation subdomain -> vmPages.
	organizationSubdomainPages map[string]vmPages
}

// Used to toggle the power state and return the old result.
func (v *vmsClient) togglePowerState(key string, fqdn bool) bool {
	if v.powerStates == nil {
		v.powerStates = map[string]bool{}
	}
	prefix := "i"
	if fqdn {
		prefix = "s"
	}
	state, ok := v.powerStates[prefix + key]
	if !ok {
		// A VM starts powered up.
		state = true
	}
	v.powerStates[prefix + key] = !state
	return state
}

func (v *vmsClient) List(
	_ context.Context,
	org core.OrganizationRef,
	opts *core.ListOptions,
) ([]*core.VirtualMachine, *katapult.Response, error) {
	// Defines the pages.
	var pages vmPages
	switch {
	case org.ID != "":
		pages = v.organisationIDPages[org.ID]
	case org.SubDomain != "":
		pages = v.organizationSubdomainPages[org.SubDomain]
	default:
		return nil, nil, katapult.ErrNotFound
	}

	// Get the VM page.
	if opts.Page > len(pages) {
		return nil, nil, katapult.ErrNotFound
	}
	page := pages[opts.Page-1]

	// Return the pages.
	totalPages := 0
	for _, v := range pages {
		totalPages += len(v)
	}
	return page, &katapult.Response{Pagination: &katapult.Pagination{
		CurrentPage: opts.Page,
		TotalPages:  len(pages),
		Total:       totalPages,
		PerPage:     len(page),
	}}, nil
}

func (v *vmsClient) Shutdown(
	_ context.Context,
	ref core.VirtualMachineRef,
) (*core.Task, *katapult.Response, error) {
	// Get the key and if it's an FQDN.
	fqdn := false
	id := ref.ID
	if id == "" {
		fqdn = true
		id = ref.FQDN
	}

	// Toggle the power state.
	poweredOn := v.togglePowerState(id, fqdn)
	if !poweredOn {
		// The VM wasn't powered on.
		return nil, nil, katapult.ErrNotAcceptable
	}
	return &core.Task{}, nil, nil
}

func (v *vmsClient) Stop(
	_ context.Context,
	ref core.VirtualMachineRef,
) (*core.Task, *katapult.Response, error) {
	// Basically acts the same as far as event mocking logic goes.
	return v.Shutdown(nil, ref)
}

func (v *vmsClient) Start(
	_ context.Context,
	ref core.VirtualMachineRef,
) (*core.Task, *katapult.Response, error) {
	// Get the key and if it's an FQDN.
	fqdn := false
	id := ref.ID
	if id == "" {
		fqdn = true
		id = ref.FQDN
	}

	// Toggle the power state.
	poweredOn := v.togglePowerState(id, fqdn)
	if poweredOn {
		// The VM was powered off.
		return nil, nil, katapult.ErrNotAcceptable
	}
	return &core.Task{}, nil, nil
}

func (v *vmsClient) Reset(
	_ context.Context,
	ref core.VirtualMachineRef,
) (*core.Task, *katapult.Response, error) {
	// Basically acts the same as far as event mocking logic goes.
	return v.Start(nil, ref)
}

func (v *vmsClient) Delete(
	_ context.Context,
	ref core.VirtualMachineRef,
) (*core.TrashObject, *katapult.Response, error) {
	switch {
	case ref.ID != "":
		delete(v.powerStates, "i" + ref.ID)
	case ref.FQDN != "":
		delete(v.powerStates, "s" + ref.FQDN)
	}
	return nil, nil, nil
}

func TestVMs_List(t *testing.T) {
	cmd := virtualMachinesCmd(&vmsClient{
		organisationIDPages: map[string]vmPages{
			"1": {
				{
					{
						ID:                  "1",
						Name:                "test",
						Hostname:            "test.example.com",
						FQDN:                "test.example.com",
						Description:         "test",
						Package: 			&core.VirtualMachinePackage{Name: "test"},
					},
				},
			},
		},
	})
	stdout := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetArgs([]string{"list", "--id=1"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, stdoutVMsList, stdout.String())
}

func TestVMs_Poweroff(t *testing.T) {
	cmd := virtualMachinesCmd(&vmsClient{})
	stdout := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetArgs([]string{"poweroff", "--id=1"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Virtual machine successfully powered down.\n", stdout.String())
}

func TestVMs_Stop(t *testing.T) {
	cmd := virtualMachinesCmd(&vmsClient{})
	stdout := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetArgs([]string{"stop", "--id=1"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Virtual machine successfully stopped.\n", stdout.String())
}

func TestVMs_Start(t *testing.T) {
	client := &vmsClient{}
	client.togglePowerState("1", false)
	cmd := virtualMachinesCmd(client)
	stdout := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetArgs([]string{"start", "--id=1"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Virtual machine successfully started.\n", stdout.String())
}

func TestVMs_Reset(t *testing.T) {
	client := &vmsClient{}
	client.togglePowerState("1", false)
	cmd := virtualMachinesCmd(client)
	stdout := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetArgs([]string{"reset", "--id=1"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Virtual machine successfully reset.\n", stdout.String())
}
