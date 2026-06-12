package kafka

import (
	"os"
	"strings"
)

type Config struct {
	Brokers []string
	Enabled bool
}

// LoadConfig loads Kafka configuration from environment variables
func LoadConfig() Config {
	brokersStr := os.Getenv("KAFKA_BROKERS")
	if brokersStr == "" {
		brokersStr = "localhost:9093"
	}

	enabled := os.Getenv("KAFKA_ENABLED")
	if enabled == "" {
		enabled = "true"
	}

	return Config{
		Brokers: strings.Split(brokersStr, ","),
		Enabled: enabled == "true",
	}
}
