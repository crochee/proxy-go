package e

import (
	"errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"

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
	resp := respWith(ctx, code, message)
	logger.FromContext(ctx.Request.Context()).Errorf("[ERROR] error_code:%s,message:%s,extra:%s",
		resp.Code, resp.Message, resp.Extra)
	ctx.JSON(code.Status(), resp)
}

// WithCode gin Response with Code
func WithCode(ctx *gin.Context, code Code) {
	resp := respWithCode(ctx, code)
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
		unwrapHandler, ok := err.(unwrap)
		if ok {
			err = unwrapHandler.Unwrap()
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
	var errorCode *ResponseError
	if errors.As(errResult, &errorCode) {
		WithCodeMessage(ctx, errorCode.Code, errorCode.Message)
		return
	}
	WithCodeMessage(ctx, Unknown, err.Error())
}

// AbortWith gin Response with with Code and message
func AbortWith(ctx *gin.Context, code Code, message string) {
	resp := respWith(ctx, code, message)
	logger.FromContext(ctx.Request.Context()).Errorf("[ABORT] error_code:%s,message:%s,extra:%s",
		resp.Code, resp.Message, resp.Extra)
	ctx.AbortWithStatusJSON(code.Status(), resp)
}

// Abort gin Response with with Code and message
func Abort(ctx *gin.Context, code Code) {
	resp := respWithCode(ctx, code)
	logger.FromContext(ctx.Request.Context()).Errorf("[ABORT] error_code:%s,message:%s", resp.Code, resp.Message)
	ctx.AbortWithStatusJSON(code.Status(), resp)
}

func respWithCode(ctx *gin.Context, code Code) *Response {
	resp := &Response{
		Code: code.ErrorCode(),
	}
	tags, _, err := language.ParseAcceptLanguage(ctx.Request.Header.Get("Accept-Language"))
	if err != nil {
		panic(err)
	}
	for _, tag := range tags {
		if tag.String() == language.Chinese.String() {
			resp.Message = code.Chinese()
			return resp
		}
	}
	resp.Message = code.English()
	return resp
}

func respWith(ctx *gin.Context, code Code, message string) *Response {
	resp := &Response{
		Code:  code.ErrorCode(),
		Extra: message,
	}
	tags, _, err := language.ParseAcceptLanguage(ctx.Request.Header.Get("Accept-Language"))
	if err != nil {
		panic(err)
	}
	for _, tag := range tags {
		if tag.String() == language.SimplifiedChinese.String() {
			resp.Message = code.Chinese()
			return resp
		}
	}
	resp.Message = code.English()
	return resp
}
