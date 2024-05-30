package entity

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Coupon struct {
	Id        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Discount  float64   `json:"discount"`
	ExpiredAt time.Time `json:"expired_at"`
	CreateAt  time.Time `json:"create_at"`
}

type CouponRepository interface {
	FindById(c context.Context, id uuid.UUID) (*Coupon, error)
	Add(c context.Context, coupon *Coupon) error
	GrantCoupon(c context.Context, userId uuid.UUID, couponId uuid.UUID) error
}
