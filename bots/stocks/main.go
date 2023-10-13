package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sebastianreh/chatroom-bots/stocks/entities"
	"github.com/sebastianreh/chatroom-bots/stocks/internal/config"
	"github.com/sebastianreh/chatroom-bots/stocks/pkg/csv"
	"github.com/sebastianreh/chatroom-bots/stocks/pkg/kafka"
	"github.com/sebastianreh/chatroom-bots/stocks/pkg/rest"
	"github.com/sebastianreh/chatroom-bots/stocks/pkg/websocket"
	"github.com/sebastianreh/chatroom/pkg/logger"
	"regexp"
)

const maxRetires = 10
const roomID = "651f0027d2c7b1f7d56ae334"
const stockCommand = "/stock=stock_code"
const BotName = "stock"

func main() {
	cfg := config.NewConfig()
	log := logger.NewLogger()
	restyClient := resty.New()
	factoryProducer, err := kafka.NewFactoryProducer(log, cfg.Kafka.Server, kafka.WithMaxRetries(maxRetires))
	if err != nil {
		log.Error("error initializing kafka producer", "main", err.Error())
	}

	socket := websocket.NewWebsocket(log, cfg, BotName, roomID)

	for {
		message, err := socket.ReadMessage()
		if err != nil {
			log.Error("error initializing kafka producer", "main", err.Error())
			continue
		}

		stockCode := parseStockCodeFromMessage(message)

		if stockCode == "" {
			continue
		}

		client := rest.NewStooqClient(log, restyClient)
		stockCSV, _ := client.GetStockCSV(stockCode)

		stockRecords, err := csv.ReadStockCsv(stockCSV)

		stockMessage, err := entities.MapRecordsToStock(roomID, stockRecords)
		if err != nil {
			log.Error("error initializing kafka producer", "main", err.Error())
		}

		stockMessageBytes, err := json.Marshal(stockMessage)

		err = factoryProducer.Send(cfg.Kafka.StocksTopic, stockMessageBytes)
		if err != nil {
			log.Error("error initializing kafka producer", "main", err.Error())
		}
	}
}

func parseStockCodeFromMessage(message string) string {
	var stockCode string
	re := regexp.MustCompile(fmt.Sprintf(`/%s=(.+)`, BotName))
	matches := re.FindStringSubmatch(message)

	if len(matches) > 1 {
		stockCode := matches[1]
		fmt.Println("Stock Code:", stockCode)
	} else {
		fmt.Println("No stock code found in the input string")
	}

	return stockCode
}
