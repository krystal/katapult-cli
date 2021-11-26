package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/augurysys/timestamp"
	"github.com/krystal/go-katapult/buildspec"
	"github.com/krystal/katapult-cli/internal/golden"
	"github.com/krystal/katapult-cli/internal/keystrokes"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
	"github.com/krystal/katapult-cli/cmd/katapult/console"
	"github.com/stretchr/testify/assert"
)

type vmPages [][]*core.VirtualMachine

type vmsClient struct {
	// Defines a ID that should be not found.
	idNotFound string

	// Defines an FQDN that should be not found.
	fqdnNotFound string

	// Defines the power state of the VM's.
	powerStates map[string]bool

	// Defines the organization ID -> vmPages.
	organizationIDPages map[string]vmPages

	// Defines the organization subdomain -> vmPages.
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
	state, ok := v.powerStates[prefix+key]
	if !ok {
		// A VM starts powered up.
		state = true
	}
	v.powerStates[prefix+key] = !state
	return state
}

func (v *vmsClient) List(_ context.Context, org core.OrganizationRef, opts *core.ListOptions) (
	[]*core.VirtualMachine, *katapult.Response, error,
) {
	// Defines the pages.
	var pages vmPages
	switch {
	case org.ID != "":
		pages = v.organizationIDPages[org.ID]
	case org.SubDomain != "":
		pages = v.organizationSubdomainPages[org.SubDomain]
	default:
		return nil, nil, core.ErrOrganizationNotFound
	}

	// Get the VM page.
	if opts.Page > len(pages) {
		return nil, nil, katapult.ErrNotFound
	}
	page := pages[opts.Page-1]

	// Return the items.
	totalItems := 0
	for _, v := range pages {
		totalItems += len(v)
	}
	return page, &katapult.Response{Pagination: &katapult.Pagination{
		CurrentPage: opts.Page,
		TotalPages:  len(pages),
		Total:       totalItems,
		PerPage:     len(page),
	}}, nil
}

func (v *vmsClient) ensureFound(ref core.VirtualMachineRef) error {
	if (ref.FQDN != "" && v.fqdnNotFound == ref.FQDN) || (ref.ID != "" && v.idNotFound == ref.ID) {
		return core.ErrVirtualMachineNotFound
	}
	return nil
}

func (v *vmsClient) Shutdown(_ context.Context, ref core.VirtualMachineRef) (*core.Task, *katapult.Response, error) {
	// Pre-execution checks.
	if err := v.ensureFound(ref); err != nil {
		return nil, nil, err
	}

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
		return nil, nil, core.NewTaskQueueingError(&katapult.ResponseError{
			Code:        "task_queueing_error",
			Description: "VM was not powered on",
		})
	}
	return &core.Task{}, nil, nil
}

func (v *vmsClient) Stop(_ context.Context, ref core.VirtualMachineRef) (*core.Task, *katapult.Response, error) {
	// Basically acts the same as far as event mocking logic goes.
	return v.Shutdown(context.TODO(), ref)
}

func (v *vmsClient) Start(_ context.Context, ref core.VirtualMachineRef) (*core.Task, *katapult.Response, error) {
	// Pre-execution checks.
	if err := v.ensureFound(ref); err != nil {
		return nil, nil, err
	}

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
		// The VM was powered on.
		return nil, nil, core.NewTaskQueueingError(&katapult.ResponseError{
			Code:        "task_queueing_error",
			Description: "VM was powered on",
		})
	}
	return &core.Task{}, nil, nil
}

func (v *vmsClient) Reset(_ context.Context, ref core.VirtualMachineRef) (*core.Task, *katapult.Response, error) {
	// Pre-execution checks.
	if err := v.ensureFound(ref); err != nil {
		return nil, nil, err
	}

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
		// The VM needs to be on to be reset.
		return nil, nil, core.NewTaskQueueingError(&katapult.ResponseError{
			Code:        "task_queueing_error",
			Description: "VM was not powered off",
		})
	}
	v.togglePowerState(id, fqdn)
	return nil, nil, nil
}

func (v *vmsClient) Delete(_ context.Context, ref core.VirtualMachineRef) (
	*core.TrashObject, *katapult.Response, error,
) {
	// Pre-execution checks.
	if err := v.ensureFound(ref); err != nil {
		return nil, nil, err
	}

	switch {
	case ref.ID != "":
		delete(v.powerStates, "i"+ref.ID)
	case ref.FQDN != "":
		delete(v.powerStates, "s"+ref.FQDN)
	}
	return nil, nil, nil
}

func TestVMs_List(t *testing.T) {
	tests := []struct {
		name string

		id         map[string]vmPages
		subdomains map[string]vmPages

		args    []string
		stderr  string
		wantErr string
	}{
		{
			name:    "no ID/subdomain provided",
			args:    []string{"list"},
			wantErr: "both ID and subdomain are unset",
		},
		{
			name: "test un-paginated vm ID list",
			id: map[string]vmPages{
				"1": {
					{
						{
							ID:          "vm_rrmEoG6CKUX0IKgX",
							Name:        "My Blog",
							Hostname:    "my-blog",
							FQDN:        "my-blog.acme-labs.katapult.cloud",
							Description: "test",
							Package:     &core.VirtualMachinePackage{Name: "test"},
						},
					},
				},
			},
			args: []string{"list", "--org-id=1"},
		},
		{
			name: "test un-paginated vm subdomain list",
			subdomains: map[string]vmPages{
				"1": {
					{
						{
							ID:          "vm_rrmEoG6CKUX0IKgX",
							Name:        "My Blog",
							Hostname:    "my-blog",
							FQDN:        "my-blog.acme-labs.katapult.cloud",
							Description: "test",
							Package:     &core.VirtualMachinePackage{Name: "test"},
						},
					},
				},
			},
			args: []string{"list", "1"},
		},
		{
			name: "test paginated vm list",
			subdomains: map[string]vmPages{
				"1": {
					{
						{
							ID:          "0",
							Name:        "test",
							Hostname:    "test.example.com",
							FQDN:        "test.example.com",
							Description: "test",
							Package:     &core.VirtualMachinePackage{Name: "test"},
						},
						{
							ID:          "1",
							Name:        "test1",
							Hostname:    "test1.example.com",
							FQDN:        "test1.example.com",
							Description: "test1",
							Package:     &core.VirtualMachinePackage{Name: "test1"},
						},
					},
					{
						{
							ID:          "2",
							Name:        "test2",
							Hostname:    "test2.example.com",
							FQDN:        "test2.example.com",
							Description: "test2",
							Package:     &core.VirtualMachinePackage{Name: "test2"},
						},
					},
				},
			},
			args: []string{"list", "1"},
		},
		{
			name:    "test not found",
			args:    []string{"list", "not_exists"},
			wantErr: "katapult: not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := virtualMachinesCmd(
				&vmsClient{organizationIDPages: tt.id, organizationSubdomainPages: tt.subdomains}, nil,
				nil, nil, nil, nil, nil,
				nil, nil, nil, nil)
			cmd.SetArgs(tt.args)
			assertCobraCommand(t, cmd, tt.wantErr, tt.stderr)
		})
	}
}

func TestVMs_Poweroff(t *testing.T) {
	tests := []struct {
		name string

		idNotFound   string
		fqdnNotFound string

		args        []string
		poweredDown *struct {
			key  string
			fqdn bool
		}
		stderr   string
		wantErr  string
		validate func(client *vmsClient) string
	}{
		{
			name:    "no ID/FQDN provided",
			args:    []string{"poweroff"},
			wantErr: "both ID and FQDN are unset",
		},
		{
			name: "test normal power down by ID",
			args: []string{"poweroff", "--id=1"},
			validate: func(client *vmsClient) string {
				state, ok := client.powerStates["i1"]
				if !ok {
					return "virtual machine not stopped"
				}
				if state {
					return "virtual machine powered on"
				}
				return ""
			},
		},
		{
			name: "test normal power down by fqdn",
			args: []string{"poweroff", "--fqdn=1"},
			validate: func(client *vmsClient) string {
				state, ok := client.powerStates["s1"]
				if !ok {
					return "virtual machine not stopped"
				}
				if state {
					return "virtual machine powered on"
				}
				return ""
			},
		},
		{
			name: "test power down of already powered down FQDN",
			args: []string{"poweroff", "--fqdn=1"},
			poweredDown: &struct {
				key  string
				fqdn bool
			}{key: "1", fqdn: true},
			wantErr: "katapult: not_acceptable: task_queueing_error: VM was not powered on",
		},
		{
			name: "test power down of already powered down ID",
			args: []string{"poweroff", "--id=1"},
			poweredDown: &struct {
				key  string
				fqdn bool
			}{key: "1", fqdn: false},
			wantErr: "katapult: not_acceptable: task_queueing_error: VM was not powered on",
		},
		{
			name:         "test fqdn not found",
			fqdnNotFound: "not_exists",
			args:         []string{"poweroff", "--fqdn=not_exists"},
			wantErr:      "unknown virtual machine",
		},
		{
			name:       "test id not found",
			idNotFound: "not_exists",
			args:       []string{"poweroff", "--id=not_exists"},
			wantErr:    "unknown virtual machine",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &vmsClient{fqdnNotFound: tt.fqdnNotFound, idNotFound: tt.idNotFound}
			if tt.poweredDown != nil {
				client.togglePowerState(tt.poweredDown.key, tt.poweredDown.fqdn)
			}
			cmd := virtualMachinesCmd(client, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
			cmd.SetArgs(tt.args)
			assertCobraCommand(t, cmd, tt.wantErr, tt.stderr)
			if tt.validate != nil {
				assert.Equal(t, "", tt.validate(client))
			}
		})
	}
}

func TestVMs_Stop(t *testing.T) {
	tests := []struct {
		name string

		idNotFound   string
		fqdnNotFound string

		args        []string
		poweredDown *struct {
			key  string
			fqdn bool
		}
		stderr   string
		wantErr  string
		validate func(client *vmsClient) string
	}{
		{
			name:    "no ID/FQDN provided",
			args:    []string{"stop"},
			wantErr: "both ID and FQDN are unset",
		},
		{
			name: "test normal stop by ID",
			args: []string{"stop", "--id=1"},
			validate: func(client *vmsClient) string {
				state, ok := client.powerStates["i1"]
				if !ok {
					return "virtual machine not stopped"
				}
				if state {
					return "virtual machine powered on"
				}
				return ""
			},
		},
		{
			name: "test normal stop by fqdn",
			args: []string{"stop", "--fqdn=1"},
			validate: func(client *vmsClient) string {
				state, ok := client.powerStates["s1"]
				if !ok {
					return "virtual machine not stopped"
				}
				if state {
					return "virtual machine powered on"
				}
				return ""
			},
		},
		{
			name: "test stop of already powered down FQDN",
			args: []string{"stop", "--fqdn=1"},
			poweredDown: &struct {
				key  string
				fqdn bool
			}{key: "1", fqdn: true},
			wantErr: "katapult: not_acceptable: task_queueing_error: VM was not powered on",
		},
		{
			name: "test stop of already powered down ID",
			args: []string{"stop", "--id=1"},
			poweredDown: &struct {
				key  string
				fqdn bool
			}{key: "1", fqdn: false},
			wantErr: "katapult: not_acceptable: task_queueing_error: VM was not powered on",
		},
		{
			name:         "test fqdn not found",
			fqdnNotFound: "not_exists",
			args:         []string{"stop", "--fqdn=not_exists"},
			wantErr:      "unknown virtual machine",
		},
		{
			name:       "test id not found",
			idNotFound: "not_exists",
			args:       []string{"stop", "--id=not_exists"},
			wantErr:    "unknown virtual machine",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &vmsClient{idNotFound: tt.idNotFound, fqdnNotFound: tt.fqdnNotFound}
			if tt.poweredDown != nil {
				client.togglePowerState(tt.poweredDown.key, tt.poweredDown.fqdn)
			}
			cmd := virtualMachinesCmd(client, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
			cmd.SetArgs(tt.args)
			assertCobraCommand(t, cmd, tt.wantErr, tt.stderr)
			if tt.validate != nil {
				assert.Equal(t, "", tt.validate(client))
			}
		})
	}
}

func TestVMs_Start(t *testing.T) {
	tests := []struct {
		name string

		idNotFound   string
		fqdnNotFound string

		args        []string
		poweredDown *struct {
			key  string
			fqdn bool
		}
		stderr   string
		wantErr  string
		validate func(client *vmsClient) string
	}{
		{
			name:    "no ID/FQDN provided",
			args:    []string{"start"},
			wantErr: "both ID and FQDN are unset",
		},
		{
			name: "test normal start by ID",
			args: []string{"start", "--id=1"},
			poweredDown: &struct {
				key  string
				fqdn bool
			}{key: "1", fqdn: false},
			validate: func(client *vmsClient) string {
				state, ok := client.powerStates["i1"]
				if !ok {
					return "virtual machine not stopped"
				}
				if !state {
					return "virtual machine powered off"
				}
				return ""
			},
		},
		{
			name: "test normal start by FQDN",
			args: []string{"start", "--fqdn=1"},
			poweredDown: &struct {
				key  string
				fqdn bool
			}{key: "1", fqdn: true},
			validate: func(client *vmsClient) string {
				state, ok := client.powerStates["s1"]
				if !ok {
					return "virtual machine not stopped"
				}
				if !state {
					return "virtual machine powered off"
				}
				return ""
			},
		},
		{
			name:    "test start of already powered up FQDN",
			args:    []string{"start", "--fqdn=1"},
			wantErr: "katapult: not_acceptable: task_queueing_error: VM was powered on",
		},
		{
			name:    "test start of already powered up ID",
			args:    []string{"start", "--id=1"},
			wantErr: "katapult: not_acceptable: task_queueing_error: VM was powered on",
		},
		{
			name:         "test fqdn not found",
			fqdnNotFound: "not_exists",
			args:         []string{"poweroff", "--fqdn=not_exists"},
			wantErr:      "unknown virtual machine",
		},
		{
			name:       "test id not found",
			idNotFound: "not_exists",
			args:       []string{"poweroff", "--id=not_exists"},
			wantErr:    "unknown virtual machine",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &vmsClient{fqdnNotFound: tt.fqdnNotFound, idNotFound: tt.idNotFound}
			if tt.poweredDown != nil {
				client.togglePowerState(tt.poweredDown.key, tt.poweredDown.fqdn)
			}
			cmd := virtualMachinesCmd(client, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
			cmd.SetArgs(tt.args)
			assertCobraCommand(t, cmd, tt.wantErr, tt.stderr)
			if tt.validate != nil {
				assert.Equal(t, "", tt.validate(client))
			}
		})
	}
}

func TestVMs_Reset(t *testing.T) {
	tests := []struct {
		name string

		idNotFound   string
		fqdnNotFound string

		args        []string
		poweredDown *struct {
			key  string
			fqdn bool
		}
		stderr   string
		wantErr  string
		validate func(client *vmsClient) string
	}{
		{
			name:    "no ID/FQDN provided",
			args:    []string{"reset"},
			wantErr: "both ID and FQDN are unset",
		},
		{
			name: "test normal reset by ID",
			args: []string{"reset", "--id=1"},
			validate: func(client *vmsClient) string {
				state, ok := client.powerStates["i1"]
				if !ok {
					return "virtual machine not stopped"
				}
				if !state {
					return "virtual machine powered off"
				}
				return ""
			},
		},
		{
			name: "test normal reset by FQDN",
			args: []string{"reset", "--fqdn=1"},
			validate: func(client *vmsClient) string {
				state, ok := client.powerStates["s1"]
				if !ok {
					return "virtual machine not stopped"
				}
				if !state {
					return "virtual machine powered off"
				}
				return ""
			},
		},
		{
			name: "test reset of already powered down FQDN",
			args: []string{"reset", "--fqdn=1"},
			poweredDown: &struct {
				key  string
				fqdn bool
			}{key: "1", fqdn: true},
			wantErr: "katapult: not_acceptable: task_queueing_error: VM was not powered off",
		},
		{
			name: "test reset of already powered down ID",
			args: []string{"reset", "--id=1"},
			poweredDown: &struct {
				key  string
				fqdn bool
			}{key: "1", fqdn: false},
			wantErr: "katapult: not_acceptable: task_queueing_error: VM was not powered off",
		},
		{
			name:         "test fqdn not found",
			fqdnNotFound: "not_exists",
			args:         []string{"poweroff", "--fqdn=not_exists"},
			wantErr:      "unknown virtual machine",
		},
		{
			name:       "test id not found",
			idNotFound: "not_exists",
			args:       []string{"poweroff", "--id=not_exists"},
			wantErr:    "unknown virtual machine",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &vmsClient{fqdnNotFound: tt.fqdnNotFound, idNotFound: tt.idNotFound}
			if tt.poweredDown != nil {
				client.togglePowerState(tt.poweredDown.key, tt.poweredDown.fqdn)
			}
			cmd := virtualMachinesCmd(client, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
			cmd.SetArgs(tt.args)
			assertCobraCommand(t, cmd, tt.wantErr, tt.stderr)
			if tt.validate != nil {
				assert.Equal(t, "", tt.validate(client))
			}
		})
	}
}

type sshPages [][]*core.AuthSSHKey

type mockSSHKeysClient struct {
	// If an error is wanted to be thrown from the client, it is set here.
	throws string

	// Defines the organization ID -> sshPages.
	organizationIDPages map[string]sshPages

	// Defines the organization subdomain -> sshPages.
	organizationSubdomainPages map[string]sshPages
}

func (v mockSSHKeysClient) List(
	_ context.Context, org core.OrganizationRef,
	opts *core.ListOptions) ([]*core.AuthSSHKey, *katapult.Response, error) {
	// Handle throwing errors.
	if v.throws != "" {
		return nil, nil, errors.New(v.throws)
	}

	// Defines the pages.
	var pages sshPages
	switch {
	case org.ID != "":
		pages = v.organizationIDPages[org.ID]
	case org.SubDomain != "":
		pages = v.organizationSubdomainPages[org.SubDomain]
	default:
		return nil, nil, core.ErrOrganizationNotFound
	}

	// Get the SSH key page.
	if opts.Page > len(pages) {
		return nil, nil, katapult.ErrNotFound
	}
	page := pages[opts.Page-1]

	// Return the items.
	totalItems := 0
	for _, v := range pages {
		totalItems += len(v)
	}
	return page, &katapult.Response{Pagination: &katapult.Pagination{
		CurrentPage: opts.Page,
		TotalPages:  len(pages),
		Total:       totalItems,
		PerPage:     len(page),
	}}, nil
}

type tagPages [][]*core.Tag

type mockTagsClient struct {
	// If an error is wanted to be thrown from the client, it is set here.
	throws string

	// Defines the organization ID -> sshPages.
	organizationIDPages map[string]tagPages

	// Defines the organization subdomain -> tagPages.
	organizationSubdomainPages map[string]tagPages
}

func (v mockTagsClient) List(_ context.Context, org core.OrganizationRef,
	opts *core.ListOptions) ([]*core.Tag, *katapult.Response, error) {
	// Handle throwing errors.
	if v.throws != "" {
		return nil, nil, errors.New(v.throws)
	}

	// Defines the pages.
	var pages tagPages
	switch {
	case org.ID != "":
		pages = v.organizationIDPages[org.ID]
	case org.SubDomain != "":
		pages = v.organizationSubdomainPages[org.SubDomain]
	default:
		return nil, nil, core.ErrOrganizationNotFound
	}

	// Get the tag page.
	if opts.Page > len(pages) {
		return nil, nil, katapult.ErrNotFound
	}
	page := pages[opts.Page-1]

	// Return the items.
	totalItems := 0
	for _, v := range pages {
		totalItems += len(v)
	}
	return page, &katapult.Response{Pagination: &katapult.Pagination{
		CurrentPage: opts.Page,
		TotalPages:  len(pages),
		Total:       totalItems,
		PerPage:     len(page),
	}}, nil
}

type mockVMPackagesClient struct {
	packages []*core.VirtualMachinePackage
	throws   string
}

func (m mockVMPackagesClient) List(context.Context, *core.ListOptions) (
	[]*core.VirtualMachinePackage, *katapult.Response, error) {
	if m.throws != "" {
		return nil, nil, errors.New(m.throws)
	}
	return m.packages, &katapult.Response{Pagination: &katapult.Pagination{
		CurrentPage: 1, TotalPages: 1, Total: len(m.packages),
	}}, nil
}

type mockDiskTemplatesClient struct {
	ref           core.OrganizationRef
	diskTemplates []*core.DiskTemplate
	throws        string
}

func (m mockDiskTemplatesClient) List(_ context.Context, org core.OrganizationRef,
	opts *core.DiskTemplateListOptions) ([]*core.DiskTemplate, *katapult.Response, error) {
	if m.throws != "" {
		return nil, nil, errors.New(m.throws)
	}
	if m.ref.ID != org.ID || m.ref.SubDomain != org.SubDomain {
		return nil, nil, fmt.Errorf("ref mismatch: expected: %s got: %s", m.ref, org)
	}
	if !opts.IncludeUniversal {
		return nil, nil, errors.New("should include universal")
	}
	return m.diskTemplates, &katapult.Response{Pagination: &katapult.Pagination{
		CurrentPage: 1, TotalPages: 1, Total: len(m.diskTemplates),
	}}, nil
}

type ipPages [][]*core.IPAddress

type mockIPAddressClient struct {
	// If an error is wanted to be thrown from the client, it is set here.
	throws string

	// Defines the organization ID -> ipPages.
	organizationIDPages map[string]ipPages

	// Defines the organization subdomain -> ipPages.
	organizationSubdomainPages map[string]ipPages
}

func (v mockIPAddressClient) List(_ context.Context, org core.OrganizationRef,
	opts *core.ListOptions) ([]*core.IPAddress, *katapult.Response, error) {
	// Handle throwing errors.
	if v.throws != "" {
		return nil, nil, errors.New(v.throws)
	}

	// Defines the pages.
	var pages ipPages
	switch {
	case org.ID != "":
		pages = v.organizationIDPages[org.ID]
	case org.SubDomain != "":
		pages = v.organizationSubdomainPages[org.SubDomain]
	default:
		return nil, nil, core.ErrOrganizationNotFound
	}

	// Get the tag page.
	if opts.Page > len(pages) {
		return nil, nil, katapult.ErrNotFound
	}
	page := pages[opts.Page-1]

	// Return the items.
	totalItems := 0
	for _, v := range pages {
		totalItems += len(v)
	}
	return page, &katapult.Response{Pagination: &katapult.Pagination{
		CurrentPage: opts.Page,
		TotalPages:  len(pages),
		Total:       totalItems,
		PerPage:     len(page),
	}}, nil
}

type mockVMBuilderClient struct {
	throws string

	OrgResult  core.OrganizationRef
	SpecResult *buildspec.VirtualMachineSpec
}

func (m *mockVMBuilderClient) CreateFromSpec(_ context.Context, org core.OrganizationRef,
	spec *buildspec.VirtualMachineSpec) (*core.VirtualMachineBuild, *katapult.Response, error) {
	if m.throws != "" {
		return nil, nil, errors.New(m.throws)
	}
	m.OrgResult = org
	m.SpecResult = spec
	return nil, nil, nil
}

var successPackages = []*core.VirtualMachinePackage{
	{ID: "DO_NOT_PICK_IGNORE_THIS_ONE"},
	{
		ID:            "vmpkg_9UVoPiUQoI1cqtRd",
		Name:          "Test",
		Permalink:     "testing",
		CPUCores:      100,
		IPv4Addresses: 10,
		MemoryInGB:    1000,
		StorageInGB:   20,
	},
}

var successDiskTemplates = []*core.DiskTemplate{
	{ID: "DO_NOT_PICK_IGNORE_THIS_ONE"},
	{
		ID:          "disk_9UVoPiUQoI1cqtRd",
		Name:        "Ubuntu 20.04",
		Description: "testing",
		Permalink:   "ubuntu-20-04",
		Universal:   true,
		LatestVersion: &core.DiskTemplateVersion{
			ID:       "versopn+9UVoPiUQoI1cqtRd",
			Number:   1,
			Stable:   true,
			SizeInGB: 5,
		},
		OperatingSystem: &core.OperatingSystem{
			ID:   "ubuntu",
			Name: "Ubuntu",
		},
	},
}

var mockIPPages = ipPages{
	{
		{ID: "DO_NOT_PICK_IGNORE_THIS_ONE"},
	},
	{
		{ID: "DO_NOT_PICK_IGNORE_THIS_ONE_2"},
		{
			ID:              "ip_UVoPiUQoI1cqtRf5",
			Address:         "8.8.8.8",
			ReverseDNS:      "ip-8-8-8-8.test.katapult.cloud",
			VIP:             true,
			Label:           "testing",
			AddressWithMask: "8.8.8.8",
			Network: &core.Network{
				ID:         "test",
				Name:       "testing",
				Permalink:  "testing-123",
				DataCenter: fixtureDataCenters[1],
			},
		},
		{
			ID:              "ip_VVoPiUQoI1cqtRf5",
			Address:         "1.1.1.1",
			ReverseDNS:      "ip-1-1-1-1.test.katapult.cloud",
			VIP:             true,
			Label:           "testing2",
			AddressWithMask: "1.1.1.1",
			Network: &core.Network{
				ID:         "test",
				Name:       "testing",
				Permalink:  "testing-123",
				DataCenter: fixtureDataCenters[1],
			},
		},
		{
			ID:              "ip_VVoPiUQoI1cqtRf5",
			Address:         "1.1.1.2",
			ReverseDNS:      "ip-1-1-1-2.test.katapult.cloud",
			VIP:             true,
			Label:           "testing3",
			AddressWithMask: "1.1.1.2",
			Network: &core.Network{
				ID:         "test",
				Name:       "testing",
				Permalink:  "testing-123",
				DataCenter: fixtureDataCenters[1],
			},
		},
		{
			ID:              "ip_VVoPiUQoI1cqtRf5",
			Address:         "1.1.1.3",
			ReverseDNS:      "ip-1-1-1-3.test.katapult.cloud",
			VIP:             true,
			Label:           "testing4",
			AddressWithMask: "1.1.1.3",
			Network: &core.Network{
				ID:         "test",
				Name:       "testing",
				Permalink:  "testing-123",
				DataCenter: fixtureDataCenters[1],
			},
		},
	},
}

var successIPPages = map[string]ipPages{
	"testing": mockIPPages,
	"loge":    mockIPPages,
}

var mockSSHPages = sshPages{
	{
		{ID: "DO_NOT_PICK_IGNORE_THIS_ONE"},
	},
	{
		{ID: "DO_NOT_PICK_IGNORE_THIS_ONE_2"},
		{
			ID:          "key_PiUQoI1cqt43Dkf",
			Name:        "testing",
			Fingerprint: "22:57:25:0d:8a:ad:00:d0:91:a2:23:7d:7b:70:39:0c",
		},
		{
			ID:          "key_PiUQoI1cqt43Dkg",
			Name:        "testing1",
			Fingerprint: "23:57:25:0d:8a:ad:00:d0:91:a2:23:7d:7b:70:39:0c",
		},
		{
			ID:          "key_PiUQoI1cqt43Dke",
			Name:        "testing2",
			Fingerprint: "24:57:25:0d:8a:ad:00:d0:91:a2:23:7d:7b:70:39:0c",
		},
		{
			ID:          "key_PiUQoI1cqt43Dkd",
			Name:        "testing3",
			Fingerprint: "25:57:25:0d:8a:ad:00:d0:91:a2:23:7d:7b:70:39:0c",
		},
		{
			ID:          "key_PiUQoI1cqt43Dkc",
			Name:        "testing4",
			Fingerprint: "26:57:25:0d:8a:ad:00:d0:91:a2:23:7d:7b:70:39:0c",
		},
		{
			ID:          "key_PiUQoI1cqt43Dkb",
			Name:        "testing5",
			Fingerprint: "27:57:25:0d:8a:ad:00:d0:91:a2:23:7d:7b:70:39:0c",
		},
		{
			ID:          "key_PiUQoI1cqt43Dka",
			Name:        "testing6",
			Fingerprint: "28:57:25:0d:8a:ad:00:d0:91:a2:23:7d:7b:70:39:0c",
		},
	},
}

var successKeyPages = map[string]sshPages{
	"testing": mockSSHPages,
	"loge":    mockSSHPages,
}

var mockTagPages = tagPages{
	{
		{Name: "A", ID: "DO_NOT_PICK_IGNORE_THIS_ONE"},
	},
	{
		{Name: "B", ID: "DO_NOT_PICK_IGNORE_THIS_ONE_2"},
		{
			ID:        "tag_PiUQoI1cqt43gei",
			Name:      "Testing",
			Color:     "fffff",
			CreatedAt: timestamp.Unix(1, 0),
		},
		{
			ID:        "tag_PiUQoI1cqt43gea",
			Name:      "Testing 1",
			Color:     "fffff",
			CreatedAt: timestamp.Unix(1, 0),
		},
		{
			ID:        "tag_PiUQoI1cqt43geb",
			Name:      "Testing 2",
			Color:     "fffff",
			CreatedAt: timestamp.Unix(1, 0),
		},
		{
			ID:        "tag_PiUQoI1cqt43gec",
			Name:      "Testing 3",
			Color:     "fffff",
			CreatedAt: timestamp.Unix(1, 0),
		},
		{
			ID:        "tag_PiUQoI1cqt43ged",
			Name:      "Testing 4",
			Color:     "fffff",
			CreatedAt: timestamp.Unix(1, 0),
		},
	},
}

var successTagPages = map[string]tagPages{
	"testing": mockTagPages,
	"loge":    mockTagPages,
}

func TestVMs_Create(t *testing.T) {
	tests := []struct {
		name string

		envs map[string]string

		orgs       []*core.Organization
		orgsThrows string

		dcs       []*core.DataCenter
		dcsThrows string

		packages       []*core.VirtualMachinePackage
		packagesThrows string

		expectedRef core.OrganizationRef

		diskTemplates       []*core.DiskTemplate
		diskTemplatesThrows string

		ipIDPages map[string]ipPages
		ipThrows  string

		keysIDPages map[string]sshPages
		keysThrows  string

		tagIDPages map[string]tagPages
		tagThrows  string

		vmCreatorThrows string

		inputs  [][]byte
		stderr  string
		wantErr string
	}{
		// Successes
		{
			name:          "success with no env",
			orgs:          fixtureOrganizations,
			dcs:           fixtureDataCenters,
			packages:      successPackages,
			expectedRef:   core.OrganizationRef{ID: "testing"},
			diskTemplates: successDiskTemplates,
			ipIDPages:     successIPPages,
			keysIDPages:   successKeyPages,
			tagIDPages:    successTagPages,
			inputs: [][]byte{
				// Organization selection.
				keystrokes.DownArrow, keystrokes.Enter,

				// Data center selection.
				keystrokes.DownArrow, keystrokes.Enter,

				// Package selection.
				keystrokes.DownArrow, keystrokes.Enter,

				// Distro selection.
				keystrokes.DownArrow, keystrokes.Enter,

				// IP address selection.
				keystrokes.DownArrow, keystrokes.DownArrow, keystrokes.Enter, keystrokes.Escape,

				// Key selection.
				keystrokes.DownArrow, keystrokes.DownArrow, keystrokes.Enter, keystrokes.Escape,

				// Tag selection.
				keystrokes.DownArrow, keystrokes.DownArrow, keystrokes.Enter, keystrokes.Escape,

				// Name field.
				{'n'},
				{'a'},
				{'m'},
				{'e'},
				{'t'},
				{'e'},
				{'s'},
				{'t'},
				keystrokes.DownArrow,

				// Hostname field.
				{'h'},
				{'o'},
				{'s'},
				{'t'},
				{'t'},
				{'e'},
				{'s'},
				{'t'},
				keystrokes.DownArrow,

				// Description field.
				{'d'},
				{'e'},
				{'s'},
				{'c'},
				{'t'},
				{'e'},
				{'s'},
				{'t'},
				keystrokes.Enter,
			},
		},
		{
			name: "success with org name env",
			envs: map[string]string{
				"KATAPULT_ORG_NAME": "Loge Enthusiasts",
			},
			orgs:          fixtureOrganizations,
			dcs:           fixtureDataCenters,
			packages:      successPackages,
			expectedRef:   core.OrganizationRef{ID: "loge"},
			diskTemplates: successDiskTemplates,
			ipIDPages:     successIPPages,
			keysIDPages:   successKeyPages,
			tagIDPages:    successTagPages,
			inputs: [][]byte{
				// Data center selection.
				keystrokes.DownArrow, keystrokes.Enter,

				// Package selection.
				keystrokes.DownArrow, keystrokes.Enter,

				// Distro selection.
				keystrokes.DownArrow, keystrokes.Enter,

				// IP address selection.
				keystrokes.DownArrow, keystrokes.DownArrow, keystrokes.Enter, keystrokes.Escape,

				// Key selection.
				keystrokes.DownArrow, keystrokes.DownArrow, keystrokes.Enter, keystrokes.Escape,

				// Tag selection.
				keystrokes.DownArrow, keystrokes.DownArrow, keystrokes.Enter, keystrokes.Escape,

				// Name field.
				{'n'},
				{'a'},
				{'m'},
				{'e'},
				{'t'},
				{'e'},
				{'s'},
				{'t'},
				keystrokes.DownArrow,

				// Hostname field.
				{'h'},
				{'o'},
				{'s'},
				{'t'},
				{'t'},
				{'e'},
				{'s'},
				{'t'},
				keystrokes.DownArrow,

				// Description field.
				{'d'},
				{'e'},
				{'s'},
				{'c'},
				{'t'},
				{'e'},
				{'s'},
				{'t'},
				keystrokes.Enter,
			},
		},
		{
			name: "success with dc name env",
			envs: map[string]string{
				"KATAPULT_DC_NAME": "hello",
			},
			orgs:          fixtureOrganizations,
			dcs:           fixtureDataCenters,
			packages:      successPackages,
			expectedRef:   core.OrganizationRef{ID: "testing"},
			diskTemplates: successDiskTemplates,
			ipIDPages:     successIPPages,
			keysIDPages:   successKeyPages,
			tagIDPages:    successTagPages,
			inputs: [][]byte{
				// Organization selection.
				keystrokes.DownArrow, keystrokes.Enter,

				// Package selection.
				keystrokes.DownArrow, keystrokes.Enter,

				// Distro selection.
				keystrokes.DownArrow, keystrokes.Enter,

				// IP address selection.
				keystrokes.DownArrow, keystrokes.DownArrow, keystrokes.Enter, keystrokes.Escape,

				// Key selection.
				keystrokes.DownArrow, keystrokes.DownArrow, keystrokes.Enter, keystrokes.Escape,

				// Tag selection.
				keystrokes.DownArrow, keystrokes.DownArrow, keystrokes.Enter, keystrokes.Escape,

				// Name field.
				{'n'},
				{'a'},
				{'m'},
				{'e'},
				{'t'},
				{'e'},
				{'s'},
				{'t'},
				keystrokes.DownArrow,

				// Hostname field.
				{'h'},
				{'o'},
				{'s'},
				{'t'},
				{'t'},
				{'e'},
				{'s'},
				{'t'},
				keystrokes.DownArrow,

				// Description field.
				{'d'},
				{'e'},
				{'s'},
				{'c'},
				{'t'},
				{'e'},
				{'s'},
				{'t'},
				keystrokes.Enter,
			},
		},
		{
			name: "success with package name env",
			envs: map[string]string{
				"KATAPULT_PACKAGE_NAME": "Test",
			},
			orgs:          fixtureOrganizations,
			dcs:           fixtureDataCenters,
			packages:      successPackages,
			expectedRef:   core.OrganizationRef{ID: "testing"},
			diskTemplates: successDiskTemplates,
			ipIDPages:     successIPPages,
			keysIDPages:   successKeyPages,
			tagIDPages:    successTagPages,
			inputs: [][]byte{
				// Organization selection.
				keystrokes.DownArrow, keystrokes.Enter,

				// Data center selection.
				keystrokes.DownArrow, keystrokes.Enter,

				// Distro selection.
				keystrokes.DownArrow, keystrokes.Enter,

				// IP address selection.
				keystrokes.DownArrow, keystrokes.DownArrow, keystrokes.Enter, keystrokes.Escape,

				// Key selection.
				keystrokes.DownArrow, keystrokes.DownArrow, keystrokes.Enter, keystrokes.Escape,

				// Tag selection.
				keystrokes.DownArrow, keystrokes.DownArrow, keystrokes.Enter, keystrokes.Escape,

				// Name field.
				{'n'},
				{'a'},
				{'m'},
				{'e'},
				{'t'},
				{'e'},
				{'s'},
				{'t'},
				keystrokes.DownArrow,

				// Hostname field.
				{'h'},
				{'o'},
				{'s'},
				{'t'},
				{'t'},
				{'e'},
				{'s'},
				{'t'},
				keystrokes.DownArrow,

				// Description field.
				{'d'},
				{'e'},
				{'s'},
				{'c'},
				{'t'},
				{'e'},
				{'s'},
				{'t'},
				keystrokes.Enter,
			},
		},
		{
			name: "success with distribution name env",
			envs: map[string]string{
				"KATAPULT_DISTRIBUTION_NAME": "Ubuntu 20.04",
			},
			orgs:          fixtureOrganizations,
			dcs:           fixtureDataCenters,
			packages:      successPackages,
			expectedRef:   core.OrganizationRef{ID: "testing"},
			diskTemplates: successDiskTemplates,
			ipIDPages:     successIPPages,
			keysIDPages:   successKeyPages,
			tagIDPages:    successTagPages,
			inputs: [][]byte{
				// Organization selection.
				keystrokes.DownArrow, keystrokes.Enter,

				// Data center selection.
				keystrokes.DownArrow, keystrokes.Enter,

				// Package selection.
				keystrokes.DownArrow, keystrokes.Enter,

				// IP address selection.
				keystrokes.DownArrow, keystrokes.DownArrow, keystrokes.Enter, keystrokes.Escape,

				// Key selection.
				keystrokes.DownArrow, keystrokes.DownArrow, keystrokes.Enter, keystrokes.Escape,

				// Tag selection.
				keystrokes.DownArrow, keystrokes.DownArrow, keystrokes.Enter, keystrokes.Escape,

				// Name field.
				{'n'},
				{'a'},
				{'m'},
				{'e'},
				{'t'},
				{'e'},
				{'s'},
				{'t'},
				keystrokes.DownArrow,

				// Hostname field.
				{'h'},
				{'o'},
				{'s'},
				{'t'},
				{'t'},
				{'e'},
				{'s'},
				{'t'},
				keystrokes.DownArrow,

				// Description field.
				{'d'},
				{'e'},
				{'s'},
				{'c'},
				{'t'},
				{'e'},
				{'s'},
				{'t'},
				keystrokes.Enter,
			},
		},
		{
			name: "success from full env",
			envs: map[string]string{
				"KATAPULT_ORG_SUBDOMAIN":   "loge",
				"KATAPULT_DC_ID":           "dc_9UVoPiUQoI1cqtRd",
				"KATAPULT_PACKAGE_ID":      "vmpkg_9UVoPiUQoI1cqtRd",
				"KATAPULT_DISTRIBUTION_ID": "Ubuntu-20-04",
				"KATAPULT_IP_ADDRESSES":    "1.1.1.1,1.1.1.2,1.1.1.3",
				"KATAPULT_SSH_KEY_IDS":     "key_PiUQoI1cqt43Dkc,key_PiUQoI1cqt43Dkd",
				"KATAPULT_SSH_KEY_NAMES":   "testing,testing1",
				"KATAPULT_SSH_KEY_FINGERPRINTS": "28:57:25:0d:8a:ad:00:d0:91:a2:23:7d:7b:70:39:0c," +
					"27:57:25:0d:8a:ad:00:d0:91:a2:23:7d:7b:70:39:0c",
				"KATAPULT_TAG_NAMES":   "Testing 2,Testing 3",
				"KATAPULT_TAG_IDS":     "tag_PiUQoI1cqt43gea,tag_PiUQoI1cqt43geb",
				"KATAPULT_NAME":        "test",
				"KATAPULT_HOSTNAME":    "testing",
				"KATAPULT_DESCRIPTION": "123",
			},
			orgs:          fixtureOrganizations,
			dcs:           fixtureDataCenters,
			packages:      successPackages,
			expectedRef:   core.OrganizationRef{ID: "loge"},
			diskTemplates: successDiskTemplates,
			ipIDPages:     successIPPages,
			keysIDPages:   successKeyPages,
			tagIDPages:    successTagPages,
		},
		{
			name: "success from full minus hostname env",
			envs: map[string]string{
				"KATAPULT_ORG_SUBDOMAIN":   "loge",
				"KATAPULT_DC_ID":           "dc_9UVoPiUQoI1cqtRd",
				"KATAPULT_PACKAGE_ID":      "vmpkg_9UVoPiUQoI1cqtRd",
				"KATAPULT_DISTRIBUTION_ID": "Ubuntu-20-04",
				"KATAPULT_IP_ADDRESSES":    "1.1.1.1,1.1.1.2,1.1.1.3",
				"KATAPULT_SSH_KEY_IDS":     "key_PiUQoI1cqt43Dkc,key_PiUQoI1cqt43Dkd",
				"KATAPULT_SSH_KEY_NAMES":   "testing,testing1",
				"KATAPULT_SSH_KEY_FINGERPRINTS": "28:57:25:0d:8a:ad:00:d0:91:a2:23:7d:7b:70:39:0c," +
					"27:57:25:0d:8a:ad:00:d0:91:a2:23:7d:7b:70:39:0c",
				"KATAPULT_TAG_NAMES":   "Testing 2,Testing 3",
				"KATAPULT_TAG_IDS":     "tag_PiUQoI1cqt43gea,tag_PiUQoI1cqt43geb",
				"KATAPULT_NAME":        "test",
				"KATAPULT_DESCRIPTION": "123",
			},
			orgs:          fixtureOrganizations,
			dcs:           fixtureDataCenters,
			packages:      successPackages,
			expectedRef:   core.OrganizationRef{ID: "loge"},
			diskTemplates: successDiskTemplates,
			ipIDPages:     successIPPages,
			keysIDPages:   successKeyPages,
			tagIDPages:    successTagPages,
			inputs: [][]byte{
				{'t'},
				{'e'},
				{'s'},
				{'t'},
				{'i'},
				{'n'},
				{'g'},
				keystrokes.Enter,
			},
		},

		// Client error throwing

		{
			name:       "orgs throws error",
			orgsThrows: "power cut at the organization",
			wantErr:    "power cut at the organization",
		},
		{
			name:      "dcs throws error",
			orgs:      []*core.Organization{{Name: "test", SubDomain: "testing"}},
			dcsThrows: "power cut at the organization",
			wantErr:   "power cut at the organization",
			inputs:    [][]byte{keystrokes.Enter},
		},
		{
			name:           "packages throws error",
			orgs:           []*core.Organization{{Name: "test", SubDomain: "testing"}},
			dcs:            []*core.DataCenter{{Country: &core.Country{}}},
			packagesThrows: "power cut at the organization",
			wantErr:        "power cut at the organization",
			inputs:         [][]byte{keystrokes.Enter, keystrokes.Enter},
		},
		{
			name:                "disk templates throws error",
			orgs:                []*core.Organization{{Name: "test", SubDomain: "testing"}},
			dcs:                 []*core.DataCenter{{Country: &core.Country{}}},
			packages:            []*core.VirtualMachinePackage{{}},
			diskTemplatesThrows: "power cut at the organization",
			wantErr:             "power cut at the organization",
			inputs:              [][]byte{keystrokes.Enter, keystrokes.Enter, keystrokes.Enter},
		},
		{
			name:          "ip listing throws error",
			orgs:          []*core.Organization{{Name: "test", SubDomain: "testing"}},
			dcs:           []*core.DataCenter{{Country: &core.Country{}}},
			packages:      []*core.VirtualMachinePackage{{}},
			diskTemplates: []*core.DiskTemplate{{}},
			ipThrows:      "power cut at the organization",
			wantErr:       "power cut at the organization",
			inputs:        [][]byte{keystrokes.Enter, keystrokes.Enter, keystrokes.Enter, keystrokes.Enter},
		},
		{
			name:          "key throws error",
			orgs:          []*core.Organization{{Name: "test", SubDomain: "testing", ID: "org"}},
			expectedRef:   core.OrganizationRef{ID: "org"},
			dcs:           []*core.DataCenter{{Country: &core.Country{}}},
			packages:      []*core.VirtualMachinePackage{{}},
			diskTemplates: []*core.DiskTemplate{{}},
			ipIDPages:     map[string]ipPages{"org": {{}}},
			keysThrows:    "power cut at the organization",
			wantErr:       "power cut at the organization",
			inputs:        [][]byte{keystrokes.Enter, keystrokes.Enter, keystrokes.Enter, keystrokes.Enter},
		},
		{
			name:          "tag listing throws error",
			orgs:          []*core.Organization{{Name: "test", SubDomain: "testing", ID: "org"}},
			expectedRef:   core.OrganizationRef{ID: "org"},
			dcs:           []*core.DataCenter{{Country: &core.Country{}}},
			packages:      []*core.VirtualMachinePackage{{}},
			diskTemplates: []*core.DiskTemplate{{}},
			keysIDPages:   map[string]sshPages{"org": {{}}},
			ipIDPages:     map[string]ipPages{"org": {{}}},
			tagThrows:     "power cut at the organization",
			wantErr:       "power cut at the organization",
			inputs:        [][]byte{keystrokes.Enter, keystrokes.Enter, keystrokes.Enter, keystrokes.Enter},
		},
		{
			name:            "vm creator throws error",
			orgs:            []*core.Organization{{Name: "test", SubDomain: "testing", ID: "org"}},
			expectedRef:     core.OrganizationRef{ID: "org"},
			dcs:             []*core.DataCenter{{Country: &core.Country{}}},
			packages:        []*core.VirtualMachinePackage{{}},
			diskTemplates:   []*core.DiskTemplate{{}},
			keysIDPages:     map[string]sshPages{"org": {{}}},
			ipIDPages:       map[string]ipPages{"org": {{}}},
			tagIDPages:      map[string]tagPages{"org": {{}}},
			vmCreatorThrows: "power cut at the organization",
			wantErr:         "power cut at the organization",
			inputs: [][]byte{
				keystrokes.Enter, keystrokes.Enter, keystrokes.Enter,
				keystrokes.Enter,
				{'a'},
				keystrokes.Enter,
			},
		},

		// Env validation handling

		{
			name:    "org name error",
			orgs:    []*core.Organization{},
			envs:    map[string]string{"KATAPULT_ORG_NAME": "test"},
			wantErr: "the org name/subdomain in your org env variable not attached to your user",
		},
		{
			name:    "org subdomain error",
			orgs:    []*core.Organization{},
			envs:    map[string]string{"KATAPULT_ORG_SUBDOMAIN": "test"},
			wantErr: "the org name/subdomain in your org env variable not attached to your user",
		},
		{
			name: "dc id error",
			orgs: []*core.Organization{{Name: "test", SubDomain: "testing", ID: "org"}},
			dcs:  []*core.DataCenter{},
			envs: map[string]string{
				"KATAPULT_ORG_NAME": "test",
				"KATAPULT_DC_ID":    "dc",
			},
			wantErr: "the dc name/id in your dc env variable not attached to your user",
		},
		{
			name: "dc name error",
			orgs: []*core.Organization{{Name: "test", SubDomain: "testing", ID: "org"}},
			dcs:  []*core.DataCenter{},
			envs: map[string]string{
				"KATAPULT_ORG_NAME": "test",
				"KATAPULT_DC_NAME":  "dc",
			},
			wantErr: "the dc name/id in your dc env variable not attached to your user",
		},
		{
			name:     "package id error",
			orgs:     []*core.Organization{{Name: "test", SubDomain: "testing", ID: "org"}},
			dcs:      []*core.DataCenter{{Name: "dc", ID: "dc"}},
			packages: []*core.VirtualMachinePackage{},
			envs: map[string]string{
				"KATAPULT_ORG_NAME":   "test",
				"KATAPULT_DC_NAME":    "dc",
				"KATAPULT_PACKAGE_ID": "package",
			},
			wantErr: "the package name/slug in your package env variable not attached to your user",
		},
		{
			name:     "package name error",
			orgs:     []*core.Organization{{Name: "test", SubDomain: "testing", ID: "org"}},
			dcs:      []*core.DataCenter{{Name: "dc", ID: "dc"}},
			packages: []*core.VirtualMachinePackage{},
			envs: map[string]string{
				"KATAPULT_ORG_NAME":     "test",
				"KATAPULT_DC_NAME":      "dc",
				"KATAPULT_PACKAGE_NAME": "package",
			},
			wantErr: "the package name/slug in your package env variable not attached to your user",
		},
		{
			name:          "distro name error",
			orgs:          []*core.Organization{{Name: "test", SubDomain: "testing", ID: "org"}},
			dcs:           []*core.DataCenter{{Name: "dc", ID: "dc"}},
			packages:      []*core.VirtualMachinePackage{{ID: "package"}},
			diskTemplates: []*core.DiskTemplate{},
			expectedRef:   core.OrganizationRef{ID: "org"},
			envs: map[string]string{
				"KATAPULT_ORG_NAME":        "test",
				"KATAPULT_DC_NAME":         "dc",
				"KATAPULT_PACKAGE_ID":      "package",
				"KATAPULT_DISTRIBUTION_ID": "testing",
			},
			wantErr: "the distribution name/slug in your distribution env variables not attached to your user",
		},
		{
			name:          "distro id error",
			orgs:          []*core.Organization{{Name: "test", SubDomain: "testing", ID: "org"}},
			dcs:           []*core.DataCenter{{Name: "dc", ID: "dc"}},
			packages:      []*core.VirtualMachinePackage{{ID: "package"}},
			diskTemplates: []*core.DiskTemplate{},
			expectedRef:   core.OrganizationRef{ID: "org"},
			envs: map[string]string{
				"KATAPULT_ORG_NAME":          "test",
				"KATAPULT_DC_NAME":           "dc",
				"KATAPULT_PACKAGE_ID":        "package",
				"KATAPULT_DISTRIBUTION_NAME": "testing",
			},
			wantErr: "the distribution name/slug in your distribution env variables not attached to your user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Defines stdin.
			stdin := &console.StdinDripFeeder{T: t, Inputs: tt.inputs}

			// Defines the mock terminal.
			mockTerminal := &console.MockTerminal{}

			// Create the clients.
			orgsClient := mockOrganizationsListClient{
				orgs:   tt.orgs,
				throws: tt.orgsThrows,
			}
			dcsClient := mockDataCentersClient{
				dcs:    tt.dcs,
				throws: tt.dcsThrows,
			}
			vmPackagesClient := mockVMPackagesClient{
				packages: tt.packages,
				throws:   tt.packagesThrows,
			}
			diskTemplatesClient := mockDiskTemplatesClient{
				diskTemplates: tt.diskTemplates,
				throws:        tt.diskTemplatesThrows,
				ref:           tt.expectedRef,
			}
			ipAddressesClient := mockIPAddressClient{
				throws:              tt.ipThrows,
				organizationIDPages: tt.ipIDPages,
			}
			sshKeysClient := mockSSHKeysClient{
				throws:              tt.keysThrows,
				organizationIDPages: tt.keysIDPages,
			}
			tags := mockTagsClient{
				throws:              tt.tagThrows,
				organizationIDPages: tt.tagIDPages,
			}

			// The VM builder client is special since it logs the result so we can process it.
			vmBuilderClient := &mockVMBuilderClient{throws: tt.vmCreatorThrows}

			// Create the command.
			cmd := virtualMachinesCmd(
				nil, orgsClient, dcsClient, vmPackagesClient, diskTemplatesClient,
				ipAddressesClient, sshKeysClient, tags, vmBuilderClient, mockTerminal,
				mapGetter{m: tt.envs})
			cmd.SetIn(stdin)
			cmd.SetArgs([]string{"create"})
			stdout := assertCobraCommandReturnStdout(t, cmd, tt.wantErr, tt.stderr)

			// Create the resulting golden data and handle it.
			buf := &bytes.Buffer{}
			buf.WriteString("-- STDOUT --\n\n")
			_, _ = mockTerminal.Buffer.WriteTo(buf)
			buf.WriteString(stdout)
			buf.WriteString("\n\n-- BUILD SPEC --\n\n")
			enc := json.NewEncoder(buf)
			enc.SetIndent("", "  ")
			if err := enc.Encode(vmBuilderClient); err != nil {
				assert.NoError(t, err)
			}
			if golden.Update() {
				golden.Set(t, buf.Bytes())
			}
			assert.Equal(t, string(golden.Get(t)), buf.String())
		})
	}
}
