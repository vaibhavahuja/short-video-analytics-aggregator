package queue

import "errors"

type MockProducer struct {
}

func (mp *MockProducer) Publish(topic string, partition int, key, msg []byte) error {
	if topic == "error_topic" {
		return errors.New("publish error")
	}
	return nil
}
