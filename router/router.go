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

func Route(ctx context.Context, cfg *config.Config) (http.Handler, error) {
	proxy, err := middlewares.BuildProxy(0, *cfg.Server.Medata[0])
	if err != nil {
		return nil, err
	}
	return middlewares.NewRecovery(ctx, proxy)
}
