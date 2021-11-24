package mosquitto

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"github.com/xiaomLee/go-plugin/mq"
	"sync"
	"time"
)

type Consumer struct {
	client    Client
	topics    map[string]byte
	callbacks map[string][]mq.MessageHandler
	mu        sync.RWMutex
}

func NewConsumer(opts ...mq.Option) (*Consumer, error) {
	c := Consumer{}
	opts = append(opts, mq.OnConnect(c.OnConnect))
	if err := c.client.Init(opts...); err != nil {
		return nil, err
	}

	if c.callbacks == nil {
		c.callbacks = make(map[string][]mq.MessageHandler)
	}
	return &c, nil
}

func (c *Consumer) Close() error {
	c.client.Disconnect(3000)
	return nil
}

func (c *Consumer) Subscribe(topics []string, callback mq.MessageHandler) error {
	dd := make(map[string]byte)
	c.mu.Lock()
	for _, v := range topics {
		if c.topics == nil {
			c.topics = make(map[string]byte)
		}
		if _, ok := c.topics[v]; !ok {
			c.topics[v] = byte(1)
		}
		dd[v] = byte(1)

		if _, ok := c.callbacks[v]; !ok {
			c.callbacks[v] = make([]mq.MessageHandler, 0)
		}
		c.callbacks[v] = append(c.callbacks[v], callback)
	}
	c.mu.Unlock()

	if token := c.client.SubscribeMultiple(dd, c.onMessageReceived); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (c *Consumer) UnSubscribe(topics []string) error {
	if token := c.client.Unsubscribe(topics...); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	c.mu.Lock()
	for _, topic := range topics {
		delete(c.callbacks, topic)
	}
	c.mu.Unlock()
	return nil
}

func (c *Consumer) OnConnect(client mqtt.Client) {
	c.mu.RLock()
	topics := c.topics
	c.mu.RUnlock()

	if len(topics) == 0 {
		return
	}
	logrus.Infoln("onConnect, Subscribe topics:", topics)
	if token := client.SubscribeMultiple(topics, c.onMessageReceived); token.WaitTimeout(time.Second*3) && token.Error() != nil {
		logrus.Error(token.Error())
	}
	logrus.Infoln("Subscribe success")
}

func (c *Consumer) onMessageReceived(client mqtt.Client, message mqtt.Message) {
	topic := message.Topic()
	logrus.WithFields(logrus.Fields{"Topic": topic, "message": string(message.Payload())}).Debug("处理消息")
	c.mu.RLock()
	callbacks, ok := c.callbacks[message.Topic()]
	c.mu.RUnlock()
	if !ok {
		return
	}
	for _, fn := range callbacks {
		fn(topic, message.Payload())
	}
}
