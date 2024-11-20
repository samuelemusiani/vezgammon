package handler

import (
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"time"
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

// @Summary Create a local game
// @Schemes
// @Description Create a local game for playing locally in the same device
// @Tags play
// @Accept json
// @Produce json
// @Success 201 {object} types.NewGame
// @Failure 400 "Already in a game"
// @Router /play/local [get]
func StartGameLocalcally(c *gin.Context) {
	user_id := c.MustGet("user_id").(int64)

	_, err := db.GetCurrentGame(user_id)
	if err != sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, "Already in a game")
		return
	}

	startdices_p1 := types.NewDices()
	startdices_p2 := types.NewDices()

	var start_player string
	if startdices_p1.Sum() >= startdices_p2.Sum() {
		start_player = types.GameCurrentPlayerP1
	} else {
		start_player = types.GameCurrentPlayerP2
	}

	firstdices := types.NewDices()

	g := types.Game{
		Player1:       user_id,
		Player2:       user_id,
		Start:         time.Now(),
		Status:        types.GameStatusOpen,
		CurrentPlayer: start_player,
		Dices:         firstdices,
	}

	_, err = db.CreateGame(g)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	newgame, err := db.GetCurrentGame(user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	ng := types.NewGame{
		DicesP1: startdices_p1,
		DicesP2: startdices_p2,
		Game:    *newgame,
	}

	c.JSON(http.StatusOK, ng)
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
	err = db.UpdateGame(g)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, "Surrended")

	// TODO: send notification to the other player that the game is over
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

	idCurrentPlayer, err := getCurrentPlayer(g.CurrentPlayer, g.Player1, g.Player2)

	if idCurrentPlayer != user_id {
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
	slog.With("turn", futureturn).Debug("turn")

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
	user_id := c.MustGet("user_id").(int64)

	// get moves from body
	buff, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var moves []types.Move
	err = json.Unmarshal(buff, &moves)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

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

	idCurrentPlayer, err := getCurrentPlayer(g.CurrentPlayer, g.Player1, g.Player2)

	if idCurrentPlayer != user_id {
		c.JSON(http.StatusBadRequest, "Not your turn")
		return
	}

	// check if moves are legal
	legalmoves, err := bgweb.GetLegalMoves(g)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	islegal := len(legalmoves) == 0 // can't make moves
	for _, m := range legalmoves {
		if reflect.DeepEqual(m, moves) {
			islegal = true
			break
		}
	}

	if !islegal {
		c.JSON(http.StatusBadRequest, "Moves not legal")
		return
	}

	g.PlayMove(moves)

	err = db.UpdateGame(g)

	// save turn
	turn := types.Turn{
		GameId: g.ID,
		User:   user_id,
		Time:   time.Now(),
		Dices:  g.Dices,
		Double: false,
		Moves:  moves,
	}

	_, err = db.CreateTurn(turn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, "Moves played")

	// TODO: send notification to the other player that it's his turn
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
	user_id := c.MustGet("user_id").(int64)

	rg, err := db.GetCurrentGame(user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, "Not in a game")
		} else {
			c.JSON(http.StatusInternalServerError, err)
		}
		return
	}

	g, err := db.GetGame(rg.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if g.WantToDouble {
		c.JSON(http.StatusBadRequest, "Double not possible")
		return
	}

	if (g.DoubleOwner != types.GameDoubleOwnerAll && g.DoubleOwner == g.CurrentPlayer) || g.DoubleValue == 64 {
		c.JSON(http.StatusBadRequest, "Double not possible")
		return
	}

	g.WantToDouble = true
	g.DoubleOwner = g.CurrentPlayer

	err = db.UpdateGame(g)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, g.DoubleValue*2)

	// TODO: send notification to the other player that he can refuse or accept the double
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
	user_id := c.MustGet("user_id").(int64)

	rg, err := db.GetCurrentGame(user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, "Not in a game")
		} else {
			c.JSON(http.StatusInternalServerError, err)
		}
		return
	}

	g, err := db.GetGame(rg.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if !g.WantToDouble {
		c.JSON(http.StatusBadRequest, "Can't refuse double")
		return
	}

	// Refuse the double is equal to surrender
	SurrendToCurrentGame(c)
	c.JSON(http.StatusCreated, "Double refused")
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
	user_id := c.MustGet("user_id").(int64)

	rg, err := db.GetCurrentGame(user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, "Not in a game")
		} else {
			c.JSON(http.StatusInternalServerError, err)
		}
		return
	}

	g, err := db.GetGame(rg.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if !g.WantToDouble {
		c.JSON(http.StatusBadRequest, "Double not possible")
		return
	}

	g.WantToDouble = false
	g.DoubleValue = 2 * g.DoubleValue

	err = db.UpdateGame(g)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// save turn as double
	doublingplayer_id, err := getCurrentPlayer(g.DoubleOwner, g.Player1, g.Player2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	turn := types.Turn{
		GameId: g.ID,
		User:   doublingplayer_id,
		Time:   time.Now(),
		Double: true,
	}

	_, err = db.CreateTurn(turn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, "Double accepted")
	// TODO: send notification to the other player that he accepts the doubleS
}
