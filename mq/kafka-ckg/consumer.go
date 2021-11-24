package kafka_ckg

import (
	"fmt"
	"github.com/xiaomLee/go-plugin/mq"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"sync"
	"time"
)

// Consumer TODO 不支持增量sub, 每次调用sub会将当前的订阅替换
type Consumer struct {
	client    *kafka.Consumer
	options   mq.Options
	callbacks map[string][]mq.MessageHandler
	mu        sync.RWMutex
	stop      chan bool
	exited    chan bool
}

func (c *Consumer) Init(option map[string]interface{}) {

}

func (c *Consumer) loop() {
	for {
		select {
		case <-c.stop:
			c.exited <- true
			return
		default:
			msg, err := c.client.ReadMessage(time.Second * 5)
			if err == nil {
				c.onMessageReceived(msg)
			} else {
				if err.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue
				}
				time.Sleep(time.Second)
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}
}

func (c *Consumer) Close() {
	c.client.Close()
	c.stop <- true
	<-c.exited
}
func (c *Consumer) Subscribe(topic []string, callback mq.MessageHandler, cover bool) error {
	c.mu.Lock()
	for _, v := range topic {
		if _, ok := c.callbacks[v]; !ok || cover {
			c.callbacks[v] = make([]mq.MessageHandler, 0)
		}
		c.callbacks[v] = append(c.callbacks[v], callback)
	}
	c.mu.Unlock()
	if err := c.client.SubscribeTopics(topic, func(consumer *kafka.Consumer, event kafka.Event) error {
		log.Println(event)
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (c *Consumer) Unsubcribe() error {
	c.client.Unsubscribe()
	return nil
}

func (c *Consumer) onMessageReceived(message *kafka.Message) {

}
