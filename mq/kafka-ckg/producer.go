package kafka_ckg

import (
	"fmt"
	ckg "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"time"
)

type Producer struct {
	producer *ckg.Producer
}

func (m *Producer) Init(option map[string]interface{}) (err error) {
	data := ckg.ConfigMap{}
	for k, v := range option {
		data[k] = v
	}
	m.producer, err = ckg.NewProducer(&data)
	if err != nil {
		log.Panicln(err)
	}
	//m.producer.Purge(ckg.PurgeNonBlocking)
	go func() {
		for e := range m.producer.Events() {
			switch ev := e.(type) {
			case *ckg.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
					return
				} else {
					fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()
	return err
}

func (m *Producer) Close() {
	m.producer.Flush(15 * 1000)
	time.Sleep(time.Second * 15)
	m.producer.Close()
}

func (m *Producer) Publish(topic string, data []byte) error {
	deliver := make(chan ckg.Event, 1)
	if err := m.producer.Produce(&ckg.Message{
		TopicPartition: ckg.TopicPartition{
			Topic:     &topic,
			Partition: ckg.PartitionAny,
		},
		Value: data,
	}, deliver); err != nil {
		return err
	}
	//e := <-deliver
	//log.Warnln(e.(*ckg.Message).TopicPartition.Partition)
	return nil
}
