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
	input  any
	output any
	err    error
}

func TestHandleMessage(t *testing.T) {
	mockProducer := &queue.MockProducer{}
	handler := NewEventHandler(mockProducer)
	testCases := []TestCase{
		{
			Name: "handle_event_success",
			input: queue.Message{
				Value: []byte(fmt.Sprintf("{\"video_id\":\"asdsa\",\"video_title\":\"Amazing Adventure\",\"genres\":[\"Action\",\"Adventure\"],\"user_id\":12345,\"platform\":\"YouTube\",\"duration\":3600,\"timestamp\":\"%d\",\"video_quality\":\"4K\"}", time.Now().Unix())),
				Topic: "random",
			},
			err: nil,
		},
		{
			Name: "handle_event_invalid_event_1",
			input: queue.Message{
				Value: []byte("{\"video_id\":\"asasa\",\"video_title\":\"Amazing Adventure\",\"genres\":[\"Action\",\"Adventure\"],\"user_id\":12345,\"platform\":\"YouTube\",\"duration\":3600,\"timestamp\":\"1725720511\",\"video_quality\":\"4K\"}"),
				Topic: "random",
			},
			err: errors.New("invalid event"),
		},
		{
			Name: "handle_event_invalid_event_2",
			input: queue.Message{
				Value: []byte("{\"video_id\":\"\",\"video_title\":\"Amazing Adventure\",\"genres\":[\"Action\",\"Adventure\"],\"user_id\":12345,\"platform\":\"YouTube\",\"duration\":3600,\"timestamp\":\"1725720511\",\"video_quality\":\"4K\"}"),
				Topic: "random",
			},
			err: errors.New("invalid event"),
		},
		{
			Name: "handle_event_invalid_event_3",
			input: queue.Message{
				Value: []byte("{\"video_id\":\"asda\",\"video_title\":\"Amazing Adventure\",\"genres\":[\"Action\",\"Adventure\"],\"user_id\":0,\"platform\":\"YouTube\",\"duration\":3600,\"timestamp\":\"1725720511\",\"video_quality\":\"4K\"}"),
				Topic: "random",
			},
			err: errors.New("invalid event"),
		},
		{
			Name: "handle_event_marshall_failure",
			input: queue.Message{
				Value:     []byte("huasgd"),
				Topic:     "random",
				Partition: 0,
				Offset:    0,
			},
			err: errors.New("invalid JSON"),
		},
		{
			Name: "handle_event_marshall_failure",
			input: queue.Message{
				Value:     []byte(fmt.Sprintf("{\"video_id\":\"asdsa\",\"video_title\":\"Amazing Adventure\",\"genres\":[\"Action\",\"Adventure\"],\"user_id\":12345,\"platform\":\"YouTube\",\"duration\":3600,\"timestamp\":\"%d\",\"video_quality\":\"4K\"}", time.Now().Unix())),
				Topic:     "error_topic",
				Partition: 0,
				Offset:    0,
			},
			err: errors.New("publish error"),
		},
	}
	for _, tc := range testCases {
		if msg, ok := tc.input.(queue.Message); ok {
			viper.Set("message_queue.kafka.reduce-consumer.topic", msg.Topic)
			err := handler.HandleMessage(context.TODO(), &msg)
			assert.Equal(t, err, tc.err)
		} else {
			t.Error("invalid input test case format")
		}
	}

}
