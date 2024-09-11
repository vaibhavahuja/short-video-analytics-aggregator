package http

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/vaibhavahuja/short-video-analytics-aggregator/internal/api/v1"
	"net/http"
)

func GetHttpServer(handler v1.VideoAggregatorHandler) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())

	engine.GET("/health", func(ctx *gin.Context) {
		ctx.AbortWithStatus(http.StatusOK)
	})

	engine.GET("/api/v1/viewer-count", handler.GetAggregatedViewsHandler)
	return engine
}
