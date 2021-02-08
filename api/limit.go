// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: l30002214
// Create: 2021/2/8

// Package main
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"proxy-go/api/response"
	"proxy-go/config/dynamic"
	"proxy-go/middlewares"
	"proxy-go/server"
)

func UpdateRateLimit(ctx *gin.Context) {
	var rateLimit dynamic.RateLimit
	if err := ctx.ShouldBindBodyWith(&rateLimit, binding.JSON); err != nil {
		response.GinError(ctx, response.ErrorWith(http.StatusBadRequest, err))
		return
	}
	if server.GlobalWatcher == nil {
		response.ErrorWithMessage(ctx, "please check server")
		return
	}
	server.GlobalWatcher.Entry() <- &server.Message{
		Name: middlewares.Switcher,
		Content: &dynamic.Config{
			Limit: &rateLimit,
		},
	}
	ctx.Status(http.StatusOK)
}
