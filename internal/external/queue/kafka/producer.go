package kafka

type Producer struct {
}

func NewProducer() *Producer {
	return &Producer{}
}

func (p *Producer) Publish(topic, partition string, msg []byte) error {
	//to add implementation
	return nil
}
