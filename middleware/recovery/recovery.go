// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package recovery

import (
	"net/http"
	"runtime/debug"

	"github.com/crochee/proxy-go/internal"
	"github.com/crochee/proxy-go/logger"
	"github.com/crochee/proxy-go/middleware"
)

type recovery struct {
	next http.Handler
}

// New creates recovery middleware
func New(next http.Handler) middleware.Handler {
	return &recovery{
		next: next,
	}
}

func (re *recovery) NameSpace() string {
	return "Recovery"
}

func (re *recovery) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log := logger.FromContext(req.Context())
			if err == http.ErrAbortHandler {
				log.Debugf("Request has been aborted [%s - %s]: %v", req.RemoteAddr, req.URL, err)
				return
			}
			log.Errorf("[Recovery] from panic in HTTP handler [%s - %s]: %+v\nStack:\n%s",
				req.RemoteAddr, req.URL, err, debug.Stack())

			http.Error(rw, internal.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}()
	re.next.ServeHTTP(rw, req)
}
