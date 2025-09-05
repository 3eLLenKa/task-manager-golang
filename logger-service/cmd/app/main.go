package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"todo/kafka-logger/internal/consumer"
)

func main() {
	brokers := []string{"kafka:9092"}
	topic := "tasks"
	groupID := "logger-group"
	logPath := "events.log"

	c, err := consumer.New(brokers, topic, groupID, logPath)

	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go c.Start(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	cancel()

	if err := c.Close(); err != nil {
		log.Printf("failed to close consumer: %v", err)
	}

	log.Println("logger-service stopped")
}
