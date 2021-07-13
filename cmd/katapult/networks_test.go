package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedNetworkResults = []struct {
	name string

	args  []string
	wants string
	err   string
}{
	{
		name: "Test listing pog-id",
		args: []string{"ls", "--id", "pog-id"},
		wants: `Networks:
 - Pognet 1 [pognet]
 - Pognet 2 [pognet2]
`,
	},
	{
		name: "Test listing pog-subdomain",
		args: []string{"ls", "--subdomain", "pog-subdomain"},
		wants: `Networks:
 - Pognet 3 [pognet3]
 - Pognet 4 [pognet4]
`,
	},
}

func TestNetworkList(t *testing.T) {
	for _, v := range expectedNetworkResults {
		cmd := networksCmd(mockAPIClient{})
		stdout := &bytes.Buffer{}
		cmd.SetOut(stdout)
		cmd.SetArgs(v.args)

		err := cmd.Execute()
		switch {
		case err == nil:
			assert.Equal(t, v.wants, stdout.String())
		case v.err != "":
			assert.Equal(t, v.err, err.Error())
		default:
			t.Fatal(err)
		}
	}
}
