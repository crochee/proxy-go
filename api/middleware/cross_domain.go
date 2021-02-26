// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/30

package middleware

import (
	"github.com/rs/cors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CrossDomain skip the cross-domain phase
func CrossDomain(ctx *gin.Context) {
	cors.Default()
	ctx.Header("Access-Control-Allow-Headers", "Content-Type")
	ctx.Header("Access-Control-Allow-Origin", origin(ctx))
	if ctx.Request.Method == http.MethodOptions {
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,DELETE,PUT,PATCH,OPTIONS")
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.Next()
}

func origin(ctx *gin.Context) string {
	origin := ctx.GetHeader("Origin")
	if origin != "" {
		return origin
	}
	return "*"
}
