package routes

import (
	"github.com/gin-gonic/gin"
	"log"

	"net/http"
	//"../database/database.go"
)

var authRoutes *gin.RouterGroup

// momentanemente qua
func checkServer(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Server running")
}

func InitRoutes(r *gin.Engine) {
	log.Printf("Init routes group")

	authRoutes = r.Group("/api/auth")

	//userRoutes = r.Group("/api/users")
	//gameRoutes = r.Group("/api/games")
	//settingRoutes = r.Group("/api/settings")
	// ...

	//testing server
	r.GET("/api/ready", checkServer)

	//api service
	authService()

	//userService
	//gameService
	//settingService
}
