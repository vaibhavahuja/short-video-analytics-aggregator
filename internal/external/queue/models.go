package queue

// Message represents a Kafka message.
type Message struct {
	Key       []byte
	Value     []byte
	Topic     string
	Partition int
	Offset    int64
}

type ConsumerConfig struct {
	Topic                string `json:"topic"`
	BootstrapServers     string `json:"bootstrap_servers"`
	ClientID             string `json:"client_id"`
	GroupID              string `json:"group_id"`
	AutoOffsetReset      string `json:"auto_offset_reset"`
	EnableAutoCommit     bool   `json:"enable_auto_commit"`
	AutoCommitIntervalMs int    `json:"auto_commit_interval_ms"`
}

type ProducerConfig struct {
	Acks         string `json:"acks"`
	Compression  string `json:"compression"`
	LingerMs     int    `json:"linger_ms"`
	BatchSize    int    `json:"batch_size"`
	RetryBackoff int    `json:"retry_backoff"`
}
