package cassandra

import (
	"context"
	"github.com/gocql/gocql"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx/v3"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/repository"
)

type CassandraRepository struct {
	session *gocqlx.Session
}

func NewCassandraRepository(cfg *repository.CassandraConfig) (*CassandraRepository, error) {
	// Create gocql cluster.
	cluster := gocql.NewCluster(cfg.Hosts...)
	cluster.Keyspace = cfg.Keyspace
	var consistency gocql.Consistency
	err := consistency.UnmarshalText([]byte(cfg.Consistency))
	if err != nil {
		consistency = gocql.Any
	}
	cluster.Consistency = consistency
	cluster.Port = cfg.Port
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.Username,
		Password: cfg.Password,
	}

	cqlSession, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		return nil, err
	}
	log.Debug().Msg("Successfully initialised cql session")

	return &CassandraRepository{
		session: &cqlSession,
	}, nil
}

func (cr *CassandraRepository) InsertAggregateByMinute(ctx context.Context, videoId string, timeStamp, views int) error {
	log.Ctx(ctx).Debug().Str("videoId", videoId).
		Int("timestamp", timeStamp).
		Int("views", views).
		Msg("inserting aggregate by minute in db")

	videoView := repository.VideoView{
		VideoId:        videoId,
		TimestampInMin: timeStamp,
		AggregateViews: views,
	}

	q := cr.session.Query(repository.VideoViewsTable.Insert()).BindStruct(videoView)
	if err := q.ExecRelease(); err != nil {
		log.Ctx(ctx).Err(err).Msg("error while inserting videoView to table")
		return err
	}
	log.Ctx(ctx).Debug().Msg("successfully inserted video to video views table")
	return nil
}

func (cr *CassandraRepository) GetViewerCountByVideoIDAndTimeRange(ctx context.Context, videoId string, timeStamp int) (int, error) {
	//to implement later
	return 0, nil
}
