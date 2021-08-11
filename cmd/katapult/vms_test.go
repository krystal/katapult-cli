package main

import (
	"context"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
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
			cmd := virtualMachinesCmd(&vmsClient{organizationIDPages: tt.id, organizationSubdomainPages: tt.subdomains})
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
			cmd := virtualMachinesCmd(client)
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
			cmd := virtualMachinesCmd(client)
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
			cmd := virtualMachinesCmd(client)
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
			cmd := virtualMachinesCmd(client)
			cmd.SetArgs(tt.args)
			assertCobraCommand(t, cmd, tt.wantErr, tt.stderr)
			if tt.validate != nil {
				assert.Equal(t, "", tt.validate(client))
			}
		})
	}
}
