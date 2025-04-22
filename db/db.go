package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	ksql "github.com/kovey/db-go/v3"
	"github.com/kovey/db-go/v3/logger"
)

var Err_Un_Support_Operate = errors.New("unsupport operate")
var Err_Not_In_Transaction = errors.New("not in transaction")
var Err_Database_Not_Initialized = errors.New("data not initialized")
var Err_Un_Support_Save_Point = errors.New("unsupport save point")

var database ksql.ConnectionInterface
var logOpen bool = false

type Config struct {
	DriverName     string
	DataSourceName string
	MaxIdleTime    time.Duration
	MaxLifeTime    time.Duration
	MaxIdleConns   int
	MaxOpenConns   int
	LogOpened      bool
	LogMax         int
}

func Database() *sql.DB {
	return database.Database()
}

func Open(conn *sql.DB, driverName string) (ksql.ConnectionInterface, error) {
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &Connection{database: conn, driverName: driverName}, nil
}

func Get() (ksql.ConnectionInterface, error) {
	if database == nil {
		return nil, Err_Database_Not_Initialized
	}

	return database.Clone(), nil
}

func InitBy(db *sql.DB, driverName string) error {
	conn, err := Open(db, driverName)
	if err != nil {
		return err
	}

	database = conn
	return nil
}

// init global connection
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
	logOpen = conf.LogOpened
	if logOpen {
		logger.Open(conf.LogMax)
	}
	return nil
}

func Close() error {
	if logOpen {
		logger.Close()
	}
	if database == nil {
		return nil
	}

	return database.Database().Close()
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

func InsertFromBy(ctx context.Context, conn ksql.ConnectionInterface, table string, columns []string, query ksql.QueryInterface) (int64, error) {
	op := NewInsert()
	op.Table(table).Columns(columns...).From(query)
	return conn.Insert(ctx, op)
}

func InsertFrom(ctx context.Context, table string, columns []string, query ksql.QueryInterface) (int64, error) {
	return InsertFromBy(ctx, database, table, columns, query)
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
	cc := NewContext(ctx)
	cc.SqlLogStart(op)
	defer cc.SqlLogEnd()

	stmt, err := conn.Prepare(cc, op)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(cc, op.Binds()...)
	if err != nil {
		return _err(err, op)
	}
	defer rows.Close()
	if rows.Err() != nil {
		return _err(rows.Err(), op)
	}

	var m T
	for rows.Next() {
		tmp := m.Clone()
		if err := tmp.Scan(rows, tmp); err != nil {
			return _err(err, op)
		}

		tmp.Sharding(op.GetSharding())
		model, ok := tmp.(T)
		if !ok {
			continue
		}

		model.WithConn(conn)
		*models = append(*models, model)
	}

	return nil
}

func Query[T ksql.RowInterface](ctx context.Context, op ksql.QueryInterface, models *[]T) error {
	return QueryBy(ctx, database, op, models)
}

func QueryRowBy[T ksql.RowInterface](ctx context.Context, conn ksql.ConnectionInterface, op ksql.QueryInterface, model T) error {
	cc := NewContext(ctx)
	cc.SqlLogStart(op)
	defer cc.SqlLogEnd()

	stmt, err := conn.Prepare(ctx, op)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, op.Binds()...)
	if row.Err() != nil {
		return _err(err, op)
	}

	if err := model.Scan(row, model); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		return _err(err, op)
	}

	model.WithConn(conn)
	model.Sharding(op.GetSharding())
	return nil
}

func QueryRow[T ksql.RowInterface](ctx context.Context, op ksql.QueryInterface, model T) error {
	if conn := model.Conn(); conn != nil {
		return QueryRowBy(ctx, conn, op, model)
	}
	return QueryRowBy(ctx, database, op, model)
}

func TransactionBy(ctx context.Context, options *sql.TxOptions, call func(ctx context.Context, db ksql.ConnectionInterface) error) ksql.TxError {
	return database.Clone().TransactionBy(ctx, options, call)
}

func Transaction(ctx context.Context, call func(ctx context.Context, db ksql.ConnectionInterface) error) ksql.TxError {
	return TransactionBy(ctx, nil, call)
}

func Find[T FindType](ctx context.Context, model ksql.ModelInterface, id T) error {
	return FindWith(ctx, database, model, id)
}

func FindWith[T FindType](ctx context.Context, conn ksql.ConnectionInterface, model ksql.ModelInterface, id T) error {
	query := NewQuery()
	query.Table(model.Table()).Columns(model.Columns()...).Where(model.PrimaryId(), "=", id)
	return QueryRowBy(ctx, conn, query, model)
}

func FindBy(ctx context.Context, model ksql.ModelInterface, call func(query ksql.QueryInterface)) error {
	return FindByWith(ctx, database, model, call)
}

func FindByWith(ctx context.Context, conn ksql.ConnectionInterface, model ksql.ModelInterface, call func(query ksql.QueryInterface)) error {
	query := NewQuery()
	query.Table(model.Table()).Columns(model.Columns()...)
	call(query)
	return QueryRowBy(ctx, conn, query, model)
}

func LockShare[T FindType](ctx context.Context, conn ksql.ConnectionInterface, model ksql.ModelInterface, id T) error {
	if conn == nil || !conn.InTransaction() {
		return Err_Not_In_Transaction
	}

	query := NewQuery()
	query.Table(model.Table()).Columns(model.Columns()...).Where(model.PrimaryId(), "=", id).For().Share()
	return QueryRowBy(ctx, conn, query, model)
}

func LockByShare(ctx context.Context, conn ksql.ConnectionInterface, model ksql.ModelInterface, call func(query ksql.QueryInterface)) error {
	if conn == nil || !conn.InTransaction() {
		return Err_Not_In_Transaction
	}

	query := NewQuery()
	query.Table(model.Table()).Columns(model.Columns()...).For().Share()
	call(query)
	return QueryRowBy(ctx, conn, query, model)
}

func Lock[T FindType](ctx context.Context, conn ksql.ConnectionInterface, model ksql.ModelInterface, id T) error {
	if conn == nil || !conn.InTransaction() {
		return Err_Not_In_Transaction
	}

	query := NewQuery()
	query.Table(model.Table()).Columns(model.Columns()...).Where(model.PrimaryId(), "=", id).For().Update()
	return QueryRowBy(ctx, conn, query, model)
}

func LockBy(ctx context.Context, conn ksql.ConnectionInterface, model ksql.ModelInterface, call func(query ksql.QueryInterface)) error {
	if conn == nil || !conn.InTransaction() {
		return Err_Not_In_Transaction
	}

	query := NewQuery()
	query.Table(model.Table()).Columns(model.Columns()...).For().Update()
	call(query)
	return QueryRowBy(ctx, conn, query, model)
}

func Table(ctx context.Context, table string, call func(table ksql.TableInterface)) error {
	ta := NewTable().Table(table).Alter()
	call(ta)
	return ta.Exec(ctx)
}

func Create(ctx context.Context, table string, call func(table ksql.TableInterface)) error {
	ta := NewTable().Table(table).Create()
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

func DropTableIfExistsBy(ctx context.Context, conn ksql.ConnectionInterface, table string) error {
	op := NewDropTable().Table(table).IfExists()
	_, err := conn.Exec(ctx, op)
	return err
}

func DropTableIfExists(ctx context.Context, table string) error {
	return DropTableIfExistsBy(ctx, database, table)
}

func ShowDDLBy(ctx context.Context, conn ksql.ConnectionInterface, table string) (string, error) {
	var tableName *string
	var ddl *string
	if err := conn.ScanRaw(ctx, Raw("SHOW CREATE TABLE "+table), &tableName, &ddl); err != nil {
		return "", err
	}

	return *ddl, nil
}

func ShowDDL(ctx context.Context, table string) (string, error) {
	return ShowDDLBy(ctx, database, table)
}

func ScanBy(ctx context.Context, conn ksql.ConnectionInterface, query ksql.QueryInterface, vals ...any) error {
	return conn.Scan(ctx, query, vals...)
}

func Scan(ctx context.Context, query ksql.QueryInterface, vals ...any) error {
	return ScanBy(ctx, database, query)
}
