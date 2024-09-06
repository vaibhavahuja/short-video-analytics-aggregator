package queue

import "context"

type ConsumerInterface interface {
	Subscribe(topics []string) error
	Consume(ctx context.Context, handler MessageHandler) error
}

// MessageHandler defines the behavior for processing a Kafka message.
type MessageHandler interface {
	HandleMessage(msg *Message) error
}

type ProducerInterface interface {
	Publish(topic, partition string, msg []byte) error
}
