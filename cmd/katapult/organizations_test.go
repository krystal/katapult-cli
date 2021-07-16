package main

import (
	"bytes"
	"context"
	"testing"

	"github.com/krystal/go-katapult"

	"github.com/krystal/go-katapult/core"
	"github.com/stretchr/testify/assert"
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

const stdoutOrganizationsList = ` - Loge Enthusiasts (loge) [loge]
 - testing, testing, 123 (test) [testing]
`

type mockOrganisationsListClient struct{}

func (mockOrganisationsListClient) List(context.Context) ([]*core.Organization, *katapult.Response, error) {
	return organizations, nil, nil
}

func TestOrganizations_List(t *testing.T) {
	cmd := organizationsCmd(mockOrganisationsListClient{})
	stdout := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetArgs([]string{"list"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, stdoutOrganizationsList, stdout.String())
}
