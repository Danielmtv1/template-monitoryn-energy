package rest

import (
	"errors"
	"net/http"

	"monitoring-energy-service/internal/domain/entities"
	domainerrors "monitoring-energy-service/internal/domain/errors"
	"monitoring-energy-service/internal/infrastructure/container"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateExampleRequest represents the request body for creating an example
type CreateExampleRequest struct {
	Name        string `json:"name" binding:"required" example:"My Example"`
	Description string `json:"description" example:"This is an example description"`
}

// UpdateExampleRequest represents the request body for updating an example
type UpdateExampleRequest struct {
	Name        string `json:"name" example:"Updated Example"`
	Description string `json:"description" example:"Updated description"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

// ListExamples godoc
// @Summary      List all examples
// @Description  Get all examples from the database
// @Tags         examples
// @Accept       json
// @Produce      json
// @Success      200  {array}   entities.ExampleEntity
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/examples [get]
func ListExamples(c *container.Container) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		examples, err := c.ExampleRepository.FindAll()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, examples)
	}
}

// GetExample godoc
// @Summary      Get an example by ID
// @Description  Get a single example by its UUID
// @Tags         examples
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Example ID (UUID)"
// @Success      200  {object}  entities.ExampleEntity
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /api/v1/examples/{id} [get]
func GetExample(c *container.Container) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		example, err := c.ExampleRepository.FindByID(id)
		if err != nil {
			if errors.Is(err, domainerrors.ErrNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "example not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		ctx.JSON(http.StatusOK, example)
	}
}

// CreateExample godoc
// @Summary      Create a new example
// @Description  Create a new example with the provided data
// @Tags         examples
// @Accept       json
// @Produce      json
// @Param        request  body      CreateExampleRequest  true  "Example data"
// @Success      201      {object}  entities.ExampleEntity
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /api/v1/examples [post]
func CreateExample(c *container.Container) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req CreateExampleRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		entity := &entities.ExampleEntity{
			Name:        req.Name,
			Description: req.Description,
		}

		created, err := c.ExampleRepository.Create(entity)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, created)
	}
}

// UpdateExample godoc
// @Summary      Update an example
// @Description  Update an existing example by ID
// @Tags         examples
// @Accept       json
// @Produce      json
// @Param        id       path      string                true  "Example ID (UUID)"
// @Param        request  body      UpdateExampleRequest  true  "Example data"
// @Success      200      {object}  entities.ExampleEntity
// @Failure      400      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /api/v1/examples/{id} [put]
func UpdateExample(c *container.Container) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		existing, err := c.ExampleRepository.FindByID(id)
		if err != nil {
			if errors.Is(err, domainerrors.ErrNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "example not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		var req UpdateExampleRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Name != "" {
			existing.Name = req.Name
		}
		if req.Description != "" {
			existing.Description = req.Description
		}

		updated, err := c.ExampleRepository.Update(existing)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, updated)
	}
}

// DeleteExample godoc
// @Summary      Delete an example
// @Description  Delete an example by ID
// @Tags         examples
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Example ID (UUID)"
// @Success      204  "No Content"
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/examples/{id} [delete]
func DeleteExample(c *container.Container) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		if err := c.ExampleRepository.Delete(id); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusNoContent, nil)
	}
}
