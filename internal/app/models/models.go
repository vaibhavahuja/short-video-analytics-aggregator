package models

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
	//val := "{\n    \"video_id\": 1,\n    \"video_title\": \"Amazing Adventure\",\n    \"genres\": [\"Action\", \"Adventure\"],\n    \"user_id\": 12345,\n    \"platform\": \"YouTube\",\n    \"duration\": 3600,\n    \"timestamp\": \"2024-09-07T15:30:00Z\",\n    \"video_quality\": \"4K\"\n}\n"
	//checks if event is valid or not
	//val := "{\"video_id\":98765,\"video_title\":\"Amazing Adventure\",\"genres\":[\"Action\",\"Adventure\"],\"user_id\":12345,\"platform\":\"YouTube\",\"duration\":3600,\"timestamp\":\"2024-09-07T15:30:00Z\",\"video_quality\":\"4K\"}"
	return true
}
