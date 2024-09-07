package handlers

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/queue"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/repository"
)

type ReducerHandler struct {
	repo repository.ShortVideoRepository
}

func NewReducerHandler(repo repository.ShortVideoRepository) *ReducerHandler {
	return &ReducerHandler{
		repo: repo,
	}
}

func (r *ReducerHandler) HandleMessage(ctx context.Context, msg *queue.Message) error {
	//handler only logs the message for now
	log.Ctx(ctx).Info().Str("key", string(msg.Key)).Int("partition", msg.Partition).Msg(string(msg.Value))
	//reducer flow
	// can have multiple goroutines or maps -> where I count the number of integers and push it later?

	return nil
}
