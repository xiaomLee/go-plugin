package mosquitto

import (
	"fmt"
	"github.com/xiaomLee/go-plugin/mq"
	"strconv"
	"testing"
	"time"
)

func TestConsumer_Subscribe(t *testing.T) {
	consumer, err := NewConsumer(mq.Auth("", ""), mq.Brokers("mydev:1883"), mq.Debug(true))
	if err != nil {
		t.Fatal(err)
	}

	// start sub 1
	go func() {
		consumer.Subscribe([]string{"topic-test-1"}, func(topic string, message []byte) bool {
			fmt.Println("this is first callback function")
			fmt.Println(topic, string(message))
			return true
		})
	}()

	// start sub 2
	go func() {
		time.Sleep(20 * time.Second)
		consumer.Subscribe([]string{"topic-test-1"}, func(topic string, message []byte) bool {
			fmt.Println("this is second callback function")
			fmt.Println(topic, string(message))
			return true
		})
	}()

	// wait 2 min
	select {
	case <-time.After(2 * time.Minute):

	}
}

func TestConsumer_Reconnect(t *testing.T) {
	consumer, err := NewConsumer(mq.Auth("", ""), mq.Brokers("127.0.0.1:1883"),
		mq.CleanSession(false),
		mq.ClientId("cccc"),
		mq.ResumeSubs(true),
		mq.Qos(1),
	)
	if err != nil {
		t.Fatal(err)
	}

	producer, err := NewProducer(mq.Auth("", ""), mq.Brokers("127.0.0.1:1883"),
		mq.CleanSession(false),
		mq.ClientId("pppp"),
		mq.Qos(1),
	)
	if err != nil {
		t.Fatal(err)
	}

	// start consumer
	if err := consumer.Subscribe([]string{"topic.test.reconnect"}, func(topic string, message []byte) bool {
		fmt.Println("consumer1", topic, string(message))
		return true
	}); err != nil {
		t.Fatal(err)
	}

	// consumer disconnect
	go func() {
		time.Sleep(time.Second * 10)
		consumer0, err := NewConsumer(mq.Auth("", ""), mq.Brokers("127.0.0.1:1883"),
			mq.CleanSession(false),
			mq.ClientId("cccc"),
			mq.ResumeSubs(true),
			mq.Qos(1),
		)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second * 5)
		consumer0.Close()
	}()

	// start producer
	for i := 0; i < 100; i++ {
		if err := producer.Publish("topic.test.reconnect", []byte(strconv.Itoa(i))); err != nil {
			t.Log(err)
		}
		time.Sleep(time.Second)
	}
}
