package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kovey/db-go/v2/config"
	"github.com/kovey/db-go/v2/itf"
	ds "github.com/kovey/db-go/v2/sql"
	"github.com/kovey/db-go/v2/sql/meta"
	"github.com/kovey/debug-go/debug"
)

var (
	database *sql.DB
	dbName   string
)

type Mysql[T itf.ModelInterface] struct {
	database *sql.DB
	tx       *Tx
	DbName   string
}

func NewMysql[T itf.ModelInterface]() *Mysql[T] {
	return &Mysql[T]{database: database, tx: nil, DbName: dbName}
}

func NewSharding[T itf.ModelInterface](database *sql.DB) *Mysql[T] {
	return &Mysql[T]{database: database, tx: nil, DbName: dbName}
}

func Init(conf config.Mysql) error {
	db, err := OpenDB(conf)
	if err != nil {
		return err
	}

	database = db
	dbName = conf.Dbname

	return nil
}

func OpenDB(conf config.Mysql) (*sql.DB, error) {
	debug.Info("connection to %s:%d, dbname: %s", conf.Host, conf.Port, conf.Dbname)
	db, err := sql.Open("mysql", conf.ToDSN())
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(conf.ActiveMax)
	db.SetMaxOpenConns(conf.ConnectionMax)
	db.SetConnMaxLifetime(time.Duration(conf.LifeTime) * time.Second)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (m *Mysql[T]) SetTx(tx *Tx) {
	m.tx = tx
}

func (m *Mysql[T]) Tx() *Tx {
	return m.tx
}

func (m *Mysql[T]) Transaction(f func(tx *Tx) error) error {
	if err := m.Begin(); err != nil {
		return err
	}

	if err := f(m.tx); err != nil {
		if er := m.tx.Rollback(); err != nil {
			debug.Erro("rollBack failure, error: %s", er)
		}

		return err
	}

	return m.tx.Commit()
}

func (m *Mysql[T]) TransactionCtx(ctx context.Context, f func(tx *Tx) error, opts *sql.TxOptions) error {
	if err := m.BeginCtx(ctx, opts); err != nil {
		return err
	}

	if err := f(m.tx); err != nil {
		if er := m.tx.Rollback(); err != nil {
			debug.Erro("rollBack failure, error: %s", er)
		}
		return err
	}

	return m.tx.Commit()
}

func (m *Mysql[T]) Begin() error {
	tx, err := m.database.Begin()
	if err != nil {
		return err
	}

	m.SetTx(NewTx(tx))
	return nil
}
func (m *Mysql[T]) BeginCtx(ctx context.Context, opts *sql.TxOptions) error {
	if opts == nil {
		opts = &sql.TxOptions{Isolation: sql.LevelDefault}
	}
	tx, err := m.database.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	m.SetTx(NewTx(tx))
	return nil
}

func (m *Mysql[T]) Commit() error {
	if m.tx == nil {
		return fmt.Errorf("transaction is not open or close")
	}

	err := m.tx.Commit()
	return err
}

func (m *Mysql[T]) RollBack() error {
	if m.tx == nil {
		return fmt.Errorf("transaction is not open or close")
	}

	err := m.tx.Rollback()
	return err
}

func (m *Mysql[T]) InTransaction() bool {
	return m.tx != nil && !m.tx.IsCompleted()
}

func (m *Mysql[T]) Query(query string, model T, args ...any) ([]T, error) {
	return Query(context.Background(), m.getDb(), query, model, args...)
}

func (m *Mysql[T]) QueryCtx(ctx context.Context, query string, model T, args ...any) ([]T, error) {
	return Query(ctx, m.getDb(), query, model, args...)
}

func (m *Mysql[T]) Exec(statement string) error {
	return Exec(context.Background(), m.getDb(), statement)
}

func (m *Mysql[T]) ExecCtx(ctx context.Context, statement string) error {
	return Exec(ctx, m.getDb(), statement)
}

func (m *Mysql[T]) Insert(insert *ds.Insert) (int64, error) {
	return Insert(context.Background(), m.getDb(), insert)
}

func (m *Mysql[T]) InsertCtx(ctx context.Context, insert *ds.Insert) (int64, error) {
	return Insert(ctx, m.getDb(), insert)
}

func (m *Mysql[T]) getDb() ConnInterface {
	if m.InTransaction() {
		return m.tx.tx
	}

	return m.database
}

func (m *Mysql[T]) Update(update *ds.Update) (int64, error) {
	return Update(context.Background(), m.getDb(), update)
}

func (m *Mysql[T]) UpdateCtx(ctx context.Context, update *ds.Update) (int64, error) {
	return Update(ctx, m.getDb(), update)
}

func (m *Mysql[T]) Delete(del *ds.Delete) (int64, error) {
	return Delete(context.Background(), m.getDb(), del)
}

func (m *Mysql[T]) DeleteCtx(ctx context.Context, del *ds.Delete) (int64, error) {
	return Delete(ctx, m.getDb(), del)
}

func (m *Mysql[T]) BatchInsert(batch *ds.Batch) (int64, error) {
	return BatchInsert(context.Background(), m.getDb(), batch)
}

func (m *Mysql[T]) BatchInsertCtx(ctx context.Context, batch *ds.Batch) (int64, error) {
	return BatchInsert(ctx, m.getDb(), batch)
}

func (m *Mysql[T]) Desc(desc *ds.Desc, model T) ([]T, error) {
	return Desc(context.Background(), m.getDb(), desc, model)
}

func (m *Mysql[T]) DescCtx(ctx context.Context, desc *ds.Desc, model T) ([]T, error) {
	return Desc(ctx, m.getDb(), desc, model)
}

func (m *Mysql[T]) ShowTables(show *ds.ShowTables, model T) ([]T, error) {
	return ShowTables(context.Background(), m.database, show, model)
}

func (m *Mysql[T]) ShowTablesCtx(ctx context.Context, show *ds.ShowTables, model T) ([]T, error) {
	return ShowTables(ctx, m.database, show, model)
}

func (m *Mysql[T]) Select(sel *ds.Select, model T) ([]T, error) {
	return Select(context.Background(), m.getDb(), sel, model)
}

func (m *Mysql[T]) SelectCtx(ctx context.Context, sel *ds.Select, model T) ([]T, error) {
	return Select(ctx, m.getDb(), sel, model)
}

func (m *Mysql[T]) FetchRow(table string, where meta.Where, model T) error {
	return FetchRow(context.Background(), m.getDb(), table, where, model)
}

func (m *Mysql[T]) FetchRowCtx(ctx context.Context, table string, where meta.Where, model T) error {
	return FetchRow(ctx, m.getDb(), table, where, model)
}

func (m *Mysql[T]) LockRow(table string, where meta.Where, model T) error {
	if !m.InTransaction() {
		return fmt.Errorf("transaction not open")
	}

	return LockRow(context.Background(), m.getDb(), table, where, model)
}

func (m *Mysql[T]) LockRowCtx(ctx context.Context, table string, where meta.Where, model T) error {
	if !m.InTransaction() {
		return fmt.Errorf("transaction not open")
	}

	return LockRow(ctx, m.getDb(), table, where, model)
}

func (m *Mysql[T]) FetchAll(table string, where meta.Where, model T) ([]T, error) {
	return FetchAll(context.Background(), m.getDb(), table, where, model)
}

func (m *Mysql[T]) FetchAllCtx(ctx context.Context, table string, where meta.Where, model T) ([]T, error) {
	return FetchAll(ctx, m.getDb(), table, where, model)
}

func (m *Mysql[T]) FetchBySelect(s *ds.Select, model T) ([]T, error) {
	return FetchBySelect(context.Background(), m.getDb(), s, model)
}

func (m *Mysql[T]) FetchBySelectCtx(ctx context.Context, s *ds.Select, model T) ([]T, error) {
	return FetchBySelect(ctx, m.getDb(), s, model)
}

func (m *Mysql[T]) FetchAllByWhere(table string, where ds.WhereInterface, model T) ([]T, error) {
	return FetchAllByWhere(context.Background(), m.getDb(), table, where, model)
}

func (m *Mysql[T]) FetchAllByWhereCtx(ctx context.Context, table string, where ds.WhereInterface, model T) ([]T, error) {
	return FetchAllByWhere(ctx, m.getDb(), table, where, model)
}

func (m *Mysql[T]) FetchPage(table string, where meta.Where, model T, page int, pageSize int, orders ...string) (*meta.Page[T], error) {
	return FetchPage(context.Background(), m.getDb(), table, where, model, page, pageSize, orders...)
}

func (m *Mysql[T]) FetchPageCtx(ctx context.Context, table string, where meta.Where, model T, page int, pageSize int, orders ...string) (*meta.Page[T], error) {
	return FetchPage(ctx, m.getDb(), table, where, model, page, pageSize, orders...)
}

func (m *Mysql[T]) FetchPageByWhere(table string, where ds.WhereInterface, model T, page int, pageSize int, orders ...string) (*meta.Page[T], error) {
	return FetchPageByWhere(context.Background(), m.getDb(), table, where, model, page, pageSize, orders...)
}

func (m *Mysql[T]) FetchPageByWhereCtx(ctx context.Context, table string, where ds.WhereInterface, model T, page int, pageSize int, orders ...string) (*meta.Page[T], error) {
	return FetchPageByWhere(ctx, m.getDb(), table, where, model, page, pageSize, orders...)
}

func (m *Mysql[T]) Count(table string, where ds.WhereInterface) (int64, error) {
	return Count(context.Background(), m.getDb(), table, where)
}

func (m *Mysql[T]) CountCtx(ctx context.Context, table string, where ds.WhereInterface) (int64, error) {
	return Count(ctx, m.getDb(), table, where)
}

func (m *Mysql[T]) FetchPageBySelect(sel *ds.Select, model T) (*meta.Page[T], error) {
	return FetchPageBySelect(context.Background(), m.getDb(), sel, model)
}

func (m *Mysql[T]) FetchPageBySelectCtx(ctx context.Context, sel *ds.Select, model T) (*meta.Page[T], error) {
	return FetchPageBySelect(ctx, m.getDb(), sel, model)
}
