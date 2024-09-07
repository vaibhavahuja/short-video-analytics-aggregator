package queue

import "context"

type ConsumerInterface interface {
	Consume(ctx context.Context, handler MessageHandler)
}

// MessageHandler defines the behavior for processing a Kafka message.
type MessageHandler interface {
	HandleMessage(ctx context.Context, msg *Message) error
}

type ProducerInterface interface {
	Publish(topic string, partition int, msg []byte) error
}
