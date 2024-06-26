package main

import (
	"context"
	"fmt"
	"log"
	"login/entity"
	"login/mq"
	"login/repository"
	"os"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib" // required for sqlx.Connect
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

const (
	CouponNumber = 100
	PG_HOST      = "PG_HOST"
	PG_PORT      = "PG_PORT"
	PG_USER      = "PG_USER"
	PG_PASSWORD  = "PG_PASSWORD"
	PG_DBNAME    = "PG_DBNAME"
)

var (
	discountCoupon = entity.Coupon{
		Id:        uuid.New(),
		Discount:  30,
		CreateAt:  time.Now(),
		ExpiredAt: time.Now().AddDate(0, 0, 7),
		Code:      "Wellcome30",
	}
)

type CampainVoucherLogin struct {
	registerMq mq.MessageQueueSubscriber
	topupMq    mq.MessageQueueSubscriber
	db         entity.CouponRepository
}

// reading from 2 queues, pop the common message and process it
func (c *CampainVoucherLogin) run() {
	couponCount := 0
	idChan := make(chan uuid.UUID, CouponNumber)
	filter := bloom.NewWithEstimates(1000000, 0.005)

	c.registerMq.Subscribe("register", func(msg interface{}) error {
		if couponCount >= CouponNumber {
			return nil
		}
		registerMSg, ok := msg.(entity.RegisterRequestMsg)
		if !ok {
			return nil
		}
		filter.Add([]byte(registerMSg.User.Id.String()))
		return nil
	})

	c.topupMq.Subscribe("topup", func(msg interface{}) error {
		if couponCount >= CouponNumber {
			return nil
		}
		topupMsg, ok := msg.(entity.TopupMsg)
		if !ok {
			return nil
		}

		if filter.Test([]byte(topupMsg.Topup.UserID.String())) {
			idChan <- topupMsg.Topup.UserID
			couponCount++
			if couponCount >= CouponNumber {
				close(idChan)
			}
		}

		return nil
	})

	err := c.db.Add(context.Background(), &discountCoupon)
	if err != nil {
		log.Fatal(err)
	}

	for id := range idChan {
		c.db.GrantCoupon(context.Background(), id, discountCoupon.Id)
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	pgldb, err := pg_conn()
	if err != nil {
		log.Fatal(err)
	}
	c := &CampainVoucherLogin{
		registerMq: mq.NewMockSubscriber(),
		topupMq:    mq.NewMockSubscriber(),
		db:         repository.NewCouponPgRepository(pgldb),
	}

	c.run()
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
