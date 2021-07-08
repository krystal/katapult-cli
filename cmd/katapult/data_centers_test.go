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

var dcs = []*core.DataCenter{
	{
		ID:        "POG1",
		Name:      "hello",
		Permalink: "Hello World!",
		Country:   &core.Country{
			ID:   "POG",
			Name: "Pogland",
		},
	},
	{
		ID:        "GB1",
		Name:      "hello",
		Permalink: "Hello World!",
		Country:   &core.Country{
			ID:   "UK",
			Name: "United Kingdom",
		},
	},
}

func mockDcServer() (net.Listener, error) {
	// Create the mock JSON.
	b, err := json.Marshal(map[string][]*core.DataCenter{
		"data_centers": dcs,
	})
	if err != nil {
		return nil, err
	}

	// Serve the mock JSON.
	ln, err := net.Listen("tcp", "127.0.0.1:3210")
	if err != nil {
		return nil, err
	}
	mux := &http.ServeMux{}
	notFound := func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(404)
		_, _ = writer.Write([]byte("Not found."))
	}
	mux.HandleFunc("/", notFound)
	mux.HandleFunc("/core/v1/data_centers", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(200)
		_, _ = writer.Write(b)
	})
	startPath := "/core/v1/data_centers/"
	mux.HandleFunc(startPath, func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == startPath {
			notFound(writer, request)
			return
		}
		dcId := request.URL.Query().Get("data_center[permalink]")
		if dcId == "" {
			notFound(writer, request)
			return
		}
		for _, v := range dcs {
			if v.ID == dcId {
				writer.Header().Add("Content-Type", "application/json")
				writer.WriteHeader(200)
				b, err := json.Marshal(map[string]*core.DataCenter{
					"data_center": v,
				})
				if err != nil {
					panic(err)
				}
				_, _ = writer.Write(b)
				return
			}
		}
		notFound(writer, request)
	})
	go func() { _ = http.Serve(ln, mux) }()
	return ln, nil
}

const stdoutDcList = ` - hello (Hello World!) [POG1] / Pogland
 - hello (Hello World!) [GB1] / United Kingdom
`

func TestDataCenters_List(t *testing.T) {
	// Create the mock server.
	ln, err := mockDcServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	client, err := newClient(&config.Config{APIURL: "http://127.0.0.1:3210", APIKey: "a"})
	if err != nil {
		t.Fatal(err)
	}

	// Validate the list command.
	cmd := dataCentersCmd(client)
	stdout := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetArgs([]string{"list"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, stdoutDcList, stdout.String())
}

var expectedDcResults = []struct{
	name string

	args  []string
	wants string
	err   string
}{
	{
		name:   "display POG1",
		args:   []string{"get", "POG1"},
		wants:  `hello (Hello World!) [POG1] / Pogland
`,
	},
	{
		name:   "display GB1",
		args:   []string{"get", "GB1"},
		wants:  `hello (Hello World!) [GB1] / United Kingdom
`,
	},
	{
		name: "display invalid DC",
		args: []string{"get", "UNPOG1"},
		err:  "unknown datacentre",
	},
}

func TestDataCenters_Get(t *testing.T) {
	// Create the mock server.
	ln, err := mockDcServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	client, err := newClient(&config.Config{APIURL: "http://127.0.0.1:3210", APIKey: "a"})
	if err != nil {
		t.Fatal(err)
	}

	// Validate the get command.
	cmd := dataCentersCmd(client)
	for _, v := range expectedDcResults {
		stdout := &bytes.Buffer{}
		cmd.SetOut(stdout)
		cmd.SetErr(&bytes.Buffer{}) // Ignore stderr, this is just testing the command framework,, not the error.
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
