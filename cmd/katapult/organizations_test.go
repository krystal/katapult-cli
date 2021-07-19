package main

import (
	"bytes"
	"context"
	"errors"
	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

var organizations = []*core.Organization{
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

		orgs   []*core.Organization
		wants  string
		stderr string
		throws string
		err    string
	}{
		{
			name: "orgamisations list",
			orgs: organizations,
			wants: ` - Loge Enthusiasts (loge) [loge]
 - testing, testing, 123 (test) [testing]
`,
		},
		{
			name: "blank organisations",
			orgs: []*core.Organization{},
		},
		{
			name:   "organization error",
			throws: "test error",
			err:    "test error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := organizationsCmd(mockOrganisationsListClient{orgs: tt.orgs, throws: tt.throws})
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			cmd.SetOut(stdout)
			cmd.SetErr(stderr)
			cmd.SetArgs([]string{"list"})
			err := cmd.Execute()
			switch {
			case err == nil:
				// Ignore this.
			case tt.err != "":
				assert.Equal(t, tt.err, err.Error())
				return
			default:
				t.Fatal(err)
			}
			assert.Equal(t, tt.wants, stdout.String())
			assert.Equal(t, tt.stderr, stderr.String())
		})
	}
}
