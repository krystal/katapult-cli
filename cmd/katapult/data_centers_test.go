package main

import (
	"encoding/json"
	"fmt"
	"github.com/krystal/go-katapult/core"
	"github.com/krystal/katapult-cli/config"
	"net"
	"net/http"
	"testing"
)

func mockDcListServer() (net.Listener, error) {
	// Create the mock JSON.
	b, err := json.Marshal(map[string][]*core.DataCenter{
		"data_centers": {
			{
				ID:        "test",
				Name:      "hello",
				Permalink: "Hello World!",
				Country:   nil,
			},
			{
				ID:        "test",
				Name:      "hello",
				Permalink: "Hello World!",
				Country:   &core.Country{
					ID:   "UK",
					Name: "United Kingdom",
				},
			},
		},
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
	_ = http.Serve(ln, mux)
	return ln, nil
}

func TestDataCenters_List(t *testing.T) {
	ln, err := mockDcListServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	client, err := newClient(&config.Config{APIURL: "http://127.0.0.1:3210"})
	if err != nil {
		t.Fatal(err)
	}
	// TODO: This.
	fmt.Println(client)
}
