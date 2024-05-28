package repository

import (
	"context"
	"errors"
	"login/entity"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExist    = errors.New("user already exist")
)

type inMemoryUserRepository struct {
	users map[string]entity.User
}

func (m *inMemoryUserRepository) GetByName(c context.Context, username string) (entity.User, error) {
	user, ok := m.users[username]
	if !ok {
		return entity.User{}, ErrUserNotFound
	}
	return user, nil
}

func (m *inMemoryUserRepository) GetByPhoneNumber(c context.Context, phone string) (entity.User, error) {
	for _, user := range m.users {
		if user.PhoneNumber == phone {
			return user, nil
		}
	}
	return entity.User{}, ErrUserNotFound
}

func (m *inMemoryUserRepository) GetByEmail(c context.Context, email string) (entity.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return entity.User{}, ErrUserNotFound
}

func (m *inMemoryUserRepository) Set(c context.Context, user entity.User) error {
	if _, ok := m.users[user.Username]; ok {
		return ErrUserExist
	}
	m.users[user.Username] = user
	return nil
}

func (m *inMemoryUserRepository) Fetch(c context.Context) ([]entity.User, error) {
	var users []entity.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func NewInMemoryUserRepository() entity.UserRepository {
	return &inMemoryUserRepository{
		users: make(map[string]entity.User),
	}
}
