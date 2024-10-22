package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	//"errors"
)

func checkServer(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Server running")
}

func main() {
	router := gin.Default()
	router.GET("/api/ready", checkServer)
	router.Run("localhost:3001")
}
