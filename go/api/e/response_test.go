// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/7/19

// Package e
package e

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/crochee/proxy/internal"
)

func TestAbort(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(func(ctx *gin.Context) {
		if ctx.Request.ContentLength > 1024*1024 {
			AbortWith(ctx, Unknown, "1024")
			return
		} else if ctx.Request.ContentLength > 1024 {
			Abort(ctx, Unknown)
			return
		}
		ctx.Next()
	})
	router.GET("/v1/bucket", func(ctx *gin.Context) {
		WithCode(ctx, Success)
	})

	header := make(http.Header)
	header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	w := internal.PerformRequest(router, http.MethodGet,
		"/v1/bucket", nil, header)
	t.Logf("%+v\nbody:%s", w.Result(), w.Body.String())
}
