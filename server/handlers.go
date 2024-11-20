package main

import (
	"net/http"
	"strings"
	"vezgammon/server/config"
	"vezgammon/server/handler"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	docs "vezgammon/server/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initHandlers(conf *config.Config) (*gin.Engine, error) {
	router := gin.Default()

	// middleware for static files (frontend)
	router.Use(static.Serve("/", static.EmbedFolder(frontend, "dist")))
	// middleware for backend API
	protected := router.Group("/api")
	protected.Use(handler.AuthMiddleware())

	// Gruppo di rotte protette per le API
	protected.POST("/register", handler.Register)
	protected.POST("/login", handler.Login)
	protected.POST("/logout", handler.Logout)
	protected.GET("/session", handler.GetSession)

	playGroup := protected.Group("/play")
	playGroup.GET("/search", handler.StartPlaySearch)
	playGroup.DELETE("/search", handler.StopPlaySearch)
	playGroup.GET("/local", handler.StartGameLocalcally)
	playGroup.GET("/", handler.GetCurrentGame)
	playGroup.DELETE("/", handler.SurrendToCurrentGame)
	playGroup.GET("/moves", handler.GetPossibleMoves)
	playGroup.POST("/moves", handler.PlayMoves)
	playGroup.POST("/double", handler.WantToDouble)
	playGroup.DELETE("/double", handler.RefuseDouble)
	playGroup.PUT("/double", handler.AcceptDouble)

	// expose swagger web console
	if conf.Swagger {
		docs.SwaggerInfo.BasePath = "/api"
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	// Read index.html into memory
	index, err := frontend.ReadFile("dist/index.html")
	if err != nil {
		return nil, err
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

	return router, nil
}
