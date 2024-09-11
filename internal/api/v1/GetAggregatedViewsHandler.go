package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/repository"
	"net/http"
	"strconv"
)

type VideoAggregatorHandler struct {
	repo repository.ShortVideoRepository
}

func NewVideoAggregatorHandler(repo repository.ShortVideoRepository) *VideoAggregatorHandler {
	return &VideoAggregatorHandler{repo: repo}
}

func (vh *VideoAggregatorHandler) GetAggregatedViewsHandler(ctx *gin.Context) {
	//handles the views
	videoId := ctx.Query("video_id")
	if videoId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "empty video id query param received",
		})
		return
	}
	timeStampInMinStr := ctx.Query("timestamp_in_min")
	timeStampInt, err := strconv.Atoi(timeStampInMinStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid timestamp_in_min received",
		})
		return
	}
	val, err := vh.repo.GetViewerCountByVideoIDAndTimeRange(ctx, videoId, timeStampInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"viewer_count": val,
	})
}
