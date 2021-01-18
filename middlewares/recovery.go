// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package middlewares

import (
	"context"
	"net/http"
	"runtime/debug"

	"proxy-go/logger"
)

type recovery struct {
	next http.Handler
	ctx  context.Context
}

// NewRecovery creates recovery middleware
func NewRecovery(ctx context.Context, next http.Handler) (http.Handler, error) {
	return &recovery{
		next: next,
		ctx:  ctx,
	}, nil
}

func (re *recovery) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
		if err := recover(); err != nil {
			log := logger.FromContext(ctx)
			if err == http.ErrAbortHandler {
				log.Debugf("Request has been aborted [%s - %s]: %v", r.RemoteAddr, r.URL, err)
				return
			}

			log.Errorf("[Recovery] from panic in HTTP handler [%s - %s]: %+v", r.RemoteAddr, r.URL, err)

			log.Errorf("Stack: %s", debug.Stack())

			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}(re.ctx, rw, req)
	re.next.ServeHTTP(rw, req)
}