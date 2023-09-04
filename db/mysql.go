package db

import (
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
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s", conf.Username, conf.Password, conf.Host, conf.Port, conf.Dbname, conf.Charset,
	))
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
		if err := m.tx.Rollback(); err != nil {
			debug.Erro("rollBack failure, error: %s", err)
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
	return Query(m.getDb(), query, model, args...)
}

func (m *Mysql[T]) Exec(statement string) error {
	return Exec(m.getDb(), statement)
}

func (m *Mysql[T]) Insert(insert *ds.Insert) (int64, error) {
	return Insert(m.getDb(), insert)
}

func (m *Mysql[T]) getDb() ConnInterface {
	if m.InTransaction() {
		return m.tx.tx
	}

	return m.database
}

func (m *Mysql[T]) Update(update *ds.Update) (int64, error) {
	return Update(m.getDb(), update)
}

func (m *Mysql[T]) Delete(del *ds.Delete) (int64, error) {
	return Delete(m.getDb(), del)
}

func (m *Mysql[T]) BatchInsert(batch *ds.Batch) (int64, error) {
	return BatchInsert(m.getDb(), batch)
}

func (m *Mysql[T]) Desc(desc *ds.Desc, model T) ([]T, error) {
	return Desc(m.getDb(), desc, model)
}

func (m *Mysql[T]) ShowTables(show *ds.ShowTables, model T) ([]T, error) {
	return ShowTables(m.database, show, model)
}

func (m *Mysql[T]) Select(sel *ds.Select, model T) ([]T, error) {
	return Select(m.getDb(), sel, model)
}

func (m *Mysql[T]) FetchRow(table string, where meta.Where, model T) error {
	return FetchRow(m.getDb(), table, where, model)
}

func (m *Mysql[T]) LockRow(table string, where meta.Where, model T) error {
	if !m.InTransaction() {
		return fmt.Errorf("transaction not open")
	}

	return LockRow(m.getDb(), table, where, model)
}

func (m *Mysql[T]) FetchAll(table string, where meta.Where, model T) ([]T, error) {
	return FetchAll(m.getDb(), table, where, model)
}

func (m *Mysql[T]) FetchBySelect(s *ds.Select, model T) ([]T, error) {
	return FetchBySelect(m.getDb(), s, model)
}

func (m *Mysql[T]) FetchAllByWhere(table string, where ds.WhereInterface, model T) ([]T, error) {
	return FetchAllByWhere(m.getDb(), table, where, model)
}

func (m *Mysql[T]) FetchPage(table string, where meta.Where, model T, page int, pageSize int) ([]T, error) {
	return FetchPage(m.getDb(), table, where, model, page, pageSize)
}

func (m *Mysql[T]) FetchPageByWhere(table string, where ds.WhereInterface, model T, page int, pageSize int) ([]T, error) {
	return FetchPageByWhere(m.getDb(), table, where, model, page, pageSize)
}
