package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"vezgammon/server/db"
	"vezgammon/server/types"

	"github.com/gin-gonic/gin"
)

type createTurnamentRequest struct {
	Name string `json:"name" example:"Tournament name"`
}

// @Summary Create a new tournament
// @Description Create a new tournament
// @Tags tournament
// @Accept  json
// @Produce  json
// @Param request body createTurnamentRequest true "createTurnamentRequest object"
// @Success 201 {object} types.ReturnTournament
// @Failure 400 "bad data, tournament alredy open"
// @Failure 500 "internal server error"
// @Router /tournament/create [post]
func CreateTournament(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	buff, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var ct createTurnamentRequest
	err = json.Unmarshal(buff, &ct)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	t := &types.Tournament{
		Name:  ct.Name,
		Owner: userID,
		Users: []int64{userID},
	}

	t, err = db.CreateTournament(*t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	rett, err := db.TournamentToReturnTournament(*t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, rett)
}

// @Summary Join a tournament
// @Description Join a tournament
// @Tags tournament
// @Accept  json
// @Produce  json
// @Param tournament_id path int true "Tournament ID"
// @Success 200 {object} types.ReturnTournament
// @Failure 404 "tournament not found"
// @Failure 400 "alredy in a tournament"
// @Router /tournament/{tournament_id} [post]
func JoinTournament(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	id := c.Param("tournament_id")
	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	tournament, err := db.GetTournament(id64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if slices.Contains(tournament.Users, userID) {
		c.JSON(http.StatusBadRequest, "alredy in a tournament")
		return
	}

	if len(tournament.Users) >= 4 {
		c.JSON(http.StatusBadRequest, "tournament is full")
		// start tournament
		err = tournamentMatchCreator(tournament)
		if err != nil {
			slog.With("error", err).Debug("at starting tournament")
		}
	} else {
		tournament.Users = append(tournament.Users, userID)
		err = db.UpdateTournament(tournament)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		returnTournament, err := db.TournamentToReturnTournament(*tournament)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, returnTournament)
	}
}

// @Summary Leave a tournament
// @Description Leave a tournament if the tournament is not started
// @Tags tournament
// @Accept  json
// @Produce  json
// @Param tournament_id path int true "Tournament ID"
// @Success 201 "leaved"
// @Failure 400 "not in this tournament"
// @Failure 400 "tournament alredy started"
// @Failure 400 "you are the owner"
// @Failure 404 "tournament not found"
// @Router /tournament/{tournament_id} [delete]
func LeaveTournament(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	id := c.Param("tournament_id")
	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	t, err := db.GetTournament(id64)
	if err != nil {
		c.JSON(http.StatusNotFound, "tournament not found")
		return
	}

	if t.Owner == userID {
		c.JSON(http.StatusBadRequest, "you are the owner")
		return
	}

	if !slices.Contains(t.Users, userID) {
		c.JSON(http.StatusBadRequest, "not in this tournament")
		return
	}

	var newUserList []int64
	for _, u := range t.Users {
		if u != userID {
			newUserList = append(newUserList, u)
		}
	}

	t.Users = newUserList

	err = db.UpdateTournament(t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, "leaved")
}

// @Summary Get a tournament
// @Description Get a tournament
// @Tags tournament
// @Accept  json
// @Produce  json
// @Param tournament_id path int true "Tournament ID"
// @Success 200 {object} types.ReturnTournament
// @Failure 404 "tournament not found"
// @Router /tournament/{tournament_id} [get]
func GetTournament(c *gin.Context) {
	id := c.Param("tournament_id")
	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	tournament, err := db.GetTournament(id64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	returnTournament, err := db.TournamentToReturnTournament(*tournament)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, returnTournament)
}

// @Summary List all tournaments
// @Description List all tournaments you can access
// @Tags tournament
// @Accept  json
// @Produce  json
// @Success 200 {object} types.TournamentList
// @Failure 500 "internal server error"
// @Router /tournament/list [get]
func ListTournaments(c *gin.Context) {
	tournamentlist, err := db.GetTournamentList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, tournamentlist)
}
