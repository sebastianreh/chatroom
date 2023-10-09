package session

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sebastianreh/chatroom/cmd/httpserver/resterror"
	"github.com/sebastianreh/chatroom/internal/config"
	"github.com/sebastianreh/chatroom/internal/entities"
	"github.com/sebastianreh/chatroom/pkg/logger"
	str "github.com/sebastianreh/chatroom/pkg/strings"
	ws "github.com/sebastianreh/chatroom/pkg/websocket"
	"log"
	"net/http"
)

const handlerName = "session.handler"

type SessionHandler interface {
	Join(c echo.Context) error
	Exit(c echo.Context) error
	GetMessages(c echo.Context) error
	OpenConnection(c echo.Context) error
}

type sessionHandler struct {
	config      config.Config
	websocket   ws.Websocket
	service     SessionService
	logs        logger.Logger
	connections map[string]*ws.Websocket
}

func NewSessionHandler(cfg config.Config, service SessionService, websocket ws.Websocket, logger logger.Logger) SessionHandler {
	return &sessionHandler{
		config:    cfg,
		websocket: websocket,
		service:   service,
		logs:      logger,
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

func (handler *sessionHandler) Exit(ctx echo.Context) error {
	request := new(entities.SessionRequest)
	if err := ctx.Bind(request); err != nil {
		handler.logs.Error(str.ErrorConcat(err, handlerName, "Exit"))
		ctx.Error(err)
		return nil
	}

	err := handler.websocket.CloseSocket(request.RoomID, request.UserID)
	if err != nil {
		err = resterror.NewBadRequestError(err.Error())
		handler.logs.Error(str.ErrorConcat(err, handlerName, "Exit"))
		ctx.Error(err)
		return nil
	}

	err = handler.service.Exit(ctx.Request().Context(), *request)
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

func (handler *sessionHandler) OpenConnection(ctx echo.Context) error {
	var sessionRequest entities.SessionRequest
	if err := ctx.Bind(&sessionRequest); err != nil {
		handler.logs.Error(str.ErrorConcat(err, handlerName, "OpenConnection"))
		ctx.Error(err)
		return nil
	}

	socket, err := handler.websocket.GetSocket(ctx.Response(), ctx.Request(), sessionRequest.RoomID, sessionRequest.UserID)
	if err != nil {
		handler.logs.Error(str.ErrorConcat(err, handlerName, "OpenConnection"))
		ctx.Error(err)
		return nil
	}

	for {
		_, msg, err := socket.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		fmt.Printf("Received message: %s\n", msg)
		err = handler.websocket.BroadCastMessage(msg, sessionRequest.RoomID)
		if err != nil {
			handler.logs.Error(str.ErrorConcat(err, handlerName, "OpenConnection"))
			continue
		}

		err = handler.service.SaveMessage(ctx.Request().Context(), sessionRequest.SessionUser, sessionRequest.RoomID, string(msg))
		if err != nil {
			err = handler.websocket.CloseSocket(sessionRequest.RoomID, sessionRequest.UserID)
			if err != nil {
				handler.logs.Error(str.ErrorConcat(err, handlerName, "OpenConnection"))
			}
			return err
		}
	}

	return nil
}
