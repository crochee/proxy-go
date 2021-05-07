// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/5/6

// Package transport
package transport

import (
	"context"
	"os"
)

type option struct {
	sigList    []os.Signal
	serverList []Server
	ctx        context.Context
}

// Signal with exit signals.
func Signal(sigList ...os.Signal) func(*option) {
	return func(o *option) { o.sigList = sigList }
}

// Servers with transport servers.
func Servers(servers ...Server) func(*option) {
	return func(o *option) { o.serverList = servers }
}

// Context with service context.
func Context(ctx context.Context) func(*option) {
	return func(o *option) { o.ctx = ctx }
}
