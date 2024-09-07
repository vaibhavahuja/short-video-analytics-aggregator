package kafka

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/queue"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(producerConfig *queue.ProducerConfig) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(producerConfig.Brokers),
			RequiredAcks: kafka.RequireNone,
			Balancer:     &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) Publish(topic string, partition int, msg []byte) error {
	message := kafka.Message{
		Topic:     topic,
		Partition: partition,
		Value:     msg,
	}
	err := p.writer.WriteMessages(context.TODO(), message)
	if err != nil {
		return err
	}
	log.Debug().Str("topic", topic).
		Int("partition", partition).
		Msg("Successfully published message")
	return nil
}
