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

type mockOrganisationsListClient struct {
	orgs   []*core.Organization
	throws string
}

func (m mockOrganisationsListClient) List(context.Context) ([]*core.Organization, *katapult.Response, error) {
	if m.throws != "" {
		return nil, nil, errors.New(m.throws)
	}
	return m.orgs, nil, nil
}

func TestOrganizations_List(t *testing.T) {
	tests := []struct {
		name string

		orgs    []*core.Organization
		want    string
		stderr  string
		throws  string
		wantErr string
	}{
		{
			name: "organizations list",
			orgs: fixtureOrganizations,
			want: ` - Loge Enthusiasts (loge) [loge]
 - testing, testing, 123 (test) [testing]
`,
		},
		{
			name: "empty organizations",
			orgs: []*core.Organization{},
		},
		{
			name:    "organization error",
			throws:  "test error",
			wantErr: "test error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := organizationsCmd(mockOrganisationsListClient{orgs: tt.orgs, throws: tt.throws})
			cmd.SetArgs([]string{"list"})
			assertCobraCommand(t, cmd, tt.wantErr, tt.want, tt.stderr)
		})
	}
}
