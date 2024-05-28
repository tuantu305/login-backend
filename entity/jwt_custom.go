package entity

import "github.com/golang-jwt/jwt/v4"

type JwtCustomClaims struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	jwt.RegisteredClaims
}
