package main

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"vezgammon/server/config"
	"vezgammon/server/db"
	"vezgammon/server/handler"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

//go:embed dist
var frontend embed.FS

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

	router := gin.Default()

	// middleware for static files (frontend)
	router.Use(static.Serve("/", static.EmbedFolder(frontend, "dist")))
	// middleware for backend API
	protected := router.Group("/api")
	protected.Use(handler.AuthMiddleware())

	// Gruppo di rotte protette per le API
	protected.GET("/ready", checkServer)
	protected.POST("/register", handler.Register)
	protected.POST("/login", handler.Login)
	protected.POST("/logout", handler.Logout)
	protected.GET("/user", handler.GetUser)

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

	log.Println("listening on ", conf.Server.Bind)
	router.Run(conf.Server.Bind)
}
