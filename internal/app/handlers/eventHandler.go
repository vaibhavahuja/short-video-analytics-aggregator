package handlers

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/queue"
)

type EventHandler struct {
	//add your event handler here - producer details maybe
}

func NewEventHandler() *EventHandler {
	return &EventHandler{}
}

func (e *EventHandler) HandleMessage(ctx context.Context, msg *queue.Message) error {
	//handler only logs the message for now
	log.Ctx(ctx).Info().Msg(string(msg.Value))
	return nil
}
