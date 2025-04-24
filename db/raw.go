package db

import (
	"context"
	"database/sql"
	"errors"

	ksql "github.com/kovey/db-go/v3"
)

var (
	Err_Sql_Not_Delete = errors.New("sql not delete")
	Err_Sql_Not_Update = errors.New("sql not update")
	Err_Sql_Not_Insert = errors.New("sql not insert")
	Err_Sql_Not_Query  = errors.New("sql not query")
	Err_Sql_Not_Exec   = errors.New("sql not exec")
)

func InsertRawBy(ctx context.Context, conn ksql.ConnectionInterface, raw ksql.ExpressInterface) (int64, error) {
	if raw.Type() != ksql.Sql_Type_Insert {
		return 0, _errRaw(Err_Sql_Not_Insert, raw)
	}

	result, err := conn.ExecRaw(ctx, raw)
	if err != nil {
		return 0, _errRaw(err, raw)
	}

	id, err := result.LastInsertId()
	return id, _errRaw(err, raw)
}

func InsertRaw(ctx context.Context, raw ksql.ExpressInterface) (int64, error) {
	return InsertRawBy(ctx, database, raw)
}

func UpdateRawBy(ctx context.Context, conn ksql.ConnectionInterface, raw ksql.ExpressInterface) (int64, error) {
	if raw.Type() != ksql.Sql_Type_Update {
		return 0, _errRaw(Err_Sql_Not_Update, raw)
	}

	result, err := conn.ExecRaw(ctx, raw)
	if err != nil {
		return 0, _errRaw(err, raw)
	}

	id, err := result.RowsAffected()
	return id, _errRaw(err, raw)
}

func UpdateRaw(ctx context.Context, raw ksql.ExpressInterface) (int64, error) {
	return UpdateRawBy(ctx, database, raw)
}

func DeleteRawBy(ctx context.Context, conn ksql.ConnectionInterface, raw ksql.ExpressInterface) (int64, error) {
	if raw.Type() != ksql.Sql_Type_Delete {
		return 0, _errRaw(Err_Sql_Not_Delete, raw)
	}

	result, err := conn.ExecRaw(ctx, raw)
	if err != nil {
		return 0, _errRaw(err, raw)
	}

	id, err := result.RowsAffected()
	return id, _errRaw(err, raw)
}

func DeleteRaw(ctx context.Context, raw ksql.ExpressInterface) (int64, error) {
	return DeleteRawBy(ctx, database, raw)
}

func QueryRawBy[T ksql.RowInterface](ctx context.Context, conn ksql.ConnectionInterface, raw ksql.ExpressInterface, models *[]T) error {
	if raw.IsExec() {
		return _errRaw(Err_Sql_Not_Query, raw)
	}

	cc := NewContext(ctx)
	cc.RawSqlLogStart(raw)
	defer cc.SqlLogEnd()

	stmt, err := conn.PrepareRaw(cc, raw)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(cc, raw.Binds()...)
	if err != nil {
		return _errRaw(err, raw)
	}

	defer rows.Close()
	if rows.Err() != nil {
		return _errRaw(rows.Err(), raw)
	}

	var m T
	for rows.Next() {
		tmp := m.Clone()
		if err := tmp.Scan(rows, tmp); err != nil {
			return _errRaw(err, raw)
		}

		model, ok := tmp.(T)
		if !ok {
			continue
		}

		model.WithConn(conn)
		*models = append(*models, model)
	}

	return nil
}

func QueryRaw[T ksql.RowInterface](ctx context.Context, raw ksql.ExpressInterface, models *[]T) error {
	return QueryRawBy(ctx, database, raw, models)
}

func QueryRowRawBy[T ksql.RowInterface](ctx context.Context, conn ksql.ConnectionInterface, raw ksql.ExpressInterface, model T) error {
	if raw.IsExec() {
		return _errRaw(Err_Sql_Not_Query, raw)
	}

	cc := NewContext(ctx)
	cc.RawSqlLogStart(raw)
	defer cc.SqlLogEnd()

	stmt, err := conn.PrepareRaw(cc, raw)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(cc, raw.Binds()...)
	if row.Err() != nil {
		return _errRaw(row.Err(), raw)
	}

	if err := model.Scan(row, model); err != nil {
		if err == sql.ErrNoRows {
			model.WithConn(conn)
			return nil
		}

		return _errRaw(err, raw)
	}

	model.WithConn(conn)
	return nil
}

func QueryRowRaw[T ksql.RowInterface](ctx context.Context, raw ksql.ExpressInterface, model T) error {
	return QueryRowRawBy(ctx, database, raw, model)
}

func _hasRaw(ctx context.Context, conn ksql.ConnectionInterface, raw ksql.ExpressInterface) (bool, error) {
	cc := NewContext(ctx)
	cc.RawSqlLogStart(raw)
	defer cc.SqlLogEnd()

	stmt, err := conn.PrepareRaw(cc, raw)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(cc, raw.Binds()...)
	if err != nil {
		return false, _errRaw(err, raw)
	}

	defer rows.Close()
	return rows.Next(), nil
}

func HasTableBy(ctx context.Context, conn ksql.ConnectionInterface, table string) (bool, error) {
	raw := Raw("SHOW TABLES LIKE '" + table + "'")
	return _hasRaw(ctx, conn, raw)
}

func HasTable(ctx context.Context, table string) (bool, error) {
	return HasTableBy(ctx, database, table)
}

func HasColumnBy(ctx context.Context, conn ksql.ConnectionInterface, table, column string) (bool, error) {
	raw := Raw("SHOW COLUMNS FROM `" + table + "` LIKE '" + column + "'")
	return _hasRaw(ctx, conn, raw)
}

func HasColumn(ctx context.Context, table, column string) (bool, error) {
	return HasColumnBy(ctx, database, table, column)
}

func HasIndexBy(ctx context.Context, conn ksql.ConnectionInterface, table, index string) (bool, error) {
	raw := Raw("SHOW INDEX FROM `"+table+"` WHERE Key_name = ?", index)
	return _hasRaw(ctx, conn, raw)
}

func HasIndex(ctx context.Context, table, index string) (bool, error) {
	return HasIndexBy(ctx, database, table, index)
}

func ExecRaw(ctx context.Context, raw ksql.ExpressInterface) (sql.Result, error) {
	return database.ExecRaw(ctx, raw)
}

func ExecRawBy(ctx context.Context, conn ksql.ConnectionInterface, raw ksql.ExpressInterface) (sql.Result, error) {
	return conn.ExecRaw(ctx, raw)
}

func DoBy(ctx context.Context, conn ksql.ConnectionInterface, raws ...ksql.ExpressInterface) (int64, error) {
	do := NewDo()
	for _, raw := range raws {
		do.Do(raw)
	}

	return conn.Exec(ctx, do)
}

func Do(ctx context.Context, raws ...ksql.ExpressInterface) (int64, error) {
	return DoBy(ctx, database, raws...)
}
