package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const stdoutDcList = ` - hello (Hello World!) [POG1] / Pogland
 - hello (Hello World!) [GB1] / United Kingdom
`

func TestDataCenters_List(t *testing.T) {
	cmd := dataCentersCmd(mockAPIClient{})
	stdout := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetArgs([]string{"list"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, stdoutDcList, stdout.String())
}

func TestDataCenters_Get(t *testing.T) {
	cmd := dataCentersCmd(mockAPIClient{})
	tests := []struct {
		name string

		args  []string
		wants string
		err   string
	}{
		{
			name: "display POG1",
			args: []string{"get", "POG1"},
			wants: `hello (Hello World!) [POG1] / Pogland
`,
		},
		{
			name: "display GB1",
			args: []string{"get", "GB1"},
			wants: `hello (Hello World!) [GB1] / United Kingdom
`,
		},
		{
			name: "display invalid DC",
			args: []string{"get", "UNPOG1"},
			err:  "unknown datacentre",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			cmd.SetOut(stdout)
			cmd.SetErr(&bytes.Buffer{}) // Ignore stderr, this is just testing the command framework, not the error.
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
