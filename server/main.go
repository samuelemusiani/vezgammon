package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	//"errors"
)

func checkServer(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Server running")
}

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
	router.GET("/api/ready", checkServer)
	router.Run(conf.Server.Bind)
}
