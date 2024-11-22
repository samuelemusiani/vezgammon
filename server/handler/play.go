package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
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
	// Placeholder, need to implement for matchmaking
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
	// Placeholder, need to implement for matchmaking
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
	userID := c.MustGet("user_id").(int64)

	_, err := db.GetCurrentGame(userID)
	if err != sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, "Already in a game")
		return
	}

	startdices_p1 := types.NewDices()
	startdices_p2 := types.NewDices()

	var startPlayer string
	if startdices_p1.Sum() >= startdices_p2.Sum() {
		startPlayer = types.GameCurrentPlayerP1
	} else {
		startPlayer = types.GameCurrentPlayerP2
	}

	firstdices := types.NewDices()

	g := types.Game{
		Player1:       userID,
		Player2:       userID,
		Start:         time.Now(),
		Status:        types.GameStatusOpen,
		CurrentPlayer: startPlayer,
		Dices:         firstdices,
	}

	_, err = db.CreateGame(g)
	if err != nil {
		slog.With("error", err).Error("GetMoves failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	newgame, err := db.GetCurrentGame(userID)
	if err != nil {
		slog.With("error", err).Error("GetMoves failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	ng := types.NewGame{
		DicesP1: startdices_p1,
		DicesP2: startdices_p2,
		Game:    *newgame,
	}

	c.JSON(http.StatusCreated, ng)
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
	userId := c.MustGet("user_id").(int64)

	retgame, err := db.GetCurrentGame(userId)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, "Game not found")
		return
	}

	if err != nil {
		slog.With("error", err).Error("GetMoves failed")
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
	userId := c.MustGet("user_id").(int64)

	rg, err := db.GetCurrentGame(userId)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, "Not in a game")
		return
	}
	if err != nil {
		slog.With("error", err).Error("GetMoves failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	g, err := db.GetGame(rg.ID)
	if err != nil {
		slog.With("error", err).Error("GetMoves failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var status string
	if g.Player1 == userId { // Player 1 surrended, player 2 wins
		status = types.GameStatusWinP2
	} else {
		status = types.GameStatusWinP1
	}

	g.Status = status
	err = db.UpdateGame(g)
	if err != nil {
		slog.With("error", err).Error("GetMoves failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, "Surrended")

	// Send notification to the other player that the game is over
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
	userId := c.MustGet("user_id").(int64)

	rg, err := db.GetCurrentGame(userId)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, "Not in a game")
		return
	}
	if err != nil {
		slog.With("error", err).Error("GetMoves failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	g, err := db.GetGame(rg.ID)
	if err != nil {
		slog.With("error", err).Error("GetMoves failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if g.WantToDouble {
		c.JSON(http.StatusBadRequest, "Double requested")
		return
	}

	idCurrentPlayer, err := getCurrentPlayer(g.CurrentPlayer, g.Player1, g.Player2)

	if idCurrentPlayer != userId {
		c.JSON(http.StatusBadRequest, "Not your turn")
		return
	}

	var futureturn types.FutureTurn
	futureturn.Dices = g.Dices
	futureturn.CanDouble = (g.CurrentPlayer == g.DoubleOwner)
	futureturn.PossibleMoves, err = bgweb.GetLegalMoves(g)
	if err != nil {
		slog.With("error", err).Error("GetMoves failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	// slog.With("turn", futureturn).Debug("turn")

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
	userId := c.MustGet("user_id").(int64)

	// get moves from body
	buff, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.With("error", err).Error("GetMoves failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var moves []types.Move
	err = json.Unmarshal(buff, &moves)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	rg, err := db.GetCurrentGame(userId)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, "Not in a game")
		return
	}
	if err != nil {
		slog.With("error", err).Error("GetMoves failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	g, err := db.GetGame(rg.ID)
	if err != nil {
		slog.With("error", err).Error("GetMoves failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if g.WantToDouble {
		c.JSON(http.StatusBadRequest, "Double requested")
		return
	}

	idCurrentPlayer, err := getCurrentPlayer(g.CurrentPlayer, g.Player1, g.Player2)

	if idCurrentPlayer != userId {
		c.JSON(http.StatusBadRequest, "Not your turn")
		return
	}

	// check if moves are legal
	legalmoves, err := bgweb.GetLegalMoves(g)
	if err != nil {
		slog.With("error", err).Error("GetMoves failed")
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
		User:   userId,
		Time:   time.Now(),
		Dices:  g.Dices,
		Double: false,
		Moves:  moves,
	}

	_, err = db.CreateTurn(turn)
	if err != nil {
		slog.With("error", err).Error("GetMoves failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Check if we are playing against a bot

	botLevel := db.GetBotLevel(g.Player2)
	// Against a bot
	if botLevel > 0 {
		var t *types.Turn
		var err error

		switch botLevel {
		case 1:
			t, err = bgweb.GetEasyMove(g)
		case 2:
			t, err = bgweb.GetMediumMove(g)
		case 3:
			t, err = bgweb.GetBestMove(g)
		default:
			slog.Error("Invalid bot level")
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		g.PlayMove(t.Moves)
		err = db.UpdateGame(g)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		_, err = db.CreateTurn(*t)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		// Send notification to the other player that it's his turn
	}

	c.JSON(http.StatusCreated, "Moves played")

	// Send notification to the other player that it's his turn
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
	userId := c.MustGet("user_id").(int64)

	rg, err := db.GetCurrentGame(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, "Not in a game")
		} else {
			slog.With("error", err).Error("GetMoves failed")
			c.JSON(http.StatusInternalServerError, err)
		}
		return
	}

	g, err := db.GetGame(rg.ID)
	if err != nil {
		slog.With("error", err).Error("GetMoves failed")
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
		slog.With("error", err).Error("GetMoves failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	botLevel := db.GetBotLevel(g.Player2)
	if botLevel > 0 {
		// Always accept the double
		slog.Debug("Bot always accept double")
		err = acceptDouble(g.Player2)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusCreated, g.DoubleValue*2)

	// Send notification to the other player that he can refuse or accept the double
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
	userId := c.MustGet("user_id").(int64)

	rg, err := db.GetCurrentGame(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, "Not in a game")
		} else {
			slog.With("error", err).Error("GetMoves failed")
			c.JSON(http.StatusInternalServerError, err)
		}
		return
	}

	g, err := db.GetGame(rg.ID)
	if err != nil {
		slog.With("error", err).Error("GetMoves failed")
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
	userId := c.MustGet("user_id").(int64)

	err := acceptDouble(userId)
	if err != nil {

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, "Not in a game")
		} else if strings.Contains(err.Error(), "Internal server error") {
			slog.With("error", err).Error("GetMoves failed")
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusBadRequest, err)
		}
		return
	}

	c.JSON(http.StatusCreated, "Double accepted")
	// Send notification to the other player that he accepts the doubleS
}

// @Summary Create a game against an easy bot
// @Schemes
// @Description Create a game against an easy bot
// @Tags play
// @Accept json
// @Produce json
// @Success 201 {object} types.NewGame
// @Failure 400 "Not in a game or double not possible"
// @Router /play/bot/easy [get]
func PlayEasyBot(c *gin.Context) {
	PlayBot("easy", c)
}

// @Summary Create a game against an medium bot
// @Schemes
// @Description Create a game against an medium bot
// @Tags play
// @Accept json
// @Produce json
// @Success 201 {object} types.NewGame
// @Failure 400 "Not in a game or double not possible"
// @Router /play/bot/medium [get]
func PlayMediumBot(c *gin.Context) {
	PlayBot("medium", c)
}

// @Summary Create a game against an hard bot
// @Schemes
// @Description Create a game against an hard bot
// @Tags play
// @Accept json
// @Produce json
// @Success 201 {object} types.NewGame
// @Failure 400 "Not in a game or double not possible"
// @Router /play/bot/hard [get]
func PlayHardBot(c *gin.Context) {
	PlayBot("hard", c)
}

func PlayBot(mod string, c *gin.Context) {
	userId := c.MustGet("user_id").(int64)

	var botId int64
	switch mod {
	case "easy":
		botId = db.GetEasyBotID()
	case "medium":
		botId = db.GetMediumBotID()
	case "hard":
		botId = db.GetHardBotID()
	default:
		slog.Error("Invalid mod on play bot")
		c.JSON(http.StatusInternalServerError, "Invalid bot")
		return
	}

	_, err := db.GetCurrentGame(userId)
	if err != sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, "Already in a game")
		return
	}

	var startdicesP1, startdicesP2 types.Dices
	for {
		startdicesP1 = types.NewDices()
		startdicesP2 = types.NewDices()

		if startdicesP1.Sum() != startdicesP2.Sum() {
			if startdicesP1.Sum() < startdicesP2.Sum() {
				startdicesP1, startdicesP2 = startdicesP2, startdicesP1
			}
			break
		}
	}

	// Against a bot the player will always start first
	var startPlayer = types.GameCurrentPlayerP1

	firstdices := types.NewDices()

	g := types.Game{
		Player1:       userId,
		Player2:       botId,
		Start:         time.Now(),
		Status:        types.GameStatusOpen,
		CurrentPlayer: startPlayer,
		Dices:         firstdices,
	}

	slog.With("game", g).Debug("Creating game")

	_, err = db.CreateGame(g)
	if err != nil {
		slog.With("error", err).Error("GetMoves failed")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	newgame, err := db.GetCurrentGame(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	ng := types.NewGame{
		DicesP1: startdicesP1,
		DicesP2: startdicesP2,
		Game:    *newgame,
	}

	c.JSON(http.StatusCreated, ng)
}

func acceptDouble(user_id int64) error {
	rg, err := db.GetCurrentGame(user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Not in a game")
		} else {
			return errors.New("Internal server error")
		}
	}

	g, err := db.GetGame(rg.ID)
	if err != nil {
		return err
	}

	if !g.WantToDouble {
		return errors.New("Double not possible")
	}

	g.WantToDouble = false
	g.DoubleValue = 2 * g.DoubleValue

	err = db.UpdateGame(g)
	if err != nil {
		slog.With("error", err).Error("GetMoves failed")
		return errors.New("Internal server error")
	}

	// save turn as double
	doublingplayer_id, err := getCurrentPlayer(g.DoubleOwner, g.Player1, g.Player2)
	if err != nil {
		slog.With("error", err).Error("GetMoves failed")
		return errors.New("Internal server error")
	}

	turn := types.Turn{
		GameId: g.ID,
		User:   doublingplayer_id,
		Time:   time.Now(),
		Double: true,
	}

	_, err = db.CreateTurn(turn)
	if err != nil {
		slog.With("error", err).Error("GetMoves failed")
		return errors.New("Internal server error")
	}
	return nil
}
