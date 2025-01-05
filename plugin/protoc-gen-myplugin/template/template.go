package template

const DBTemplate = `
package mysql

import (
	"context"
	"my-proto-plugin/gen/go/{{.ProtoPackage}}"

	"github.com/gocraft/dbr/v2"
)

const Table{{.Entity}} = "{{.GoPackage}}"

var {{.GoPackage}}Columns = []string{
	{{range .Fields}}
	"{{.Name}}",
	{{end}}
}

func Get{{.Entity}}Columns() []string {
	return {{.GoPackage}}Columns
}

func Create{{.Entity}}(ctx context.Context, sess *dbr.Session, {{.GoPackage}} *{{.ProtoPackage}}.{{.Entity}}) error {
	_, err := sess.InsertInto(Table{{.Entity}}).Columns({{.GoPackage}}Columns...).ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func Get{{.Entity}}(ctx context.Context, sess *dbr.Session, id string) (*{{.ProtoPackage}}.{{.Entity}}, error) {
	var rec *{{.ProtoPackage}}.{{.Entity}}
	err := sess.Select({{.GoPackage}}Columns).From(Table{{.Entity}}).
		Where(dbr.Eq("{{.PK}}", id)).LoadOneContext(ctx, rec)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func Update{{.Entity}}(ctx context.Context, sess *dbr.Session, {{.GoPackage}} *{{.ProtoPackage}}.{{.Entity}}) error {
	_, err := sess.Update(Table{{.Entity}}).SetMap(map[string]interface{}{
		{{range .Fields}}
		"{{.Name}}": {{$.GoPackage}}.{{.Name}},
		{{end}}
	}).Where(dbr.Eq("{{.PK}}", {{.GoPackage}}.Id)).ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func Delete{{.Entity}}(ctx context.Context, sess *dbr.Session, id string) error {
	_, err := sess.DeleteFrom(Table{{.Entity}}).Where(dbr.Eq("{{.PK}}", id)).ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}
`
