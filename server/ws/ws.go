package ws

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	//CORS per locale
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var users = make(map[int64]*websocket.Conn)

func WSHandler(w http.ResponseWriter, r *http.Request, user_id int64) {
	slog.Info("Starting WebSocket connection", "user_id", user_id)
	slog.Debug("Request headers", "headers", r.Header)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Add client connection to clients connention array
	clients[conn] = true
	users[user_id] = conn
	SendMessage(user_id, "Connection established")
	slog.Debug("test")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			delete(clients, conn)
			return
		}

		fmt.Println("Received message:", string(message))
	}
}

// Send messsage to user
func SendMessage(user int64, message string) error {
	conn, ok := users[user]
	if !ok {
		return errors.New("Connetion not found for user")
	}

	active, ok := clients[conn]
	if !ok {
		return errors.New("Connetion not found")
	}

	if !active {
		return errors.New("Connetion not active")
	}

	conn.WriteJSON(message)

	return nil
}
