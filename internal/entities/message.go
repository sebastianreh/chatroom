package entities

import "time"

type Message struct {
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
}
