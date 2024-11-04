package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/

// checkserver godoc
// @Summary test if server is running
// @Schemes
// @Description do a ping to api/ready to test if server is running
// @Tags test
// @Accept json
// @Produce json
// @Success 200 {string} Server running
// @Router /ready [get]
func checkServer(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Server running")
}
