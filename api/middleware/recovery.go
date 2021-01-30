// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/30

package middleware

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"proxy-go/logger"
	"runtime/debug"
	"strings"
)

// Recovery panic log
func Recovery(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			var brokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				if se, ok := ne.Err.(*os.SyscallError); ok {
					brokenPipe = strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
						strings.Contains(strings.ToLower(se.Error()), "connection reset by peer")
				}
			}
			httpRequest, _ := httputil.DumpRequest(ctx.Request, false)
			logger.Errorf("[Recovery] %v\n%v\n%v", string(httpRequest), err, string(debug.Stack()))
			if brokenPipe {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
				return
			}
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}()
	ctx.Next()
}
