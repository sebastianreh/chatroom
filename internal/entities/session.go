package entities

type SessionRequest struct {
	RoomID string      `json:"room_id" validate:"required"`
	User   SessionUser `json:"user" validate:"required"`
}

type Session struct {
	RoomID string  `json:"room_id"`
	Events []Event `json:"users"`
}

type SessionUser struct {
	Username string `json:"username" validate:"required"`
	ID       string `json:"id" validate:"required"`
}
