// Package transport
package transport

import (
	"context"
	"os"
)

type option struct {
	signals    []os.Signal
	serverList []AppServer
	ctx        context.Context
}

// Signal with exit signals.
func Signal(sigList ...os.Signal) func(*option) {
	return func(o *option) { o.signals = sigList }
}

// Servers with transport servers.
func Servers(servers ...AppServer) func(*option) {
	return func(o *option) { o.serverList = servers }
}

// Context with service context.
func Context(ctx context.Context) func(*option) {
	return func(o *option) { o.ctx = ctx }
}
