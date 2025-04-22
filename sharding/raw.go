package sharding

import (
	"context"
	"database/sql"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
)

func InsertRaw(key any, ctx context.Context, raw ksql.ExpressInterface) (int64, error) {
	return db.InsertRawBy(ctx, database.conn(key), raw)
}

func UpdateRaw(key any, ctx context.Context, raw ksql.ExpressInterface) (int64, error) {
	return db.UpdateRawBy(ctx, database.conn(key), raw)
}

func DeleteRaw(key any, ctx context.Context, raw ksql.ExpressInterface) (int64, error) {
	return db.DeleteRawBy(ctx, database.conn(key), raw)
}

func QueryRaw[T ksql.RowInterface](key any, ctx context.Context, raw ksql.ExpressInterface, models *[]T) error {
	if err := db.QueryRawBy(ctx, database.conn(key), raw, models); err != nil {
		return err
	}

	for _, model := range *models {
		var tmp any = model
		if t, ok := tmp.(ShardingInterface); ok {
			t.WithKey(key)
		}
	}

	return nil
}

func QueryRowRaw[T ksql.RowInterface](key any, ctx context.Context, raw ksql.ExpressInterface, model T) error {
	var tmp any = model
	if t, ok := tmp.(ShardingInterface); ok {
		t.WithKey(key)
	}

	return db.QueryRowRawBy(ctx, database.conn(key), raw, model)
}

func HasTable(ctx context.Context, table string) (bool, error) {
	if database == nil {
		return false, db.Err_Database_Not_Initialized
	}

	has := true
	err := database.Range(func(index int, conn ksql.ConnectionInterface) error {
		if ok, err := db.HasTableBy(ctx, conn, table); err != nil || !ok {
			has = ok
			return err
		}

		return nil
	})

	return has, err
}

func HasColumn(ctx context.Context, table, column string) (bool, error) {
	if database == nil {
		return false, db.Err_Database_Not_Initialized
	}

	has := true
	err := database.Range(func(index int, conn ksql.ConnectionInterface) error {
		if ok, err := db.HasColumnBy(ctx, conn, table, column); err != nil || !ok {
			has = ok
			return err
		}
		return nil
	})

	return has, err
}

func HasIndex(ctx context.Context, table, index string) (bool, error) {
	if database == nil {
		return false, db.Err_Database_Not_Initialized
	}

	has := true
	err := database.Range(func(i int, conn ksql.ConnectionInterface) error {
		if ok, err := db.HasIndexBy(ctx, conn, table, index); err != nil || !ok {
			has = ok
			return err
		}

		return nil
	})

	return has, err
}

func ExecRaw(key any, ctx context.Context, raw ksql.ExpressInterface) (sql.Result, error) {
	return db.ExecRawBy(ctx, database.conn(key), raw)
}
