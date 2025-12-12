package kafkaconf

import (
	"log"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func deliveryReport(deliveryChan chan kafka.Event) {
	for e := range deliveryChan {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
			} else {
				log.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
		}
	}
}

type KafkaFactory struct {
	brokerList string
	autoOffset string
}

func NewKafkaFactory(brokers []string, autoOffset string) *KafkaFactory {
	brokerList := strings.Join(brokers, ",")

	return &KafkaFactory{
		brokerList: brokerList,
		autoOffset: autoOffset,
	}
}

func (kf *KafkaFactory) NewProducer() *kafka.Producer {
	deliveryChan := make(chan kafka.Event)

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":   kf.brokerList,
		"retries":             1,
		"socket.timeout.ms":   5000,
		"go.delivery.reports": false,
	})

	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}

	go deliveryReport(deliveryChan)

	return p
}

func (kf *KafkaFactory) NewConsumer(groupID string) *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":       kf.brokerList,
		"group.id":                groupID,
		"auto.offset.reset":       kf.autoOffset,
		"session.timeout.ms":      10000,
		"heartbeat.interval.ms":   3000,
		"enable.auto.commit":      true,
		"auto.commit.interval.ms": 5000,
		"max.poll.interval.ms":    300000,
	})
	if err != nil {
		panic(err)
	}
	return c
}
