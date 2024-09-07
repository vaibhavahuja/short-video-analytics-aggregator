package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/app/models"
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

	var event models.ShortVideoAnalyticsEvent
	//val := "{\"video_id\":98765,\"video_title\":\"Amazing Adventure\",\"genres\":[\"Action\",\"Adventure\"],\"user_id\":12345,\"platform\":\"YouTube\",\"duration\":3600,\"timestamp\":\"2024-09-07T15:30:00Z\",\"video_quality\":\"4K\"}"
	//log.Ctx(ctx).Debug().Str("value", string(msg.Value)).Msg("received message")
	if err := json.Unmarshal([]byte(msg.Value), &event); err != nil {
		log.Ctx(ctx).Err(err).Msg("error while unmarshalling event")
		return err
	}
	if !event.IsValid() {
		log.Ctx(ctx).Error().Msg("received invalid event")
		return errors.New("invalid event")
	}

	log.Ctx(ctx).Debug().Any("event", event).Msg("successfully unmarshalled the event")
	//get partitionId
	partition := getPartitionByVideoId(event.VideoId)
	// To keep things simple, the system only cares about the video_id and when was it viewed - the timestamp
	topic := viper.GetString("message_queue.kafka.reduce-consumer.topic")
	//will have a partition selection logic such that all messages for video_id -> x go to the same partition always
	if err := e.producer.Publish(topic, partition, []byte(event.VideoId), []byte(event.Timestamp)); err != nil {
		return err
	}
	return nil
}

// getPartitionByVideoId Calculates the partition to which we need to send videoId using hashing techniques
func getPartitionByVideoId(videoId string) int {
	numPartitions := viper.GetInt("message_queue.kafka.number_of_partitions")
	if numPartitions == 0 {
		//default number of partitions is one
		numPartitions = 1
	}
	videoIdBytes := []byte(videoId)
	hash := sha256.Sum256(videoIdBytes)
	hashInt := binary.BigEndian.Uint32(hash[:4])
	partition := int(hashInt) % numPartitions
	return partition
}
