package e

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/crochee/proxy-go/pkg/logger"
)

// 封装成统一的字段给前端处理
type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Extra   string `json:"extra"`
}

// WithCodeMessage gin Response with with Code and message
func WithCodeMessage(ctx *gin.Context, code Code, message string) {
	resp := &Response{
		Code:    code.String(),
		Message: code.English(),
		Extra:   message,
	}
	if strings.Contains(ctx.Request.Header.Get("accept-language"), "zh") {
		resp.Message = code.Chinese()
	}
	logger.FromContext(ctx.Request.Context()).Errorf("[ERROR] error_code:%s,message:%s,extra:%s",
		resp.Code, resp.Message, resp.Extra)
	ctx.JSON(code.Status(), resp)
}

// WithCode gin Response with Code
func WithCode(ctx *gin.Context, code Code) {
	resp := &Response{
		Code:    code.String(),
		Message: code.English(),
	}
	if strings.Contains(ctx.Request.Header.Get("accept-language"), "zh") {
		resp.Message = code.Chinese()
	}
	logger.FromContext(ctx.Request.Context()).Errorf("[ERROR] error_code:%s,message:%s", resp.Code, resp.Message)
	ctx.JSON(code.Status(), resp)
}

func Unwrap(err error) error {
	type causer interface {
		Cause() error
	}

	type unwrap interface {
		Unwrap() error
	}

	for err != nil {
		unwrap, ok := err.(unwrap)
		if ok {
			err = unwrap.Unwrap()
			continue
		}
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

// Errors gin Response with error
func Errors(ctx *gin.Context, err error) {
	var errResult error
	if errResult = Unwrap(err); errResult == nil {
		WithCodeMessage(ctx, Success, err.Error())
		return
	}
	var errorCode *serviceError
	if errors.As(errResult, &errorCode) {
		WithCodeMessage(ctx, errorCode.Code, errorCode.Message)
		return
	}
	WithCodeMessage(ctx, Unknown, err.Error())
}

// AbortWith gin Response with with Code and message
func AbortWith(ctx *gin.Context, code Code, message string) {
	resp := &Response{
		Code:    code.String(),
		Message: code.English(),
		Extra:   message,
	}
	if strings.Contains(ctx.Request.Header.Get("accept-language"), "zh") {
		resp.Message = code.Chinese()
	}
	logger.FromContext(ctx.Request.Context()).Errorf("[ABORT] error_code:%s,message:%s,extra:%s",
		resp.Code, resp.Message, resp.Extra)
	ctx.AbortWithStatusJSON(code.Status(), resp)
}

// Abort gin Response with with Code and message
func Abort(ctx *gin.Context, code Code) {
	resp := &Response{
		Code:    code.String(),
		Message: code.English(),
	}
	if strings.Contains(ctx.Request.Header.Get("accept-language"), "zh") {
		resp.Message = code.Chinese()
	}
	logger.FromContext(ctx.Request.Context()).Errorf("[ABORT] error_code:%s,message:%s", resp.Code, resp.Message)
	ctx.AbortWithStatusJSON(code.Status(), resp)
}
