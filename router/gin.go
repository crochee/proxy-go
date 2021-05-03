// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/30

package router

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/crochee/proxy-go/api"
)

// @title obs Swagger API
// @version 1.0
// @description This is a obs server.

// NewGinEngine gin router
func NewGinEngine() *gin.Engine {
	router := gin.New()

	if gin.Mode() == gin.DebugMode {
		// swagger
		url := ginSwagger.URL("/swagger/doc.json")
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

		// 增加性能测试
		pprof.Register(router)
	}

	prof := router.Group("/debug/pprof")
	{
		prof.GET("/index", api.Index)
		prof.GET("/profile", api.Profile)
		prof.GET("/trace", api.Trace)
		prof.GET("/heap", api.Heap)
	}

	routerV1 := router.Group("/api/v1")

	mid := routerV1.Group("/mid")
	{
		mid.POST("/switch", api.UpdateSwitch)
		mid.GET("/switch", api.ListSwitch)
		mid.PUT("/limit", api.UpdateRateLimit)
		mid.GET("/limit", api.GetRateLimit)
	}

	return router
}
