package ws

import (
	"log/slog"
	"vezgammon/server/db"

	"github.com/gorilla/websocket"
)

func chat(conn *websocket.Conn, user_id int64) {
	defer conn.Close()
	var m Message
	for {
		err := conn.ReadJSON(&m)

		if err != nil {
			slog.With("err", err).Error("Error reading message")

			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				slog.Debug("Normal closure")
				f, ok := disconnect[user_id]
				if !ok {
					break
				}
				slog.Debug("Removing player from matchmaking queue")

				err := f(user_id)
				slog.With("err", err).Error("Removing player from matchmaking queue")
			}
			break
		}

		slog.With("message", m).Debug("Received message")

		rg, err := db.GetCurrentGame(user_id)
		if err != nil {
			slog.With("err", err).Error("Error getting current game in chat")
			continue
		}

		g, err := db.GetGame(rg.ID)
		if err != nil {
			slog.With("err", err).Error("Error getting game in chat")
			continue
		}

		var oppoonentId int64
		if g.Player1 == user_id {
			oppoonentId = g.Player2
		} else {
			oppoonentId = g.Player1
		}

		err = SendMessage(oppoonentId, m)

		if err != nil {
			slog.With("err", err).Error("Error writing message")
			break
		}
	}
}
