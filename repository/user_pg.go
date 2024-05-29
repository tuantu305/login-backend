package repository

import (
	"context"
	"login/entity"

	"github.com/jmoiron/sqlx"
)

type UserPgRepository struct {
	db *sqlx.DB
}

func NewUserPgRepository(db *sqlx.DB) *UserPgRepository {
	return &UserPgRepository{db}
}

// Fetch implements entity.UserRepository.
func (u *UserPgRepository) Fetch(c context.Context) ([]entity.User, error) {
	var users []entity.User
	if err := u.db.GetContext(c, &users, fetchUsersQuery); err != nil {
		return nil, err
	}

	return users, nil
}

// GetByEmail implements entity.UserRepository.
func (u *UserPgRepository) GetByEmail(c context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := u.db.GetContext(c, &user, getUserByEmailQuery, email); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByName implements entity.UserRepository.
func (u *UserPgRepository) GetByName(c context.Context, username string) (*entity.User, error) {
	var user entity.User
	if err := u.db.GetContext(c, &user, getUserByNameQuery, username); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByPhoneNumber implements entity.UserRepository.
func (u *UserPgRepository) GetByPhoneNumber(c context.Context, phone string) (*entity.User, error) {
	var user entity.User
	if err := u.db.GetContext(c, &user, getUserByPhoneQuery, phone); err != nil {
		return nil, err
	}

	return &user, nil
}

// Set implements entity.UserRepository.
func (u *UserPgRepository) Set(c context.Context, user entity.User) error {
	if row := u.db.QueryRowContext(
		c,
		setUserQuery,
		user.Fullname,
		user.PhoneNumber,
		user.Email,
		user.Username,
		user.Password,
		user.Birthdate,
		user.LastLogin); row.Err() != nil {
		return row.Err()
	}

	return nil
}
