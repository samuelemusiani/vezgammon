package handler

import (
	"database/sql"
	"fmt"
	"log/slog"
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

func TestMain(m *testing.M) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	db.Conn, _ = sql.Open(
		"postgres", fmt.Sprintf("postgres://%s:%s@%s/vezgammon?sslmode=disable", conf.Database.User, conf.Database.Password, conf.Database.Address))

	bgweb.Init(&conf)

	router, _ = InitHandlers(&conf)

	m.Run()
}
