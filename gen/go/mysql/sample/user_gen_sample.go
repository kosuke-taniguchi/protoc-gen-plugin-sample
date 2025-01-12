package mysql

import (
	"context"
	"my-proto-plugin/gen/go/user"

	"github.com/gocraft/dbr/v2"
)

const TableUser = "user"

var userColumns = []string{
	"id",
	"name",
	"email",
	"password",
	"birthday",
}

func GetUserColumns() []string {
	return userColumns
}

func CreateUser(ctx context.Context, sess *dbr.Session, user *user.User) error {
	_, err := sess.InsertInto(TableUser).Columns(userColumns...).ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(ctx context.Context, sess *dbr.Session, id string) (*user.User, error) {
	var rec *user.User
	err := sess.Select(userColumns).From(TableUser).
		Where(dbr.Eq("id", id)).LoadOneContext(ctx, rec)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func UpdateUser(ctx context.Context, sess *dbr.Session, user *user.User) error {
	_, err := sess.Update(TableUser).SetMap(map[string]interface{}{
		"name":     user.Name,
		"password": user.Password,
	}).Where(dbr.Eq("id", user.Id)).ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(ctx context.Context, sess *dbr.Session, id string) error {
	_, err := sess.DeleteFrom(TableUser).Where(dbr.Eq("id", id)).ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}
