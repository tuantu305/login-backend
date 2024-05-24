package main

import (
	"login/repository"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	JWT_SECRET         = ""
	JWT_EXPIRY         = 1000
	JWT_REFRESH_SECRET = ""
	JWT_REFRESH_EXPIRY = 2000
)

type LoginUser struct {
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Password    string `json:"password"`
}

type JwtCustomClaims struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	jwt.RegisteredClaims
}

type LoginHandler struct {
	db repository.UserRepository
}

func (h *LoginHandler) handle(c *gin.Context) {
	loginUser := LoginUser{}
	err := c.ShouldBind(&loginUser)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid request",
		})
	}

	if loginUser.Username == "" && loginUser.Email == "" && loginUser.PhoneNumber == "" {
		c.JSON(400, gin.H{
			"message": "invalid request",
		})
	}

	var user repository.User
	if loginUser.Username != "" {
		user, err = h.db.GetUserByName(loginUser.Username)
	} else if loginUser.Email != "" {
		user, err = h.db.GetUserByEmail(loginUser.Email)
	} else if loginUser.PhoneNumber != "" {
		user, err = h.db.GetUserByPhoneNumber(loginUser.PhoneNumber)
	}

	if err != nil {
		c.JSON(404, gin.H{
			"message": "user not found",
		})
	}

	// Use bcrypt to compare the password
	// to prevent timing attack
	// but bcrypt is insanely slow
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		c.JSON(401, gin.H{
			"message": "invalid password",
		})
	}

	accessToken, err := CreateAccessToken(&user, JWT_SECRET, JWT_EXPIRY)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "internal server error",
		})
	}

	refreshToken, err := CreateAccessToken(&user, JWT_REFRESH_SECRET, JWT_REFRESH_EXPIRY)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "internal server error",
		})
	}

	c.JSON(200, gin.H{
		"message":       "login",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func newLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

func CreateAccessToken(user *repository.User, secret string, expiry int) (string, error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry))
	claims := &JwtCustomClaims{
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, err
}
