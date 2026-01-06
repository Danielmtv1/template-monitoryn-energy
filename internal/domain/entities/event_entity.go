package entities

import (
	"time"

	"github.com/google/uuid"
)

// EventEntity representa un evento de energía capturado desde Kafka
//
// PROPÓSITO:
// Esta entidad almacena todos los eventos que llegan desde Kafka en la base de datos
// PostgreSQL, permitiendo un histórico completo de eventos para análisis posterior.
//
// CAMPOS:
// - ID: Identificador único UUID generado automáticamente por PostgreSQL
// - EventType: Tipo de evento (power_reading, status_update, efficiency_report, alert)
// - Source: Fuente del evento (nombre de la planta de energía)
// - Data: Datos completos del evento en formato JSON almacenado como texto
// - Metadata: Metadatos adicionales opcionales en formato JSON
// - CreatedAt: Timestamp de cuando se guardó el evento en la base de datos
//
// CAMBIO REALIZADO: Archivo creado desde cero
// RAZÓN: Necesitábamos una entidad para persistir eventos de Kafka en PostgreSQL
type EventEntity struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	EventType     string    `gorm:"type:varchar(100);index:idx_event_type;not null" json:"event_type"`
	PlantSourceId uuid.UUID `gorm:"type:uuid;index:idx_plant_source_id;not null" json:"plant_source_id"`
	Source        string    `gorm:"type:varchar(255)" json:"source"`
	Data          string    `gorm:"type:text" json:"data"` // Cambiado de jsonb a text para compatibilidad con GORM string type
	Metadata      string    `gorm:"type:text" json:"metadata,omitempty"`
	CreatedAt     time.Time `gorm:"autoCreateTime;index:idx_created_at" json:"created_at"`
	// Relaciones
	PlantSource EnergyPlants `gorm:"foreignKey:PlantSourceId;references:ID" json:"plant_source,omitempty"`
}

func (EventEntity) TableName() string {
	return "events"
}
