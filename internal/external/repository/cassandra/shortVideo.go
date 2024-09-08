package cassandra

import (
	"context"
)

type CassandraRepository struct {
	//todo add cassandra client and other things here
}

func NewCassandraRepository() *CassandraRepository {
	return &CassandraRepository{}
}

func (cr *CassandraRepository) InsertAggregateByMinute(ctx context.Context, videoId, timeStamp string, views int) error {
	return nil
}

func (cr *CassandraRepository) GetTopNPopularVideo(ctx context.Context) ([]int, error) {
	return []int{}, nil
}

func (cr *CassandraRepository) GetViewerCountByVideoIDAndTimeRange(ctx context.Context) (int, error) {
	return 0, nil
}
