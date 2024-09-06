package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetHttpServer() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())

	engine.GET("/health", func(ctx *gin.Context) {
		ctx.AbortWithStatus(http.StatusOK)
	})
	return engine
}
