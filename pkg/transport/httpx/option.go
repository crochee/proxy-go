// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/6/3

// Package httpx
package httpx

import (
	"crypto/tls"

	"github.com/crochee/proxy-go/pkg/logger"
)

type Option func(*option)

type option struct {
	tlsConfig  *tls.Config
	requestLog logger.Builder

	beforeStart []func() error
	afterStart  []func() error

	beforeStop []func() error
	afterStop  []func() error
}

// TlsConfig
func TlsConfig(cfg *tls.Config) Option {
	return func(o *option) { o.tlsConfig = cfg }
}

// RequestLog
func RequestLog(log logger.Builder) Option {
	return func(o *option) { o.requestLog = log }
}

// BeforeStart
func BeforeStart(fs ...func() error) Option {
	return func(o *option) { o.beforeStart = fs }
}

// AfterStart
func AfterStart(fs ...func() error) Option {
	return func(o *option) { o.afterStart = fs }
}

// BeforeStop
func BeforeStop(fs ...func() error) Option {
	return func(o *option) { o.beforeStop = fs }
}

func AfterStop(fs ...func() error) Option {
	return func(o *option) { o.afterStop = fs }
}
