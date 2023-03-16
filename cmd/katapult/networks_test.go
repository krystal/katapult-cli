package main

import (
	"context"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
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

var idVirtualNetworks = []*core.VirtualNetwork{
	{
		ID:   "pognet-virtual-1",
		Name: "Pognet Virtual Network 1",
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

type mockNetworkList struct{}

func (mockNetworkList) List(
	_ context.Context, org core.OrganizationRef, reqOpts ...katapult.RequestOption,
) ([]*core.Network, []*core.VirtualNetwork, *katapult.Response, error) {
	switch {
	case org.SubDomain == "pog-subdomain":
		return subdomainNetworks, nil, nil, nil
	case org.ID == "pog-id":
		return idNetworks, idVirtualNetworks, nil, nil
	default:
		return nil, nil, nil, katapult.ErrNotFound
	}
}

func TestNetworks_List(t *testing.T) {
	tests := []struct {
		name string

		args    []string
		output  string
		stderr  string
		wantErr string
	}{
		{
			name: "Test listing pog-id human readable",
			args: []string{"ls", "--id", "pog-id"},
		},
		{
			name:   "Test listing pog-id json",
			args:   []string{"ls", "--id", "pog-id"},
			output: "json",
		},
		{
			name: "Test listing pog-subdomain human readable",
			args: []string{"ls", "--subdomain", "pog-subdomain"},
		},
		{
			name:   "Test listing pog-subdomain json",
			args:   []string{"ls", "--subdomain", "pog-subdomain"},
			output: "json",
		},
		{
			name:    "No flags provided",
			args:    []string{"ls"},
			stderr:  "Error: both ID and subdomain are unset\n",
			wantErr: "both ID and subdomain are unset",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := networksCmd(mockNetworkList{})
			cmd.SetArgs(tt.args)
			outputFlag = tt.output
			assertCobraCommand(t, cmd, tt.wantErr, tt.stderr)
			outputFlag = ""
		})
	}
}
