package rest

import (
	"monitoring-energy-service/internal/infrastructure/container"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup, c *container.Container) {
	api := router.Group("/api/v1")
	{
		examples := api.Group("/examples")
		{
			examples.GET("", ListExamples(c))
			examples.POST("", CreateExample(c))
			examples.GET("/:id", GetExample(c))
			examples.PUT("/:id", UpdateExample(c))
			examples.DELETE("/:id", DeleteExample(c))
		}
	}
}
