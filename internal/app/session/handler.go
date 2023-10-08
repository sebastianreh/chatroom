package session

import (
	"github.com/labstack/echo/v4"
	"github.com/sebastianreh/chatroom/cmd/httpserver/resterror"
	"github.com/sebastianreh/chatroom/internal/config"
	"github.com/sebastianreh/chatroom/internal/entities"
	"github.com/sebastianreh/chatroom/pkg/logger"
	str "github.com/sebastianreh/chatroom/pkg/strings"
	"net/http"
)

const handlerName = "session.handler"

type SessionHandler interface {
	Join(c echo.Context) error
	GetMessages(c echo.Context) error
}

type sessionHandler struct {
	config  config.Config
	service SessionService
	logs    logger.Logger
}

func NewSessionHandler(cfg config.Config, service SessionService, logger logger.Logger) SessionHandler {
	return &sessionHandler{
		config:  cfg,
		service: service,
		logs:    logger,
	}
}

func (handler *sessionHandler) Join(ctx echo.Context) error {
	request := new(entities.SessionRequest)
	if err := ctx.Bind(request); err != nil {
		handler.logs.Error(str.ErrorConcat(err, handlerName, "Join"))
		ctx.Error(err)
		return nil
	}

	err := handler.service.Join(ctx.Request().Context(), *request)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusOK)
}

func (handler *sessionHandler) GetMessages(ctx echo.Context) error {
	roomID := ctx.Param("room_id")
	if roomID == str.Empty {
		err := resterror.NewBadRequestError("empty room id")
		handler.logs.Error(str.ErrorConcat(err, handlerName, "Join"))
		ctx.Error(err)
		return nil
	}

	messages, err := handler.service.GetMessages(ctx.Request().Context(), roomID)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.JSON(http.StatusOK, messages)
}
