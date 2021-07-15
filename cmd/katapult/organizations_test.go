package main

import (
	"bytes"
	"testing"

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

func TestOrganizations_List(t *testing.T) {
	mock := singleResponse(t, "/core/v1/organizations", "organizations", organizations)
	cmd := organizationsCmd(mock)
	stdout := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetArgs([]string{"list"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, stdoutOrganizationsList, stdout.String())
}
