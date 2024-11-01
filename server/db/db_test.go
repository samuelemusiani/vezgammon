package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"testing"
	"vezgammon/server/config"
)

var conf = config.Config{
	Docker:   false,
	Database: config.Database{User: "test", Password: "test", Address: "localhost:5432"},
	Server:   config.Server{Bind: "8080"},
}

func TestMain(m *testing.M) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	conn, _ = sql.Open(
		"postgres", fmt.Sprintf("postgres://%s:%s@%s/vezgammon?sslmode=disable", conf.Database.User, conf.Database.Password, conf.Database.Address))

	m.Run()
}

func TestInit(t *testing.T) {

	err := Init(&conf)
	if err != nil {
		t.Fatalf("can't connect to db")
	}
}
