package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"monitoring-energy-service/docs"
	"monitoring-energy-service/internal/infrastructure/adapters/rest"
	"monitoring-energy-service/internal/infrastructure/conf"
	"monitoring-energy-service/internal/infrastructure/conf/database"
	"monitoring-energy-service/internal/infrastructure/container"

	"github.com/caarlos0/env/v11"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	buildDate string
	gitCommit string
)

// @title           Monitoring Energy Service API
// @version         1.0
// @description     Go microservice template with hexagonal architecture
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9000
// @BasePath  /

// @schemes http https
func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Printf("build date: %s", buildDate)
	log.Printf("build git commit: %s", gitCommit)

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("WARNING - Couldn't load .env file, error: %v", err)
	}

	cfg := &conf.Config{}
	opts := env.Options{OnSet: conf.OnSetConfig}
	if err := env.ParseWithOptions(cfg, opts); err != nil {
		log.Fatalf("%+v\n", err)
	}

	port := cfg.Port
	environment := cfg.Env

	// Update swagger host dynamically
	docs.SwaggerInfo.Host = "localhost:" + port

	kafkaBrokers := strings.Split(cfg.ListKafkaBrokers, ",")
	consumerGroup := cfg.ConsumeGroup
	autoOffsetKafka := "latest"

	db, err := database.SetupDatabasePsql()
	if err != nil {
		log.Fatalf("Error when initializing database, error: %v", err)
	}

	err = database.PerformMigrations(db)
	if err != nil {
		log.Fatalf("Error when performing migrations to database, error: %v", err)
	}

	timeoutSeconds := cfg.HttpClientTimeout
	httpClient := &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}

	c := container.NewContainer(
		db,
		kafkaBrokers,
		consumerGroup,
		httpClient,
		autoOffsetKafka,
		container.WithConfig(*cfg),
	)

	// Start Kafka consumer in background
	go c.KafkaService.ConsumeEvents()

	router := gin.New()
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/healthz", "/readyz", "/swagger/*any"},
	}))
	router.Use(gin.Recovery())

	suffixesFromEnv := strings.Split(cfg.AllowedCorsSuffixes, ",")
	allowedSuffixes := make([]string, 0, len(suffixesFromEnv))
	for _, suffix := range suffixesFromEnv {
		trimmedSuffix := strings.TrimSpace(suffix)
		if trimmedSuffix != "" {
			allowedSuffixes = append(allowedSuffixes, trimmedSuffix)
		}
	}

	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost:") {
				return true
			}

			for _, suffix := range allowedSuffixes {
				if strings.HasSuffix(origin, suffix) {
					return true
				}
			}

			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	mainGroup := router.Group("")
	rest.SetupRoutes(mainGroup, c)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/healthz", gin.WrapF(HealthCheck))
	router.GET("/readyz", gin.WrapF(HealthCheck))

	if environment == "dev" {
		log.Printf("Running in development mode")
		log.Printf("Swagger UI available at http://localhost:%s/swagger/index.html", port)
	}

	log.Printf("Server starting on port %s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
