package api

import (
	"encoding/json"
	"log"
	"strings"

	"monitoring-energy-service/internal/domain/ports/input"
	"monitoring-energy-service/internal/domain/ports/output"
)

type KafkaService struct {
	kafkaAdapter  output.KafkaAdapterInterface
	topicHandlers map[string]input.MessageHandler
	stopChan      chan struct{}
}

var _ input.KafkaServiceInterface = &KafkaService{}

func NewKafkaService(adapter output.KafkaAdapterInterface) *KafkaService {
	return &KafkaService{
		kafkaAdapter:  adapter,
		topicHandlers: make(map[string]input.MessageHandler),
		stopChan:      make(chan struct{}),
	}
}

func (ks *KafkaService) SendEvent(topic string, key string, event any) error {
	value, err := json.Marshal(event)
	if err != nil {
		return err
	}
	log.Printf("sending to kafka topic %s message %s", topic, value)
	return ks.kafkaAdapter.SendMessage(topic, key, value)
}

func (ks *KafkaService) RegisterHandler(topic string, handler input.MessageHandler) {
	ks.topicHandlers[topic] = handler
}

func (ks *KafkaService) ConsumeEvents() {
	log.Printf("Starting to consume events from Kafka")

	topics := make([]string, 0, len(ks.topicHandlers))
	for topic := range ks.topicHandlers {
		topics = append(topics, topic)
	}

	if len(topics) == 0 {
		log.Printf("No topics registered, skipping Kafka consumer")
		return
	}

	if err := ks.kafkaAdapter.SubscribeTopics(topics); err != nil {
		log.Fatalf("Error subscribing to topics: %s", err)
		return
	}

	disconnectedCount := 0
	connectionRefusedCount := 0

	for {
		select {
		case <-ks.stopChan:
			log.Println("Stopping Kafka event consumption.")
			return
		default:
			message, topic, err := ks.kafkaAdapter.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %s", err)

				if strings.Contains(err.Error(), "Disconnected") {
					disconnectedCount++
					if disconnectedCount > 10 {
						panic("disconnected from kafka too many times")
					}
				}

				if strings.Contains(err.Error(), "Connection refused") {
					connectionRefusedCount++
					if connectionRefusedCount > 10 {
						panic("connection refused from kafka too many times")
					}
				}

				continue
			}

			if handler, ok := ks.topicHandlers[topic]; ok {
				log.Printf("Handling message for topic %s", topic)
				err = handler.HandleMessage(message)
				if err != nil {
					log.Printf("Error handling message for topic %s: %s", topic, err)
				}
			} else {
				log.Printf("No handler registered for topic %s", topic)
			}
		}
	}
}

func (ks *KafkaService) StopConsuming() {
	close(ks.stopChan)
}
