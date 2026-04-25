package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/r200a/notification-service/internal/model"
	"github.com/r200a/notification-service/internal/service"
	"github.com/r200a/notification-service/pkg/config"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader  *kafka.Reader
	service *service.NotificationService
}

func NewKafkaConsumer(cfg *config.Config, svc *service.NotificationService) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.KafkaBroker},
		Topic:    cfg.KafkaTopic,
		GroupID:  cfg.KafkaGroupID,
		MinBytes: 1,
		MaxBytes: 10e6, // 10MB
	})

	return &KafkaConsumer{
		reader:  reader,
		service: svc,
	}
}

func (c *KafkaConsumer) Start(ctx context.Context) {
	log.Println("Kafka consumer started — waiting for events...")

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("consumer shutting down")
				break
			}
			log.Printf("error reading message: %v", err)
			continue
		}

		log.Printf("received event: topic=%s partition=%d offset=%d",
			msg.Topic, msg.Partition, msg.Offset)

		var event model.ApplicationEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("failed to parse event: %v — raw: %s", err, string(msg.Value))
			continue
		}

		c.service.HandleEvent(event)
	}

	if err := c.reader.Close(); err != nil {
		log.Printf("error closing reader: %v", err)
	}
}
