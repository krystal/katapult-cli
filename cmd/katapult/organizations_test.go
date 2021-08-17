package main

import (
	"testing"

	"github.com/krystal/go-katapult/core"
)

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
			cmd := organizationsCmd(mockOrganizationsListClient{orgs: tt.orgs, throws: tt.throws})
			cmd.SetArgs([]string{"list"})
			assertCobraCommand(t, cmd, tt.wantErr, tt.want, tt.stderr)
		})
	}
}
