package api

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"monitoring-energy-service/internal/domain/ports/input"

	"github.com/google/uuid"
)

// EventGenerator genera automáticamente eventos de monitoreo de energía
//
// PROPÓSITO:
// Simula un sistema de monitoreo de plantas de energía enviando eventos periódicos
// a Kafka. Esto permite demostrar el flujo completo: Kafka → Consumer → PostgreSQL → API REST
//
// CAMBIO REALIZADO: Archivo creado desde cero
// RAZÓN: Necesitábamos generar eventos automáticamente cada 5 minutos según requisitos
type EventGenerator struct {
	kafkaService input.KafkaServiceInterface // Servicio para enviar mensajes a Kafka
	topic        string                      // Tópico de Kafka donde se envían los eventos ("intake")
	stopChan     chan struct{}               // Canal para detener el generador de forma segura
}

// EnergyMonitoringEvent estructura de datos para eventos de monitoreo de energía
//
// CAMBIO: Estructura nueva
// RAZÓN: Define el formato de los eventos que se envían a Kafka con datos realistas
// de plantas de energía (potencia generada, eficiencia, temperatura, etc.)
type EnergyMonitoringEvent struct {
	PlantID        string    `json:"plant_id"`            // Identificador de la planta (plant-1, plant-2, etc.)
	PlantSourceId  uuid.UUID `json:"plant_source_id"`     // UUID real de la planta en la base de datos
	PlantName      string    `json:"plant_name"`          // Nombre descriptivo de la planta
	EventType      string    `json:"event_type"`          // Tipo: power_reading, status_update, efficiency_report, alert
	PowerGenerated float64   `json:"power_generated_mw"`  // Potencia generada en megavatios (0-1000 MW)
	PowerConsumed  float64   `json:"power_consumed_mw"`   // Potencia consumida en megavatios (0-50 MW)
	Efficiency     float64   `json:"efficiency_percent"`  // Eficiencia de la planta (75-95%)
	Temperature    float64   `json:"temperature_celsius"` // Temperatura de operación (20-50°C)
	Status         string    `json:"status"`              // Estado: operational, maintenance, standby, peak_load
	Timestamp      time.Time `json:"timestamp"`           // Timestamp del evento
}

// NewEventGenerator crea una nueva instancia del generador de eventos
// PARÁMETROS:
// - kafkaService: Servicio de Kafka para enviar mensajes
// - topic: Nombre del tópico de Kafka (normalmente "intake")
func NewEventGenerator(kafkaService input.KafkaServiceInterface, topic string) *EventGenerator {
	return &EventGenerator{
		kafkaService: kafkaService,
		topic:        topic,
		stopChan:     make(chan struct{}),
	}
}

// Start inicia el generador de eventos en modo continuo
//
// FUNCIONAMIENTO:
// 1. Envía el primer lote de 30 eventos inmediatamente al iniciar
// 2. Luego envía 30 eventos cada 5 minutos de forma automática
// 3. Se ejecuta en un goroutine separado (llamado con 'go')
//
// CAMBIO: Método nuevo
// RAZÓN: Implementa el requisito de enviar 30 eventos cada 5 minutos automáticamente
func (eg *EventGenerator) Start() {
	log.Println("Starting Event Generator - will send 30 messages every 5 minutes")

	// CAMBIO: Envía el primer lote inmediatamente
	// RAZÓN: Permite ver eventos de inmediato sin esperar 60 minutos
	eg.generateAndSendEvents()

	// CAMBIO: Configura ticker para enviar cada 60 minutos
	// RAZÓN: Cumple con el requisito de enviar eventos periódicamente
	ticker := time.NewTicker(60 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Cada 60 minutos, genera y envía 30 eventos
			eg.generateAndSendEvents()
		case <-eg.stopChan:
			// Permite detener el generador de forma limpia
			log.Println("Stopping Event Generator")
			return
		}
	}
}

// generateAndSendEvents genera y envía exactamente 30 eventos a Kafka
//
// FUNCIONAMIENTO:
// 1. Define arrays de nombres de plantas, estados y tipos de eventos
// 2. En un loop genera 30 eventos con datos aleatorios pero realistas
// 3. Envía cada evento a Kafka usando KafkaService
// 4. Espera 100ms entre eventos para no saturar Kafka
//
// CAMBIO: Método nuevo
// RAZÓN: Implementa la lógica de generación de 30 eventos con datos simulados realistas
func (eg *EventGenerator) generateAndSendEvents() {
	log.Println("Generating and sending 30 events to Kafka...")

	// CAMBIO: Define plantas con sus UUIDs reales de la base de datos
	// RAZÓN: Usa los UUIDs exactos del archivo de seed para relacionar eventos con plantas reales
	type PlantInfo struct {
		ID   uuid.UUID
		Name string
	}

	plants := []PlantInfo{
		{
			ID:   uuid.MustParse("1e2d3c4b-5a6f-7e8d-9c0b-1a2b3c4d5e6f"),
			Name: "Solar Plant Alpha",
		},
		{
			ID:   uuid.MustParse("2f3e4d5c-6b7a-8c9d-0e1f-2b3c4d5e6f7a"),
			Name: "Wind Farm Beta",
		},
		{
			ID:   uuid.MustParse("c2e78b94-76f4-49a4-b6e8-d62c8d1d23ea"),
			Name: "Hydro Plant Gamma",
		},
	}

	statuses := []string{"operational", "maintenance", "standby", "peak_load"}
	eventTypes := []string{"power_reading", "status_update", "efficiency_report", "alert"}

	// CAMBIO: Loop que genera exactamente 30 eventos
	// RAZÓN: Cumple con el requisito específico de "30 mensajes cada 60 minutos"
	for i := 0; i < 30; i++ {
		// CAMBIO: Selecciona una planta aleatoria de las reales
		// RAZÓN: Usa plantas reales de la base de datos con sus UUIDs
		selectedPlant := plants[rand.Intn(len(plants))]

		// CAMBIO: Genera datos aleatorios pero dentro de rangos realistas
		// RAZÓN: Simula datos reales de plantas de energía para testing
		event := EnergyMonitoringEvent{
			PlantID:        fmt.Sprintf("plant-%d", rand.Intn(len(plants))+1),
			PlantSourceId:  selectedPlant.ID,
			PlantName:      selectedPlant.Name,
			EventType:      eventTypes[rand.Intn(len(eventTypes))],
			PowerGenerated: rand.Float64() * 1000,  // 0-1000 MW
			PowerConsumed:  rand.Float64() * 50,    // 0-50 MW
			Efficiency:     75 + rand.Float64()*20, // 75-95%
			Temperature:    20 + rand.Float64()*30, // 20-50°C
			Status:         statuses[rand.Intn(len(statuses))],
			Timestamp:      time.Now(),
		}

		// CAMBIO: Crea una key única para el mensaje de Kafka
		// RAZÓN: Kafka usa keys para particionar mensajes y garantizar orden
		key := fmt.Sprintf("%s-%d", event.PlantID, time.Now().Unix())
		err := eg.kafkaService.SendEvent(eg.topic, key, event)
		if err != nil {
			log.Printf("Error sending event %d to Kafka: %v", i+1, err)
		} else {
			log.Printf("Event %d sent: PlantID=%s, Type=%s, Power=%.2fMW",
				i+1, event.PlantID, event.EventType, event.PowerGenerated)
		}

		// CAMBIO: Pausa de 100ms entre eventos
		// RAZÓN: Evita saturar Kafka y permite ver el flujo de eventos en los logs
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("Finished sending 30 events to Kafka")
}

// Stop detiene el generador de eventos de forma segura
// CAMBIO: Método nuevo
// RAZÓN: Permite apagar el generador sin memory leaks cerrando el canal stopChan
func (eg *EventGenerator) Stop() {
	close(eg.stopChan)
}
