package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	KafkaBroker    string
	KafkaTopic     string
	KafkaGroupID   string
	SendGridAPIKey string
	FromEmail      string
	FromName       string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, reading from environment")
	}

	return &Config{
		KafkaBroker:    getEnv("KAFKA_BROKER", "localhost:9092"),
		KafkaTopic:     getEnv("KAFKA_TOPIC", "application.events"),
		KafkaGroupID:   getEnv("KAFKA_GROUP_ID", "notification-service"),
		SendGridAPIKey: getEnv("SENDGRID_API_KEY", ""),
		FromEmail:      getEnv("FROM_EMAIL", "noreply@vcplatform.com"),
		FromName:       getEnv("FROM_NAME", "VC Platform"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
