package router

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title PROXY Swagger API
// @version 1.0
// @description This is a server API.

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

	return router
}
