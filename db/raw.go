package db

import (
	"context"
	"database/sql"

	ksql "github.com/kovey/db-go/v3"
)

func InsertRawBy(ctx context.Context, conn ksql.ConnectionInterface, raw ksql.ExpressInterface) (int64, error) {
	result, err := conn.ExecRaw(ctx, raw)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func InsertRaw(ctx context.Context, raw ksql.ExpressInterface) (int64, error) {
	return InsertRawBy(ctx, database, raw)
}

func UpdateRawBy(ctx context.Context, conn ksql.ConnectionInterface, raw ksql.ExpressInterface) (int64, error) {
	result, err := conn.ExecRaw(ctx, raw)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func UpdateRaw(ctx context.Context, raw ksql.ExpressInterface) (int64, error) {
	return UpdateRawBy(ctx, database, raw)
}

func DeleteRawBy(ctx context.Context, conn ksql.ConnectionInterface, raw ksql.ExpressInterface) (int64, error) {
	result, err := conn.ExecRaw(ctx, raw)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func DeleteRaw(ctx context.Context, raw ksql.ExpressInterface) (int64, error) {
	return DeleteRawBy(ctx, database, raw)
}

func QueryRawBy[T ksql.RowInterface](ctx context.Context, conn ksql.ConnectionInterface, raw ksql.ExpressInterface, models *[]T) error {
	stmt, err := conn.PrepareRaw(ctx, raw)
	if err != nil {
		return err
	}

	rows, err := stmt.QueryContext(ctx, raw.Binds()...)
	if err != nil {
		return _errRaw(err, raw)
	}

	if rows.Err() != nil {
		return _errRaw(rows.Err(), raw)
	}
	defer rows.Close()

	var m T
	for rows.Next() {
		tmp := m.Clone()
		if err := rows.Scan(tmp.Values()...); err != nil {
			return err
		}

		model, ok := tmp.(T)
		if !ok {
			continue
		}

		model.FromFetch()
		model.SetConn(conn)
		*models = append(*models, model)
	}

	return nil
}

func QueryRaw[T ksql.RowInterface](ctx context.Context, raw ksql.ExpressInterface, models *[]T) error {
	return QueryRawBy(ctx, database, raw, models)
}

func QueryRowRawBy[T ksql.RowInterface](ctx context.Context, conn ksql.ConnectionInterface, raw ksql.ExpressInterface, model T) error {
	stmt, err := database.PrepareRaw(ctx, raw)
	if err != nil {
		return err
	}

	row := stmt.QueryRowContext(ctx, raw.Binds()...)
	if row.Err() != nil {
		return _errRaw(err, raw)
	}

	if err := row.Scan(model.Values()...); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		return _errRaw(err, raw)
	}

	model.FromFetch()
	model.SetConn(conn)
	return nil
}

func QueryRowRaw[T ksql.RowInterface](ctx context.Context, raw ksql.ExpressInterface, model T) error {
	return QueryRowRawBy(ctx, database, raw, model)
}

func _hasRaw(ctx context.Context, conn ksql.ConnectionInterface, raw ksql.ExpressInterface) (bool, error) {
	stmt, err := conn.PrepareRaw(ctx, raw)
	if err != nil {
		return false, err
	}

	rows, err := stmt.QueryContext(ctx, raw.Binds()...)
	if err != nil {
		return false, _errRaw(err, raw)
	}
	if rows.Err() != nil {
		return false, _errRaw(rows.Err(), raw)
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
