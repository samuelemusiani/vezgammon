package routes

import (
	"github.com/gin-gonic/gin"
	//"log"
	//"net/http"
)

func handleLogin(c *gin.Context) {
	//TODO:
	//check username/email and password
	//respond: yay/nay

	c.JSON(200, gin.H{
		"message": "Login successful",
		"token":   "...",
	})

}

func authService() {
	authRoutes.POST("/login", handleLogin)

	//authRoutes.POST("/register", handleRegister)
	//authGroup.POST("/logout", handleLogout)
}
