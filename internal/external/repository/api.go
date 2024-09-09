package repository

import "context"

// ShortVideoRepository Defines the methods for interacting with short video aggregated views data
type ShortVideoRepository interface {
	InsertAggregateByMinute(ctx context.Context, videoId string, timeStamp, views int) error
	GetViewerCountByVideoIDAndTimeRange(ctx context.Context, videoId string, timeStamp int) (int, error)
}
