package repository

import (
	"context"
	"errors"
)

// MockRepository Using this for unit tests
type MockRepository struct{}

func (mr *MockRepository) InsertAggregateByMinute(ctx context.Context, videoId string, timeStamp, views int) error {
	if videoId == "video_id_1" {
		return errors.New("qa error")
	}
	return nil
}

func (mr *MockRepository) GetViewerCountByVideoIDAndTimeRange(ctx context.Context, videoId string, timeStamp int) (int, error) {
	return 7, nil
}
