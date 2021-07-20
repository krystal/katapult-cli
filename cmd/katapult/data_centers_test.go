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
)

var fixtureDataCenters = []*core.DataCenter{
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
		if v.ID == ref.Permalink {
			return v, nil, nil
		}
	}
	return nil, nil, fmt.Errorf("unknown datacentre")
}

func TestDataCenters_List(t *testing.T) {
	tests := []struct {
		name string

		dcs    []*core.DataCenter
		want   string
		stderr string
		throws string
		err    string
	}{
		{
			name: "data center list",
			dcs:  fixtureDataCenters,
			want: ` - hello (POG1) [POG1-ID] / Pogland
 - hello (GB1) [GB1-ID] / United Kingdom
`,
		},
		{
			name: "empty data centers",
			dcs:  []*core.DataCenter{},
		},
		{
			name:   "data center error",
			throws: "test error",
			err:    "test error",
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
			switch {
			case err == nil:
				// Ignore this.
			case tt.err != "":
				assert.Equal(t, tt.err, err.Error())
				return
			default:
				t.Fatal(err)
			}
			assert.Equal(t, tt.want, stdout.String())
			assert.Equal(t, tt.stderr, stderr.String())
		})
	}
}

func TestDataCenters_Get(t *testing.T) {
	tests := []struct {
		name string

		args   []string
		dc     string
		want   string
		stderr string
		err    string
	}{
		{
			name: "display POG1",
			args: []string{"get", "POG1-ID"},
			want: "hello (POG1) [POG1-ID] / Pogland\n",
		},
		{
			name: "display GB1",
			args: []string{"get", "GB1-ID"},
			want: "hello (GB1) [GB1-ID] / United Kingdom\n",
		},
		{
			name:   "display invalid DC",
			args:   []string{"get", "UNPOG1"},
			stderr: "Error: unknown datacentre\n",
			err:    "unknown datacentre",
		},
	}

	cmd := dataCentersCmd(mockDataCentersClient{dcs: fixtureDataCenters})
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
				assert.Equal(t, tt.want, stdout.String())
			case tt.err != "":
				assert.Equal(t, tt.err, err.Error())
			default:
				t.Fatal(err)
			}
		})
	}
}
