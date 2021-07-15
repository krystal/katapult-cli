package mocks

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/krystal/go-katapult"
)

func MockOKJSON(b []byte) *katapult.Response {
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
