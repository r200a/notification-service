package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/r200a/notification-service/internal/consumer"
	"github.com/r200a/notification-service/internal/service"
	"github.com/r200a/notification-service/pkg/config"
)

func main() {
	cfg := config.Load()

	svc := service.NewNotificationService(cfg)
	kafkaSvc := consumer.NewKafkaConsumer(cfg, svc)

	ctx, cancel := context.WithCancel(context.Background())

	// graceful shutdown on Ctrl+C or SIGTERM
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("shutdown signal received")
		cancel()
	}()

	kafkaSvc.Start(ctx)
	log.Println("notification service stopped")
}
