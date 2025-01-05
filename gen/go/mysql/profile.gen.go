package mysql

import (
	"context"
	"my-proto-plugin/proto/user"

	"github.com/gocraft/dbr/v2"
)

const TableProfile = "profile"

var profileColumns = []string{

	"UserId",

	"Birthday",

	"Email",

	"Gender",
}

func GetProfileColumns() []string {
	return profileColumns
}

func CreateProfile(ctx context.Context, sess *dbr.Session, profile *user.Profile) error {
	_, err := sess.InsertInto(TableProfile).Columns(profileColumns...).ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func GetProfile(ctx context.Context, sess *dbr.Session, id string) (*user.Profile, error) {
	var rec *user.Profile
	err := sess.Select(profileColumns).From(TableProfile).
		Where(dbr.Eq("", id)).LoadOneContext(ctx, rec)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func UpdateProfile(ctx context.Context, sess *dbr.Session, profile *user.Profile) error {
	_, err := sess.Update(TableProfile).SetMap(map[string]interface{}{

		"UserId": profile.UserId,

		"Birthday": profile.Birthday,

		"Email": profile.Email,

		"Gender": profile.Gender,
	}).Where(dbr.Eq("", profile.Id)).ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func DeleteProfile(ctx context.Context, sess *dbr.Session, id string) error {
	_, err := sess.DeleteFrom(TableProfile).Where(dbr.Eq("", id)).ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}
