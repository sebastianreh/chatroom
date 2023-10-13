package entities

import "time"

type Stock struct {
	Symbol   string
	DateTime time.Time
	Open     string
	High     string
	Low      string
	Close    string
	Volume   string
}
