package main

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
	"github.com/krystal/katapult-cli/cmd/katapult/mocks"

	"github.com/stretchr/testify/assert"
)

var idNetworks = []*core.Network{
	{
		ID:        "pognet",
		Name:      "Pognet 1",
		Permalink: "pog-1",
		DataCenter: &core.DataCenter{
			ID:        "POG1",
			Name:      "Pogland 1",
			Permalink: "pog1",
			Country: &core.Country{
				ID:   "pog",
				Name: "Pogland",
			},
		},
	},
	{
		ID:        "pognet2",
		Name:      "Pognet 2",
		Permalink: "pog-2",
		DataCenter: &core.DataCenter{
			ID:        "POG1",
			Name:      "Pogland 1",
			Permalink: "pog1",
			Country: &core.Country{
				ID:   "pog",
				Name: "Pogland",
			},
		},
	},
}

var subdomainNetworks = []*core.Network{
	{
		ID:        "pognet3",
		Name:      "Pognet 3",
		Permalink: "pog-3",
		DataCenter: &core.DataCenter{
			ID:        "POG1",
			Name:      "Pogland 1",
			Permalink: "pog1",
			Country: &core.Country{
				ID:   "pog",
				Name: "Pogland",
			},
		},
	},
	{
		ID:        "pognet4",
		Name:      "Pognet 4",
		Permalink: "pog-4",
		DataCenter: &core.DataCenter{
			ID:        "POG1",
			Name:      "Pogland 1",
			Permalink: "pog1",
			Country: &core.Country{
				ID:   "pog",
				Name: "Pogland",
			},
		},
	},
}

func TestNetworkList(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockRequestMaker(ctrl)
	matcher := mocks.KatapultRequestMatcher{
		Path: "/core/v1/organizations/_/available_networks",
		ExpectedParams: []mocks.URLParamMatcher{
			mocks.URLParamOr(
				mocks.URLParamContains("organization[id]"),
				mocks.URLParamContains("organization[sub_domain]")),
		},
	}

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

	mock.
		EXPECT().
		Do(gomock.Any(), matcher, gomock.Any()).
		DoAndReturn(func(_ context.Context, req *katapult.Request, iface interface{}) (*katapult.Response, error) {
			params := req.URL.Query()
			x := params.Get("organization[id]")
			if x == "pog-id" {
				b, err := json.Marshal(map[string]interface{}{
					"networks": idNetworks,
				})
				if err != nil {
					return nil, err
				}
				if err = json.Unmarshal(b, iface); err != nil {
					return nil, err
				}
				return mocks.MockOKJSON(b), nil
			}
			x = params.Get("organization[sub_domain]")
			if x == "pog-subdomain" {
				b, err := json.Marshal(map[string]interface{}{
					"networks": subdomainNetworks,
				})
				if err != nil {
					return nil, err
				}
				if err = json.Unmarshal(b, iface); err != nil {
					return nil, err
				}
				return mocks.MockOKJSON(b), nil
			}
			return nil, katapult.ErrNotFound
		}).
		Times(len(tests))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := networksCmd(mock)
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
