package mosquitto

import (
	"github.com/xiaomLee/go-plugin/mq"
	"testing"
)

func TestProducer_Publish(t *testing.T) {
	p, err := NewProducer(
		mq.Auth("", ""),
		mq.Brokers("127.0.0.1:1883"),
		mq.ClientId("test.common.test.producer"),
		mq.Qos(0),
	)
	if err != nil {
		t.Fatal(err)
	}
	p.Publish("test.hello", []byte("hello world"))
}
