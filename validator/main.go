package main

import (
	"errors"
	"login/mq"
	"login/repository"
	"time"
)

var (
	errMustNotEmpty  = errors.New("must not empty")
	errPasswordIsReq = errors.New("password is required")
	errPasswordWeak  = errors.New("password is weak")
	errBirthdate     = errors.New("invalid birthdate")
	errLastLogin     = errors.New("invalid last login")
)

type RegisterRequest struct {
	Fullname    string `json:"fullname,omitempty"`
	PhoneNumner string `json:"phone_number,omitempty"`
	Email       string `json:"email,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Birthdate   string `json:"birthdate,omitempty"`
	LastLogin   string `json:"last_login,omitempty"`
}

type RegisterResponse struct {
	Code   int    `json:"code,omitempty"`
	Status string `json:"status,omitempty"`
}

type RegisterRequestMsg struct {
	Id      string          `json:"id"`
	Request RegisterRequest `json:"register_request"`
}

type RegisterResponseMsg struct {
	Id       string           `json:"id"`
	Response RegisterResponse `json:"register_response"`
}

// Consider to use validator library
// Sanitize input to prevent SQL Injection
func validateRegisterRequest(
	req RegisterRequest,
	sanitizer Sanitizer,
	db repository.LoginRepository,
) error {
	if req.PhoneNumner == "" && req.Email == "" && req.Username == "" {
		return errMustNotEmpty
	}

	if req.Password == "" {
		return errPasswordIsReq
	}

	if !sanitizer.Verify(req.Password) {
		return errPasswordWeak
	}

	user := repository.User{}

	if req.Birthdate != "" {
		birth, err := time.Parse("2006-01-02", req.Birthdate)
		if err != nil {
			return errBirthdate
		}
		user.Birthdate = birth
	}

	if req.LastLogin != "" {
		lastLogin, err := time.Parse("2006-01-02", req.LastLogin)
		if err != nil {
			return errLastLogin
		}
		user.LastLogin = lastLogin
	}

	user.Fullname = sanitizer.Sanitize(req.Fullname)
	user.PhoneNumber = sanitizer.Sanitize(req.PhoneNumner)
	user.Email = sanitizer.Sanitize(req.Email)
	user.Username = req.Username
	user.Password = req.Password

	err := db.SetUser(user)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// Get message from queue and validate it
	// If message is valid, write to database
	// Then send response to client via message queue

	publiser := mq.NewMockPublisher()
	subscriber := mq.NewMockSubscriber()
	sanitizer := NewMockSanitizer()
	db := repository.NewInMemoryLoginRepository()

	subscriber.Subscribe("register", func(msg interface{}) error {
		registerMsg, ok := msg.(RegisterRequestMsg)
		if !ok {
			return errors.New("invalid message")
		}

		err := validateRegisterRequest(
			registerMsg.Request,
			sanitizer,
			db,
		)
		if err != nil {
			publiser.Public("register_response", RegisterResponseMsg{
				Id: registerMsg.Id,
				Response: RegisterResponse{
					Code:   500,
					Status: err.Error(),
				},
			})
			return err
		}

		publiser.Public("register_response", RegisterResponseMsg{
			Id: registerMsg.Id,
			Response: RegisterResponse{
				Code:   200,
				Status: "success",
			},
		})
		return nil
	})

	// Block forever
	<-make(chan struct{})
}
