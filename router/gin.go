// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/30

package router

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"proxy-go/api"
	"proxy-go/api/middleware"
)

// @title obs Swagger API
// @version 1.0
// @description This is a obs server.

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name ak

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name sk

// NewGinEngine gin router
func NewGinEngine() *gin.Engine {
	router := gin.New()
	router.Use(middleware.CrossDomain)

	if gin.Mode() != gin.ReleaseMode {
		// swagger
		url := ginSwagger.URL("/swagger/doc.json")
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

		// 增加性能测试
		pprof.Register(router)
	}

	routerV1 := router.Group("/api/v1")
	{
		routerV1.POST("/switch", api.UpdateSwitch)
		routerV1.PUT("/limit", api.UpdateRateLimit)
	}

	return router
}
