package router

import (
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/crochee/proxy-go/api"
	_ "github.com/crochee/proxy-go/docs"
)

// @title PROXY Swagger API
// @version 1.0
// @description This is a server API.

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-Auth-Token

// ApiHandler gin 路由处理
func ApiHandler() http.Handler {
	router := gin.New()

	if gin.Mode() == gin.DebugMode {
		url := ginSwagger.URL("/swagger/doc.json")
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

		pprof.Register(router)
	}

	router.GET("/metrics", func(ctx *gin.Context) {
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
	})

	v1router := router.Group("/v1")
	{
		v1router.GET("/nodes", api.GetBalanceNode)
	}
	return router
}
