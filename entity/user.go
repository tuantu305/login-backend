package entity

import (
	"context"

	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID `json:"id,omitempty"`
	Fullname    string    `json:"fullname,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty" binding:"e164"`
	Email       string    `json:"email,omitempty" binding:"email"`
	Username    string    `json:"username,omitempty"`
	Password    string    `json:"password" binding:"required"`
	Birthdate   string    `json:"birthdate,omitempty" time_format:"2006-01-02" time_utc:"true"`
	LastLogin   string    `json:"last_login,omitempty" time_format:"2006-01-02 15:04:05" time_utc:"true"`
}

type UserRepository interface {
	GetByName(c context.Context, username string) (*User, error)
	GetByPhoneNumber(c context.Context, phone string) (*User, error)
	GetByEmail(c context.Context, email string) (*User, error)
	Set(c context.Context, user User) error
	Fetch(c context.Context) ([]User, error)
}
