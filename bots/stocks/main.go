package main

import (
	"encoding/json"
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
			container.Logs.Error("error reading message", "ProcessMessages", err.Error())
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
			container.Logs.Error("error reading stock csv", "ProcessMessages", err.Error())
			continue
		}

		stockMessage, err := entities.MapRecordsToStock(container.RoomID, stockRecords)
		if err != nil {
			container.Logs.Error("error mapping records for stock", "ProcessMessages", err.Error())
			continue
		}

		stockMessageBytes, err := json.Marshal(stockMessage)
		if err != nil {
			container.Logs.Error("error marshaling message", "ProcessMessages", err.Error())
			continue
		}

		err = container.Producer.Send(container.Config.Kafka.StocksTopic, stockMessageBytes)
		if err != nil {
			container.Logs.Error("error initializing kafka producer", "ProcessMessages", err.Error())
		}
	}
}
