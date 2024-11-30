package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"reflect"
	"time"
	"vezgammon/server/bgweb"
	"vezgammon/server/matchmaking"
	"vezgammon/server/ws"

	"vezgammon/server/db"
	"vezgammon/server/types"

	"github.com/gin-gonic/gin"
)

var ErrInternal = errors.New("Internal server error")
var ErrNotInGame = errors.New("Not in a game")
var ErrDoubleNotPossible = errors.New("Double not possible")

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
	slog.Debug("Inizio a cercare un game")
	user_id := c.MustGet("user_id").(int64)

	//send to db the user [searching]
	err := matchmaking.SearchGame(user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	ws.AddDisconnectHandler(user_id, matchmaking.StopSearch)

	c.JSON(http.StatusOK, "Search started")
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
	user_id := c.MustGet("user_id").(int64)

	err := matchmaking.StopSearch(user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, "Search stopped")
}

// @Summary Create a game with a link
// @Schemes
// @Description Create a game with a link
// @Tags play
// @Accept json
// @Produce json
// @Success 201 "Link created"
// @Router /play/invite [get]
func StartPlayInviteSearch(c *gin.Context) {
	user_id := c.MustGet("user_id").(int64)

	link, err := matchmaking.GenerateLink(user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	type Link struct {
		Link string
	}

	c.JSON(http.StatusCreated, Link{Link: link})
}

// @Summary Join a game with a link
// @Schemes
// @Description Join a game with a link
// @Tags play
// @Accept json
// @Produce json
// @Param id path string true "Link ID"
// @Success 200 "Link generated"
// @Failure 400 "Already in a game"
// @Failure 404 "Link not found"
// @Router /play/invite/{id} [get]
func PlayInvite(c *gin.Context) {
	user_id := c.MustGet("user_id").(int64)

	uuid := c.Param("id")

	err := matchmaking.JoinLink(uuid, user_id)
	if err != nil {
		slog.With("error", err).Error("Joining link")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, "Link joined")
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

	startdicesP1 := types.NewDices()
	startdicesP2 := types.NewDices()

	var startPlayer string
	if startdicesP1.Sum() >= startdicesP2.Sum() {
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
		slog.With("error", err).Error("Creaiting local game")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	newgame, err := db.GetCurrentGame(userID)
	if err != nil {
		slog.With("error", err).Error("Getting current game on local")
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
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, "Game not found")
			return
		} else {
			slog.With("error", err).Error("Getting current game in /play")
			c.JSON(http.StatusInternalServerError, err)
			return
		}
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
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, ErrNotInGame.Error())
			return
		} else {
			slog.With("error", err).Error("Getting current game in /play [delete]")
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	g, err := db.GetGame(rg.ID)
	if err != nil {
		slog.With("error", err).Error("Getting game from DB in /play [delete]")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var status string
	var opponentID int64
	if g.Player1 == userId { // Player 1 surrended, player 2 wins
		status = types.GameStatusWinP2
		opponentID = g.Player2
	} else {
		status = types.GameStatusWinP1
		opponentID = g.Player1
	}

	g.Status = status
	err = db.UpdateGame(g)
	if err != nil {
		slog.With("error", err).Error("Updating game in /play [delete]")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Send notification to the other player that the game is over
	err = ws.GameEnd(opponentID)
	if err != nil {
		slog.With("error", err).Error("Sending message to player")
	}

	c.JSON(http.StatusCreated, "Surrended")
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
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, ErrNotInGame.Error())
			return
		} else {
			slog.With("error", err).Error("Getting current game in /moves [get]")
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	g, err := db.GetGame(rg.ID)
	if err != nil {
		slog.With("error", err).Error("Getting game from DB in /moves [get]")
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
		slog.With("error", err).Error("Getting legal moves in /moves [get]")
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
		slog.With("error", err).Error("Reading body in /moves [post]")
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
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, ErrNotInGame.Error())
			return
		} else {
			slog.With("error", err).Error("Getting current game in /moves [post]")
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	g, err := db.GetGame(rg.ID)
	if err != nil {
		slog.With("error", err).Error("Getting game from DB in /moves [post]")
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
		slog.With("error", err).Error("Getting legal moves in /moves [post]")
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
		slog.With("error", err).Error("Creating turn in /moves [post]")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	isEnded, winner := isGameEnded(g)
	if isEnded {
		err := endGame(g, winner)
		slog.With("error", err).Error("Ending game")
		c.JSON(http.StatusCreated, "Moves played; Game ended")
		return
	}

	// Check if we are playing against a bot

	botLevel := db.GetBotLevel(g.Player2)
	// Against a bot
	if botLevel > 0 {
		var m []types.Move
		var err error

		switch botLevel {
		case 1:
			m, err = bgweb.GetEasyMove(g)
		case 2:
			m, err = bgweb.GetMediumMove(g)
		case 3:
			m, err = bgweb.GetBestMove(g)
		default:
			slog.Error("Invalid bot level")
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		g.PlayMove(m)
		err = db.UpdateGame(g)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		t := types.Turn{
			GameId: g.ID,
			User:   g.Player2,
			Time:   time.Now(),
			Moves:  m,
		}

		_, err = db.CreateTurn(t)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		isEnded, winner := isGameEnded(g)
		if isEnded {
			err := endGame(g, winner)
			slog.With("error", err).Error("Ending game")

			botWin := db.GetBotLevel(winner) != 0

			messagesWinUser := []string{
				"Bella partita! Torna quando vuoi per una rivincita.",
				"Non è stata una partita facile, ma alla fine hai vinto. Complimenti!",
				"Bravo! Hai vinto contro un avversario molto forte.",
				"Posso direi che ti ho lasciato vincere?",
				"Okay hai vinto, ma non illuderti. Non succederà più.",
			}

			messagesWinBot := []string{
				"Non sottovalutare l'importanza di pianificare. Bravo comunque",
				"Ehm... scusa, mi è scappata la mano. Rivincita?",
				"Non è stata una partita facile, ma alla fine ho vinto. Che peccato!",
				"Ti avevo avvisato, non perdo mai!",
				"Skill issues?",
			}

			var messages []string

			if botWin {
				messages = messagesWinBot
			} else {
				messages = messagesWinUser
			}

			m := messages[rand.Intn(len(messages))]
			err = ws.SendBotMessage(userId, m)
			if err != nil {
				slog.With("error", err).Error("Sending message to player")
			}

			c.JSON(http.StatusCreated, "Moves played; Game ended")
			return
		} else {

			messages := []string{
				"Bella mossa!",
				"Giocata ben fatta!",
				"Muovo il mio pedone...",
			}

			send := rand.Intn(3)
			if send == 0 {
				m := messages[rand.Intn(len(messages))]
				err = ws.SendBotMessage(userId, m)
				if err != nil {
					slog.With("error", err).Error("Sending message to player")
				}
			}
		}

	} else {
		// We are playing against another player but current player is already inverted
		opponentID, err := getCurrentPlayer(g.CurrentPlayer, g.Player1, g.Player2)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		slog.With("opponentID", opponentID).Debug("Turn made")
		ws.TurnMade(opponentID)
	}

	c.JSON(http.StatusCreated, "Moves played")
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
			c.JSON(http.StatusNotFound, ErrNotInGame.Error())
		} else {
			slog.With("error", err).Error("Getting current game in /double [post]")
			c.JSON(http.StatusInternalServerError, err)
		}
		return
	}

	g, err := db.GetGame(rg.ID)
	if err != nil {
		slog.With("error", err).Error("Getting game from DB in /double [post]")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if g.WantToDouble {
		c.JSON(http.StatusBadRequest, ErrDoubleNotPossible.Error())
		return
	}

	currentPlayerID, err := getCurrentPlayer(g.CurrentPlayer, g.Player1, g.Player2)
	if err != nil {
		slog.With("error", err).Error("Getting current player in /double [post]")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if (g.DoubleOwner != types.GameDoubleOwnerAll &&
		g.DoubleOwner != g.CurrentPlayer) ||
		g.DoubleValue == 64 ||
		(g.DoubleOwner == types.GameDoubleOwnerAll && currentPlayerID != userId) {
		c.JSON(http.StatusBadRequest, ErrDoubleNotPossible.Error())
		return
	}

	g.WantToDouble = true
	g.DoubleOwner = invertPlayer(g.CurrentPlayer)

	err = db.UpdateGame(g)
	if err != nil {
		slog.With("error", err).Error("Updating game in /double [post]")
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
	} else {
		id, err := getOpponentID(g.CurrentPlayer, g.Player1, g.Player2)
		if err != nil {
			slog.With("error", err).Error("Getting opponent ID in /double [post]")
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		ws.WantToDouble(id)
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
			c.JSON(http.StatusNotFound, ErrNotInGame.Error())
		} else {
			slog.With("error", err).Error("Getting current game in /double [delete]")
			c.JSON(http.StatusInternalServerError, err)
		}
		return
	}

	g, err := db.GetGame(rg.ID)
	if err != nil {
		slog.With("error", err).Error("Getting game from DB in /double [delete]")
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
			c.JSON(http.StatusNotFound, ErrNotInGame.Error())
		} else if errors.Is(err, ErrInternal) {
			slog.With("error", err).Error("Accepting double")
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusBadRequest, err)
		}
		return
	}

	c.JSON(http.StatusCreated, "Double accepted")
	// Send notification to the other player that he accepts the doubleS
}

// @Summary Get last game status
// @Schemes
// @Description Get last fame status
// @Tags play
// @Accept json
// @Produce json
// @Success 200 {string} string "Status of the last game"
// @Failure 404 "No games or no status found"
// @Router /play/last/winner [get]
func GetLastGameWinner(c *gin.Context) {
	userId := c.MustGet("user_id").(int64)

	status, err := db.GetLastGameWinner(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, "No games found")
		} else {
			slog.With("error", err).Error("Getting last game")
			c.JSON(http.StatusInternalServerError, err)
		}
		return
	}

	c.JSON(http.StatusOK, status)
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
		slog.With("error", err).Error("Creaiting bot game")
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

	messages := []string{
		"Pronto per giocare? Buona fortuna!",
		"Vediamo chi è il più bravo!",
		"Preparati a perdere!",
		"Pronto per la sfida?",
		"Finalmente un avversario all'altezza!",
		"Finalmente qualcuno che pensa di potermi battere. Che coraggio!",
	}

	m := messages[rand.Intn(len(messages))]
	go func() {
		time.Sleep(1 * time.Second)
		err := ws.SendBotMessage(userId, m)
		if err != nil {
			slog.With("error", err).Error("Sending message to player")
		}
	}()

	c.JSON(http.StatusCreated, ng)
}

func acceptDouble(userId int64) error {
	rg, err := db.GetCurrentGame(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotInGame
		} else {
			return ErrInternal
		}
	}

	g, err := db.GetGame(rg.ID)
	if err != nil {
		return err
	}

	if !g.WantToDouble {
		return ErrDoubleNotPossible
	}

	g.WantToDouble = false
	g.DoubleValue = 2 * g.DoubleValue

	err = db.UpdateGame(g)
	if err != nil {
		slog.With("error", err).Error("Updating game after accepting double")
		return ErrInternal
	}

	// save turn as double
	doublingplayer_id, err := getCurrentPlayer(g.DoubleOwner, g.Player1, g.Player2)
	if err != nil {
		slog.With("error", err).Error("Getting current player in acceptDouble")
		return ErrInternal
	}

	turn := types.Turn{
		GameId: g.ID,
		User:   doublingplayer_id,
		Time:   time.Now(),
		Double: true,
	}

	_, err = db.CreateTurn(turn)
	if err != nil {
		slog.With("error", err).Error("Creating turn in acceptDouble")
		return ErrInternal
	}

	id, err := getOpponentID(g.DoubleOwner, g.Player1, g.Player2)
	if err != nil {
		slog.With("error", err).Error("Getting opponent ID in acceptDouble")
		return ErrInternal
	}
	ws.DoubleAccepted(id)

	return nil

}

// True if the game is ended. Return id of the winner
func isGameEnded(g *types.Game) (bool, int64) {
	s := func(a [25]int8) int {
		sum := 0
		for i := 0; i < 25; i++ {
			sum += int(a[i])
		}
		return sum
	}

	if s(g.P1Checkers) == 0 {
		return true, g.Player1
	}

	if s(g.P2Checkers) == 0 {
		return true, g.Player2
	}

	return false, 0
}

func endGame(g *types.Game, winnerID int64) error {
	if winnerID == g.Player1 {
		g.Status = types.GameStatusWinP1
	} else {
		g.Status = types.GameStatusWinP2
	}

	g.End = time.Now()

	err := db.UpdateGame(g)
	if err != nil {
		return err
	}

	rg := db.GameToReturnGame(g)

	if rg.GameType == types.GameTypeOnline {
		u1, err := db.GetUser(g.Player1)
		if err != nil {
			slog.With("error", err).Error("Getting user in endGame")
			return err
		}

		u2, err := db.GetUser(g.Player2)
		if err != nil {
			slog.With("error", err).Error("Getting user in endGame")
			return err
		}

		elo1, elo2 := calculateElo(u1.Elo, u2.Elo, winnerID == g.Player1)

		err = db.UpdateUserElo(g.Player1, elo1)
		if err != nil {
			slog.With("error", err).Error("Updating elo in endGame")
			return err
		}

		err = db.UpdateUserElo(g.Player2, elo2)
		if err != nil {
			slog.With("error", err).Error("Updating elo in endGame")
			return err
		}
	}

	ws.GameEnd(g.Player1)
	ws.GameEnd(g.Player2)

	return err
}
