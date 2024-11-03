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

func Login(c *gin.Context) {
	// Definisci una struttura per il login
	type loginUserType struct {
		Password string `json:"password"`
		Username string `json:"username"`
	}

	var loginUser loginUserType
	log.Println("login test")
	if err := c.BindJSON(&loginUser); err != nil {
		slog.With("err", err).Error("Bad request unmarshal")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Verifica le credenziali
	user, err := db.LoginUser(loginUser.Username, loginUser.Password)
	slog.With(err).Debug("user return")
	if err != nil {
		slog.With("err", err).Error("Login failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Session creation failed"})
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
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Mail,
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

		slog.Debug("validiamo il token")
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

func Logout(c *gin.Context) {
	// Cancella il token di sessione
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err = db.Logout(sessionToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Prendo lo user che è attualmente connesso
func GetUser(c *gin.Context) {
	user_id := c.MustGet("user_id")
	slog.With("user_id", user_id)
	user, err := db.GetUser(user_id)
	if err != nil {
		slog.With("err", err).Error("User not found")
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, user)
	return
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
