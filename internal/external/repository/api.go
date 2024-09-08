package repository

import "context"

// ShortVideoRepository Defines the methods for interacting with short video aggregated views data
type ShortVideoRepository interface {
	InsertAggregateByMinute(ctx context.Context, videoId, timeStamp string, views int) error
	GetTopNPopularVideo(ctx context.Context) ([]int, error)
	GetViewerCountByVideoIDAndTimeRange(ctx context.Context) (int, error)
}
