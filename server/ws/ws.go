package ws

import (
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

func WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Aggiungi il nuovo client alla lista dei clients connessi
	clients[conn] = true

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			delete(clients, conn)
			return
		}

		fmt.Println("Received message:", string(message))

		// Invia il messaggio a tutti i client connessi
	}
}

func SendMessage(message any) error {
	for client := range clients {
		err := client.WriteJSON(message)
		if err != nil {
			log.Println(err)
			delete(clients, client)
			return err
		}
	}
	return nil
}

func searchGame() (interface{}, error) {
	// Logica per trovare un game
	game := map[string]interface{}{
		"id":   "123",
		"name": "Example Game",
	}
	return game, nil
}
