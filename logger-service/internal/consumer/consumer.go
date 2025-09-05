package consumer

import (
	"context"
	"log"
	"todo/kafka-logger/internal/logger"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	logger *logger.Logger
}

func New(brokers []string, topic, groupID, logPath string) (*Consumer, error) {
	logWriter, err := logger.New(logPath)
	if err != nil {
		return nil, err
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	return &Consumer{
		reader: reader,
		logger: logWriter,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("failed to read message: %v", err)
			continue
		}

		msg := string(m.Value)
		c.logger.Write(msg)

		log.Printf("logged event: %s", msg)
	}
}

func (c *Consumer) Close() error {
	if err := c.reader.Close(); err != nil {
		return err
	}
	return c.logger.Close()
}
