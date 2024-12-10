package ws

import (
	"errors"
	"log"
	"log/slog"
	"net/http"
	"sync"

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

// var clients = make(map[*websocket.Conn]bool)
// var users = make(map[int64]*websocket.Conn)
// var disconnect = make(map[int64]func(int64) error)
var clients sync.Map
var users sync.Map
var disconnect sync.Map

type Websocket struct{}

func GetWebsocket() Websocket {
	return Websocket{}
}

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
	clients.Store(conn, true)
	users.Store(user_id, conn)
	SendMessage(user_id, Message{Type: "connection_esatablished"})
	slog.Debug("test")
}

// Send messsage to user
func SendMessage(user_id int64, message Message) error {
	value, ok := users.Load(user_id)
	if !ok {
		slog.Debug("SendMessage error: connection not found for user",
			"user_id", user_id,
			"error", ErrConnNotFoundForUser)
		return ErrConnNotFoundForUser
	}
	conn := value.(*websocket.Conn)

	value, ok = clients.Load(conn)
	if !ok {
		slog.Debug("SendMessage error: connection not found in clients map",
			"user_id", user_id,
			"error", ErrConnNotFound)
		return ErrConnNotFound
	}
	active := value.(bool)

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

func AddDisconnectHandler(user_id int64, f func(int64) error) {
	disconnect.Store(user_id, f)
}

func SendBotMessage(user_id int64, message string) error {
	return SendMessage(user_id, Message{Type: "chat_message", Payload: message})
}

func GameTournamentReady(user_id int64) error {
	return SendMessage(user_id, Message{Type: "game_tournament_ready"})
}

func TournamentEnded(user_id int64) error {
	return SendMessage(user_id, Message{Type: "tournament_ended"})
}

func TournamentCancelled(user_id int64) error {
	return SendMessage(user_id, Message{Type: "tournament_cancelled"})
}

func TournamentNewUserEnrolled(user_id int64) error {
	return SendMessage(user_id, Message{Type: "tournament_new_user_enrolled"})
}

func TournamentUserLeft(user_id int64) error {
	return SendMessage(user_id, Message{Type: "tournament_user_left"})
}

func (ws Websocket) SendGameFound(user_id int64) error {
	return SendGameFound(user_id)
}
