// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/7/1

// Package client
package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestNewMockClient(t *testing.T) {
	type Arg struct {
		method  string
		url     string
		body    []byte
		headers map[string]string
	}
	testList := []struct {
		name     string
		input    Arg
		expected *http.Response
	}{
		{
			name: "mock",
			input: Arg{
				method:  http.MethodPost,
				url:     "https://www.baidu.com",
				body:    nil,
				headers: nil,
			},
			expected: nil,
		},
	}
	for _, tt := range testList {
		t.Run(tt.name, func(t *testing.T) {
			req, err := NewRequest(context.Background(), tt.input.method, tt.input.url, tt.input.body, tt.input.headers)
			require.NoError(t, err)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := NewMockClient(ctrl)

			c.EXPECT().Do(gomock.Any()).DoAndReturn(func(reqC *http.Request) (*http.Response, error) {
				return nil, nil
			})

			resp, err := c.Do(req)
			require.NoError(t, err)
			require.Equal(t, tt.expected, resp)
		})
	}
}
