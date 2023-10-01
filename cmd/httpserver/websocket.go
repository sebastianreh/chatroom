package httpserver

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"sync"
)

var (
	upgrader         = websocket.Upgrader{}
	connections      = make(map[*websocket.Conn]struct{})
	connectionsMutex = sync.Mutex{}
)

func Hello(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	// Add the new connection to the pool
	connectionsMutex.Lock()
	connections[ws] = struct{}{}
	connectionsMutex.Unlock()

	for {
		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Printf("Received message: %s\n", msg)
		broadcastMessage(msg)
	}

	// Remove the closed connection from the pool
	connectionsMutex.Lock()
	delete(connections, ws)
	connectionsMutex.Unlock()

	return nil
}

// broadcastMessage sends a message to all connected clients in the pool
func broadcastMessage(message []byte) {
	connectionsMutex.Lock()
	defer connectionsMutex.Unlock()

	for conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println(err)
		}
	}
}
