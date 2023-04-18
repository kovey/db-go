package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kovey/db-go/v2/config"
	ds "github.com/kovey/db-go/v2/sql"
	"github.com/kovey/debug-go/debug"
)

var (
	database *sql.DB
	dev      string
)

type Mysql[T any] struct {
	database        *sql.DB
	tx              *sql.Tx
	isInTransaction bool
}

func NewMysql[T any]() *Mysql[T] {
	return &Mysql[T]{database: database, tx: nil, isInTransaction: false}
}

func NewSharding[T any](database *sql.DB) *Mysql[T] {
	return &Mysql[T]{database: database, tx: nil, isInTransaction: false}
}

func Init(conf config.Mysql) error {
	db, err := OpenDB(conf)
	if err != nil {
		return err
	}

	dev = conf.Dev
	database = db
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

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (m *Mysql[T]) Begin() error {
	tx, err := m.database.Begin()
	if err != nil {
		return err
	}

	m.tx = tx
	m.isInTransaction = true
	return nil
}

func (m *Mysql[T]) Commit() error {

	if m.tx == nil {
		return fmt.Errorf("transaction is not open or close")
	}

	m.isInTransaction = false
	err := m.tx.Commit()
	m.tx = nil
	return err
}

func (m *Mysql[T]) RollBack() error {

	if m.tx == nil {
		return fmt.Errorf("transaction is not open or close")
	}

	m.isInTransaction = false
	err := m.tx.Rollback()
	m.tx = nil
	return err
}

func (m *Mysql[T]) InTransaction() bool {
	return m.isInTransaction
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
	if m.isInTransaction {
		return m.tx
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

func (m *Mysql[T]) Select(sel *ds.Select, modal T) ([]T, error) {
	return Select(m.getDb(), sel, modal)
}

func (m *Mysql[T]) FetchRow(table string, where map[string]any, modal T) (T, error) {
	return FetchRow(m.getDb(), table, where, modal)
}

func (m *Mysql[T]) FetchAll(table string, where map[string]any, modal T) ([]T, error) {
	return FetchAll(m.getDb(), table, where, modal)
}

func (m *Mysql[T]) FetchAllByWhere(table string, where ds.WhereInterface, modal T) ([]T, error) {
	return FetchAllByWhere(m.getDb(), table, where, modal)
}

func (m *Mysql[T]) FetchPage(table string, where map[string]any, modal T, page int, pageSize int) ([]T, error) {
	return FetchPage(m.getDb(), table, where, modal, page, pageSize)
}

func (m *Mysql[T]) FetchPageByWhere(table string, where ds.WhereInterface, modal T, page int, pageSize int) ([]T, error) {
	return FetchPageByWhere(m.getDb(), table, where, modal, page, pageSize)
}
