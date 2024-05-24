package repository

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExist    = errors.New("user already exist")
)

type User struct {
	Fullname    string
	PhoneNumber string
	Email       string
	Username    string
	Password    string
	Birthdate   time.Time
	LastLogin   time.Time
}

type LoginRepository interface {
	GetUserByName(username string) (User, error)
	GetUserByPhoneNumber(phone string) (User, error)
	GetUserByEmail(email string) (User, error)
	SetUser(user User) error
}

type inMemoryLoginRepository struct {
	users map[string]User
}

func (m *inMemoryLoginRepository) GetUserByName(username string) (User, error) {
	user, ok := m.users[username]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

func (m *inMemoryLoginRepository) GetUserByPhoneNumber(phone string) (User, error) {
	for _, user := range m.users {
		if user.PhoneNumber == phone {
			return user, nil
		}
	}
	return User{}, ErrUserNotFound
}

func (m *inMemoryLoginRepository) GetUserByEmail(email string) (User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, ErrUserNotFound
}

func (m *inMemoryLoginRepository) SetUser(user User) error {
	if _, ok := m.users[user.Username]; ok {
		return ErrUserExist
	}
	m.users[user.Username] = user
	return nil
}

func NewInMemoryLoginRepository() LoginRepository {
	return &inMemoryLoginRepository{
		users: make(map[string]User),
	}
}
