package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

var (
	BACKEND_ADDRESS string
)

func init() {
	BACKEND_ADDRESS = os.Getenv("BACKEND_ADDRESS")
}

func main() {
	router := gin.Default()

	router.POST("/register", registerHandler)
	router.POST("/login", loginHandler)

	err := router.Run(BACKEND_ADDRESS)
	if err != nil {
		panic(err)
	}
}

func registerHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "register",
	})
}

func loginHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "login",
	})
}
