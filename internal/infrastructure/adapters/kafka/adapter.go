package kafka

import (
	"monitoring-energy-service/internal/domain/ports/output"
	"monitoring-energy-service/internal/infrastructure/conf/kafkaconf"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaAdapter struct {
	producer *kafka.Producer
	consumer *kafka.Consumer
}

var _ output.KafkaAdapterInterface = &KafkaAdapter{}

func NewKafkaAdapter(factory *kafkaconf.KafkaFactory, groupID string) *KafkaAdapter {
	return &KafkaAdapter{
		producer: factory.NewProducer(),
		consumer: factory.NewConsumer(groupID),
	}
}

func (ka *KafkaAdapter) SubscribeTopics(topics []string) error {
	return ka.consumer.SubscribeTopics(topics, nil)
}

func (ka *KafkaAdapter) SendMessage(topic, key string, message []byte) error {
	ka.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          message,
	}, nil)

	ka.producer.Flush(15 * 1000)

	return nil
}

func (ka *KafkaAdapter) ReadMessage() (message []byte, topic string, err error) {
	msg, err := ka.consumer.ReadMessage(-1)
	if err != nil {
		return nil, "", err
	}

	topic = *msg.TopicPartition.Topic
	return msg.Value, topic, nil
}
