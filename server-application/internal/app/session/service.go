package session

import (
	"context"
	"errors"
	"fmt"
	"github.com/sebastianreh/chatroom/internal/config"
	"github.com/sebastianreh/chatroom/internal/entities"
	"github.com/sebastianreh/chatroom/pkg/logger"
	str "github.com/sebastianreh/chatroom/pkg/strings"
	"time"
)

const (
	serviceName     = "session.service"
	messageMaxLimit = 50
)

type SessionService interface {
	Join(ctx context.Context, sessionJoin entities.SessionChatRequest) (entities.JoinResponse, error)
	Exit(ctx context.Context, sessionExit entities.SessionChatRequest) error
	SaveMessage(ctx context.Context, message entities.ChatMessage, roomID string) error
	GetMessages(ctx context.Context, roomID string) ([]entities.ChatMessage, error)
}

type sessionService struct {
	config     config.Config
	repository SessionRepository
	logs       logger.Logger
}

func NewSessionService(cfg config.Config, repository SessionRepository, logger logger.Logger) SessionService {
	return &sessionService{
		config:     cfg,
		repository: repository,
		logs:       logger,
	}
}

func (service *sessionService) Join(ctx context.Context, sessionJoin entities.SessionChatRequest) (entities.JoinResponse, error) {
	var joinResponse entities.JoinResponse
	session, err := service.repository.Get(ctx, sessionJoin.RoomID)
	if err != nil {
		return joinResponse, err
	}

	if session.RoomID == str.Empty {
		session.RoomID = sessionJoin.RoomID
	}
	currentUsers, err := service.addUserToSession(sessionJoin.Username, session.CurrentUsers)
	if err != nil {
		joinResponse.Users = session.CurrentUsers
		return joinResponse, nil
	}
	session.CurrentUsers = currentUsers
	session.Events = append(session.Events, entities.Event{
		SessionUser: sessionJoin.SessionUser,
		Type:        entities.RoomActionEventType,
		CreatedAt:   time.Now().UTC(),
		Content:     entities.JoinContent,
	})

	err = service.repository.Set(ctx, session)
	joinResponse.Users = session.CurrentUsers

	if err != nil {
		return joinResponse, err
	}

	return joinResponse, nil
}

func (service *sessionService) addUserToSession(usernameLookUp string, sessionUsers []string) ([]string, error) {
	for _, user := range sessionUsers {
		if user == usernameLookUp {
			err := errors.New(fmt.Sprintf("user %s already in session", usernameLookUp))
			service.logs.Warn(str.ErrorConcat(err, serviceName, "addUserToSession"))
			return sessionUsers, err
		}
	}

	sessionUsers = append(sessionUsers, usernameLookUp)
	return sessionUsers, nil
}

func (service *sessionService) Exit(ctx context.Context, sessionExit entities.SessionChatRequest) error {
	session, err := service.repository.Get(ctx, sessionExit.RoomID)
	if err != nil {
		return err
	}

	if session.RoomID == str.Empty {
		err = errors.New("session is empty, user was not inside room")
		service.logs.Error(str.ErrorConcat(err, serviceName, "Set"))
		return err
	}

	session.CurrentUsers = removeUserFromSession(sessionExit.Username, session.CurrentUsers)
	session.Events = append(session.Events, entities.Event{
		SessionUser: sessionExit.SessionUser,
		Type:        entities.RoomActionEventType,
		CreatedAt:   time.Now().UTC(),
		Content:     entities.ExitContent,
	})

	err = service.repository.Set(ctx, session)
	if err != nil {
		return err
	}

	return nil
}

func removeUserFromSession(usernameLookUp string, sessionUsers []string) []string {
	for i, username := range sessionUsers {
		if username == usernameLookUp {
			sessionUsers = append(sessionUsers[:i], sessionUsers[i+1:]...)
			break
		}
	}
	return sessionUsers
}

func (service *sessionService) SaveMessage(ctx context.Context, message entities.ChatMessage, roomID string) error {
	session, err := service.repository.Get(ctx, roomID)
	if err != nil {
		return err
	}

	if session.RoomID == str.Empty {
		err = errors.New("session is empty, user was not logged in")
		service.logs.Error(str.ErrorConcat(err, serviceName, "Set"))
		return err
	}

	session.Events = append(session.Events, entities.Event{
		SessionUser: message.SessionUser,
		Type:        entities.UserMessageEventType,
		CreatedAt:   message.CreatedAt,
		Content:     message.Content,
	})

	err = service.repository.Set(ctx, session)
	if err != nil {
		return err
	}

	return nil
}

func (service *sessionService) GetMessages(ctx context.Context, roomID string) ([]entities.ChatMessage, error) {
	var messages []entities.ChatMessage
	session, err := service.repository.Get(ctx, roomID)
	if err != nil {
		return messages, err
	}

	if len(session.Events) == 0 {
		return messages, nil
	}

	index := len(session.Events) - messageMaxLimit
	if index < 0 {
		index = 0
	}

	for _, event := range session.Events[index:] {
		if event.Type == entities.UserMessageEventType {
			message := entities.ChatMessage{
				SessionUser: event.SessionUser,
				CreatedAt:   event.CreatedAt,
				Content:     event.Content,
			}
			messages = append(messages, message)
		}
	}

	return messages, nil
}
