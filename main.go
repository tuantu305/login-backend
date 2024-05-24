package main

import (
	"os"

	"login/mq"

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
	messageQueue := mq.NewMockMQ()
	idGenerator := newMockIdGenerator(BACKEND_NODE)

	registerHandler := newRegisterHandler(
		messageQueue,
		idGenerator,
	)

	loginHandler := newLoginHandler()

	router.POST("/register", registerHandler.handle)
	router.POST("/login", loginHandler.handle)

	err := router.Run(BACKEND_ADDRESS)
	if err != nil {
		panic(err)
	}
}
