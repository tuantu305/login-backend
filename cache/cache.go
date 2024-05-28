package cache

import (
	"context"
	"errors"
	"login/entity"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type inMemoryCache struct {
	cache map[string]entity.User
}

func (i *inMemoryCache) SetUser(c context.Context, key string, value entity.User) error {
	i.cache[key] = value
	return nil
}

func (i *inMemoryCache) GetUser(c context.Context, key string) (entity.User, error) {
	user, ok := i.cache[key]
	if !ok {
		return entity.User{}, ErrUserNotFound
	}
	return user, nil
}

func NewInMemoryCache() entity.Cache {
	return &inMemoryCache{
		cache: make(map[string]entity.User),
	}
}
