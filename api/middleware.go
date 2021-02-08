// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/8

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"proxy-go/api/response"
	"proxy-go/config/dynamic"
	"proxy-go/logger"
	"proxy-go/middlewares"
	"proxy-go/server"
)

// UpdateSwitch godoc
// @Summary update switch
// @Description update switch middleware config
// @Tags middleware
// @Accept application/json
// @Produce  application/json
// @Param request body model.Switch true "switch config"
// @Success 200
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/mid/switch [post]
func UpdateSwitch(ctx *gin.Context) {
	dynamicSwitch := &dynamic.Switch{
		Node: dynamic.BalanceNode{
			Scheme: "http", // 默认值 http
			Weight: 1.0,    // 1.0
		},
	}
	if err := ctx.ShouldBindBodyWith(dynamicSwitch, binding.JSON); err != nil {
		response.GinError(ctx, response.ErrorWith(http.StatusBadRequest, err))
		return
	}
	if server.GlobalWatcher == nil {
		response.ErrorWithMessage(ctx, "please check server")
		return
	}
	logger.FromContext(ctx.Request.Context()).Debugf("%+v", dynamicSwitch)
	server.GlobalWatcher.Entry() <- &server.Message{
		Name: middlewares.Switcher,
		Content: &dynamic.Config{
			Switcher: dynamicSwitch,
		},
	}
	ctx.Status(http.StatusOK)
}

// UpdateRateLimit godoc
// @Summary update rate limit
// @Description update rate limit middleware config
// @Tags middleware
// @Accept application/json
// @Produce  application/json
// @Param request body model.RateLimit true "rate limit config"
// @Success 200
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/mid/limit [post]
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
