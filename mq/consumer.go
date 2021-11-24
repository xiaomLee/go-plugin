package mq

type Consumer interface {
	Close() error
	Subscribe(topics []string, callback MessageHandler) error
}

type MessageHandler func(topic string, message []byte) bool

type Option func(*Options)
