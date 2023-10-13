package main

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/sebastianreh/chatroom-bots/stocks/internal/config"
	"github.com/sebastianreh/chatroom-bots/stocks/pkg/csv"
	"github.com/sebastianreh/chatroom-bots/stocks/pkg/kafka"
	"github.com/sebastianreh/chatroom-bots/stocks/pkg/rest"
	"github.com/sebastianreh/chatroom/pkg/logger"
)

const maxRetires = 10

func main() {
	cfg := config.NewConfig()
	log := logger.NewLogger()
	restyClient := resty.New()

	client := rest.NewStooqClient(log, restyClient)
	stockCSV, _ := client.GetStockCSV("aapl.us")

	stock, err := csv.ReadStockCsv(stockCSV)
	if err != nil {
		log.Error("error initializing kafka producer", "main", err)
	}

	factoryProducer, err := kafka.NewFactoryProducer(log, cfg.Kafka.Server, kafka.WithMaxRetries(maxRetires))
	if err != nil {
		log.Error("error initializing kafka producer", "main", err)
	}

	stockBytes, err := json.Marshal(stock)

	err = factoryProducer.Send(cfg.Kafka.StocksTopic, stockBytes)
	if err != nil {
		log.Error("error initializing kafka producer", "main", err)
	}
}
