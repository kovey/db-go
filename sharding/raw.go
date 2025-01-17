package sharding

import (
	"context"
	"database/sql"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/db"
)

func InsertRaw(key any, ctx context.Context, raw ksql.ExpressInterface) (int64, error) {
	return db.InsertRawBy(ctx, _getConn(key), raw)
}

func UpdateRaw(key any, ctx context.Context, raw ksql.ExpressInterface) (int64, error) {
	return db.UpdateRawBy(ctx, _getConn(key), raw)
}

func DeleteRaw(key any, ctx context.Context, raw ksql.ExpressInterface) (int64, error) {
	return db.DeleteRawBy(ctx, _getConn(key), raw)
}

func QueryRaw[T ksql.RowInterface](key any, ctx context.Context, raw ksql.ExpressInterface, models *[]T) error {
	if err := db.QueryRawBy(ctx, _getConn(key), raw, models); err != nil {
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

	return db.QueryRowRawBy(ctx, _getConn(key), raw, model)
}

func HasTable(ctx context.Context, table string) (bool, error) {
	if connsCount < 1 {
		return false, db.Err_Database_Not_Initialized
	}

	for index := 0; index < connsCount; index++ {
		if has, err := db.HasTableBy(ctx, baseConns[index], table); err != nil || !has {
			return has, err
		}
	}

	return true, nil
}

func HasColumn(ctx context.Context, table, column string) (bool, error) {
	if connsCount < 1 {
		return false, db.Err_Database_Not_Initialized
	}

	for index := 0; index < connsCount; index++ {
		if has, err := db.HasColumnBy(ctx, baseConns[index], table, column); err != nil || !has {
			return has, err
		}
	}

	return true, nil
}

func HasIndex(ctx context.Context, table, index string) (bool, error) {
	if connsCount < 1 {
		return false, db.Err_Database_Not_Initialized
	}

	for i := 0; i < connsCount; i++ {
		if has, err := db.HasIndexBy(ctx, baseConns[i], table, index); err != nil || !has {
			return has, err
		}
	}

	return true, nil
}

func ExecRaw(key any, ctx context.Context, raw ksql.ExpressInterface) (sql.Result, error) {
	return db.ExecRawBy(ctx, _getConn(key), raw)
}
