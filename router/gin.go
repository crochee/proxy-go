// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/30

package router

import (
	"net/http"
	"runtime/pprof"
)

// @title obs Swagger API
// @version 1.0
// @description This is a obs server.

// NewGinEngine gin router
func NewGinEngine() http.Handler {

	pprof.StartCPUProfile()

}
