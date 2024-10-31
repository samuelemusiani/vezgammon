package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"vezgammon/server/db"
	"vezgammon/server/types"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	buff, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.With("err", err).Error("Reading body")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	type customUser struct {
		Password  string `json:"password"`
		Username  string `json:"username"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Mail      string `json:"mail"`
	}

	var tempu customUser
	err = json.Unmarshal(buff, &tempu)
	if err != nil {
		slog.With("err", err).Debug("Bad request unmarshal request body")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(tempu.Password), bcrypt.DefaultCost)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			c.JSON(http.StatusBadRequest, err)
		} else {
			slog.With("err", err).Error("Hashing password")
			c.JSON(http.StatusInternalServerError, err)
		}
		return
	}

	u := types.User{
		Username:  tempu.Username,
		Firstname: tempu.Firstname,
		Lastname:  tempu.Lastname,
		Mail:      tempu.Mail,
	}

	retu, err := db.CreateUser(u, string(hash))
	if err != nil {
		slog.With("err", err).Error("Creating user in db")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, retu)
}

func GetAllUsers(c *gin.Context) {
	users, err := db.GetUsers()
	if err != nil {
		slog.With("err", err).Error("Getting users")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, users)
}
