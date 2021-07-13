package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const stdoutOrganizationsList = ` - Loge Enthusiasts (loge) [loge]
 - testing, testing, 123 (test) [testing]
`

func TestOrganizations_List(t *testing.T) {
	cmd := organizationsCmd(mockAPIClient{})
	stdout := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetArgs([]string{"list"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, stdoutOrganizationsList, stdout.String())
}
