package main

import "github.com/gin-gonic/gin"

func routes(e *gin.Engine) {
	v1 := e.Group("/v1")

	v1.POST("/orders")
}
