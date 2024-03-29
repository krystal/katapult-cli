package main

import (
	"bytes"
	"testing"

	"github.com/krystal/go-katapult/core"
	"github.com/krystal/katapult-cli/internal/golden"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataCenters_List(t *testing.T) {
	tests := []struct {
		name string

		output  string
		dcs     []*core.DataCenter
		stderr  string
		throws  string
		wantErr string
	}{
		{
			name: "data center human readable list",
			dcs:  fixtureDataCenters,
		},
		{
			name:   "data center json list",
			output: "json",
			dcs:    fixtureDataCenters,
		},
		{
			name: "empty data centers human readable",
			dcs:  []*core.DataCenter{},
		},
		{
			name:   "empty data centers json",
			output: "json",
			dcs:    []*core.DataCenter{},
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
			outputFlag = tt.output
			err := cmd.Execute()
			outputFlag = ""

			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)

			if golden.Update() {
				golden.Set(t, stdout.Bytes())
			}
			assert.Equal(t, string(golden.Get(t)), stdout.String())
			assert.Equal(t, tt.stderr, stderr.String())
		})
	}
}

func TestDataCenters_Get(t *testing.T) {
	tests := []struct {
		name string

		args    []string
		output  string
		dc      string
		stderr  string
		wantErr string
	}{
		{
			name: "display POG1 human readable",
			args: []string{"get", "POG1"},
		},
		{
			name:   "display POG1 json",
			args:   []string{"get", "POG1"},
			output: "json",
		},
		{
			name:   "display GB1 json",
			args:   []string{"get", "GB1"},
			output: "json",
		},
		{
			name: "display GB1 human readable",
			args: []string{"get", "GB1"},
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
			outputFlag = tt.output
			assertCobraCommand(t, cmd, tt.wantErr, tt.stderr)
			outputFlag = ""
		})
	}
}
