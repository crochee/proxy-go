// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/4

package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// GinError gin response with format err
func GinError(ctx *gin.Context, err error) {
	switch value := err.(type) {
	case *ProxyError:
		ctx.JSON(value.Code.Status(), &Response{
			Code:    value.Code.Value(),
			Message: value.Message,
		})
	default:
		ctx.JSON(http.StatusInternalServerError, &Response{
			Code:    "500",
			Message: value.Error(),
		})
	}
}

// ErrorWithMessage gin response with message
func ErrorWithMessage(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, &Response{
		Code:    "500",
		Message: message,
	})
}
