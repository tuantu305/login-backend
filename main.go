package main

import (
	"fmt"
	"log"
	"os"

	"login/cache"
	"login/internal/utility"
	"login/mq"
	"login/repository"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

const (
	BACKEND_ADDRESS = "BACKEND_ADDRESS"
	BACKEND_NODE    = "BACKEND_NODE"
	PG_HOST         = "PG_HOST"
	PG_PORT         = "PG_PORT"
	PG_USER         = "PG_USER"
	PG_PASSWORD     = "PG_PASSWORD"
	PG_DBNAME       = "PG_DBNAME"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	dbConn, err := pg_conn()
	if err != nil {
		log.Fatal(err)
	}

	backendNode := os.Getenv(BACKEND_NODE)
	backendAddress := os.Getenv(BACKEND_ADDRESS)

	router := gin.Default()
	v1 := router.Group("/v1")
	messageQueue := mq.NewMockMQ()
	idGenerator := utility.NewMockIdGenerator(backendNode)
	db := repository.NewUserPgRepository(dbConn)
	cache := cache.NewInMemoryCache()

	registerHandler := newRegisterHandler(messageQueue, idGenerator)
	loginHandler := newLoginHandler(db, cache)

	v1.POST("/register", registerHandler.handle)
	v1.POST("/login", loginHandler.handle)

	err = router.Run(backendAddress)
	if err != nil {
		panic(err)
	}
}

func pg_conn() (*sqlx.DB, error) {
	pg_host := os.Getenv(PG_HOST)
	pg_port := os.Getenv(PG_PORT)
	pg_user := os.Getenv(PG_USER)
	pg_password := os.Getenv(PG_PASSWORD)
	pg_dbname := os.Getenv(PG_DBNAME)
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pg_host,
		pg_port,
		pg_user,
		pg_password,
		pg_dbname)

	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
