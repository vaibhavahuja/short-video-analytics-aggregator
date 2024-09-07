package handlers

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/queue"
)

type EventHandler struct {
	producer queue.ProducerInterface
}

func NewEventHandler(producer queue.ProducerInterface) *EventHandler {
	return &EventHandler{
		producer: producer,
	}
}

func (e *EventHandler) HandleMessage(ctx context.Context, msg *queue.Message) error {
	//handler only logs the message for now
	log.Ctx(ctx).Info().Msg(string(msg.Value))
	//now it pushes the message to some other topic
	topic := viper.GetString("message_queue.kafka.reduce-consumer.topic")
	//will have a partition selection logic such that all messages for video_id -> x go to the same partition always
	if err := e.producer.Publish(topic, 0, msg.Value); err != nil {
		return err
	}
	return nil
}
