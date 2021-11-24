package mq

type Producer interface {
	Publish(topic string, data []byte) error
	Close()
}
