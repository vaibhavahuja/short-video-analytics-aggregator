package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/spf13/viper"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/queue"
	"testing"
	"time"
)

type TestCase struct {
	Name   string
	Input  any
	Output any
	Err    error
}

func TestHandleMessage(t *testing.T) {
	mockProducer := &queue.MockProducer{}
	handler := NewEventHandler(mockProducer)
	testCases := []TestCase{
		{
			Name: "handle_event_success",
			Input: queue.Message{
				Value: []byte(fmt.Sprintf("{\"video_id\":\"asdsa\",\"video_title\":\"Amazing Adventure\",\"genres\":[\"Action\",\"Adventure\"],\"user_id\":12345,\"platform\":\"YouTube\",\"duration\":3600,\"timestamp\":\"%d\",\"video_quality\":\"4K\"}", time.Now().Unix())),
				Topic: "random",
			},
			Err: nil,
		},
		{
			Name: "handle_event_invalid_event_1",
			Input: queue.Message{
				Value: []byte("{\"video_id\":\"asasa\",\"video_title\":\"Amazing Adventure\",\"genres\":[\"Action\",\"Adventure\"],\"user_id\":12345,\"platform\":\"YouTube\",\"duration\":3600,\"timestamp\":\"1725720511\",\"video_quality\":\"4K\"}"),
				Topic: "random",
			},
			Err: errors.New("invalid event"),
		},
		{
			Name: "handle_event_invalid_event_2",
			Input: queue.Message{
				Value: []byte("{\"video_id\":\"\",\"video_title\":\"Amazing Adventure\",\"genres\":[\"Action\",\"Adventure\"],\"user_id\":12345,\"platform\":\"YouTube\",\"duration\":3600,\"timestamp\":\"1725720511\",\"video_quality\":\"4K\"}"),
				Topic: "random",
			},
			Err: errors.New("invalid event"),
		},
		{
			Name: "handle_event_invalid_event_3",
			Input: queue.Message{
				Value: []byte("{\"video_id\":\"asda\",\"video_title\":\"Amazing Adventure\",\"genres\":[\"Action\",\"Adventure\"],\"user_id\":0,\"platform\":\"YouTube\",\"duration\":3600,\"timestamp\":\"1725720511\",\"video_quality\":\"4K\"}"),
				Topic: "random",
			},
			Err: errors.New("invalid event"),
		},
		{
			Name: "handle_event_marshall_failure",
			Input: queue.Message{
				Value:     []byte("huasgd"),
				Topic:     "random",
				Partition: 0,
				Offset:    0,
			},
			Err: errors.New("invalid JSON"),
		},
		{
			Name: "handle_event_marshall_failure",
			Input: queue.Message{
				Value:     []byte(fmt.Sprintf("{\"video_id\":\"asdsa\",\"video_title\":\"Amazing Adventure\",\"genres\":[\"Action\",\"Adventure\"],\"user_id\":12345,\"platform\":\"YouTube\",\"duration\":3600,\"timestamp\":\"%d\",\"video_quality\":\"4K\"}", time.Now().Unix())),
				Topic:     "error_topic",
				Partition: 0,
				Offset:    0,
			},
			Err: errors.New("publish error"),
		},
	}
	for _, tc := range testCases {
		if msg, ok := tc.Input.(queue.Message); ok {
			viper.Set("message_queue.kafka.reduce-consumer.topic", msg.Topic)
			err := handler.HandleMessage(context.TODO(), &msg)
			assert.Equal(t, err, tc.Err)
		} else {
			t.Error("invalid Input test case format ", tc.Name)
		}
	}

}
