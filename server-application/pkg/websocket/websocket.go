package websocket

import (
	"errors"
	"fmt"
	ws "github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type Websocket interface {
	GetSocket(responseWriter http.ResponseWriter, request *http.Request, groupID, userID string) (*ws.Conn, error)
	CloseSocket(groupID, userID string) error
	BroadCastMessage(message []byte, groupID string) error
}

type websocket struct {
	upgrader    ws.Upgrader
	connections map[string]map[string]*ws.Conn
	connMutex   sync.Mutex
}

func NewWebsocket() Websocket {
	return &websocket{
		upgrader: ws.Upgrader{
			CheckOrigin: GetCheckFunc(),
		},
		connections: make(map[string]map[string]*ws.Conn),
		connMutex:   sync.Mutex{},
	}
}

func (w *websocket) GetSocket(responseWriter http.ResponseWriter, request *http.Request, groupID, userID string) (*ws.Conn, error) {
	socket, ok := w.connections[userID][groupID]
	if ok {
		return socket, nil
	}

	socket, err := w.upgrader.Upgrade(responseWriter, request, nil)
	if err != nil {
		return nil, err
	}

	w.connMutex.Lock()
	_, ok = w.connections[groupID]
	if !ok {
		usersMap := make(map[string]*ws.Conn)
		usersMap[userID] = socket
		w.connections[groupID] = usersMap
	} else {
		w.connections[groupID][userID] = socket
	}
	w.connMutex.Unlock()

	return socket, nil
}

func (w *websocket) CloseSocket(groupID, userID string) error {
	w.connMutex.Lock()
	defer w.connMutex.Unlock()
	socket, ok := w.connections[groupID][userID]
	if !ok {
		return errors.New(fmt.Sprintf("no sockets found for groupID %s", groupID))
	}
	err := socket.Close()
	if err != nil {
		return err
	}
	delete(w.connections[groupID], userID)

	return nil
}

func (w *websocket) BroadCastMessage(message []byte, groupID string) error {
	w.connMutex.Lock()
	defer w.connMutex.Unlock()

	socketsList, ok := w.connections[groupID]
	if !ok {
		return nil
	}

	for _, socket := range socketsList {
		if err := socket.WriteMessage(ws.TextMessage, message); err != nil {
			return err
		}
	}

	return nil
}

// Here we should implement validation with JWT
func GetCheckFunc() func(r *http.Request) bool {
	return func(r *http.Request) bool {
		return true
	}
}
