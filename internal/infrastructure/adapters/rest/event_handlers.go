package rest

// event_handlers.go - Handlers REST para consultar eventos guardados en PostgreSQL
//
// PROPÓSITO:
// Expone endpoints HTTP para consultar los eventos que fueron enviados a Kafka
// y guardados en PostgreSQL por el IntakeHandler.
//
// CAMBIO REALIZADO: Archivo creado desde cero
// RAZÓN: Necesitábamos una API REST para consultar eventos desde cualquier cliente HTTP
//
// ENDPOINTS CREADOS:
// - GET /api/v1/events           - Lista todos los eventos (ordenados por fecha DESC)
// - GET /api/v1/events/:id       - Obtiene un evento específico por UUID
// - GET /api/v1/events/type/:type - Filtra eventos por tipo (power_reading, alert, etc.)

import (
	"errors"
	"net/http"

	domainerrors "monitoring-energy-service/internal/domain/errors"
	"monitoring-energy-service/internal/infrastructure/container"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListEvents obtiene todos los eventos de la base de datos
// CAMBIO: Handler nuevo
// RAZÓN: Permite consultar el histórico completo de eventos vía HTTP
//
// ListEvents godoc
// @Summary      List all events
// @Description  Get all events from the database ordered by creation time
// @Tags         events
// @Accept       json
// @Produce      json
// @Success      200  {array}   entities.EventEntity
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/events [get]
func ListEvents(c *container.Container) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		events, err := c.EventRepository.FindAll()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, events)
	}
}

// GetEvent obtiene un evento específico por su ID
// CAMBIO: Handler nuevo
// RAZÓN: Permite recuperar un evento individual para análisis detallado
// CAMBIO: Ahora distingue entre errores 404 (not found) y 500 (internal)
// RAZÓN: Proporciona respuestas HTTP más precisas según el tipo de error
//
// GetEvent godoc
// @Summary      Get an event by ID
// @Description  Get a single event by its UUID
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Event ID (UUID)"
// @Success      200  {object}  entities.EventEntity
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /api/v1/events/{id} [get]
func GetEvent(c *container.Container) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		event, err := c.EventRepository.FindByID(id)
		if err != nil {
			if errors.Is(err, domainerrors.ErrNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		ctx.JSON(http.StatusOK, event)
	}
}

// GetEventsByType filtra eventos por tipo
// CAMBIO: Handler nuevo
// RAZÓN: Permite análisis de eventos específicos (ej: solo "alerts" o solo "power_reading")
//
// GetEventsByType godoc
// @Summary      Get events by type
// @Description  Get all events filtered by event type
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        type   path      string  true  "Event Type"
// @Success      200  {array}   entities.EventEntity
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/events/type/{type} [get]
func GetEventsByType(c *container.Container) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		eventType := ctx.Param("type")

		events, err := c.EventRepository.FindByEventType(eventType)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, events)
	}
}
