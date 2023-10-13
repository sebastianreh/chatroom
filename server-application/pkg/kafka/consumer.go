package kafka

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sebastianreh/chatroom/internal/config"
	"github.com/sebastianreh/chatroom/pkg/logger"
)

type Consumer struct {
	topic    string
	groupID  string
	consumer *kafka.Consumer
}

// NewConsumer initializes a new Consumer struct.
func NewConsumer(cfg config.Config, factoryConsumer *kafka.Consumer) *Consumer {
	return &Consumer{
		topic:    cfg.Kafka.StocksTopic,
		groupID:  cfg.Kafka.GroupID,
		consumer: factoryConsumer,
	}
}

func (c *Consumer) Listen(process func(context.Context, []byte)) error {
	if err := c.consumer.Subscribe(c.topic, nil); err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	for {
		event := c.consumer.Poll(100)
		if event == nil {
			continue
		}
		switch e := event.(type) {
		case *kafka.Message:
			process(context.Background(), e.Value)
		case kafka.Error:
			return fmt.Errorf("consumer error: %w", e)
		default:
			// No-op for other events
		}
	}
}

func NewKafkaConsumer(cfg config.Config, logger logger.Logger) (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.Server,  // Assuming your config has a Kafka struct with a Server field
		"group.id":          cfg.Kafka.GroupID, // And a GroupID field
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		logger.Error("failed to create Kafka consumer", err)
		return nil, err
	}

	return consumer, nil
}
