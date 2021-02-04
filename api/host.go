// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/30

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"proxy-go/api/response"
	"proxy-go/model"
)

func UpdateHost(ctx *gin.Context) {
	var replaceHost model.ReplaceHost
	if err := ctx.ShouldBindBodyWith(&replaceHost, binding.JSON); err != nil {
		response.GinError(ctx, response.ErrorWith(http.StatusBadRequest, err))
		return
	}
	ctx.JSON(http.StatusOK, replaceHost)
}
