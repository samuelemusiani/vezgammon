package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"strings"

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

	// Read index.html into memory
	index, err := frontend.ReadFile("dist/index.html")
	if err != nil {
		log.Fatal(err)
	}

	// If no route match is probably a vue route. So we return the index.html
	// and the vue-router takes from here
	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.String(), "/api") {
			c.JSON(http.StatusNotFound, "")
			return
		}
		c.Data(http.StatusOK, "text/html", index)
	})

	router.Run(conf.Server.Bind)
}
