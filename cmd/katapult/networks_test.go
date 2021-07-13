package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNetworkList(t *testing.T) {
	tests := []struct {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := networksCmd(mockAPIClient{})
			stdout := &bytes.Buffer{}
			cmd.SetOut(stdout)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			switch {
			case err == nil:
				assert.Equal(t, tt.wants, stdout.String())
			case tt.err != "":
				assert.Equal(t, tt.err, err.Error())
			default:
				t.Fatal(err)
			}
		})
	}
}
