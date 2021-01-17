// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package router

import (
	"context"
	"net/http"

	"proxy-go/config"
	"proxy-go/middlewares"
)

func Route(ctx context.Context) (http.Handler, error) {
	proxy, err := middlewares.BuildProxy(0)
	if err != nil {
		return nil, err
	}
	var rh http.Handler
	if rh, err = middlewares.NewReplaceHost(ctx, proxy, *config.Cfg.Middleware.ReplaceHost); err != nil {
		return nil, err
	}
	return middlewares.NewRecovery(ctx, rh)
}
