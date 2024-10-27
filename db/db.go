package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	ksql "github.com/kovey/db-go/v3"
)

var Err_Un_Support_Operate = errors.New("unsupport operate")

var database ksql.ConnectionInterface

type Config struct {
	DriverName     string
	DataSourceName string
	MaxIdleTime    time.Duration
	MaxLifeTime    time.Duration
	MaxIdleConns   int
	MaxOpenConns   int
}

func Database() *sql.DB {
	return database.Database()
}

func GDB() ksql.ConnectionInterface {
	return database
}

func Open(conn *sql.DB, driverName string) (ksql.ConnectionInterface, error) {
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &Connection{database: conn, driverName: driverName}, nil
}

func Init(conf Config) error {
	db, err := sql.Open(conf.DriverName, conf.DataSourceName)
	if err != nil {
		return err
	}

	db.SetConnMaxIdleTime(conf.MaxIdleTime)
	db.SetConnMaxLifetime(conf.MaxLifeTime)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetMaxOpenConns(conf.MaxIdleConns)
	conn, err := Open(db, conf.DriverName)
	if err != nil {
		return err
	}

	database = conn
	return nil
}

func InsertBy(ctx context.Context, conn ksql.ConnectionInterface, table string, data *Data) (int64, error) {
	op := NewInsert()
	op.Table(table)
	data.Range(func(key string, val any) {
		op.Add(key, val)
	})

	return conn.Insert(ctx, op)
}

func Insert(ctx context.Context, table string, data *Data) (int64, error) {
	return InsertBy(ctx, database, table, data)
}

func UpdateBy(ctx context.Context, conn ksql.ConnectionInterface, table string, data *Data, where ksql.WhereInterface) (int64, error) {
	op := NewUpdate()
	op.Table(table)
	data.Range(func(key string, val any) {
		op.Set(key, val)
	})
	op.Where(where)

	return conn.Exec(ctx, op)
}

func Update(ctx context.Context, table string, data *Data, where ksql.WhereInterface) (int64, error) {
	return UpdateBy(ctx, database, table, data, where)
}

func DeleteBy(ctx context.Context, conn ksql.ConnectionInterface, table string, where ksql.WhereInterface) (int64, error) {
	op := NewDelete()
	op.Table(table).Where(where)

	return conn.Exec(ctx, op)
}

func Delete(ctx context.Context, table string, where ksql.WhereInterface) (int64, error) {
	return DeleteBy(ctx, database, table, where)
}

func ExecBy(ctx context.Context, conn ksql.ConnectionInterface, op ksql.SqlInterface) (int64, error) {
	return conn.Exec(ctx, op)
}

func Exec(ctx context.Context, op ksql.SqlInterface) (int64, error) {
	return ExecBy(ctx, database, op)
}

func QueryBy[T ksql.RowInterface](ctx context.Context, conn ksql.ConnectionInterface, op ksql.QueryInterface, models *[]T) error {
	stmt, err := conn.Prepare(ctx, op)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, op.Binds()...)
	if err != nil {
		return _err(err, op)
	}
	if rows.Err() != nil {
		return _err(rows.Err(), op)
	}
	defer rows.Close()

	var m T
	for rows.Next() {
		tmp := m.Clone()
		if err := rows.Scan(tmp.Values()...); err != nil {
			return _err(err, op)
		}

		model, ok := tmp.(T)
		if !ok {
			continue
		}

		model.SetConn(conn)
		model.FromFetch()
		*models = append(*models, model)
	}

	return nil
}

func Query[T ksql.RowInterface](ctx context.Context, op ksql.QueryInterface, models *[]T) error {
	return QueryBy(ctx, database, op, models)
}

func QueryRowBy[T ksql.RowInterface](ctx context.Context, conn ksql.ConnectionInterface, op ksql.QueryInterface, model T) error {
	stmt, err := conn.Prepare(ctx, op)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, op.Binds()...)
	if row.Err() != nil {
		return _err(err, op)
	}

	if err := row.Scan(model.Values()...); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		return _err(err, op)
	}

	model.FromFetch()
	model.SetConn(conn)
	return nil
}

func QueryRow[T ksql.RowInterface](ctx context.Context, op ksql.QueryInterface, model T) error {
	return QueryRowBy(ctx, database, op, model)
}

func Transaction(ctx context.Context, call func(ctx context.Context, db *Connection) error) error {
	tx, err := database.Database().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	conn := &Connection{Tx: tx}
	if err := call(ctx, conn); err != nil {
		var commitErr = err
		if err := tx.Rollback(); err != nil {
			return &TxErr{CommitErr: commitErr, RollbackErr: err}
		}

		return err
	}

	return tx.Commit()
}

func Find[T FindType](ctx context.Context, model ksql.ModelInterface, id T) error {
	query := NewQuery()
	query.Table(model.Table()).Columns(model.Columns()...).Where(model.PrimaryId(), "=", id)
	return QueryRow(ctx, query, model)
}

func FindBy(ctx context.Context, model ksql.ModelInterface, call func(query ksql.QueryInterface)) error {
	query := NewQuery()
	query.Table(model.Table()).Columns(model.Columns()...)
	call(query)
	return QueryRow(ctx, query, model)
}

func Table(ctx context.Context, table string, call func(table ksql.TableInterface)) error {
	ta := NewTable().Table(table)
	call(ta)
	return ta.Exec(ctx)
}

func Schema(ctx context.Context, schema string, call func(schema ksql.SchemaInterface)) error {
	sc := NewSchema().Schema(schema)
	call(sc)
	_, err := Exec(ctx, sc)
	return err
}

func DropTableBy(ctx context.Context, conn ksql.ConnectionInterface, table string) error {
	op := NewDropTable().Table(table)
	_, err := conn.Exec(ctx, op)
	return err
}

func DropTable(ctx context.Context, table string) error {
	return DropTableBy(ctx, database, table)
}
