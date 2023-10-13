package session

import (
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
	"strings"
)

const handlerName = "session.handler"

type SessionHandler interface {
	Join(c echo.Context) error
	Exit(c echo.Context) error
	GetMessages(c echo.Context) error
	HandleChatConnection(c echo.Context) error
	HandleBotConnection(c echo.Context) error
	Listen()
}

type sessionHandler struct {
	config    config.Config
	websocket ws.Websocket
	listener  kafka.Consumer
	service   SessionService
	logs      logger.Logger
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

func (handler *sessionHandler) ReadStockMessage(message []byte) {
	var stock entities.StockMessage
	err := json.Unmarshal(message, &stock)
	if err != nil {
		handler.logs.Error(str.ErrorConcat(err, handlerName, "ReadStockMessage"))
		return
	}

	handler.logs.Info("New message received", handlerName, "Listen")

	err = handler.websocket.BroadCastMessage(message, stock.RoomID)
	if err != nil {
		return
	}
}

func (handler *sessionHandler) Join(ctx echo.Context) error {
	var joinResponse entities.JoinResponse
	request := new(entities.SessionChatRequest)
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

	}

	return ctx.JSON(http.StatusOK, joinResponse)
}

func (handler *sessionHandler) Exit(ctx echo.Context) error {
	request := new(entities.SessionChatRequest)
	if err := ctx.Bind(request); err != nil {
		handler.logs.Error(str.ErrorConcat(err, handlerName, "Exit"))
		ctx.Error(err)
		return nil
	}

	err := handler.websocket.CloseSocket(request.RoomID, request.UserID)
	if err != nil {
		err = resterror.NewBadRequestError(err.Error())
		ctx.Error(err)
		return nil
	}

	exitAction := entities.GetExitAction(request.SessionUser)
	err = handler.websocket.BroadCastMessage(exitAction.ToBytes(), request.RoomID)
	if err != nil {
		ctx.Error(err)
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
		handler.logs.Error(str.ErrorConcat(err, handlerName, "GetMessages"))
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

func (handler *sessionHandler) HandleChatConnection(ctx echo.Context) error {
	var sessionChatRequest entities.SessionChatRequest
	if err := ctx.Bind(&sessionChatRequest); err != nil {
		handler.logs.Error(str.ErrorConcat(err, handlerName, "HandleChatConnection"))
		ctx.Error(err)
		return nil
	}

	socket, err := handler.websocket.GetSocket(ctx.Response(), ctx.Request(), sessionChatRequest.RoomID, sessionChatRequest.UserID)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	messageChan := make(chan []byte)

	go func() {
		for msg := range messageChan {

			var decodedMessage = new(entities.ChatMessage)
			err := json.Unmarshal(msg, decodedMessage)
			if err != nil {
				handler.logs.Error(str.ErrorConcat(err, handlerName, "HandleChatConnection"))
				continue
			}

			if strings.HasPrefix(decodedMessage.Content, "/") {
				err = handler.websocket.BroadCastMessage(msg, sessionChatRequest.RoomID)
				if err != nil {
					continue
				}
			}

			fmt.Printf("Received message: %s\n", msg)
			err = handler.websocket.BroadCastMessage(msg, sessionChatRequest.RoomID)
			if err != nil {
				continue
			}

			err = handler.service.SaveMessage(ctx.Request().Context(), *decodedMessage, sessionChatRequest.RoomID)
			if err != nil {
				err = handler.websocket.CloseSocket(sessionChatRequest.RoomID, sessionChatRequest.UserID)
				if err != nil {
					continue
				}
			}
		}
	}()

	for {
		if socket != nil {
			msgType, msg, err := socket.ReadMessage()
			if msgType == -1 {
				break
			}

			if err != nil {
				handler.logs.Error(str.ErrorConcat(err, handlerName, "HandleChatConnection"))
				continue
			}
			messageChan <- msg
		}
	}

	return nil
}

func (handler *sessionHandler) HandleBotConnection(ctx echo.Context) error {
	var botSessionRequest entities.BotSessionRequest
	if err := ctx.Bind(&botSessionRequest); err != nil {
		handler.logs.Error(str.ErrorConcat(err, handlerName, "HandleBotConnection"))
		ctx.Error(err)
		return nil
	}

	socket, err := handler.websocket.GetSocket(ctx.Response(), ctx.Request(), botSessionRequest.RoomID, botSessionRequest.BotName)
	if err != nil {
		handler.logs.Error(str.ErrorConcat(err, handlerName, "HandleBotConnection"))
		ctx.Error(err)
		return nil
	}

	messageChan := make(chan []byte)

	go func() {
		for {
			_, msg, err := socket.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}

			var decodedMessage = new(entities.ChatMessage)
			err = json.Unmarshal(msg, decodedMessage)
			if err != nil {
				handler.logs.Error(str.ErrorConcat(err, handlerName, "HandleBotConnection"))
				continue
			}

			if strings.HasPrefix(decodedMessage.Content, fmt.Sprintf("/%s", botSessionRequest.BotName)) {
				contentBytes, err := json.Marshal(decodedMessage.Content)
				if err != nil {
					handler.logs.Error(str.ErrorConcat(err, handlerName, "HandleBotConnection"))
					continue
				}

				err = handler.websocket.SendMessageToSocket(contentBytes, socket)
				if err != nil {
					handler.logs.Error(str.ErrorConcat(err, handlerName, "HandleBotConnection"))
					continue
				}
			}
		}
	}()

	for {
		if socket != nil {
			msgType, msg, err := socket.ReadMessage()
			if msgType == -1 {
				break
			}

			if err != nil {
				handler.logs.Error(str.ErrorConcat(err, handlerName, "HandleChatConnection"))
				continue
			}
			messageChan <- msg
		}
	}

	return nil
}
