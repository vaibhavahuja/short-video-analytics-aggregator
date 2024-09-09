package repository

import "github.com/scylladb/gocqlx/v3/table"

// CassandraConfig holds the configuration for the Cassandra connection
type CassandraConfig struct {
	Hosts       []string `json:"hosts"`
	Port        int      `json:"port"`
	Keyspace    string   `json:"keyspace"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Consistency string   `json:"consistency"`
	Timeout     int      `json:"timeout"`
}

var videoViewsMetaData = table.Metadata{
	Name:    "video_views",
	Columns: []string{"video_id", "timestamp_in_min", "aggregate_views"},
	PartKey: []string{"video_id"},
	SortKey: []string{"timestamp_in_min"},
}

// VideoViewsTable allows for simple CRUD operations based on videoViewsMetaData
var VideoViewsTable = table.New(videoViewsMetaData)

// VideoView represents a row in videoViewsTable.
type VideoView struct {
	VideoId        string
	TimestampInMin int
	AggregateViews int
}
