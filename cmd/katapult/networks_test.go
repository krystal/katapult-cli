package main

import (
	"bytes"
	"context"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"

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

type mockNetworkList struct {}

func (mockNetworkList) List(
	_ context.Context, org core.OrganizationRef,
) ([]*core.Network, []*core.VirtualNetwork, *katapult.Response, error) {
	// TODO
}

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
			cmd := networksCmd(mockNetworkList{})
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
