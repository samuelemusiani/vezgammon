package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"vezgammon/server/db"
	"vezgammon/server/types"

	"github.com/gin-gonic/gin"
)

// @Summary Create a new tournament
// @Description Create a new tournament
// @Tags tournament
// @Accept  json
// @Produce  json
// @Param request body types.ReturnTournament true "Tournament object"
// @Success 201 {object} types.ReturnTournament
// @Failure 400 "bad data, tournament alredy open"
// @Failure 500 "internal server error"
// @Router /play/tournament [post]
func CreateTournament(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	buff, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var rt types.ReturnTournament
	err = json.Unmarshal(buff, &rt)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	t, err := db.ReturnTournamentToTournament(rt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if t.Owner != userID {
		c.JSON(http.StatusBadRequest, "you are not the owner")
		return
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
// @Router /play/tournament/{tournament_id} [post]
func JoinTournament(c *gin.Context) {

}

// @Summary Get a tournament
// @Description Get a tournament
// @Tags tournament
// @Accept  json
// @Produce  json
// @Param tournament_id path int true "Tournament ID"
// @Success 200 {object} types.ReturnTournament
// @Failure 404 "tournament not found"
// @Router /play/tournament/{tournament_id} [get]
func GetTournament(c *gin.Context) {

}

// @Summary List all tournaments
// @Description List all tournaments you can access
// @Tags tournament
// @Accept  json
// @Produce  json
// @Success 200 {object} types.TournamentList
// @Failure 500 "internal server error"
// @Router /play/tournament/list [get]
func ListTournaments(c *gin.Context) {

}
