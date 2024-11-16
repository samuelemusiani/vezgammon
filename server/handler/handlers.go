package handler

import (
	"embed"
	"net/http"
	"strings"
	"vezgammon/server/config"
	"vezgammon/server/handler"
	"vezgammon/server/ws"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	docs "vezgammon/server/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:embed dist
var frontend embed.FS

func InitHandlers(conf *config.Config) (*gin.Engine, error) {
	router := gin.Default()

	// middleware for static files (frontend)
	router.Use(static.Serve("/", static.EmbedFolder(frontend, "dist")))
	// middleware for backend API
	protected := router.Group("/api")
	protected.Use(AuthMiddleware())

	// Gruppo di rotte protette per le API
	protected.POST("/register", Register)
	protected.POST("/login", Login)
	protected.POST("/logout", Logout)
	protected.GET("/session", GetSession)

	playGroup := protected.Group("/play")
	playGroup.GET("/search", StartPlaySearch)
	playGroup.DELETE("/search", StopPlaySearch)
	playGroup.GET("/local", StartGameLocalcally)
	playGroup.GET("/", GetCurrentGame)
	playGroup.DELETE("/", SurrendToCurrentGame)
	playGroup.GET("/moves", GetPossibleMoves)
	playGroup.POST("/moves", PlayMoves)
	playGroup.POST("/double", WantToDouble)
	playGroup.DELETE("/double", RefuseDouble)
	playGroup.PUT("/double", AcceptDouble)
	//aggiugere la route per il ws
	playGroup.GET("/ws", func(c *gin.Context) {
		ws.WSHandler(c.Writer, c.Request, c.MustGet("user_id").(int64))
	})

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
