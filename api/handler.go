package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/crochee/proxy-go/api/e"
	"github.com/crochee/proxy-go/config/dynamic"
	"github.com/crochee/proxy-go/pkg/middleware"
	"github.com/crochee/proxy-go/pkg/resource/api"
)

// Handlers godoc
// @Summary Handlers
// @Description 更新中间件
// @Security ApiKeyAuth
// @Tags Handler
// @Accept application/json
// @Produce  application/json
// @Success 200
// @Failure 400 {object} e.Response
// @Failure 500 {object} e.Response
// @Router /v1/handlers [post]
func Handlers(ctx *gin.Context) {
	var cfg dynamic.Config
	err := ctx.ShouldBindBodyWith(&cfg, binding.JSON)
	if err != nil {
		e.Errors(ctx, err)
		return
	}
	if err = middleware.Register(ctx.Request.Context(), api.Handlers(cfg)...); err != nil {
		e.Errors(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}
