package main

import (
	"embed"
	"log"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

//go:embed dist
var frontend embed.FS

func main() {
	path := "./config.toml"

	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	err := parseConf(path)
	if err != nil {
		log.Fatal(err)
	}

	conf := getConf()
	router := gin.Default()

	// middleware for static files (frontend)
	router.Use(static.Serve("/", static.EmbedFolder(frontend, "dist")))

	router.GET("/api/ready", checkServer)

	router.Run(conf.Server.Bind)
}
