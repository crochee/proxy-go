// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/7/1

// Package client
package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

var defaultClient Client

type Option func(*option)

type option struct {
	t *http.Transport
}

// TlsConfig
func TlsConfig(cfg *tls.Config) Option {
	return func(o *option) { o.t.TLSClientConfig = cfg }
}

// Timeout
func Timeout(t time.Duration) Option {
	return func(o *option) { o.t.ResponseHeaderTimeout = t }
}

// DefaultClient
func DefaultClient(opts ...Option) {
	defaultClient = NewStandardClient(opts...)
}

func NewStandardClient(opts ...Option) *standardClient {
	o := &option{
		t: &http.Transport{
			MaxIdleConnsPerHost: 100,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 300 * time.Second,
			ForceAttemptHTTP2:     true,
		},
	}
	for _, opt := range opts {
		opt(o)
	}
	return &standardClient{client: &http.Client{Transport: o.t}}
}

type standardClient struct {
	client *http.Client
}

func (s *standardClient) Do(req *http.Request) (*http.Response, error) {
	return s.client.Do(req)
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
