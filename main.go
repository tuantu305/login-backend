package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

var (
	BACKEND_ADDRESS string
	BACKEND_NODE    string
)

func init() {
	BACKEND_ADDRESS = os.Getenv("BACKEND_ADDRESS")
	BACKEND_NODE = os.Getenv("BACKEND_NODE")
}

func main() {
	router := gin.Default()
	messageQueue := newMockMQ()
	idGenerator := newMockIdGenerator(BACKEND_NODE)

	registerHandler := newRegisterHandler(
		messageQueue,
		idGenerator,
	)

	router.POST("/register", registerHandler.handle)
	router.POST("/login", loginHandler)

	err := router.Run(BACKEND_ADDRESS)
	if err != nil {
		panic(err)
	}
}

func loginHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "login",
	})
}
