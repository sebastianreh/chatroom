package entities

type SessionRequest struct {
	RoomID string `json:"room_id" validate:"required" query:"room_id"`
	SessionUser
}

type Session struct {
	RoomID string  `json:"room_id"  validate:"required"`
	Events []Event `json:"users"  validate:"required"`
}

type SessionUser struct {
	Username string `json:"username" validate:"required" query:"username"`
	UserID   string `json:"user_id" validate:"required" query:"user_id"`
}
