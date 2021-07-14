package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
	"github.com/krystal/katapult-cli/cmd/katapult/mocks"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const stdoutDcList = ` - hello (Hello World!) [POG1] / Pogland
 - hello (Hello World!) [GB1] / United Kingdom
`

//func TestDataCenters_List(t *testing.T) {
//	cmd := dataCentersCmd(mockAPIClient{})
//	stdout := &bytes.Buffer{}
//	cmd.SetOut(stdout)
//	cmd.SetArgs([]string{"list"})
//	if err := cmd.Execute(); err != nil {
//		t.Fatal(err)
//	}
//	assert.Equal(t, stdoutDcList, stdout.String())
//}

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

type katapultRequestMatcher struct {
	Path string `json:"path"`
	ExpectedParams []string `json:"expected_params"`
}

func (k katapultRequestMatcher) String() string {
	return "check if this is a valid request to " + k.Path
}

func (k katapultRequestMatcher) Matches(iface interface{}) bool {
	req, ok := iface.(*katapult.Request)
	if !ok {
		// This should never happen.
		return false
	}
	if req.URL.Path != k.Path {
		// This is for the wrong path.
		return false
	}
	q := req.URL.Query()
	for _, k := range k.ExpectedParams {
		if q.Get(k) == "" {
			// This query param is not set.
			return false
		}
	}
	return true
}

// TODO: Move!
func okJSON(b []byte) *katapult.Response {
	return &katapult.Response{
		Response: &http.Response{
			StatusCode: 200,
			Header: http.Header{
				"Content-Type": {"application/json"},
			},
			Body:          ioutil.NopCloser(bytes.NewReader(b)),
			ContentLength: int64(len(b)),
		},
	}
}

func TestDataCenters_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockRequestMaker(ctrl)
	matcher := katapultRequestMatcher{
		Path:           "/core/v1/data_centers/_",
		ExpectedParams: []string{"data_center[permalink]"},
	}
	tests := []struct {
		name string

		args  []string
		dc    string
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
					return okJSON(b), nil
				}
			}
			return nil, fmt.Errorf("unknown datacentre")
		}).
		Times(len(tests))
	cmd := dataCentersCmd(mock)
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
