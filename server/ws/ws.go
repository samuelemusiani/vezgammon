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

func WSHandler(w http.ResponseWriter, r *http.Request, userID int64) {
	slog.Info("Starting WebSocket connection", "user_id", userID)
	slog.Debug("Request headers", "headers", r.Header)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go chat(conn, userID)

	// Add client connection to clients connention array
	clients.Store(conn, true)
	users.Store(userID, conn)
	SendMessage(userID, Message{Type: "connection_esatablished"})
	slog.Debug("test")
}

// Send messsage to user
func SendMessage(userID int64, message Message) error {
	value, ok := users.Load(userID)
	if !ok {
		slog.Debug("SendMessage error: connection not found for user",
			"user_id", userID,
			"error", ErrConnNotFoundForUser)
		return ErrConnNotFoundForUser
	}
	conn := value.(*websocket.Conn)

	value, ok = clients.Load(conn)
	if !ok {
		slog.Debug("SendMessage error: connection not found in clients map",
			"user_id", userID,
			"error", ErrConnNotFound)
		return ErrConnNotFound
	}
	active := value.(bool)

	if !active {
		slog.Debug("SendMessage error: inactive connection",
			"user_id", userID,
			"error", ErrConnNotFound)
		return ErrConnNotFound
	}

	err := conn.WriteJSON(message)
	if err != nil {
		slog.Debug("SendMessage error: failed to write JSON",
			"user_id", userID,
			"error", err)
		return err
	}

	slog.Debug("SendMessage: message sent successfully",
		"user_id", userID)
	return nil
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func SendGameFound(userID int64) error {
	return SendMessage(userID, Message{Type: "game_found"})
}

func GameNotFound(userID int64) error {
	return SendMessage(userID, Message{Type: "game_not_found"})
}

func TurnMade(userID int64) error {
	return SendMessage(userID, Message{Type: "turn_made"})
}

func WantToDouble(userID int64) error {
	return SendMessage(userID, Message{Type: "want_to_double"})
}

func DoubleAccepted(userID int64) error {
	return SendMessage(userID, Message{Type: "double_accepted"})
}

func GameEnd(userID int64) error {
	return SendMessage(userID, Message{Type: "game_end"})
}

func AddDisconnectHandler(userID int64, f func(int64) error) {
	disconnect.Store(userID, f)
}

func SendBotMessage(userID int64, message string) error {
	return SendMessage(userID, Message{Type: "chat_message", Payload: message})
}

func GameTournamentReady(userID int64) error {
	return SendMessage(userID, Message{Type: "game_tournament_ready"})
}

func TournamentEnded(userID int64) error {
	return SendMessage(userID, Message{Type: "tournament_ended"})
}

func TournamentCancelled(userID int64) error {
	return SendMessage(userID, Message{Type: "tournament_cancelled"})
}

func TournamentNewUserEnrolled(userID int64) error {
	return SendMessage(userID, Message{Type: "tournament_new_user_enrolled"})
}

func TournamentNewBotEnrolled(userID int64) error {
    return SendMessage(userID, Message{Type: "tournament_new_bot_enrolled"})
}

func TournamentUserLeft(userID int64) error {
	return SendMessage(userID, Message{Type: "tournament_user_left"})
}

func TournamentBotLeft (userID int64) error {
    return SendMessage(userID, Message{Type: "tournament_bot_left"})
}

func (ws Websocket) SendGameFound(userID int64) error {
	return SendGameFound(userID)
}
