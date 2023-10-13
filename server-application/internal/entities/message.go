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
