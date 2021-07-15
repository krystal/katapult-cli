package mocks

import "net/url"

// URLParamMatcher is used to match a URL parameter.
type URLParamMatcher interface {
	// MatchParam is used to match a query param.
	MatchParam(v url.Values) bool
}

type matchContains string

func (m matchContains) MatchParam(v url.Values) bool {
	x := v.Get(string(m))
	return x != ""
}

// URLParamContains is used to check if the URL parameter contains a key.
func URLParamContains(key string) URLParamMatcher {
	return matchContains(key)
}

type paramOr []URLParamMatcher

func (p paramOr) MatchParam(v url.Values) bool {
	for _, x := range p {
		if x.MatchParam(v) {
			return true
		}
	}
	return false
}

// URLParamOr is used to check if one or more of the params are true.
func URLParamOr(params ...URLParamMatcher) URLParamMatcher {
	x := make(paramOr, len(params))
	copy(x, params)
	return x
}
