package entities

import "encoding/json"

const (
	JoinAction = "join"
	ExitAction = "exit"
)

type SessionRequest struct {
	RoomID string `json:"room_id" validate:"required" query:"room_id"`
	SessionUser
}

type Session struct {
	RoomID       string   `json:"room_id"  validate:"required"`
	CurrentUsers []string `json:"current_users"`
	Events       []Event  `json:"events"  validate:"required"`
}

type SessionUser struct {
	Username string `json:"username" validate:"required" query:"username"`
	UserID   string `json:"user_id" validate:"required" query:"user_id"`
}

type JoinResponse struct {
	Users []string `json:"users"`
}

type SessionAction struct {
	SessionUser
	Type string `json:"type"`
}

func GetJoinAction(sessionUser SessionUser) SessionAction {
	return getSessionAction(sessionUser, JoinAction)
}

func GetExitAction(sessionUser SessionUser) SessionAction {
	return getSessionAction(sessionUser, ExitAction)
}

func getSessionAction(sessionUser SessionUser, actionType string) SessionAction {
	return SessionAction{
		Type:        actionType,
		SessionUser: sessionUser,
	}
}

func (s SessionAction) ToBytes() []byte {
	sBytes, _ := json.Marshal(s)
	return sBytes
}
