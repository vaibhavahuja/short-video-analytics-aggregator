package models

import (
	"github.com/rs/zerolog/log"
	"strconv"
	"sync"
	"time"
)

var (
	AggregatorMap ViewAggregator
	once          sync.Once
)

type ShortVideoAnalyticsEvent struct {
	VideoId      string   `json:"video_id"`
	VideoTitle   string   `json:"video_title"`
	Genres       []string `json:"genres"`
	UserId       int      `json:"user_id"`
	Platform     string   `json:"platform"`
	Duration     int      `json:"duration"`
	Timestamp    string   `json:"timestamp"`
	VideoQuality string   `json:"video_quality"`
}

func (event *ShortVideoAnalyticsEvent) IsValid() bool {
	if event.UserId == 0 {
		log.Error().Msg("invalid event. empty user id.")
		return false
	}

	if event.VideoId == "" {
		log.Error().Msg("invalid event. received empty video_id")
		return false
	}

	//checks if event was sent within +- 15 seconds of actual range
	if !isValidTimestamp(event.Timestamp) {
		log.Error().Msg("invalid event. received invalid timestamp")
		return false
	}

	return true
}

func isValidTimestamp(timestamp string) bool {
	if timestamp == "" {
		return false
	}
	currTime := time.Now()
	minTime := currTime.Add(-30 * time.Second)
	maxTime := currTime.Add(30 * time.Second)

	// Convert Unix timestamp to time.Time
	timeStampInt, err := strconv.Atoi(timestamp)
	if err != nil {
		return false
	}
	t := time.Unix(int64(timeStampInt), 0)
	// checking if timestamp falls in given range
	return t.After(minTime) && t.Before(maxTime)
}

type timeStamp string
type videoId string

type ViewAggregator struct {
	aggregator map[timeStamp]map[videoId]int
	mu         sync.Mutex
}

// InitAggregatorMap Initialize the singleton instance for the aggregatorMap
func InitAggregatorMap() {
	once.Do(func() {

		AggregatorMap = ViewAggregator{aggregator: make(map[timeStamp]map[videoId]int)}
	})

}

// AddView Thread-safe method to update the aggregator map
func (va *ViewAggregator) AddView(id string, ts string) {
	va.mu.Lock()
	defer va.mu.Unlock()

	if va.aggregator[timeStamp(ts)] == nil {
		va.aggregator[timeStamp(ts)] = make(map[videoId]int)
	}

	va.aggregator[timeStamp(ts)][videoId(id)]++
}

// GetViews Thread-safe method to get the view count for a specific id
func (va *ViewAggregator) GetViews(id string, ts string) int {
	va.mu.Lock()
	defer va.mu.Unlock()
	if va.aggregator[timeStamp(ts)] == nil {
		return 0
	}
	return va.aggregator[timeStamp(ts)][videoId(id)]
}

func (va *ViewAggregator) GetVideoIdsByTimeStamp(ts string) []string {
	va.mu.Lock()
	defer va.mu.Unlock()
	subMap := va.aggregator[timeStamp(ts)]
	var resp []string
	for id, _ := range subMap {
		resp = append(resp, string(id))
	}
	return resp
}
