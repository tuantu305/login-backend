package main

import (
	"os"

	"login/cache"
	"login/internal/utility"
	"login/mq"
	"login/repository"

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
	idGenerator := utility.NewMockIdGenerator(BACKEND_NODE)
	db := repository.NewInMemoryUserRepository()
	cache := cache.NewInMemoryCache()

	registerHandler := newRegisterHandler(messageQueue, idGenerator)
	loginHandler := newLoginHandler(db, cache)

	router.POST("/register", registerHandler.handle)
	router.POST("/login", loginHandler.handle)

	err := router.Run(BACKEND_ADDRESS)
	if err != nil {
		panic(err)
	}
}
