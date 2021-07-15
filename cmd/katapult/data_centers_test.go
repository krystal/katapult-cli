package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
	"github.com/krystal/katapult-cli/cmd/katapult/mocks"

	"github.com/stretchr/testify/assert"
)

var dcs = []*core.DataCenter{
	{
		ID:        "POG1",
		Name:      "hello",
		Permalink: "Hello World!",
		Country: &core.Country{
			ID:   "POG",
			Name: "Pogland",
		},
	},
	{
		ID:        "GB1",
		Name:      "hello",
		Permalink: "Hello World!",
		Country: &core.Country{
			ID:   "UK",
			Name: "United Kingdom",
		},
	},
}

const stdoutDcList = ` - hello (Hello World!) [POG1] / Pogland
 - hello (Hello World!) [GB1] / United Kingdom
`

func TestDataCenters_List(t *testing.T) {
	mock := singleResponse(t, "/core/v1/data_centers", "data_centers", dcs)
	cmd := dataCentersCmd(mock)
	stdout := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetArgs([]string{"list"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, stdoutDcList, stdout.String())
}

func TestDataCenters_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockRequestMaker(ctrl)
	matcher := mocks.KatapultRequestMatcher{
		Path: "/core/v1/data_centers/_",
		ExpectedParams: []mocks.URLParamMatcher{
			mocks.URLParamContains("data_center[permalink]"),
		},
	}
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
			name:   "display invalid DC",
			args:   []string{"get", "UNPOG1"},
			stderr: "Error: unknown datacentre\n",
			err:    "unknown datacentre",
		},
	}

	mock.
		EXPECT().
		Do(gomock.Any(), matcher, gomock.Any()).
		DoAndReturn(func(_ context.Context, req *katapult.Request, iface interface{}) (*katapult.Response, error) {
			permalink := req.URL.Query().Get("data_center[permalink]")
			for _, v := range dcs {
				if v.ID == permalink {
					b, err := json.Marshal(map[string]interface{}{
						"data_center": v,
					})
					if err != nil {
						return nil, err
					}
					if err = json.Unmarshal(b, iface); err != nil {
						return nil, err
					}
					return mocks.MockOKJSON(b), nil
				}
			}
			return nil, fmt.Errorf("unknown datacentre")
		}).
		Times(len(tests))

	cmd := dataCentersCmd(mock)
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
