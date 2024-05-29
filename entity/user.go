package entity

import "context"

type User struct {
	Fullname    string `json:"fullname"`
	PhoneNumber string `json:"phone_number" binding:"e164"`
	Email       string `json:"email" binding:"email"`
	Username    string `json:"username"`
	Password    string `json:"password" binding:"required"`
	Birthdate   string `json:"birthdate" time_format:"2006-01-02" time_utc:"true"`
	LastLogin   string `json:"last_login" time_format:"2006-01-02 15:04:05" time_utc:"true"`
}

type UserRepository interface {
	GetByName(c context.Context, username string) (*User, error)
	GetByPhoneNumber(c context.Context, phone string) (*User, error)
	GetByEmail(c context.Context, email string) (*User, error)
	Set(c context.Context, user User) error
	Fetch(c context.Context) ([]User, error)
}
