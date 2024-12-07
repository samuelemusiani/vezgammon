package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"

	//"time"
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
	Avatar    string `json:"avatar" example:"robot"`
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

	avatar := "https://api.dicebear.com/6.x/avataaars/svg?seed=" + tempu.Username // default avatar from username

	u := types.User{
		Username:  tempu.Username,
		Firstname: tempu.Firstname,
		Lastname:  tempu.Lastname,
		Mail:      tempu.Mail,
		Avatar:    avatar,
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
	Token    string `json:"token"`
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
		false, // solo HTTPS
		false, // httpOnly
	)

	// Risposta di successo
	c.JSON(http.StatusOK, loginResponseType{
		Message: "Login successful",
		User: loginResponseUser{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Mail,
			Token:    sessionToken,
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

// Return the statistics of the current user loggged in
// @Summary Get users' stats
// @Schemes
// @Description Get users' stats
// @Accept json
// @Produce json
// @Success 200 {object} types.Stats
// @Failure 500 "error"
// @Router /stats [get]
func GetStats(c *gin.Context) {
	user_id := c.MustGet("user_id").(int64)

	userstats, err := db.GetStats(user_id)
	if err != nil {
		slog.With("err", err).Error("User not found")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, userstats)
	return
}

// Return the statistics of the current user loggged in
// @Summary Get users' stats WITHOUT AUTH
// @Schemes
// @Description Get users' stats WITHOUT AUTH
// @Tags public
// @Accept json
// @Produce json
// @Success 200 {object} types.Stats
// @Failure 500 "error"
// @Router /player/{username} [get]
func GetPlayer(c *gin.Context) {
	username := c.Param("username")
	u, err := db.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, db.UserNotFound) {
			c.JSON(http.StatusNotFound, "User not found")
			return
		}
		slog.With("err", err).Error("Getting user")
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	c.Set("user_id", u.ID)
	GetStats(c)
}

// @Summary Return the player avatar
// @Description Return the player avatar
// @Tags public
// @Accept json
// @Produce json
// @Param username path string true "username string"
// @Success 200  string https://api.dicebear.com/9.x/adventurer/svg?seed=Maria
// @Failure 404  "User not found"
// @Failure 500 "error"
// @Router /player/{username}/avatar [get]
func GetPlayerAvatar(c *gin.Context) {
	username := c.Param("username")

	u, err := db.GetUserByUsername(username)
	if err != nil {
		if err == db.UserNotFound {
			c.JSON(http.StatusNotFound, "User not found")
			return
		}
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	c.JSON(http.StatusOK, u.Avatar)
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

// Return all the badges user acquired
// @Summary Get user's badges
// @Schemes
// @Description Get user's badges
// @Tags
// @Accept json
// @Produce json
// @Success 200 {object} types.Badge
// @Failure 500 "error"
// @Router /badge [get]
func GetBadge(c *gin.Context) {
	user_id := c.MustGet("user_id").(int64)

	badge, err := db.GetBadge(user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	slog.With("badge", badge).Debug("Badge")
	c.JSON(http.StatusOK, badge)

}

type Avatar struct {
	Avatar string `json:"avatar"`
}

// @Summary Change user avatar image
// @Schemes
// @Description Change user avatar
// @Accept json
// @Produce json
// @Success 200
// @Failure 500 "error"
// @Router /avatar [patch]
func ChangeAvatar(c *gin.Context) {
	user_id := c.MustGet("user_id").(int64)

	buff, err := io.ReadAll(c.Request.Body)
	slog.With("buff", buff).Debug("Ohh allora")
	if err != nil {
		slog.With("err", err).Error("Reading body")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var a Avatar
	err = json.Unmarshal(buff, &a)
	if err != nil {
		slog.With("err", err).Debug("Bad request unmarshal request body")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	slog.With("avatar reading", a).Debug("Avatar")
	err = db.ChangeAvatar(user_id, a.Avatar)
	if err != nil {
		slog.With("err", err).Error("Chaning avatar")
		c.JSON(http.StatusBadRequest, "Error chaning avatar")
	}

	c.JSON(http.StatusOK, "Avatar has been changed successfuly")
}

type changePasswordType struct {
	NewPass string `json:"new_pass"`
	OldPass string `json:"old_pass"`
}

// Change password
// @Summary Change password of the user
// @Schemes
// @Description Change password given the old and new pass
// @Tags authentication
// @Accept json
// @Param request body changePasswordType true "old and new password"
// @Produce json
// @Success 200
// @Failure 500 "error"
// @Router /pass [patch]
func ChangePass(c *gin.Context) {
	user_id := c.MustGet("user_id").(int64)

	user, err := db.GetUser(user_id)
	if err != nil {
		slog.With("err", err).Error("Bad request")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	var s changePasswordType
	if err := c.BindJSON(&s); err != nil {
		slog.With("err", err).Error("Bad request unmarshal")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	err = db.ChangePass(user.Username, s.NewPass, s.OldPass)
	if err != nil {
		slog.With("err", err).Debug("Changing password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password change failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password has been changed successfuly"})
	return
}
