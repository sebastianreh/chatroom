package session

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sebastianreh/chatroom/cmd/httpserver/resterror"
	"github.com/sebastianreh/chatroom/internal/config"
	"github.com/sebastianreh/chatroom/internal/entities"
	"github.com/sebastianreh/chatroom/pkg/kafka"
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
	HandleConnection(c echo.Context) error
}

type sessionHandler struct {
	config      config.Config
	websocket   ws.Websocket
	listener    kafka.Consumer
	service     SessionService
	logs        logger.Logger
	connections map[string]*ws.Websocket
}

func NewSessionHandler(cfg config.Config, service SessionService, websocket ws.Websocket, listener kafka.Consumer, logger logger.Logger) SessionHandler {
	return &sessionHandler{
		config:    cfg,
		websocket: websocket,
		listener:  listener,
		service:   service,
		logs:      logger,
	}
}

func (handler *sessionHandler) Listen() {
	err := handler.listener.Listen(handler.ReadStockMessage)
	if err != nil {
		handler.logs.Error(str.ErrorConcat(err, handlerName, "Listen"))
		return
	}
}

func (handler *sessionHandler) ReadStockMessage(context.Context, []byte) {
	err := json.Unmarshal(message, &chargebackRequest)

}

func (handler *sessionHandler) Join(ctx echo.Context) error {
	var joinResponse entities.JoinResponse
	request := new(entities.SessionRequest)
	if err := ctx.Bind(request); err != nil {
		handler.logs.Error(str.ErrorConcat(err, handlerName, "Join"))
		ctx.Error(err)
		return nil
	}

	joinResponse, err := handler.service.Join(ctx.Request().Context(), *request)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	joinAction := entities.GetJoinAction(request.SessionUser)
	err = handler.websocket.BroadCastMessage(joinAction.ToBytes(), request.RoomID)
	if err != nil {
		handler.logs.Error(str.ErrorConcat(err, handlerName, "Join"))
		ctx.Error(err)
	}

	return ctx.JSON(http.StatusOK, joinResponse)
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

	exitAction := entities.GetJoinAction(request.SessionUser)
	err = handler.websocket.BroadCastMessage(exitAction.ToBytes(), request.RoomID)
	if err != nil {
		handler.logs.Error(str.ErrorConcat(err, handlerName, "Join"))
		ctx.Error(err)
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

func (handler *sessionHandler) HandleConnection(ctx echo.Context) error {
	var sessionRequest entities.SessionRequest
	if err := ctx.Bind(&sessionRequest); err != nil {
		handler.logs.Error(str.ErrorConcat(err, handlerName, "HandleConnection"))
		ctx.Error(err)
		return nil
	}

	socket, err := handler.websocket.GetSocket(ctx.Response(), ctx.Request(), sessionRequest.RoomID, sessionRequest.UserID)
	if err != nil {
		handler.logs.Error(str.ErrorConcat(err, handlerName, "HandleConnection"))
		ctx.Error(err)
		return nil
	}

	for {
		_, msg, err := socket.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		var decodedMessage = new(entities.ChatMessage)
		err = json.Unmarshal(msg, decodedMessage)
		if err != nil {
			handler.logs.Error(str.ErrorConcat(err, handlerName, "HandleConnection"))
			continue
		}

		fmt.Printf("Received message: %s\n", msg)
		err = handler.websocket.BroadCastMessage(msg, sessionRequest.RoomID)
		if err != nil {
			handler.logs.Error(str.ErrorConcat(err, handlerName, "HandleConnection"))
			continue
		}

		err = handler.service.SaveMessage(ctx.Request().Context(), *decodedMessage, sessionRequest.RoomID)
		if err != nil {
			err = handler.websocket.CloseSocket(sessionRequest.RoomID, sessionRequest.UserID)
			if err != nil {
				handler.logs.Error(str.ErrorConcat(err, handlerName, "HandleConnection"))
			}
			return err
		}
	}

	return nil
}
