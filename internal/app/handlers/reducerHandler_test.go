package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/spf13/viper"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/app/models"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/queue"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/repository"
	"testing"
	"time"
)

func TestReducerHandler_HandleMessage(t *testing.T) {
	//initialising dependencies
	mockRepo := &repository.MockRepository{}
	reducerHandler := NewReducerHandler(mockRepo)
	models.InitAggregatorMap()

	testCases := []TestCase{
		{
			Name: "handle_event_success",
			Input: queue.Message{
				Key:   []byte("video_id_1"),
				Value: []byte(fmt.Sprintf("%d", time.Now().Unix())),
				Topic: "random",
			},
			Err: nil,
		},
		{
			Name: "handle_event_error_invalid_timestamp",
			Input: queue.Message{
				Key:   []byte("video_id_1"),
				Value: []byte(fmt.Sprintf("%d", time.Now().Add(-2*time.Hour).Unix())),
				Topic: "random",
			},
			Err: errors.New("expired event time"),
		},
		{
			Name: "handle_event_error_invalid_timestamp",
			Input: queue.Message{
				Key:   []byte(""),
				Value: []byte(fmt.Sprintf("%d", time.Now().Add(-2*time.Hour).Unix())),
				Topic: "random",
			},
			Err: errors.New("empty video id. invalid event"),
		},
	}

	for _, tc := range testCases {
		if msg, ok := tc.Input.(queue.Message); ok {
			viper.Set("message_queue.kafka.reduce-consumer.topic", msg.Topic)
			err := reducerHandler.HandleMessage(context.TODO(), &msg)
			assert.Equal(t, err, tc.Err)
		} else {
			t.Error("invalid Input test case format ", tc.Name)
		}
	}
}

func TestReducerHandler_ProcessCurrentTimeData(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	reducerHandler := NewReducerHandler(mockRepo)
	models.InitAggregatorMap()
	currTime := time.Now().Add(-1 * time.Minute).Truncate(time.Minute).Unix()
	videoIds := []string{"video_id_1", "video_id_2"}
	for _, id := range videoIds {
		models.AggregatorMap.AddView(id, fmt.Sprintf("%d", currTime))
	}

	viper.Set("overrides.process-current-time", 50*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()
	reducerHandler.ProcessCurrentTimeData(ctx)
}
