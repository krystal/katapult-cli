package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/core"
	"github.com/krystal/katapult-cli/cmd/katapult/mocks"
)

func singleResponse(t *testing.T, path, key string, body interface{}) core.RequestMaker {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockRequestMaker(ctrl)
	matcher := mocks.KatapultRequestMatcher{
		Path: path,
	}

	mock.EXPECT().
		Do(gomock.Any(), matcher, gomock.Any()).
		DoAndReturn(func(_ context.Context, _ *katapult.Request, iface interface{}) (*katapult.Response, error) {
			b, err := json.Marshal(map[string]interface{}{
				key: body,
			})
			if err != nil {
				return nil, err
			}
			if err := json.Unmarshal(b, iface); err != nil {
				return nil, err
			}
			return mocks.MockOKJSON(b), nil
		}).
		Times(1)

	return mock
}
