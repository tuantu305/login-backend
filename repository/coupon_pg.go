package repository

import (
	"context"
	"login/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CouponPgRepository struct {
	db *sqlx.DB
}

// Add implements entity.CouponRepository.
func (cr *CouponPgRepository) Add(c context.Context, coupon *entity.Coupon) error {
	row := cr.db.QueryRowContext(c, addCouponQuery, coupon.Code, coupon.Discount, coupon.ExpiredAt, coupon.CreateAt)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

// FindById implements entity.CouponRepository.
func (cr *CouponPgRepository) FindById(c context.Context, id uuid.UUID) (*entity.Coupon, error) {
	var coupon entity.Coupon
	err := cr.db.GetContext(c, &coupon, getCouponByIdQuery, id)
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

// GrantCoupon implements entity.CouponRepository.
func (cr *CouponPgRepository) GrantCoupon(c context.Context, userId uuid.UUID, couponId uuid.UUID) error {
	_, err := cr.db.ExecContext(c, grantCouponQuery, couponId, userId)
	if err != nil {
		return err
	}

	return nil
}

func NewCouponPgRepository(db *sqlx.DB) *CouponPgRepository {
	return &CouponPgRepository{db}
}

var _ entity.CouponRepository = &CouponPgRepository{}
