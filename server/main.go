package main

import (
	"log"
	"log/slog"
	"os"
	"vezgammon/server/bgweb"
	"vezgammon/server/config"
	"vezgammon/server/db"
	"vezgammon/server/handler"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	path := "./config.toml"

	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	err := config.Parse(path)
	if err != nil {
		log.Fatal(err)
	}

	conf := config.Get()

	err = db.Init(conf)
	if err != nil {
		log.Fatal(err)
	}

	bgweb.Init(conf)

	router, err := handler.InitHandlers(conf)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("listening on ", conf.Server.Bind)
	router.Run(conf.Server.Bind)
}
