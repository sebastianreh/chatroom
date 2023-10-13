package entities

import "time"

type StockMessage struct {
	RoomID    string    `json:"room_id"`
	Message   string    `json:"bot_message"`
	CreatedAt time.Time `json:"created_at"`
}
