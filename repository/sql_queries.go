package repository

const (
	getUserByNameQuery  = `SELECT * FROM users WHERE username = $1`
	getUserByPhoneQuery = `SELECT * FROM users WHERE phone_number = $1`
	getUserByEmailQuery = `SELECT * FROM users WHERE email = $1`
	setUserQuery        = `INSERT INTO users (fullname, phone_number, email, username, password, birthdate, last_login) 
									VALUES ($1, $2, $3, $4, $5, $6, $7)
									RETURNING *`
	fetchUsersQuery = `SELECT * FROM users`

	getCouponByIdQuery = `SELECT * FROM coupons WHERE id = $1`
	grantCouponQuery   = `UPDATE users 
			SET coupon_id =  coupon_id || $1
			WHERE id = $2`
	addCouponQuery = `INSERT INTO coupons (code, discount, expired_at, create_at) 
									VALUES ($1, $2, $3, $4) 
									RETURNING *`
)
