package mysql

import (
	"context"
	"my-proto-plugin/gen/go/mysql"
	userpb "my-proto-plugin/gen/go/user"
	"my-proto-plugin/pkg/database"

	"github.com/gocraft/dbr/v2"
)

type user struct {
	Sess *dbr.Session
}

func NewUser(sess *dbr.Session) database.User {
	return &user{Sess: sess}
}

func (u *user) Create(ctx context.Context, user *userpb.User) error {
	return mysql.CreateUser(ctx, u.Sess, user)
}

func (u *user) Get(ctx context.Context, id string) (*userpb.User, error) {
	return mysql.GetUser(ctx, u.Sess, id)
}

func (u *user) Delete(ctx context.Context, id string) error {
	return mysql.DeleteUser(ctx, u.Sess, id)
}

func (u *user) Update(ctx context.Context, user *userpb.User) error {
	return mysql.UpdateUser(ctx, u.Sess, user)
}
