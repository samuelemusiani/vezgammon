package ws

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gotest.tools/v3/assert"
)

var router *gin.Engine

func TestSendMessage(t *testing.T) {
	// Create test server
	router = gin.Default()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router.GET("/api/ws", func(c *gin.Context) {
			c.Set("user_id", int64(1))
			WSHandler(c.Writer, c.Request, c.MustGet("user_id").(int64))
		})
		router.ServeHTTP(w, r)
	}))

	AddDisconnectHandler(1, func(id int64) error {
		return nil
	})

	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/ws"

	// Connect websocket client
	headers := http.Header{}

	ws, _, err := websocket.DefaultDialer.Dial(url, headers)
	assert.NilError(t, err)

	// Give the server a moment to register the connection
	time.Sleep(100 * time.Millisecond)

	// Test sending a message
	err = SendMessage(int64(1), Message{Type: "connection_established"})
	assert.NilError(t, err)

	// Verify the message was received
	messageType, _, err := ws.ReadMessage()
	assert.NilError(t, err)
	assert.Equal(t, messageType, websocket.TextMessage)

	SendGameFound(int64(1))
	GameNotFound(int64(1))
	TurnMade(int64(1))
	WantToDouble(int64(1))
	DoubleAccepted(int64(1))
	GameEnd(int64(1))
	SendBotMessage(int64(1), "")
	GameTournamentReady(int64(1))
	TournamentEnded(int64(1))
	TournamentCancelled(int64(1))
	TournamentNewUserEnrolled(int64(1))
	TournamentUserLeft(int64(1))

	// Graceful shutdown
	err = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	assert.NilError(t, err)

	// Close the WebSocket connection
	ws.Close()

	// Close the test server
	server.Close()
}

func TestSendMessageError(t *testing.T) {
	err := SendMessage(int64(249), Message{Type: "connection_established"})
	assert.ErrorIs(t, err, ErrConnNotFoundForUser)

	w := websocket.Conn{}
	users.Store(int64(1), &w)
	err = SendMessage(int64(1), Message{Type: "game_found"})
	assert.ErrorIs(t, err, ErrConnNotFound)

	clients.Store(&w, false)
	SendMessage(int64(1), Message{Type: "game_found"})
	assert.ErrorIs(t, err, ErrConnNotFound)
}
