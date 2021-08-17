package main

import (
	"bytes"
	"testing"

	"github.com/krystal/go-katapult/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataCenters_List(t *testing.T) {
	tests := []struct {
		name string

		dcs     []*core.DataCenter
		want    string
		stderr  string
		throws  string
		wantErr string
	}{
		{
			name: "data center list",
			dcs:  fixtureDataCenters,
			want: ` - hello (POG1) [dc_9UVoPiUQoI1cqtRd] / Pogland
 - hello (GB1) [dc_9UVoPiUQoI1cqtR0] / United Kingdom
`,
		},
		{
			name: "empty data centers",
			dcs:  []*core.DataCenter{},
		},
		{
			name:    "data center error",
			throws:  "test error",
			wantErr: "test error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := dataCentersCmd(mockDataCentersClient{dcs: tt.dcs, throws: tt.throws})
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			cmd.SetOut(stdout)
			cmd.SetErr(stderr)
			cmd.SetArgs([]string{"list"})
			err := cmd.Execute()

			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tt.want, stdout.String())
			assert.Equal(t, tt.stderr, stderr.String())
		})
	}
}

func TestDataCenters_Get(t *testing.T) {
	tests := []struct {
		name string

		args    []string
		dc      string
		want    string
		stderr  string
		wantErr string
	}{
		{
			name: "display POG1",
			args: []string{"get", "POG1"},
			want: "hello (POG1) [dc_9UVoPiUQoI1cqtRd] / Pogland\n",
		},
		{
			name: "display GB1",
			args: []string{"get", "GB1"},
			want: "hello (GB1) [dc_9UVoPiUQoI1cqtR0] / United Kingdom\n",
		},
		{
			name:    "display invalid DC",
			args:    []string{"get", "UNPOG1"},
			stderr:  "Error: unknown datacentre\n",
			wantErr: "unknown datacentre",
		},
	}

	cmd := dataCentersCmd(mockDataCentersClient{dcs: fixtureDataCenters})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd.SetArgs(tt.args)
			assertCobraCommand(t, cmd, tt.wantErr, tt.want, tt.stderr)
		})
	}
}
