package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/app/models"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/queue"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/external/repository"
	"github.com/vaibhavahuja/short-video-analytics-aggregator/internal/utils"
	"time"
)

type ReducerHandler struct {
	repo repository.ShortVideoRepository
}

func NewReducerHandler(repo repository.ShortVideoRepository) *ReducerHandler {
	//also starts a separate goroutine which picks up the values from the map and inserts that into the cassandra database
	return &ReducerHandler{
		repo: repo,
	}
}

func (r *ReducerHandler) HandleMessage(ctx context.Context, msg *queue.Message) error {
	//handler only logs the message for now
	log.Ctx(ctx).Info().Str("key", string(msg.Key)).Int("partition", msg.Partition).Msg(string(msg.Value))
	//reducer flow
	videoId := string(msg.Key)
	if videoId == "" {
		return errors.New("empty video id. invalid event")
	}
	err := r.process(ctx, videoId, string(msg.Value))
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Error while processing message - reducer")
		return err
	}
	return nil
}

func (r *ReducerHandler) process(ctx context.Context, videoId string, timestamp string) error {
	//trim the timestamp to nearest minute and keep storing in a global hash map
	timeStampInMinutes, err := validateAndTrimTimeStampToMin(timestamp)
	if err != nil {
		return err
	}

	//aggregation
	models.AggregatorMap.AddView(videoId, timeStampInMinutes)
	log.Ctx(ctx).Debug().Str("video_id", videoId).Msg("Successfully added view in the aggregator map")
	return nil
}

// validateAndTrimTimeStampToMin Checks if timestamp received is within my range to process and if yes, trim to minutes and return
func validateAndTrimTimeStampToMin(timestamp string) (string, error) {
	t, err := utils.ConvertUnixTimestampToTime(timestamp)
	if err != nil {
		return "", err
	}
	//validate if time is within range
	minTime := time.Now().Add(-1 * time.Minute)
	if !t.After(minTime) {
		return "", errors.New("expired event time")
	}

	return fmt.Sprintf("%d", t.Truncate(time.Minute).Unix()), nil
}

// ProcessCurrentTimeData processes the current time data in the hash map and inserts into the database
func (r *ReducerHandler) ProcessCurrentTimeData(ctx context.Context) {
	ticker := time.Tick(1 * time.Minute)
	duration := viper.GetDuration("overrides.process-current-time")
	if duration > 0 {
		ticker = time.Tick(duration)
	}
	for {
		select {
		case <-ticker:
			prevMinuteTime := time.Now().Add(-1 * time.Minute).Truncate(time.Minute)
			prevMinuteTimeString := utils.ConvertTimeToStringUnix(prevMinuteTime)
			//get all entries for this time
			videoIds := models.AggregatorMap.GetVideoIdsByTimeStamp(prevMinuteTimeString)
			for _, videoId := range videoIds {
				// todo - do this via a worker pool, in order to control the number of goroutines spawned
				go func(id string) {
					views := models.AggregatorMap.GetViews(id, prevMinuteTimeString)
					err := r.repo.InsertAggregateByMinute(ctx, id, int(prevMinuteTime.Unix()), views)
					if err != nil {
						log.Ctx(ctx).Err(err).Msg("error while inserting aggregate by minute")
					}
				}(videoId)
			}
		case <-ctx.Done():
			//breaking the infinite loop if context is done
			return
		}
	}
}
