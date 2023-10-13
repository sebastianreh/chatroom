package kafka

import (
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sebastianreh/chatroom/pkg/logger"
	"github.com/sebastianreh/chatroom/pkg/strings"
	"time"
)

const (
	defaultFlushMs    = 1000
	defaultWaitSecs   = 2
	defaultMaxRetries = 10
)

type FactoryProducer interface {
	Send(topic string, message []byte) (err error)
}

type producer struct {
	logs          logger.Logger
	serverAddress string
	flushMs       int
	waitSecs      int
	maxRetries    int
}

type ProducerOption func(*producer)

func NewFactoryProducer(logger logger.Logger, serverAddress string, opts ...ProducerOption) (FactoryProducer, error) {
	if strings.IsEmpty(serverAddress) {
		return &producer{}, errors.New("error, the serverAddress variable is empty")
	}

	producer := &producer{
		logs:          logger,
		serverAddress: serverAddress,
		flushMs:       defaultFlushMs,
		waitSecs:      defaultWaitSecs,
		maxRetries:    defaultMaxRetries,
	}

	for _, opt := range opts {
		opt(producer)
	}

	return producer, nil
}

func (producer *producer) Send(topic string, message []byte) (err error) {
	retries := 1
	for retries <= producer.maxRetries {
		err = producer.sendMessage(topic, message)
		if err == nil {
			return nil
		}
		time.Sleep(time.Duration(producer.waitSecs) * time.Second)
		retries++
	}
	return err
}

func (producer *producer) sendMessage(topic string, message []byte) error {
	prd, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": producer.serverAddress,
	})
	if err != nil {
		producer.logs.Error("error creating kafka kafkaProducer", err, "sendMessage")
		return err
	}
	defer prd.Close()
	var anyErr error
	producer.logs.Info(fmt.Sprintf("sending -> topic [%s] message %s", topic, message), "sendMessage")
	err = prd.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: message,
	}, nil)
	if err != nil {
		anyErr = err
	}

	prd.Flush(producer.flushMs)
	if anyErr != nil {
		producer.logs.Error("Error sending message: %+v", anyErr)
		return err
	}
	return nil
}

func WithMaxRetries(maxRetries int) ProducerOption {
	return func(h *producer) {
		h.maxRetries = maxRetries
	}
}
