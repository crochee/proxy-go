// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/30

package router

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

// GinRun gin router
func GinRun() *gin.Engine {
	router := gin.New()

	if gin.Mode() != gin.ReleaseMode {
		// swagger
		url := ginSwagger.URL("/swagger/doc.json")
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

		// 增加性能测试
		pprof.Register(router)
	}
	return router
}
