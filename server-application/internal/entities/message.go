package entities

import "time"

const (
	RoomActionEventType  = "room_action"
	UserMessageEventType = "message"
	JoinContent          = "user_join"
	ExitContent          = "user_exit"
)

type Event struct {
	SessionUser
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
}

type ChatMessage struct {
	SessionUser
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
}

type BotMessage struct {
	Command string `json:"command"`
	Value   string `json:"value"`
}

type StockMessage struct {
	RoomID    string    `json:"room_id"`
	Message   string    `json:"bot_message"`
	CreatedAt time.Time `json:"created_at"`
}
