package api

import (
	"encoding/json"
	"log"

	"monitoring-energy-service/internal/domain/ports/input"
)

type IntakeHandler struct {
	// Add any dependencies you need (repositories, services, etc.)
}

var _ input.MessageHandler = &IntakeHandler{}

func NewIntakeHandler() *IntakeHandler {
	return &IntakeHandler{}
}

func (h *IntakeHandler) HandleMessage(message []byte) error {
	log.Printf("Received message on intake topic: %s", string(message))

	// TODO: Implement your business logic here
	// Example: Parse the message and process it
	var data map[string]interface{}
	if err := json.Unmarshal(message, &data); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return err
	}

	// Process the data as needed
	log.Printf("Processed data: %+v", data)
	// save data to database or trigger other actions
	return nil
}
