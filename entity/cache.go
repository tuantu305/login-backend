package entity

import "context"

type Cache interface {
	SetUser(c context.Context, key string, value User) error
	GetUser(c context.Context, key string) (User, error)
}
