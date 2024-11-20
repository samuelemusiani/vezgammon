package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"vezgammon/server/config"
	"vezgammon/server/db"
	"vezgammon/server/types"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type customUser struct {
	Password  string `json:"password" example:"1234"`
	Username  string `json:"username" example:"gio"`
	Firstname string `json:"firstname" example:"giorgio"`
	Lastname  string `json:"lastname" example:"rossi"`
	Mail      string `json:"mail" example:"giorossi@mail.it"`
}

// @Summary Register new user
// @Schemes
// @Description Register new user
// @Tags authentication
// @Accept json
// @Param request body customUser true "user with password"
// @Produce json
// @Success 201 {object} types.User "user created"
// @Router /register [post]
func Register(c *gin.Context) {
	buff, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.With("err", err).Error("Reading body")
		c.JSON(http.StatusInternalServerError, err)
		return
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

type loginUserType struct {
	Password string `json:"password" example:"1234"`
	Username string `json:"username" example:"gio"`
}
type loginResponseUser struct {
	ID       int64  `json:"id" example:"1"`
	Username string `json:"username" example:"gio"`
	Email    string `json:"email" example:"giorossi@mail.it"`
}

type loginResponseType struct {
	Message string            `json:"message" example:"Login successful"`
	User    loginResponseUser `json:"user"`
}

// @Summary Login
// @Schemes
// @Description Login with a user
// @Tags authentication
// @Accept json
// @Param request body loginUserType true "username and password"
// @Produce json
// @Success 200 {object} loginResponseType
// @Router /login [post]
func Login(c *gin.Context) {

	var loginUser loginUserType
	log.Println("login test")
	if err := c.BindJSON(&loginUser); err != nil {
		slog.With("err", err).Error("Bad request unmarshal")
		c.JSON(http.StatusBadRequest, loginResponseType{Message: "Bad request"})
		return
	}

	// Verifica le credenziali
	user, err := db.LoginUser(loginUser.Username, loginUser.Password)
	if err != nil {
		slog.With("err", err).Error("Login failed")
		c.JSON(http.StatusUnauthorized, loginResponseType{Message: "Invalid credentials"})
		return
	}

	// Genera un token di sessione
	slog.Debug("Generating Session Token")
	sessionToken := db.GenerateSessionToken()
	slog.With("token", sessionToken).Debug("Token")

	// Salva il token in un sistema di sessione
	err = db.SaveSessionToken(user.ID, sessionToken)
	if err != nil {
		slog.With("err", err).Error("Failed to save session")
		c.JSON(http.StatusInternalServerError, loginResponseType{Message: "Session creation failed"})
		return
	}

	// Imposta il cookie di sessione
	c.SetCookie(
		"session_token",
		sessionToken,
		3600, // scadenza in secondi (1 ora)
		"/",
		config.Get().Server.Domain,
		true,  // solo HTTPS
		false, // httpOnly
	)

	// Risposta di successo
	c.JSON(http.StatusOK, loginResponseType{
		Message: "Login successful",
		User: loginResponseUser{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Mail,
		},
	})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Non controlliamo se stiamo facendo la login || la register
		if c.Request.URL.String() == "/api/login" || c.Request.URL.String() == "/api/register" {
			c.Next()
			return
		}

		// Ottieni il token di sessione dal cookie
		sessionToken, err := c.Cookie("session_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Verifica il token di sessione
		userID, err := db.ValidateSessionToken(sessionToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			c.Abort()
			return
		}

		// Aggiunge l'ID utente al contesto per uso successivo
		c.Set("user_id", userID)

		// Vai avanti con la richiesta
		c.Next()
	}
}

// @Summary Logout
// @Schemes
// @Description Logout
// @Tags authentication
// @Accept json
// @Produce json
// @Success 200 "Logged out successfully"
// @Failure 401 "Unauthorized"
// @Failure 500 "Logout failed"
// @Router /logout [post]
func Logout(c *gin.Context) {
	// Cancella il token di sessione
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = db.Logout(sessionToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, "Logout failed")
		return
	}

	c.JSON(http.StatusOK, "Logged out successfully")
}

// Return the session of the current user logged in
// @Summary Get current auth session
// @Schemes
// @Description Get auth session
// @Tags authentication
// @Accept json
// @Produce json
// @Success 200 {object} types.User
// @Failure 500 "error"
// @Router /session [get]
func GetSession(c *gin.Context) {
	userId := c.MustGet("user_id").(int64)
	user, err := db.GetUser(userId)
	if err != nil {
		slog.With("err", err).Error("User not found")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, user)
	return
}

/*
func GetAllUsers(c *gin.Context) {
	users, err := db.GetUsers()
	if err != nil {
		slog.With("err", err).Error("Getting users")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, users)
}
*/
