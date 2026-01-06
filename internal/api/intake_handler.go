package api

import (
	"encoding/json"
	"log"

	"monitoring-energy-service/internal/domain/entities"
	"monitoring-energy-service/internal/domain/ports/input"
	"monitoring-energy-service/internal/domain/ports/output"
)

// IntakeHandler procesa mensajes consumidos desde Kafka
//
// PROPÓSITO:
// Actúa como consumer de Kafka para el topic "intake", recibiendo eventos y
// guardándolos en PostgreSQL para análisis posterior.
//
// CAMBIO REALIZADO: Se agregó EventRepository como dependencia
// RAZÓN: Necesitábamos persistir los eventos consumidos en la base de datos
type IntakeHandler struct {
	eventRepository output.EventRepositoryInterface // CAMBIO: Agregado para guardar eventos en DB
}

var _ input.MessageHandler = &IntakeHandler{}

// NewIntakeHandler crea una nueva instancia del handler de Kafka
// CAMBIO: Ahora recibe eventRepository como parámetro
// RAZÓN: Necesita acceso a la DB para guardar los eventos consumidos
func NewIntakeHandler(eventRepository output.EventRepositoryInterface) *IntakeHandler {
	return &IntakeHandler{
		eventRepository: eventRepository,
	}
}

// HandleMessage procesa cada mensaje recibido desde Kafka
//
// FLUJO:
// 1. Recibe mensaje como bytes desde Kafka
// 2. Deserializa el JSON a un mapa
// 3. Extrae event_type y plant_name
// 4. Convierte los datos a JSON string
// 5. Crea una entidad EventEntity
// 6. Guarda en PostgreSQL usando el repositorio
//
// CAMBIO REALIZADO: Completamente reescrito desde el TODO inicial
// RAZÓN: Implementar la persistencia de eventos en PostgreSQL
func (h *IntakeHandler) HandleMessage(message []byte) error {
	log.Printf("Received message on intake topic: %s", string(message))

	// CAMBIO: Parse del mensaje JSON
	// RAZÓN: Necesitamos extraer campos específicos (event_type, plant_name)
	var data map[string]interface{}
	if err := json.Unmarshal(message, &data); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return err
	}

	// CAMBIO: Extrae event_type del mensaje
	// RAZÓN: Indexamos por event_type para filtrado rápido en queries
	eventType := "unknown"
	if et, ok := data["event_type"].(string); ok {
		eventType = et
	}

	// CAMBIO: Extrae plant_name como source
	// RAZÓN: Permite identificar de qué planta viene cada evento
	source := "kafka-intake"
	if plantName, ok := data["plant_name"].(string); ok {
		source = plantName
	}

	// CAMBIO: Convierte data completo a JSON string
	// RAZÓN: PostgreSQL almacena el JSON completo como texto para consultas posteriores
	dataJSON, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling data: %v", err)
		return err
	}

	// CAMBIO: Crea entidad de evento
	// RAZÓN: Mapea el mensaje de Kafka a nuestra estructura de base de datos
	event := &entities.EventEntity{
		EventType: eventType,
		Source:    source,
		Data:      string(dataJSON),
	}

	// CAMBIO: Guarda en PostgreSQL
	// RAZÓN: Persiste el evento para consultas posteriores via API REST o DBeaver
	savedEvent, err := h.eventRepository.Create(event)
	if err != nil {
		log.Printf("Error saving event to database: %v", err)
		return err
	}

	log.Printf("Event saved to database with ID: %s, Type: %s", savedEvent.ID, savedEvent.EventType)
	return nil
}
