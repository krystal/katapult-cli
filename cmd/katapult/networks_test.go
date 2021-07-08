package main

import (
	"bytes"
	"encoding/json"
	"github.com/krystal/go-katapult/core"
	"github.com/krystal/katapult-cli/config"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"testing"
)

var idNetworks = []*core.Network{
	{
		ID:         "pognet",
		Name:       "Pognet 1",
		Permalink:  "pog-1",
		DataCenter: &core.DataCenter{
			ID:        "POG1",
			Name:      "Pogland 1",
			Permalink: "pog1",
			Country:   &core.Country{
				ID:   "pog",
				Name: "Pogland",
			},
		},
	},
	{
		ID:         "pognet2",
		Name:       "Pognet 2",
		Permalink:  "pog-2",
		DataCenter: &core.DataCenter{
			ID:        "POG1",
			Name:      "Pogland 1",
			Permalink: "pog1",
			Country:   &core.Country{
				ID:   "pog",
				Name: "Pogland",
			},
		},
	},
}

var subdomainNetworks = []*core.Network{
	{
		ID:         "pognet3",
		Name:       "Pognet 3",
		Permalink:  "pog-3",
		DataCenter: &core.DataCenter{
			ID:        "POG1",
			Name:      "Pogland 1",
			Permalink: "pog1",
			Country:   &core.Country{
				ID:   "pog",
				Name: "Pogland",
			},
		},
	},
	{
		ID:         "pognet4",
		Name:       "Pognet 4",
		Permalink:  "pog-4",
		DataCenter: &core.DataCenter{
			ID:        "POG1",
			Name:      "Pogland 1",
			Permalink: "pog1",
			Country:   &core.Country{
				ID:   "pog",
				Name: "Pogland",
			},
		},
	},
}

func mockNetworksServer() (net.Listener, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:7654")
	if err != nil {
		return nil, err
	}
	mux := &http.ServeMux{}
	notFound := func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(404)
		_, _ = writer.Write([]byte("Not found."))
	}
	mux.HandleFunc("/", notFound)
	mux.HandleFunc("/core/v1/organizations/_/available_networks", func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/core/v1/organizations/_/available_networks" {
			notFound(writer, request)
			return
		}
		var networks []*core.Network
		if id := request.URL.Query().Get("organization[id]"); id == "pog-id" {
			networks = idNetworks
		} else if id = request.URL.Query().Get("organization[sub_domain]"); id == "pog-subdomain" {
			networks = subdomainNetworks
		} else {
			notFound(writer, request)
			return
		}
		b, err := json.Marshal(map[string]interface{}{"networks": networks})
		if err != nil {
			panic(err)
		}
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(200)
		_, _ = writer.Write(b)
	})
	go func() { _ = http.Serve(ln, mux) }()
	return ln, nil
}

var expectedNetworkResults = []struct{
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

func TestNetworkList(t *testing.T) {
	// Create the mock server.
	ln, err := mockNetworksServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	client, err := newClient(&config.Config{APIURL: "http://127.0.0.1:7654", APIKey: "a"})
	if err != nil {
		t.Fatal(err)
	}

	// Validate the get command.
	for _, v := range expectedNetworkResults {
		cmd := networksCmd(client)
		stdout := &bytes.Buffer{}
		cmd.SetOut(stdout)
		cmd.SetArgs(v.args)
		if err := cmd.Execute(); err == nil {
			assert.Equal(t, v.wants, stdout.String())
		} else if v.err != "" {
			assert.Equal(t, v.err, err.Error())
		} else {
			t.Fatal(err)
		}
	}
}
