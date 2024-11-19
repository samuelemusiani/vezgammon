package handler

import (
	"log/slog"
	"net/http"
	"testing"
	"vezgammon/server/bgweb"
	"vezgammon/server/config"
	"vezgammon/server/db"

	"github.com/gin-gonic/gin"
)

var conf = config.Config{
	Docker:   false,
	Database: config.Database{User: "test", Password: "test", Address: "localhost:5432"},
	Server:   config.Server{Bind: "8080"},
	Bgweb:    config.Bgweb{Url: "localhost:3030/api/v1/"},
}

var router *gin.Engine

var session_token *http.Cookie

func TestMain(m *testing.M) {
	config.Set(&conf)

	slog.SetLogLoggerLevel(slog.LevelDebug)
	db.Init(config.Get())

	q := "DROP TABLE IF EXISTS games CASCADE"
	db.Conn.Exec(q)

	q = "DROP TABLE IF EXISTS users CASCADE"
	db.Conn.Exec(q)

	db.Init(config.Get()) // recreate tables

	bgweb.Init(config.Get())

	router, _ = InitHandlers(config.Get())

	m.Run()
}
