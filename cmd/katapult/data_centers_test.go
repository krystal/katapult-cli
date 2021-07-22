package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fixtureDataCenters = []*core.DataCenter{
	{
		ID:        "dc_9UVoPiUQoI1cqtRd",
		Name:      "hello",
		Permalink: "POG1",
		Country: &core.Country{
			ID:   "POG",
			Name: "Pogland",
		},
	},
	{
		ID:        "dc_9UVoPiUQoI1cqtR0",
		Name:      "hello",
		Permalink: "GB1",
		Country: &core.Country{
			ID:   "UK",
			Name: "United Kingdom",
		},
	},
}

type mockDataCentersClient struct {
	dcs    []*core.DataCenter
	throws string
}

func (m mockDataCentersClient) List(context.Context) ([]*core.DataCenter, *katapult.Response, error) {
	if m.throws != "" {
		return nil, nil, errors.New(m.throws)
	}
	return m.dcs, nil, nil
}

func (m mockDataCentersClient) Get(
	_ context.Context, ref core.DataCenterRef) (*core.DataCenter, *katapult.Response, error) {
	if m.throws != "" {
		return nil, nil, errors.New(m.throws)
	}
	for _, v := range m.dcs {
		if v.Permalink == ref.Permalink {
			return v, nil, nil
		}
	}
	return nil, nil, fmt.Errorf("unknown datacentre")
}

func TestDataCenters_List(t *testing.T) {
	tests := []struct {
		name string

		args    []string
		dcs     []*core.DataCenter
		want    string
		stderr  string
		throws  string
		wantErr string
	}{
		{
			name: "data center human readable list",
			dcs:  fixtureDataCenters,
			want: ` - hello (POG1) [dc_9UVoPiUQoI1cqtRd] / Pogland
 - hello (GB1) [dc_9UVoPiUQoI1cqtR0] / United Kingdom
`,
		},
		{
			name: "data center json list",
			args: []string{"list", "-o", "json"},
			dcs:  fixtureDataCenters,
			want: `[{"id":"dc_9UVoPiUQoI1cqtRd","name":"hello","permalink":"POG1","country":{"id":"POG","name":"Pogland"}},{"id":"dc_9UVoPiUQoI1cqtR0","name":"hello","permalink":"GB1","country":{"id":"UK","name":"United Kingdom"}}]
`,
		},
		{
			name: "empty data centers human readable",
			dcs:  []*core.DataCenter{},
		},
		{
			name: "empty data centers json",
			args: []string{"list", "-o", "json"},
			dcs:  []*core.DataCenter{},
			want: "[]\n",
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
			if tt.args == nil {
				cmd.SetArgs([]string{"list"})
			} else {
				cmd.SetArgs(tt.args)
			}
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
			name: "display POG1 human readable",
			args: []string{"get", "POG1"},
			want: "hello (POG1) [dc_9UVoPiUQoI1cqtRd] / Pogland\n",
		},
		{
			name: "display POG1 json",
			args: []string{"get", "POG1", "-o", "json"},
			want: `{"id":"dc_9UVoPiUQoI1cqtRd","name":"hello","permalink":"POG1","country":{"id":"POG","name":"Pogland"}}
`,
		},
		{
			name: "display GB1 json",
			args: []string{"get", "GB1", "-o", "json"},
			want: `{"id":"dc_9UVoPiUQoI1cqtR0","name":"hello","permalink":"GB1","country":{"id":"UK","name":"United Kingdom"}}
`,
		},
		{
			name: "display GB1 human readable",
			args: []string{"get", "GB1", "-o", "yaml"},
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
