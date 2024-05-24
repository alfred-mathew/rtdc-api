package main

import (
	"github.com/gin-gonic/gin"
	"signzy.com/rtdc-api/health"
)

func createServer() *gin.Engine {
	router := gin.Default()

	router.SetTrustedProxies(nil)

	router.GET("/", health.Check)

	return router
}
