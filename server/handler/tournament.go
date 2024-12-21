package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"time"
	"vezgammon/server/db"
	"vezgammon/server/types"
	"vezgammon/server/ws"

	"github.com/gin-gonic/gin"
)

var ErrTournamentNotFound = "tournament not found"

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
		Name:         ct.Name,
		Owner:        userID,
		Users:        []int64{userID},
		CreationDate: time.Now(),
		Status:       types.TournamentStatusWaiting,
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

		// send message to all users
		for _, u := range tournament.Users {
			if u != userID {
				ws.TournamentNewUserEnrolled(u)
			}
		}

		c.JSON(http.StatusOK, returnTournament)
	}
}

type invite struct {
	Username string `json:"username" example:"username"`
}

type inviteList []invite

// @Summary Invite a user or a bot a tournament
// @Description Invite a user or a bot a tournament, if it is a bot it accepts the invitation automatically, same bot can be invited multiple times
// @Tags tournament
// @Accept  json
// @Produce  json
// @Param tournament_id path int true "Tournament ID"
// @Param request body inviteList true "Invite object"
// @Success 201 "invited"
// @Failure 404 "user not found"
// @Failure 400 "user alredy in the tournament"
// @Failure 400 "you are not in the owner"
// @Failure 500 "internal server error"
// @Router /tournament/{tournament_id}/invite [post]
func InviteTournament(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	id := c.Param("tournament_id")

	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	tournament, err := db.GetTournament(id64)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrTournamentNotFound)
		return
	}

	if tournament.Owner != userID {
		c.JSON(http.StatusBadRequest, "you are not the owner")
		return
	}

	if tournament.Status != types.TournamentStatusWaiting {
		c.JSON(http.StatusBadRequest, "tournament alredy started")
		return
	}

	buff, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var invites inviteList
	err = json.Unmarshal(buff, &invites)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	for _, inviteuser := range invites {
		user, err := db.GetUserByUsername(inviteuser.Username)
		if err != nil {
			c.JSON(http.StatusNotFound, "user not found")
			return
		}

		if !user.IsBot {
			c.JSON(http.StatusBadRequest, "only bots can be invited")
			return
		}

		if len(tournament.Users) >= 4 {
			c.JSON(http.StatusBadRequest, "tournament is full")
			return
		} else {
			tournament.Users = append(tournament.Users, user.ID)
		}
	}

	err = db.UpdateTournament(tournament)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// send message to all users
	for _, u := range tournament.Users {
		// if u != userID {}
		ws.TournamentNewBotEnrolled(u)
	}
}

// @Summary Delete users and bots from a tournament
// @Description Delete users and bots from a tournament
// @Tags tournament
// @Accept  json
// @Produce  json
// @Param tournament_id path int true "Tournament ID"
// @Param request body inviteList true "Delete user list object"
// @Success 201 "deleted"
// @Failure 404 "tournament not found"
// @Failure 400 "user not in the tournament"
// @Failure 400 "you are not in the owner"
// @Failure 400 "tournament alredy started"
// @Failure 500 "internal server error"
// @Router /tournament/{tournament_id}/deletebot [delete]
func TournamentDeleteUers(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	id := c.Param("tournament_id")

	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	tournament, err := db.GetTournament(id64)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrTournamentNotFound)
		return
	}

	if tournament.Owner != userID {
		c.JSON(http.StatusBadRequest, "you are not the owner")
		return
	}

	if tournament.Status != types.TournamentStatusWaiting {
		c.JSON(http.StatusBadRequest, "tournament alredy started")
		return
	}

	buff, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var userlist inviteList
	err = json.Unmarshal(buff, &userlist)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	for _, user := range userlist {
		user, err := db.GetUserByUsername(user.Username)
		if err != nil {
			c.JSON(http.StatusNotFound, "user not found")
			return
		}

		if !user.IsBot {
			c.JSON(http.StatusBadRequest, "only bots can be deleted")
			return
		}

		if !slices.Contains(tournament.Users, user.ID) {
			c.JSON(http.StatusBadRequest, "user not in the tournament")
			return
		}

		for i, u := range tournament.Users {
			if u == user.ID {
				tournament.Users = append(tournament.Users[:i], tournament.Users[i+1:]...)
				break
			}
		}
	}

	err = db.UpdateTournament(tournament)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// send message to all users
	for _, u := range tournament.Users {
		ws.TournamentBotLeft(u)
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
		c.JSON(http.StatusNotFound, ErrTournamentNotFound)
		return
	}

	if t.Status != types.TournamentStatusWaiting {
		c.JSON(http.StatusBadRequest, "tournament alredy started")
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

	// send message to all users that the user left
	for _, u := range t.Users {
		if u != userID {
			ws.TournamentUserLeft(u)
		}
	}

	c.JSON(http.StatusCreated, "leaved")
}

// @Summary Start a tournament
// @Description Start a tournament, only the owner can start it
// @Tags tournament
// @Accept  json
// @Produce  json
// @Param tournament_id path int true "Tournament ID"
// @Success 201 "tournament started"
// @Failure 400 "tournament alredy started"
// @Failure 400 "not enough players"
// @Failure 400 "you are not the owner"
// @Failure 404 "tournament not found"
// @Router /tournament/{tournament_id}/start [post]
func StartTournament(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	id := c.Param("tournament_id")
	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	tournament, err := db.GetTournament(id64)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}

	if tournament.Owner != userID {
		c.JSON(http.StatusBadRequest, "you are not the owner")
		return
	}

	if tournament.Status != types.TournamentStatusWaiting {
		c.JSON(http.StatusBadRequest, "tournament alredy started")
		return
	}

	if len(tournament.Users) == 4 {
		// start tournament
		tournament.Status = types.TournamentStatusInProgress
		err = db.UpdateTournament(tournament)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		err = tournamentMatchCreator(tournament)
		if err != nil {
			slog.With("error", err).Debug("at starting tournament")
		}

		c.JSON(http.StatusCreated, "tournament started")
	} else {
		c.JSON(http.StatusBadRequest, "not enough players")
	}
}

// @Summary Cancel a tournament
// @Description Cancel a waiting tournament, only the owner can cancel it
// @Tags tournament
// @Accept  json
// @Produce  json
// @Param tournament_id path int true "Tournament ID"
// @Success 201 "tournament canceled"
// @Failure 400 "tournament alredy started"
// @Failure 400 "you are not the owner"
// @Failure 404 "tournament not found"
// @Router /tournament/{tournament_id}/cancel [post]
func CancelTournament(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	id := c.Param("tournament_id")
	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.With("error", err).Debug("at cancel tournament 1")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	tournament, err := db.GetTournament(id64)
	if err != nil {
		slog.With("error", err).Debug("at cancel tournament 2")
		c.JSON(http.StatusNotFound, err)
	}

	if tournament.Owner != userID {
		c.JSON(http.StatusBadRequest, "you are not the owner")
		return
	}

	if tournament.Status != types.TournamentStatusWaiting {
		c.JSON(http.StatusBadRequest, "tournament alredy started")
		return
	}

	err = db.DeleteTournament(id64)
	if err != nil {
		slog.With("error", err).Debug("at cancel tournament 3")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, "tournament canceled")

	// send cancel message to all users
	for _, u := range tournament.Users {
		ws.TournamentCancelled(u)
	}
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
