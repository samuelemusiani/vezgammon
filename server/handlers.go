package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func checkServer(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Server running")
}
