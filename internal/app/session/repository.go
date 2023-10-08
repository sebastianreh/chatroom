package session

import (
	"context"
	"encoding/json"
	"github.com/sebastianreh/chatroom/internal/config"
	"github.com/sebastianreh/chatroom/internal/entities"
	"github.com/sebastianreh/chatroom/pkg/logger"
	"github.com/sebastianreh/chatroom/pkg/redis"
	str "github.com/sebastianreh/chatroom/pkg/strings"
	"time"
)

const (
	repositoryName  = "session.repository"
	InactiveTimeTTL = time.Duration(24) * time.Hour
)

type SessionRepository interface {
	Set(ctx context.Context, session entities.Session) error
	Get(ctx context.Context, roomID string) (entities.Session, error)
}

type sessionRepository struct {
	config         config.Config
	redis          redis.Redis
	logs           logger.Logger
	collectionName string
}

func NewSessionRepository(cfg config.Config, redis redis.Redis, logger logger.Logger) SessionRepository {
	return &sessionRepository{
		config:         cfg,
		redis:          redis,
		logs:           logger,
		collectionName: cfg.MongoDB.Collections.Rooms,
	}
}

func (repository *sessionRepository) Set(ctx context.Context, session entities.Session) error {
	eventsBytes, err := json.Marshal(session.Events)
	if err != nil {
		repository.logs.Error(str.ErrorConcat(err, repositoryName, "Set"))
	}
	err = repository.redis.Set(ctx, session.RoomID, string(eventsBytes), InactiveTimeTTL)
	if err != nil {
		repository.logs.Error(str.ErrorConcat(err, repositoryName, "Set"))
	}
	return err
}

func (repository *sessionRepository) Get(ctx context.Context, roomID string) (entities.Session, error) {
	var session entities.Session
	var events []entities.Event
	eventsString, err := repository.redis.Get(ctx, roomID)
	if err != nil {
		repository.logs.Error(str.ErrorConcat(err, repositoryName, "Get"))
		return session, err
	}

	if eventsString == str.Empty {
		return session, nil
	}

	sessionBytes := []byte(eventsString)
	err = json.Unmarshal(sessionBytes, &events)
	if err != nil {
		repository.logs.Error(str.ErrorConcat(err, repositoryName, "Get"))
		return session, err
	}

	session.RoomID = roomID
	session.Events = events

	return session, err
}
