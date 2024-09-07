package kafka

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/queue"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/repository"
)

type Consumer struct {
	kafkaReader *kafka.Reader
	repository  repository.ShortVideoRepository
}

func NewConsumer(shortVideoRepo repository.ShortVideoRepository, cfg *queue.ConsumerConfig) *Consumer {
	readerConfig := kafka.ReaderConfig{
		Brokers: []string{cfg.BootstrapServers},
		GroupID: cfg.GroupID,
		Topic:   cfg.Topic,
	}

	if cfg.AutoOffsetReset == "earliest" {
		readerConfig.StartOffset = kafka.FirstOffset
	} else if cfg.AutoOffsetReset == "latest" {
		readerConfig.StartOffset = kafka.LastOffset
	}

	return &Consumer{
		repository:  shortVideoRepo,
		kafkaReader: kafka.NewReader(readerConfig),
	}
}

func (c *Consumer) Consume(ctx context.Context, handler queue.MessageHandler) {
	go func() {
		for {
			msg, err := c.kafkaReader.ReadMessage(ctx)
			if err != nil {
				log.Ctx(ctx).Err(err).Msg("error while consuming message")
			}
			message := queue.Message{
				Key:       msg.Key,
				Value:     msg.Value,
				Topic:     msg.Topic,
				Partition: msg.Partition,
				Offset:    msg.Offset,
			}
			if err = handler.HandleMessage(ctx, &message); err != nil {
				log.Ctx(ctx).Err(err).Msg("error while handling message")
			}
		}
	}()
}
