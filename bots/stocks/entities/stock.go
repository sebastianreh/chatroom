package entities

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"
)

const (
	MessageFormat = "%s quote is $%v per share"
)

type StockMessage struct {
	RoomID    string    `json:"room_id"`
	Message   string    `json:"bot_message"`
	CreatedAt time.Time `json:"created_at"`
}

func MapRecordsToStock(roomID string, records [][]string) (StockMessage, error) {
	var stock StockMessage
	if len(records) != 2 {
		return stock, errors.New("records len is not 2")
	}

	data := records[1]
	high, err := strconv.ParseFloat(data[4], 64)
	if err != nil {
		return stock, err
	}

	low, err := strconv.ParseFloat(data[5], 64)
	if err != nil {
		return stock, err
	}

	average := roundTo2Decimals((high + low) / 2)
	stockMessage := StockMessage{
		RoomID:    roomID,
		Message:   fmt.Sprintf(MessageFormat, data[0], average),
		CreatedAt: time.Now().UTC(),
	}

	return stockMessage, nil
}

func roundTo2Decimals(value float64) float64 {
	return math.Round(value*100) / 100
}
