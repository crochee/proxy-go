// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/7/1

// Package client
package client

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/url"
)

var defaultClient Client

// DefaultClient
func DefaultClient(opts ...Option) {
	defaultClient = NewStandardClient(opts...)
}

// Send
func Send(ctx context.Context, method string, uri string,
	body []byte, headers map[string]string) (*http.Response, error) {
	req, err := NewRequest(ctx, method, uri, body, headers)
	if err != nil {
		return nil, err
	}
	return Do(req)
}

// Do
func Do(req *http.Request) (*http.Response, error) {
	return defaultClient.Do(req)
}

// NewRequest
func NewRequest(ctx context.Context, method string, uri string,
	body []byte, headers map[string]string) (*http.Request, error) {
	tempUri, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if tempUri.Hostname() == "" {
		return nil, errors.New("the url hasn't ip or domain name")
	}
	tempUri.RawQuery = tempUri.Query().Encode()
	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, method, tempUri.String(), bytes.NewReader(body)); err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	return req, nil
}
