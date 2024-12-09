package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
	"vezgammon/server/config"

	_ "github.com/lib/pq"
)

var Conn *sql.DB

func Init(conf *config.Config) error {
	var user, password, address string
	var err error

	if conf.Docker {
		user = "postgres"
		buff, err := os.ReadFile("/run/secrets/db-password")
		if err != nil {
			return err
		}
		password = string(buff)
		address = "db:5432"
	} else {
		user = conf.Database.User
		password = conf.Database.Password
		address = conf.Database.Address
	}

	password = strings.TrimSpace(password)

	url := fmt.Sprintf("postgres://%s:%s@%s/vezgammon?sslmode=disable", user, password, address)
	slog.With("url", url).Debug("connecting to db")

	Conn, err = sql.Open("postgres", url)
	if err != nil {
		return err
	}

	for range 20 {
		err = Conn.Ping()
		if err != nil {
			slog.With("err", err).Debug("Waiting for DB")
			time.Sleep(time.Second * 5)
		} else {
			break
		}
	}

	if Conn == nil {
		return err
	}

	// database initialization
	err = InitUser()
	if err != nil {
		slog.With("err", err).Debug("init users")
		return err
	}
	// cookie initializazion
	err = initCookie()
	if err != nil {
		slog.With("err", err).Debug("init cookie")
		return err
	}

	err = InitTournament()
	if err != nil {
		slog.With("err", err).Debug("init tournament")
		return err
	}

	// game initialization
	err = initGame()
	if err != nil {
		slog.With("err", err).Debug("init game")
		return err
	}

	return err
}
