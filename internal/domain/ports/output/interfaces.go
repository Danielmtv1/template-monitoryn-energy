package output

import (
	"monitoring-energy-service/internal/domain/entities"

	"github.com/google/uuid"
)

// KafkaAdapterInterface defines the contract for Kafka adapter operations
type KafkaAdapterInterface interface {
	SendMessage(topic, key string, message []byte) error
	ReadMessage() (message []byte, topic string, err error)
	SubscribeTopics(topics []string) error
}

// WebhookAdapterInterface defines the contract for webhook operations
type WebhookAdapterInterface interface {
	SendPayload(url string, payload any) error
}

// ExampleRepositoryInterface defines the contract for example data persistence
type ExampleRepositoryInterface interface {
	FindByID(id uuid.UUID) (*entities.ExampleEntity, error)
	FindAll() ([]*entities.ExampleEntity, error)
	Create(entity *entities.ExampleEntity) (*entities.ExampleEntity, error)
	Update(entity *entities.ExampleEntity) (*entities.ExampleEntity, error)
	Delete(id uuid.UUID) error
}
