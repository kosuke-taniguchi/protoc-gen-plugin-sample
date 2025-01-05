package database

import (
	"context"
	"my-proto-plugin/gen/go/user"
)

type User interface {
	Create(ctx context.Context, user *user.User) error
	Get(ctx context.Context, id string) (*user.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, user *user.User) error
}
