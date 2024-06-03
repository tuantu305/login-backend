package main

import (
	"log"
	"login/entity"
	"login/internal/utility"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	JWT_SECRET         string
	JWT_EXPIRY         int
	JWT_REFRESH_SECRET string
	JWT_REFRESH_EXPIRY int
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
		var userDb *entity.User
		if loginUser.Username != "" {
			userDb, err = h.db.GetByName(c, loginUser.Username)
		} else if loginUser.Email != "" {
			userDb, err = h.db.GetByEmail(c, loginUser.Email)
		} else if loginUser.PhoneNumber != "" {
			userDb, err = h.db.GetByPhoneNumber(c, loginUser.PhoneNumber)
		}

		user = *userDb

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

	c.JSON(http.StatusOK, resp)
}

func newLoginHandler(db entity.UserRepository, cache entity.Cache) *LoginHandler {
	var err error
	JWT_SECRET = os.Getenv("JWT_SECRET")
	JWT_EXPIRY, err = strconv.Atoi(os.Getenv("JWT_EXPIRY"))
	if err != nil {
		log.Fatal(err)
	}
	JWT_REFRESH_SECRET = os.Getenv("JWT_REFRESH")
	JWT_REFRESH_EXPIRY, err = strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRY"))
	if err != nil {
		log.Fatal(err)
	}
	return &LoginHandler{
		db:    db,
		cache: cache,
	}
}
