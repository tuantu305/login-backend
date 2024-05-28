package main

import (
	"login/entity"
	"login/internal/utility"
	"net/http"

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
	db    entity.UserRepository
	cache entity.Cache
}

func (h *LoginHandler) handle(c *gin.Context) {
	loginUser := entity.LoginRequest{}
	err := c.ShouldBind(&loginUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{Message: "invalid request"})
		return
	}

	if loginUser.Username == "" && loginUser.Email == "" && loginUser.PhoneNumber == "" {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{Message: "invalid request"})
		return
	}
	var user entity.User

	user, err = h.cache.GetUser(c, loginUser.Username)
	if err != nil {
		if loginUser.Username != "" {
			user, err = h.db.GetByName(c, loginUser.Username)
		} else if loginUser.Email != "" {
			user, err = h.db.GetByEmail(c, loginUser.Email)
		} else if loginUser.PhoneNumber != "" {
			user, err = h.db.GetByPhoneNumber(c, loginUser.PhoneNumber)
		}

		if err != nil {
			c.JSON(http.StatusNotFound, entity.ErrorResponse{Message: "user not found"})
			return
		}

		h.cache.SetUser(c, user.Username, user)
	}

	// Use bcrypt to compare the password
	// to prevent timing attack
	// but bcrypt is insanely slow
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		c.JSON(http.StatusForbidden, entity.ErrorResponse{Message: "invalid password"})
		return
	}

	accessToken, err := utility.CreateAccessToken(&user, JWT_SECRET, JWT_EXPIRY)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Message: "internal server error"})
		return
	}

	refreshToken, err := utility.CreateAccessToken(&user, JWT_REFRESH_SECRET, JWT_REFRESH_EXPIRY)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Message: "internal server error"})
		return
	}

	resp := entity.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusAccepted, resp)
}

func newLoginHandler(
	db entity.UserRepository,
	cache entity.Cache,
) *LoginHandler {
	return &LoginHandler{
		db:    db,
		cache: cache,
	}
}
