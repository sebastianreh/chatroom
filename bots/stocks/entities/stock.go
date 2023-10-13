package entities

import (
	"errors"
	"time"
)

const (
	timeLayout = "2006-01-02 15:04:05"
)

type Stock struct {
	Symbol   string
	DateTime time.Time
	Open     string
	High     string
	Low      string
	Close    string
	Volume   string
}

func MapRecordsToStock(records [][]string) (Stock, error) {
	var stock Stock
	if len(records) != 2 {
		return stock, errors.New("records len is not 2")
	}

	data := records[1]
	dateTimeStr := data[1] + " " + data[2]
	dateTime, err := time.Parse(timeLayout, dateTimeStr)
	if err != nil {
		return stock, err
	}

	stock.DateTime = dateTime
	stock.Symbol = data[0]
	stock.Open = data[3]
	stock.High = data[4]
	stock.Low = data[5]
	stock.Close = data[6]
	stock.Volume = data[7]

	return stock, nil
}
