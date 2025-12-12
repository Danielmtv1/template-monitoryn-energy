package conf

import (
	"fmt"
	"log"
	"strings"
)

type Config struct {
	Port             string `env:"PORT" envDefault:"9000"`
	Env              string `env:"ENVIRONMENT" envDefault:"dev"`
	ListKafkaBrokers string `env:"LIST_KAFKA_BROKERS,required"`
	ConsumeGroup     string `env:"CONSUMER_GROUP,required"`
	HttpClientTimeout int    `env:"HTTP_CLIENT_TIMEOUT" envDefault:"30"`

	// Kafka Topics
	ConsumerTopic string `env:"CONSUMER_TOPIC" envDefault:"events.default"`
	ProducerTopic string `env:"PRODUCER_TOPIC" envDefault:"events.output"`

	// Webhook
	WebhookEnabled bool   `env:"WEBHOOK_ENABLED" envDefault:"false"`
	WebhookUrl     string `env:"WEBHOOK_URL"`

	// CORS
	AllowedCorsSuffixes string `env:"ALLOWED_CORS_SUFFIXES" envDefault:".spotcloud.io"`
}

func OnSetConfig(tag string, value interface{}, isDefault bool) {
	if strings.Contains(tag, "SECRET") ||
		strings.Contains(tag, "PASSWORD") ||
		strings.Contains(tag, "API_KEY") {
		if s, ok := value.(string); ok && len(s) > 4 {
			rune := []rune(s)
			repeatedStars := strings.Repeat("*", len(rune)-4)
			value = string(rune[:2]) + repeatedStars + string(rune[len(rune)-2:])
		}
	}

	msg := fmt.Sprintf("%s=%v", tag, value)
	if isDefault {
		msg += " (default)"
	}

	log.Println(msg)
}
