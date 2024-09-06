package queue

// Message represents a Kafka message.
type Message struct {
	Key       []byte
	Value     []byte
	Topic     string
	Partition int32
	Offset    int64
}
