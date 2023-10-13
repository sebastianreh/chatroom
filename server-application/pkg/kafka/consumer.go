package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sebastianreh/chatroom/internal/config"
	"github.com/sebastianreh/chatroom/pkg/logger"
)

type Consumer interface {
	Listen(process func([]byte)) error
}

type consumer struct {
	topic    string
	listener *kafka.Consumer
}

func NewKafkaConsumer(cfg config.Config, logger logger.Logger) (Consumer, error) {
	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.Server,  // Assuming your config has a Kafka struct with a Server field
		"group.id":          cfg.Kafka.GroupID, // And a GroupID field
		"auto.offset.reset": "latest",
	})

	csr := new(consumer)
	csr.listener = kafkaConsumer
	csr.topic = cfg.Kafka.StocksTopic

	if err != nil {
		logger.Error("failed to create Kafka consumer", err)
		return csr, err
	}

	return csr, nil
}

func (c *consumer) Listen(process func([]byte)) error {
	if err := c.listener.Subscribe(c.topic, nil); err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}
	for {
		event := c.listener.Poll(100)
		if event == nil {
			continue
		}
		switch e := event.(type) {
		case *kafka.Message:
			process(e.Value)
		case kafka.Error:
			return fmt.Errorf("consumer error: %w", e)
		default:
		}
	}
}
