package ws

import (
	"errors"
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

var ErrConnNotFound = errors.New("Connection not found")
var ErrConnNotActive = errors.New("Connetion not active")
var ErrConnNotFoundForUser = errors.New("Connetion not found for user")

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

	go chat(conn, user_id)

	// Add client connection to clients connention array
	clients[conn] = true
	users[user_id] = conn
	SendMessage(user_id, Message{Type: "connection_esatablished"})
	slog.Debug("test")
}

// Send messsage to user
func SendMessage(user_id int64, message Message) error {
	conn, ok := users[user_id]
	if !ok {
		slog.Debug("SendMessage error: connection not found for user",
			"user_id", user_id,
			"error", ErrConnNotFoundForUser)
		return ErrConnNotFoundForUser
	}

	active, ok := clients[conn]
	if !ok {
		slog.Debug("SendMessage error: connection not found in clients map",
			"user_id", user_id,
			"error", ErrConnNotFound)
		return ErrConnNotFound
	}

	if !active {
		slog.Debug("SendMessage error: inactive connection",
			"user_id", user_id,
			"error", ErrConnNotFound)
		return ErrConnNotFound
	}

	err := conn.WriteJSON(message)
	if err != nil {
		slog.Debug("SendMessage error: failed to write JSON",
			"user_id", user_id,
			"error", err)
		return err
	}

	slog.Debug("SendMessage: message sent successfully",
		"user_id", user_id)
	return nil
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func SendGameFound(user_id int64) error {
	return SendMessage(user_id, Message{Type: "game_found"})
}

func GameNotFound(user_id int64) error {
	return SendMessage(user_id, Message{Type: "game_not_found"})
}

func TurnMade(user_id int64) error {
	return SendMessage(user_id, Message{Type: "turn_made"})
}

func WantToDouble(user_id int64) error {
	return SendMessage(user_id, Message{Type: "want_to_double"})
}

func DoubleAccepted(user_id int64) error {
	return SendMessage(user_id, Message{Type: "double_accepted"})
}

func GameEnd(user_id int64) error {
	return SendMessage(user_id, Message{Type: "game_end"})
}

func GameTournamentReady(user_id int64) error {
	return SendMessage(user_id, Message{Type: "game_tournament_ready"})
}

func TournamentEnded(user_id int64) error {
	return SendMessage(user_id, Message{Type: "tournament_ended"})
}

func TournamentCanceled(user_id int64) error {
	return SendMessage(user_id, Message{Type: "tournament_canceled"})
}
