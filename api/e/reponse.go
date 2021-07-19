package e

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/crochee/proxy-go/pkg/logger"
)

var (
	bundle   *i18n.Bundle
	messages = &i18n.Message{
		Description: "The content of the error message",
		Other:       "{{.}}",
	}
)

func init() {
	bundle = i18n.NewBundle(language.English)
}

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
	resp.Message = i18n.NewLocalizer(bundle, ctx.Request.Header.Get("Accept-Language")).
		MustLocalize(&i18n.LocalizeConfig{
			TemplateData:   code.Detail(),
			DefaultMessage: messages,
		})
	return resp
}

func respWith(ctx *gin.Context, code Code, message string) *Response {
	resp := &Response{
		Code: code.ErrorCode(),
	}
	loc := i18n.NewLocalizer(bundle, ctx.Request.Header.Get("Accept-Language"))
	// Message
	cfg := &i18n.LocalizeConfig{
		TemplateData:   code.Detail(),
		DefaultMessage: messages,
	}
	resp.Message = loc.MustLocalize(cfg)
	// Extra
	cfg.TemplateData = message
	resp.Extra = loc.MustLocalize(cfg)
	return resp
}
