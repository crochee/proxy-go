// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/8

package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"proxy-go/api/response"
	"proxy-go/config/dynamic"
	"proxy-go/internal"
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
// @Param request body dynamic.Switch true "switch config"
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
		Name: middlewares.CompleteAction(middlewares.Switcher, middlewares.Update),
		Content: &dynamic.Config{
			Switcher: dynamicSwitch,
		},
	}
	ctx.Status(http.StatusOK)
}

// ListSwitch godoc
// @Summary list switch
// @Description list switch middleware config
// @Tags middleware
// @Accept application/json
// @Produce  application/json
// @Success 200 {array} dynamic.Switch
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/mid/switch [get]
func ListSwitch(ctx *gin.Context) {
	if server.GlobalWatcher == nil {
		response.ErrorWithMessage(ctx, "please check server")
		return
	}
	server.GlobalWatcher.Entry() <- &server.Message{
		Name: middlewares.CompleteAction(middlewares.Switcher, middlewares.List),
	}
	tc := internal.AcquireTimer(30 * time.Second)
	var (
		err  error
		resp interface{}
		ok   bool
	)
	select {
	case <-ctx.Request.Context().Done():
		err = ctx.Request.Context().Err()
	case resp, ok = <-server.GlobalWatcher.Out():
		if !ok {
			err = errors.New("chan is closed")
		}
	case <-tc.C:
		err = errors.New("time out")
	}
	internal.ReleaseTimer(tc)
	if err != nil {
		response.GinError(ctx, response.ErrorWith(http.StatusInternalServerError, err))
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// UpdateRateLimit godoc
// @Summary update rate limit
// @Description update rate limit middleware config
// @Tags middleware
// @Accept application/json
// @Produce  application/json
// @Param request body dynamic.RateLimit true "rate limit config"
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
		Name: middlewares.CompleteAction(middlewares.RateLimiter, middlewares.Update),
		Content: &dynamic.Config{
			Limit: &rateLimit,
		},
	}
	ctx.Status(http.StatusOK)
}
