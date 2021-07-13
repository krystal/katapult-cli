package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
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

var organizations = []*core.Organization{
	{
		ID:        "loge",
		Name:      "Loge Enthusiasts",
		SubDomain: "loge",
	},
	{
		ID:        "testing",
		Name:      "testing, testing, 123",
		SubDomain: "test",
	},
}

// Used to mock the API client for the purpose of unit tests.
type mockAPIClient struct{}

func (mockAPIClient) Do(
	_ context.Context,
	request *katapult.Request,
	v interface{},
) (*katapult.Response, error) {
	// Defines URL starts for multiple result routes starting with the same string.
	dcStart := "/core/v1/data_centers"

	// Handle creating un-paginated responses for JSON 200 OK's.
	okJSON := func(b []byte) *katapult.Response {
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

	// Process the path.
	path := request.URL.Path
	switch {
	case path == "/core/v1/organizations":
		// Handle listing organizations.
		b, err := json.Marshal(map[string]interface{}{
			"organizations": organizations,
		})
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(b, v); err != nil {
			return nil, err
		}
		return okJSON(b), nil
	case path == "/core/v1/organizations/_/available_networks":
		// Handle listing all networks.
		var networks []*core.Network
		if id := request.URL.Query().Get("organization[id]"); id == "pog-id" {
			networks = idNetworks
		} else if id = request.URL.Query().Get("organization[sub_domain]"); id == "pog-subdomain" {
			networks = subdomainNetworks
		} else {
			return nil, katapult.ErrNotFound
		}
		b, err := json.Marshal(map[string]interface{}{
			"networks": networks,
		})
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(b, v); err != nil {
			return nil, err
		}
		return okJSON(b), nil
	case strings.HasPrefix(path, dcStart):
		// Handle getting all DC's.
		if path == dcStart {
			b, err := json.Marshal(map[string]interface{}{
				"data_centers": dcs,
			})
			if err != nil {
				return nil, err
			}
			if err := json.Unmarshal(b, v); err != nil {
				return nil, err
			}
			return okJSON(b), nil
		}

		// Validate the next part is a slash.
		if path[len(dcStart)] != '/' {
			return nil, katapult.ErrNotFound
		}

		// Handle getting one DC.
		dcID := request.URL.Query().Get("data_center[permalink]")
		if dcID == "" {
			return nil, katapult.ErrBadRequest
		}
		for _, x := range dcs {
			if x.ID == dcID {
				b, err := json.Marshal(map[string]interface{}{
					"data_center": x,
				})
				if err != nil {
					return nil, err
				}
				if err := json.Unmarshal(b, v); err != nil {
					return nil, err
				}
				return okJSON(b), nil
			}
		}
		return nil, errors.New("unknown datacentre")
	default:
		return nil, katapult.ErrNotFound
	}
}
