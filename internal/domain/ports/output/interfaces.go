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

// EventRepositoryInterface define el contrato para la persistencia de eventos
//
// CAMBIO REALIZADO: Interface agregada (líneas 30-35)
// RAZÓN: Sigue el patrón de arquitectura hexagonal, definiendo el contrato que
// debe implementar el repositorio de eventos (EventRepository)
//
// MÉTODOS:
// - Create: Guarda un evento consumido desde Kafka
// - FindAll: Lista todos los eventos (para API REST)
// - FindByID: Obtiene un evento específico
// - FindByEventType: Filtra eventos por tipo (power_reading, alert, etc.)
type EventRepositoryInterface interface {
	Create(entity *entities.EventEntity) (*entities.EventEntity, error)
	FindAll() ([]*entities.EventEntity, error)
	FindByID(id uuid.UUID) (*entities.EventEntity, error)
	FindByEventType(eventType string) ([]*entities.EventEntity, error)
}
