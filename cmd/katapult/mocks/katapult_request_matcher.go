package mocks

import "github.com/krystal/go-katapult"

type KatapultRequestMatcher struct {
	Path           string   `json:"path"`
	ExpectedParams []string `json:"expected_params"`
}

func (k KatapultRequestMatcher) String() string {
	return "check if this is a valid request to " + k.Path
}

func (k KatapultRequestMatcher) Matches(iface interface{}) bool {
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
