package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"
	"vezgammon/server/config"

	_ "github.com/lib/pq"
)

var conn *sql.DB

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

	url := fmt.Sprintf("postgres://%s:%s@%s/vezgammon?sslmode=disable", user, password, address)
	slog.With("url", url).Debug("connecting to db")

	conn, err = sql.Open("postgres", url)
	if err != nil {
		return err
	}

	for range 20 {
		slog.Debug("enterd")

		err = conn.Ping()
		slog.With("err", err).Debug("db errrrr")
		if err != nil {
			slog.With("err", err).Debug("Waiting for DB")
			time.Sleep(time.Second * 5)
		} else {
			break
		}
	}

	if conn == nil {
		return err
	}

	// database initialization
	err = initUser()

	return err
}
