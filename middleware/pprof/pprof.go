// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/5/7

// Package pprof
package pprof

import (
	"net/http"
	"net/http/pprof"

	"github.com/crochee/proxy-go/middleware"
)

type index struct {
}

func NewIndex() middleware.Handler {
	return index{}
}

func (index) NameSpace() string {
	return "index"
}

func (index) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	pprof.Index(writer, request)
}

type cmdline struct {
}

func NewCmdline() middleware.Handler {
	return cmdline{}
}

func (cmdline) NameSpace() string {
	return "cmdline"
}

func (cmdline) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	pprof.Cmdline(writer, request)
}

type profile struct {
}

func NewProfile() middleware.Handler {
	return profile{}
}

func (profile) NameSpace() string {
	return "profile"
}

func (profile) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	pprof.Profile(writer, request)
}

type symbol struct {
}

func NewSymbol() middleware.Handler {
	return symbol{}
}

func (symbol) NameSpace() string {
	return "symbol"
}

func (symbol) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	pprof.Symbol(writer, request)
}

type trace struct {
}

func NewTrace() middleware.Handler {
	return trace{}
}

func (trace) NameSpace() string {
	return "trace"
}

func (trace) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	pprof.Trace(writer, request)
}
