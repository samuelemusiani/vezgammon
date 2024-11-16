package ws

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]bool)
var users = make(map[int64]*websocket.Conn)

//autenticazione

func WSHandler(w http.ResponseWriter, r *http.Request, user_id int64) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Add client connection to clients connention array
	clients[conn] = true
	users[user_id] = conn

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
