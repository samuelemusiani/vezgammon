package handler

import (
	"database/sql"
	"net/http"
	"vezgammon/server/bgweb"
	"vezgammon/server/db"
	"vezgammon/server/types"

	"github.com/gin-gonic/gin"
)

// @Summary Start a matchmaking search for a new game
// @Schemes
// @Description Start a matchmaking search for a new game
// @Tags play
// @Accept json
// @Produce json
// @Success 201 "Search started"
// @Failure 400 "Already searching or in a game"
// @Router /play/search [get]
func StartPlaySearch(c *gin.Context) {
}

// @Summary Stop a running matchmaking search
// @Schemes
// @Description Stop a running matchmaking search
// @Tags play
// @Accept json
// @Produce json
// @Success 204 "Search stopped"
// @Failure 400 "Not searching"
// @Router /play/search [delete]
func StopPlaySearch(c *gin.Context) {
}

// @Summary Get current game
// @Schemes
// @Description Get current game
// @Tags play
// @Accept json
// @Produce json
// @Success 200 {object} types.ReturnGame
// @Failure 404 "Game not found"
// @Router /play [get]
func GetCurrentGame(c *gin.Context) {
	user_id := c.MustGet("user_id").(int64)

	retgame, err := db.GetCurrentGame(user_id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, "Game not found")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, retgame)
}

// @Summary Surrend to current game
// @Schemes
// @Description Surrend to current game
// @Tags play
// @Accept json
// @Produce json
// @Success 201 "Surrended"
// @Failure 404 "Not in a game"
// @Router /play [delete]
func SurrendToCurrentGame(c *gin.Context) {
	user_id := c.MustGet("user_id").(int64)

	rg, err := db.GetCurrentGame(user_id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, "Not in a game")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	g, err := db.GetGame(rg.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var status string
	if g.Player1 == user_id { // Player 1 surrended, player 2 wins
		status = types.GameStatusWinP2
	} else {
		status = types.GameStatusWinP1
	}

	g.Status = status
	err = db.UpdateGame(*g)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, "Surrended")
}

// @Summary Get possible moves for next turn
// @Schemes
// @Description Get possible moves for next turn
// @Tags play
// @Accept json
// @Produce json
// @Success 200 {object} types.FutureTurn "Dice with all possible moves and the ability to double"
// @Failure 400 "Not in a game, not your turn or double requested"
// @Router /play/moves [get]
func GetPossibleMoves(c *gin.Context) {
	user_id := c.MustGet("user_id").(int64)

	rg, err := db.GetCurrentGame(user_id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, "Not in a game")
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	g, err := db.GetGame(rg.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if g.WantToDouble {
		c.JSON(http.StatusBadRequest, "Double requested")
		return
	}

	var currentplayer string
	if g.Player1 == user_id {
		currentplayer = types.GameCurrentPlayerP1
	} else {
		currentplayer = types.GameCurrentPlayerP1
	}

	if currentplayer != g.CurrentPlayer {
		c.JSON(http.StatusBadRequest, "Not your turn")
		return
	}

	var futureturn types.FutureTurn
	futureturn.Dices = g.Dices
	futureturn.CanDouble = (g.CurrentPlayer == g.DoubleOwner)
	futureturn.PossibleMoves, err = bgweb.GetLegalMoves(g)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, futureturn)
}

// @Summary Play all the moves for the current turn
// @Schemes
// @Description Play all the moves for the current turn
// @Tags play
// @Accept json
// @Param request body []types.Move true "Moves to play"
// @Produce json
// @Success 201 "Moves played"
// @Failure 400 "Moves not legal, not your turn or not in a game"
// @Router /play/moves [post]
func PlayMoves(c *gin.Context) {
}

// @Summary The player want to double
// @Schemes
// @Description The player want to double
// @Tags play
// @Accept json
// @Produce json
// @Success 201 {integer} int "Value of the red dice after the double"
// @Failure 400 "Not in a game or double not possible"
// @Router /play/double [post]
func WantToDouble(c *gin.Context) {
}

// @Summary Refuse the double
// @Schemes
// @Description Refuse the double
// @Tags play
// @Accept json
// @Produce json
// @Success 201 "Double refused"
// @Failure 400 "Not in a game or can't refuse double"
// @Router /play/double [delete]
func RefuseDouble(c *gin.Context) {
}

// @Summary Accept the double
// @Schemes
// @Description Accept the double
// @Tags play
// @Accept json
// @Produce json
// @Success 201 "Double accepted"
// @Failure 400 "Not in a game or double not possible"
// @Router /play/double [put]
func AcceptDouble(c *gin.Context) {
}
