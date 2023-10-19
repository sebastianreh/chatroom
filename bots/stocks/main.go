package main

import (
	"encoding/json"
	"fmt"
	"github.com/sebastianreh/chatroom-bots/stocks/entities"
	ctr "github.com/sebastianreh/chatroom-bots/stocks/internal/container"
	"github.com/sebastianreh/chatroom-bots/stocks/pkg/csv"
)

const (
	BotName     = "stock"
	emptyString = ""
)

func main() {
	container := ctr.Build(BotName)
	ProcessMessages(container)
}

func ProcessMessages(container ctr.Container) {
	for {
		message, err := container.Socket.ReadMessage()
		if err != nil {
			container.Logs.Error(fmt.Sprintf("error reading message: %s", err.Error()), "ProcessMessages")
			continue
		}

		if message.Value == emptyString {
			continue
		}

		stockCSV, err := container.Client.GetStockCSV(message.Value)
		if err != nil {
			continue
		}

		stockRecords, err := csv.ReadStockCsv(stockCSV)
		if err != nil {
			container.Logs.Error(fmt.Sprintf("error reading stock csv: %s", err.Error()), "ProcessMessages")
			continue
		}

		stockMessage, err := entities.MapRecordsToStock(container.RoomID, stockRecords)
		if err != nil {
			container.Logs.Error(fmt.Sprintf("error mapping records for stock: %s", err.Error()), "ProcessMessages")
			continue
		}

		stockMessageBytes, err := json.Marshal(stockMessage)
		if err != nil {
			container.Logs.Error(fmt.Sprintf("error marshaling message: %s", err.Error()), "ProcessMessages")
			continue
		}

		err = container.Producer.Send(container.Config.Kafka.StocksTopic, stockMessageBytes)
		if err != nil {
			container.Logs.Error(fmt.Sprintf("error initializing kafka producer: %s", err.Error()), "ProcessMessages")
		}
	}
}
