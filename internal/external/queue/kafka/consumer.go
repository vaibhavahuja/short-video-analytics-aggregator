package kafka

import (
	"context"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/queue"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/repository"
)

type Consumer struct {
	//add kafka config
	repository repository.ShortVideoRepository
}

func NewConsumer(shortVideoRepo repository.ShortVideoRepository) *Consumer {
	return &Consumer{
		repository: shortVideoRepo,
	}
}

func (c *Consumer) Subscribe(topics []string) error {
	return nil
}

func (c *Consumer) Consume(ctx context.Context, handler queue.MessageHandler) error {
	return nil
}
