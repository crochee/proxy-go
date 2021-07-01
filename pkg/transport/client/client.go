// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/7/1

// Package client
package client

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

//go:generate mockgen -destination client_mock.go -package client -source client.go

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

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
