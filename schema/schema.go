package schema

import (
	"context"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
)

func Schema(ctx context.Context, schema string, call func(schema ksql.SchemaInterface)) error {
	sc := db.NewSchema().Schema(schema).Alter()
	call(sc)
	_, err := db.Exec(ctx, sc)
	return err
}

func Create(ctx context.Context, schema string, call func(schema ksql.SchemaInterface)) error {
	sc := db.NewSchema().Schema(schema).Create()
	call(sc)
	_, err := db.Exec(ctx, sc)
	return err
}
