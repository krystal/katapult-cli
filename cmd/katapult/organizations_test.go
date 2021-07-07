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

var organizations = []*core.Organization{
	{
		ID: "loge",
		Name: "Loge Enthusiasts",
		SubDomain: "loge",
	},
	{
		ID: "testing",
		Name: "testing, testing, 123",
		SubDomain: "test",
	},
}

func mockOrganizationsServer() (net.Listener, error) {
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
	mux.HandleFunc("/core/v1/organizations", func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/core/v1/organizations" {
			notFound(writer, request)
			return
		}
		b, err := json.Marshal(map[string]interface{}{"organizations": organizations})
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

const stdoutOrganizationsList = ` - Loge Enthusiasts (loge) [loge]
 - testing, testing, 123 (test) [testing]
`

func TestOrganizations_List(t *testing.T) {
	// Create the mock server.
	ln, err := mockOrganizationsServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	client, err := newClient(&config.Config{APIURL: "http://127.0.0.1:3210", APIKey: "a"})
	if err != nil {
		t.Fatal(err)
	}

	// Validate the list command.
	cmd := organizationsCmd(client)
	stdout := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetArgs([]string{"list"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, stdoutOrganizationsList, stdout.String())
}
