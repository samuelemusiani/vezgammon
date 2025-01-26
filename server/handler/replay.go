package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"vezgammon/server/db"
	"vezgammon/server/types"

	"github.com/gin-gonic/gin"
)

type gameReqPos struct {
	GameID int64 `json:"game_id" example:"1"`
	Move   int64 `json:"move" example:"1"`
}

var defaultCheckers = [25]int8{
	0, 0, 0, 0, 0, 0, 5, 0, 3, 0, 0, 0, 0,
	5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}

// @Summary Get position of a game
// @Schemes
// @Description Get position of a game based on id and move number
// @Accept json
// @Param request body gameReqPos true "game id and move number"
// @Produce json
// @Success 201 {object} types.ReturnReplay "game position"
// @Router /replay [post]
func GetReplay(c *gin.Context) {

	userID := c.MustGet("user_id").(int64)

	buff, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.With("err", err).Error("Reading body")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var tmp gameReqPos
	err = json.Unmarshal(buff, &tmp)
	if err != nil {
		slog.With("err", err).Debug("Bad request unmarshal request body")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	g, err := db.GetGame(tmp.GameID)
	if err != nil {
		slog.With("err", err).Error("Getting game")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if g.Player1 != userID && g.Player2 != userID {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	g.P1Checkers = defaultCheckers
	g.P2Checkers = defaultCheckers

	turns, err := db.GetTurns(tmp.GameID)
	if err != nil {
		slog.With("err", err).Error("GetTurns")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if len(turns) > 0 {
		if turns[0].User == g.Player1 {
			g.CurrentPlayer = types.GameCurrentPlayerP1
		} else {
			g.CurrentPlayer = types.GameCurrentPlayerP2
		}
	}

	if tmp.Move > int64(len(turns)) {
		c.JSON(http.StatusBadRequest, "Move number out of range")
		return
	}

	var dices types.Dices
	g, dices = reconstructGameFromTurns(turns, g, tmp.Move)

	u1, err := db.GetUser(g.Player1)
	if err != nil {
		slog.With("err", err).Error("Getting user")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	u2, err := db.GetUser(g.Player2)
	if err != nil {
		slog.With("err", err).Error("Getting user")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	rg := g.ToReturnGame(u1.Username, u2.Username)

	c.JSON(http.StatusOK, types.ReturnReplay{
		Game:  rg,
		Dices: dices,
	})
}
