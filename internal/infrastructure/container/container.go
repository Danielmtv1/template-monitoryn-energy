package container

import (
	"net/http"

	"monitoring-energy-service/internal/api"
	"monitoring-energy-service/internal/domain/ports/input"
	"monitoring-energy-service/internal/domain/ports/output"
	"monitoring-energy-service/internal/infrastructure/adapters/http/webhook"
	"monitoring-energy-service/internal/infrastructure/adapters/kafka"
	"monitoring-energy-service/internal/infrastructure/adapters/repositories"
	"monitoring-energy-service/internal/infrastructure/conf"
	"monitoring-energy-service/internal/infrastructure/conf/kafkaconf"

	"gorm.io/gorm"
)

type ContainerOption func(*Container)

// Container mantiene todas las dependencias de la aplicación (Dependency Injection)
//
// CAMBIOS REALIZADOS:
// - Agregado EventRepository (línea 26): Para acceso a base de datos de eventos
// - Agregado EventGenerator (línea 27): Para generar eventos automáticamente
type Container struct {
	db                *gorm.DB
	cfg               conf.Config
	KafkaService      input.KafkaServiceInterface
	WebhookAdapter    output.WebhookAdapterInterface
	ExampleRepository output.ExampleRepositoryInterface
	EventRepository   output.EventRepositoryInterface // CAMBIO: Agregado para gestionar eventos en DB
	EventGenerator    *api.EventGenerator             // CAMBIO: Agregado para generar eventos cada 5 min
}

func NewContainer(
	db *gorm.DB,
	kafkaBrokers []string,
	consumerGroup string,
	httpClient *http.Client,
	autoOffset string,
	opts ...ContainerOption,
) *Container {
	container := &Container{db: db}

	for _, opt := range opts {
		opt(container)
	}

	// Initialize repositories
	exampleRepository := repositories.NewExampleRepository(db)
	container.ExampleRepository = exampleRepository

	// CAMBIO: Inicializa EventRepository
	// RAZÓN: Necesario para que IntakeHandler y REST API puedan acceder a eventos en DB
	eventRepository := repositories.NewEventRepository(db)
	container.EventRepository = eventRepository

	// Initialize Kafka
	kafkaFactory := kafkaconf.NewKafkaFactory(kafkaBrokers, autoOffset)
	kafkaAdapter := kafka.NewKafkaAdapter(kafkaFactory, consumerGroup)
	kafkaService := api.NewKafkaService(kafkaAdapter)
	container.KafkaService = kafkaService

	// Initialize Webhook adapter
	webhookAdapter := webhook.NewAdapter(httpClient)
	container.WebhookAdapter = webhookAdapter

	// Register Kafka handlers here
	// CAMBIO: IntakeHandler ahora recibe eventRepository como parámetro
	// RAZÓN: Necesita el repositorio para guardar eventos consumidos desde Kafka en PostgreSQL
	intakeHandler := api.NewIntakeHandler(eventRepository)
	kafkaService.RegisterHandler(container.cfg.ConsumerTopic, intakeHandler)

	// CAMBIO: Inicializa Event Generator con topic "intake"
	// RAZÓN: Genera automáticamente 30 eventos cada 5 minutos enviándolos a Kafka
	eventGenerator := api.NewEventGenerator(kafkaService, "intake")
	container.EventGenerator = eventGenerator

	return container
}

func WithConfig(config conf.Config) ContainerOption {
	return func(c *Container) {
		c.cfg = config
	}
}

func (c *Container) GetConfig() conf.Config {
	return c.cfg
}
