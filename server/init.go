package main

import (
	"kebabgomon/server/server/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	//"errors"
)

var config *Config
var router *gin.Engine

func initServer() {
	path := "./config.toml"

	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	err := parseConf(path)
	if err != nil {
		log.Fatal(err)
	}

	config = getConf()
	router = gin.Default()

	//init all routes
	routes.InitRoutes(router)

	router.Run(conf.Server.Bind)
}
