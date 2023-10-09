package session

import (
	"context"
	"errors"
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
	Join(ctx context.Context, sessionJoin entities.SessionRequest) error
	Exit(ctx context.Context, sessionExit entities.SessionRequest) error
	SaveMessage(ctx context.Context, user entities.SessionUser, roomID, message string) error
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

func (service *sessionService) Join(ctx context.Context, sessionJoin entities.SessionRequest) error {
	session, err := service.repository.Get(ctx, sessionJoin.RoomID)
	if err != nil {
		return err
	}

	if session.RoomID == str.Empty {
		session.RoomID = sessionJoin.RoomID
		session.Events = append(session.Events, entities.Event{
			SessionUser: sessionJoin.SessionUser,
			Type:        entities.RoomActionEventType,
			CreatedAt:   time.Now().UTC(),
			Content:     entities.JoinContent,
		})
	} else if session.RoomID == sessionJoin.RoomID {
		session.Events = append(session.Events, entities.Event{
			SessionUser: sessionJoin.SessionUser,
			Type:        entities.RoomActionEventType,
			CreatedAt:   time.Now().UTC(),
			Content:     entities.JoinContent,
		})
	}

	err = service.repository.Set(ctx, session)
	if err != nil {
		return err
	}

	return nil
}

func (service *sessionService) Exit(ctx context.Context, sessionExit entities.SessionRequest) error {
	session, err := service.repository.Get(ctx, sessionExit.RoomID)
	if err != nil {
		return err
	}

	if session.RoomID == str.Empty {
		err = errors.New("session is empty, user was not inside room")
		service.logs.Error(str.ErrorConcat(err, serviceName, "Set"))
		return err
	}

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

func (service *sessionService) SaveMessage(ctx context.Context, user entities.SessionUser, roomID, message string) error {
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
		SessionUser: user,
		Type:        entities.UserMessageEventType,
		CreatedAt:   time.Now().UTC(),
		Content:     message,
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

	messageCounter := 0
	for _, event := range session.Events {
		if messageCounter == messageMaxLimit {
			break
		}
		if event.Type == entities.UserMessageEventType {
			message := entities.ChatMessage{
				SessionUser: event.SessionUser,
				CreatedAt:   event.CreatedAt,
				Content:     event.Content,
			}
			messages = append(messages, message)
			messageCounter++
		}
	}

	return messages, nil
}
