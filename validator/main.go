package main

import (
	"context"
	"errors"
	"login/entity"
	"login/mq"
	"login/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	errMustNotEmpty  = errors.New("must not empty")
	errPasswordIsReq = errors.New("password is required")
	errPasswordWeak  = errors.New("password is weak")
	errBirthdate     = errors.New("invalid birthdate")
	errLastLogin     = errors.New("invalid last login")
)

// Consider to use validator library
// Sanitize input to prevent SQL Injection
func validateRegisterRequest(
	ctx context.Context,
	req entity.User,
	sanitizer Sanitizer,
	db entity.UserRepository,
) error {
	if req.PhoneNumber == "" && req.Email == "" && req.Username == "" {
		return errMustNotEmpty
	}

	if !sanitizer.Verify(req.Password) {
		return errPasswordWeak
	}

	req.Fullname = sanitizer.Sanitize(req.Fullname)
	req.PhoneNumber = sanitizer.Sanitize(req.PhoneNumber)
	req.Email = sanitizer.Sanitize(req.Email)

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	req.Password = string(encryptedPassword)

	err = db.Set(ctx, req)
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
	db := repository.NewInMemoryUserRepository()

	subscriber.Subscribe("register", func(msg interface{}) error {
		registerMsg, ok := msg.(entity.RegisterRequestMsg)
		if !ok {
			return errors.New("invalid message")
		}

		err := validateRegisterRequest(
			context.Background(),
			registerMsg.User,
			sanitizer,
			db,
		)
		if err != nil {
			publiser.Public("register_response", entity.RegisterResponseMsg{
				Id: registerMsg.Id,
				Response: entity.RegisterResponse{
					Code:   500,
					Status: err.Error(),
				},
			})
			return err
		}

		publiser.Public("register_response", entity.RegisterResponseMsg{
			Id: registerMsg.Id,
			Response: entity.RegisterResponse{
				Code:   200,
				Status: "success",
			},
		})
		return nil
	})

	// Block forever
	<-make(chan struct{})
}
