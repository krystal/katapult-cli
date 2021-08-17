package main

import (
	"context"
	"errors"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
)

var fixtureOrganizations = []*core.Organization{
	{
		ID:        "loge",
		Name:      "Loge Enthusiasts",
		SubDomain: "loge",
	},
	{
		ID:        "testing",
		Name:      "testing, testing, 123",
		SubDomain: "test",
	},
}

type mockOrganizationsListClient struct {
	orgs   []*core.Organization
	throws string
}

func (m mockOrganizationsListClient) List(context.Context) ([]*core.Organization, *katapult.Response, error) {
	if m.throws != "" {
		return nil, nil, errors.New(m.throws)
	}
	return m.orgs, nil, nil
}

func TestOrganizations_List(t *testing.T) {
	tests := []struct {
		name string

		output  string
		orgs    []*core.Organization
		stderr  string
		throws  string
		wantErr string
	}{
		{
			name: "organizations list human readable",
			orgs: fixtureOrganizations,
		},
		{
			name:   "organizations list json",
			orgs:   fixtureOrganizations,
			output: "json",
		},
		{
			name: "empty organizations human readable",
			orgs: []*core.Organization{},
		},
		{
			name:   "empty organizations json",
			orgs:   []*core.Organization{},
			output: "json",
		},
		{
			name:    "organization error",
			throws:  "test error",
			wantErr: "test error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := organizationsCmd(mockOrganizationsListClient{orgs: tt.orgs, throws: tt.throws})
			outputFlag = tt.output
			cmd.SetArgs([]string{"list"})
			assertCobraCommand(t, cmd, tt.wantErr, tt.stderr)
			outputFlag = ""
		})
	}
}
