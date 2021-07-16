package main

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"

	"github.com/stretchr/testify/assert"
)

var dcs = []*core.DataCenter{
	{
		ID:        "POG1-ID",
		Name:      "hello",
		Permalink: "POG1",
		Country: &core.Country{
			ID:   "POG",
			Name: "Pogland",
		},
	},
	{
		ID:        "GB1-ID",
		Name:      "hello",
		Permalink: "GB1",
		Country: &core.Country{
			ID:   "UK",
			Name: "United Kingdom",
		},
	},
}

type mockDataCentersClient struct{}

func (mockDataCentersClient) List(context.Context) ([]*core.DataCenter, *katapult.Response, error) {
	return dcs, nil, nil
}

func (mockDataCentersClient) Get(_ context.Context, ref core.DataCenterRef) (*core.DataCenter, *katapult.Response, error) {
	for _, v := range dcs {
		if v.ID == ref.Permalink {
			return v, nil, nil
		}
	}
	return nil, nil, fmt.Errorf("unknown datacentre")
}

const stdoutDcList = ` - hello (POG1) [POG1-ID] / Pogland
 - hello (GB1) [GB1-ID] / United Kingdom
`

func TestDataCenters_List(t *testing.T) {
	cmd := dataCentersCmd(mockDataCentersClient{})
	stdout := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetArgs([]string{"list"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, stdoutDcList, stdout.String())
}

func TestDataCenters_Get(t *testing.T) {
	tests := []struct {
		name string

		args   []string
		dc     string
		wants  string
		stderr string
		err    string
	}{
		{
			name: "display POG1",
			args: []string{"get", "POG1-ID"},
			wants: "hello (POG1) [POG1-ID] / Pogland\n",
		},
		{
			name: "display GB1",
			args: []string{"get", "GB1-ID"},
			wants: "hello (GB1) [GB1-ID] / United Kingdom\n",
		},
		{
			name:   "display invalid DC",
			args:   []string{"get", "UNPOG1"},
			stderr: "Error: unknown datacentre\n",
			err:    "unknown datacentre",
		},
	}

	cmd := dataCentersCmd(mockDataCentersClient{})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			cmd.SetOut(stdout)
			cmd.SetErr(stderr)
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			assert.Equal(t, tt.stderr, stderr.String())
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
