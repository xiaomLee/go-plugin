package mosquitto

import (
	"github.com/xiaomLee/go-plugin/mq"
	"time"
)

type Producer struct {
	client Client
}

const defaultTimeout = 3

func NewProducer(opts ...mq.Option) (*Producer, error) {
	p := Producer{}
	if err := p.client.Init(opts...); err != nil {
		return nil, err
	}

	return &p, nil
}

func (p *Producer) Publish(topic string, data []byte) error {
	token := p.client.Publish(topic, byte(p.client.options.Qos), p.client.options.Retained, data)
	timeout := defaultTimeout
	if p.client.options.PublishTimeout > 0 {
		timeout = p.client.options.PublishTimeout
	}
	token.WaitTimeout(time.Second * time.Duration(timeout))
	return token.Error()
}

func (p *Producer) Close() {
	p.client.Disconnect(3000)
}
