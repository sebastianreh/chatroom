package websocket

import (
	"fmt"
	ws "github.com/gorilla/websocket"
	"github.com/sebastianreh/chatroom-bots/stocks/internal/config"
	"github.com/sebastianreh/chatroom/pkg/logger"
)

type Websocket interface {
	ReadMessage() (string, error)
}

type websocket struct {
	logger.Logger
	socket *ws.Conn
}

func NewWebsocket(logs logger.Logger, cfg config.Config, botName, roomID string) Websocket {
	url := fmt.Sprintf("%s?bot_name=%s&room_id%s", cfg.Websocket.Endpoint, botName, roomID)
	socket, err := getSocket(url)
	if err != nil {
		logs.Fatal(err.Error())
		panic(err)
	}

	return websocket{
		Logger: logs,
		socket: socket,
	}
}

func getSocket(url string) (*ws.Conn, error) {
	conn, _, err := ws.DefaultDialer.Dial(url, nil)
	if err != nil {
		return conn, err
	}

	return conn, nil
}

func (w websocket) ReadMessage() (string, error) {
	_, msgBytes, err := w.socket.ReadMessage()
	if err != nil {
		w.Logger.Error("error reading message", "ReadMessage", err.Error())
	}
	return string(msgBytes), nil
}
