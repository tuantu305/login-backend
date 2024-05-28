package main

import (
	"login/entity"
	"login/internal/utility"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	JWT_SECRET         = ""
	JWT_EXPIRY         = 1000
	JWT_REFRESH_SECRET = ""
	JWT_REFRESH_EXPIRY = 2000
)

type LoginHandler struct {
	db entity.UserRepository
}

func (h *LoginHandler) handle(c *gin.Context) {
	loginUser := entity.LoginRequest{}
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

	var user entity.User
	if loginUser.Username != "" {
		user, err = h.db.GetByName(c, loginUser.Username)
	} else if loginUser.Email != "" {
		user, err = h.db.GetByEmail(c, loginUser.Email)
	} else if loginUser.PhoneNumber != "" {
		user, err = h.db.GetByPhoneNumber(c, loginUser.PhoneNumber)
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

	accessToken, err := utility.CreateAccessToken(&user, JWT_SECRET, JWT_EXPIRY)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "internal server error",
		})
	}

	refreshToken, err := utility.CreateAccessToken(&user, JWT_REFRESH_SECRET, JWT_REFRESH_EXPIRY)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "internal server error",
		})
	}

	resp := entity.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(200, resp)
}

func newLoginHandler() *LoginHandler {
	return &LoginHandler{}
}
